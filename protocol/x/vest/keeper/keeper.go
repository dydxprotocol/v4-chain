package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	"math/big"
	"time"

	sdklog "cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	gometrics "github.com/armon/go-metrics"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/maps"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

type (
	Keeper struct {
		cdc             codec.BinaryCodec
		storeKey        storetypes.StoreKey
		bankKeeper      types.BankKeeper
		blockTimeKeeper types.BlockTimeKeeper
		authorities     map[string]struct{}
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper types.BankKeeper,
	blockTimeKeeper types.BlockTimeKeeper,
	authorities []string,
) *Keeper {
	return &Keeper{
		cdc:             cdc,
		storeKey:        storeKey,
		bankKeeper:      bankKeeper,
		blockTimeKeeper: blockTimeKeeper,
		authorities:     maps.ArrayToMapInterface(authorities),
	}
}

func (k Keeper) HasAuthority(authority string) bool {
	_, ok := k.authorities[authority]
	return ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With(sdklog.ModuleKey, fmt.Sprintf("x/%s", types.ModuleName))
}

// Process vesting for all vest entries. Intended to be called in BeginBlocker.
// For each vest entry:
// 1. Return if `block_time <= vest_entry.start_time` (vesting has not started yet)
// 2. Return if `prev_block_time >= vest_entry.end_time` (vesting has ended)
// 3. Transfer the following amount of tokens from vester account to treasury account:
//
//		  min(
//			(block_time - last_vest_time) / (end_time - last_vest_time),
//		 	1
//		  ) * vester_account_balance
//
//	  where `last_vest_time = max(start_time, prev_block_time)`
func (k Keeper) ProcessVesting(ctx sdk.Context) {
	// Convert timestamps to milliseconds for algebraic operations.
	blockTimeMilli := ctx.BlockTime().UnixMilli()
	prevBlockInfo := k.blockTimeKeeper.GetPreviousBlockInfo(ctx)
	prevBlockTimeMilli := prevBlockInfo.Timestamp.UnixMilli()

	// Process each vest entry.
	for _, entry := range k.GetAllVestEntries(ctx) {
		startTimeMilli := entry.StartTime.UnixMilli()
		endTimeMilli := entry.EndTime.UnixMilli()
		// `block_time` <= `start_time`. Vesting has not started.
		if blockTimeMilli <= startTimeMilli {
			continue
		}
		// `end_time` <= `prev_block_time`. Vesting has ended.
		if endTimeMilli <= prevBlockTimeMilli {
			continue
		}

		// last_vest_time = max(start_time, prev_block_time)
		lastVestTimeMilli := lib.Max(startTimeMilli, prevBlockTimeMilli)
		// Calculate (block_time - last_vest_time) / (end_time - last_vest_time)
		// Given `block_time > prev_block_time` and `block_time > start_time` ===> `block_time > last_vest_time`
		// Given `end_time > prev_block_time` and `end_time > start_time` ===> `end_time > last_vest_time`
		// Therefore, both numerator and denominator are positive.
		bigRatVestProportion := big.NewRat(blockTimeMilli-lastVestTimeMilli, endTimeMilli-lastVestTimeMilli)

		// Get vester account remaining balance.
		vesterBalance := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(entry.VesterAccount), entry.Denom)
		vestAmount := vesterBalance.Amount
		if bigRatVestProportion.Cmp(lib.BigRat1()) < 0 {
			// vestProportion < 1, so vest_amount = vester_balance * vestProportion
			bigRatBalance := new(big.Rat).SetInt(vesterBalance.Amount.BigInt())
			bigRatVestAmount := new(big.Rat).Mul(
				bigRatBalance,
				bigRatVestProportion,
			)
			vestAmount = sdkmath.NewIntFromBigInt(lib.BigRatRound(bigRatVestAmount, false))
		}

		if !vestAmount.IsZero() {
			// Transfer vest_amount from vester_account to treasury_account.
			// Since `vest_amount = min(vest_proportion, 1) * vester_balance`,
			// we must have `vest_amount <= vester_balance`
			if err := k.bankKeeper.SendCoinsFromModuleToModule(
				ctx,
				entry.VesterAccount,
				entry.TreasuryAccount,
				sdk.NewCoins(sdk.NewCoin(entry.Denom, vestAmount)),
			); err != nil {
				panic(err)
			}
		}

		// Report vest amount.
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.VestAmount},
			float32(vestAmount.Int64()),
			[]gometrics.Label{metrics.GetLabelForStringValue(metrics.VesterAccount, entry.VesterAccount)},
		)
		// Report vester account balance after vest event.
		balanceAfterVest := k.bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(entry.VesterAccount), entry.Denom)
		telemetry.SetGaugeWithLabels(
			[]string{types.ModuleName, metrics.BalanceAfterVestEvent},
			float32(balanceAfterVest.Amount.Int64()),
			[]gometrics.Label{metrics.GetLabelForStringValue(metrics.VesterAccount, entry.VesterAccount)},
		)
	}
}

func (k Keeper) GetAllVestEntries(ctx sdk.Context) (
	list []types.VestEntry,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestEntryKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.VestEntry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return list
}

func (k Keeper) GetVestEntry(ctx sdk.Context, vesterAccount string) (
	val types.VestEntry,
	err error,
) {
	defer telemetry.ModuleMeasureSince(
		types.ModuleName,
		time.Now(),
		metrics.GetVestEntry,
		metrics.Latency,
	)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestEntryKeyPrefix))
	b := store.Get(types.VestEntryKey(vesterAccount))

	// If VestEntry does not exist in state, return error
	if b == nil {
		return types.VestEntry{}, errorsmod.Wrapf(types.ErrVestEntryNotFound, "vesterAccount: %s", vesterAccount)
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, nil
}

func (k Keeper) SetVestEntry(
	ctx sdk.Context,
	entry types.VestEntry,
) (
	err error,
) {
	if err := entry.Validate(); err != nil {
		return err
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestEntryKeyPrefix))
	b := k.cdc.MustMarshal(&entry)
	store.Set(types.VestEntryKey(
		entry.VesterAccount,
	), b)
	return nil
}

func (k Keeper) DeleteVestEntry(
	ctx sdk.Context,
	vesterAccount string,
) (
	err error,
) {
	if _, err := k.GetVestEntry(ctx, vesterAccount); err != nil {
		return errorsmod.Wrap(err, "failed to delete vest entry")
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.VestEntryKeyPrefix))
	store.Delete(types.VestEntryKey(vesterAccount))

	return nil
}
