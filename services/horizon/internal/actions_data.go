package horizon

import (
	"github.com/stellar/go/services/horizon/internal/actions"
	"github.com/stellar/go/services/horizon/internal/db2/core"
	"github.com/stellar/go/services/horizon/internal/render/sse"
	"github.com/stellar/go/support/render/hal"
)

// DataShowAction renders a account summary found by its address.
type DataShowAction struct {
	Action
	Address string
	Key     string
	Data    core.AccountData
}

// JSON is a method for actions.JSON
func (action *DataShowAction) JSON() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {

			hal.Render(action.W, map[string]string{
				"value": action.Data.Value,
			})
		},
	)
}

// Raw is a method for actions.Raw
func (action *DataShowAction) Raw() {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			var raw []byte
			raw, action.Err = action.Data.Raw()
			if action.Err != nil {
				return
			}

			action.W.Write(raw)
		},
	)
}

// SSE is a method for actions.SSE
func (action *DataShowAction) SSE(stream sse.Stream) {
	action.Do(
		action.loadParams,
		action.loadRecord,
		func() {
			stream.Send(sse.Event{Data: action.Data.Value})
		},
	)
}

// GetTopic is a method for actions.SSE
//
// There is no value in this action for specific account_id, so registration topic is a general
// change in the ledger.
func (action *DataShowAction) GetTopic() string {
	return action.GetString("account_id")
}

func (action *DataShowAction) loadParams() {
	action.Address = action.GetAddress("account_id", actions.RequiredParam)
	action.Key = action.GetString("key")
}

func (action *DataShowAction) loadRecord() {
	action.Err = action.CoreQ().
		AccountDataByKey(&action.Data, action.Address, action.Key)
	if action.Err != nil {
		return
	}
}
