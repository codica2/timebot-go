package handlers

import (
	"github.com/nlopes/slack"
	"regexp"
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
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
	var from, to models.Date

	switch interval {
	case "week":
		from = models.Today().StartOfWeek()
		to = models.Today().EndOfWeek()
	case "last week":
		from = models.Today().StartOfWeek().Minus(7)
		to = models.Today().EndOfWeek().Minus(7)
	case "month":
		from = models.BeginningOfMonth()
		to = models.EndOfMonth()
	case "last month":
		from = models.BeginningOfLastMonth()
		to = models.EndOfLastMonth()
	}

	timeEntries, err = models.GetTimeEntriesInPeriodWithProjectAndUser(user, project, from, to)

	if err != nil {
		handleError(user.UID, err)
		return
	}

	displayTimeEntries(timeEntries, from, to, user)
}

func displayTimeEntries(timeEntries []*models.TimeEntry, from models.Date, to models.Date, user *models.User) {
	stringArray := models.NewStringArray()

	today := models.Today()

	for d := from; to.CompareTo(&d) >= 0 && d.CompareTo(&today) <= 0; d = d.Plus(1) {
		entries := findEntries(timeEntries, &d)

		line := fmt.Sprintf("`%s:", d.Format("02.01.06` (Mon)"))

		if len(entries) == 0 {
			line += " No entries"
			stringArray.Add(line)
		} else {
			stringArray.Add(line)
			for _, entry := range entries {
				stringArray.Add(entry.Description())
			}
		}
	}

	sender.SendMessage(user.UID, stringArray.Join("\n"))
}

func findEntries(entries []*models.TimeEntry, d *models.Date) []*models.TimeEntry {
	out := []*models.TimeEntry{}

	for _, entry := range entries {
		if entry.Date.Equal(d) {
			out = append(out, entry)
		}
	}

	return out
}
