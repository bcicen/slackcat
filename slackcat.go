package main

import (
	"bufio"
	"fmt"
	//	"github.com/bluele/slack"
	"github.com/codegangsta/cli"
	"os"
	"os/user"
)

var channel string

func readConfig(path string) (string, error) {
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

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("unable to determine current user")
		os.Exit(1)
	}

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
			fmt.Println("no channel provided!")
			os.Exit(1)
		}
		token, err := readConfig(usr.HomeDir + "/.slackcat")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(token)
	}

	app.Run(os.Args)

}
