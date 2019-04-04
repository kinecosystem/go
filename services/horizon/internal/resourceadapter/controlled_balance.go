package resourceadapter

import (
	"context"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
)

func PopulateControlledBalance(ctx context.Context, dest *core.ControlledBalance, row core.ControlledBalance) {
	dest.AccountId = row.AccountId
	dest.Balance = row.Balance
	return
}
