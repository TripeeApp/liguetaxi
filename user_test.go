package liguetaxi

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestUserStatusUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		b    []byte
		want userStatus
	}{
		{[]byte(`"24"`), UserStatusActive},
		{[]byte(`"25"`), UserStatusInactive},
		{[]byte(`"46"`), UserStatusSynching},
		{[]byte(`"46"`), UserStatusSynching},
	}

	for _, tc := range testCases {
		var status userStatus

		if err := status.UnmarshalJSON(tc.b); err != nil {
			t.Fatalf("got error calling userStatus.UnmarshalJSON(%+v): %s; want nil.", tc.b, err.Error())
		}

		if status != tc.want {
			t.Errorf("got userStatus.UnmarshalJSON(%s): %v; want %v.", tc.b, status, tc.want)
		}
	}
}

func TestUserStatusMarshalJSON(t *testing.T) {
	testCases := []struct {
		status *userStatus
		want   []byte
	}{
		{UserStatusActive.New(), []byte(`"24"`)},
		{UserStatusInactive.New(), []byte(`"25"`)},
		{UserStatusSynching.New(), []byte(`"46"`)},
		{nil, []byte(`null`)},
	}

	for _, tc := range testCases {
		got, err := tc.status.MarshalJSON()
		if err != nil {
			t.Fatalf("got error calling userStatus.MarshalJSON(): %s; want nil.", err.Error())
		}

		if !bytes.Equal(got, tc.want) {
			t.Errorf("got userStatus.MarshalJSON(): %s; want %s.", got, tc.want)
		}
	}
}

func TestEmptyObjToStrUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		b    []byte
		want string
	}{
		{[]byte(`"non-empty string"`), "non-empty string"},
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

type testRequester struct {
	body   interface{}
	ctx    context.Context
	err    error
	method string
	output reflect.Value
	path   endpoint
}

func (t *testRequester) Request(ctx context.Context, method string, path endpoint, body, output interface{}) error {
	t.ctx = ctx
	t.method = method
	t.path = path
	t.body = body

	if t.output.IsValid() {
		out := reflect.ValueOf(output)
		if !out.IsNil() && out.Elem().CanSet() {
			out.Elem().Set(t.output)
		}
	}

	return t.err
}

func TestUser(t *testing.T) {
	testCases := []struct {
		name    string
		call    func(ctx context.Context, req requester) (resp interface{}, err error)
		ctx     context.Context
		method  string
		path    endpoint
		body    interface{}
		wantRes interface{}
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
			userFilter{"123", "test"},
			&UserResponse{
				Status: ReqStatusOK,
				Data: DataUser{
					Status: UserStatusActive.New(),
				},
			},
		},
		{
			"Create()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).Create(ctx, &User{Name: "Test"})
				return
			},
			context.Background(),
			http.MethodPost,
			createUserEndpoint,
			&User{Name: "Test"},
			&OperationResponse{
				Status: ReqStatusOK,
			},
		},
		{
			"Update()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).Update(ctx, &User{Name: "Test"})
				return
			},
			context.Background(),
			http.MethodPost,
			updateUserEndpoint,
			&User{Name: "Test"},
			&OperationResponse{
				Status: ReqStatusOK,
			},
		},
		{
			"UpdateStatus()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).UpdateStatus(ctx, &UserStatus{Name: "Test", Status: UserStatusInactive})
				return
			},
			context.Background(),
			http.MethodPost,
			updateUserStatusEndpoint,
			&UserStatus{Name: "Test", Status: UserStatusInactive},
			&OperationResponse{
				Status: ReqStatusOK,
			},
		},
		{
			"ReadClassifier()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).ReadClassifier(ctx, "1", "test")
				return
			},
			context.Background(),
			http.MethodPost,
			readClassifierEndpoint,
			classifierFilter{Field: "1", Value: "test"},
			&ClassifierResponse{
				Status: ReqStatusOK,
			},
		},
		{
			"CreateClassifier()",
			func(ctx context.Context, req requester) (resp interface{}, err error) {
				resp, err = (&UserService{req}).CreateClassifier(ctx, &Classifier{Field: "test", Value: "test2"})
				return
			},
			context.Background(),
			http.MethodPost,
			createClassifierEndpoint,
			&Classifier{Field: "test", Value: "test2"},
			&ClassifierOperationResponse{
				OperationResponse: OperationResponse{
					Status: ReqStatusOK,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // creates scoped test case
		t.Run(tc.name, func(t *testing.T) {
			req := &testRequester{output: reflect.ValueOf(tc.wantRes).Elem()}

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

func TestUserError(t *testing.T) {
	testCases := []struct {
		name string
		call func(req requester) error
		err  error
	}{
		{
			"Read()",
			func(req requester) error {
				_, err := (&UserService{req}).Read(context.Background(), "123", "test")
				return err
			},
			errors.New("Error"),
		},
		{
			"Create()",
			func(req requester) error {
				_, err := (&UserService{req}).Create(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"Update()",
			func(req requester) error {
				_, err := (&UserService{req}).Update(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
		{
			"ReadClassifier()",
			func(req requester) error {
				_, err := (&UserService{req}).ReadClassifier(context.Background(), "1", "test")
				return err
			},
			errors.New("Error"),
		},
		{
			"CreateClassifier()",
			func(req requester) error {
				_, err := (&UserService{req}).CreateClassifier(context.Background(), nil)
				return err
			},
			errors.New("Error"),
		},
	}

	for _, tc := range testCases {
		tc := tc // creates scoped test case
		t.Run(tc.name, func(t *testing.T) {
			req := &testRequester{err: tc.err}

			err := tc.call(req)
			if !reflect.DeepEqual(err, tc.err) {
				t.Errorf("got error: %s; want %s.", err, tc.err)
			}
		})
	}
}
