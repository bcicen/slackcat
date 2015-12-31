package main

import (
	"bufio"
	"fmt"
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

func getConfigPath() string {
	usr, err := user.Current()
	failOnError(err, "unable to determine current user", false)

	return usr.HomeDir + "/.slackcat"
}

func readConfig() string {
	path := getConfigPath()
	file, err := os.Open(path)
	failOnError(err, "missing config: "+path, false)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines[0]
}

func readIn(tee bool) *os.File {
	var line string
	var lines []string

	tmp, err := ioutil.TempFile(os.TempDir(), "multivac-")
	failOnError(err, "failed to create tempfile", false)

	for {
		_, err := fmt.Scan(&line)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
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

	return tmp
}

func postToSlack(token, path, name, channelName string) error {
	defer os.Remove(path)

	api := slack.New(token)
	channel, err := api.FindChannelByName(channelName)
	if err != nil {
		return err
	}

	err = api.FilesUpload(&slack.FilesUploadOpt{
		Filepath: path,
		Filename: name,
		Title:    name,
		Channels: []string{channel.Id},
	})
	if err != nil {
		return err
	}

	fmt.Printf("file %s uploaded to %s\n", name, channelName)
	return nil
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
	fmt.Println(err)
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "slackcat"
	app.Usage = "redirect a file to slack"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel to post to",
		},
	}

	app.Action = func(c *cli.Context) {
		token := readConfig()

		if c.String("channel") == "" {
			exit(fmt.Errorf("no channel provided!"))
		}

		tmpPath := readIn(c.Bool("tee"))
		fileName := strconv.FormatInt(time.Now().Unix(), 10)

		err := postToSlack(token, tmpPath.Name(), fileName, c.String("channel"))
		failOnError(err, "error uploading file to Slack", true)
	}

	app.Run(os.Args)

}
