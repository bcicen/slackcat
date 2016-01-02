# slackcat
Slackcat is a simple commandline utility to post snippets to Slack.

<p align="center">
  <img src="https://raw.githubusercontent.com/bcicen/slackcat/master/demo.gif" alt="slackcat"/>
</p>

## Usage
Pipe command output:
```bash
$ echo -e "hi\nthere" | slackcat --channel general --filename hello
file hello uploaded to general
```

Post an existing file:
```bash
$ slackcat -c general /home/user/bot.png
file bot.png uploaded to general
```

## Installing

Download the latest release for your platform:

#### OS X

```bash
wget https://github.com/bcicen/slackcat/releases/download/v0.3/slackcat-0.3-darwin-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

#### Linux

```bash
wget https://github.com/bcicen/slackcat/releases/download/v0.3/slackcat-0.3-linux-amd64 -O slackcat
sudo mv slackcat /usr/local/bin/
sudo chmod +x /usr/local/bin/slackcat
```

and create a Slackcat config file:
```bash
echo '<your-slack-token>' > ~/.slackcat
```

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--noop | Skip posting file to Slack. Useful for testing
--channel, -c | Slack channel to post to
--filename, -n | Filename for upload. Defaults to given filename or current timestamp if reading from stdin.
