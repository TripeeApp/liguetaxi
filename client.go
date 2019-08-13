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

func New(host, token string, c *http.Client) *Client {
	return &Client{
		host,
		&http.Client{
			Transport: &Transport{
				token,
				http.DefaultClient.Transport,
			},
		},
	}
}
