package bridge

import (
	"github.com/kinecosystem/go/amount"
	b "github.com/kinecosystem/go/build"
	shared "github.com/kinecosystem/go/services/internal/bridge-compliance-shared"
	"github.com/kinecosystem/go/services/internal/bridge-compliance-shared/http/helpers"
	"github.com/kinecosystem/go/services/internal/bridge-compliance-shared/protocols"
)

// PaymentOperationBody represents payment operation
type PaymentOperationBody struct {
	Source      *string
	Destination string
	Amount      string
	Asset       protocols.Asset
}

// ToTransactionMutator returns go-stellar-base TransactionMutator
func (op PaymentOperationBody) ToTransactionMutator() b.TransactionMutator {
	mutators := []interface{}{
		b.Destination{op.Destination},
	}

	if op.Asset.Code != "" && op.Asset.Issuer != "" {
		mutators = append(
			mutators,
			b.CreditAmount{op.Asset.Code, op.Asset.Issuer, op.Amount},
		)
	} else {
		mutators = append(
			mutators,
			b.NativeAmount{op.Amount},
		)
	}

	if op.Source != nil {
		mutators = append(mutators, b.SourceAccount{*op.Source})
	}

	return b.Payment(mutators...)
}

// Validate validates if operation body is valid.
func (op PaymentOperationBody) Validate() error {
	if !shared.IsValidAccountID(op.Destination) {
		return helpers.NewInvalidParameterError("destination", "Destination must be a public key (starting with `G`).")
	}

	_, err := amount.Parse(op.Amount)
	if err != nil {
		return helpers.NewInvalidParameterError("amount", "Invalid amount.")
	}

	err = op.Asset.Validate()
	if err != nil {
		return helpers.NewInvalidParameterError("asset", "Invalid asset.")
	}

	if op.Source != nil && !shared.IsValidAccountID(*op.Source) {
		return helpers.NewInvalidParameterError("source", "Source must be a public key (starting with `G`).")
	}

	return nil
}
