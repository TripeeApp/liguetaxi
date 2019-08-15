package liguetaxi

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestUserStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct{
		b	[]byte
		want	userStatus
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
		b	[]byte
		want	string
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
	body	interface{}
	ctx	context.Context
	err	error
	method	string
	output	interface{}
	path	endpoint

}

func (t *testRequester) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	t.ctx = ctx
	t.method = method
	t.path = path
	t.body = body

	out := reflect.ValueOf(output)
	if !out.IsNil() && out.Elem().CanSet() && t.output != nil {
		out.Elem().Set(reflect.ValueOf(t.output))
	}

	return t.err
}

func TestUser(t *testing.T) {
	testCases := []struct{
		name	string
		call	func(ctx context.Context, req requester) (resp interface{}, err error)
		ctx	context.Context
		method	string
		path	endpoint
		body	interface{}
		wantRes	interface{}
	}{
		{
			"Read()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).Read(ctx, "123", "test")
				return
			},
			context.Background(),
			http.MethodPost,
			readUserEndpoint,
			userFilter{ "123", "test"},
			UserResponse{
				status: status{ReqStatusOK},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // create scoped test case
		t.Run(tc.name, func(t *testing.T) {
			req := &testRequester{output: tc.wantRes}

			res, err := tc.call(tc.ctx, req)
			if err != nil {
				t.Fatalf("got error while calling User %s: %s, want nil", tc.name, err.Error())
			}

			if !reflect.DeepEqual(req.ctx, tc.ctx) {
				t.Errorf("got Requester Context %+v; want %+v.", req.ctx, tc.ctx)
			}

			if req.method != tc.method {
				t.Errorf("got request method: %s; want %s.", req.method, tc.method)
			}

			if req.path != tc.path {
				t.Errorf("got request path: %s; want %s.", req.path, tc.path)
			}

			if !reflect.DeepEqual(req.body, tc.body) {
				t.Errorf("got request body: %+v; want %+v.", req.body, tc.body)
			}

			if !reflect.DeepEqual(res, tc.wantRes) {
				t.Errorf("got response: %+v; want %+v.", res, tc.wantRes)
			}
		})
	}
}

func TestUserReadError(t *testing.T) {
	testCases := []struct{
		err	error
	}{
		{errors.New("Error")},
	}

	for _, tc := range testCases {
		req := &testRequester{err: tc.err}
		u := &UserService{req}

		_, err := u.Read(nil, "123", "Test")
		if err == nil {
			t.Error("got error nil while calling User.Read; want not nil.")
		}
	}
}
