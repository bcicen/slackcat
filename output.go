// output/err convenience methods
package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	bold = color.New(color.Bold).SprintFunc()
	red  = color.New(color.FgRed).SprintFunc()
	cyan = color.New(color.FgCyan).SprintFunc()
)

func output(s string) {
	fmt.Printf("%s %s\n", cyan("slackcat"), s)
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
