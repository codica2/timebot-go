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
	"time"
)

var publicChannelsMap map[string]bool
var timebotId string
var usersMap map[string]string

func main() {
	fmt.Printf("Number of logical processors: %d\n", runtime.GOMAXPROCS(0))
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)
	sender.Api = api

	var err error

	timebotId, usersMap, err = getUsers(api)

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

func getUsers(api *slack.Client) (string, map[string]string, error) {
	users, err := api.GetUsers()

	if err != nil {
		return "", nil, err
	}

	usersMap := make(map[string]string)
	timebotId := ""

	for _, user := range users {
		usersMap[user.ID] = user.Name
		if user.Name == "timebot" {
			timebotId = user.ID
		}
	}

	return timebotId, usersMap, nil
}

func startBot(api *slack.Client) {
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack")
		case *slack.MessageEvent:
			if messageIsProcessable(&event.Msg) && underDevelopment(&event.Msg) {
				go sender.SendMessage(event.Msg.User, "Sorry, I am under maintenance now")
			} else if messageIsProcessable(&event.Msg) {
				go (func() {
					logMessage(&event.Msg)
					handlers.HandleMessage(&event.Msg)
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

func logMessage(msg *slack.Msg) {
	location, _ := time.LoadLocation("Europe/Kiev")
	t := time.Now().In(location).Format("02.01.06 15:04:05")
	fmt.Printf("%s - %s - %s\n", t, usersMap[msg.User], msg.Text)
}
