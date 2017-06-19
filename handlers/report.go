package handlers

import (
	"github.com/nlopes/slack"
	"regexp"
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"time"
)

const reportRegexpString = "^ *show (week|last week|month|last month)(?: (.*?))? *$"

func handleReport(message *slack.Msg) {
	reportRegexp := regexp.MustCompile(reportRegexpString)

	matchData := reportRegexp.FindStringSubmatch(message.Text)

	interval := matchData[1]
	projectName := matchData[2]

	user, err := models.FindUser(message.User)

	if err != nil {
		handleError(message.User, err)
		return
	}

	var project *models.Project


	if len(projectName) != 0 {
		project, err = models.FindProjectByNameOrAlias(projectName)

		if _, ok := err.(models.NotFoundError); ok {
			sender.SendMessage(user.UID, fmt.Sprintf("The project with name \"%s\" was not found.", projectName))
			return
		} else if err != nil {
			handleError(user.UID, err)
			return
		}
	}

	timeEntries := []*models.TimeEntry{}
	var from, to time.Time

	switch interval {
	case "week":
		from = startOfWeek()
		to = endOfWeek()
	case "last week":
		from = time.Unix(startOfWeek().Unix() - 7 * 24 * 60 * 60, 0)
		to = time.Unix(endOfWeek().Unix() - 7 * 24 * 60 * 60, 0)
	case "month":
	case "last month":
	}

	timeEntries, err = models.GetTimeEntriesInPeriodWithProjectAndUser(user, project, from, to)

	if err != nil {
		handleError(user.UID, err)
	}

	displayTimeEntries(timeEntries, from, to)
}

func displayTimeEntries(timeEntries []*models.TimeEntry, from time.Time, to time.Time) {
	fmt.Println(timeEntries)
}

func startOfWeek() time.Time {
	now := time.Now()
	unix := now.Unix()
	weekday := int64(now.Weekday())

	if weekday == 0 {
		return time.Unix(unix - 24 * 60 * 60 * 6, 0)
	} else {
		return time.Unix(unix - 24 * 60 * 60 * (weekday - 1), 0)
	}
}

func endOfWeek() time.Time {
	now := time.Now()
	unix := now.Unix()
	weekday := int64(now.Weekday())

	if weekday == 0 {
		return now
	} else {
		return time.Unix(unix - 24 * 60 * 60 * (7 - weekday), 0)
	}

}
