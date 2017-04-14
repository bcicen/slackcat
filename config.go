package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/skratchdot/open-golang/open"
)

const configURL = "http://slackcat.chat/configure"

//Slack team and channel read from file
type Config struct {
	Teams          map[string]string `toml:"teams"`
	DefaultTeam    string            `toml:"default_team"`
	DefaultChannel string            `toml:"default_channel"`
}

func (c *Config) parseChannelOpt(channel string) (string, string, error) {
	//use default channel if none provided
	if channel == "" {
		if c.DefaultChannel == "" {
			return "", "", fmt.Errorf("no channel provided")
		}
		return c.DefaultTeam, c.DefaultChannel, nil
	}
	//if channel is prefixed with a team
	if strings.Contains(channel, ":") {
		s := strings.Split(channel, ":")
		return s[0], s[1], nil
	}
	//use default team with provided channel
	return c.DefaultTeam, channel, nil
}

func readConfig() Config {
	var config Config

	lines := readLines(getConfigPath())

	//simple config file
	if len(lines) == 1 {
		config.Teams["default"] = lines[0]
		config.DefaultTeam = "default"
		return config
	}

	//advanced config file
	body := strings.Join(lines, "\n")
	if _, err := toml.Decode(body, &config); err != nil {
		exitErr(fmt.Errorf("failed to parse config: %s", err))
	}
	return config
}

func getConfigPath() string {
	userHome, ok := os.LookupEnv("HOME")
	if !ok {
		exitErr(fmt.Errorf("$HOME not set"))
	}

	if xdgSupport() {
		xdgHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			xdgHome = fmt.Sprintf("%s/.config", userHome)
		}
		return fmt.Sprintf("%s/slackcat/config", xdgHome)
	}

	return fmt.Sprintf("%s/.slackcat", userHome)
}

func xdgSupport() bool {
	re := regexp.MustCompile("^XDG_*")
	for _, e := range os.Environ() {
		if re.FindAllString(e, 1) != nil {
			return true
		}
	}
	return false
}

func strip(s string) string { return strings.Replace(s, " ", "", -1) }

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
	output("Creating token request for Slackcat")
	open.Run(configURL)
	output("Use the below URL to authorize slackcat if browser fails to launch")
	output(configURL)
}
