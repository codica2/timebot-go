package handlers

import (
	"os"
	"io/ioutil"
	"fmt"
	"github.com/alex-bogomolov/timebot_go/sender"
)

const showHelpRegexp = "^ *help *$"

func handleShowHelp(uid string) {
	projectPath := os.Getenv("TIMEBOT_GO_PATH")

	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/messages/help.txt", projectPath))

	if err != nil {
		handleError(uid, err)
		return
	}

	text := string(bytes)

	sender.SendMessage(uid, text)
}
