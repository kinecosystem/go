package horizon

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2"
	"github.com/stellar/go/services/horizon/internal/db2/core"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/test"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/hal"
)

var defaultPage db2.PageQuery = db2.PageQuery{
	Order:  db2.OrderAscending,
	Limit:  db2.DefaultPageSize,
	Cursor: "",
}

func TestGetAccountInfo(t *testing.T) {
	tt := test.Start(t).Scenario("allow_trust")
	defer tt.Finish()

	w := mustInitWeb(context.Background(), &history.Q{tt.HorizonSession()}, &core.Q{tt.CoreSession()}, time.Duration(5), 0, true, true, true)

	res, err := w.getAccountInfo(tt.Ctx, &showActionQueryParams{AccountID: "GCXKG6RN4ONIEPCMNFB732A436Z5PNDSRLGWK7GBLCMQLIFO4S7EYWVU"})
	tt.Assert.NoError(err)

	account, ok := res.(*horizon.Account)
	if !ok {
		tt.Assert.FailNow("type assertion failed when getting account info")
	}

	tt.Assert.Equal("8589934593", account.Sequence)
	tt.Assert.NotEqual(0, account.LastModifiedLedger)

	for _, balance := range account.Balances {
		if balance.Type == "native" {
			tt.Assert.Equal(uint32(0), balance.LastModifiedLedger)
		} else {
			tt.Assert.NotEqual(uint32(0), balance.LastModifiedLedger)
		}
	}

	_, err = w.getAccountInfo(tt.Ctx, &showActionQueryParams{AccountID: "GDBAPLDCAEJV6LSEDFEAUDAVFYSNFRUYZ4X75YYJJMMX5KFVUOHX46SQ"})
	tt.Assert.Equal(errors.Cause(err), sql.ErrNoRows)
}

func TestGetTransactionPage(t *testing.T) {
	tt := test.Start(t).Scenario("base")
	defer tt.Finish()

	ctx := context.Background()
	w := mustInitWeb(ctx, &history.Q{tt.HorizonSession()}, &core.Q{tt.CoreSession()}, time.Duration(5), 0, true, true, true)

	// filter by account
	params := &indexActionQueryParams{
		AccountID:        "GBRPYHIL2CI3FNQ4BXLFMNDLFJUNPU2HY3ZMFSHONUCEOASW7QC7OX2H",
		PagingParams:     defaultPage,
		IncludeFailedTxs: true,
	}

	page, err := w.getTransactionPage(ctx, params)
	pageVal, ok := page.(hal.Page)
	if !ok {
		tt.Assert.FailNow("returned type mismatch")
	}
	tt.Assert.NoError(err)
	tt.Assert.Equal(3, len(pageVal.Embedded.Records))

	// filter by ledger
	params = &indexActionQueryParams{
		LedgerID:         3,
		PagingParams:     defaultPage,
		IncludeFailedTxs: true,
	}

	page, err = w.getTransactionPage(ctx, params)
	pageVal, ok = page.(hal.Page)
	if !ok {
		tt.Assert.FailNow("returned type mismatch")
	}
	tt.Assert.NoError(err)
	tt.Assert.Equal(1, len(pageVal.Embedded.Records))

	// no filter
	params = &indexActionQueryParams{
		PagingParams:     defaultPage,
		IncludeFailedTxs: true,
	}

	page, err = w.getTransactionPage(ctx, params)
	pageVal, ok = page.(hal.Page)
	if !ok {
		tt.Assert.FailNow("returned type mismatch")
	}
	tt.Assert.NoError(err)
	tt.Assert.Equal(4, len(pageVal.Embedded.Records))
}
