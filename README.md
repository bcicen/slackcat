# slackcat
Slackcat is a simple commandline utility to post snippets to Slack.

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

## Options

Option | Description
--- | ---
--tee, -t | Print stdin to screen before posting
--noop | Skip posting file to Slack. Useful for testing
--channel, -c | Slack channel to post to
--filename, -n | Filename for upload. Defaults to current timestamp
