package harbor

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func (receiver *Client) fetchInitialCSRF() error {
	req, err := http.NewRequest("GET", receiver.endpoint+"/c/login", nil)
	if err != nil {
		return fmt.Errorf("creating initial CSRF request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("initial CSRF request failed: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	csrf := resp.Header.Get("X-Harbor-CSRF-Token")
	if csrf == "" {
		return fmt.Errorf("no X-Harbor-CSRF-Token header in initial response")
	}
	receiver.csrf = csrf

	receiver.mergeCookies(resp.Cookies())
	return nil
}

func (receiver *Client) login() error {
	if err := receiver.fetchInitialCSRF(); err != nil {
		return fmt.Errorf("failed to fetch initial CSRF token: %w", err)
	}

	form := url.Values{}
	form.Set("principal", receiver.username)
	form.Set("password", receiver.password)

	req, err := http.NewRequest("POST", receiver.endpoint+"/c/login", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("creating login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("X-Harbor-CSRF-Token", receiver.csrf)
	if receiver.cookie != "" {
		req.Header.Set("Cookie", receiver.cookie)
	}
	req.Header.Set("Referer", receiver.endpoint+"/account/sign-in")

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 302 {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	if len(resp.Cookies()) == 0 {
		return fmt.Errorf("login succeeded but no session cookie returned")
	}
	receiver.mergeCookies(resp.Cookies())
	return receiver.refreshCSRFToken()
}

func (receiver *Client) refreshCSRFToken() error {
	req, err := receiver.newRequest("GET", "/api/v2.0/configurations")
	if err != nil {
		return err
	}
	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("refreshing CSRF token: %w", err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()

	if csrf := resp.Header.Get("X-Harbor-CSRF-Token"); csrf != "" {
		receiver.csrf = csrf
		return nil
	}
	for _, ck := range resp.Cookies() {
		if ck.Name == "__csrf" {
			receiver.csrf = ck.Value
			return nil
		}
	}
	return nil
}
