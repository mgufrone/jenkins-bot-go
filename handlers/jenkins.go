package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type CrumbResponse struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

type jenkinsClient struct {
	host, repo, branch, username, password string
	buildID                                int
	client                                 *http.Client
	crumbRes                               *CrumbResponse
	keepCrumb                              bool
}

func newJenkinsClient(host, repo, branch string, buildID int) *jenkinsClient {
	jar, _ := cookiejar.New(nil)
	return &jenkinsClient{
		host: host, repo: repo, branch: branch,
		username: os.Getenv("JENKINS_USERNAME"), password: os.Getenv("JENKINS_PASSWORD"),
		buildID: buildID,
		client:  &http.Client{Jar: jar},
	}
}

func (cli *jenkinsClient) request(req *http.Request) (*http.Response, error) {
	timer, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	req = req.WithContext(timer)
	req.URL.Host = cli.host
	req.URL.Scheme = "https"
	req.SetBasicAuth(cli.username, cli.password)
	log.Println("requesting: ", req.URL.String())
	res, err := cli.client.Do(req)
	if err == nil {
		log.Println("response", res.StatusCode)
	}
	return res, err
}
func (cli *jenkinsClient) do(req *http.Request) (res *http.Response, err error) {
	if cli.crumbRes == nil {
		cli.crumbRes, err = cli.crumb()
		if err != nil {
			return nil, err
		}
	}
	req.Header = http.Header{
		cli.crumbRes.CrumbRequestField: []string{cli.crumbRes.Crumb},
		"Content-Type":                 []string{"application/json"},
	}
	req.URL.Host = cli.host
	req.URL.Scheme = "https"
	req.SetBasicAuth(cli.username, cli.password)
	res, err = cli.client.Do(req)
	if !cli.keepCrumb {
		cli.crumbRes = nil
	}
	return
}
func (cli *jenkinsClient) crumb() (res *CrumbResponse, err error) {
	crumbReq, err := http.NewRequest("GET", "/crumbIssuer/api/json", nil)
	if err != nil {
		return
	}
	response, err := cli.request(crumbReq)
	if err != nil {
		return
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&res)
	return
}

func (cli *jenkinsClient) getPipelineInputID() (string, string, error) {
	path := fmt.Sprintf(
		"/blue/rest/organizations/jenkins/pipelines/%s/branches/%s/runs/%d/steps",
		cli.repo,
		cli.branch,
		cli.buildID,
	)
	pipeReq, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", "", err
	}
	pipeRes, err := cli.request(pipeReq)
	if err != nil {
		return "", "", err
	}
	var steps []Step
	defer pipeRes.Body.Close()
	if err = json.NewDecoder(pipeRes.Body).Decode(&steps); err != nil {
		return "", "", err
	}
	step := steps[len(steps)-1]
	input := steps[len(steps)-1].Input
	if input == nil {
		return "", "", errors.New("no input provided")
	}
	return step.ID, input.ID, nil
}

func (cli *jenkinsClient) submitInput(action string) error {
	stepID, inputID, err := cli.getPipelineInputID()
	if err != nil {
		return err
	}
	path := fmt.Sprintf(
		"/blue/rest/organizations/jenkins/pipelines/%s/branches/%s/runs/%d/steps/%s/",
		cli.repo,
		cli.branch,
		cli.buildID,
		stepID,
	)
	payload := map[string]interface{}{
		"id":         inputID,
		"parameters": []string{},
	}
	if action == "reject" {
		payload["abort"] = true
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	pipeReq, err := http.NewRequest("POST", path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	_, err = cli.do(pipeReq)
	return err
	// fallback to curl. go automatically lowercase header name
}

func NewJenkinsHandler() *EventRegistrar {
	return &EventRegistrar{
		EventType:       socketmode.EventTypeInteractive,
		InteractionType: slack.InteractionTypeBlockActions,
		Handler:         jenkinsHandler,
	}
}

func parseJenkinsURL(str string) (repo, branch string, buildID int) {
	rep := strings.Replace(str, "/job", "", -1)
	fmt.Println("rep", rep)
	strArr := strings.Split(rep, "/")
	repo = strArr[1]
	branch = strArr[2]
	buildID, _ = strconv.Atoi(strArr[3])
	return
}

func jenkinsHandler(webClient *slack.Client, client *socketmode.Client, req *socketmode.Request) error {
	by, _ := req.Payload.MarshalJSON()
	var pl slack.InteractionCallback
	if err := json.Unmarshal(by, &pl); err != nil {
		return err
	}
	splits := strings.SplitAfterN(pl.ActionCallback.BlockActions[0].Value, ":", 2)
	action := strings.ReplaceAll(splits[0], ":", "")
	jenkinsURI, _ := url.Parse(splits[1])
	repo, branch, buildID := parseJenkinsURL(jenkinsURI.Path)
	cli := newJenkinsClient(jenkinsURI.Host, repo, branch, buildID)
	err := cli.submitInput(action)
	client.Ack(*req)
	if err != nil {
		return err
	}
	blocks := pl.Message.Blocks.BlockSet
	blck := slack.NewSectionBlock(&slack.TextBlockObject{
		Type: "mrkdwn",
		Text: fmt.Sprintf("_%s by <@%s>_", action, pl.User.ID),
		Emoji:    false,
		Verbatim: false,
	}, nil, nil)
	blocks[len(blocks)-1] = blck
	msgID, blockID, replID, err := webClient.UpdateMessage(pl.Channel.ID, pl.Message.Timestamp, slack.MsgOptionBlocks(pl.Message.Blocks.BlockSet...))
	log.Println("replaced", msgID, blockID, replID, err)
	return err
	// we will limit to only multi branch pipeline for now
	// get latest path
}
