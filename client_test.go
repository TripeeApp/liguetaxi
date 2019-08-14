package liguetaxi

import (
	"context"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func newMockServer(handler func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestRequest(t *testing.T) {
	testCases := []struct{
		ctx context.Context
		path	string
		method	string
		body	interface{}
		server	*httptest.Server
		want	*http.Response
	}{
		{
			context.Background(),
			"/",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
		{
			context.Background(),
			"",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("go Request.Method %s; want %s.", r.Method, http.MethodGet)
				}
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
		{
			context.Background(),
			"foo/",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.Path; got != "/foo/" {
					t.Errorf("got Request.URL: %s; want foo/.", got)
				}
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
		{
			context.Background(),
			"",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Body != http.NoBody {
					t.Errorf("got Request.Body: %+v, want empty.", r.Body)
				}
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
		{
			context.Background(),
			"",
			http.MethodPost,
			struct{Name string `json:"name"`}{"Testing"},
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				if r.Body == http.NoBody {
					t.Error("got Request.Body empty, want not empty.")
				}

				got, _ := ioutil.ReadAll(r.Body)

				if want := []byte(`{"name":"Testing"}`); !bytes.Contains(got, want) {
					t.Errorf("got body: %s, want %s.", got, want)
				}
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
	}

	for _, tc := range testCases {
		u, _ := url.Parse(tc.server.URL)
		c := New(u, "abc", nil)

		res, err := c.Request(tc.ctx, tc.method, tc.path, tc.body)
		if err != nil {
			t.Fatalf("got error %s; want nil.", err.Error())
		}

		if res.StatusCode != tc.want.StatusCode {
			t.Errorf("got Response StatusCode %d; want %d.", res.StatusCode, tc.want.StatusCode)
		}

		tc.server.Close()
	}
}

func TestRequestError(t *testing.T) {
	testCases := []struct{
		ctx context.Context
		path	string
		method	string
		body	interface{}
		server	*httptest.Server
		want	*http.Response
	}{
		{
			context.Background(),
			":",
			http.MethodGet,
			nil,
			newMockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			&http.Response{StatusCode: http.StatusOK},
		},
	}

	for _, tc := range testCases {
		u, _ := url.Parse(tc.server.URL)
		c := New(u, "abc", nil)

		_, err := c.Request(tc.ctx, tc.method, tc.path, tc.body)
		if err == nil {
			t.Errorf("got error nil; want not nil.")
		}

		tc.server.Close()
	}
}
