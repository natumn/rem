package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

type Config struct {
	token string `json:"token"`
}

func main() {
	api := slack.New(Config.token)
	os.Exit(run(api))
}
