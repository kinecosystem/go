package horizon

import (
	"context"
	"net/http"

<<<<<<< HEAD
	gctx "github.com/goji/context"
	"github.com/kinecosystem/go/services/horizon/internal/context/requestid"
	"github.com/kinecosystem/go/services/horizon/internal/httpx"
	"github.com/zenazn/goji/web"
	"golang.org/x/net/context"
=======
	"github.com/stellar/go/services/horizon/internal/httpx"
	"github.com/stellar/go/support/context/requestid"
>>>>>>> horizon-v0.15.3
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
