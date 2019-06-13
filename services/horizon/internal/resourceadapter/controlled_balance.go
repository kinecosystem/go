package resourceadapter

import (
	"context"
	"github.com/kinecosystem/go/services/horizon/internal/db2/core"
)

func PopulateControlledAccount(ctx context.Context, dest *core.ControlledAccountString, row core.ControlledAccountString) {
	dest.AccountId = row.AccountId
	dest.Balance = (row.Balance)
	return
}
