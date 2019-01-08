package horizon

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kinecosystem/go/network"
	"github.com/kinecosystem/go/services/horizon/internal/ingest"
	"github.com/kinecosystem/go/services/horizon/internal/ledger"
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"github.com/kinecosystem/go/services/horizon/internal/test"
)

// Test subscription gets updates, and sse.Tick() doesnt trigger channel.
func TestSSEPubsub(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	subscription := sse.Subscribe("a")
	defer sse.Unsubscribe(subscription, "a")


	var wg sync.WaitGroup
	wg.Add(1)
	go func(subscription chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		select {
		case <-subscription:
			return
		case <-time.After(2 * time.Second):
			t.Fatal("subscription did not trigger")
		}
	}(subscription, &wg)
	sse.Publish("a")
	wg.Wait()

	sse.Tick()
	select {
	case <-subscription:
		t.Fatal("subscription shouldn't trigger after tick")
	case <-time.After(2 * time.Second): // no-op. Success!
	}

}

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
	time.Sleep(10 * time.Millisecond)
	sse.Publish("a")
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
			case <-time.After(2 * time.Second):
				t.Fatal("Subscription did not trigger within 2 seconds")
			}
		}(subscriptions[i], &wg)
	}

	for i := 0; i < 100; i++ {
		sse.Publish(strconv.Itoa(i))
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
			case <-time.After(2 * time.Second):
				t.Fatal("Subscription did not trigger within 2 seconds")
			}
		}(subscriptions[i], &wg)
	}
	time.Sleep(10 * time.Millisecond)

	sse.Publish("a")

	wg.Wait()
}

// Test SSE subscription get message when ingest to Horizon happens.
func TestSSEPubsubTransactions(t *testing.T) {
	SCENARIO_NAME := "kahuna"
	TX_HASH := "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB"

	tt := test.Start(t).ScenarioWithoutHorizon(SCENARIO_NAME)
	defer tt.Finish()

	subscription := sse.Subscribe(TX_HASH)
	defer sse.Unsubscribe(subscription, TX_HASH)

	var wg sync.WaitGroup
	wg.Add(10)

	go func(subscription chan interface{}, wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			select {
			case <-subscription:
				wg.Done()
			case <-time.After(10 * time.Second):
				t.Fatal("subscription did not trigger within fast enough")
			}
		}
	}(subscription, &wg)

	ingestHorizon(tt)

	wg.Wait()
}

// Helpers from ingest/main_test.go

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
	)
}
