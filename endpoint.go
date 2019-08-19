package liguetaxi

import (
	"context"
	"fmt"
)

// Request return types.
const (
	Json	= `json`
	Xml	= `xml`
)

// ContextKey is just an empty struct. It exists so ResType
// can be used as unique key for context.
type contextKey struct{}

// ResType is the context key to use with context's
// WithValue function to set the type of request's
// payload that will be received from the API.
var ResType contextKey

type endpoint string

// ContextType returns the type of API response (JSON or XML)
// associated with the Context.
func (e endpoint) ContextType(ctx context.Context) string {
	suffix := Json
	if ctx != nil {
		if t, ok := ctx.Value(ResType).(string); ok && t != "" {
			suffix = t
		}
	}
	return suffix
}

// String reads the Context and returns the endpoint suffixed with the type
// of the request: json or xml.
func (e endpoint) String(ctx context.Context) string {
	return fmt.Sprintf("%s/%s", e, e.ContextType(ctx))
}
