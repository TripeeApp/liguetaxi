package liguetaxi

import (
	"context"
	"fmt"
)

// ContextKey is just an empty struct. It exists so ResType
// can be used as unique key for context.
type contextKey struct{}

// ResType is the context key to use with context's
// WithValue function to set the type of request's
// payload that will be received from the API.
var ResType contextKey

// Request return types.
const (
	Json	= `json`
	Xml	= `xml`
)

type endpoint string

// String returns the endpoint suffixed with the type
// of the request: json or xml.
// The suffix is read from the context. If there's
// no suffix set in context, defaults to `json`.
func (e *endpoint) String(ctx context.Context) string {
	suffix := Json
	if t, ok := ctx.Value(ResType).(string); ok && t != "" {
		suffix = t
	}
	return fmt.Sprintf("%s/%s", *e, suffix)
}
