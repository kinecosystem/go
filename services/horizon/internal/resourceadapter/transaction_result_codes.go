package resourceadapter

import (
	"context"

<<<<<<< HEAD
	"github.com/kinecosystem/go/services/horizon/internal/txsub"
	. "github.com/kinecosystem/go/protocols/horizon"

=======
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/txsub"
>>>>>>> stellar/master
)

// Populate fills out the details
func PopulateTransactionResultCodes(ctx context.Context,
	dest *TransactionResultCodes,
	fail *txsub.FailedTransactionError,
) (err error) {

	dest.TransactionCode, err = fail.TransactionResultCode()
	if err != nil {
		return
	}

	dest.OperationCodes, err = fail.OperationResultCodes()
	if err != nil {
		return
	}

	return
}
