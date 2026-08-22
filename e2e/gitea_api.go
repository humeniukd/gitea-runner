// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:build e2e

package e2e

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const pollInterval = 250 * time.Millisecond

type GiteaAPI struct {
	baseURL string
	token   string
}

type ActionRun struct {
	ID         int64  `json:"id"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HeadBranch string `json:"head_branch"`
	Path       string `json:"path"`
	Event      string `json:"event"`
}

type actionRunsResponse struct {
	WorkflowRuns []ActionRun `json:"workflow_runs"`
	TotalCount   int64       `json:"total_count"`
}

type ActionJob struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

type actionJobsResponse struct {
	Jobs []ActionJob `json:"jobs"`
}

type Artifact struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type artifactsResponse struct {
	Artifacts []Artifact `json:"artifacts"`
}

func (a *GiteaAPI) CreateRepo(ctx context.Context, name string, autoInit bool) error {
	body := map[string]any{"name": name, "auto_init": autoInit}
	return a.doJSON(ctx, http.MethodPost, "/api/v1/user/repos", body, nil)
}

func (a *GiteaAPI) CreateFile(ctx context.Context, owner, repo, path, content, message string) error {
	body := map[string]any{
		"content": base64.StdEncoding.EncodeToString([]byte(content)),
		"message": message,
	}
	url := fmt.Sprintf("/api/v1/repos/%s/%s/contents/%s", owner, repo, path)
	return a.doJSON(ctx, http.MethodPost, url, body, nil)
}

func (a *GiteaAPI) CreateSecret(ctx context.Context, owner, repo, name, value string) error {
	body := map[string]any{"data": value}
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/secrets/%s", owner, repo, name)
	return a.doJSON(ctx, http.MethodPut, url, body, nil)
}

func (a *GiteaAPI) DefaultBranch(ctx context.Context, owner, repo string) (string, error) {
	var resp struct {
		DefaultBranch string `json:"default_branch"`
	}
	url := fmt.Sprintf("/api/v1/repos/%s/%s", owner, repo)
	if err := a.doJSON(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return "", err
	}
	return resp.DefaultBranch, nil
}

func (a *GiteaAPI) CreateVariable(ctx context.Context, owner, repo, name, value string) error {
	body := map[string]any{"value": value}
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/variables/%s", owner, repo, name)
	return a.doJSON(ctx, http.MethodPost, url, body, nil)
}

func (a *GiteaAPI) DispatchWorkflow(ctx context.Context, owner, repo, workflowID, ref string, inputs map[string]string) error {
	body := map[string]any{"ref": ref}
	if len(inputs) > 0 {
		body["inputs"] = inputs
	}
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflowID)
	return a.doJSON(ctx, http.MethodPost, url, body, nil)
}

var ErrCancelUnsupported = errors.New("run cancellation is unsupported by this gitea version")

func (a *GiteaAPI) CancelRun(ctx context.Context, owner, repo string, runID int64) error {
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs/%d/cancel", owner, repo, runID)
	err := a.doJSON(ctx, http.MethodPost, url, nil, nil)
	var status statusError
	if errors.As(err, &status) && status.routeAbsent() {
		return ErrCancelUnsupported
	}
	return err
}

var ErrNoRun = errors.New("no workflow run yet")

func (a *GiteaAPI) LatestRun(ctx context.Context, owner, repo string) (*ActionRun, error) {
	var resp actionRunsResponse
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs?limit=1", owner, repo)
	if err := a.doJSON(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return nil, err
	}
	if len(resp.WorkflowRuns) == 0 {
		return nil, ErrNoRun
	}
	return &resp.WorkflowRuns[0], nil
}

func (a *GiteaAPI) Runs(ctx context.Context, owner, repo string) ([]ActionRun, error) {
	var resp actionRunsResponse
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs?limit=50", owner, repo)
	if err := a.doJSON(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return nil, err
	}
	return resp.WorkflowRuns, nil
}

// Filter by path "<workflow>@<ref>"; Gitea has no per-workflow runs route.
func (a *GiteaAPI) RunsForWorkflow(ctx context.Context, owner, repo, workflow string) ([]ActionRun, error) {
	runs, err := a.Runs(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	matched := make([]ActionRun, 0, len(runs))
	for _, run := range runs {
		file, _, _ := strings.Cut(run.Path, "@")
		if file == workflow || strings.HasSuffix(file, "/"+workflow) {
			matched = append(matched, run)
		}
	}
	return matched, nil
}

func (a *GiteaAPI) WaitForRunConclusion(ctx context.Context, owner, repo string, runID int64, timeout time.Duration) (*ActionRun, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs/%d", owner, repo, runID)
	for {
		var run ActionRun
		if err := a.doJSON(ctx, http.MethodGet, url, nil, &run); err != nil {
			return nil, err
		}
		if run.Status == "completed" {
			return &run, nil
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("run %d did not complete within %s (last status %q): %w", runID, timeout, run.Status, ctx.Err())
		case <-time.After(pollInterval):
		}
	}
}

func (a *GiteaAPI) Jobs(ctx context.Context, owner, repo string, runID int64) ([]ActionJob, error) {
	var resp actionJobsResponse
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID)
	if err := a.doJSON(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Jobs, nil
}

func (a *GiteaAPI) JobLogs(ctx context.Context, owner, repo string, jobID int64) (string, error) {
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.baseURL+url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "token "+a.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("GET %s: %d: %s", url, resp.StatusCode, body)
	}
	return string(body), nil
}

func (a *GiteaAPI) Artifacts(ctx context.Context, owner, repo string, runID int64) ([]Artifact, error) {
	var resp artifactsResponse
	url := fmt.Sprintf("/api/v1/repos/%s/%s/actions/runs/%d/artifacts", owner, repo, runID)
	if err := a.doJSON(ctx, http.MethodGet, url, nil, &resp); err != nil {
		return nil, err
	}
	return resp.Artifacts, nil
}

func (a *GiteaAPI) doJSON(ctx context.Context, method, path string, reqBody, respBody any) error {
	var bodyReader io.Reader
	if reqBody != nil {
		encoded, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+a.token)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return statusError{method: method, path: path, code: resp.StatusCode, body: string(body)}
	}
	if respBody != nil && len(body) > 0 {
		return json.Unmarshal(body, respBody)
	}
	return nil
}

type statusError struct {
	method, path string
	code         int
	body         string
}

func (e statusError) Error() string {
	return fmt.Sprintf("%s %s: %d: %s", e.method, e.path, e.code, e.body)
}

// Unregistered routes return plain "404 page not found".
func (e statusError) routeAbsent() bool {
	return e.code == http.StatusNotFound && strings.Contains(e.body, "page not found")
}
