package jenkins

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/route"
	"github.com/mgufrone/go-httpclient"
	"github.com/mgufrone/go-httpclient/interceptor"
	"mgufrone.dev/job-tracking/packages/jenkins/contracts"
	"mgufrone.dev/job-tracking/packages/jenkins/handlers"
	"net/http"
	"net/url"
)

const Binding = "jenkins.service"
const BindingHandler = "jenkins.handler"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		logger := app.MakeLog()
		cfg := app.MakeConfig()
		cli := httpclient.Standard().
			AddInterceptor(func(req *http.Request) {
				go func() {
					logger.Debug("sending request to", req.Method, req.URL.String())
				}()
			}).
			AddInterceptor(interceptor.Header("Content-Type", "application/json")).
			AddInterceptor(func(req *http.Request) {
				parsedUrl, _ := url.Parse(cfg.GetString("jenkins.url"))
				req.URL.Host = parsedUrl.Host
				req.URL.Scheme = parsedUrl.Scheme
				req.SetBasicAuth(cfg.GetString("jenkins.username"), cfg.GetString("jenkins.token"))
			})
		return NewJenkins(cli), nil
	})
	app.Bind(BindingHandler, func(app foundation.Application) (any, error) {
		svc, _ := app.Make(Binding)
		logger := app.MakeLog()
		return handlers.NewJenkins(svc.(contracts.Jenkins), logger), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	rt := app.MakeRoute()
	app.Publishes("mgufrone.dev/job-tracking/packages/jenkins", map[string]string{
		"config/jenkins.go": app.ConfigPath("jenkins.go"),
	})
	rt.Prefix("/approval").Group(func(router route.Router) {
		requestApproval := handlers.NewRequestApproval()
		router.Post("/jenkins", requestApproval.SendApproval)
	})
}
