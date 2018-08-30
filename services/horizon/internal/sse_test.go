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
