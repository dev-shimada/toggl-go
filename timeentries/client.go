// Package timeentries provides a client for interacting with the Toggl API.
package timeentries

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// Client represents a Toggl API client.
type Client struct {
	HttpClient httpClient
	Token      string
}

// NewClient creates a new Client with the given API token.
func NewClient(token string) Client {
	return Client{
		HttpClient: &http.Client{},
		Token:      token,
	}
}

func (c Client) newRequest(u url.URL) http.Request {
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	// header.Add("charset", "utf-8")
	toggl := http.Request{
		URL: &url.URL{
			Scheme:   "https",
			Host:     "api.track.toggl.com",
			User:     url.UserPassword(c.Token, "api_token"),
			RawQuery: u.RawQuery,
		},
		Header: header,
	}
	return toggl
}

// Get creates a GET request to the specified URL.
func (c Client) Get(u url.URL) http.Request {
	toggl := c.newRequest(u)
	toggl.Method = http.MethodGet
	return toggl
}

// Post creates a POST request to the specified URL with the given body.
func (c Client) Post(u url.URL, body []byte) http.Request {
	toggl := c.newRequest(u)
	toggl.Method = http.MethodPost
	toggl.Body = io.ReadCloser(io.NopCloser(bytes.NewBuffer(body)))
	return toggl
}

// Patch creates a PATCH request to the specified URL with the given body.
func (c Client) Patch(u url.URL, body []byte) http.Request {
	toggl := c.newRequest(u)
	toggl.Method = http.MethodPatch
	toggl.Body = io.ReadCloser(io.NopCloser(bytes.NewBuffer(body)))
	return toggl
}

// Put creates a PUT request to the specified URL with the given body.
func (c Client) Put(u url.URL, body []byte) http.Request {
	toggl := c.newRequest(u)
	toggl.Method = http.MethodPut
	toggl.Body = io.ReadCloser(io.NopCloser(bytes.NewBuffer(body)))
	return toggl
}

// Delete creates a DELETE request to the specified URL.
func (c Client) Delete(u url.URL) http.Request {
	toggl := c.newRequest(u)
	toggl.Method = http.MethodDelete
	return toggl
}
