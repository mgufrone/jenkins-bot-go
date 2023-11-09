package tests

import (
	"github.com/goravel/framework/testing"

	"mgufrone.dev/jenkins-bot-go/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
