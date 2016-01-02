package main

import (
	"bufio"
	"fmt"
	"github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

func getConfigPath() string {
	usr, err := user.Current()
	failOnError(err, "unable to determine current user", false)
	return usr.HomeDir + "/.slackcat"
}

func readConfig() string {
	path := getConfigPath()
	file, err := os.Open(path)
	failOnError(err, "unable to read config: "+path, false)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines[0]
}

func readIn(tee bool) string {
	var line string
	var lines []string

	tmp, err := ioutil.TempFile(os.TempDir(), "slackcat-")
	failOnError(err, "unable to create tmpfile", false)

	for {
		_, err := fmt.Scan(&line)
		if err != nil {
			if err != io.EOF {
				exit(err)
			}
			break
		}
		if tee {
			fmt.Println(line)
		}
		lines = append(lines, line)
	}

	w := bufio.NewWriter(tmp)
	for _, line := range lines {
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
			exit(fmt.Errorf("%s:\n%s", msg, err))
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
	app.Version = "0.2"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to Slack. Useful for testing",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel to post to",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
	}

	app.Action = func(c *cli.Context) {
		var filePath string
		var fileName string

		token := readConfig()
		api := slack.New(token)

		if c.String("channel") == "" {
			exit(fmt.Errorf("no channel provided!"))
		}

		channel, err := api.FindChannelByName(c.String("channel"))
		failOnError(err, "Slack API error", true)

		if len(c.Args()) > 0 {
			filePath = c.Args()[0]
			fileName = filepath.Base(filePath)
		} else {
			filePath = readIn(c.Bool("tee"))
			fileName = strconv.FormatInt(time.Now().Unix(), 10)
			defer os.Remove(filePath)
		}

		//override default filename with provided option value
		if c.String("filename") != "" {
			fileName = c.String("filename")
		}

		if c.Bool("noop") {
			output(fmt.Sprintf("skipping upload of file %s to %s", fileName, channel.Name))
		} else {
			err = api.FilesUpload(&slack.FilesUploadOpt{
				Filepath: filePath,
				Filename: fileName,
				Title:    fileName,
				Channels: []string{channel.Id},
			})
			failOnError(err, "error uploading file to Slack", true)
			output(fmt.Sprintf("file %s uploaded to %s", fileName, channel.Name))
		}
	}

	app.Run(os.Args)

}
