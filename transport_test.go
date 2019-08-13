package liguetaxi

import (
	"net/http/httptest"
	"net/http"
	"testing"
)

type testRoundTripper func(r *http.Request) (*http.Response, error)

func (rt testRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return rt(r)
}

func TestRoundTrip(t *testing.T) {
	var reqSent *http.Request

	baseTripperCalled := false
	expectedAuth := "Basic abc"

	rt := testRoundTripper(func(r *http.Request) (*http.Response, error) {
		baseTripperCalled = true
		reqSent = r
		return nil, nil
	})

	tr := &Transport{"abc", rt}

	req := httptest.NewRequest(http.MethodGet, "", nil)

	tr.RoundTrip(req)

	if !baseTripperCalled {
		t.Error("expected Transport.RoundTrip() to call base RoundTripper.")
	}

	if auth := reqSent.Header.Get("Authorization"); auth != expectedAuth {
		t.Errorf("got Authorization Header: %s; want %s", auth, expectedAuth)
	}
}
