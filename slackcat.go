package main

import (
	"bufio"
	"fmt"
	//	"github.com/bluele/slack"
	"github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"time"
)

var channel string

type Payload struct {
	token    string
	filename string
	filepath string
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("%s: %s", msg, err)
		os.Exit(1)
	}
}

func getConfigPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("unable to determine current user")
	}
	return usr.HomeDir + "/.slackcat", nil
}

func readConfig() (string, error) {
	path, err := getConfigPath()

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("missing config: %s", path)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines[0], nil
}

func readIn() *os.File {
	var line string
	var lines []string

	tmp, err := ioutil.TempFile(os.TempDir(), "multivac-")
	if err != nil {
		panic(err)
	}

	for {
		_, err := fmt.Scan(&line)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
		fmt.Println(line)
		lines = append(lines, line)
	}

	w := bufio.NewWriter(tmp)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp
}

func postToSlack(token, tmpPath, fileName, channelName string) error {
	api := slack.New(token)
	channel, err := api.FindChannelByName(channelName)
	if err != nil {
		return fmt.Errorf("Error uploading file to Slack: %s", err)
	}

	err = api.FilesUpload(&slack.FilesUploadOpt{
		Filepath: tmpPath,
		Filename: fileName,
		Title:    fileName,
		Channels: []string{channel.Id},
	})
	if err != nil {
		return fmt.Errorf("Error uploading file to Slack: %s", err)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "slackcat"
	app.Usage = "redirect a file to slack"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel to post to",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.String("channel") == "" {
			panic(fmt.Errorf("no channel provided!"))
		}

		token, err := readConfig()
		if err != nil {
			panic(err)
		}

		tmpPath := readIn()
		fileName := strconv.FormatInt(time.Now().Unix(), 10)

		err = postToSlack(token, tmpPath.Name(), fileName, c.String("channel"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Remove(tmpPath.Name())

	}

	app.Run(os.Args)

}
