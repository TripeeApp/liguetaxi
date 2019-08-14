package liguetaxi

import "context"

type ResponseTypeKey struct{}

type Endpoint string

func (e *Endpoint) String(ctx context.Context) string {
	return "/test/json"
}
