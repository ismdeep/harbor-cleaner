package harbor

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ismdeep/log"
)

type Client struct {
	ctx         context.Context
	endpoint    string
	host        string
	username    string
	password    string
	cookie      string
	csrf        string
	concurrency int
	httpClient  *http.Client
}

type TagInfo struct {
	Repository string
	Name       string
	Digest     string
	Image      string
	Size       int64
}

func NewClient(ctx context.Context, endpoint string, username string, password string, concurrency int) (*Client, error) {

	client := Client{
		ctx:         ctx,
		endpoint:    endpoint,
		host:        "",
		username:    username,
		password:    password,
		cookie:      "",
		csrf:        "",
		concurrency: concurrency,
		httpClient:  &http.Client{Transport: &http.Transport{}},
	}

	u, err := url.Parse(client.endpoint)
	if err != nil {
		return nil, errors.Join(errors.New("failed to extract host from endpoint"), err)
	}
	client.host = u.Host

	log.WithContext(ctx).Info("Logging in...")
	if err := client.login(); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}
	log.WithContext(ctx).Info("Logged in.")

	return &client, nil
}

func (receiver *Client) Concurrency() int {
	return receiver.concurrency
}

func (receiver *Client) Host() string {
	return receiver.host
}

func escapeRepoPath(repo string) string {
	return strings.ReplaceAll(repo, "/", "%252F")
}

func (receiver *Client) newRequest(method, path string) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%v%v", receiver.endpoint, path), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Cookie", receiver.cookie)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if receiver.csrf != "" {
		req.Header.Set("X-Harbor-CSRF-Token", receiver.csrf)
	}
	return req, nil
}

func (receiver *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("unauthorized (401) — check credentials")
	}
	if resp.StatusCode == http.StatusForbidden {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("forbidden (403) — check permissions")
	}
	return resp, nil
}

func (receiver *Client) setAuth(req *http.Request) {
	req.Header.Set("Cookie", receiver.cookie)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	if receiver.csrf != "" {
		req.Header.Set("X-Harbor-CSRF-Token", receiver.csrf)
	}
}

func (receiver *Client) mergeCookies(newCookies []*http.Cookie) {
	existing := parseCookieString(receiver.cookie)
	for _, ck := range newCookies {
		existing[ck.Name] = ck.Value
	}
	var parts []string
	for name, value := range existing {
		parts = append(parts, name+"="+value)
	}
	receiver.cookie = strings.Join(parts, "; ")
}

func parseCookieString(raw string) map[string]string {
	m := make(map[string]string)
	if raw == "" {
		return m
	}
	for _, pair := range strings.Split(raw, "; ") {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return m
}
