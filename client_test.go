package liguetaxi

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct{
		host	*url.URL
		token	string
		client	*http.Client
	}{
		{&url.URL{}, "abc", nil},
		{&url.URL{}, "abc", &http.Client{}},
		{&url.URL{}, "abc", &http.Client{
			Transport: testRoundTripper(func(r *http.Request) (*http.Response, error) {
				return nil, nil
			}),
		}},
	}

	for _, tc := range testCases {
		c := New(tc.host, tc.token, tc.client)

		if c.host != tc.host {
			t.Errorf("got c.host : %s; want %s.", c.host, tc.host)
		}

		tr := http.RoundTripper(&Transport{tc.token, http.DefaultClient.Transport})
		if tc.client != nil {
			tr = tc.client.Transport
		}

		if !reflect.DeepEqual(c.client.Transport, tr) {
			t.Errorf("got Transport %+v; want %+v.", c.client.Transport, tr)
		}
	}
}
