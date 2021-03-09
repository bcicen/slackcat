package main

import (
	"fmt"

	"github.com/slack-go/slack"
)

var (
	api     *slack.Client
	msgOpts = &slack.PostMessageParameters{AsUser: true}
)

// InitAPI for Slack
func InitAPI(token string) {
	api = slack.New(token)
	res, err := api.AuthTest()
	failOnError(err, "Slack API Error")
	output(fmt.Sprintf("connected to %s as %s", res.Team, res.User))
}

// Return list of all channels by name
func listChannels() (names []string) {
	list, err := getConversations("public_channel", "private_channel")
	failOnError(err)
	for _, c := range list {
		names = append(names, c.Name)
	}
	return names
}

// Return list of all groups by name
func listGroups() (names []string) {
	list, err := getConversations("mpim")
	failOnError(err)
	for _, c := range list {
		names = append(names, c.Name)
	}
	return names
}

// Return list of all ims by name
func listIms() (names []string) {
	users, err := api.GetUsers()
	failOnError(err)

	list, err := getConversations("im")
	failOnError(err)
	for _, c := range list {
		for _, u := range users {
			if u.ID == c.User {
				names = append(names, u.Name)
				continue
			}
		}
	}
	return names
}

// Lookup Slack id for channel, group, or im by name
func lookupSlackID(name string) string {
	list, err := getConversations("public_channel", "private_channel", "mpim")
	if err == nil {
		for _, c := range list {
			if c.Name == name {
				return c.ID
			}
		}
	}
	users, err := api.GetUsers()
	if err == nil {
		list, err := getConversations("im")
		if err == nil {
			for _, c := range list {
				for _, u := range users {
					if u.ID == c.User {
						return c.ID
					}
				}
			}
		}
	}
	exitErr(fmt.Errorf("No such channel, group, or im"))
	return ""
}

func getConversations(types ...string) (list []slack.Channel, err error) {
	cursor := ""
	for {
		param := &slack.GetConversationsParameters{
			Cursor:          cursor,
			ExcludeArchived: "true",
			Types:           types,
		}
		channels, cur, err := api.GetConversations(param)
		if err != nil {
			return list, err
		}
		list = append(list, channels...)
		if cur == "" {
			break
		}
		cursor = cur
	}
	return list, nil
}
