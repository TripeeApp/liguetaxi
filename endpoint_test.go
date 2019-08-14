package liguetaxi

import (
	"testing"
	"context"
)

func TestEndpointString(t *testing.T) {
	testCases := []struct{
		ctx	 context.Context
		endpoint Endpoint
		want	 string
	}{
		{
			context.Background(),
			Endpoint("/test"),
			"/test/json",
		},
		{

			context.WithValue(context.Background(), ResTypeKey, "json"),
			Endpoint("/test"),
			"/test/json",
		},
	}

	for _, tc := range testCases {
		e := tc.endpoint.String(tc.ctx)

		if e != tc.want {
			t.Errorf("got Endpoint.String(%+v): %s; want %s.", tc.ctx, e, tc.want)
		}
	}
}
