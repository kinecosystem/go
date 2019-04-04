package resourceadapter

import (
	"context"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
	"github.com/kinecosystem/go/xdr"
)

func PopulateAggregateBalance(ctx context.Context, dest *core.AggregateBalance, accountId string, balance xdr.Int64) {
	dest.AccountId = accountId
	dest.Balance = balance
	return
}
