package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	slack2 "github.com/mgufrone/jenkins-slackbot/pkg/slack"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"net/http"
	"strings"
)

type slackInteraction struct {
	handlers []slack2.ISlackInteractionHandler
	logger   logrus.FieldLogger
}

func (s *slackInteraction) handler(ctx *gin.Context) {
	var (
		callback slack.InteractionCallback
		response slack.Message
		err      error
	)
	py := ctx.PostForm("payload")
	payloadReader := strings.NewReader(py)
	s.logger.Infoln(py)
	if py == "" {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{})
		return
	}
	if err = json.NewDecoder(payloadReader).Decode(&callback); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	//_ = json.NewEncoder(os.Stdout).Encode(callback)
	//s.logger.Infoln("received payload", callback)
	for _, v := range s.handlers {
		if v.When(callback) {
			response, err = v.Run(ctx, callback)
			break
		}
	}
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (s *slackInteraction) Mount(srv gin.IRouter) {
	grp := srv.Group("/slack/interaction")
	{
		grp.POST("/", s.handler)
	}
}

func slackInteractionHandler(handlers []slack2.ISlackInteractionHandler, logger logrus.FieldLogger) *slackInteraction {
	return &slackInteraction{
		handlers: handlers,
		logger:   logger,
	}
}
