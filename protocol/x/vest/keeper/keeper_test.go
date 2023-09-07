package keeper_test

import (
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	cometbfttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
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
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	_, err := k.GetVestEntry(ctx, "not_existing_vest_entry")
	require.ErrorIs(t, err, types.ErrVestEntryNotFound)
}

func TestVestEntryStorage_Exists(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	err := k.SetVestEntry(ctx, TestValidEntry)
	require.NoError(t, err)

	got, err := k.GetVestEntry(ctx, TestVesterAccount)
	require.NoError(t, err)
	require.Equal(t, TestValidEntry, got)
}

func TestGetAllVestEntries(t *testing.T) {
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry))
	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry2))
	require.NoError(t, k.SetVestEntry(ctx, TestValidEntry3))

	gotEntries := k.GetAllVestEntries(ctx)
	expectedEntries := []types.VestEntry{
		types.DefaultGenesis().VestEntries[0],
		TestValidEntry,
		TestValidEntry2,
		TestValidEntry3,
	}

	// 1 default from genesis + 3 added
	require.Len(t, gotEntries, 4)
	for i := range gotEntries {
		require.Equal(t, expectedEntries[i], gotEntries[i])
	}
}

func TestProcessVesting(t *testing.T) {
	testVestTokenDenom := "testdenom"
	testVesterAccount := rewardstypes.VesterAccountName
	testTreasuryAccount := rewardstypes.TreasuryAccountName

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
			vesterBalance: sdk.NewInt(1_000_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(0, 0).In(time.UTC),
				EndTime:         time.Unix(1, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0).In(time.UTC),
			blockTime:               time.Unix(1001, 0).In(time.UTC),
			expectedVesterBalance:   sdk.NewInt(1_000_000),
			expectedTreasuryBalance: sdk.NewInt(0),
		},
		"vesting has ended": {
			vesterBalance: sdk.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(2000, 0).In(time.UTC),
				EndTime:         time.Unix(2001, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0).In(time.UTC),
			blockTime:               time.Unix(1001, 0).In(time.UTC),
			expectedVesterBalance:   sdk.NewInt(0),
			expectedTreasuryBalance: sdk.NewInt(0),
		},
		"vesting in progress, start_time < prev_block_time < block_time < end_time": {
			vesterBalance: sdk.NewInt(2_000_000),
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
			expectedTreasuryBalance: sdk.NewInt(2_000),
			expectedVesterBalance:   sdk.NewInt(1_998_000),
		},
		"vesting in progress, start_time < prev_block_time < block_time < end_time, rounds down": {
			vesterBalance: sdk.NewInt(2_005_000),
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
			expectedTreasuryBalance: sdk.NewInt(3_007),
			expectedVesterBalance:   sdk.NewInt(2_001_993),
		},
		"vesting in progress,  start_time < prev_block_time < block_time < end_time, vester has empty balance": {
			vesterBalance: sdk.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1000, 0),
			blockTime:               time.Unix(1001, 500_000_000),
			expectedTreasuryBalance: sdk.NewInt(0),
			expectedVesterBalance:   sdk.NewInt(0),
		},
		"vesting about to end, start_time < prev_block_time < end_time < block_time, vest all balance": {
			vesterBalance: sdk.NewInt(2_005_000),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(1999, 0).In(time.UTC),
			blockTime:               time.Unix(2001, 0).In(time.UTC),
			expectedTreasuryBalance: sdk.NewInt(2_005_000),
			expectedVesterBalance:   sdk.NewInt(0),
		},
		"vesting just started, prev_block_time < start_time < block_time < end_time": {
			vesterBalance: sdk.NewInt(2_005_000),
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
			expectedTreasuryBalance: sdk.NewInt(668),
			expectedVesterBalance:   sdk.NewInt(2_004_332),
		},
		"vesting just started, prev_block_time < start_time < block_time < end_time, vester has empty balance": {
			vesterBalance: sdk.NewInt(0),
			vestEntry: types.VestEntry{
				VesterAccount:   testVesterAccount,
				TreasuryAccount: testTreasuryAccount,
				Denom:           testVestTokenDenom,
				StartTime:       time.Unix(500, 0).In(time.UTC),
				EndTime:         time.Unix(2000, 0).In(time.UTC),
			},
			prevBlockTime:           time.Unix(499, 0),
			blockTime:               time.Unix(500, 500_000_000),
			expectedTreasuryBalance: sdk.NewInt(0),
			expectedVesterBalance:   sdk.NewInt(0),
		},
	} {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder().WithGenesisDocFn(func() (genesis cometbfttypes.GenesisDoc) {
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
			}).WithTesting(t).Build()
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
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
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
	tApp := testapp.NewTestAppBuilder().WithTesting(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VestKeeper

	err := k.DeleteVestEntry(ctx, "not_existing_vest_entry")
	require.ErrorIs(t, err, types.ErrVestEntryNotFound)
}
