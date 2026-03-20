package rest

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/charlires/go-base-backend-service/internal/adapters/input/rest/gen"
	"github.com/charlires/go-base-backend-service/internal/pkg/logger"
)

// NewLoggingMiddleware returns a gen.StrictMiddlewareFunc that enriches the
// request-scoped logger with the following fields and stores it in the context:
//
//   - operation_id — the oapi-codegen operation name (e.g. "GetUserById")
//   - request_id   — value of X-Request-ID header, or a generated UUID if absent
//   - method       — HTTP method (e.g. "GET")
//   - path         — URL path (e.g. "/users/123")
//
// Downstream layers retrieve the logger via logger.FromCtx(ctx).
func NewLoggingMiddleware(base *slog.Logger) gen.StrictMiddlewareFunc {
	return func(next gen.StrictHandlerFunc, operationID string) gen.StrictHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			enriched := base.With(
				slog.String("operation_id", operationID),
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			ctx = logger.ContextWithLogger(ctx, enriched)

			return next(ctx, w, r, request)
		}
	}
}
