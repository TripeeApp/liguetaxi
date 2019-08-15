package liguetaxi

import (
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestUserStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct{
		b []byte
		want userStatus
	}{
		{[]byte(`24`), UserStatusActive},
		{[]byte(`25`), UserStatusInactive},
	}

	for _, tc := range testCases {
		var status userStatus

		status.UnmarshalJSON(tc.b)

		if status != tc.want {
			t.Errorf("got userStatus.UnmarshalJSON(%s): %v; want %v.", tc.b, status, tc.want)
		}
	}
}

func TestEmptyObjToStrUnmarshalJSON(t *testing.T) {
	testCases := []struct{
		b []byte
		want string
	}{
		{[]byte(`non-empty string`), "non-empty string"},
		{[]byte(`{}`), ""},
	}

	for _, tc := range testCases {
		var str emptyObjToStr

		str.UnmarshalJSON(tc.b)

		if string(str) != tc.want {
			t.Errorf("got emptyObjToStr.UnmarshalJSON(%s): %v; want %s.", tc.b, str, tc.want)
		}
	}
}

type testRequester struct{
	body interface{}
	ctx context.Context
	method string
	path string

	res *http.Response
	err error
}

func (t *testRequester) Request(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	t.ctx = ctx
	t.method = method
	t.path = path
	t.body = body

	return t.res, t.err
}

func TestUserRead(t *testing.T) {
	testCases := []struct{
		ctx		context.Context
		name		string
		id		string
		res		*http.Response
		method		string
		path		string
		body		userFilter
		wantRes		UserResponse
	}{
		{
			context.Background(),
			"test",
			"123",
			&http.Response{
				StatusCode: http.StatusOK,
				Body: ioutil.NopCloser(strings.NewReader(`{"status":1}`)),
			},
			http.MethodPost,
			readUserEndpoint.String(`json`),
			userFilter{ "123", "test"},
			UserResponse{
				status: status{ReqStatusOK},
			},
		},
	}

	for _, tc := range testCases {
		req := &testRequester{res: tc.res}
		u := &UserService{req}

		res, err := u.Read(tc.ctx, tc.id, tc.name)
		if err != nil {
			t.Fatalf("got error while calling User.Read(%+v, %s, %s): %s, want nil", tc.ctx, tc.id, tc.name, err.Error())
		}

		if !reflect.DeepEqual(req.ctx, tc.ctx) {
			t.Errorf("got Requester Context %+v; want %+v.", req.ctx, tc.ctx)
		}

		if req.method != tc.method {
			t.Errorf("got Requester Method: %s; want %s.", req.method, tc.method)
		}

		if req.path != tc.path {
			t.Errorf("got Requester Path: %s; want %s.", req.path, tc.path)
		}

		if !reflect.DeepEqual(req.body, tc.body) {
			t.Errorf("got Requester Bodt: %+v; want %+v.", req.body, tc.body)
		}

		if !reflect.DeepEqual(res, tc.wantRes) {
			t.Errorf("got UserResponse: %+v; want %+v.", res, tc.wantRes)
		}
	}
}
