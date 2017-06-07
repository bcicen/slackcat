package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var (
	noop    = false
	build   = ""
	version = "dev-build"
)

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
	failOnError(err, "unable to create tmpfile")

	w := bufio.NewWriter(tmp)
	for line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()

	return tmp.Name()
}

func handleUsageError(c *cli.Context, err error, _ bool) error {
	fmt.Fprintf(c.App.Writer, "%s %s\n\n", "Incorrect Usage.", err.Error())
	cli.ShowAppHelp(c)
	return cli.NewExitError("", 1)
}

func printFullVersion(c *cli.Context) {
	fmt.Fprintf(c.App.Writer, "%v version %v, build %v\n", c.App.Name, c.App.Version, build)
}

func main() {
	cli.VersionPrinter = printFullVersion

	app := cli.NewApp()
	app.Name = "slackcat"
	app.Usage = "redirect a file to slack"
	app.Version = version
	app.OnUsageError = handleUsageError
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "channel, c",
			Usage: "Slack channel or group to post to",
		},
		cli.StringFlag{
			Name:  "comment",
			Usage: "Initial comment for snippet",
		},
		cli.BoolFlag{
			Name:  "configure",
			Usage: "Configure Slackcat via oauth",
		},
		cli.StringFlag{
			Name:  "filename, n",
			Usage: "Filename for upload. Defaults to current timestamp",
		},
		cli.StringFlag{
			Name:  "filetype",
			Usage: "Specify filetype for syntax highlighting",
		},
		cli.BoolFlag{
			Name:  "list",
			Usage: "List team channel names",
		},
		cli.BoolFlag{
			Name:  "noop",
			Usage: "Skip posting file to Slack. Useful for testing",
		},
		cli.BoolFlag{
			Name:  "stream, s",
			Usage: "Stream messages to Slack continuously instead of uploading a single snippet",
		},
		cli.BoolFlag{
			Name:  "tee, t",
			Usage: "Print stdin to screen before posting",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("configure") {
			configureOA()
			os.Exit(0)
		}

		configPath, exists := getConfigPath()
		if !exists {
			exitErr(fmt.Errorf("missing config file at %s\nuse --configure to create", configPath))
		}
		config := ReadConfig(configPath)

		team, channel, err := config.parseChannelOpt(c.String("channel"))
		failOnError(err)

		noop = c.Bool("noop")
		fileName := c.String("filename")
		fileType := c.String("filetype")
		fileComment := c.String("comment")

		token := config.Teams[team]
		if token == "" {
			exitErr(fmt.Errorf("no such team: %s", team))
		}

		slackcat := newSlackcat(token, channel)

		if c.Bool("list") {
			fmt.Println("channels:")
			for _, n := range slackcat.listChannels() {
				fmt.Printf("  %s\n", n)
			}
			fmt.Println("groups:")
			for _, n := range slackcat.listGroups() {
				fmt.Printf("  %s\n", n)
			}
			fmt.Println("ims:")
			for _, n := range slackcat.listIms() {
				fmt.Printf("  %s\n", n)
			}
			os.Exit(0)
		}

		if len(c.Args()) > 0 {
			if c.Bool("stream") {
				output("filepath provided, ignoring stream option")
			}
			filePath := c.Args()[0]
			if fileName == "" {
				fileName = filepath.Base(filePath)
			}
			slackcat.postFile(filePath, fileName, fileType, fileComment)
			os.Exit(0)
		}

		lines := make(chan string)
		go readIn(lines, c.Bool("tee"))

		if c.Bool("stream") {
			slackcat.stream(lines)
		} else {
			filePath := writeTemp(lines)
			defer os.Remove(filePath)
			slackcat.postFile(filePath, fileName, fileType, fileComment)
			os.Exit(0)
		}
	}

	app.Run(os.Args)

}
