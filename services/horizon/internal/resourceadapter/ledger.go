package resourceadapter

import (
	"context"
	"fmt"

	"github.com/stellar/go/amount"
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2/history"
	"github.com/stellar/go/services/horizon/internal/httpx"
	"github.com/stellar/go/support/render/hal"
	"github.com/stellar/go/xdr"
)

func PopulateLedger(ctx context.Context, dest *Ledger, row history.Ledger, shouldPopulateHalCustomLinks bool) {
	dest.ID = row.LedgerHash
	dest.PT = row.PagingToken()
	dest.Hash = row.LedgerHash
	dest.PrevHash = row.PreviousLedgerHash.String
	dest.Sequence = row.Sequence
	// Default to `transaction_count`
	dest.SuccessfulTransactionCount = row.TransactionCount
	if row.SuccessfulTransactionCount != nil {
		dest.SuccessfulTransactionCount = *row.SuccessfulTransactionCount
	}
	dest.FailedTransactionCount = row.FailedTransactionCount
	dest.OperationCount = row.OperationCount
	dest.ClosedAt = row.ClosedAt
	dest.TotalCoins = amount.String(xdr.Int64(row.TotalCoins))
	dest.FeePool = amount.String(xdr.Int64(row.FeePool))
	dest.BaseFee = row.BaseFee
	dest.BaseReserve = row.BaseReserve
	dest.MaxTxSetSize = row.MaxTxSetSize
	dest.ProtocolVersion = row.ProtocolVersion

	if row.LedgerHeaderXDR.Valid {
		dest.HeaderXDR = row.LedgerHeaderXDR.String
	} else {
		dest.HeaderXDR = ""
	}

	if shouldPopulateHalCustomLinks {
		self := fmt.Sprintf("/ledgers/%d", row.Sequence)
		lb := hal.LinkBuilder{httpx.BaseURL(ctx)}
		dest.Links = new(LedgerLinks)
		dest.Links.Self = lb.LinkPtr(self)
		dest.Links.Transactions = lb.PagedLinkPtr(self, "transactions")
		dest.Links.Operations = lb.PagedLinkPtr(self, "operations")
		dest.Links.Payments = lb.PagedLinkPtr(self, "payments")
		dest.Links.Effects = lb.PagedLinkPtr(self, "effects")
	}

	return
}
