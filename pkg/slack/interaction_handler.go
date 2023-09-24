package slack

import (
	"context"
	"github.com/slack-go/slack"
	"go.uber.org/fx"
)

type ISlackInteractionHandler interface {
	When(callback slack.InteractionCallback) bool
	Run(ctx context.Context, callback slack.InteractionCallback) (slack.Message, error)
}

func AsSlackInteractionHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ISlackInteractionHandler)),
		fx.ResultTags(`group:"slack_interaction_handlers"`),
	)

}
