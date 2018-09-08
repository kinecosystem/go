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

// Check subscription gets updates, and sse.Tick() doesnt trigger channel.
func TestSSEPubsub(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	subscription := sse.Subscribe("a")
	defer sse.Unsubscribe(subscription, "a")

	sse.Publish("a")

	select {
	case <-subscription:
		// no-op.  Success!
	case <-time.After(2 * time.Second):
		t.Error("channel did not trigger")
		t.FailNow()
	}

	sse.Tick()
	select {
	case <-subscription:
		t.Error("channel shouldn't trigger after tick")
		t.FailNow()
	case <-time.After(2 * time.Second):
		// no-op. Success!
	}

}

// Check 2 subscriptions to different topics. Make sure that one topic doesnt
// raise both channels.
func TestSSEPubsubMultipleChannels(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	channelA := sse.Subscribe("a")
	channelB := sse.Subscribe("b")
	defer sse.Unsubscribe(channelA, "a")
	defer sse.Unsubscribe(channelB, "b")

	var wg sync.WaitGroup
	wg.Add(1)
	go func(channelA chan interface{}, channelB chan interface{}, wg *sync.WaitGroup) {
		defer wg.Done()
		select {
		case <-channelA:
			select {
			case <-channelA:
				t.Error("channelA shouldnt trigger")
				t.FailNow()
			case <-channelB:
				t.Error("channelB shouldnt trigger")
				t.FailNow()
			case <-time.After(2 * time.Second):
				// no-op. Success!
			}
		case <-channelB:
			t.Error("channelB shouldnt trigger")
			t.FailNow()
		case <-time.After(2 * time.Second):
			t.Error("channelA did not trigger")
			t.FailNow()
		}
	}(channelA, channelB, &wg)

	sse.Publish("a")
	wg.Wait()
}

// Check multiple number of topics handled
func TestSSEPubsubManyChannels(t *testing.T) {
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
				t.Error("Subscription did not trigger within 2 seconds")
				t.FailNow()
			}
		}(subscriptions[i], &wg)
	}

	for j := 0; j < 100; j++ {
		sse.Publish(strconv.Itoa(j))
	}

	wg.Wait()
}

// Check multiple subscriptions to the same topic.
func TestSSEPubsubManyListeners(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	var wg sync.WaitGroup
	wg.Add(100)
	subscriptions := make([]chan interface{}, 100)
	//var subscriptions [100]chan interface{}
	for i := 0; i < 100; i++ {
		subscriptions[i] = sse.Subscribe("a")
		defer sse.Unsubscribe(subscriptions[i], "a")
		go func(subscription chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			select {
			case <-subscription:
				return
			case <-time.After(2 * time.Second):
				t.Error("Subscription did not trigger within 2 seconds")
				t.FailNow()
			}
		}(subscriptions[i], &wg)
	}

	sse.Publish("a")

	wg.Wait()
}

// Check sse subscription get message when ingest to Horizon happens.
func TestSSEPubsubTransactions(t *testing.T) {
	tt := test.Start(t).ScenarioWithoutHorizon("kahuna")
	defer tt.Finish()
	subscription := sse.Subscribe("GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB")
	defer sse.Unsubscribe(subscription, "GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB")
	var wg sync.WaitGroup
	wg.Add(10)

	go func(subscription chan interface{}, wg *sync.WaitGroup) {
		for i := 0; i < 10; i++ {
			select {
			case <-subscription:
				wg.Done()
			case <-time.After(10 * time.Second):
				t.Error("Subscription did not trigger within 10 seconds")
				t.FailNow()
			}
		}
	}(subscription, &wg)
	ingestHorizon(tt)

	wg.Wait()
}

// helpers from ingest/main_test.go
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
