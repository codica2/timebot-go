package handlers

import (
	"github.com/nlopes/slack"
	"regexp"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"database/sql"
	"time"
)

const createEntryForDayRegexp = "^ *update (\\d?\\d\\.\\d?\\d(?:\\.(?:\\d\\d)?\\d\\d)?) (.*) (\\d?\\d:[0-5]\\d) ([^\\s](?:.|\\s)*[^\\s]) *$"

func handleCreateEntryForDay(msg *slack.Msg) {
	r := regexp.MustCompile(createEntryForDayRegexp)

	matchData := r.FindStringSubmatch(msg.Text)

	date := matchData[1]
	projectName := matchData[2]
	t := matchData[3]
	details := matchData[4]

	user, err := models.FindUser(msg.User)

	if err != nil {
		handleError(msg.User, err)
		return
	}

	project, err := models.FindProjectByNameOrAlias(projectName)

	if _, ok := err.(models.NotFoundError); ok {
		sender.SendMessage("Project with name \"%s\" was not found.", projectName)
		return
	}

	minutes, err := parseTime(t)

	if err != nil {
		handleError(msg.User, err)
		return
	}

	d, err := models.ParseDate(date)

	if err != nil {
		handleError(msg.User, err)
		return
	}

	timeEntry := models.TimeEntry{}

	now := time.Now()

	timeEntry.UserId = user.ID
	timeEntry.ProjectId = project.ID
	timeEntry.Date = *d
	timeEntry.Time = t
	timeEntry.Minutes = minutes
	timeEntry.Details = sql.NullString{details, true}
	timeEntry.CreatedAt = now
	timeEntry.UpdatedAt = now

	err = timeEntry.Create()

	if err != nil {
		handleError(msg.User, err)
		return
	}

	sender.SendMessage(user.UID, "Time entry was successfully created")
}
