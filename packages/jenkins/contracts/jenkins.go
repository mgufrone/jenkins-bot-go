package contracts

import "context"

type Jenkins interface {
	GetPendingAction(jobName string, buildID int) (PendingAction, error)
	Approve(ctx context.Context, act PendingAction) error
	Reject(ctx context.Context, act PendingAction) error
}

type Handler interface {
	Run(ctx context.Context, payload string) error
}
