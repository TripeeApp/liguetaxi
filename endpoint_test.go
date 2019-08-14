package liguetaxi

import (
	"testing"
	"context"
)

func TestEndpointString(t *testing.T) {
	testCases := []struct{
		context func() context.Context
		endpoint Endpoint
		want string
	}{
		{
			func() context.Context {
				return context.Background()
			},
			Endpoint("/test"),
			"/test/json",
		},
	}

	for _, tc := range testCases {
		e := tc.endpoint.String(tc.context())

		if e != tc.want {
			t.Errorf("got Endpoint.String(%+v): %s; want %s.", tc.context(), e, tc.want)
		}
	}
}
