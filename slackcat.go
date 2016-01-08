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
	"github.com/skratchdot/open-golang/open"
)

var version = "dev-build"

const (
	base_url  = "https://slack.com/oauth/authorize"
	client_id = "7065709201.17699618306"
	scope     = "channels%3Aread+groups%3Aread+im%3Aread+users%3Aread+chat%3Awrite%3Auser+files%3Awrite%3Auser+files%3Aread"
)

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exit(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.slackcat"
}

func readConfig() string {
	path := getConfigPath()
	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines[0]
}

func configureOA() {
	oa_url := base_url + "?scope=" + scope + "&client_id=" + client_id
	output("Creating token request for Slackcat")
	err := open.Run(oa_url)
	if err != nil {
		output("Please open the below URL in your browser to authorize SlackCat")
		output(oa_url)
	}
	//	_, err := fmt.Scanf("%s", &i)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	failOnError(err, "", true)
	fmt.Println(dir)
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
		var channelId string

		if c.Bool("configure") {
			configureOA()
		}

		token := readConfig()
		api := slack.New(token)

		if c.String("channel") == "" {
			exit(fmt.Errorf("no channel provided!"))
		}

		channel, err := api.FindChannelByName(c.String("channel"))
		if err != nil {
			group, err := api.FindGroupByName(c.String("channel"))
			failOnError(err, "Slack API error", true)
			channelId = group.Id
		} else {
			channelId = channel.Id
		}

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
