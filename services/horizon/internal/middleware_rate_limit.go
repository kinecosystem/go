package horizon

import (
	"net/http"
)

func (web *Web) RateLimitMiddleware(next http.Handler) http.Handler {
	if web.rateLimiter == nil {
		return next
	}
	return web.rateLimiter.RateLimit(next)
}
