package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// DefaultHTTPClient is the default HTTP client to use for API calls
var DefaultHTTPClient = http.DefaultClient

// Client uses Grafana HTTP API for interacting with Grafana Server
type Client struct {
	baseURL   string
	key       string
	basicAuth bool
	client    *http.Client
}

// NewClient initializes client for interacting with Grafana Server
// authString can be either username:password or a Grafana API key
func NewClient(apiURL, authString string, client *http.Client) *Client {
	key := ""
	basicAuth := strings.Contains(authString, ":")
	baseURL, _ := url.Parse(apiURL)
	if !basicAuth {
		key = fmt.Sprintf("Bearer %s", authString)
	} else {
		parts := strings.Split(authString, ":")
		baseURL.User = url.UserPassword(parts[0], parts[1])
	}
	return &Client{
		baseURL:   baseURL.String(),
		basicAuth: basicAuth,
		key:       key,
		client:    client,
	}
}

func (r *Client) get(query string, params url.Values) ([]byte, int, error) {
	return r.doRequest("GET", query, params, nil)
}

func (r *Client) patch(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PATCH", query, params, bytes.NewBuffer(body))
}

func (r *Client) put(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("PUT", query, params, bytes.NewBuffer(body))
}

func (r *Client) post(query string, params url.Values, body []byte) ([]byte, int, error) {
	return r.doRequest("POST", query, params, bytes.NewBuffer(body))
}

func (r *Client) delete(query string) ([]byte, int, error) {
	return r.doRequest("DELETE", query, nil, nil)
}

func (r *Client) doRequest(method, query string, params url.Values, buf io.Reader) ([]byte, int, error) {
	u, _ := url.Parse(r.baseURL)
	u.Path = path.Join(u.Path, query)
	if params != nil {
		u.RawQuery = params.Encode()
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if !r.basicAuth {
		req.Header.Set("Authorization", r.key)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "grafana-client")
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return data, resp.StatusCode, err
}
