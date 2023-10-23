package keeper_test

import (
	"math/big"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	big_testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/big"
	blocktimetypes "github.com/dydxprotocol/v4-chain/protocol/x/blocktime/types"
	rewardstypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
	"github.com/stretchr/testify/require"
)

const (
	TestVesterAccount   = "test_vester"
	TestTreasuryAccount = "test_treasury"
	TestDenom           = "testdenom"
)

var (
	TestValidEntry = types.VestEntry{
		TreasuryAccount: TestTreasuryAccount,
		VesterAccount:   TestVesterAccount,
		Denom:           "testdenom",
		StartTime:       time.Unix(0, 0).In(time.UTC),
		EndTime:         time.Unix(1, 0).In(time.UTC),
	}
	TestValidEntry2 = types.VestEntry{
		TreasuryAccount: "test_treasury2",
		VesterAccount:   "test_vester2",
		Denom:           "testdenom2",
		StartTime:       time.Unix(0, 0).In(time.UTC),
		EndTime:         time.Unix(1, 0).In(time.UTC),
	}
	TestValidEntry3 = types.VestEntry{
		TreasuryAccount: "test_treasury3",
		VesterAccount:   "test_vester3",
		Denom:           "testdenom3",
		StartTime:       time.Unix(0, 0).In(time.UTC),
		EndTime:         time.Unix(1, 0).In(time.UTC),
	}
)

func TestVestEntryStorage_NotFound(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	_, err := k.GetVestEntry(ctx, "not_existing_vest_entry")
	require.ErrorIs(t, err, types.ErrVestEntryNotFound)
}

func TestVestEntryStorage_Exists(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	err := k.SetVestEntry(ctx, TestValidEntry)
	require.NoError(t, err)

	got, err := k.GetVestEntry(ctx, TestVesterAccount)
	require.NoError(t, err)
	require.Equal(t, TestValidEntry, got)
}

func TestGetAllVestEntries(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry))
	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry2))
	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry3))

	gotEntries := k.GetAllVestEntries(ctx)
	expectedEntries := []types.VestEntry{
		types.DefaultGenesis().VestEntries[0],
		types.DefaultGenesis().VestEntries[1],
		TestValidEntry,
		TestValidEntry2,
		TestValidEntry3,
	}

	// 2 default from genesis + 3 added
	require.Len(t, gotEntries, 5)
	for i := range gotEntries {
		require.Equal(t, expectedEntries[i], gotEntries[i])
	}
}

func TestProcessVesting(t *testing.T) {
	testVestTokenDenom := "testdenom"
	testVesterAccount := rewardstypes.VesterAccountName
	testTreasuryAccount := rewardstypes.TreasuryAccountName
	testPrevBlockTime := time.Date(2023, 11, 5, 8, 55, 20, 0, time.UTC).In(time.UTC)
	testCurrBlockTime := time.Date(2023, 11, 5, 8, 55, 22, 0, time.UTC).In(time.UTC)

	for name, tc := range map[string]struct {
		vesterBalance           sdkmath.Int
		vestEntry               types.VestEntry
		prevBlockTime           time.Time
		blockTime               time.Time
		err                     error
		expectedVesterBalance   sdkmath.Int
		expectedTreasuryBalance sdkmath.Int
	}{
		"vesting has not started": {
			vesterBalance: sdkmath.NewInt(1_000_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(0, 0).In(time.UTC),
				EndTime:         time.Unix(1, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0).In(time.UTC),
			blockTime:               time.Unix(1001, 0).In(time.UTC),
			expectedVesterBalance:   sdkmath.NewInt(1_000_000),
			expectedTreasuryBalance: sdkmath.NewInt(0),
		},
		"vesting has ended": {
			vesterBalance: sdkmath.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(2000, 0).In(time.UTC),
				EndTime:         time.Unix(2001, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0).In(time.UTC),
			blockTime:               time.Unix(1001, 0).In(time.UTC),
			expectedVesterBalance:   sdkmath.NewInt(0),
			expectedTreasuryBalance: sdkmath.NewInt(0),
		},
		"vesting in progress, start_time < prev_block_time < block_time < end_time": {
			vesterBalance: sdkmath.NewInt(2_000_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime: time.Unix(1000, 0),
			blockTime:     time.Unix(1001, 0),
			// (1001 - 1000) / (2000 - 1000) * 1_000_000 = 1_000
			expectedTreasuryBalance: sdkmath.NewInt(2_000),
			expectedVesterBalance:   sdkmath.NewInt(1_998_000),
		},
		"vesting in progress, realistic values, start_time < prev_block_time < block_time < end_time": {
			vesterBalance: sdkmath.NewIntFromBigInt(
				big_testutil.Int64MulPow10(20_000_000, 18), // 20 million full coin, 2e26 in base denom.
			),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       types.DefaultVestingStartTime,
				EndTime:         types.DefaultVestingEndTime,
			},
			prevBlockTime: testPrevBlockTime,
			blockTime:     testCurrBlockTime,
			expectedTreasuryBalance: sdkmath.NewIntFromBigInt(
				big_testutil.MustFirst(
					new(big.Int).SetString("1095437830069111172", 10), // 1.09e18
				),
			),
			expectedVesterBalance: sdkmath.NewIntFromBigInt(
				big_testutil.MustFirst(
					new(big.Int).SetString("19999998904562169930888828", 10), // 1.99e25
				),
			),
		},
		"vesting in progress, start_time < prev_block_time < block_time < end_time, rounds down": {
			vesterBalance: sdkmath.NewInt(2_005_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime: time.Unix(1000, 0),
			blockTime:     time.Unix(1001, 500_000_000),
			// (1001.5 - 1000) / (2000 - 1000) * 2_005_000 = 3007
			expectedTreasuryBalance: sdkmath.NewInt(3_007),
			expectedVesterBalance:   sdkmath.NewInt(2_001_993),
		},
		"vesting in progress,  start_time < prev_block_time < block_time < end_time, vester has empty balance": {
			vesterBalance: sdkmath.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0),
			blockTime:               time.Unix(1001, 500_000_000),
			expectedTreasuryBalance: sdkmath.NewInt(0),
			expectedVesterBalance:   sdkmath.NewInt(0),
		},
		"vesting about to end, start_time < prev_block_time < end_time < block_time, vest all balance": {
			vesterBalance: sdkmath.NewInt(2_005_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1999, 0).In(time.UTC),
			blockTime:               time.Unix(2001, 0).In(time.UTC),
			expectedTreasuryBalance: sdkmath.NewInt(2_005_000),
			expectedVesterBalance:   sdkmath.NewInt(0),
		},
		"vesting just started, prev_block_time < start_time < block_time < end_time": {
			vesterBalance: sdkmath.NewInt(2_005_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime: time.Unix(499, 0),
			blockTime:     time.Unix(500, 500_000_000),
			// 0.5 / (2000 - 500) * 2_005_000 = 668
			expectedTreasuryBalance: sdkmath.NewInt(668),
			expectedVesterBalance:   sdkmath.NewInt(2_004_332),
		},
		"vesting just started, prev_block_time < start_time < block_time < end_time, vester has empty balance": {
			vesterBalance: sdkmath.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(499, 0),
			blockTime:               time.Unix(500, 500_000_000),
			expectedTreasuryBalance: sdkmath.NewInt(0),
			expectedVesterBalance:   sdkmath.NewInt(0),
		},
	} {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Update x/vest genesis state with test vest entry
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *types.GenesisState) {
						genesisState.VestEntries = []types.VestEntry{tc.vestEntry}
					},
				)
				// Set up vester account balance in genesis state
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *banktypes.GenesisState) {
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: authtypes.NewModuleAddress(tc.vestEntry.VesterAccount).String(),
							Coins: []sdk.Coin{
								sdk.NewCoin(testVestTokenDenom, tc.vesterBalance),
							},
						})
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()

			// Set previous block time
			tApp.App.BlockTimeKeeper.SetPreviousBlockInfo(ctx, &blocktimetypes.BlockInfo{
				Timestamp: tc.prevBlockTime,
			})

			k := tApp.App.VestKeeper

			k.ProcessVesting(ctx.WithBlockTime(tc.blockTime))

			require.Equal(t,
				tc.expectedVesterBalance,
				tApp.App.BankKeeper.GetBalance(
					ctx, authtypes.NewModuleAddress(testVesterAccount),
					testVestTokenDenom,
				).Amount,
			)

			require.Equal(t,
				tc.expectedTreasuryBalance,
				tApp.App.BankKeeper.GetBalance(
					ctx, authtypes.NewModuleAddress(testTreasuryAccount),
					testVestTokenDenom,
				).Amount,
			)
		})
	}
}

func TestDeleteVestEntry_Success(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	err := k.SetVestEntry(ctx, TestValidEntry)
	require.NoError(t, err)
	_, err = k.GetVestEntry(ctx, TestVesterAccount)
	require.NoError(t, err)

	err = k.DeleteVestEntry(ctx, TestVesterAccount)
	require.NoError(t, err)
	_, err = k.GetVestEntry(ctx, TestVesterAccount)
	require.ErrorIs(t, err, types.ErrVestEntryNotFound)
}

func TestDeleteVestEntry_NotFound(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	err := k.DeleteVestEntry(ctx, "not_existing_vest_entry")
	require.ErrorIs(t, err, types.ErrVestEntryNotFound)
}
