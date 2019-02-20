package resourceadapter

import (
	"context"

<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/db2/history"
	. "github.com/kinecosystem/go/protocols/horizon"
=======
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2/history"
>>>>>>> stellar/master
)

func PopulateHistoryAccount(ctx context.Context, dest *HistoryAccount, row history.Account) {
	dest.ID = row.Address
	dest.AccountID = row.Address
}
