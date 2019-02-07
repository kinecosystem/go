package resourceadapter

import (
	"context"

	"github.com/kinecosystem/go/xdr"
	. "github.com/kinecosystem/go/protocols/horizon"

)

func PopulateAsset(ctx context.Context, dest *Asset, asset xdr.Asset) error {
	return asset.Extract(&dest.Type, &dest.Code, &dest.Issuer)
}
