package resourceadapter

import (
	"context"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
)

func PopulateAggregateBalance(ctx context.Context, dest *core.AggregateBalanceString, accountId string, balance string) {
	dest.AccountId = accountId
	dest.Balance = balance
	return
}
