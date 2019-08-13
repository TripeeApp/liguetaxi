package liguetaxi

import (
	"fmt"
	"net/http"
)

// Transport is the RountTripper for injecting
// the authorization header for every request to
// Ligue Taxi player.
type Transport struct {
	// token is the string to be injected to the
	// outgoing request header
	token string

	// Base is the base RoundTripper to make HTTP request.
	Base http.RoundTripper
}

// RoundTrip injects the Authorization Header with the Token
func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	// We should not modify the origin request
	// per RoundTripper contract. See 
	// https://golang.org/pkg/net/http/#RoundTripper
	req := cloneReq(r)
	// Injects the Authorization Header
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", t.token))

	return t.Base.RoundTrip(req)
}

// cloneReq returns a clone of the *http.Request.
// the clone is a shallow copy of the struct and its Header map.
func cloneReq(r *http.Request) *http.Request {
	// shalow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
