package main

import "gopkg.in/jmcvetta/napping.v1"

type slackNotification struct {
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Channel   string `json:"channel"`
	Text      string `json:"text"`
}

func notifyUser(message string, user Member) {
	notif := slackNotification{}

	notif.Channel = "@" + user.Name
	notif.Username = "Kudo Panda"
	notif.IconEmoji = ":panda_face:"
	notif.Text = message

	napping.Post(config.SlackWebookURL, notif, nil, nil)
}

func notifyChannel(message string) {
	notif := slackNotification{}

	notif.Channel = "#" + config.Channel
	notif.Username = "Kudo Panda"
	notif.IconEmoji = ":panda_face:"
	notif.Text = message

	napping.Post(config.SlackWebookURL, notif, nil, nil)
}
