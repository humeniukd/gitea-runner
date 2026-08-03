// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package artifactcache

import (
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// The results service is one origin serving every github.actions.results.api.v1 service. Gitea
// implements the artifact half of it and this server the cache half, so neither address is the
// results service on its own, and a client told otherwise reaches for a cache service that is not
// there. That is what breaks the clients this runner cannot patch, docker buildx among them.
//
// Forwarding the artifact service from here makes this origin the whole results service, so
// ACTIONS_RESULTS_URL can point at it truthfully. The instance to forward to travels with the job
// registration rather than with configuration, so a cache server shared between runners serves
// each of their instances. The Host header is rewritten to the address the runner knows Gitea by,
// which is the one a job container would have used itself, so Gitea keeps minting the same upload
// and download URLs and the artifact bodies never travel through here.
const artifactServicePath = "/twirp/github.actions.results.api.v1.ArtifactService/"

// forwardOrNotFound is the router's fallback: the artifact service of the instance the job
// registered with, and the 404 the router would have written otherwise.
func (h *Handler) forwardOrNotFound(w http.ResponseWriter, r *http.Request) {
	cred, ok := h.lookupCredential(bearerToken(r))
	if !ok || cred.Results == "" || !strings.HasPrefix(r.URL.Path, artifactServicePath) {
		http.NotFound(w, r)
		return
	}
	target, err := url.Parse(strings.TrimSuffix(cred.Results, "/"))
	if err != nil {
		h.logger.Errorf("artifact service forward to %q: %v", cred.Results, err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	h.logger.Debugf("%s %s: forwarding to %s", r.Method, r.URL.Path, target)

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(target)
			// Gitea derives the URLs it hands back from the Host, and their scheme from the
			// connection unless a forwarded header overrides it. Presenting the target's host and
			// forwarding no headers of our own is what keeps those URLs pointing at Gitea, over
			// the scheme the runner reaches it by.
			r.Out.Host = target.Host
		},
		ErrorHandler: func(w http.ResponseWriter, _ *http.Request, err error) {
			h.logger.Warnf("artifact service forward to %s: %v", target, err)
			w.WriteHeader(http.StatusBadGateway)
		},
	}
	if cred.InsecureTLS {
		proxy.Transport = insecureTransport
	}
	proxy.ServeHTTP(w, r)
}

// insecureTransport is shared, because a transport per request would pool no connections.
var insecureTransport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}} //nolint:gosec // the runner reaches its instance on the operator's say-so
