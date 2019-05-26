package resourceadapter

import (
	"context"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/amount"
)

func PopulateControlledBalance(ctx context.Context, dest *core.ControlledBalanceString, row core.ControlledBalance) {
	dest.AccountId = row.AccountId
	dest.Balance = amount.String(row.Balance)
	return
}
