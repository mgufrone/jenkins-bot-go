package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
	"strconv"
)

type RequestDeploy struct {
	Channel     string `form:"channel" json:"channel"`
	BuildName   string `form:"build_name" json:"build_name"`
	BuildNumber int    `form:"build_number" json:"build_number"`
	Text        string `form:"message" json:"message"`
}

func (r *RequestDeploy) Authorize(ctx http.Context) error {
	return nil
}

func (r *RequestDeploy) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"channel":      "required",
		"message":      "required|max_len:255",
		"build_name":   "required",
		"build_number": "required|int",
	}
}

func (r *RequestDeploy) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RequestDeploy) Attributes(ctx http.Context) map[string]string {
	return map[string]string{
		"build_name":   "build name",
		"build_number": "build id/number",
	}
}

func (r *RequestDeploy) PrepareForValidation(ctx http.Context, data validation.Data) error {
	if val, exists := data.Get("build_number"); exists {
		switch v := val.(type) {
		case string:
			value, _ := strconv.Atoi(v)
			_ = data.Set("build_number", value)
		case float64:
			_ = data.Set("build_number", int(v))
		case float32:
			_ = data.Set("build_number", int(v))
		case int64:
			_ = data.Set("build_number", int(v))
		}
	}
	if val, exist := data.Get("channel"); !exist || val == "" {
		return data.Set("channel", facades.Config().GetString("slack.channel"))
	}
	return nil
}
