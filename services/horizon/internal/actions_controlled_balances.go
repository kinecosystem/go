package horizon

import (
	"github.com/kinecosystem/go/services/horizon/internal/actions"
	"github.com/kinecosystem/go/services/horizon/internal/db2"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/services/horizon/internal/resourceadapter"
	"github.com/kinecosystem/go/support/render/hal"
)

// Interface verifications
var _ actions.JSONer = (*ControlledBalancesAction)(nil)

// ControlledBalancesAction returns a map of accounts and their balance for the given master accountid
type ControlledBalancesAction struct {
	Action
	Address            string
	ControlledBalances []*core.ControlledBalance
	PagingParams       db2.PageQuery
	Page               hal.Page
}

// JSON is a method for actions.JSON
func (action *ControlledBalancesAction) JSON() error {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecord,
		action.loadPage,
		func() { hal.Render(action.W, action.Page) })
	return action.Err
}

func (action *ControlledBalancesAction) loadParams() {
	action.Address = action.GetAddress("account_id", actions.RequiredParam)
}

func (action *ControlledBalancesAction) loadRecord() {
	action.ControlledBalances = make([]*core.ControlledBalance, 0)
	action.Err = action.CoreQ().ControlledBalancesByAccountId(&action.ControlledBalances, action.Address)
	if action.Err != nil {
		return
	}
}

func (action *ControlledBalancesAction) loadPage() {
	for _, cb := range action.ControlledBalances {
		var res core.ControlledBalance
		resourceadapter.PopulateControlledBalance(action.R.Context(), &res, *cb)
		action.Page.Add(res)
	}

	action.Page.FullURL = action.FullURL()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
