package logger

import (
"context"
"log/slog"
)

type contextKey struct{}

// ContextWithLogger stores a *slog.Logger enriched with request-scoped fields
// into the context. Must be called once per request in the HTTP adapter layer.
func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, l)
}

// FromCtx retrieves the *slog.Logger from the context. Falls back to
// slog.Default() if no logger has been stored in the context.
func FromCtx(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}
