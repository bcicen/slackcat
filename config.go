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

type Config struct {
	teams          map[string]string
	defaultTeam    string
	defaultChannel string
}

func getConfigPath() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		exitErr(fmt.Errorf("$HOME not set"))
	}
	return homedir + "/.slackcat"
}

func (c *Config) parseChannelOpt(channel string) (string, string, error) {
	//use default channel if none provided
	if channel == "" {
		if c.defaultChannel == "" {
			return "", "", fmt.Errorf("no channel provided!")
		} else {
			return c.defaultTeam, c.defaultChannel, nil
		}
	}
	//if channel is prefixed with a team
	if strings.Contains(channel, ":") {
		s := strings.Split(channel, ":")
		return s[0], s[1], nil
	}
	//use default team with provided channel
	return c.defaultTeam, channel, nil
}

func readConfig() *Config {
	config := &Config{
		teams:          make(map[string]string),
		defaultTeam:    "",
		defaultChannel: "",
	}
	lines := readLines(getConfigPath())

	//simple config file
	if len(lines) == 1 {
		config.teams["default"] = lines[0]
		config.defaultTeam = "default"
		return config
	}

	//advanced config file
	for _, line := range lines {
		s := strings.Split(line, "=")
		if len(s) != 2 {
			exitErr(fmt.Errorf("failed to parse config at: %s", line))
		}
		key, val := strip(s[0]), strip(s[1])
		switch key {
		case "default_team":
			config.defaultTeam = val
		case "default_channel":
			config.defaultChannel = val
		default:
			config.teams[key] = val
		}
	}
	return config
}

func strip(s string) string {
	return strings.Replace(s, " ", "", -1)
}

func readLines(path string) []string {
	var lines []string

	file, err := os.Open(path)
	failOnError(err, "unable to read config", true)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines
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
