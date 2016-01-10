package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

var version = "dev-build"

//Lookup Slack id for channel, group, or im
func lookupSlackId(api *slack.Slack, name string) (string, error) {
	channel, err := api.FindChannelByName(name)
	if err == nil {
		return channel.Id, nil
	}
	group, err := api.FindGroupByName(name)
	if err == nil {
		return group.Id, nil
	}
	im, err := api.FindImByName(name)
	if err == nil {
		return im.Id, nil
	}
	return "", fmt.Errorf("No such channel, group, or im")
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
			exit(fmt.Errorf("%s: %s", msg, err))
		} else {
			exit(fmt.Errorf("%s", msg))
		}
	}
}

func exit(err error) {
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
		var filePath string
		var fileName string

		if c.Bool("configure") {
			configureOA()
		}

		token := readConfig()
		api := slack.New(token)

		if c.String("channel") == "" {
			exit(fmt.Errorf("no channel provided!"))
		}

		channelId, err := lookupSlackId(api, c.String("channel"))
		failOnError(err, "", true)

		if len(c.Args()) > 0 {
			filePath = c.Args()[0]
			fileName = filepath.Base(filePath)
		} else {
			lines := make(chan string)
			go readIn(lines, c.Bool("tee"))
			filePath = writeTemp(lines)
			fileName = strconv.FormatInt(time.Now().Unix(), 10)
			defer os.Remove(filePath)
		}

		//override default filename with provided option value
		if c.String("filename") != "" {
			fileName = c.String("filename")
		}

		if c.Bool("noop") {
			output(fmt.Sprintf("skipping upload of file %s to %s", fileName, c.String("channel")))
		} else {
			start := time.Now()
			err = api.FilesUpload(&slack.FilesUploadOpt{
				Filepath: filePath,
				Filename: fileName,
				Title:    fileName,
				Channels: []string{channelId},
			})
			failOnError(err, "error uploading file to Slack", true)
			duration := strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64)
			output(fmt.Sprintf("file %s uploaded to %s (%ss)", fileName, c.String("channel"), duration))
		}
	}

	app.Run(os.Args)

}
