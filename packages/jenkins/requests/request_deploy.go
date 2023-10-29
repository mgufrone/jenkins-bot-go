package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type RequestDeploy struct {
	Channel   string `form:"channel" json:"channel"`
	BuildName string `form:"build_name" json:"build_name"`
	BuildID   string `form:"build_id" json:"build_id"`
	Message   string `form:"message" json:"message"`
}

func (r *RequestDeploy) Authorize(ctx http.Context) error {
	return nil
}

func (r *RequestDeploy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"channel":    "required",
		"message":    "required|max_len:255",
		"build_name": "required",
		"build_id":   "required|int",
	}
}

func (r *RequestDeploy) Messages(ctx http.Context) map[string]string {
	return map[string]string{
	}
}

func (r *RequestDeploy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
	}
}

func (r *RequestDeploy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	if val, exist := data.Get("channel"); !exist || val == "" {
		return data.Set("channel", facades.Config().GetString("slack.channel"))
	}
	return nil
}
