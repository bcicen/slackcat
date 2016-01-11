package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var version = "dev-build"

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
			output("abort")
			os.Exit(0)
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

func (sc *SlackCat) processStreamQ(noop bool) {
	if !(sc.queue.isEmpty()) {
		msglines := sc.queue.flush()
		if noop {
			output(fmt.Sprintf("skipped posting of %s message lines to %s", strconv.Itoa(len(msglines)), sc.channelName))
		} else {
			sc.postMsg(msglines)
		}
	}
	time.Sleep(3 * time.Second)
	sc.processStreamQ(noop)
}

func (sc *SlackCat) postMsg(msglines []string) {
	msg := fmt.Sprintf("```%s```", strings.Join(msglines, "\n"))
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
		os.Exit(0)
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

func readIn(lines chan string, tee bool) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines <- scanner.Text()
		if tee {
			fmt.Println(scanner.Text())
		}
	}
	close(lines)
}

func writeTemp(lines chan string) string {
	tmp, err := ioutil.TempFile(os.TempDir(), "slackcat-")
	failOnError(err, "unable to create tmpfile", false)

	w := bufio.NewWriter(tmp)
	for line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp.Name()
}

func output(s string) {
	cyan := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s %s\n", cyan("slackcat"), s)
}

func failOnError(err error, msg string, appendErr bool) {
	if err != nil {
		if appendErr {
			exitErr(fmt.Errorf("%s: %s", msg, err))
		} else {
			exitErr(fmt.Errorf("%s", msg))
		}
	}
}

func exitErr(err error) {
	output(color.RedString(err.Error()))
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "slackcat"
	app.Usage = "redirect a file to slack"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to Slack continuously instead of uploading a single snippet",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to Slack. Useful for testing",
		},
		cli.BoolFlag{
			Name:  "configure",
			Usage: "Configure Slackcat via oauth",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("configure") {
			configureOA()
		}

		token := readConfig()
		fileName := c.String("filename")

		if c.String("channel") == "" {
			exitErr(fmt.Errorf("no channel provided!"))
		}

		slackcat, err := newSlackCat(token, c.String("channel"))
		failOnError(err, "Slack API Error", true)

		if len(c.Args()) > 0 {
			if c.Bool("stream") {
				output("filepath provided, ignoring stream option")
			}
			filePath := c.Args()[0]
			if fileName == "" {
				fileName = filepath.Base(filePath)
			}
			slackcat.postFile(filePath, fileName, c.Bool("noop"))
		}

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
			output("starting stream")
			go slackcat.addToStreamQ(lines)
			go slackcat.processStreamQ(c.Bool("noop"))
			go slackcat.trap()
			select {}
		} else {
			filePath := writeTemp(lines)
			defer os.Remove(filePath)
			slackcat.postFile(filePath, fileName, c.Bool("noop"))
		}
	}

	app.Run(os.Args)

}
