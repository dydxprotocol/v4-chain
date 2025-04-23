package keeper

import (
	"fmt"
	"math"
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	"github.com/dydxprotocol/v4-chain/protocol/lib/slinky"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/dydxprotocol/slinky/x/marketmap/types/tickermetadata"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/listing/types"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// Function to set hard cap on listed markets in module store
func (k Keeper) SetMarketsHardCap(ctx sdk.Context, hardCap uint32) error {
	store := ctx.KVStore(k.storeKey)
	value := gogotypes.UInt32Value{Value: hardCap}
	store.Set([]byte(types.HardCapForMarketsKey), k.cdc.MustMarshal(&value))
	return nil
}

// Function to get hard cap on listed markets from module store
func (k Keeper) GetMarketsHardCap(ctx sdk.Context) (hardCap uint32) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.HardCapForMarketsKey))
	var result gogotypes.UInt32Value
	k.cdc.MustUnmarshal(b, &result)
	return result.Value
}

// Function to wrap the creation of a new market
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreateMarket(
	ctx sdk.Context,
	ticker string,
) (marketId uint32, err error) {
	marketId = k.PricesKeeper.AcquireNextMarketID(ctx)

	// Get market details from marketmap
	// TODO: change to use util from marketmap when available
	marketMapPair, err := slinky.MarketPairToCurrencyPair(ticker)
	if err != nil {
		return 0, err
	}
	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, marketMapPair.String())
	if err != nil {
		return 0, types.ErrMarketNotFound
	}

	// Create a new market
	market, err := k.PricesKeeper.CreateMarket(
		ctx,
		pricestypes.MarketParam{
			Id:                marketId,
			Pair:              ticker,
			MinPriceChangePpm: types.MinPriceChangePpm_LongTail,
		},
		pricestypes.MarketPrice{
			Id:       marketId,
			Exponent: int32(marketMapDetails.Ticker.Decimals) * -1,
			Price:    0,
		},
	)
	if err != nil {
		return 0, err
	}

	return market.Id, nil
}

// Function to wrap the creation of a new clob pair
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreateClobPair(
	ctx sdk.Context,
	perpetualId uint32,
) (clobPairId uint32, err error) {
	clobPairId = k.ClobKeeper.AcquireNextClobPairID(ctx)

	clobPair, err := k.ClobKeeper.CreatePerpetualClobPair(
		ctx,
		clobPairId,
		perpetualId,
		satypes.BaseQuantums(types.DefaultStepBaseQuantums),
		types.DefaultQuantumConversionExponent,
		types.SubticksPerTick_LongTail,
		clobtypes.ClobPair_STATUS_ACTIVE,
	)
	if err != nil {
		return 0, err
	}

	return clobPair.Id, nil
}

// Function to wrap the creation of a new perpetual
// Note: This will only list long-tail/isolated markets
func (k Keeper) CreatePerpetual(
	ctx sdk.Context,
	marketId uint32,
	ticker string,
) (perpetualId uint32, err error) {
	perpetualId = k.PerpetualsKeeper.AcquireNextPerpetualID(ctx)

	// Get reference price from market map
	// TODO: change to use util from marketmap when available
	marketMapPair, err := slinky.MarketPairToCurrencyPair(ticker)
	if err != nil {
		return 0, err
	}
	marketMapDetails, err := k.MarketMapKeeper.GetMarket(ctx, marketMapPair.String())
	if err != nil {
		return 0, types.ErrMarketNotFound
	}
	metadata, err := tickermetadata.DyDxFromJSONString(marketMapDetails.Ticker.Metadata_JSON)
	if err != nil {
		return 0, types.ErrInvalidMarketMapTickerMetadata
	}
	if metadata.ReferencePrice == 0 {
		return 0, types.ErrReferencePriceZero
	}

	// calculate atomic resolution from reference price
	// atomic resolution = -6 - (log10(referencePrice) - decimals)
	atomicResolution := types.ResolutionOffset -
		(int32(math.Floor(math.Log10(float64(metadata.ReferencePrice)))) -
			int32(marketMapDetails.Ticker.Decimals))

	marketType := perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED
	if metadata.CrossLaunch {
		marketType = perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS
	}

	// Create a new perpetual
	perpetual, err := k.PerpetualsKeeper.CreatePerpetual(
		ctx,
		perpetualId,
		ticker,
		marketId,
		atomicResolution,
		types.DefaultFundingPpm,
		types.LiquidityTier_IML_5x,
		marketType,
	)
	if err != nil {
		return 0, err
	}

	return perpetual.GetId(), nil
}

func (k Keeper) UpgradeIsolatedPerpetualToCross(
	ctx sdk.Context,
	perpetualId uint32,
) error {
	// Validate perpetual exists and is in isolated mode
	perpetual, err := k.PerpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return err
	}
	if perpetual.Params.GetMarketType() != perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED {
		return fmt.Errorf("perpetual %d is not an isolated perpetual and cannot be upgraded to cross", perpetualId)
	}

	err = k.SubaccountsKeeper.TransferIsolatedInsuranceFundToCross(ctx, perpetualId)
	if err != nil {
		return err
	}

	err = k.SubaccountsKeeper.TransferIsolatedCollateralToCross(ctx, perpetualId)
	if err != nil {
		return err
	}

	_, err = k.PerpetualsKeeper.SetPerpetualMarketType(
		ctx,
		perpetualId,
		perpetualtypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS,
	)
	if err != nil {
		return err
	}

	// Fetch updated perpetual data after the upgrade
	perpetual, err = k.PerpetualsKeeper.GetPerpetual(ctx, perpetualId)
	if err != nil {
		return err
	}

	// Emit indexer event.
	k.GetIndexerEventManager().AddTxnEvent(
		ctx,
		indexerevents.SubtypeUpdatePerpetual,
		indexerevents.UpdatePerpetualEventVersion,
		indexer_manager.GetBytes(
			indexerevents.NewUpdatePerpetualEvent(
				perpetual.Params.Id,
				perpetual.Params.Ticker,
				perpetual.Params.MarketId,
				perpetual.Params.AtomicResolution,
				perpetual.Params.LiquidityTier,
				perpetual.Params.MarketType,
				perpetual.Params.DefaultFundingPpm,
			),
		),
	)

	return nil
}

// Function to set listing vault deposit params in module store
func (k Keeper) SetListingVaultDepositParams(
	ctx sdk.Context,
	params types.ListingVaultDepositParams,
) error {
	// Validate the params
	if err := params.Validate(); err != nil {
		return err
	}

	// Store the params in the module store
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.ListingVaultDepositParamsKey)
	store.Set(key, k.cdc.MustMarshal(&params))
	return nil
}

// Function to get listing vault deposit params from module store
func (k Keeper) GetListingVaultDepositParams(
	ctx sdk.Context,
) (vaultDepositParams types.ListingVaultDepositParams) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.ListingVaultDepositParamsKey))
	k.cdc.MustUnmarshal(b, &vaultDepositParams)
	return vaultDepositParams
}

// Function to deposit to the megavault for a new PML market
// This function deposits money to the megavault, transfers the new vault
// deposit amount to the new market vault and locks the shares for the deposit
func (k Keeper) DepositToMegavaultforPML(
	ctx sdk.Context,
	fromSubaccount satypes.SubaccountId,
	clobPairId uint32,
) error {
	// Get the listing vault deposit params
	vaultDepositParams := k.GetListingVaultDepositParams(ctx)

	// Deposit to the megavault
	totalDepositAmount := new(big.Int).Add(
		vaultDepositParams.NewVaultDepositAmount.BigInt(),
		vaultDepositParams.MainVaultDepositAmount.BigInt(),
	)
	mintedShares, err := k.VaultKeeper.DepositToMegavault(
		ctx,
		fromSubaccount,
		totalDepositAmount,
	)
	if err != nil {
		return err
	}

	vaultId := vaulttypes.VaultId{
		Type:   vaulttypes.VaultType_VAULT_TYPE_CLOB,
		Number: clobPairId,
	}

	// Transfer the new vault deposit amount to the new market vault
	err = k.VaultKeeper.AllocateToVault(
		ctx,
		vaultId,
		vaultDepositParams.NewVaultDepositAmount.BigInt(),
	)
	if err != nil {
		return err
	}

	// Lock the shares for the new vault deposit amount
	err = k.VaultKeeper.LockShares(
		ctx,
		fromSubaccount.Owner,
		vaulttypes.BigIntToNumShares(mintedShares),
		uint32(ctx.BlockHeight())+vaultDepositParams.NumBlocksToLockShares,
	)
	if err != nil {
		return err
	}

	// Activate vault to quoting status
	err = k.VaultKeeper.SetVaultStatus(
		ctx,
		vaultId,
		vaulttypes.VaultStatus_VAULT_STATUS_QUOTING,
	)
	if err != nil {
		return err
	}

	return nil
}
