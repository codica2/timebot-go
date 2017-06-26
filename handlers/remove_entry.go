package handlers

import (
	"regexp"
	"strconv"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	"fmt"
)

const removeEntryRegexp = "^ *remove entry (\\d+) *$"

func handleRemoveEntry(uid, text string) {
	r := regexp.MustCompile(removeEntryRegexp)

	id, err := strconv.ParseInt(r.FindStringSubmatch(text)[1], 10, 64)

	if err != nil {
		handleError(uid, err)
		return
	}

	user, err := models.FindUser(uid)

	if err != nil {
		handleError(uid, err)
		return
	}

	entry, err := models.FindTimeEntryByID(int(id))

	if err != nil {
		handleError(uid, err)
		return
	}

	if user.ID != entry.UserId {
		sender.SendMessage(uid, "You are not allowed to remove other user's entries.")
		return
	}

	entry.Delete()

	sender.SendMessage(uid, fmt.Sprintf("Time entry with id %d was successfully deleted.", entry.ID))
}
