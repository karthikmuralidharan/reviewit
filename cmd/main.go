package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"

	"github.com/etherlabsio/reviewit"
	"github.com/etherlabsio/reviewit/provider/gitlab"
	"github.com/etherlabsio/reviewit/transport/slack"
)

func main() {
	viper.AutomaticEnv()

	var (
		providerToken = viper.GetString("GITLAB_TOKEN")
		group         = viper.GetString("GITLAB_GROUP")
		slackWebhook  = viper.GetString("SLACK_WEBHOOK_URL")
		slackChannel  = viper.GetString("SLACK_CHANNEL")
	)

	requiredInput := []string{"GITLAB_TOKEN", "GITLAB_GROUP", "SLACK_WEBHOOK_URL", "SLACK_CHANNEL"}
	for _, k := range requiredInput {
		if !viper.IsSet(k) {
			fmt.Println(k + " is not set")
			os.Exit(1)
		}
	}

	provider := gitlab.NewGroupFilter(group, providerToken)
	transport, err := slack.New(slackChannel, slackWebhook)
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
