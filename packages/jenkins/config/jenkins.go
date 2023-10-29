package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("jenkins", map[string]any{
		"url":             config.Env("JENKINS_URL", "http://jenkins.localhost"),
		"username":        config.Env("JENKINS_USERNAME", ""),
		"token":           config.Env("JENKINS_USER_API_TOKEN", ""),
		"acknowledgeText": ":gh-loading: Submitting action",
	})
}
