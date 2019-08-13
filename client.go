package liguetaxi

import "net/http"

// Client encapsulates the requests to the
// Ligue Taxi endpoints.
type Client struct {
	// host is the Ligue Taxi URL Host.
	host string
	// client is the http client.
	client *http.Client
}

// New returns a Client for requests Ligue Taxi API.
func New(host, token string, c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	c.Transport = &Transport{
		token,
		c.Transport,
	}
	return &Client{host, c}
}
