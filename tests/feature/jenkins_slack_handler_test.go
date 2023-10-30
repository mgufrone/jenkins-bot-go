package feature

import (
	"context"
	"encoding/json"
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	mock2 "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"mgufrone.dev/job-tracking/packages/jenkins"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
	facades2 "mgufrone.dev/job-tracking/packages/jenkins/facades"
	"mgufrone.dev/job-tracking/packages/jenkins/handlers"
	"mgufrone.dev/job-tracking/tests"
	"testing"
)

type mockJenkins struct {
	mock2.Mock
}

func (m *mockJenkins) GetPendingAction(jobName string, buildID int) (contracts.PendingAction, error) {
	vals := m.Called(jobName, buildID)
	return vals.Get(0).(contracts.PendingAction), vals.Error(1)
}

func (m *mockJenkins) Approve(ctx context.Context, act contracts.PendingAction) error {
	vals := m.Called(ctx, act)
	return vals.Error(0)
}

func (m *mockJenkins) Reject(ctx context.Context, act contracts.PendingAction) error {
	vals := m.Called(ctx, act)
	return vals.Error(0)
}

type mockInteractionSubmitted struct {
	mock2.Mock
}

func (m *mockInteractionSubmitted) Signature() string {
	return "on_mock_interaction"
}

func (m *mockInteractionSubmitted) Queue(args ...any) event.Queue {
	return m.Called(args...).Get(0).(event.Queue)
}

func (m *mockInteractionSubmitted) Handle(args ...any) error {
	return m.Called(args...).Error(0)
}

type JenkinsSlackHandlerSuite struct {
	suite.Suite
	tests.TestCase
	svc *mockJenkins
}

func TestJenkinsHandlerSuite(t *testing.T) {
	suite.Run(t, new(JenkinsSlackHandlerSuite))
}

// SetupTest will run before each test in the suite.
func (s *JenkinsSlackHandlerSuite) SetupTest() {
	//s.ht = httptest.NewServer(facades.Route())
	//s.evtMock, s.task = mock.Event()
	s.svc = &mockJenkins{}
	facades.App().Bind(jenkins.BindingHandler, func(app foundation.Application) (any, error) {
		return handlers.NewJenkins(s.svc, facades.Log()), nil
	})
}

// TearDownTest will run after each test in the suite.
func (s *JenkinsSlackHandlerSuite) TearDownTest() {
	s.svc = nil
}

func (s *JenkinsSlackHandlerSuite) TestFail0() {
	handler := facades2.Handler()
	err := handler.Run(context.TODO(), "")
	s.Assert().NotNil(err)
}
func (s *JenkinsSlackHandlerSuite) TestFail1() {
	handler := facades2.Handler()
	err := handler.Run(context.TODO(), "{}")
	s.Assert().Nil(err)
}
func (s *JenkinsSlackHandlerSuite) TestFail2() {
	s.svc.On("GetPendingAction", mock2.AnythingOfType("string"), mock2.AnythingOfType("int")).Return(contracts.PendingAction{}, errors.New("failed to retrieve"))
	handler := facades2.Handler()
	var callback slack.InteractionCallback
	callback.Type = slack.InteractionTypeInteractionMessage
	callback.CallbackID = "devops/random:1"
	callback.ActionCallback.AttachmentActions = []*slack.AttachmentAction{
		{
			Value: "abort",
		},
	}
	by, _ := json.Marshal(callback)
	err := handler.Run(context.TODO(), string(by))
	s.Assert().NotNil(err)
}
func (s *JenkinsSlackHandlerSuite) TestFail3() {
	facades.Event().Register(map[event.Event][]event.Listener{})
	s.svc.On("GetPendingAction", mock2.AnythingOfType("string"), mock2.AnythingOfType("int")).Return(contracts.PendingAction{
		AbortURL: "http://localhost/somewhere",
	}, nil)
	s.svc.On("Approve", mock2.Anything, mock2.Anything).Return(errors.New("failed to trigger"))
	handler := facades2.Handler()
	var callback slack.InteractionCallback
	callback.Type = slack.InteractionTypeInteractionMessage
	callback.CallbackID = "devops/random:1"
	callback.ActionCallback.AttachmentActions = []*slack.AttachmentAction{
		{
			Value: "approve",
		},
	}
	by, _ := json.Marshal(callback)
	err := handler.Run(context.TODO(), string(by))
	s.Assert().NotNil(err)
}
func (s *JenkinsSlackHandlerSuite) TestFail4() {
	//facades.Event().Register(map[event.Event][]event.Listener{})
	s.svc.On("GetPendingAction", mock2.AnythingOfType("string"), mock2.AnythingOfType("int")).Return(contracts.PendingAction{
		AbortURL: "http://localhost/somewhere",
	}, nil)
	s.svc.On("Reject", mock2.Anything, mock2.Anything).Return(errors.New("failed to trigger"))
	handler := facades2.Handler()
	var callback slack.InteractionCallback
	callback.Type = slack.InteractionTypeInteractionMessage
	callback.CallbackID = "devops/random:1"
	callback.ActionCallback.AttachmentActions = []*slack.AttachmentAction{
		{
			Value: "abort",
		},
	}
	by, _ := json.Marshal(callback)
	err := handler.Run(context.TODO(), string(by))
	s.Assert().NotNil(err)
}
func (s *JenkinsSlackHandlerSuite) TestSuccessApprove() {
	//facades.Event().Register(map[event.Event][]event.Listener{})
	s.svc.On("GetPendingAction", mock2.AnythingOfType("string"), mock2.AnythingOfType("int")).Return(contracts.PendingAction{
		AbortURL: "http://localhost/somewhere",
	}, nil)
	s.svc.On("Approve", mock2.Anything, mock2.Anything).Return(nil)
	handler := facades2.Handler()
	var callback slack.InteractionCallback
	callback.Type = slack.InteractionTypeInteractionMessage
	callback.CallbackID = "devops/random:1"
	callback.ActionCallback.AttachmentActions = []*slack.AttachmentAction{
		{
			Value: "approve",
		},
	}
	by, _ := json.Marshal(callback)
	err := handler.Run(context.TODO(), string(by))
	s.Assert().NotNil(err)
}
