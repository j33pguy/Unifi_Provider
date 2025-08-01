// client/unifi/client.go
package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	apiPath    string
	username   string
	password   string
	site       string
	csrfToken  string
}

func NewClient(host, username, password, site string, verifySSL bool) (*Client, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s:8443", host))
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{}
	if !verifySSL {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	httpClient := &http.Client{
		Jar:       jar,
		Transport: tr,
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    u,
		apiPath:    "/proxy/network",
		username:   username,
		password:   password,
		site:       site,
	}

	err = c.login()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) login() error {
	loginURL := c.baseURL.JoinPath("/api/auth/login")

	body := map[string]string{
		"username": c.username,
		"password": c.password,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", loginURL.String(), bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %d", resp.StatusCode)
	}

	c.csrfToken = resp.Header.Get("X-CSRF-Token")

	return nil
}

func (c *Client) SitePath(p string) string {
	return fmt.Sprintf("/api/s/%s/%s", c.site, p)
}

func (c *Client) DoRequest(method, path string, body interface{}, v interface{}) error {
	fullPath := c.apiPath + path
	u := c.baseURL.JoinPath(fullPath)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), reqBody)
	if err != nil {
		return err
	}

	if body != nil || method != "GET" {
		req.Header.Set("Content-Type", "application/json")
	}
	if method != "GET" && c.csrfToken != "" {
		req.Header.Set("X-CSRF-Token", c.csrfToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %d %s", resp.StatusCode, resp.Status)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil
}
