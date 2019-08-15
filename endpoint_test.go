package liguetaxi

import (
	"testing"
	"context"
)

func TestEndpointContextType(t *testing.T) {
	tetCases := []struct{
		ctx	 context.Context
		endpoint endpoint
		want	 string
	}{
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
		}
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
			context.Background(),
			endpoint("/test"),
			"/test/json",
		},
		{

			context.WithValue(context.Background(), ResType, Json),
			endpoint("/test2"),
			"/test2/json",
		},
		{

			context.WithValue(context.Background(), ResType, Xml),
			endpoint("/test3"),
			"/test3/xml",
		},
		{

			context.WithValue(context.Background(), ResType, ""),
			endpoint("/test4"),
			"/test4/json",
		},
	}

	for _, tc := range testCases {
		e := tc.endpoint.String(tc.ctx)

		if e != tc.want {
			t.Errorf("got Endpoint.String(%+v): %s; want %s.", tc.ctx, e, tc.want)
		}
	}
}
