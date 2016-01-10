package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
)

const (
	base_url  = "https://slack.com/oauth/authorize"
	client_id = "7065709201.17699618306"
	scope     = "channels%3Aread+groups%3Aread+im%3Aread+chat%3Awrite%3Auser+files%3Awrite%3Auser"
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
	os.Exit(0)
}
