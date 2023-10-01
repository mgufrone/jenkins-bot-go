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

type ManagerParams struct {
	fx.In
	Slck     *slack.Client
	Socket   *socketmode.Client
	Handlers []ISocketSubscriber `group:"socket_subscribers"`
	Logger   logrus.FieldLogger
}

func NewSocketManager(param ManagerParams, lc fx.Lifecycle) *SocketManager {
	mgr := &SocketManager{web: param.Slck, socket: param.Socket, handlers: param.Handlers, logger: param.Logger}
	lc.Append(fx.StartHook(func() error {
		go mgr.Loop()
		go func() {
			_ = mgr.socket.RunContext(context.TODO())
		}()
		return nil
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
				s.logger.Debugf("invoke %s on %s", handler.Name(), evt.Type)
				_ = handler.Run(s.web, s.socket, evt)
			}(hndl, event)
		}
	}
}
