package timebot

import ("fmt"
	    "github.com/nlopes/slack")

type OutgoingMessage struct {
	Text string
	User string
}


func HandleMessage(ev *slack.MessageEvent, outgoingMessages chan *OutgoingMessage, users map[string]string) {
	fmt.Printf("Received message from %s with text \"%s\"", users[ev.Msg.User], ev.Msg.Text)
	outgoingMessage := OutgoingMessage{Text: "I'm maintained now", User: ev.Msg.User}

	outgoingMessages <- &outgoingMessage
}