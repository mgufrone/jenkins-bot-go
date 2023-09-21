package env

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in     string
		before func()
		out    string
	}{
		{"", nil, ""},
		{"set", nil, ""},
		{"set", func() {
			os.Setenv("set", "something")
		}, "something"},
	}
	for _, c := range testCases {
		if c.before != nil {
			c.before()
		}
		res := Get(c.in)
		assert.Equal(t, c.out, res)
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in     string
		before func()
		def    string
		out    string
	}{
		{"", nil, "", ""},
		{"k", nil, "", ""},
		{"k", nil, "v", "v"},
		{"k", func() {
			os.Setenv("k", "c")
		}, "v", "c"},
	}
	for _, c := range testCases {
		if c.before != nil {
			c.before()
		}
		Default(c.in, c.def)
		assert.Equal(t, c.out, Get(c.in))
	}
}

func TestRequires(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in     string
		before func()
		fail   bool
	}{
		{"", nil, false},
		{"k1", nil, true},
		{"k1", func() {
			Default("k1", "Test")
		}, false},
	}
	for _, c := range testCases {
		if c.before != nil {
			c.before()
		}
		err := Requires(c.in)
		if c.fail {
			assert.NotNil(t, err)
			continue
		}
		assert.Nil(t, err)
	}
}

func TestGetOr(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		in     string
		before func()
		def    string
		out    string
	}{
		{"", nil, "", ""},
		{"k2", nil, "", ""},
		{"k2", nil, "v", "v"},
		{"k2", func() {
			os.Setenv("k2", "c")
		}, "v", "c"},
	}
	for _, c := range testCases {
		if c.before != nil {
			c.before()
		}
		Default(c.in, c.def)
		assert.Equal(t, c.out, Get(c.in))
	}
}
