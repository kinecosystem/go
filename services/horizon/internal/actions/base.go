package actions

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	horizonContext "github.com/kinecosystem/go/services/horizon/internal/context"
	"github.com/kinecosystem/go/services/horizon/internal/render"
	hProblem "github.com/kinecosystem/go/services/horizon/internal/render/problem"
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"github.com/kinecosystem/go/support/errors"
	"github.com/kinecosystem/go/support/render/problem"
)

// Base is a helper struct you can use as part of a custom action via
// composition.
//
// TODO: example usage
type Base struct {
	W   http.ResponseWriter
	R   *http.Request
	Err error

	appCtx  context.Context
	isSetup bool
}

// Prepare established the common attributes that get used in nearly every
// action.  "Child" actions may override this method to extend action, but it
// is advised you also call this implementation to maintain behavior.
func (base *Base) Prepare(w http.ResponseWriter, r *http.Request, appCtx context.Context) {
	base.W = w
	base.R = r
	base.appCtx = appCtx
}

// Execute trigger content negotiation and the actual execution of one of the
// action's handlers.
func (base *Base) Execute(action interface{}) {
	ctx := base.R.Context()
	contentType := render.Negotiate(base.R)

	switch contentType {
	case render.MimeHal, render.MimeJSON:
		action, ok := action.(JSONer)
		if !ok {
			goto NotAcceptable
		}

		err := action.JSON()
		if err != nil {
			problem.Render(ctx, base.W, err)
			return
		}

	case render.MimeEventStream:
		var notification chan interface{}

		switch ac := action.(type) {
		case EventStreamer:
			// Subscribe this handler to the topic if the SSE request is related to a specific topic (tx_id, account_id, etc.).
			// This causes action.SSE to only be triggered by this topic. Unsubscribe when done.
			topic := ac.GetPubsubTopic()
			if topic != "" {
				notification = sse.Subscribe(topic)
				defer sse.Unsubscribe(notification, topic)
			}
		case SingleObjectStreamer:
		default:
			goto NotAcceptable
		}

		stream := sse.NewStream(ctx, base.W)

		var oldHash [32]byte
		for {
			// Rate limit the request if it's a call to stream since it queries the DB every second. See
			// https://github.com/stellar/go/issues/715 for more details.
			app := base.R.Context().Value(&horizonContext.AppContextKey)
			rateLimiter := app.(RateLimiterProvider).GetRateLimiter()
			if rateLimiter != nil {
				limited, _, err := rateLimiter.RateLimiter.RateLimit(rateLimiter.VaryBy.Key(base.R), 1)
				if err != nil {
					stream.Err(errors.Wrap(err, "RateLimiter error"))
					return
				}
				if limited {
					stream.Err(sse.ErrRateLimited)
					return
				}
			}

			switch ac := action.(type) {
			case EventStreamer:
				err := ac.SSE(stream)
				if err != nil {
					stream.Err(err)
					return
				}

			case SingleObjectStreamer:
				newEvent, err := ac.LoadEvent()
				if err != nil {
					stream.Err(err)
					return
				}
				resource, err := json.Marshal(newEvent.Data)
				if err != nil {
					stream.Err(errors.Wrap(err, "unable to marshal next action resource"))
					return
				}

				nextHash := sha256.Sum256(resource)
				if bytes.Equal(nextHash[:], oldHash[:]) {
					break
				}

				oldHash = nextHash
				stream.SetLimit(10)
				stream.Send(newEvent)
			}

			// Manually send the preamble in case there are no data events in SSE to trigger a stream.Send call.
			// This method is called every iteration of the loop, but is protected by a sync.Once variable so it's
			// only executed once.
			stream.Init()

			if stream.IsDone() {
				return
			}

			select {
			case <-notification:
				// No-op, continue onto the next iteration.
				continue
			case <-ctx.Done():
			case <-base.appCtx.Done():
			}

			stream.Done()
			return
		}
	case render.MimeRaw:
		action, ok := action.(RawDataResponder)
		if !ok {
			goto NotAcceptable
		}

		err := action.Raw()
		if err != nil {
			problem.Render(ctx, base.W, err)
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
