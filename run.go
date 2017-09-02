package main

import (
	"github.com/nlopes/slack"
	"log"
)

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
