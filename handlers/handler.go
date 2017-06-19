package handlers

import (
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
	"time"
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
	entryTime := matches[2]
	minutes := parseTime(entryTime)
	details := matches[3]
	user, err := models.FindUser(message.User)

	if err != nil {
		sender.SendMessage(user.UID, "Sorry. An error occurred.")
		return
	}

	project, err := models.FindProjectByNameOrAlias(projectName)

	if err != nil && err.Error() == "The project was not found" {
		fmt.Println(err)
		sender.SendMessage(user.UID, "The project with name \""+projectName+"\" was not found.")
		return
	} else if err != nil {
		fmt.Println(err)
		sender.SendMessage(user.UID, "Sorry. An error occurred.")
		return
	}

	timeEntry := models.TimeEntry{
		UserId:    user.ID,
		Date:      time.Now(),
		ProjectId: project.ID,
		Details:   details,
		Minutes:   minutes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Time:      entryTime,
	}

	err = timeEntry.Create()

	if err != nil {
		fmt.Println(err)
		sender.SendMessage(user.UID, "Sorry. An error occurred.")
		return
	}

	sender.SendMessage(user.UID, "The time entry was successfully created.")
}

func parseTime(time string) int {
	regex := regexp.MustCompile("^(\\d?\\d):(\\d\\d)$")

	matchData := regex.FindStringSubmatch(time)

	hours := int(matchData[1])
	minutes := int(matchData[2])

	return hours*60 + minutes
}
