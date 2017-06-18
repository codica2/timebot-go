package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nlopes/slack"
	"os"
	"github.com/alex-bogomolov/timebot_go/timebot"
)

type User struct {
	Id   int
	Name int
}

func main() {
	slackToken := os.Getenv("SLACK_TOKEN")

	api := slack.New(slackToken)
	timebotId, users, err := getUsers(api)

	if err != nil {
		fmt.Println(err)
		return
	}

	db, err := connectToDatabase()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db)

	startBot(api, timebotId, users)
}

func getUsers(api *slack.Client) (string, map[string]string, error) {
	users, err := api.GetUsers()

	if err != nil {
		return "", nil, err
	}

	var timebotId string
	usersMap := map[string]string{}

	if err == nil {
		for _, v := range users {
			if v.Name == "timebot" {
				timebotId = v.ID
			}

			usersMap[v.ID] = v.Name
		}
	}

	return timebotId, usersMap, nil
}

func connectToDatabase() (*sql.DB, error) {
	db, dbError := sql.Open("postgres", "user=postgres password=postgres database=timebot_development sslmode=disable")

	if dbError != nil {
		return nil, dbError
	} else {
		return db, nil
	}
}

func startBot(api *slack.Client, timebotId string, users map[string]string) {
	rtm := api.NewRTM()

	outgoingMessages := make(chan *timebot.OutgoingMessage)

	go rtm.ManageConnection()
	go ListenToOutgoingMessages(outgoingMessages, api)

	for msg := range rtm.IncomingEvents {

		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			fmt.Println("Connected to Slack")
		case *slack.MessageEvent:
			if ev.Msg.User == timebotId {
				break
			}

			timebot.HandleMessage(ev, outgoingMessages, users)

		}
	}
}

func ListenToOutgoingMessages(channel chan *timebot.OutgoingMessage, api *slack.Client) {
	for msg := range channel {
		api.SendMessage(msg.User, slack.MsgOptionText(msg.Text, true), slack.MsgOptionAsUser(true))
	}
}
