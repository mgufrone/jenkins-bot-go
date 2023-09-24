package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
)

type slackOauth struct {
	api    *slack.Client
	logger logrus.FieldLogger
}

func (s *slackOauth) Mount(srv gin.IRouter) {
	grp := srv.Group("/slack")
	{
		grp.GET("/redirect", s.redirectHandler)
	}
}

func (s *slackOauth) redirectHandler(ctx *gin.Context) {
	redirectURI := "https://" + ctx.Request.Host + ctx.Request.URL.Path
	s.logger.Infoln("redirect URI", redirectURI)
	resp, err := slack.GetOAuthV2ResponseContext(ctx, &http.Client{}, env.Get("SLACK_CLIENT_ID"), env.Get("SLACK_CLIENT_SECRET"), ctx.Query("code"), redirectURI)
	s.logger.Infoln("token", resp.AccessToken)
	if err != nil {
		s.logger.Errorln(err)
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token":  "ok",
		"scopes": resp.Scope,
		"data":   resp.IncomingWebhook,
	})
}

func newSlackOauthHandler(api *slack.Client, logger logrus.FieldLogger) *slackOauth {
	return &slackOauth{api: api, logger: logger}
}
