package keeper_test

import (
	"math/big"
	"testing"

	"cosmossdk.io/math"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()

	logger := tApp.App.AffiliatesKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func TestRegisterAffiliate(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	tests := []struct {
		name        string
		referee     string
		affiliate   string
		expectError bool
	}{
		{
			name:        "Register new affiliate",
			referee:     "referee1",
			affiliate:   "affiliate1",
			expectError: false,
		},
		{
			name:        "Register existing referee",
			referee:     "referee1",
			affiliate:   "affiliate2",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.RegisterAffiliate(ctx, tc.referee, tc.affiliate)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				affiliate, exists := k.GetReferredBy(ctx, tc.referee)
				require.True(t, exists)
				require.Equal(t, tc.affiliate, affiliate)
			}
		})
	}
}

func TestAddReferredVolume(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliate := "affiliate1"
	initialVolume := dtypes.NewInt(1000)
	addedVolume := dtypes.NewInt(500)

	err := k.AddReferredVolume(ctx, affiliate, initialVolume)
	require.NoError(t, err)

	volume, exists := k.GetReferredVolume(ctx, affiliate)
	require.True(t, exists)
	require.Equal(t, initialVolume, volume)

	err = k.AddReferredVolume(ctx, affiliate, addedVolume)
	require.NoError(t, err)

	updatedVolume, exists := k.GetReferredVolume(ctx, affiliate)
	require.True(t, exists)
	require.Equal(t, dtypes.NewInt(1500), updatedVolume)
}

func TestGetTakerFeeShareViaReferredVolume(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	// Set up affiliate tiers
	affiliateTiers := types.AffiliateTiers{
		Tiers: []types.AffiliateTiers_Tier{
			{Level: 0, ReqReferredVolume: 0, TakerFeeSharePpm: 1000, ReqStakedWholeCoins: 200},
			{Level: 1, ReqReferredVolume: 1000, TakerFeeSharePpm: 2000, ReqStakedWholeCoins: 1000},
		},
	}
	err := k.UpdateAffiliateTiers(ctx, affiliateTiers)
	require.NoError(t, err)

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(big.NewInt(100))))
	require.NoError(t, err)
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Add referred volume for affiliate
	err = k.AddReferredVolume(ctx, affiliate, dtypes.NewInt(500))
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, uint32(1000), feeSharePpm)

	// Add more referred volume to upgrade tier
	err = k.AddReferredVolume(ctx, affiliate, dtypes.NewInt(500))
	require.NoError(t, err)

	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, uint32(2000), feeSharePpm)
}

func TestGetTakerFeeShareViaStakedAmount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	// Set up affiliate tiers
	affiliateTiers := types.AffiliateTiers{
		Tiers: []types.AffiliateTiers_Tier{
			{Level: 0, ReqReferredVolume: 0, TakerFeeSharePpm: 1000, ReqStakedWholeCoins: 200},
			{Level: 1, ReqReferredVolume: 1000, TakerFeeSharePpm: 2000, ReqStakedWholeCoins: 1000},
		},
	}
	err := k.UpdateAffiliateTiers(ctx, affiliateTiers)
	require.NoError(t, err)

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(big.NewInt(1000))))
	require.NoError(t, err)
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, uint32(2000), feeSharePpm)

	// Add more staked amount to upgrade tier
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(big.NewInt(2000))))
	require.NoError(t, err)
	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee)
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, uint32(2000), feeSharePpm)
}

func TestUpdateAffiliateTiers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	// Set up valid affiliate tiers
	validTiers := types.AffiliateTiers{
		Tiers: []types.AffiliateTiers_Tier{
			{Level: 0, ReqReferredVolume: 0, TakerFeeSharePpm: 1000},
			{Level: 1, ReqReferredVolume: 1000, TakerFeeSharePpm: 2000},
		},
	}
	err := k.UpdateAffiliateTiers(ctx, validTiers)
	require.NoError(t, err)

	// Retrieve and validate updated tiers
	updatedTiers, err := k.GetAllAffiliateTiers(ctx)
	require.NoError(t, err)
	require.Equal(t, validTiers, updatedTiers)

	// Set up invalid affiliate tiers (not sorted by level)
	invalidTiers := types.AffiliateTiers{
		Tiers: []types.AffiliateTiers_Tier{
			{Level: 1, ReqReferredVolume: 1000, TakerFeeSharePpm: 2000},
			{Level: 0, ReqReferredVolume: 0, TakerFeeSharePpm: 1000},
		},
	}
	err = k.UpdateAffiliateTiers(ctx, invalidTiers)
	require.Error(t, err)
}
