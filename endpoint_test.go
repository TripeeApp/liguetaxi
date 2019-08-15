package liguetaxi

import (
	"testing"
	"context"
)

func TestEndpointContextType(t *testing.T) {
	testCases := []struct{
		ctx	 context.Context
		endpoint endpoint
		want	 string
	}{
		{
			nil,
			endpoint("/test"),
			Json,
		},
		{
			context.Background(),
			endpoint("/test"),
			Json,
		},
		{
			context.WithValue(context.Background(), ResType, Json),
			endpoint("/test"),
			Json,
		},
		{
			context.WithValue(context.Background(), ResType, Xml),
			endpoint("/test"),
			Xml,
		},
	}


	for _, tc := range testCases {
		e := tc.endpoint.ContextType(tc.ctx)

		if e != tc.want {
			t.Errorf("got Endpoint.ContextType(%+v): %s; want %s.", tc.ctx, e, tc.want)
		}
	}

}

func TestEndpointString(t *testing.T) {
	testCases := []struct{
		suffix	 string
		endpoint endpoint
		want	 string
	}{
		{
			Json,
			endpoint("/test"),
			"/test/json",
		},
		{
			Xml,
			endpoint("/test2"),
			"/test2/xml",
		},
	}

	for _, tc := range testCases {
		e := tc.endpoint.String(tc.suffix)

		if e != tc.want {
			t.Errorf("got Endpoint.String(%s): %s; want %s.", tc.suffix, e, tc.want)
		}
	}
}
