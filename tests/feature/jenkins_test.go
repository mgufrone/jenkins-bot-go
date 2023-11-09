package feature

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
	"github.com/pkg/errors"
	slack2 "github.com/slack-go/slack"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"mgufrone.dev/jenkins-bot-go/packages/slack"
	"mgufrone.dev/jenkins-bot-go/tests"
	"mgufrone.dev/jenkins-bot-go/tests/mocks"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type JenkinsTestSuite struct {
	suite.Suite
	tests.TestCase
	ht *httptest.Server
}

func TestJenkinsSuite(t *testing.T) {
	suite.Run(t, new(JenkinsTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *JenkinsTestSuite) SetupTest() {
	s.ht = httptest.NewServer(facades.Route())
}

// TearDownTest will run after each test in the suite.
func (s *JenkinsTestSuite) TearDownTest() {
	s.ht = nil
}

func (s *JenkinsTestSuite) TestApprovalFail0() {
	res, _ := http.Post(s.ht.URL+"/approval/jenkins", "application/json", nil)
	defer res.Body.Close()
	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
}
func (s *JenkinsTestSuite) TestApprovalFail1() {
	res, _ := http.Post(s.ht.URL+"/approval/jenkins", "application/json", strings.NewReader(`{}`))
	defer res.Body.Close()
	by, _ := io.ReadAll(res.Body)
	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	s.Require().Contains(string(by), "errors")
}
func (s *JenkinsTestSuite) TestApprovalFail2() {
	res, _ := http.Post(s.ht.URL+"/approval/jenkins", "application/json", strings.NewReader(`{"message":"something"}`))
	defer res.Body.Close()
	by, _ := io.ReadAll(res.Body)
	s.Require().Equal(http.StatusBadRequest, res.StatusCode)
	s.Require().Contains(string(by), "errors")
}
func (s *JenkinsTestSuite) TestApprovalFail4() {
	mHttp := &mocks.MockHttpClient{}
	cli := slack2.New("nothing", slack2.OptionHTTPClient(mHttp))
	rd := strings.NewReader(`{"status":"ok"}`)
	body := io.NopCloser(rd)
	mHttp.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, errors.New("failed to reach out"))
	facades.App().Bind(slack.BindingClient, func(app foundation.Application) (any, error) {
		return cli, nil
	})
	res, _ := http.Post(s.ht.URL+"/approval/jenkins", "application/json", strings.NewReader(`{"message":"something", "build_number":1, "build_name": "connect"}`))
	defer res.Body.Close()
	mHttp.AssertNumberOfCalls(s.T(), "Do", 1)
	s.Require().Equal(http.StatusInternalServerError, res.StatusCode)
}
func (s *JenkinsTestSuite) TestApprovalSuccess0() {
	mHttp := &mocks.MockHttpClient{}
	log := facades.Log()
	cli := slack2.New("nothing", slack2.OptionHTTPClient(mHttp))
	rd := strings.NewReader(`{"status":"ok"}`)
	body := io.NopCloser(rd)
	mHttp.On("Do", mock.Anything).Return(&http.Response{StatusCode: http.StatusOK, Body: body}, nil)
	facades.App().Bind(slack.BindingClient, func(app foundation.Application) (any, error) {
		return cli, nil
	})
	res, _ := http.Post(s.ht.URL+"/approval/jenkins", "application/json", strings.NewReader(`{"message":"something", "build_number":1, "build_name": "connect"}`))
	defer res.Body.Close()
	mHttp.AssertNumberOfCalls(s.T(), "Do", 1)
	by, _ := io.ReadAll(res.Body)
	arg := mHttp.Calls[0].Arguments[0].(*http.Request)
	reqBody, _ := io.ReadAll(arg.Body)
	form, _ := url.ParseQuery(string(reqBody))
	s.Assert().Equal("nothing", form.Get("token"))
	s.Assert().Equal("notification", form.Get("channel"))
	s.Assert().Contains(form.Get("attachments"), "connect:1")
	log.Info(string(by))
	s.Require().Equal(http.StatusOK, res.StatusCode)
}
