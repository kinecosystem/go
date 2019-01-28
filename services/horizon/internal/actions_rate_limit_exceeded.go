package horizon

import (
	"net/http"

<<<<<<< HEAD
	"github.com/zenazn/goji/web"

	hProblem "github.com/kinecosystem/go/services/horizon/internal/render/problem"
	"github.com/kinecosystem/go/support/render/problem"
=======
	hProblem "github.com/stellar/go/services/horizon/internal/render/problem"
	"github.com/stellar/go/support/render/problem"
>>>>>>> horizon-v0.15.3
)

// RateLimitExceededAction renders a 429 response
type RateLimitExceededAction struct {
	Action
	App *App
}

func (action RateLimitExceededAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ap := &action.Action
	ap.Prepare(w, r)
	ap.App = action.App
	problem.Render(action.R.Context(), action.W, hProblem.RateLimitExceeded)
}
