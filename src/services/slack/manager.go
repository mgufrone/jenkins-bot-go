package slack

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/fx"
)

type SocketManager struct {
	web      *slack.Client
	socket   *socketmode.Client
	logger   logrus.FieldLogger
	handlers []ISocketSubscriber
}

func NewSocketManager(web *slack.Client, socket *socketmode.Client, handlers []ISocketSubscriber, logger logrus.FieldLogger, lc fx.Lifecycle) *SocketManager {
	mgr := &SocketManager{web: web, socket: socket, handlers: handlers, logger: logger}
	lc.Append(fx.StartHook(func() error {
		go mgr.Loop()
		return mgr.socket.RunContext(context.TODO())
	}))
	return mgr
}

func (s *SocketManager) Loop() {
	s.logger.Infoln("start listening incoming events")
	for event := range s.socket.Events {
		for _, hndl := range s.handlers {
			go func(handler ISocketSubscriber, evt socketmode.Event) {
				s.logger.Debugf("check if %s can run the event %s", handler.Name(), evt.Type)
				if !handler.When(evt) {
					return
				}
				_ = handler.Run(s.web, s.socket, evt)
			}(hndl, event)
		}
	}
}
