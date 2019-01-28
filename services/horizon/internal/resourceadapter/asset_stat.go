package resourceadapter

import (
<<<<<<< HEAD:services/horizon/internal/resource/asset_stat.go
	"github.com/kinecosystem/go/amount"
	"github.com/kinecosystem/go/services/horizon/internal/db2/assets"
	"github.com/kinecosystem/go/services/horizon/internal/render/hal"
	"github.com/kinecosystem/go/xdr"
	"golang.org/x/net/context"
=======
	"context"

	"github.com/stellar/go/amount"
	. "github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/services/horizon/internal/db2/assets"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/support/render/hal"
	"github.com/stellar/go/xdr"
>>>>>>> horizon-v0.15.3:services/horizon/internal/resourceadapter/asset_stat.go
)

// PopulateAssetStat fills out the details
//func PopulateAssetStat(
func PopulateAssetStat(
	ctx context.Context,
	res *AssetStat,
	row assets.AssetStatsR,
) (err error) {
	res.Asset.Type = row.Type
	res.Asset.Code = row.Code
	res.Asset.Issuer = row.Issuer
	res.Amount, err = amount.IntStringToAmount(row.Amount)
	if err != nil {
		return errors.Wrap(err, "Invalid amount in PopulateAssetStat")
	}
	res.NumAccounts = row.NumAccounts
	res.Flags = AccountFlags{
		(row.Flags & int8(xdr.AccountFlagsAuthRequiredFlag)) != 0,
		(row.Flags & int8(xdr.AccountFlagsAuthRevocableFlag)) != 0,
		(row.Flags & int8(xdr.AccountFlagsAuthImmutableFlag)) != 0,
	}
	res.PT = row.SortKey

	res.Links.Toml = hal.NewLink(row.Toml)
	return
}
