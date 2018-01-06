package main

import (
	"log"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	gentleman "gopkg.in/h2non/gentleman.v2"
)

type envConfig struct {
	BotToken string `envconfig:"BOT_TOKEN" required:"true"`
}

type workingTime struct {
	startWorking  time.Time
	finishWorking time.Time
	weekTotal     time.Duration
}

func run(args []string) int {
	cli := gentleman.New()
	var env envConfig
	var wt workingTime
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}

	log.Printf("[INFO] Start slack event listening")
	slackCli := slack.New(env.BotToken)

	rtm := slackCli.NewRTM()
	go rtm.ManageConnection()
	weekAlert := time.NewTicker(168 * time.Hour)
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("[INFO]working bot...")
			case *slack.MessageEvent:
				log.Print("[INFO]get message")
				if strings.Contains(ev.Text, "出勤ed") {
					setToggl()
					rtm.SendMessage(rtm.NewOutgoingMessage("出勤を確認しました。\n今日もお仕事頑張ってください！", ev.Channel))
				}
				if strings.Contains(ev.Text, "退勤ed") {
					// endToggl()
					rtm.SendMessage(rtm.NewOutgoingMessage("退勤を確認しました。\n今日働いた時間は です!\nお疲れ様でした。", ev.Channel))
				}
				if ev.Text == "todo" {
					// ttl := getTodayTodo()
					rtm.SendMessage(rtm.NewOutgoingMessage("今日のtodoは", ev.Channel))
				}
				if ev.Text == "今日のプログラミングの時間" {

				}
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		case <-weekAlert.C:
			log.Println("[INFO]weekend alert")
			rtm.SendMessage(rtm.NewOutgoingMessage("今週もお疲れ様です。", "test_natsumen"))
		}
	}
}
