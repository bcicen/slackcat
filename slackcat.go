package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/slack"
)

type SlackCat struct {
	api         *slack.Slack
	opts        *slack.ChatPostMessageOpt
	queue       *StreamQ
	shutdown    chan os.Signal
	channelName string
	channelId   string
}

func newSlackCat(token, channelName string) (*SlackCat, error) {
	sc := &SlackCat{
		api:         slack.New(token),
		opts:        &slack.ChatPostMessageOpt{AsUser: true},
		queue:       newStreamQ(),
		shutdown:    make(chan os.Signal, 1),
		channelName: channelName,
	}
	err := sc.lookupSlackId()
	if err != nil {
		return nil, err
	}
	signal.Notify(sc.shutdown, os.Interrupt)
	return sc, nil
}

func (sc *SlackCat) trap() {
	sigcount := 0
	for sig := range sc.shutdown {
		if sigcount > 0 {
			exitErr(fmt.Errorf("aborted"))
		}
		output(fmt.Sprintf("got signal: %s", sig.String()))
		output("press ctrl+c again to exit immediately")
		sigcount++
		go sc.exit()
	}
}

func (sc *SlackCat) exit() {
	for {
		if sc.queue.isEmpty() {
			os.Exit(0)
		} else {
			output("flushing remaining messages to Slack...")
			time.Sleep(3 * time.Second)
		}
	}
}

//Lookup Slack id for channel, group, or im
func (sc *SlackCat) lookupSlackId() error {
	api := sc.api
	channel, err := api.FindChannelByName(sc.channelName)
	if err == nil {
		sc.channelId = channel.Id
		return nil
	}
	group, err := api.FindGroupByName(sc.channelName)
	if err == nil {
		sc.channelId = group.Id
		return nil
	}
	im, err := api.FindImByName(sc.channelName)
	if err == nil {
		sc.channelId = im.Id
		return nil
	}
	fmt.Println(err)
	return fmt.Errorf("No such channel, group, or im")
}

func (sc *SlackCat) addToStreamQ(lines chan string) {
	for line := range lines {
		sc.queue.add(line)
	}
	sc.exit()
}

//TODO: handle messages with length exceeding maximum for Slack chat
func (sc *SlackCat) processStreamQ(noop bool, plain bool) {
	if !(sc.queue.isEmpty()) {
		msglines := sc.queue.flush()
		if noop {
			output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), sc.channelName))
		} else {
			sc.postMsg(msglines, plain)
		}
	}
	time.Sleep(3 * time.Second)
	sc.processStreamQ(noop, plain)
}

func (sc *SlackCat) postMsg(msglines []string, plain bool) {
	fmtStr := "```%s```"
	if plain {
		fmtStr = "%s"
	}
	msg := fmt.Sprintf(fmtStr, strings.Join(msglines, "\n"))
	err := sc.api.ChatPostMessage(sc.channelId, msg, sc.opts)
	failOnError(err, "", true)
	output(fmt.Sprintf("posted %s message lines to %s", strconv.Itoa(len(msglines)), sc.channelName))
}

func (sc *SlackCat) postFile(filePath, fileName string, noop bool) {
	//default to timestamp for filename
	if fileName == "" {
		fileName = strconv.FormatInt(time.Now().Unix(), 10)
	}

	if noop {
		output(fmt.Sprintf("skipping upload of file %s to %s", fileName, sc.channelName))
		return
	}

	start := time.Now()
	err := sc.api.FilesUpload(&slack.FilesUploadOpt{
		Filepath: filePath,
		Filename: fileName,
		Title:    fileName,
		Channels: []string{sc.channelId},
	})
	failOnError(err, "error uploading file to Slack", true)
	duration := strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)
	output(fmt.Sprintf("file %s uploaded to %s (%ss)", fileName, sc.channelName, duration))
	os.Exit(0)
}
