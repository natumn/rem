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

type workingTime struct {
	startWorking  time.Time
	finishWorking time.Time
	weekTotal time.Duration
}

func run(args []string) int {
	var env envConfig
	var wt workingTime
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		return 1
	}
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)

	rtm := client.NewRTM()
	go rtm.ManageConnection()
	t := time.NewTicker(10 * time.Second)
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Print("[INFO] ")
			case *slack.MessageEvent:
				log.Print("[INFO]get message")
				if ev.Text == "出勤" {
					wt.startWorking = time.Now()
					log.Println("[INFO]start working:", wt.startWorking)
					rtm.SendMessage(rtm.NewOutgoingMessage("出勤を確認しました。\n今日もお仕事頑張ってください！", ev.Channel))
				}
				if ev.Text == "退勤" {
					wt.finishWorking = time.Now()
					log.Println("[INFO]finish working:", wt.finishWorking)
					totalTime := wt.finishWorking.Sub(wt.startWorking)
					log.Println("[INFO]total time:",totalTime.String())
					rtm.SendMessage(rtm.NewOutgoingMessage("退勤を確認しました。\n今日働いた時間は"+totalTime.String()+"です。\nお疲れ様でした！", ev.Channel))
					wt.weekTotal += totalTime
				}
			case *slack.InvalidAuthEvent:
				log.Print("Invalid credentials")
				return 1
			}
		case <-t.C:
			log.Println("[INFO]working time.ticker 10s...")
		}
	}
}

func main() {
	os.Exit(run(os.Args[1:]))
}
