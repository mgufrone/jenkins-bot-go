package tests

import (
	"github.com/goravel/framework/testing"

	"mgufrone.dev/job-tracking/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
