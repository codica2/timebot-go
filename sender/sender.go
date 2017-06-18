package sender

import "github.com/nlopes/slack"

var Api *slack.Client

func SendMessage(receiver, text string) {
	Api.SendMessage(receiver, slack.MsgOptionText(text, true), slack.MsgOptionAsUser(true))
}
