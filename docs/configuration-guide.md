# Configuration Guide
Slackcat may be configured via a simple or advanced configuration.

## Default Path
If your environment specifies an [XDG Base Directory](https://specifications.freedesktop.org/basedir-spec/latest/index.html), `slackcat` will use the configuration file at `~/.config/slackcat/config`; otherwise, will fallback to `~/.slackcat`

## Simple Configuration

Generate a new Slack token with:
```bash
slackcat --configure
```
A new browser window will be opened for you to confirm the request via Slack, and you'll be returned a token.

Create a Slackcat config file and you're ready to go!
```bash
echo '<your-slack-token>' > ~/.slackcat
```

## Advanced Configuration

Advanced configuration allows for multiple Slack teams, a default team, and default channel in [TOML](https://github.com/toml-lang/toml) format.

#### Example ~/.config/slackcat Config
```bash
default_team = "team1"
default_channel = "general"

[teams]
  team1 = "<team1-slack-token>"
  team2 = "<team2-slack-token>"
```
By default, all messages will be sent to the team1 general channel.

#### Example Usage

Post a file to team1 #general channel:
```bash
slackcat /path/to/file.txt
```

Post a file to team1 #testing channel:
```bash
slackcat -c testing /path/to/file.txt
```

Post a file to team2 #testing channel:
```bash
slackcat -c team2:testing /path/to/file.txt
```
