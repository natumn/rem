package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"log"

	"github.com/nlopes/slack"
)

// Config is slack token data.
type Config struct {
	Token string `json:"token"`
}

func run(api *slack.Client) int {
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("hello, event")
			case *slack.MessageEvent:
				log.Print("message, event")
				rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ev.Channel))
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func main() {
	bytes, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	api := slack.New(config.Token)
	os.Exit(run(api))
}
