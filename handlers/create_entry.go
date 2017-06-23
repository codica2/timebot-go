package handlers

import (
	"database/sql"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
	"time"
)

const newEntryStringRegexp = "^ *(.*) (\\d?\\d:[0-5]\\d) ([^\\s](?:.|\\s)*[^\\s])$"

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

	if _, ok := err.(models.NotFoundError); ok {
		sender.SendMessage(user.UID, "The project with name \""+projectName+"\" was not found.")
		return
	} else if err != nil {
		handleError(user.UID, err)
		return
	}

	timeEntry := models.TimeEntry{
		UserId:    user.ID,
		Date:      models.NewDate(time.Now()),
		ProjectId: project.ID,
		Details:   sql.NullString{details, true},
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
