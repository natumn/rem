package main

import (
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
)

type envConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
}

func run(args []string) int {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)

	rtm := client.NewRTM()
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
	os.Exit(run(os.Args[1:]))
}
