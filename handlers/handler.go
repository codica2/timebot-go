package handlers

import (
	"fmt"
	"github.com/alex-bogomolov/timebot_go/sender"
	"github.com/nlopes/slack"
	"regexp"
	"strconv"
	"strings"
)

func HandleMessage(message *slack.Msg) {
	if matched, err := regexp.MatchString(createEntryForDayRegexp, message.Text); matched && err == nil {
		handleCreateEntryForDay(message)
	} else if matched, err = regexp.MatchString(newEntryStringRegexp, message.Text); matched && err == nil {
		fmt.Printf("Message \"%s\" is create new entry\n", message.Text)
		handleNewEntry(message)
	} else if matched, err = regexp.MatchString(reportRegexpString, message.Text); matched && err == nil {
		handleReport(message)
	} else if matched, err = regexp.MatchString(showProjectsRegexp, message.Text); matched && err == nil {
		handleShowProjects(message.User)
	} else if matched, err = regexp.MatchString(addProjectRegexp, message.Text); matched && err == nil {
		handleAddProject(message.User, message.Text)
	} else if matched, err = regexp.MatchString(showHelpRegexp, strings.ToLower(message.Text)); matched && err == nil {
		handleShowHelp(message.User)
	} else if matched, err = regexp.MatchString(removeEntryRegexp, strings.ToLower(message.Text)); matched && err == nil {
		handleRemoveEntry(message.User, message.Text)
	} else {
		sender.SendMessage(message.User, "Sorry. I don't understand you.")
	}
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
