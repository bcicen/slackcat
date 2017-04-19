// Helper and convenience methods
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	bold = color.New(color.Bold).SprintFunc()
	red  = color.New(color.FgRed).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
)

func strip(s string) string { return strings.Replace(s, " ", "", -1) }

func basedir(path string) string {
	parts := strings.Split(path, "/")
	return strings.Join((parts[0 : len(parts)-1]), "/")
}

func output(s string) {
	fmt.Printf("%s %s\n", bold(cyan("slackcat")), s)
}

func failOnError(err error, msg ...string) {
	if err != nil {
		if msg != nil {
			err = fmt.Errorf("%s: %s", msg[0], err)
		}
		exitErr(err)
	}
}

func appendErr(msg string, err error) error {
	return fmt.Errorf("%s: %s", msg, err)
}

func exitErr(err error) {
	output(red(err.Error()))
	os.Exit(1)
}
