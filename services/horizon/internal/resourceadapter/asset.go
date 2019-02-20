package resourceadapter

import (
	"context"

<<<<<<< HEAD
	"github.com/kinecosystem/go/xdr"
	. "github.com/kinecosystem/go/protocols/horizon"

=======
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
>>>>>>> stellar/master
)

func PopulateAsset(ctx context.Context, dest *Asset, asset xdr.Asset) error {
	return asset.Extract(&dest.Type, &dest.Code, &dest.Issuer)
}
