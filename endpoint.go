package liguetaxi

import (
	"context"
	"fmt"
)

type ResponseTypeKey struct{}

type Endpoint string

// String returns the endpoint suffixed with the type
// of the request: json or xml.
func (e *Endpoint) String(ctx context.Context) string {
	return fmt.Sprintf("%s/%s", *e, `json`)
}
