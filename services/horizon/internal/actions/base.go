package actions

import (
	"database/sql"
	"net/http"
	"time"

	horizonContext "github.com/stellar/go/services/horizon/internal/context"
	"github.com/stellar/go/services/horizon/internal/ledger"
	"github.com/stellar/go/services/horizon/internal/render"
	hProblem "github.com/stellar/go/services/horizon/internal/render/problem"
	"github.com/stellar/go/services/horizon/internal/render/sse"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/support/render/problem"
)

// Base is a helper struct you can use as part of a custom action via
// composition.
//
// TODO: example usage
type Base struct {
	W   http.ResponseWriter
	R   *http.Request
	Err error

	sseUpdateFrequency time.Duration
	isSetup            bool
}

// Prepare established the common attributes that get used in nearly every
// action.  "Child" actions may override this method to extend action, but it
// is advised you also call this implementation to maintain behavior.
func (base *Base) Prepare(w http.ResponseWriter, r *http.Request, sseUpdateFrequency time.Duration) {
	base.W = w
	base.R = r
	base.sseUpdateFrequency = sseUpdateFrequency
}

// Execute trigger content negotiation and the actual execution of one of the
// action's handlers.
func (base *Base) Execute(action interface{}) {
	ctx := base.R.Context()
	contentType := render.Negotiate(base.R)

	switch contentType {
	case render.MimeHal, render.MimeJSON:
		action, ok := action.(JSON)

		if !ok {
			goto NotAcceptable
		}

		action.JSON()

		if base.Err != nil {
			problem.Render(ctx, base.W, base.Err)
			return
		}

	case render.MimeEventStream:
		var pumped chan interface{}

		action, ok := action.(SSE)
		if !ok {
			goto NotAcceptable
		}

<<<<<<< HEAD
		// Subscribe this handler to the topic if the SSE request is related to a specific topic (tx_id, account_id, etc.).
		// This causes action.SSE to only be triggered by this topic. Unsubscribe when done.
		topic := action.GetTopic()
		if topic != "" {
			pumped = sse.Subscribe(topic)
			defer sse.Unsubscribe(pumped, topic)
		}

		stream := sse.NewStream(base.Ctx, base.W, base.R)
=======
		stream := sse.NewStream(ctx, base.W)
>>>>>>> horizon-v0.15.3

		for {
			lastLedgerState := ledger.CurrentState()

			// Rate limit the request if it's a call to stream since it queries the DB every second. See
			// https://github.com/stellar/go/issues/715 for more details.
			app := base.R.Context().Value(&horizonContext.AppContextKey)
			rateLimiter := app.(RateLimiterProvider).GetRateLimiter()
			if rateLimiter != nil {
				limited, _, err := rateLimiter.RateLimiter.RateLimit(rateLimiter.VaryBy.Key(base.R), 1)
				if err != nil {
					log.Ctx(ctx).Error(errors.Wrap(err, "RateLimiter error"))
					stream.Err(errors.New("Unexpected stream error"))
					return
				}
				if limited {
					stream.Err(errors.New("rate limit exceeded"))
					return
				}
			}

			action.SSE(stream)

			if base.Err != nil {
				// In the case that we haven't yet sent an event, is also means we
				// haven't sent the preamble, meaning we should simply return the normal HTTP
				// error.
				if stream.SentCount() == 0 {
					problem.Render(ctx, base.W, base.Err)
					return
				}

				if errors.Cause(base.Err) == sql.ErrNoRows {
					base.Err = errors.New("Object not found")
				} else {
					log.Ctx(ctx).Error(base.Err)
					base.Err = errors.New("Unexpected stream error")
				}

				// Send errors through the stream and then close the stream.
				stream.Err(base.Err)
			}

			// Manually send the preamble in case there are no data events in SSE to trigger a stream.Send call.
			// This method is called every iteration of the loop, but is protected by a sync.Once variable so it's
			// only executed once.
			stream.Init()

			if stream.IsDone() {
				return
			}

			// Make sure this is buffered channel of size 1. Otherwise, the go routine below
			// will never return if `newLedgers` channel is not read. From Effective Go:
			// > If the channel is unbuffered, the sender blocks until the receiver has received the value.
			newLedgers := make(chan bool, 1)
			go func() {
				for {
					time.Sleep(base.sseUpdateFrequency)
					currentLedgerState := ledger.CurrentState()
					if currentLedgerState.HistoryLatest >= lastLedgerState.HistoryLatest+1 {
						newLedgers <- true
						return
					}
				}
			}()

			select {
			case <-ctx.Done():
				stream.Done()
				return
<<<<<<< HEAD
			case <-pumped:
				//no-op, continue onto the next iteration
=======
			case <-newLedgers:
				continue
>>>>>>> horizon-v0.15.3
			}
		}
	case render.MimeRaw:
		action, ok := action.(Raw)

		if !ok {
			goto NotAcceptable
		}

		action.Raw()

		if base.Err != nil {
			problem.Render(ctx, base.W, base.Err)
			return
		}
	default:
		goto NotAcceptable
	}
	return

NotAcceptable:
	problem.Render(ctx, base.W, hProblem.NotAcceptable)
	return
}

// Do executes the provided func iff there is no current error for the action.
// Provides a nicer way to invoke a set of steps that each may set `action.Err`
// during execution
func (base *Base) Do(fns ...func()) {
	for _, fn := range fns {
		if base.Err != nil {
			return
		}

		fn()
	}
}

// Setup runs the provided funcs if and only if no call to Setup() has been
// made previously on this action.
func (base *Base) Setup(fns ...func()) {
	if base.isSetup {
		return
	}
	base.Do(fns...)
	base.isSetup = true
}
