package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mgufrone/go-httpclient"
	"github.com/pkg/errors"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
	"net/http"
	"net/url"
	"strings"
)

type Jenkins struct {
	cli httpclient.Client
}

func NewJenkins(cli httpclient.Client) *Jenkins {
	return &Jenkins{cli: cli}
}

func (j *Jenkins) GetPendingAction(jobName string, buildID int) (res contracts.PendingAction, err error) {
	path := fmt.Sprintf(
		"/job/%s/%d/wfapi/pendingInputActions",
		jobName,
		buildID,
	)
	pipeReq, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return
	}
	pipeRes, err := j.cli.Do(pipeReq)
	if err != nil {
		return
	}
	var runResponse []contracts.PendingAction
	defer pipeRes.Body.Close()
	if err = json.NewDecoder(pipeRes.Body).Decode(&runResponse); err != nil {
		return
	}
	if len(runResponse) == 0 {
		err = errors.New("no input provided")
		return
	}
	return runResponse[0], nil
}

func (j *Jenkins) Approve(ctx context.Context, act contracts.PendingAction) (err error) {
	targetURL, _ := url.Parse(act.ProceedURL)
	dt := targetURL.Query()
	dt.Add("proceed", act.ProceedText)
	values := map[string]string{
		"proceed": act.ProceedText,
	}
	by, _ := json.Marshal(values)
	pipeReq, err := http.NewRequest("POST", act.ProceedURL, strings.NewReader(dt.Encode()))
	pipeReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	dt.Set("json", string(by))
	pipeReq.URL = targetURL
	pipeReq.URL.RawQuery = dt.Encode()
	return j.submitAction(ctx, pipeReq)
}

func (j *Jenkins) Reject(ctx context.Context, act contracts.PendingAction) (err error) {
	pipeReq, err := http.NewRequest("POST", act.AbortURL, nil)
	if err == nil {
		return j.submitAction(ctx, pipeReq)
	}
	return
}

func (j *Jenkins) submitAction(ctx context.Context, req *http.Request) (err error) {
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	defer req.Body.Close()
	_, err = j.cli.Do(req)
	return
}
