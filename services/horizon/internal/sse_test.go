package horizon

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/kinecosystem/go/network"
	"github.com/kinecosystem/go/services/horizon/internal/actions"
	"github.com/kinecosystem/go/services/horizon/internal/ingest"
	"github.com/kinecosystem/go/services/horizon/internal/ledger"
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"github.com/kinecosystem/go/services/horizon/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test 2 subscriptions to different topics. Make sure that one topic doesnt
// raise both channels.
func TestSSEPubsubMultipleChannels(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	subA := sse.Subscribe("a")
	subB := sse.Subscribe("b")
	defer sse.Unsubscribe(subA, "a")
	defer sse.Unsubscribe(subB, "b")

	var wg sync.WaitGroup
	wg.Add(1)
	go func(subA chan interface{}, subB chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()

		select {
		case <-subA: // no-op. Success!
		case <-subB:
			t.Fatal("subscription B shouldn't trigger")
		case <-time.After(2 * time.Second):
			t.Fatal("subscription A did not trigger")
		}

		select {
		case <-subA:
			t.Fatal("subscription A shouldn't trigger")
		case <-subB:
			t.Fatal("subscription B shouldn't trigger")
		case <-time.After(2 * time.Second): // no-op. Success!
		}
	}(subA, subB, &wg)
	sse.Publish("a", true)
	wg.Wait()
}

// Test multiple number of topics handled
func TestSSEPubsubManyTopics(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	var wg sync.WaitGroup
	wg.Add(100)
	subscriptions := make([]chan interface{}, 100)

	for i := 0; i < 100; i++ {
		subscriptions[i] = sse.Subscribe(strconv.Itoa(i))
		defer sse.Unsubscribe(subscriptions[i], strconv.Itoa(i))

		go func(subscription chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			select {
			case <-subscription:
				return
			case <-time.After(10 * time.Second):
				t.Fatal("Subscription did not trigger within 2 seconds")
			}
		}(subscriptions[i], &wg)
	}

	for i := 0; i < 100; i++ {
		sse.Publish(strconv.Itoa(i), true)
	}

	wg.Wait()
}

// Test multiple subscriptions to the same topic.
func TestSSEPubsubManySubscribers(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	var wg sync.WaitGroup
	wg.Add(100)
	subscriptions := make([]chan interface{}, 100)
	for i := 0; i < 100; i++ {
		subscriptions[i] = sse.Subscribe("a")
		defer sse.Unsubscribe(subscriptions[i], "a")

		go func(subscription chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			select {
			case <-subscription:
				return
			case <-time.After(10 * time.Second):
				t.Fatal("Subscription did not trigger within 2 seconds")
			}
		}(subscriptions[i], &wg)
	}
	sse.Publish("a", true)

	wg.Wait()
}

// Test Actions which implement actions.EventStreamer return the correct URL path parameter value
// which is used for SSE PubSub subscription.
//
// Actions and parameters are defined in services/horizon/internal/init_web.go
func TestSSEPubsubTopic(t *testing.T) {
	for _, tt := range []struct {
		path,
		paramKey,
		paramValue string
		streamEventType actions.EventStreamer
	}{
		// Constructs a URL using the structure as in init_web.go
		// e.g. "/accounts/{account_id=12345}" ==> GetTopic() should return "12345"

		// Unspecific topics.
		{"/ledgers/", "ledger", "ledger", &LedgerIndexAction{}},
		{"/order_book/", "order_book", "order_book", &OrderBookShowAction{}},

		// Topics using a specific parameter e.g. and account address.
		//
		// NOTE the parameter we compare (e.g. the account address) doesn't have to be a valid
		// address for the purpose of this test. We only care the SSE PubSub implementation fetches
		// the correct parameter. The validity is already checked elsewhere and is out of the scope
		// of SSE handling in general.
		{"/accounts/", "account_id", "1", &AccountShowAction{}},
		{"/transactions/", "account_id", "2", &TransactionIndexAction{}},
		{"/operations/", "account_id", "3", &OperationIndexAction{}},
		{"/payments/", "account_id", "4", &PaymentsIndexAction{}},
		{"/effects/", "account_id", "5", &EffectIndexAction{}},
		{"/offers/", "account_id", "6", &OffersByAccountAction{}},
		{"/trades/", "account_id", "7", &TradeIndexAction{}},

		// "2nd-level" parameters e.g. /accounts/{account_id}/data/{key}
		{"/data/", "account_id", "8", &DataShowAction{}},
	} {
		// Initialize the Action implemented by the current actions.EventStreamer
		action := reflect.ValueOf(tt.streamEventType).Elem()
		base := action.FieldByName("Base")
		require.True(t, base.IsValid()) // Sanity check: Base should be a struct member of action
		base.Set(reflect.ValueOf(*makeAction(tt.path, map[string]string{tt.paramKey: tt.paramValue})))

		// Test that it's EventStreamer.GetTopic() returns the correct expected PubSub topic
		assert.Equal(t, tt.paramValue, tt.streamEventType.GetTopic(), "Path: %s/%s", tt.path, "/", tt.paramValue)
	}
}

// Test various SSE subscriptions to various topics receive notifications when an ingestion
// that relates to these topics occur.
func TestSSEPubsubSubscribeToTopics(t *testing.T) {
	SCENARIO_NAME := "kahuna"

	tt := test.Start(t).ScenarioWithoutHorizon(SCENARIO_NAME)
	defer tt.Finish()

	var wg sync.WaitGroup
	for _, topic := range []string{
		"transactions",
		// account in "kahuna" memo test
		"GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB",
	} {
		subscription := sse.Subscribe(topic)
		defer sse.Unsubscribe(subscription, topic)

		wg.Add(1)
		go func(topic string, subscription chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()

			select {
			case <-subscription:
			case <-time.After(10 * time.Second):
				t.Fatalf("subscription did not trigger within 10s for topic \"%s\"", topic)
			}
		}(topic, subscription, &wg)
	}

	ingestHorizon(tt)

	wg.Wait()
}

// Helpers from actions/helpers_test.go and ingest/main_test.go

func ingestHorizon(tt *test.T) *ingest.Session {
	sys := sys(tt)
	s := ingest.NewSession(sys)
	s.Cursor = ingest.NewCursor(1, ledger.CurrentState().CoreLatest, sys)
	s.Run()

	return s
}

func sys(tt *test.T) *ingest.System {
	return ingest.New(
		network.TestNetworkPassphrase,
		"",
		tt.CoreSession(),
		tt.HorizonSession(),
		"HORIZON",
		ingest.Config{},
	)
}

func makeAction(path string, body map[string]string) *actions.Base {
	rctx := chi.NewRouteContext()
	for k, v := range body {
		rctx.URLParams.Add(k, v)
	}

	r, _ := http.NewRequest("GET", path, nil)
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	action := &actions.Base{
		R: r,
	}
	return action
}
