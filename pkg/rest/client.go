package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client defines the structure for API requests
type Client struct {
	BaseUrl    string
	HttpClient *http.Client
	Headers    map[string]string
}

type Config struct {
	// BaseUrl is the base URL the requests are sent to.
	// Any instance of Client must at least contain a valid BaseUrl.
	BaseUrl string
	// Timeout specifies a time limit for requests made by this
	// client as described in https://golang.org/pkg/net/http/#Client.
	// A zero value means no timeout.
	Timeout time.Duration
}

// NewClient initializes and returns a new Client.
func NewClient(config *Config) *Client {
	return &Client{
		BaseUrl: config.BaseUrl,
		HttpClient: &http.Client{
			Timeout: config.Timeout,
		},
		Headers: make(map[string]string),
	}
}

// SetHeader is a convenience method that allows setting custom headers
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

// Request makes an HTTP request with dynamic method, path, and body
func (c *Client) Request(method, path string, body interface{}) (*http.Response, error) {
	url := c.BaseUrl + path

	// Convert body to JSON if provided
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body to JSON failed: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
