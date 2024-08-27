package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	statstypes "github.com/dydxprotocol/v4-chain/protocol/x/stats/types"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.AffiliatesKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestRegisterAffiliate_GetReferredBy(t *testing.T) {
	tests := []struct {
		name        string
		referee     string
		affiliate   string
		expectError bool
		setup       func(t *testing.T, ctx sdk.Context, k *keeper.Keeper)
	}{
		{
			name:        "Register new affiliate",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.BobAccAddress.String(),
			expectError: false,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				// No setup needed for this test case
			},
		},
		{
			name:        "Register existing referee",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.CarlAccAddress.String(),
			expectError: true,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				err := k.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tests {
		tApp := testapp.NewTestAppBuilder(t).Build()
		ctx := tApp.InitChain()
		k := tApp.App.AffiliatesKeeper
		tc.setup(t, ctx, &k)
		t.Run(tc.name, func(t *testing.T) {
			err := k.RegisterAffiliate(ctx, tc.referee, tc.affiliate)
			if tc.expectError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			affiliate, exists := k.GetReferredBy(ctx, tc.referee)
			require.True(t, exists)
			require.Equal(t, tc.affiliate, affiliate)
		})
	}
}

func TestGetReferredByEmptyAffiliate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliate, exists := k.GetReferredBy(ctx, constants.AliceAccAddress.String())
	require.False(t, exists)
	require.Equal(t, "", affiliate)
}

func TestAddReferredVolume(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliate := "affiliate1"
	initialVolume := big.NewInt(1000)
	addedVolume := big.NewInt(500)

	err := k.AddReferredVolume(ctx, affiliate, initialVolume)
	require.NoError(t, err)

	volume, err := k.GetReferredVolume(ctx, affiliate)
	require.NoError(t, err)
	require.Equal(t, initialVolume, volume)

	err = k.AddReferredVolume(ctx, affiliate, addedVolume)
	require.NoError(t, err)

	updatedVolume, err := k.GetReferredVolume(ctx, affiliate)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1500), updatedVolume)
}

func TestGetReferredVolumeInvalidAffiliate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliate := "malformed_address"
	_, exists := k.GetReferredBy(ctx, affiliate)
	require.False(t, exists)

	affiliate = constants.AliceAccAddress.String()
	_, exists = k.GetReferredBy(ctx, affiliate)
	require.False(t, exists)
}

func TestGetTakerFeeShareViaReferredVolume(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper
	// Set up affiliate tiers
	affiliateTiers := constants.DefaultAffiliateTiers
	k.UpdateAffiliateTiers(ctx, affiliateTiers)
	stakingKeeper := tApp.App.StakingKeeper

	err := stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(constants.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	)
	require.NoError(t, err)

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, constants.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm, feeSharePpm)

	// Add more referred volume to upgrade tier
	err = k.AddReferredVolume(ctx, affiliate, big.NewInt(
		int64(constants.DefaultAffiliateTiers.Tiers[1].ReqReferredVolume),
	))
	require.NoError(t, err)

	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, constants.DefaultAffiliateTiers.Tiers[1].TakerFeeSharePpm, feeSharePpm)
}

func TestGetTakerFeeShareViaStakedAmount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper
	ctx = ctx.WithBlockTime(time.Now())
	// Set up affiliate tiers
	affiliateTiers := constants.DefaultAffiliateTiers
	k.UpdateAffiliateTiers(ctx, affiliateTiers)

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper
	err := stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(constants.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	)
	require.NoError(t, err)
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, constants.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm, feeSharePpm)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(
		time.Duration(statstypes.StakedAmountCacheDurationSeconds+1) * time.Second,
	))
	// Add more staked amount to upgrade tier
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(new(big.Int).Mul(
				big.NewInt(int64(constants.DefaultAffiliateTiers.Tiers[1].ReqStakedWholeCoins)),
				big.NewInt(1e18),
			))))
	require.NoError(t, err)
	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, constants.DefaultAffiliateTiers.Tiers[1].TakerFeeSharePpm, feeSharePpm)
}

// Test volume qualifies for tier 2 and stake qualifies for tier 3
// should return tier 3
func TestGetTierForAffiliate_VolumeAndStake(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliateTiers := constants.DefaultAffiliateTiers
	k.UpdateAffiliateTiers(ctx, affiliateTiers)
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper

	err := stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(constants.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	)
	require.NoError(t, err)

	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	reqReferredVolume := big.NewInt(int64(affiliateTiers.Tiers[2].ReqReferredVolume))
	err = k.AddReferredVolume(ctx, affiliate, reqReferredVolume)
	require.NoError(t, err)

	stakedAmount := new(big.Int).Mul(
		big.NewInt(int64(affiliateTiers.Tiers[3].ReqStakedWholeCoins)),
		big.NewInt(1e18),
	)
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(affiliate,
			constants.AliceValAddress.String(),
			math.LegacyNewDecFromBigInt(stakedAmount),
		),
	)
	require.NoError(t, err)

	tierLevel, feeSharePpm, err := k.GetTierForAffiliate(ctx, affiliate)
	require.NoError(t, err)

	require.Equal(t, uint32(3), tierLevel)
	require.Equal(t, affiliateTiers.Tiers[3].TakerFeeSharePpm, feeSharePpm)
}

func TestUpdateAffiliateTiers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	// Set up valid affiliate tiers
	validTiers := constants.DefaultAffiliateTiers
	k.UpdateAffiliateTiers(ctx, validTiers)

	// Retrieve and validate updated tiers
	updatedTiers, err := k.GetAllAffiliateTiers(ctx)
	require.NoError(t, err)
	require.Equal(t, validTiers, updatedTiers)
}
