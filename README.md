# slackcat
Slackcat is a simple commandline utility to post snippets to Slack.


  <img width="500px" src="https://raw.githubusercontent.com/bcicen/slackcat/master/demo.gif" alt="slackcat"/>


## Installing
Download the [latest release](https://github.com/bcicen/slackcat/releases) for your platform:

```bash
curl -Lo slackcat https://github.com/bcicen/slackcat/releases/download/1.7.2/slackcat-1.7.2-$(uname -s)-amd64
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

`slackcat` is also available via homebrew:
```brew
brew install slackcat
```

## Building
To optionally build `slackcat` from source, ensure you have [dep](https://github.com/golang/dep) installed and run:
```
$ git clone https://github.com/bcicen/slackcat.git
$ cd slackcat
$ go get github.com/bcicen/slackcat
$ make build
```

You must use GNU make for the build to work correctly. If your platform does not install GNU make as `make` (i.e. OpenBSD) then you will need to install gmake and run:
```
$ git clone https://github.com/bcicen/slackcat.git
$ cd slackcat
$ go get github.com/bcicen/slackcat
$ gmake build
```

## Configuration

Generate an initial config, or add a new team token with:
```bash
slackcat --configure
```
You'll be prompted for a team nickname and a new browser window will be opened for you to confirm the request via Slack. Provide the returned token to slackcat when prompted, and you're ready to go!

For configuring multiple teams and default channels, see [Configuration Guide](https://github.com/bcicen/slackcat/blob/master/docs/configuration-guide.md).

## Usage
Pipe command output as a text snippet:
```bash
$ echo -e "hi\nthere" | slackcat --channel general --filename hello
*slackcat* file hello uploaded to general
```

Post an existing file:
```bash
$ slackcat --channel general /home/user/bot.png
*slackcat* file bot.png uploaded to general
```

Stream input continuously:
```bash
$ tail -F -n0 /path/to/log | slackcat --channel general --stream
*slackcat* posted 5 message lines to general
*slackcat* posted 2 message lines to general
...
```

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--stream, -s | Stream messages to Slack continuously instead of uploading a single snippet
--noop | Skip posting file to Slack. Useful for testing
--configure | Configure Slackcat via oauth
--iconemoji, -i | Specify emoji icon for message (e.g. ":+1:")
--channel, -c | Slack channel, group, or user to post to
--filename, -n | Filename for upload. Defaults to given filename or current timestamp if reading from stdin
--filetype | Specify filetype for syntax highlighting. Defaults to autodetect
--comment | Initial comment for snippet
--username | Stream messages as given bot user. Defaults to auth user
--thread | Stream messages to thread after initial comment message
