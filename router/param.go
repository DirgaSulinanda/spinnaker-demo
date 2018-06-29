package router

import "context"

type param string

// Param returns param p for the context.
func Param(ctx context.Context, p string) string {
	return ctx.Value(param(p)).(string)
}

// WithParam returns a new context with param p set to v.
func WithParam(ctx context.Context, p, v string) context.Context {
	return context.WithValue(ctx, param(p), v)
}
