package handlers

import (
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
	"strconv"
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
	minutes, err := parseTime(entryTime)

	if err != nil {
		handleError(message.User, err)
		return
	}

	details := matches[3]
	user, err := models.FindUser(message.User)

	if err != nil {
		handleError(message.User, err)
		return
	}

	project, err := models.FindProjectByNameOrAlias(projectName)

	if err != nil && err.Error() == "The project was not found" {
		sender.SendMessage(user.UID, "The project with name \""+projectName+"\" was not found.")
		return
	} else if err != nil {
		handleError(user.UID, err)
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
		handleError(message.User, err)
		return
	}

	sender.SendMessage(user.UID, "The time entry was successfully created.")
}

func parseTime(time string) (int, error) {
	regex := regexp.MustCompile("^(\\d?\\d):(\\d\\d)$")

	matchData := regex.FindStringSubmatch(time)

	hours, err := strconv.ParseInt(matchData[1], 10, 64)

	if err != nil {
		return 0, err
	}

	minutes, err := strconv.ParseInt(matchData[2], 10, 64)

	if err != nil {
		return 0, err
	}

	return int(hours)*60 + int(minutes), nil
}

func handleError(uid string, err error) {
	sender.SendMessage(uid, fmt.Sprintf("An error occured: %s", err.Error()))
}
