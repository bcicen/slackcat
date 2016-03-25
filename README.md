# slackcat
Slackcat is a simple commandline utility to post snippets to Slack.


  <img width="500px" src="https://raw.githubusercontent.com/vektorlab/slackcat/master/demo.gif" alt="slackcat"/>


## Quickstart

### Installing
Download the latest release for your platform:

#### OS X

```brew
brew install slackcat
```
or
```bash
curl -Lo slackcat https://github.com/vektorlab/slackcat/releases/download/v1.0/slackcat-1.0-darwin-amd64
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

#### Linux

```bash
wget https://github.com/vektorlab/slackcat/releases/download/v1.0/slackcat-1.0-linux-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

### Configuration

Generate a new Slack token with:
```bash
slackcat --configure
```
A new browser window will be opened for you to confirm the request via Slack, and you'll be returned a token.

Create a Slackcat config file and you're ready to go!
```bash
echo '<your-slack-token>' > ~/.slackcat
```

For configuring multiple teams and default channels, see [Configuration Guide](https://github.com/vektorlab/slackcat/blob/master/docs/configuration-guide.md).

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

Stream input continously as a formatted message:
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
--plain, -p | When streaming, write messages as plain text instead of code blocks
--noop | Skip posting file to Slack. Useful for testing
--configure | Configure Slackcat via oauth
--channel, -c | Slack channel, group, or user to post to
--filename, -n | Filename for upload. Defaults to given filename or current timestamp if reading from stdin
--filetype | Specify filetype for synax highlighting. Defaults to autodetect
--comment | Initial comment for snippet
