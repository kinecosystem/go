package horizon

import (
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"testing"
	"time"
)

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

func TestSsePubsubMultipleChannels(t *testing.T) {
	ht := StartHTTPTest(t, "base")
	defer ht.Finish()

	ht.App.ticks.Stop()

	channel_a := sse.Subscribe("a")
	channel_b := sse.Subscribe("b")

	sse.Publish("a")

	success := false
	for {
		select {
		case <-channel_a:
			success = true
			// no-op.  Success!
		case <-channel_b:
			t.Error("channel_b shouldnt trigger")
			t.FailNow()
		case <-time.After(2 * time.Second):
			if !success {
				t.Error("channel_a did not trigger")
				t.FailNow()
			}
			return
		}
	}

}
