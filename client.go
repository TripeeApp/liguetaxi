package liguetaxi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// reqStatus is the request status.
// Success = 1
// Fail = 0
type reqStatus int

const (
	// Request statuses.
	ReqStatusFail reqStatus = iota
	ReqStatusOK

	// Error message format.
	errFmt = `Error while request the LigueTaxi API: %s; Status Code: %d; Body: %s.`
)

// status is the request status.
type status struct {
	Status reqStatus `json:"status"`
}

// ApiError implements the error interface
// and returns infos from the request
type ApiError struct {
	statusCode int
	body       []byte
	msg        string
}

func (e *ApiError) Error() string {
	return fmt.Sprintf(errFmt, e.msg, e.statusCode, e.body)
}

// requester is the interface that performs a request
// to the server and parses the payload.
type requester interface {
	Request(ctx context.Context, method string, path endpoint, body, output interface{}) error
}

// Client encapsulates the requests to the
// Ligue Taxi endpoints.
type Client struct {
	// host is the Ligue Taxi URL Host.
	host *url.URL
	// client is the http client.
	client *http.Client

	// User is the service that handles http logic for requests
	// related to the user.
	User *UserService
}

// New returns a Client for requests Ligue Taxi API.
func NewClient(host *url.URL, token string, c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}
	c.Transport = &Transport{
		token,
		c.Transport,
	}

	client := &Client{host: host, client: c}
	client.User = &UserService{client}
	return client
}

// Request created an API request. A relative path can be providaded
// in which case it is resolved relative to the host of the Client.
func (c *Client) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	u, err := c.host.Parse(path.String(ctx))
	if err != nil {
		return err
	}

	var b io.ReadWriter
	if body != nil {
		b = new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(method, u.String(), b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")

	req = req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO: add tests for error on reading body
	r, _ := ioutil.ReadAll(res.Body)

	// TODO: Implements the XML decoding based on the
	// endpoint's ContextType(ctx) value.
	// For now the JSON decoding will work.
	if err := json.Unmarshal(r, output); err != nil {
		return &ApiError{
			statusCode: res.StatusCode,
			body:       r,
			msg:        err.Error(),
		}
	}

	return nil
}
