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

const AuthURL = "http://slackcat.chat/configure"

// Slack team and channel read from file
type Config struct {
	Teams          map[string]string `toml:"teams"`
	DefaultTeam    string            `toml:"default_team"`
	DefaultChannel string            `toml:"default_channel"`
}

// NewConfig returns new default config
func NewConfig() *Config {
	return &Config{
		Teams: make(map[string]string),
	}
}

// ReadConfig returns config read from file
func ReadConfig(path string) *Config {
	config := NewConfig()
	lines, err := readLines(path)
	failOnError(err, "unable to read config")

	// simple config file
	if len(lines) == 1 {
		config.Teams["default"] = lines[0]
		config.DefaultTeam = "default"
		return config
	}

	// advanced config file
	body := strings.Join(lines, "\n")
	_, err = toml.Decode(body, &config)
	failOnError(err, "failed to parse config")

	return config
}

func (c *Config) Write(path string) {
	cfgdir := basedir(path)
	// create config dir if not exist
	if _, err := os.Stat(cfgdir); err != nil {
		err = os.MkdirAll(cfgdir, 0755)
		if err != nil {
			exitErr(fmt.Errorf("failed to initialize config dir [%s]: %s", cfgdir, err))
		}
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		exitErr(fmt.Errorf("failed to open config for writing: %s", err))
	}

	writer := toml.NewEncoder(file)
	err = writer.Encode(c)
	if err != nil {
		exitErr(fmt.Errorf("failed to write config: %s", err))
	}
}

func (c *Config) parseChannelOpt(channel string) (string, string, error) {
	// use default channel if none provided
	if channel == "" {
		if c.DefaultChannel == "" {
			return "", "", fmt.Errorf("no channel provided")
		}
		return c.DefaultTeam, c.DefaultChannel, nil
	}
	// if channel is prefixed with a team
	if strings.Contains(channel, ":") {
		s := strings.Split(channel, ":")
		return s[0], s[1], nil
	}
	// use default team with provided channel
	return c.DefaultTeam, channel, nil
}

// determine config path from environment
func getConfigPath() (path string, exists bool) {
	userHome, ok := os.LookupEnv("HOME")
	if !ok {
		exitErr(fmt.Errorf("$HOME not set"))
	}

	path = fmt.Sprintf("%s/.slackcat", userHome) // default path

	if xdgSupport() {
		xdgHome, ok := os.LookupEnv("XDG_CONFIG_HOME")
		if !ok {
			xdgHome = fmt.Sprintf("%s/.config", userHome)
		}
		path = fmt.Sprintf("%s/slackcat/config", xdgHome)
	}

	if _, err := os.Stat(path); err == nil {
		exists = true
	}

	return path, exists
}

func basedir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join((parts[0 : len(parts)-1]), "/")
}

// Test for environemnt supporting XDG spec
func xdgSupport() bool {
	re := regexp.MustCompile("^XDG_*")
	for _, e := range os.Environ() {
		if re.FindAllString(e, 1) != nil {
			return true
		}
	}
	return false
}

func configureOA() {
	var nick, token string
	var config *Config

	cfgPath, cfgExists := getConfigPath()
	if !cfgExists {
		config = NewConfig()
	} else {
		config = ReadConfig(cfgPath)
	}

	fmt.Printf("nickname for team: ")
	fmt.Scanf("%s", &nick)
	if nick == "" {
		exitErr(fmt.Errorf("no name provided"))
	}

	output("creating token request for slackcat")
	open.Run(AuthURL)
	output("Use the below URL to authorize slackcat if browser fails to launch")
	output(AuthURL)

	fmt.Printf("token issued: ")
	fmt.Scanf("%s", &token)
	if token == "" {
		exitErr(fmt.Errorf("no token provided"))
	}

	// creating a new config file
	if !cfgExists {
		config.DefaultTeam = nick
	}
	config.Teams[nick] = token
	config.Write(cfgPath)

	output(fmt.Sprintf("added team to config file at %s", cfgPath))
}

func readLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			lines = append(lines, scanner.Text())
		}
	}
	return lines, nil
}
