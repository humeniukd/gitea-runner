// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package container

import (
	"bytes"
	"cmp"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coder/websocket"
	"go.yaml.in/yaml/v4"
)

const (
	serviceAccountDir = "/var/run/secrets/kubernetes.io/serviceaccount"
	execProtocol      = "v5.channel.k8s.io"
)

type kubeClient struct {
	server, namespace, token string
	tokenFile                string // re-read on every request, service account tokens rotate
	http                     *http.Client
}

type apiError struct {
	code    int
	message string
}

func (e *apiError) Error() string { return e.message }

func isNotFound(err error) bool {
	var apiErr *apiError
	return errors.As(err, &apiErr) && apiErr.code == http.StatusNotFound
}

func newKubeClient(kubeconfig, namespace string) (*kubeClient, error) {
	client := &kubeClient{namespace: namespace}
	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	if host := os.Getenv("KUBERNETES_SERVICE_HOST"); host != "" && kubeconfig == "" {
		if err := client.loadInCluster(host, tlsConfig); err != nil {
			return nil, err
		}
	} else if err := client.loadKubeconfig(kubeconfig, tlsConfig); err != nil {
		return nil, err
	}
	client.http = &http.Client{Transport: &http.Transport{Proxy: http.ProxyFromEnvironment, TLSClientConfig: tlsConfig}}
	return client, nil
}

func (c *kubeClient) loadInCluster(host string, tlsConfig *tls.Config) error {
	c.server = "https://" + net.JoinHostPort(host, os.Getenv("KUBERNETES_SERVICE_PORT"))
	c.tokenFile = serviceAccountDir + "/token"
	ca, err := os.ReadFile(serviceAccountDir + "/ca.crt")
	if err != nil {
		return err
	}
	tlsConfig.RootCAs = x509.NewCertPool()
	tlsConfig.RootCAs.AppendCertsFromPEM(ca)
	if c.namespace == "" {
		namespace, err := os.ReadFile(serviceAccountDir + "/namespace")
		if err != nil {
			return err
		}
		c.namespace = strings.TrimSpace(string(namespace))
	}
	return nil
}

func (c *kubeClient) loadKubeconfig(path string, tlsConfig *tls.Config) error {
	if paths := filepath.SplitList(os.Getenv("KUBECONFIG")); path == "" && len(paths) > 0 {
		path = paths[0]
	}
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		path = filepath.Join(home, ".kube", "config")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("no in-cluster service account and no kubeconfig: %w", err)
	}
	var config struct {
		CurrentContext string `yaml:"current-context"`
		Contexts       []struct {
			Name    string
			Context struct{ Cluster, User, Namespace string }
		}
		Clusters []struct {
			Name    string
			Cluster struct {
				Server                   string
				CertificateAuthority     string `yaml:"certificate-authority"`
				CertificateAuthorityData string `yaml:"certificate-authority-data"`
			}
		}
		Users []struct {
			Name string
			User struct {
				Token                 string
				ClientCertificate     string `yaml:"client-certificate"`
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKey             string `yaml:"client-key"`
				ClientKeyData         string `yaml:"client-key-data"`
			}
		}
	}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return fmt.Errorf("parse kubeconfig %s: %w", path, err)
	}
	var current struct{ Cluster, User, Namespace string }
	for _, entry := range config.Contexts {
		if entry.Name == config.CurrentContext {
			current = entry.Context
		}
	}
	if current.Cluster == "" {
		return fmt.Errorf("kubeconfig %s has no current context", path)
	}
	c.namespace = cmp.Or(c.namespace, current.Namespace, "default")

	read := func(data, file string) ([]byte, error) {
		if data != "" {
			return base64.StdEncoding.DecodeString(data)
		}
		if file == "" {
			return nil, nil
		}
		if !filepath.IsAbs(file) {
			file = filepath.Join(filepath.Dir(path), file)
		}
		return os.ReadFile(file)
	}

	for _, cluster := range config.Clusters {
		if cluster.Name != current.Cluster {
			continue
		}
		c.server = strings.TrimSuffix(cluster.Cluster.Server, "/")
		ca, err := read(cluster.Cluster.CertificateAuthorityData, cluster.Cluster.CertificateAuthority)
		if err != nil {
			return err
		}
		if ca != nil {
			tlsConfig.RootCAs = x509.NewCertPool()
			tlsConfig.RootCAs.AppendCertsFromPEM(ca)
		}
	}
	for _, user := range config.Users {
		if user.Name != current.User {
			continue
		}
		c.token = user.User.Token
		cert, err := read(user.User.ClientCertificateData, user.User.ClientCertificate)
		if err != nil {
			return err
		}
		key, err := read(user.User.ClientKeyData, user.User.ClientKey)
		if err != nil {
			return err
		}
		if cert != nil {
			pair, err := tls.X509KeyPair(cert, key)
			if err != nil {
				return err
			}
			tlsConfig.Certificates = []tls.Certificate{pair}
		}
	}
	if c.server == "" {
		return fmt.Errorf("kubeconfig %s names no server for context %q", path, config.CurrentContext)
	}
	return nil
}

func (c *kubeClient) authorize(header http.Header) error {
	token := c.token
	if c.tokenFile != "" {
		content, err := os.ReadFile(c.tokenFile)
		if err != nil {
			return err
		}
		token = strings.TrimSpace(string(content))
	}
	if token != "" {
		header.Set("Authorization", "Bearer "+token)
	}
	return nil
}

func (c *kubeClient) path(resource string) string {
	return "/api/v1/namespaces/" + c.namespace + "/" + resource
}

func (c *kubeClient) do(ctx context.Context, method, path string, in, out any) error {
	var body io.Reader
	if in != nil {
		content, err := json.Marshal(in)
		if err != nil {
			return err
		}
		body = bytes.NewReader(content)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.server+path, body)
	if err != nil {
		return err
	}
	if err := c.authorize(req.Header); err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusMultipleChoices {
		return statusError(resp)
	}
	if out == nil {
		out = io.Discard // draining lets the connection be reused
	}
	if writer, ok := out.(io.Writer); ok {
		_, err = io.Copy(writer, resp.Body)
		return err
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func statusError(resp *http.Response) error {
	content, _ := io.ReadAll(resp.Body)
	var status struct{ Message string }
	if json.Unmarshal(content, &status) != nil || status.Message == "" {
		status.Message = strings.TrimSpace(string(content))
	}
	return &apiError{code: resp.StatusCode, message: status.Message}
}

func (c *kubeClient) exec(ctx context.Context, pod, containerName string, command []string, stdin io.Reader, stdout, stderr io.Writer) error {
	query := url.Values{"container": {containerName}, "command": command, "stdout": {"true"}, "stderr": {"true"}}
	if stdin != nil {
		query.Set("stdin", "true")
	}
	header := http.Header{}
	if err := c.authorize(header); err != nil {
		return err
	}
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	conn, resp, err := websocket.Dial(ctx, c.server+c.path("pods/"+pod+"/exec?")+query.Encode(), &websocket.DialOptions{
		HTTPClient:   c.http,
		HTTPHeader:   header,
		Subprotocols: []string{execProtocol},
	})
	if err != nil {
		if resp != nil && resp.StatusCode >= http.StatusMultipleChoices {
			defer resp.Body.Close()
			err = statusError(resp)
		}
		return fmt.Errorf("exec in pod %s: %w", pod, err)
	}
	defer func() { _ = conn.CloseNow() }()
	if conn.Subprotocol() != execProtocol {
		return fmt.Errorf("the Kubernetes API server does not speak %s, which needs Kubernetes 1.30 or later", execProtocol)
	}
	conn.SetReadLimit(-1)

	if stdin != nil {
		go func() {
			if err := sendStdin(ctx, conn, stdin); err != nil {
				cancel(err)
			}
		}()
	}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if conn.Ping(ctx) != nil {
					return
				}
			}
		}
	}()

	for {
		_, frame, err := conn.Read(ctx)
		if err != nil {
			return cmp.Or(context.Cause(ctx), fmt.Errorf("exec in pod %s ended without an exit status: %w", pod, err))
		}
		if len(frame) < 2 {
			continue
		}
		switch frame[0] {
		case 1:
			_, err = stdout.Write(frame[1:])
		case 2:
			_, err = stderr.Write(frame[1:])
		case 3:
			return cmp.Or(context.Cause(ctx), execStatus(frame[1:]))
		}
		if err != nil {
			return err
		}
	}
}

func sendStdin(ctx context.Context, conn *websocket.Conn, stdin io.Reader) error {
	frame := make([]byte, 1+32*1024) // frame[0] is the stdin channel, 0
	for {
		n, err := stdin.Read(frame[1:])
		if n > 0 {
			if err := conn.Write(ctx, websocket.MessageBinary, frame[:1+n]); err != nil {
				return err
			}
		}
		if errors.Is(err, io.EOF) {
			return conn.Write(ctx, websocket.MessageBinary, []byte{255, 0}) // v5 close signal for stdin
		}
		if err != nil {
			return err
		}
	}
}

func execStatus(content []byte) error {
	var status struct {
		Status, Message string
		Details         struct {
			Causes []struct{ Reason, Message string }
		}
	}
	if err := json.Unmarshal(content, &status); err != nil {
		return fmt.Errorf("decode exec status: %w", err)
	}
	if status.Status == "Success" {
		return nil
	}
	for _, cause := range status.Details.Causes {
		if code, err := strconv.Atoi(cause.Message); cause.Reason == "ExitCode" && err == nil {
			return ExitCodeError(code)
		}
	}
	return fmt.Errorf("exec failed: %s", status.Message)
}
