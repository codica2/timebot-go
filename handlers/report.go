package handlers

import (
	"fmt"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
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
	days := make(map[string][]*models.TimeEntry)

	today := models.Today()

	for _, timeEntry := range timeEntries {
		if entries, ok := days[timeEntry.Date.String()]; ok {
			days[timeEntry.Date.String()] = append(entries, timeEntry)
		} else {
			days[timeEntry.Date.String()] = []*models.TimeEntry{timeEntry}
		}
	}

	stringArray := models.NewStringArray()

	for d := from; to.CompareTo(&d) >= 0 && d.CompareTo(&today) <= 0; d = d.Plus(1) {
		entries := days[d.String()]

		line := fmt.Sprintf("`%s:", d.Format("02.01.06` (Monday)"))

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

	stringArray.Add(fmt.Sprintf("*Total*: %s.", totalTime(timeEntries)))

	sender.SendMessage(user.UID, stringArray.Join("\n"))
}

func totalTime(entries []*models.TimeEntry) string {
	total := 0

	for _, entry := range entries {
		total += entry.Minutes
	}

	minutes := total % 60
	hours := total / 60

	return fmt.Sprintf("%d:%02d", hours, minutes)
}
