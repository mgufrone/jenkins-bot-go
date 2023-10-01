package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mgufrone/go-httpclient"
	"github.com/mgufrone/go-httpclient/interceptor"
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

type Client struct {
	cli httpclient.Client
}

func NewClient(logger logrus.FieldLogger) *Client {
	cli := httpclient.Standard().
		AddInterceptor(func(req *http.Request) {
			go func() {
				logger.Debugln("sending request to", req.Method, req.URL.String())
			}()
		}).
		AddInterceptor(interceptor.Header("Content-Type", "application/json")).
		AddInterceptor(func(req *http.Request) {
			parsedUrl, _ := url.Parse(env.Get("JENKINS_URL"))
			req.URL.Host = parsedUrl.Host
			req.URL.Scheme = parsedUrl.Scheme
			req.SetBasicAuth(env.Get("JENKINS_USERNAME"), env.Get("JENKINS_PASSWORD"))
		})
	return &Client{
		cli,
	}
}

func (cli *Client) GetPendingAction(jobName string, buildID int) (*PendingAction, error) {
	path := fmt.Sprintf(
		"/job/%s/%d/wfapi/pendingInputActions",
		jobName,
		buildID,
	)
	pipeReq, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	pipeRes, err := cli.cli.Do(pipeReq)
	if err != nil {
		return nil, err
	}
	var runResponse []*PendingAction
	defer pipeRes.Body.Close()
	if err = json.NewDecoder(pipeRes.Body).Decode(&runResponse); err != nil {
		return nil, err
	}
	if len(runResponse) == 0 {
		return nil, errors.New("no input provided")
	}
	return runResponse[0], nil
}

func (cli *Client) Approve(ctx context.Context, act *PendingAction, submitter string) (err error) {
	targetURL, _ := url.Parse(act.ProceedURL)
	dt := targetURL.Query()
	dt.Add("proceed", act.ProceedText)
	dt.Add("json", fmt.Sprintf(`{"parameter":[{"name":"SLACK_SUBMITTER", "value":"%s"}]}`, submitter))
	targetURL.RawQuery = dt.Encode()
	pipeReq, err := http.NewRequest("POST", act.ProceedURL, nil)
	if err != nil {
		return err
	}
	pipeReq.URL = targetURL
	return cli.submitAction(ctx, pipeReq)
}
func (cli *Client) Abort(ctx context.Context, act *PendingAction) (err error) {
	pipeReq, err := http.NewRequest("POST", act.AbortURL, nil)
	if err == nil {
		return cli.submitAction(ctx, pipeReq)
	}
	return
}
func (cli *Client) submitAction(ctx context.Context, req *http.Request) (err error) {
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	_, err = cli.cli.Do(req)
	return
}
