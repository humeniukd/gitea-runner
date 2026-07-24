// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"os"
	"testing"

	"gitea.com/gitea/runner/internal/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestJobProxyEnv(t *testing.T) {
	t.Run("sets both spellings of every variable", func(t *testing.T) {
		env := jobProxyEnv(config.Proxy{
			HTTPProxy:  "http://proxy:3128",
			HTTPSProxy: "http://proxy:3128",
			NoProxy:    "example.com",
		})

		assert.Equal(t, map[string]string{
			"http_proxy":  "http://proxy:3128",
			"HTTP_PROXY":  "http://proxy:3128",
			"https_proxy": "http://proxy:3128",
			"HTTPS_PROXY": "http://proxy:3128",
			"no_proxy":    "example.com",
			"NO_PROXY":    "example.com",
		}, env)
	})

	// A job container reaching the runner's cache server on the local network must not be
	// sent through an internet proxy.
	t.Run("adds the runner's own servers to no_proxy", func(t *testing.T) {
		env := jobProxyEnv(
			config.Proxy{HTTPProxy: "http://proxy:3128", NoProxy: "example.com"},
			"http://192.168.1.10:34567/",
			"http://192.168.1.10:34568",
		)

		assert.Equal(t, "example.com,192.168.1.10", env["no_proxy"])
		assert.Equal(t, "example.com,192.168.1.10", env["NO_PROXY"])
	})

	t.Run("adds the runner's servers even without a configured no_proxy", func(t *testing.T) {
		env := jobProxyEnv(config.Proxy{HTTPSProxy: "http://proxy:3128"}, "http://cache.local:8088/")

		assert.Equal(t, "cache.local", env["no_proxy"])
		assert.NotContains(t, env, "http_proxy")
	})

	t.Run("is empty when no proxy is configured", func(t *testing.T) {
		assert.Nil(t, jobProxyEnv(config.Proxy{}, "http://cache.local:8088/"))
	})
}

func TestNoProxyWithRunnerHosts(t *testing.T) {
	assert.Equal(t, "a.example,b.example",
		noProxyWithRunnerHosts(" a.example , b.example "))
	// A host already listed is not repeated.
	assert.Equal(t, "cache.local",
		noProxyWithRunnerHosts("cache.local", "http://cache.local:8088/"))
	// Unparseable or empty URLs are skipped rather than producing empty entries.
	assert.Empty(t, noProxyWithRunnerHosts("", "", "://nope"))
}

func TestApplyProcessProxyEnv(t *testing.T) {
	t.Run("sets both spellings for the runner process", func(t *testing.T) {
		t.Setenv("http_proxy", "")
		t.Setenv("HTTP_PROXY", "")
		t.Setenv("no_proxy", "")
		t.Setenv("NO_PROXY", "")

		ApplyProcessProxyEnv(config.Proxy{HTTPProxy: "http://proxy:3128", NoProxy: "example.com"})

		assert.Equal(t, "http://proxy:3128", os.Getenv("http_proxy"))
		assert.Equal(t, "http://proxy:3128", os.Getenv("HTTP_PROXY"))
		assert.Equal(t, "example.com", os.Getenv("no_proxy"))
		assert.Equal(t, "example.com", os.Getenv("NO_PROXY"))
	})

	// The process environment wins, so a systemd unit or container environment stays in
	// control of the runner it launched.
	t.Run("keeps an environment variable that is already set", func(t *testing.T) {
		t.Setenv("http_proxy", "http://from-env:3128")
		t.Setenv("HTTP_PROXY", "")

		ApplyProcessProxyEnv(config.Proxy{HTTPProxy: "http://from-config:3128"})

		assert.Equal(t, "http://from-env:3128", os.Getenv("http_proxy"))
		assert.Empty(t, os.Getenv("HTTP_PROXY"))
	})
}
