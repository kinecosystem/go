// Package hchi provides functions to support embedded and retrieving
// a request id from a go context tree
package hchi

import (
	"context"

	"github.com/go-chi/chi/middleware"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

const (
	reqidKey key = iota
)

// WithRequestID create a context from the provided parent and the provided request id string.
func WithRequestID(ctx context.Context, reqid string) context.Context {
	return context.WithValue(ctx, reqidKey, reqid)
}

// WithChiRequestID returns a new context bound with the value of the request id.
func WithChiRequestID(ctx context.Context) context.Context {
	reqid := middleware.GetReqID(ctx)
	return WithRequestID(ctx, reqid)
}

// RequestID returns the set request id, if one has been set, from the
// provided context returns "" if no requestid has been set
func RequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	val, _ := ctx.Value(reqidKey).(string)
	return val
}
