package main

import (
	"fmt"

	"github.com/bluele/slack"
)

var (
	api     *slack.Slack
	msgOpts = &slack.ChatPostMessageOpt{AsUser: true}
)

func InitAPI(token string) {
	api = slack.New(token)
	res, err := api.AuthTest()
	failOnError(err, "Slack API Error")
	output(fmt.Sprintf("connected to %s as %s", res.Team, res.User))
}

// Return list of all channels by name
func listChannels() (names []string) {
	list, err := api.ChannelsList()
	failOnError(err)
	for _, c := range list {
		names = append(names, c.Name)
	}
	return names
}

// Return list of all groups by name
func listGroups() (names []string) {
	list, err := api.GroupsList()
	failOnError(err)
	for _, c := range list {
		names = append(names, c.Name)
	}
	return names
}

// Return list of all ims by name
func listIms() (names []string) {
	users, err := api.UsersList()
	failOnError(err)

	list, err := api.ImList()
	failOnError(err)
	for _, c := range list {
		for _, u := range users {
			if u.Id == c.User {
				names = append(names, u.Profile.RealName)
				continue
			}
		}
	}
	return names
}

// Lookup Slack id for im by real name
func findImByRealName(realName string) string {
	users, err := api.UsersList()
	failOnError(err)
	for _, u := range users {
		if u.Profile.RealName == realName {
			return u.Id
		}
	}
	exitErr(fmt.Errorf("No such channel, group, or im"))
	return ""
}

// Lookup Slack id for channel, group, or im by name or real name
func lookupSlackID(name string) string {
	if channel, err := api.FindChannelByName(name); err == nil {
		return channel.Id
	}
	if group, err := api.FindGroupByName(name); err == nil {
		return group.Id
	}
	if im, err := api.FindImByName(name); err == nil {
		return im.Id
	}
	return findImByRealName(name)
}
