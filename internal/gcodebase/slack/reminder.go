package slack

import "galaxy_walker/internal/github.com/nlopes/slack"

type SlackReminder struct {
	Client           *slack.Client
	SlackChannelName string
}

var _slack_reminder map[string]*slack.Client

func GetSlackReminder(token, channelName string) *SlackReminder {
	if _, exist := _slack_reminder[token]; !exist {
		_slack_reminder[token] = slack.New(token)
	}
	return &SlackReminder{
		Client:           _slack_reminder[token],
		SlackChannelName: channelName,
	}
}

func init() {
	_slack_reminder = make(map[string]*slack.Client)
}
