package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/skratchdot/open-golang/open"
)

const (
	base_url  = "https://slack.com/oauth/authorize"
	client_id = "7065709201.17699618306"
	scope     = "channels%3Aread+groups%3Aread+im%3Aread+users%3Aread+chat%3Awrite%3Auser+files%3Awrite%3Auser"
)

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exitErr(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.slackcat"
}

func getToken(teamName string) string {
	token := readConfig()[teamName]
	if token == "" {
		exitErr(fmt.Errorf("no such team configured: %s", teamName))
	}
	return token
}

func readConfig() map[string]string {
	var line string

	path := getConfigPath()
	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	teams := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, ":") {
			s := strings.Split(line, ":")
			teams[s[0]] = strings.Replace(s[1], " ", "", -1)
		} else {
			teams["default"] = line
		}
	}
	return teams
}

func configureOA() {
	oa_url := base_url + "?scope=" + scope + "&client_id=" + client_id
	output("Creating token request for Slackcat")
	err := open.Run(oa_url)
	if err != nil {
		output("Please open the below URL in your browser to authorize SlackCat")
		output(oa_url)
	}
}
