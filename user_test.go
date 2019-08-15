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

func TestUserRead(t *testing.T) {
	testCases := []struct{
		ctx	context.Context
		name	string
		id	string
		method	string
		path	endpoint
		body	userFilter
		wantRes	UserResponse
	}{
		{
			context.Background(),
			"test",
			"123",
			http.MethodPost,
			readUserEndpoint,
			userFilter{ "123", "test"},
			UserResponse{
				status: status{ReqStatusOK},
			},
		},
	}

	for _, tc := range testCases {
		req := &testRequester{output: tc.wantRes}
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
