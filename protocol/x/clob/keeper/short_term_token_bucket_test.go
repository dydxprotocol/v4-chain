package keeper_test

import (
	"bytes"
	"encoding/json"
	"math"
	"strings"
	"testing"

	cosmoslog "cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	stattypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/require"
)

// Alice has $10,000. The bucket holds two orders and refills nothing until a test sets a rate.
func setupShortTermTokenBucket(t *testing.T) (keepertest.ClobKeepersTestContext, sdk.Context) {
	ks, ctx := setupShortTermOrderKeeper(t, constants.Alice_Num0_10_000USD, nil)
	ks.ClobKeeper.Flags.ShortTermTokenBucketEnabled = true
	ks.ClobKeeper.Flags.ShortTermTokenBucketShadow = false
	ks.ClobKeeper.Flags.ShortTermTokenBucketMinCapacity = 2
	ks.ClobKeeper.Flags.ShortTermTokenBucketBurstBlocks = 0
	ks.ClobKeeper.Flags.ShortTermTokenBucketVolumePerOrder = 0
	ks.ClobKeeper.Flags.ShortTermTokenBucketCollateralPerOrder = 0
	return ks, ctx.WithBlockHeight(100)
}

func TestShortTermTokenBucket_RejectsWhenEmpty(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper

	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 0)))
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 1)))
	err := k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 2))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	// The bucket is per address, so another of Alice's subaccounts is empty too.
	err = k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num1, 0))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	// With no volume or collateral term nothing ever refills.
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(1_000))
	err = k.ValidateShortTermOrderTokenBucket(ctx.WithBlockHeight(1_000), restingBuyOrder(constants.Alice_Num0, 2))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	// A rejected order still counts as activity, and the bucket is only dropped 10,000 blocks after the last one.
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(10_900))
	err = k.ValidateShortTermOrderTokenBucket(ctx.WithBlockHeight(10_900), restingBuyOrder(constants.Alice_Num0, 2))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(20_900))
	err = k.ValidateShortTermOrderTokenBucket(ctx.WithBlockHeight(20_900), restingBuyOrder(constants.Alice_Num0, 2))
	require.NoError(t, err)

	// Bob has his own bucket.
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Bob_Num0, 0)))
}

func TestShortTermTokenBucket_FractionalCollateralRefill(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	// $10,000 at $20,000 per order is half an order per block.
	k.Flags.ShortTermTokenBucketCollateralPerOrder = 20_000_000_000

	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 0)))
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 1)))
	err := k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 2))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	err = k.ValidateShortTermOrderTokenBucket(ctx.WithBlockHeight(101), restingBuyOrder(constants.Alice_Num0, 2))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	later := ctx.WithBlockHeight(102)
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(later, restingBuyOrder(constants.Alice_Num0, 2)))
	err = k.ValidateShortTermOrderTokenBucket(later, restingBuyOrder(constants.Alice_Num0, 3))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)

	// A long idle stretch refills to capacity and no further.
	idle := ctx.WithBlockHeight(1_000)
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(idle, restingBuyOrder(constants.Alice_Num0, 4)))
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(idle, restingBuyOrder(constants.Alice_Num0, 5)))
	err = k.ValidateShortTermOrderTokenBucket(idle, restingBuyOrder(constants.Alice_Num0, 6))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)
}

// Alice's $1M of volume and $10,000 of collateral each refill one order per block, and her bucket holds five
// blocks of refill.
func setupShortTermTokenBucketWithBurst(t *testing.T) (keepertest.ClobKeepersTestContext, sdk.Context) {
	ks, ctx := setupShortTermTokenBucket(t)
	ks.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &stattypes.UserStats{
		TakerNotional: 600_000_000_000,
		MakerNotional: 400_000_000_000,
	})
	ks.ClobKeeper.Flags.ShortTermTokenBucketVolumePerOrder = 1_000_000_000_000
	ks.ClobKeeper.Flags.ShortTermTokenBucketCollateralPerOrder = 10_000_000_000
	ks.ClobKeeper.Flags.ShortTermTokenBucketBurstBlocks = 5
	return ks, ctx
}

// requirePlaces places n orders at ctx and requires that exactly those go through.
func requirePlaces(t *testing.T, k *keeper.Keeper, ctx sdk.Context, subaccountId satypes.SubaccountId, n uint32) {
	for i := range n {
		require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(subaccountId, i)))
	}
	err := k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(subaccountId, n))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)
}

func TestShortTermTokenBucket_NewBucketGrowsToBurst(t *testing.T) {
	ks, ctx := setupShortTermTokenBucketWithBurst(t)
	k := ks.ClobKeeper

	// A new bucket holds only the two-order minimum, so moving money to a fresh address buys no burst.
	requirePlaces(t, k, ctx, constants.Alice_Num0, 2)

	// It grows by two orders a block, up to its capacity of ten.
	requirePlaces(t, k, ctx.WithBlockHeight(103), constants.Alice_Num0, 6)
	requirePlaces(t, k, ctx.WithBlockHeight(1_000), constants.Alice_Num0, 10)
}

func TestShortTermTokenBucket_FullBucketKeptUntilIdle(t *testing.T) {
	ks, ctx := setupShortTermTokenBucketWithBurst(t)
	k := ks.ClobKeeper
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 0)))

	// A full bucket bigger than a new one survives pruning, so the burst it saved is still there.
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(1_000))
	requirePlaces(t, k, ctx.WithBlockHeight(1_000), constants.Alice_Num0, 10)

	// After 10,000 blocks without an order it is dropped, and the address starts again from the minimum.
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(11_000))
	requirePlaces(t, k, ctx.WithBlockHeight(11_000), constants.Alice_Num0, 2)
}

func TestShortTermTokenBucket_LowerRateClipsSavedBurst(t *testing.T) {
	ks, ctx := setupShortTermTokenBucketWithBurst(t)
	k := ks.ClobKeeper
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, 0)))

	// Her volume ages out, so the next rating halves her refill and capacity, and the ten orders she saved are cut
	// to five.
	ks.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &stattypes.UserStats{})
	requirePlaces(t, k, ctx.WithBlockHeight(1_000), constants.Alice_Num0, 5)
}

func TestShortTermTokenBucket_CollateralAcrossSubaccounts(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	ks.SubaccountsKeeper.SetSubaccount(ctx, constants.Alice_Num1_10_000USD)
	// Each $10,000 refills one order a block, and both of Alice's subaccounts count whichever one places.
	k.Flags.ShortTermTokenBucketCollateralPerOrder = 10_000_000_000

	requirePlaces(t, k, ctx, constants.Alice_Num1, 2)
	next := ctx.WithBlockHeight(101)
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, restingBuyOrder(constants.Alice_Num1, 2)))
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, restingBuyOrder(constants.Alice_Num0, 3)))
	err := k.ValidateShortTermOrderTokenBucket(next, restingBuyOrder(constants.Alice_Num0, 4))
	require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)
}

func TestShortTermTokenBucket_RateRefreshedEvery100Blocks(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	k.Flags.ShortTermTokenBucketVolumePerOrder = 1_000_000_000_000
	requirePlaces(t, k, ctx, constants.Alice_Num0, 2)

	// Volume recorded after the bucket is rated only counts from its next rating, 100 blocks on, and refills
	// from there.
	ks.StatsKeeper.SetUserStats(ctx, constants.Alice_Num0.Owner, &stattypes.UserStats{TakerNotional: 1_000_000_000_000})
	for _, height := range []int64{101, 199, 200} {
		err := k.ValidateShortTermOrderTokenBucket(ctx.WithBlockHeight(height), restingBuyOrder(constants.Alice_Num0, 2))
		require.ErrorIs(t, err, types.ErrShortTermOrderRateExceeded)
	}
	requirePlaces(t, k, ctx.WithBlockHeight(201), constants.Alice_Num0, 1)
}

func TestShortTermTokenBucket_ReduceOnlyShadowAndDisabled(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	k.Flags.ShortTermTokenBucketMinCapacity = 1
	order := restingBuyOrder(constants.Alice_Num0, 0)
	reduceOnly := order
	reduceOnly.ReduceOnly = true

	// Reduce-only orders spend tokens like any other.
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, reduceOnly))
	require.ErrorIs(t, k.ValidateShortTermOrderTokenBucket(ctx, order), types.ErrShortTermOrderRateExceeded)

	// Once the bucket is empty, one reduce-only order a block still goes through.
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, reduceOnly))
	require.ErrorIs(t, k.ValidateShortTermOrderTokenBucket(ctx, reduceOnly), types.ErrShortTermOrderRateExceeded)
	next := ctx.WithBlockHeight(101)
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, reduceOnly))
	require.ErrorIs(t, k.ValidateShortTermOrderTokenBucket(next, order), types.ErrShortTermOrderRateExceeded)

	// Shadow mode lets an empty bucket through, and enforcing again finds it still empty.
	k.Flags.ShortTermTokenBucketShadow = true
	for range 3 {
		require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, order))
	}
	k.Flags.ShortTermTokenBucketShadow = false
	require.ErrorIs(t, k.ValidateShortTermOrderTokenBucket(next, order), types.ErrShortTermOrderRateExceeded)

	k.Flags.ShortTermTokenBucketEnabled = false
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, order))
}

func TestShortTermTokenBucket_WindowReport(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	k.Flags.ShortTermTokenBucketShadow = true
	var buf bytes.Buffer
	ctx = ctx.WithLogger(cosmoslog.NewLogger(&buf, cosmoslog.OutputJSONOption()))

	// Five orders at height 100 against a bucket of two, then one more a block later.
	for i := range uint32(5) {
		require.NoError(t, k.ValidateShortTermOrderTokenBucket(ctx, restingBuyOrder(constants.Alice_Num0, i)))
	}
	next := ctx.WithBlockHeight(101)
	require.NoError(t, k.ValidateShortTermOrderTokenBucket(next, restingBuyOrder(constants.Alice_Num0, 5)))

	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(150))
	require.Empty(t, tokenBucketWindowReports(t, &buf))

	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(200))
	reports := tokenBucketWindowReports(t, &buf)
	require.Len(t, reports, 1)
	require.Equal(t, constants.Alice_Num0.Owner, reports[0]["owner"])
	require.EqualValues(t, 100, reports[0]["window_blocks"])
	require.EqualValues(t, 6, reports[0]["placed"])
	require.EqualValues(t, 4, reports[0]["rejected"])
	require.EqualValues(t, 5, reports[0]["max_per_block"])

	// The window resets, so a quiet window reports nothing.
	buf.Reset()
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(300))
	require.Empty(t, tokenBucketWindowReports(t, &buf))
}

func TestShortTermTokenBucket_WindowReportOnEviction(t *testing.T) {
	ks, ctx := setupShortTermTokenBucket(t)
	k := ks.ClobKeeper
	k.Flags.ShortTermTokenBucketShadow = true
	k.Flags.ShortTermTokenBucketCollateralPerOrder = 10_000_000_000
	var buf bytes.Buffer
	ctx = ctx.WithLogger(cosmoslog.NewLogger(&buf, cosmoslog.OutputJSONOption()))

	opened := ctx.WithBlockHeight(150)
	for i := range uint32(3) {
		require.NoError(t, k.ValidateShortTermOrderTokenBucket(opened, restingBuyOrder(constants.Alice_Num0, i)))
	}

	// One order a block refills the bucket of two by height 152, so it is dropped when the window closes. Its
	// window counts from the order that opened it, and the rejection is reported first.
	k.PruneShortTermTokenBuckets(ctx.WithBlockHeight(200))
	reports := tokenBucketWindowReports(t, &buf)
	require.Len(t, reports, 1)
	require.EqualValues(t, 50, reports[0]["window_blocks"])
	require.EqualValues(t, 3, reports[0]["placed"])
	require.EqualValues(t, 1, reports[0]["rejected"])
}

func tokenBucketWindowReports(t *testing.T, buf *bytes.Buffer) []map[string]any {
	var reports []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if line == "" {
			continue
		}
		var entry map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		if entry["message"] == "Short-term order token bucket window" {
			reports = append(reports, entry)
		}
	}
	return reports
}

// Cold rates the address on every order, as for a new address or one due for a new rate. Warm is every other
// order.
func BenchmarkValidateShortTermOrderTokenBucket(b *testing.B) {
	for _, sa := range benchSubaccounts {
		for _, cold := range []bool{true, false} {
			name := sa.name + "/warm"
			if cold {
				name = sa.name + "/cold"
			}
			b.Run(name, func(b *testing.B) {
				ks, ctx := setupShortTermOrderKeeper(b, sa.subaccount, nil)
				ks.ClobKeeper.Flags.ShortTermTokenBucketEnabled = true
				ks.ClobKeeper.Flags.ShortTermTokenBucketMinCapacity = math.MaxUint32
				order := restingBuyOrder(*sa.subaccount.Id, 0)
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					orderCtx := ctx
					if cold {
						orderCtx = ctx.WithBlockHeight(ctx.BlockHeight() + 100*int64(i))
					}
					if err := ks.ClobKeeper.ValidateShortTermOrderTokenBucket(orderCtx, order); err != nil {
						b.Fatal(err)
					}
				}
			})
		}
	}
}
