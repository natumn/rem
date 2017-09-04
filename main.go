package main

import (
	"os"
	"log"

	"github.com/nlopes/slack"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
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
				rtm.SendMessage(rtm.NewOutgoingMessage("ok!", ev.Channel))
			
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		}
	}
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}
	log.Printf("[INFO] Start slack event listening")

	api := slack.New(env.BotToken)
	os.Exit(run(api))
}
