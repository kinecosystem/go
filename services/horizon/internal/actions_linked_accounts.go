package horizon

import (
	"github.com/kinecosystem/go/services/horizon/internal/actions"
	"github.com/kinecosystem/go/services/horizon/internal/render/sse"
	"github.com/kinecosystem/go/support/render/hal"
)

// Interface verifications
var _ actions.JSONer = (*AggregateBalanceAction)(nil)
var _ actions.RawDataResponder = (*AggregateBalanceAction)(nil)
var _ actions.EventStreamer = (*AggregateBalanceAction)(nil)

// AggregateBalanceAction renders a account summary found by its address.
type AggregateBalanceAction struct {
	Action
	Address string
	Balance string
}

// JSON is a method for actions.JSON
func (action *AggregateBalanceAction) JSON() error {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {hal.Render(action.W, action.Balance)})
	return action.Err
}

// Raw is a method for actions.Raw
func (action *AggregateBalanceAction) Raw() error {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {hal.Render(action.W, action.Balance)})
	return action.Err
}

// SSE is a method for actions.SSE
func (action *AggregateBalanceAction) SSE(stream *sse.Stream) error {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			stream.Send(sse.Event{Data: action.Balance})
		},
	)
	return action.Err
}

// GetPubsubTopic is a method for actions.SSE
//
// There is no value in this action for specific account_id, so registration topic is a general
// change in the ledger.
func (action *AggregateBalanceAction) GetPubsubTopic() string {
	return action.GetString("account_id")
}

func (action *AggregateBalanceAction) loadParams() {
	action.Address = action.GetAddress("account_id", actions.RequiredParam)
}

func (action *AggregateBalanceAction) loadRecord() {
	 action.Balance, action.Err = action.CoreQ().AggregateBalanceByAccountId(&action.Balance, action.Address)
}
