package problem

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stellar/go/services/horizon/internal/hchi"
	"github.com/stellar/go/support/render/problem"
	"github.com/stretchr/testify/assert"
)

var ctx = context.Background()
var testRender = func(ctx context.Context, p interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	problem.Render(ctx, w, p)
	return w
}

func TestCommonProblems(t *testing.T) {
	testCases := []struct {
		testName     string
		p            problem.P
		expectedCode int
	}{
		{"NotFound", problem.NotFound, 404},
		{"RateLimitExceeded", RateLimitExceeded, 429},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			w := testRender(ctx, tc.p)
			assert.Equal(t, tc.expectedCode, w.Code)
		})
	}
}

func TestInflate(t *testing.T) {
	// Sets Instance to the request id based on the context
	ctx2 := hchi.WithRequestID(ctx, "2")
	p := problem.P{}

	Inflate(ctx2, &p)
	assert.Equal(t, "2", p.Instance)

	// when no request id is set, instance should be ""
	Inflate(ctx, &p)
	assert.Equal(t, "", p.Instance)
}
