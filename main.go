package main

import (
	"log"
	"os"
	"time"

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
				log.Print("[INFO] ")
			case *slack.MessageEvent:
				log.Print("[INFO]get message")
				if ev.Text == "出勤" {
					log.Print("[INFO]get going to work event")
					st := time.Now()
					rtm.SendMessage(rtm.NewOutgoingMessage("出勤を確認しました。\n今日もお仕事頑張ってください！",ev.Channel))
				}

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
