package handlers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJenkinsURL(t *testing.T) {
	repo, branch, buildID := parseJenkinsURL("/job/ba/job/bc/1")
	assert.Equal(t, "ba", repo)
	assert.Equal(t, "bc", branch)
	assert.Equal(t, 1, buildID)
}
