// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package run

import (
	"net/url"
	"os"
	"slices"
	"strings"

	"gitea.com/gitea/runner/internal/pkg/config"

	log "github.com/sirupsen/logrus"
)

// Proxy variables are conventionally read in both cases, and tools disagree on which one
// they look at, so both spellings are always set together.
var proxyVarNames = []struct {
	lower string
	upper string
}{
	{"http_proxy", "HTTP_PROXY"},
	{"https_proxy", "HTTPS_PROXY"},
	{"no_proxy", "NO_PROXY"},
}

// ApplyProcessProxyEnv puts the configured proxy in the runner process's own environment,
// which is where Go's HTTP client reads it from for everything the runner fetches itself:
// the Gitea API, and the action repositories it clones. An environment variable that is
// already set wins, so a systemd unit or container environment stays in control.
//
// It must run before the first outbound request, because the standard library resolves
// the proxy environment once.
func ApplyProcessProxyEnv(proxy config.Proxy) {
	for _, value := range []struct {
		names [2]string
		value string
	}{
		{[2]string{"http_proxy", "HTTP_PROXY"}, proxy.HTTPProxy},
		{[2]string{"https_proxy", "HTTPS_PROXY"}, proxy.HTTPSProxy},
		{[2]string{"no_proxy", "NO_PROXY"}, proxy.NoProxy},
	} {
		if value.value == "" {
			continue
		}
		if os.Getenv(value.names[0]) != "" || os.Getenv(value.names[1]) != "" {
			continue
		}
		for _, name := range value.names {
			if err := os.Setenv(name, value.value); err != nil {
				log.Warnf("cannot set %s for the runner process: %v", name, err)
			}
		}
	}
}

// jobProxyEnv returns the proxy variables to give to every job, so that steps and the
// actions they run reach the internet the same way the runner does.
//
// The runner's own servers are added to no_proxy: a job container talking to the cache
// server on the local network must not be sent through an internet proxy, which is the
// failure this saves an operator from debugging. Note that image pulls are made by the
// docker daemon, not by the runner, so they follow the daemon's own proxy configuration.
func jobProxyEnv(proxy config.Proxy, runnerURLs ...string) map[string]string {
	if proxy.HTTPProxy == "" && proxy.HTTPSProxy == "" && proxy.NoProxy == "" {
		return nil
	}

	noProxy := noProxyWithRunnerHosts(proxy.NoProxy, runnerURLs...)
	env := map[string]string{}
	for _, name := range proxyVarNames {
		var value string
		switch name.lower {
		case "http_proxy":
			value = proxy.HTTPProxy
		case "https_proxy":
			value = proxy.HTTPSProxy
		case "no_proxy":
			value = noProxy
		}
		if value == "" {
			continue
		}
		env[name.lower] = value
		env[name.upper] = value
	}
	return env
}

// noProxyWithRunnerHosts appends the hosts of the runner's own servers to a no_proxy list,
// keeping the operator's entries and adding none twice.
func noProxyWithRunnerHosts(noProxy string, runnerURLs ...string) string {
	entries := []string{}
	for entry := range strings.SplitSeq(noProxy, ",") {
		if entry = strings.TrimSpace(entry); entry != "" {
			entries = append(entries, entry)
		}
	}

	for _, raw := range runnerURLs {
		host := hostOf(raw)
		if host == "" || slices.Contains(entries, host) {
			continue
		}
		entries = append(entries, host)
	}
	return strings.Join(entries, ",")
}

// hostOf returns the host of a URL without its port, which is the form a no_proxy entry
// takes.
func hostOf(raw string) string {
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(strings.TrimSuffix(raw, "/"))
	if err != nil || parsed.Hostname() == "" {
		return ""
	}
	return parsed.Hostname()
}
