package main

import (
	"fmt"

	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("token")
	os.Exit(run(api))
}
