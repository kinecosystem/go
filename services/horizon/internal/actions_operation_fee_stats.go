package horizon

import (
	"fmt"

	"github.com/stellar/go/services/horizon/internal/operationfeestats"
	"github.com/stellar/go/support/render/hal"
)

// This file contains the actions:
//
// OperationFeeStatsAction: stats representing current state of network fees

// OperationFeeStatsAction renders a few useful statistics that describe the
// current state of operation fees on the network.
type OperationFeeStatsAction struct {
	Action
	Min         int64
	Mode        int64
	LastBaseFee int64
	LastLedger  int64
}

// JSON is a method for actions.JSON
func (action *OperationFeeStatsAction) JSON() {
	action.Do(
		action.loadRecords,
		func() {
			hal.Render(action.W, map[string]string{
				"min_accepted_fee":     fmt.Sprint(action.Min),
				"mode_accepted_fee":    fmt.Sprint(action.Mode),
				"last_ledger_base_fee": fmt.Sprint(action.LastBaseFee),
				"last_ledger":          fmt.Sprint(action.LastLedger),
			})
		},
	)
}

func (action *OperationFeeStatsAction) loadRecords() {
	cur := operationfeestats.CurrentState()

	action.Min = cur.Min
	action.Mode = cur.Mode
	action.LastBaseFee = cur.LastBaseFee
	action.LastLedger = cur.LastLedger
}
