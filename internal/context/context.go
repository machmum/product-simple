package ctxin

import (
	"context"
)

const RequestID = "id"

type ReqIDContextKey string

// NewRequestIDContext returns a new Context carrying requestID.
func NewRequestIDContext(ctx context.Context, requestID string) context.Context {
	k := ReqIDContextKey(RequestID)
	return context.WithValue(ctx, k, requestID)
}

// FromRequestIDContext extracts the user requestID from ctx, if present.
func FromRequestIDContext(ctx context.Context) (value string) {
	if ctx == nil {
		return
	}

	k := ReqIDContextKey(RequestID)
	if v := ctx.Value(k); v != nil {
		value = v.(string)
	}
	return
}

const Skip = "skip"

type SkipContextKey string

// NewSkipContext returns a new Context carrying skip.
func NewSkipContext(ctx context.Context, skip int) context.Context {
	k := SkipContextKey(Skip)
	return context.WithValue(ctx, k, skip)
}

// FromSkipContext extracts the user skip from ctx, if present.
func FromSkipContext(ctx context.Context) (value int) {
	if ctx == nil {
		return
	}

	k := SkipContextKey(Skip)
	if v := ctx.Value(k); v != nil {
		value = v.(int)
	}
	return
}
