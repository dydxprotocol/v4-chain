package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	canonicaltypes "github.com/dydxprotocol/v4-chain/protocol/x/canonicalusdc/types"
	ratelimittypes "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func uusdcLimiters(hourly, daily int64) ratelimittypes.LimitParams {
	return ratelimittypes.LimitParams{
		Denom: types.UusdcDenom,
		Limiters: []ratelimittypes.Limiter{
			{Period: time.Hour, BaselineMinimum: dtypes.NewInt(hourly)},
			{Period: 24 * time.Hour, BaselineMinimum: dtypes.NewInt(daily)},
		},
	}
}

func uusdcCapacity(hourly, daily int64) ratelimittypes.DenomCapacity {
	return ratelimittypes.DenomCapacity{
		Denom:        types.UusdcDenom,
		CapacityList: []dtypes.SerializableInt{dtypes.NewInt(hourly), dtypes.NewInt(daily)},
	}
}

func withdraw(t *testing.T, ctx sdk.Context, k *Keeper, sender sdk.AccAddress, channel string, amount int64) uint64 {
	t.Helper()
	resp, err := NewTransferDecorator(k).Transfer(ctx, &transfertypes.MsgTransfer{
		SourcePort: transfertypes.PortID, SourceChannel: channel,
		Token: sdk.NewInt64Coin(types.UusdcDenom, amount), Sender: sender.String(),
		Receiver: "remote1receiver", TimeoutTimestamp: 1,
	})
	require.NoError(t, err)
	return resp.Sequence
}

func refillable(k *Keeper, ctx sdk.Context) string {
	return k.GetLedger(ctx).RefillableCapacity.String()
}

// setupRefill enters GRADUAL with a holder owning 300 logical USDC, 100 of it
// Injective-backed.
func setupRefill(t *testing.T, params ratelimittypes.LimitParams) (sdk.Context, *Keeper, *bankMock, sdk.AccAddress) {
	t.Helper()
	ctx, k, bank, transfer, _ := setupKeeper(t)
	bindOpenRoutes(k, activeControls())
	rateLimit(k).params = params
	holder := sdk.AccAddress("holder")
	bank.set(holder, types.UusdcDenom, math.NewInt(200))
	require.NoError(t, updateControls(t, ctx, k, activeControls()))
	transfer.logicalOwner = holder
	require.NoError(t, k.MintLogicalAndSend(ctx, holder, sdk.NewCoins(sdk.NewInt64Coin(types.UusdcDenom, 100))))
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(100))
	ledger := k.GetLedger(ctx)
	ledger.InjectiveBacking = dtypes.NewInt(100)
	require.NoError(t, k.SetLedger(ctx, ledger))
	return ctx, k, bank, holder
}

func TestInjectiveWithdrawalsAccrueRefillableCapacityUpToBaseline(t *testing.T) {
	ctx, k, bank, holder := setupRefill(t, uusdcLimiters(30, 70))

	withdraw(t, ctx, k, holder, nobleChannel, 60)
	require.Equal(t, "0", refillable(k, ctx), "Noble withdrawals accrue nothing")

	first := withdraw(t, ctx, k, holder, injChannel, 40)
	require.Equal(t, "40", refillable(k, ctx))
	withdraw(t, ctx, k, holder, injChannel, 40)
	require.Equal(t, "70", refillable(k, ctx), "capped at the largest limiter baseline")

	// ICS20 refunds the physical voucher to the module before the module settles.
	bank.set(canonicaltypes.ModuleAddress, physicalDenom, math.NewInt(100))
	require.NoError(t, k.Settle(
		ctx, mustPending(t, k, ctx, injChannel, first), false, canonicaltypes.AttributeValueError,
	))
	require.Equal(t, "30", refillable(k, ctx), "an undone charge comes off the refillable amount")
	require.Equal(t, int64(40), rateLimit(k).chargedFor(types.UusdcDenom), "x/ratelimit refunded the undone charge")
}

func TestWithdrawalsWithoutLimitersAccrueNothing(t *testing.T) {
	ctx, k, _, holder := setupRefill(t, ratelimittypes.LimitParams{Denom: types.UusdcDenom})
	withdraw(t, ctx, k, holder, injChannel, 60)
	require.Equal(t, "0", refillable(k, ctx))
}

func TestRefillCapacityIsBoundedByRefillable(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	rateLimit(k).params = uusdcLimiters(30, 100)
	rateLimit(k).capacity = uusdcCapacity(0, 50)
	ledger := canonicaltypes.DefaultLedger()
	ledger.RefillableCapacity = dtypes.NewInt(10)
	require.NoError(t, k.SetLedger(ctx, ledger))

	require.NoError(t, k.RefillCapacity(ctx, types.UusdcDenom, math.NewInt(25)))
	require.Equal(t, int64(10), rateLimit(k).creditedFor(types.UusdcDenom))
	require.Equal(t, "0", refillable(k, ctx))

	require.NoError(t, k.RefillCapacity(ctx, types.UusdcDenom, math.NewInt(25)))
	require.Equal(t, int64(10), rateLimit(k).creditedFor(types.UusdcDenom))
}

func TestRefillCapacityIsClampedToCurrentCap(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	rateLimit(k).params = uusdcLimiters(30, 60)
	rateLimit(k).capacity = uusdcCapacity(0, 0)
	ledger := canonicaltypes.DefaultLedger()
	ledger.RefillableCapacity = dtypes.NewInt(200)
	require.NoError(t, k.SetLedger(ctx, ledger))

	require.NoError(t, k.RefillCapacity(ctx, types.UusdcDenom, math.NewInt(200)))
	require.Equal(t, int64(30), rateLimit(k).creditedFor(types.UusdcDenom), "no limiter is lifted above its baseline")
	require.Equal(t, "0", refillable(k, ctx), "a lowered limiter applies retroactively")
}

func TestRefillCapacityOnlyRestoresWhatIsMissing(t *testing.T) {
	ctx, k, _, _, _ := setupKeeper(t)
	require.NoError(t, k.SetControls(ctx, activeControls()))
	rateLimit(k).params = uusdcLimiters(30, 100)
	rateLimit(k).capacity = uusdcCapacity(20, 90)
	ledger := canonicaltypes.DefaultLedger()
	ledger.RefillableCapacity = dtypes.NewInt(30)
	require.NoError(t, k.SetLedger(ctx, ledger))

	require.NoError(t, k.RefillCapacity(ctx, types.UusdcDenom, math.NewInt(30)))
	require.Equal(t, int64(10), rateLimit(k).creditedFor(types.UusdcDenom), "the bucket recovered the rest itself")
	require.Equal(t, "0", refillable(k, ctx), "the deposit consumes its refillable amount regardless")

	rateLimit(k).capacity = uusdcCapacity(30, 100)
	ledger.RefillableCapacity = dtypes.NewInt(30)
	require.NoError(t, k.SetLedger(ctx, ledger))
	require.NoError(t, k.RefillCapacity(ctx, types.UusdcDenom, math.NewInt(30)))
	require.Equal(t, int64(10), rateLimit(k).creditedFor(types.UusdcDenom), "a full bucket takes nothing")
	require.Equal(t, "0", refillable(k, ctx))
}
