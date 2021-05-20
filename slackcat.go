package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

//Slackcat client
type Slackcat struct {
	queue       *StreamQ
	shutdown    chan os.Signal
	username    string
	iconEmoji   string
	channelID   string
	channelName string
}

func newSlackcat(username, iconEmoji, channelname string) *Slackcat {
	sc := &Slackcat{
		queue:       newStreamQ(),
		shutdown:    make(chan os.Signal, 1),
		username:    username,
		iconEmoji:   iconEmoji,
		channelName: channelname,
	}

	sc.channelID = lookupSlackID(sc.channelName)

	signal.Notify(sc.shutdown, os.Interrupt)
	return sc
}

func (sc *Slackcat) trap() {
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

func (sc *Slackcat) exit() {
	for {
		if sc.queue.IsEmpty() {
			os.Exit(0)
		} else {
			output("flushing remaining messages to Slack...")
			time.Sleep(3 * time.Second)
		}
	}
}

func (sc *Slackcat) stream(lines chan string) {
	output("starting stream")

	go func() {
		for line := range lines {
			sc.queue.Add(line)
		}
		sc.exit()
	}()

	go sc.processStreamQ()
	go sc.trap()
	select {}
}

func (sc *Slackcat) processStreamQ() {
	if !(sc.queue.IsEmpty()) {
		msglines := sc.queue.Flush()
		if noop {
			output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), sc.channelName))
		} else {
			sc.postMsg(msglines)
		}
		sc.queue.Ack()
	}
	time.Sleep(3 * time.Second)
	sc.processStreamQ()
}

var CurMsgTS string

func (sc *Slackcat) postMsg(msglines []string) {
	msg := strings.Join(msglines, "\n")
	msg = strings.Replace(msg, "&", "%26amp%3B", -1)
	msg = strings.Replace(msg, "<", "%26lt%3B", -1)
	msg = strings.Replace(msg, ">", "%26gt%3B", -1)

	msgOpts := []slack.MsgOption{slack.MsgOptionText(msg, false)}
	if sc.username != "" {
		msgOpts = append(msgOpts, slack.MsgOptionAsUser(false))
		msgOpts = append(msgOpts, slack.MsgOptionUsername(sc.username))
	} else {
		msgOpts = append(msgOpts, slack.MsgOptionAsUser(true))
	}
	if sc.iconEmoji != "" {
		msgOpts = append(msgOpts, slack.MsgOptionIconEmoji(sc.iconEmoji))
	}

	if thread {
		if CurMsgTS != "" {
			msgOpts = append(msgOpts, slack.MsgOptionTS(CurMsgTS))
		}
	}

	var err error

	_, CurMsgTS, err = api.PostMessage(sc.channelID, msgOpts...)

	output(CurMsgTS)

	failOnError(err)
	count := strconv.Itoa(len(msglines))
	output(fmt.Sprintf("posted %s message lines to %s", count, sc.channelName))
}

func (sc *Slackcat) postFile(filePath, fileName, fileType, fileComment string) {
	//default to timestamp for filename
	if fileName == "" {
		fileName = strconv.FormatInt(time.Now().Unix(), 10)
	}

	if noop {
		output(fmt.Sprintf("skipping upload of file %s to %s", fileName, sc.channelName))
		return
	}

	start := time.Now()
	_, err := api.UploadFile(slack.FileUploadParameters{
		File:           filePath,
		Filename:       fileName,
		Filetype:       fileType,
		Title:          fileName,
		InitialComment: fileComment,
		Channels:       []string{sc.channelID},
	})
	failOnError(err, "error uploading file to Slack")
	duration := strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)
	output(fmt.Sprintf("file %s uploaded to %s (%ss)", fileName, sc.channelName, duration))
	os.Exit(0)
}
