package main

import (
	"database/sql"
	"fmt"
	"github.com/alex-bogomolov/timebot_go/handlers"
	"github.com/alex-bogomolov/timebot_go/models"
	"github.com/alex-bogomolov/timebot_go/sender"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
	"os"
)

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)
	sender.Api = api

	timebotId, err := getTimebotId(api)

	if err != nil {
		fmt.Println(err)
		return
	}

	models.DB, err = connectToDatabase()

	if err != nil {
		fmt.Println(err)
		return
	}

	startBot(api, timebotId)
}

func getTimebotId(api *slack.Client) (string, error) {
	users, err := api.GetUsers()

	if err != nil {
		return "", err
	}

	var timebotId string

	for _, user := range users {
		if user.Name == "timebot" {
			timebotId = user.ID
			break
		}
	}

	return timebotId, nil
}

func startBot(api *slack.Client, timebotId string) {
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack")
		case *slack.MessageEvent:
			if event.Msg.User != timebotId && os.Getenv("GOLANG_ENV") == "development" && event.Msg.User != "U0L1X3Q4D" {
				go sender.SendMessage(event.Msg.User, "Sorry, I am under maintenance now")
			} else if event.Msg.User != timebotId {
				go handlers.HandleMessage(&event.Msg)
			}
		}
	}
}

// "user=postgres password=postgres database=timebot_development sslmode=disable"

func connectToDatabase() (*sql.DB, error) {
	db, dbError := sql.Open("postgres", os.Getenv("TIMEBOT_GO_DB_CONNECTION_STRING"))

	if dbError != nil {
		return nil, dbError
	} else {
		return db, nil
	}
}
