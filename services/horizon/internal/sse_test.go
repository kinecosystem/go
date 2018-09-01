package horizon

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kinecosystem/go/network"
	. "github.com/kinecosystem/go/services/horizon/internal/ingest"
	"github.com/kinecosystem/go/services/horizon/internal/ledger"
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"github.com/kinecosystem/go/services/horizon/internal/test"
)

// Check subscription gets updates, and sse.Tick() doesnt raise channel.
func TestSsePubsub(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	ch := sse.Subscribe("a")

	sse.Publish("a")

	select {
	case <-ch:
		// no-op.  Success!
	case <-time.After(2 * time.Second):
		t.Error("channel did not trigger")
		t.FailNow()
	}

	sse.Tick()
	select {
	case <-ch:
		t.Error("channel shouldn't trigger after tick")
		t.FailNow()
	case <-time.After(2 * time.Second):
		// no-op. Success!
	}

}

// Check 2 subscriptions to different topics. Make sure that one topic doesnt
// raise both channels.
func TestSsePubsubMultipleChannels(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	channelA := sse.Subscribe("a")
	channelB := sse.Subscribe("b")

	sse.Publish("a")

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
}

// Check multiple number of topics handled
func TestSsePubsubManyChannels(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	var wg sync.WaitGroup
	wg.Add(99)
	var chans [100]chan interface{}
	for i := 0; i < 99; i++ {
		chans[i] = sse.Subscribe(strconv.Itoa(i))
		go func(channel chan interface{}) {
			defer wg.Done()
			select {
			case <-channel:
				return
			}
		}(chans[i])
	}

	for j := 0; j < 99; j++ {
		sse.Publish(strconv.Itoa(j))
	}

	wg.Wait()
}

// Check multiple subscriptions to the same topic.
func TestSsePubsubManyListeners(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	var wg sync.WaitGroup
	wg.Add(99)
	var chans [100]chan interface{}
	for i := 0; i < 99; i++ {
		chans[i] = sse.Subscribe("a")
		go func(channel chan interface{}) {
			defer wg.Done()
			select {
			case <-channel:
				return
			}
		}(chans[i])
	}

	sse.Publish("a")

	wg.Wait()
}

// Check sse subscription get message when ingest to Horizon happens.
func TestSsePubsubTransactions(t *testing.T) {
	tt := test.Start(t).ScenarioWithoutHorizon("kahuna")
	defer tt.Finish()
	channel := sse.Subscribe("GA46VRKBCLI2X6DXLX7AIEVRFLH3UA7XBE3NGNP6O74HQ5LXHMGTV2JB")

	var wg sync.WaitGroup
	wg.Add(10)

	go func(channel chan interface{}) {
		for i := 0; i < 10; i++ {
			select {
			case <-channel:
				wg.Done()
			}
		}
	}(channel)
	ingestHorizon(tt)

	wg.Wait()
}

// helpers from ingest/main_test.go
func ingestHorizon(tt *test.T) *Session {
	sys := sys(tt)
	s := NewSession(sys)
	s.Cursor = NewCursor(1, ledger.CurrentState().CoreLatest, sys)
	s.Run()

	return s
}

func sys(tt *test.T) *System {
	return New(
		network.TestNetworkPassphrase,
		"",
		tt.CoreSession(),
		tt.HorizonSession(),
	)
}
