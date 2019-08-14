package liguetaxi

import (
	"context"
	"encoding/json"
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// Client encapsulates the requests to the
// Ligue Taxi endpoints.
type Client struct {
	// host is the Ligue Taxi URL Host.
	host *url.URL
	// client is the http client.
	client *http.Client
}

// New returns a Client for requests Ligue Taxi API.
func New(host *url.URL, token string, c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	c.Transport = &Transport{
		token,
		c.Transport,
	}
	return &Client{host, c}
}

// Request created an API request. A relative path can be providaded
// in which case it is resolved relative to the host of the Client.
func (c *Client) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	u, _ := c.host.Parse(path)
	var b io.ReadWriter
	if body != nil {
		b = new(bytes.Buffer)
		json.NewEncoder(b).Encode(body)
	}
	req, _ := http.NewRequest(method, u.String(), b)
	return c.client.Do(req)
}
