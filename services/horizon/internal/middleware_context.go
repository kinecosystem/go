package horizon

import (
	"context"
	"net/http"

	"github.com/kinecosystem/go/services/horizon/internal/httpx"
	"github.com/kinecosystem/go/support/context/requestid"
)

func contextMiddleware(parent context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = requestid.ContextFromChi(ctx)
			ctx, cancel := httpx.RequestContext(ctx, w, r)

			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
