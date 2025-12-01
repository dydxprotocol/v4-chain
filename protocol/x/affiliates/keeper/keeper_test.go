package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gogoproto/proto"
	indexerevents "github.com/dydxprotocol/v4-chain/protocol/indexer/events"
	"github.com/dydxprotocol/v4-chain/protocol/indexer/indexer_manager"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/affiliates/types"
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
		expectError error
		setup       func(t *testing.T, ctx sdk.Context, k *keeper.Keeper)
	}{
		{
			name:        "Register new affiliate",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.BobAccAddress.String(),
			expectError: nil,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
				require.NoError(t, err)
			},
		},
		{
			name:        "Register existing referee",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.CarlAccAddress.String(),
			expectError: types.ErrAffiliateAlreadyExistsForReferee,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
				require.NoError(t, err)
				err = k.RegisterAffiliate(ctx, constants.AliceAccAddress.String(), constants.BobAccAddress.String())
				require.NoError(t, err)
			},
		},
		{
			name:        "Invalid referee address",
			referee:     "invalid_address",
			affiliate:   constants.BobAccAddress.String(),
			expectError: types.ErrInvalidAddress,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
				require.NoError(t, err)
			},
		},
		{
			name:        "Invalid affiliate address",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   "invalid_address",
			expectError: types.ErrInvalidAddress,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
				require.NoError(t, err)
			},
		},
		{
			name:        "No tiers set",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.BobAccAddress.String(),
			expectError: types.ErrAffiliateTiersNotSet,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				// No setup needed for this test case
			},
		},
		{
			name:        "Self referral",
			referee:     constants.AliceAccAddress.String(),
			affiliate:   constants.AliceAccAddress.String(),
			expectError: types.ErrSelfReferral,
			setup: func(t *testing.T, ctx sdk.Context, k *keeper.Keeper) {
				// No setup needed for this test case
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
			if tc.expectError != nil {
				require.ErrorIs(t, err, tc.expectError)
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

func TestGetTakerFeeShareViaReferredVolume(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper
	// Set up affiliate tiers
	affiliateTiers := types.DefaultAffiliateTiers
	err := k.UpdateAffiliateTiers(ctx, affiliateTiers)
	require.NoError(t, err)
	stakingKeeper := tApp.App.StakingKeeper

	require.NoError(t, k.UpdateAffiliateParameters(ctx, &types.MsgUpdateAffiliateParameters{
		Authority: constants.GovAuthority,
		AffiliateParameters: types.AffiliateParameters{
			Maximum_30DAttributableVolumePerReferredUserQuoteQuantums: 100_000_000_000_000,
		},
	}))

	require.NoError(t, stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	))

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee, map[string]bool{})
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, types.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm, feeSharePpm)

	tApp.App.StatsKeeper.SetUserStats(ctx, affiliate, &statstypes.UserStats{
		TakerNotional:                            100_000_000_000_000,
		MakerNotional:                            100_000_000_000_000,
		Affiliate_30DReferredVolumeQuoteQuantums: 1_000_000_000_001,
	})

	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee, map[string]bool{})
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, types.DefaultAffiliateTiers.Tiers[1].TakerFeeSharePpm, feeSharePpm)
}

func TestGetTakerFeeShareViaStakedAmount(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper
	ctx = ctx.WithBlockTime(time.Now())
	// Set up affiliate tiers
	affiliateTiers := types.DefaultAffiliateTiers
	err := k.UpdateAffiliateTiers(ctx, affiliateTiers)
	require.NoError(t, err)

	// Register affiliate and referee
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	)
	require.NoError(t, err)
	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	// Get taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, referee, map[string]bool{})
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, types.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm, feeSharePpm)

	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(
		time.Duration(statstypes.StakedBaseTokensCacheDurationSeconds+1) * time.Second,
	))
	// Add more staked amount to upgrade tier
	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(new(big.Int).Mul(
				big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[1].ReqStakedWholeCoins)),
				big.NewInt(1e18),
			))))
	require.NoError(t, err)
	// Get updated taker fee share for referee
	affiliateAddr, feeSharePpm, exists, err = k.GetTakerFeeShare(ctx, referee, map[string]bool{})
	require.NoError(t, err)
	require.True(t, exists)
	require.Equal(t, affiliate, affiliateAddr)
	require.Equal(t, types.DefaultAffiliateTiers.Tiers[0].TakerFeeSharePpm, feeSharePpm)
}

// Test volume qualifies for tier 2 and stake qualifies for tier 3
// should return tier 3
func TestGetTierForAffiliate_VolumeAndStake(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	affiliateTiers := types.DefaultAffiliateTiers
	err := k.UpdateAffiliateTiers(ctx, affiliateTiers)
	require.NoError(t, err)
	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	stakingKeeper := tApp.App.StakingKeeper

	err = stakingKeeper.SetDelegation(ctx,
		stakingtypes.NewDelegation(constants.AliceAccAddress.String(),
			constants.AliceValAddress.String(), math.LegacyNewDecFromBigInt(
				new(big.Int).Mul(
					big.NewInt(int64(types.DefaultAffiliateTiers.Tiers[0].ReqStakedWholeCoins)),
					big.NewInt(1e18),
				),
			),
		),
	)
	require.NoError(t, err)

	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)

	tApp.App.StatsKeeper.SetUserStats(ctx, affiliate, &statstypes.UserStats{
		TakerNotional:                            100_000_000_000_000,
		MakerNotional:                            100_000_000_000_000,
		Affiliate_30DReferredVolumeQuoteQuantums: affiliateTiers.Tiers[2].ReqReferredVolumeQuoteQuantums,
	})

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

	tierLevel, feeSharePpm, err := k.GetTierForAffiliate(ctx, affiliate, map[string]bool{})
	require.NoError(t, err)

	require.Equal(t, uint32(2), tierLevel)
	require.Equal(t, affiliateTiers.Tiers[2].TakerFeeSharePpm, feeSharePpm)
}

func TestUpdateAffiliateTiers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	tests := []struct {
		name           string
		affiliateTiers types.AffiliateTiers
		expectedError  error
	}{
		{
			name:           "Valid tiers",
			affiliateTiers: types.DefaultAffiliateTiers,
			expectedError:  nil,
		},
		{
			name: "Invalid tiers - decreasing volume requirement",
			affiliateTiers: types.AffiliateTiers{
				Tiers: []types.AffiliateTiers_Tier{
					{ReqReferredVolumeQuoteQuantums: 1000, ReqStakedWholeCoins: 100, TakerFeeSharePpm: 100},
					{ReqReferredVolumeQuoteQuantums: 500, ReqStakedWholeCoins: 200, TakerFeeSharePpm: 200},
				},
			},
			expectedError: types.ErrInvalidAffiliateTiers,
		},
		{
			name: "Taker fee share ppm greater than cap",
			affiliateTiers: types.AffiliateTiers{
				Tiers: []types.AffiliateTiers_Tier{
					{ReqReferredVolumeQuoteQuantums: 1000, ReqStakedWholeCoins: 100, TakerFeeSharePpm: 100},
					{ReqReferredVolumeQuoteQuantums: 2000, ReqStakedWholeCoins: 200, TakerFeeSharePpm: 550_000}, // 55%
				},
			},
			expectedError: types.ErrRevShareSafetyViolation,
		},
		{
			name: "Valid tiers with req_staked_whole_coins set to 0 (deprecated field)",
			affiliateTiers: types.AffiliateTiers{
				Tiers: []types.AffiliateTiers_Tier{
					{ReqReferredVolumeQuoteQuantums: 0, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 50_000},            // 5%
					{ReqReferredVolumeQuoteQuantums: 1_000_000, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 100_000},   // 10%
					{ReqReferredVolumeQuoteQuantums: 10_000_000, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 150_000},  // 15%
					{ReqReferredVolumeQuoteQuantums: 100_000_000, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 200_000}, // 20%
				},
			},
			expectedError: nil,
		},
		{
			name: "Valid tiers with mixed req_staked_whole_coins values including 0",
			affiliateTiers: types.AffiliateTiers{
				Tiers: []types.AffiliateTiers_Tier{
					{ReqReferredVolumeQuoteQuantums: 0, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 50_000},
					{ReqReferredVolumeQuoteQuantums: 1_000_000, ReqStakedWholeCoins: 0, TakerFeeSharePpm: 100_000},
					{ReqReferredVolumeQuoteQuantums: 10_000_000, ReqStakedWholeCoins: 100, TakerFeeSharePpm: 150_000},
					{ReqReferredVolumeQuoteQuantums: 100_000_000, ReqStakedWholeCoins: 500, TakerFeeSharePpm: 200_000},
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := k.UpdateAffiliateTiers(ctx, tc.affiliateTiers)

			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				// Retrieve and validate updated tiers
				updatedTiers, err := k.GetAllAffiliateTiers(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.affiliateTiers, updatedTiers)
			}
		})
	}
}

func getRegisterAffiliateEventsFromIndexerBlock(
	ctx sdk.Context,
	affiliatesKeeper *keeper.Keeper,
) []*indexerevents.RegisterAffiliateEventV1 {
	block := affiliatesKeeper.GetIndexerEventManager().ProduceBlock(ctx)
	var registerAffiliateEvents []*indexerevents.RegisterAffiliateEventV1
	for _, event := range block.Events {
		if event.Subtype != indexerevents.SubtypeRegisterAffiliate {
			continue
		}
		if _, ok := event.OrderingWithinBlock.(*indexer_manager.IndexerTendermintEvent_TransactionIndex); ok {
			var registerAffiliateEvent indexerevents.RegisterAffiliateEventV1
			err := proto.Unmarshal(event.DataBytes, &registerAffiliateEvent)
			if err != nil {
				panic(err)
			}
			registerAffiliateEvents = append(registerAffiliateEvents, &registerAffiliateEvent)
		}
	}
	return registerAffiliateEvents
}

func TestRegisterAffiliateEmitEvent(t *testing.T) {
	ctx, k, _, _ := keepertest.AffiliatesKeepers(t, true)

	affiliate := constants.AliceAccAddress.String()
	referee := constants.BobAccAddress.String()
	err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
	require.NoError(t, err)

	err = k.RegisterAffiliate(ctx, referee, affiliate)
	require.NoError(t, err)
	expectedEvent := &indexerevents.RegisterAffiliateEventV1{
		Referee:   referee,
		Affiliate: affiliate,
	}

	events := getRegisterAffiliateEventsFromIndexerBlock(ctx, k)
	require.Equal(t, 1, len(events))
	require.Equal(t, expectedEvent, events[0])
}

func TestSetAffiliateWhitelist(t *testing.T) {
	ctx, k, _, _ := keepertest.AffiliatesKeepers(t, true)

	testCases := []struct {
		name          string
		whitelist     types.AffiliateWhitelist
		expectedError error
	}{
		{
			name: "Single tier with single address",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String()},
						TakerFeeSharePpm: 100_000, // 10%
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Multiple tiers with multiple addresses",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String(), constants.BobAccAddress.String()},
						TakerFeeSharePpm: 200_000, // 20%
					},
					{
						Addresses:        []string{constants.CarlAccAddress.String()},
						TakerFeeSharePpm: 300_000, // 30%
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "Duplicate address across tiers",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String()},
						TakerFeeSharePpm: 250_000, // 25%
					},
					{
						Addresses:        []string{constants.AliceAccAddress.String()},
						TakerFeeSharePpm: 350_000, // 35%
					},
				},
			},
			expectedError: types.ErrDuplicateAffiliateAddressForWhitelist,
		},
		{
			name: "Taker fee share ppm greater than cap",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String()},
						TakerFeeSharePpm: 550_000, // 55%
					},
				},
			},
			expectedError: types.ErrRevShareSafetyViolation,
		},
		{
			name: "Invalid bech32 address present",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses: []string{
							constants.AliceAccAddress.String(),
							"dydxinvalidaddress",
						},
						TakerFeeSharePpm: 500_000, // 50%
					},
				},
			},
			expectedError: types.ErrInvalidAddress,
		},
		{
			name: "Validator operator address not accepted",
			whitelist: types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses: []string{
							constants.AliceAccAddress.String(),
							"dydxvaloper1et2kxktzr6tav65uhrxsyr8gx82xvan6gl78xd",
						},
						TakerFeeSharePpm: 500_000, // 50%
					},
				},
			},
			expectedError: types.ErrInvalidAddress,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := k.SetAffiliateWhitelist(ctx, tc.whitelist)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
			} else {
				require.NoError(t, err)

				storedWhitelist, err := k.GetAffiliateWhitelist(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.whitelist, storedWhitelist)
			}
		})
	}
}

func TestGetAffiliateWhiteListMap(t *testing.T) {
	testCases := []struct {
		name           string
		whitelist      *types.AffiliateWhitelist
		expectedLength int
		expectedMap    map[string]uint32
	}{
		{
			name: "Multiple tiers with multiple addresses",
			whitelist: &types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String(), constants.CarlAccAddress.String()},
						TakerFeeSharePpm: 100_000, // 10%
					},
					{
						Addresses:        []string{constants.BobAccAddress.String()},
						TakerFeeSharePpm: 200_000, // 20%
					},
				},
			},
			expectedLength: 3,
			expectedMap: map[string]uint32{
				constants.AliceAccAddress.String(): 100_000, // 10%
				constants.CarlAccAddress.String():  100_000, // 10%
				constants.BobAccAddress.String():   200_000, // 20%
			},
		},
		{
			name: "Single tier with single address",
			whitelist: &types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{
					{
						Addresses:        []string{constants.AliceAccAddress.String()},
						TakerFeeSharePpm: 150_000, // 15%
					},
				},
			},
			expectedLength: 1,
			expectedMap: map[string]uint32{
				constants.AliceAccAddress.String(): 150_000, // 15%
			},
		},
		{
			name: "Empty tiers",
			whitelist: &types.AffiliateWhitelist{
				Tiers: []types.AffiliateWhitelist_Tier{},
			},
			expectedLength: 0,
			expectedMap:    map[string]uint32{},
		},
		{
			name:           "tiers not set",
			whitelist:      nil,
			expectedLength: 0,
			expectedMap:    map[string]uint32{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, _, _ := keepertest.AffiliatesKeepers(t, true)

			if tc.whitelist != nil {
				err := k.SetAffiliateWhitelist(ctx, *tc.whitelist)
				require.NoError(t, err)
			}

			whitelistMap, err := k.GetAffiliateWhitelistMap(ctx)
			require.NoError(t, err)
			require.Equal(t, tc.expectedLength, len(whitelistMap))
			require.Equal(t, tc.expectedMap, whitelistMap)
		})
	}
}

func TestGetTakerFeeShareViaWhitelist(t *testing.T) {
	tiers := types.DefaultAffiliateTiers

	testCases := []struct {
		name                string
		affiliateAddr       string
		refereeAddr         string
		overrides           *types.AffiliateOverrides
		expectedFeeSharePpm uint32
		expectedExists      bool
	}{
		{
			name:          "Affiliate in whitelist",
			affiliateAddr: constants.AliceAccAddress.String(),
			refereeAddr:   constants.BobAccAddress.String(),
			overrides: &types.AffiliateOverrides{
				Addresses: []string{constants.AliceAccAddress.String()},
			},
			expectedFeeSharePpm: 250_000, // 25%
			expectedExists:      true,
		},
		{
			name:                "Affiliate not in whitelist",
			affiliateAddr:       constants.AliceAccAddress.String(),
			refereeAddr:         constants.BobAccAddress.String(),
			overrides:           &types.AffiliateOverrides{},
			expectedFeeSharePpm: tiers.Tiers[0].TakerFeeSharePpm,
			expectedExists:      true,
		},
		{
			name:          "Referee not registered",
			affiliateAddr: "",
			refereeAddr:   constants.BobAccAddress.String(),
			overrides: &types.AffiliateOverrides{
				Addresses: []string{constants.AliceAccAddress.String()},
			},
			expectedFeeSharePpm: 0,
			expectedExists:      false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, k, _, _ := keepertest.AffiliatesKeepers(t, true)
			err := k.UpdateAffiliateTiers(ctx, tiers)
			require.NoError(t, err)

			if tc.overrides != nil {
				err := k.SetAffiliateOverrides(ctx, *tc.overrides)
				require.NoError(t, err)
			}
			if tc.affiliateAddr != "" {
				err := k.RegisterAffiliate(ctx, tc.refereeAddr, tc.affiliateAddr)
				require.NoError(t, err)
			}
			affiliateOverridesMap, err := k.GetAffiliateOverridesMap(ctx)
			require.NoError(t, err)

			affiliateAddr, feeSharePpm, exists, err := k.GetTakerFeeShare(ctx, tc.refereeAddr, affiliateOverridesMap)
			require.NoError(t, err)
			require.Equal(t, tc.affiliateAddr, affiliateAddr)
			require.Equal(t, tc.expectedFeeSharePpm, feeSharePpm)
			require.Equal(t, tc.expectedExists, exists)
		})
	}
}

func TestGetTierForAffiliateEmptyTiers(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	tierLevel, feeSharePpm, err := k.GetTierForAffiliate(ctx, constants.AliceAccAddress.String(), map[string]bool{})
	require.NoError(t, err)
	require.Equal(t, uint32(0), tierLevel)
	require.Equal(t, uint32(0), feeSharePpm)
}

func TestUpdateAffiliateParameters(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	err := k.UpdateAffiliateParameters(ctx, &types.MsgUpdateAffiliateParameters{
		Authority:           constants.GovAuthority,
		AffiliateParameters: types.DefaultAffiliateParameters,
	})
	require.NoError(t, err)

	affiliateParameters, err := k.GetAffiliateParameters(ctx)
	require.NoError(t, err)
	require.Equal(
		t,
		uint64(100_000_000_000_000),
		affiliateParameters.GetMaximum_30DAttributableVolumePerReferredUserQuoteQuantums(),
	)
	require.Equal(t, uint32(2), affiliateParameters.GetRefereeMinimumFeeTierIdx())
	require.Equal(
		t,
		uint64(10_000_000_000),
		affiliateParameters.GetMaximum_30DAffiliateRevenuePerReferredUserQuoteQuantums(),
	)
}

func TestGetTierForAffiliateOverrides(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.AffiliatesKeeper

	err := k.UpdateAffiliateTiers(ctx, types.DefaultAffiliateTiers)
	require.NoError(t, err)

	tierLevel, feeSharePpm, err := k.GetTierForAffiliate(ctx, constants.AliceAccAddress.String(), map[string]bool{})
	require.NoError(t, err)
	require.Equal(t, uint32(0), tierLevel)
	require.Equal(t, uint32(50_000), feeSharePpm)

	tierLevel, feeSharePpm, err = k.GetTierForAffiliate(ctx, constants.AliceAccAddress.String(), map[string]bool{
		constants.AliceAccAddress.String(): true,
	})
	require.NoError(t, err)
	require.Equal(t, uint32(4), tierLevel)
	require.Equal(t, uint32(250_000), feeSharePpm)
}
