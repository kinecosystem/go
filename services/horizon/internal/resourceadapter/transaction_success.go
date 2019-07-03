package resourceadapter

import (
	"context"

	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/httpx"
	"github.com/stellar/go/services/horizon/internal/txsub"
	"github.com/stellar/go/support/render/hal"
)

// Populate fills out the details
func PopulateTransactionSuccess(ctx context.Context, dest *TransactionSuccess, result txsub.Result, shouldPopulateHalCustomLinks bool) {
	dest.Hash = result.Hash
	dest.Ledger = result.LedgerSequence
	dest.Env = result.EnvelopeXDR
	dest.Result = result.ResultXDR
	dest.Meta = result.ResultMetaXDR

	if shouldPopulateHalCustomLinks {
		lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
		dest.Links = new(TransactionSuccessLinks)
		dest.Links.Transaction = lb.LinkPtr("/transactions", result.Hash)
	}

	return
}
