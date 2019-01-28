package horizon

import (
	"compress/flate"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/rcrowley/go-metrics"
	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/stellar/go/services/horizon/internal/db2"
	hProblem "github.com/stellar/go/services/horizon/internal/render/problem"
	"github.com/stellar/go/services/horizon/internal/txsub/sequence"
	"github.com/stellar/go/support/render/problem"
	"github.com/throttled/throttled"
)

// Web contains the http server related fields for horizon: the router,
// rate limiter, etc.
type Web struct {
	router      *chi.Mux
	rateLimiter *throttled.HTTPRateLimiter

	requestTimer metrics.Timer
	failureMeter metrics.Meter
	successMeter metrics.Meter
}

// initWeb installed a new Web instance onto the provided app object.
func initWeb(app *App) {
	app.web = &Web{
		router:       chi.NewRouter(),
		requestTimer: metrics.NewTimer(),
		failureMeter: metrics.NewMeter(),
		successMeter: metrics.NewMeter(),
	}

	// register problems
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
	problem.RegisterError(sequence.ErrNoMoreRoom, hProblem.ServerOverCapacity)
	problem.RegisterError(db2.ErrInvalidCursor, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidLimit, problem.BadRequest)
	problem.RegisterError(db2.ErrInvalidOrder, problem.BadRequest)
}

// initWebMiddleware installs the middleware stack used for horizon onto the
// provided app.
func initWebMiddleware(app *App) {

	r := app.web.router
	r.Use(chimiddleware.Timeout(app.config.ConnectionTimeout))
	r.Use(chimiddleware.StripSlashes)
	r.Use(app.Middleware)
	r.Use(requestCacheHeadersMiddleware)
	r.Use(chimiddleware.RequestID)
	r.Use(contextMiddleware(app.ctx))
	r.Use(xff.Handler)
	r.Use(LoggerMiddleware)
	r.Use(requestMetricsMiddleware)
	r.Use(RecoverMiddleware)
	r.Use(chimiddleware.Compress(flate.DefaultCompression, "application/hal+json"))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	r.Use(c.Handler)

	r.Use(app.web.RateLimitMiddleware)
}

// initWebActions installs the routing configuration of horizon onto the
// provided app.  All route registration should be implemented here.
func initWebActions(app *App) {

	r := app.web.router
	r.Get("/", RootAction{}.Handle)
	r.Get("/metrics", MetricsAction{}.Handle)

	// ledger actions
	r.Route("/ledgers", func(r chi.Router) {
		r.Get("/", LedgerIndexAction{}.Handle)
		r.Route("/{ledger_id}", func(r chi.Router) {
			r.Get("/", LedgerShowAction{}.Handle)
			r.Get("/transactions", TransactionIndexAction{}.Handle)
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", PaymentsIndexAction{}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
		})
	})

	// account actions
	r.Route("/accounts", func(r chi.Router) {
		r.Route("/{account_id}", func(r chi.Router) {
			r.Get("/", AccountShowAction{}.Handle)
			r.Get("/transactions", TransactionIndexAction{}.Handle)
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", PaymentsIndexAction{}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
			r.Get("/offers", OffersByAccountAction{}.Handle)
			r.Get("/trades", TradeIndexAction{}.Handle)
			r.Get("/data/{key}", DataShowAction{}.Handle)
		})
	})

	// transaction history actions
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", TransactionIndexAction{}.Handle)
		r.Route("/{tx_id}", func(r chi.Router) {
			r.Get("/", TransactionShowAction{}.Handle)
			r.Get("/operations", OperationIndexAction{}.Handle)
			r.Get("/payments", PaymentsIndexAction{}.Handle)
			r.Get("/effects", EffectIndexAction{}.Handle)
		})
	})

	// operation actions
	r.Route("/operations", func(r chi.Router) {
		r.Get("/", OperationIndexAction{}.Handle)
		r.Get("/{id}", OperationShowAction{}.Handle)
		r.Get("/{op_id}/effects", EffectIndexAction{}.Handle)
	})

	// payment actions
	r.Get("/payments", PaymentsIndexAction{}.Handle)

	// effect actions
	r.Get("/effects", EffectIndexAction{}.Handle)

	// trading related endpoints
	r.Get("/trades", TradeIndexAction{}.Handle)
	r.Get("/trade_aggregations", TradeAggregateIndexAction{}.Handle)
	r.Route("/offers", func(r chi.Router) {
		r.Get("/{id}", NotImplementedAction{}.Handle)
		r.Get("/{offer_id}/trades", TradeIndexAction{}.Handle)
	})
	r.Get("/order_book", OrderBookShowAction{}.Handle)

	// Transaction submission API
	r.Post("/transactions", TransactionCreateAction{}.Handle)
	r.Get("/paths", PathIndexAction{}.Handle)

	if app.config.EnableAssetStats {
		// Asset related endpoints
		r.Get("/assets", AssetsAction{}.Handle)
	}

	// Network state related endpoints
	r.Get("/operation_fee_stats", OperationFeeStatsAction{}.Handle)

	// friendbot
	if app.config.FriendbotURL != nil {
		redirectFriendbot := func(w http.ResponseWriter, r *http.Request) {
			redirectURL := app.config.FriendbotURL.String() + "?" + r.URL.RawQuery
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		}
		r.Post("/friendbot", redirectFriendbot)
		r.Get("/friendbot", redirectFriendbot)
	}

	r.NotFound(NotFoundAction{}.Handle)
}

func initWebRateLimiter(app *App) {
	// Disabled
	if app.config.RateLimit == nil {
		return
	}

	rateLimiter, err := throttled.NewGCRARateLimiter(50000, *app.config.RateLimit)
	if err != nil {
		panic(fmt.Errorf("unable to create RateLimiter"))
	}

	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter:   rateLimiter,
		DeniedHandler: &RateLimitExceededAction{App: app, Action: Action{}},
	}
	httpRateLimiter.VaryBy = VaryByRemoteIP{}
	app.web.rateLimiter = &httpRateLimiter
}

type VaryByRemoteIP struct{}

func (v VaryByRemoteIP) Key(r *http.Request) string {
	return remoteAddrIP(r)
}

func remoteAddrIP(r *http.Request) string {
	// To support IPv6
	lastSemicolon := strings.LastIndex(r.RemoteAddr, ":")
	if lastSemicolon == -1 {
		return r.RemoteAddr
	} else {
		return r.RemoteAddr[0:lastSemicolon]
	}
}

func firstXForwardedFor(r *http.Request) string {
	return strings.TrimSpace(strings.SplitN(r.Header.Get("X-Forwarded-For"), ",", 2)[0])
}

func init() {
	appInit.Add(
		"web.init",
		initWeb,

		"app-context",
	)

	appInit.Add(
		"web.rate-limiter",
		initWebRateLimiter,

		"web.init",
	)
	appInit.Add(
		"web.middleware",
		initWebMiddleware,

		"web.init",
		"web.rate-limiter",
	)
	appInit.Add(
		"web.actions",
		initWebActions,

		"web.init",
	)
}
