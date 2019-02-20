package resourceadapter

import (
	"context"

<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/httpx"
	"github.com/kinecosystem/go/services/horizon/internal/txsub"
	. "github.com/kinecosystem/go/protocols/horizon"
	"github.com/kinecosystem/go/support/render/hal"
=======
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/httpx"
	"github.com/stellar/go/services/horizon/internal/txsub"
	"github.com/stellar/go/support/render/hal"
>>>>>>> stellar/master
)

// Populate fills out the details
func PopulateTransactionSuccess(ctx context.Context, dest *TransactionSuccess, result txsub.Result) {
	dest.Hash = result.Hash
	dest.Ledger = result.LedgerSequence
	dest.Env = result.EnvelopeXDR
	dest.Result = result.ResultXDR
	dest.Meta = result.ResultMetaXDR

	lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
	dest.Links.Transaction = lb.Link("/transactions", result.Hash)
	return
}
