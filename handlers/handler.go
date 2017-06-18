package handlers

import (
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
)

const newEntryStringRegexp = "^ *(.*) (\\d?\\d:[0-5]\\d) ([^\\s](?:.|\\s)*[^\\s])$"

func HandleMessage(message *slack.Msg) {
	if matched, err := regexp.MatchString(newEntryStringRegexp, message.Text); matched && err == nil {
		fmt.Printf("Message \"%s\" is create new entry\n", message.Text)
		handleNewEntry(message)
	} else {
		sender.SendMessage(message.User, "Sorry. I don't understand you.")
	}
}

func handleNewEntry(message *slack.Msg) {
	newEntryRegexp := regexp.MustCompile(newEntryStringRegexp)
	matches := newEntryRegexp.FindStringSubmatch(message.Text)

	projectName := matches[1]
	time := matches[2]
	description := matches[3]
	user, err := models.FindUser(message.User)

	if err != nil {
		sender.SendMessage(user.UID, "Sorry. An error occurred.")
	} else {
		sender.SendMessage(user.UID, "The time entry will be created soon.")
	}

	fmt.Printf("Project: %s; time: %s; description: %s\n", projectName, time, description)
}
