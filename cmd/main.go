package main

import (
	"fmt"
	"log"
	"os"

	"github.com/etherlabsio/reviewit"
	"github.com/etherlabsio/reviewit/provider/gitlab"
	"github.com/etherlabsio/reviewit/transport/slack"
	"github.com/spf13/pflag"
)

func main() {
	var (
		providerToken = pflag.String("gitlab.token", "", "Version control provider access token")
		group         = pflag.String("gitlab.group", "", "Name of the group to search projects under gitlab, github etc.")
		slackWebhook  = pflag.String("slack.webhook", "", "Webhook URL for slack's bot")
		slackChannel  = pflag.String("slack.channel", "", "Name of the slack channel to publish")
	)
	pflag.Parse()

	pflag.VisitAll(func(f *pflag.Flag) {
		if f.DefValue == "" {
			fmt.Println(f.Name + " is not set. \n")
			os.Exit(1)
		}
	})

	provider := gitlab.NewGroupFilter(*group, *providerToken)
	transport, err := slack.New(*slackChannel, *slackWebhook)
	if err != nil {
		log.Fatal(err)
	}
	pending, err := provider.Filter()
	if err != nil {
		log.Fatal(err)
	}
	transport.Send(pending)
}

// Runner runs the application
type Runner struct {
	filter reviewit.Filterer
	sender reviewit.Sender
	err    error
}
