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
			endpoint("test"),
			Json,
		},
		{
			context.Background(),
			endpoint("test"),
			Json,
		},
		{
			context.WithValue(context.Background(), ResType, Json),
			endpoint("test"),
			Json,
		},
		{
			context.WithValue(context.Background(), ResType, Xml),
			endpoint("test"),
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
		ctx	 context.Context
		endpoint endpoint
		want	 string
	}{
		{
			nil,
			endpoint("test"),
			"api/test/json",
		},
		{
			context.Background(),
			endpoint("test2"),
			"api/test2/json",
		},
		{
			context.WithValue(context.Background(), ResType, Json),
			endpoint("test2"),
			"api/test2/json",
		},
		{
			context.WithValue(context.Background(), ResType, Xml),
			endpoint("test2"),
			"api/test2/xml",
		},
	}

	for _, tc := range testCases {
		e := tc.endpoint.String(tc.ctx)

		if e != tc.want {
			t.Errorf("got Endpoint.String(%+v): %s; want %s.", tc.ctx, e, tc.want)
		}
	}
}
