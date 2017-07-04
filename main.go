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
	"runtime"
)

var publicChannelsMap map[string]bool
var timebotId string

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)
	sender.Api = api

	var err error

	timebotId, err = getTimebotId(api)

	if err != nil {
		fmt.Println(err)
		return
	}

	models.DB, err = connectToDatabase()

	if err != nil {
		fmt.Println(err)
		return
	}

	publicChannels, err := api.GetChannels(true)

	if err != nil {
		fmt.Println(err)
		return
	}

	publicChannelsMap = make(map[string]bool)

	for _, channel := range publicChannels {
		publicChannelsMap[channel.ID] = true
	}

	startBot(api)
}

func getTimebotId(api *slack.Client) (string, error) {
	users, err := api.GetUsers()

	if err != nil {
		return "", err
	}

	for _, user := range users {
		if user.Name == "timebot" {
			timebotId = user.ID
			break
		}
	}

	return timebotId, nil
}

func startBot(api *slack.Client) {
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	semaphore := make(chan int, runtime.NumCPU())

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack")
		case *slack.MessageEvent:
			if messageIsProcessable(&event.Msg) && underDevelopment(&event.Msg) {
				semaphore <- 1
				go (func() {
					sender.SendMessage(event.Msg.User, "Sorry, I am under maintenance now")
					<- semaphore
				})()
			} else if messageIsProcessable(&event.Msg) {
				semaphore <- 1
				go (func() {
					handlers.HandleMessage(&event.Msg)
					<- semaphore
				})()
			}
		}
	}
}

func messageIsProcessable(msg *slack.Msg) bool {
	return msg.User != timebotId && messageIsNotFromPublicChannel(msg.Channel)
}

func messageIsNotFromPublicChannel(channelId string) bool {
	if _, ok := publicChannelsMap[channelId]; ok {
		return false
	} else {
		return true
	}
}

func underDevelopment(msg *slack.Msg) bool {
	return os.Getenv("GOLANG_ENV") == "development" && msg.User != "U0L1X3Q4D"
}

func connectToDatabase() (*sql.DB, error) {
	db, dbError := sql.Open("postgres", os.Getenv("TIMEBOT_GO_DB_CONNECTION_STRING"))

	if dbError != nil {
		return nil, dbError
	} else {
		return db, nil
	}
}
