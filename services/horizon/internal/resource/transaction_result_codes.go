package resource

import (
	"github.com/kinecosystem/go/services/horizon/internal/txsub"
	"golang.org/x/net/context"
)

// Populate fills out the details
func (res *TransactionResultCodes) Populate(ctx context.Context,
	fail *txsub.FailedTransactionError,
) (err error) {

	res.TransactionCode, err = fail.TransactionResultCode()
	if err != nil {
		return
	}

	res.OperationCodes, err = fail.OperationResultCodes()
	if err != nil {
		return
	}

	return
}