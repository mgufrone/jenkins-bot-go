package facades

import (
	slack2 "github.com/slack-go/slack"
	"log"

	"mgufrone.dev/jenkins-bot-go/packages/slack"
	"mgufrone.dev/jenkins-bot-go/packages/slack/contracts"
)

func Slack() contracts.Slack {
	instance, err := slack.App.Make(slack.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Slack)
}

func SlackClient() *slack2.Client {
	instance, err := slack.App.Make(slack.BindingClient)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(*slack2.Client)
}
