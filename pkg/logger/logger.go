package logger

import (
	"github.com/mgufrone/jenkins-slackbot/pkg/env"
	log "github.com/sirupsen/logrus"
	"os"
)

type stderrHook struct {
	logger *log.Logger
}

func (h *stderrHook) Levels() []log.Level {
	return []log.Level{log.PanicLevel, log.FatalLevel, log.ErrorLevel}
}

func (h *stderrHook) Fire(entry *log.Entry) error {
	entry.Logger = h.logger
	return nil
}
func parseLogLevel(level string) log.Level {
	switch level {
	case "info":
		return log.InfoLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "warn", "warning":
		return log.WarnLevel
	case "trace":
		return log.TraceLevel
	default:
		return log.DebugLevel
	}
}

func Default() log.FieldLogger {
	return NewLogger(
		env.GetOr("APP_NAME", "app"),
		env.GetOr("APP_ENV", "dev"),
		env.GetOr("APP_LOG_LEVEL", "info"),
	)
}
func NewLogger(appName, env, level string) *log.Entry {
	logger := log.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)
	errLogger := log.StandardLogger()
	errLogger.SetOutput(os.Stdout)
	errLogger.SetReportCaller(false)
	logger.AddHook(&stderrHook{logger: errLogger})
	logger.SetFormatter(&log.TextFormatter{})
	logLevel := parseLogLevel(level)
	logger.SetLevel(logLevel)
	return logger.WithFields(log.Fields{
		"app_name": appName,
		"env":      env,
	})
}
func WithAppLogger(appName, env, level string) func() *log.Entry {
	return func() *log.Entry {
		return NewLogger(appName, env, level)
	}
}
