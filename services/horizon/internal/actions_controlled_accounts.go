package horizon

import (
	"github.com/kinecosystem/go/services/horizon/internal/actions"
	"github.com/kinecosystem/go/services/horizon/internal/db2"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/services/horizon/internal/resourceadapter"
	"github.com/kinecosystem/go/support/render/hal"
)

// Interface verifications
var _ actions.JSONer = (*ControlledAccountsAction)(nil)

// ControlledAccountsAction returns a map of accounts and their balance for the given master accountid
type ControlledAccountsAction struct {
	Action
	Address            string
	ControlledAccounts []*core.ControlledAccountString
	PagingParams       db2.PageQuery
	Page               hal.Page
}

// JSON is a method for actions.JSON
func (action *ControlledAccountsAction) JSON() error {
	action.Do(
		action.EnsureHistoryFreshness,
		action.loadParams,
		action.ValidateCursorWithinHistory,
		action.loadRecord,
		action.loadPage,
		func() { hal.Render(action.W, action.Page) })
	return action.Err
}

func (action *ControlledAccountsAction) loadParams() {
	action.Address = action.GetAddress("account_id", actions.RequiredParam)
}

func (action *ControlledAccountsAction) loadRecord() {
	action.ControlledAccounts = make([]*core.ControlledAccountString, 0)
	action.Err = action.CoreQ().ControlledAccountsByAccountId(&action.ControlledAccounts, action.Address)
	if action.Err != nil {
		return
	}
}

func (action *ControlledAccountsAction) loadPage() {
	for _, cb := range action.ControlledAccounts {
		var res core.ControlledAccountString
		resourceadapter.PopulateControlledAccount(action.R.Context(), &res, *cb)
		action.Page.Add(res)
	}

	action.Page.FullURL = action.FullURL()
	action.Page.Limit = action.PagingParams.Limit
	action.Page.Cursor = action.PagingParams.Cursor
	action.Page.Order = action.PagingParams.Order
	action.Page.PopulateLinks()
}
