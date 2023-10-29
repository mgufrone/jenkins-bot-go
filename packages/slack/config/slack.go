package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("slack", map[string]any{
		"app_token": config.Env("SLACK_APP_TOKEN"),
		"bot_token": config.Env("SLACK_BOT_TOKEN"),
		"channel":   config.Env("SLACK_DEFAULT_CHANNEL", "notification"),
	})
}
