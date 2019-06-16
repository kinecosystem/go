package horizon

import (
	"github.com/kinecosystem/go/services/horizon/internal/actions"
	"github.com/kinecosystem/go/services/horizon/internal/db2"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/services/horizon/internal/resourceadapter"
	"github.com/kinecosystem/go/support/render/hal"
)

// Interface verifications
var _ actions.JSONer = (*AggregateBalanceAction)(nil)

// AggregateBalanceAction renders a account summary found by its address.
type AggregateBalanceAction struct {
	Action
	ControlledAccount core.ControlledAccountString
	Page              hal.Page
	PagingParams      db2.PageQuery
}

// JSON is a method for actions.JSON
func (action *AggregateBalanceAction) JSON() error {
	action.Do(
		action.loadParams,
		action.loadRecord,
		action.loadPage,
		func() { hal.Render(action.W, action.Page) })
	return action.Err
}

func (action *AggregateBalanceAction) loadParams() {
	action.ControlledAccount.AccountId = action.GetAddress("account_id", actions.RequiredParam)
}

func (action *AggregateBalanceAction) loadRecord() {
	action.Err = action.CoreQ().AggregateBalanceByAccountId(&action.ControlledAccount)
	if action.Err != nil {
		return
	}
}

func (action *AggregateBalanceAction) loadPage() {
	var res core.AggregateBalanceString
	resourceadapter.PopulateAggregateBalance(action.R.Context(), &res, action.ControlledAccount.AccountId, action.ControlledAccount.Balance)
	action.Page.Add(res)

	action.Page.FullURL = action.FullURL()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
