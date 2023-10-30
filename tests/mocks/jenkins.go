package mocks

import (
	"context"
	mock2 "github.com/stretchr/testify/mock"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
)

type MockJenkins struct {
	mock2.Mock
}

func (m *MockJenkins) GetPendingAction(jobName string, buildID int) (contracts.PendingAction, error) {
	vals := m.Called(jobName, buildID)
	return vals.Get(0).(contracts.PendingAction), vals.Error(1)
}

func (m *MockJenkins) Approve(ctx context.Context, act contracts.PendingAction) error {
	vals := m.Called(ctx, act)
	return vals.Error(0)
}

func (m *MockJenkins) Reject(ctx context.Context, act contracts.PendingAction) error {
	vals := m.Called(ctx, act)
	return vals.Error(0)
}
