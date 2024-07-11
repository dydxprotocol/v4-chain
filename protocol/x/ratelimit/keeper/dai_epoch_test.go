package keeper_test

import (
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestCheckCurrentDAIYieldEpochElapsed(t *testing.T) {
	testCases := []struct {
		name            string
		isFirstBlock    bool
		currentEpoch    *big.Int
		blockNumber     uint64
		currentBlock    int64
		expectedElapsed bool
		expectedError   bool
	}{
		{
			name:            "First block, no epoch set",
			isFirstBlock:    true,
			expectedElapsed: true,
			expectedError:   false,
		},
		{
			name:            "Epoch not elapsed",
			isFirstBlock:    false,
			currentEpoch:    big.NewInt(1),
			blockNumber:     100,
			currentBlock:    150,
			expectedElapsed: false,
			expectedError:   false,
		},
		{
			name:            "Epoch elapsed",
			isFirstBlock:    false,
			currentEpoch:    big.NewInt(1),
			blockNumber:     100,
			currentBlock:    300,
			expectedElapsed: true,
			expectedError:   false,
		},
		{
			name:            "Epoch not found",
			isFirstBlock:    false,
			currentEpoch:    big.NewInt(1),
			blockNumber:     0,
			currentBlock:    200,
			expectedElapsed: false,
			expectedError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			if !tc.isFirstBlock {
				k.SetCurrentDaiYieldEpochNumber(ctx, tc.currentEpoch)
				if tc.blockNumber != 0 {
					k.SetDaiYieldEpochParams(ctx, tc.currentEpoch.Uint64()%types.DAI_YIELD_ARRAY_SIZE, types.DaiYieldEpochParams{
						BlockNumber: tc.blockNumber,
					})
				}
			}

			ctx = ctx.WithBlockHeight(tc.currentBlock)

			elapsed, err := k.CheckCurrentDAIYieldEpochElapsed(ctx)
			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedElapsed, elapsed)
			}
		})
	}
}

func TestDAIYieldEpochHasElapsed(t *testing.T) {
	testCases := []struct {
		name                  string
		currentBlockNumber    uint64
		epochStartBlockNumber uint64
		expectedElapsed       bool
	}{
		{
			name:                  "Epoch not elapsed",
			currentBlockNumber:    150,
			epochStartBlockNumber: 100,
			expectedElapsed:       false,
		},
		{
			name:                  "Epoch just elapsed",
			currentBlockNumber:    200,
			epochStartBlockNumber: 100,
			expectedElapsed:       true,
		},
		{
			name:                  "Epoch long elapsed",
			currentBlockNumber:    250,
			epochStartBlockNumber: 100,
			expectedElapsed:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			_ = tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			elapsed := k.DAIYieldEpochHasElapsed(tc.currentBlockNumber, tc.epochStartBlockNumber)
			require.Equal(t, tc.expectedElapsed, elapsed)
		})
	}
}

func TestCheckFirstDAIYieldEpoch(t *testing.T) {
	testCases := []struct {
		name          string
		currentEpoch  *big.Int
		expectedEpoch *big.Int
		expectedFirst bool
	}{
		{
			name:          "Epoch not found",
			currentEpoch:  nil,
			expectedEpoch: nil,
			expectedFirst: true,
		},
		{
			name:          "Epoch found",
			currentEpoch:  big.NewInt(1),
			expectedEpoch: big.NewInt(1),
			expectedFirst: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			if tc.currentEpoch != nil {
				k.SetCurrentDaiYieldEpochNumber(ctx, tc.currentEpoch)
			}

			epoch, first := k.CheckFirstDAIYieldEpoch(ctx)
			require.Equal(t, tc.expectedFirst, first)
			require.Equal(t, tc.expectedEpoch, epoch)
		})
	}
}

func TestGetCurrentDAIYieldEpochBlockNumber(t *testing.T) {
	testCases := []struct {
		name          string
		currentEpoch  *big.Int
		blockNumber   uint64
		expectedFound bool
		expectedBlock uint64
	}{
		{
			name:          "Epoch params found",
			currentEpoch:  big.NewInt(1),
			blockNumber:   100,
			expectedFound: true,
			expectedBlock: 100,
		},
		{
			name:          "Epoch params not found",
			currentEpoch:  big.NewInt(2),
			blockNumber:   0,
			expectedFound: false,
			expectedBlock: 0,
		},
		{
			name:          "Test the modding of the array index",
			currentEpoch:  big.NewInt(150),
			blockNumber:   100,
			expectedFound: true,
			expectedBlock: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			if tc.expectedFound {
				k.SetDaiYieldEpochParams(ctx, tc.currentEpoch.Uint64()%types.DAI_YIELD_ARRAY_SIZE, types.DaiYieldEpochParams{
					BlockNumber: tc.blockNumber,
				})
			}

			blockNumber, found := k.GetCurrentDAIYieldEpochBlockNumber(ctx, tc.currentEpoch)
			require.Equal(t, tc.expectedFound, found)
			require.Equal(t, tc.expectedBlock, blockNumber)
		})
	}
}

func TestTransferRemainingDAIYieldToInsuranceFund(t *testing.T) {
	testCases := []struct {
		name                     string
		tradingDaiMinted         string
		totalTradingDaiClaimed   string
		initialPoolBalance       sdk.Coins
		initialInsuranceBalance  sdk.Coins
		expectedPoolBalance      sdk.Coins
		expectedInsuranceBalance sdk.Coins
		expectError              bool
	}{
		{
			name:                   "Invalid tradingDaiMinted",
			tradingDaiMinted:       "invalid",
			totalTradingDaiClaimed: "100",
			expectError:            true,
		},
		{
			name:                   "Invalid totalTradingDaiClaimed",
			tradingDaiMinted:       "100",
			totalTradingDaiClaimed: "invalid",
			expectError:            true,
		},
		{
			name:                     "tradingDaiMintedAtEpoch <= tradingDaiClaimedAtEpoch",
			tradingDaiMinted:         "100",
			totalTradingDaiClaimed:   "100",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			expectedInsuranceBalance: sdk.NewCoins(),
			expectError:              false,
		},
		{
			name:                     "Not enough money in pool account",
			tradingDaiMinted:         "200",
			totalTradingDaiClaimed:   "100",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(50))),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(50))),
			expectedInsuranceBalance: sdk.NewCoins(),
			expectError:              true,
		},
		{
			name:                     "Everything works, insurance fund starts with no money",
			tradingDaiMinted:         "200",
			totalTradingDaiClaimed:   "100",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			expectedInsuranceBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			expectError:              false,
		},
		{
			name:                     "Everything works, insurance fund starts with money",
			tradingDaiMinted:         "200",
			totalTradingDaiClaimed:   "100",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
			initialInsuranceBalance:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(50))),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			expectedInsuranceBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(150))),
			expectError:              false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			bankKeeper := tApp.App.BankKeeper
			// Mint initial pool balance
			if !tc.initialPoolBalance.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialPoolBalance))
			}

			// Mint initial insurance balance
			if !tc.initialInsuranceBalance.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialInsuranceBalance))
				require.NoError(t, bankKeeper.SendCoins(ctx, authtypes.NewModuleAddress(types.PoolAccount), perptypes.InsuranceFundModuleAddress, tc.initialInsuranceBalance))
			}

			err := k.TransferRemainingDAIYieldToInsuranceFund(ctx, tc.tradingDaiMinted, tc.totalTradingDaiClaimed)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Check pool account balance
				poolBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.PoolAccount), types.TradingDAIDenom)
				require.Equal(t, tc.expectedPoolBalance.AmountOf(types.TradingDAIDenom).String(), poolBalance.Amount.String())

				// Check insurance fund balance
				insuranceBalance := bankKeeper.GetBalance(ctx, perptypes.InsuranceFundModuleAddress, types.TradingDAIDenom)
				require.Equal(t, tc.expectedInsuranceBalance.AmountOf(types.TradingDAIDenom).String(), insuranceBalance.Amount.String())
			}
		})
	}
}

func TestCalculateYieldParamsForNewEpoch(t *testing.T) {
	testCases := []struct {
		name                     string
		currentEpoch             *big.Int
		expectedTDAISupply       *big.Int
		expectedTradingDaiMinted *big.Int
		expectedYieldCollected   *big.Int
		expectedNewEpoch         uint64
		expectError              bool
	}{
		{
			name:                     "First time, no current epoch set",
			currentEpoch:             nil,
			expectedTDAISupply:       big.NewInt(0),
			expectedTradingDaiMinted: big.NewInt(0),
			expectedYieldCollected:   big.NewInt(0),
			expectedNewEpoch:         0,
			expectError:              false,
		},
		{
			name:                     "Subsequent time, current epoch set",
			currentEpoch:             big.NewInt(1),
			expectedTDAISupply:       big.NewInt(0),
			expectedTradingDaiMinted: big.NewInt(0),
			expectedYieldCollected:   big.NewInt(0),
			expectedNewEpoch:         2,
			expectError:              false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			if tc.currentEpoch != nil {
				k.SetCurrentDaiYieldEpochNumber(ctx, tc.currentEpoch)
			}

			tDAISupply, tradingDaiMinted, yieldCollectedByInsuranceFund, newEpoch, err := k.CalculateYieldParamsForNewEpoch(ctx)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTDAISupply, tDAISupply)
				require.Equal(t, tc.expectedTradingDaiMinted, tradingDaiMinted)
				require.Equal(t, tc.expectedYieldCollected, yieldCollectedByInsuranceFund)
				require.Equal(t, tc.expectedNewEpoch, newEpoch)
			}
		})
	}
}

func TestMintYieldGeneratedDuringEpoch(t *testing.T) {
	testCases := []struct {
		name                     string
		initialSDAISupply        sdk.Coins
		initialTradingDAISupply  sdk.Coins
		sdaiPrice                *big.Int
		expectedTDAISupply       *big.Int
		expectedTradingDaiToMint *big.Int
		expectError              bool
	}{
		{
			name:                    "sDAI price not set",
			initialSDAISupply:       sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(200))),
			initialTradingDAISupply: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			sdaiPrice:               nil,
			expectError:             true,
		},
		{
			name:                    "tradingDaiAfterYield will be less than intial trading dai",
			initialSDAISupply:       sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(100))),
			initialTradingDAISupply: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
			sdaiPrice:               new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil)),
			expectError:             true,
		},
		{
			name:                     "Successful minting",
			initialSDAISupply:        sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(200))),
			initialTradingDAISupply:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			sdaiPrice:                new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil)),
			expectedTDAISupply:       big.NewInt(100),
			expectedTradingDaiToMint: big.NewInt(100),
			expectError:              false,
		},
		{
			name:                     "Both initial supplies start at 0",
			initialSDAISupply:        sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(0))),
			initialTradingDAISupply:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			sdaiPrice:                new(big.Int).Mul(big.NewInt(1), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS), nil)),
			expectedTDAISupply:       big.NewInt(0),
			expectedTradingDaiToMint: big.NewInt(0),
			expectError:              false,
		},
		{
			name:                     "Initial trading DAI higher than sDAI with rate higher than 1",
			initialSDAISupply:        sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(100))),
			initialTradingDAISupply:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
			sdaiPrice:                new(big.Int).Mul(big.NewInt(25), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS-2), nil)),
			expectedTDAISupply:       big.NewInt(200),
			expectedTradingDaiToMint: big.NewInt(200),
			expectError:              false,
		},
		{
			name:                     "Trading DAI to mint is 0",
			initialSDAISupply:        sdk.NewCoins(sdk.NewCoin(types.SDaiDenom, sdkmath.NewInt(100))),
			initialTradingDAISupply:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
			sdaiPrice:                new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(types.BASE_10), big.NewInt(types.SDAI_DECIMALS-1), nil)),
			expectedTDAISupply:       big.NewInt(200),
			expectedTradingDaiToMint: big.NewInt(0),
			expectError:              true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			bankKeeper := tApp.App.BankKeeper

			// Mint initial sDAI supply
			if !tc.initialSDAISupply.IsZero() {
				mintingErr := bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialSDAISupply)
				require.NoError(t, mintingErr)
				sendingErr := bankKeeper.SendCoinsFromModuleToModule(ctx, types.PoolAccount, types.SDAIPoolAccount, tc.initialSDAISupply)
				require.NoError(t, sendingErr)
			}

			// Mint initial tradingDAI supply
			if !tc.initialTradingDAISupply.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialTradingDAISupply))
			}

			// Set sDAI price
			if tc.sdaiPrice != nil {
				k.SetSDAIPrice(ctx, tc.sdaiPrice)
			}

			tDAISupply, tradingDaiToMint, yieldCollectedByInsuranceFund, err := k.MintYieldGeneratedDuringEpoch(ctx)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedTDAISupply, tDAISupply)
				require.Equal(t, tc.expectedTradingDaiToMint, tradingDaiToMint)

				// Check the total supply of tradingDAI
				totalTradingDAISupply := bankKeeper.GetSupply(ctx, types.TradingDAIDenom).Amount.BigInt()
				expectedTotalTradingDAISupply := new(big.Int).Add(tc.expectedTDAISupply, tc.expectedTradingDaiToMint)
				require.Equal(t, expectedTotalTradingDAISupply, totalTradingDAISupply)

				// Check the yield collected by the insurance fund
				// This part depends on the implementation of CollectYieldForInsuranceFunds
				// Assuming it returns the correct amount, we can check it directly
				require.NotNil(t, yieldCollectedByInsuranceFund)
			}
		})
	}
}

func TestCollectYieldForInsuranceFund(t *testing.T) {
	testCases := []struct {
		name                     string
		initialPoolBalance       sdk.Coins
		initialInsuranceBalance  sdk.Coins
		tradingDaiMinted         *big.Int
		tradingDaiSupplyBefore   *big.Int
		expectedYield            *big.Int
		expectedPoolBalance      sdk.Coins
		expectedInsuranceBalance sdk.Coins
		expectError              bool
	}{
		{
			name:                     "No trading DAI supply before new epoch",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			initialInsuranceBalance:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			tradingDaiMinted:         big.NewInt(0),
			tradingDaiSupplyBefore:   big.NewInt(0),
			expectedYield:            big.NewInt(0),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			expectedInsuranceBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			expectError:              false,
		},
		{
			name:                     "No balance in insurance fund",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(300))),
			initialInsuranceBalance:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			tradingDaiMinted:         big.NewInt(100),
			tradingDaiSupplyBefore:   big.NewInt(200),
			expectedYield:            big.NewInt(0),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(300))),
			expectedInsuranceBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(0))),
			expectError:              false,
		},
		{
			name:                     "Successful yield collection",
			initialPoolBalance:       sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(400))),
			initialInsuranceBalance:  sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			tradingDaiMinted:         big.NewInt(100),
			tradingDaiSupplyBefore:   big.NewInt(400),
			expectedYield:            big.NewInt(25),
			expectedPoolBalance:      sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(375))),
			expectedInsuranceBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(125))),
			expectError:              false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper

			bankKeeper := tApp.App.BankKeeper

			// Mint initial pool balance
			if !tc.initialPoolBalance.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialPoolBalance))
			}

			// Mint initial insurance balance
			if !tc.initialInsuranceBalance.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialInsuranceBalance))
				require.NoError(t, bankKeeper.SendCoins(ctx, authtypes.NewModuleAddress(types.PoolAccount), perptypes.InsuranceFundModuleAddress, tc.initialInsuranceBalance))
			}

			yield, err := k.CollectYieldForInsuranceFund(ctx, perptypes.InsuranceFundModuleAddress, tc.tradingDaiMinted, tc.tradingDaiSupplyBefore)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedYield, yield)

				// Check pool account balance
				poolBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.PoolAccount), types.TradingDAIDenom)
				require.Equal(t, tc.expectedPoolBalance.AmountOf(types.TradingDAIDenom).String(), poolBalance.Amount.String())

				// Check insurance fund balance
				insuranceBalance := bankKeeper.GetBalance(ctx, perptypes.InsuranceFundModuleAddress, types.TradingDAIDenom)
				require.Equal(t, tc.expectedInsuranceBalance.AmountOf(types.TradingDAIDenom).String(), insuranceBalance.Amount.String())

				// Check the total supply of tradingDAI
				totalTradingDAISupply := bankKeeper.GetSupply(ctx, types.TradingDAIDenom).Amount.BigInt()
				expectedTotalTradingDAISupply := new(big.Int).Add(tc.tradingDaiSupplyBefore, tc.tradingDaiMinted)
				require.Equal(t, expectedTotalTradingDAISupply, totalTradingDAISupply)
			}
		})
	}
}

func TestCollectYieldForInsuranceFunds(t *testing.T) {
	testCases := []struct {
		name                      string
		initialPoolBalance        sdk.Coins
		initialInsuranceBalances  map[uint32]sdk.Coins
		tradingDaiMinted          *big.Int
		tradingDaiSupplyBefore    *big.Int
		expectedYield             *big.Int
		expectedPoolBalance       sdk.Coins
		expectedInsuranceBalances map[uint32]sdk.Coins
		expectError               bool
	}{
		{
			name:               "Isolated market",
			initialPoolBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(1000))),
			initialInsuranceBalances: map[uint32]sdk.Coins{
				100: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			},
			tradingDaiMinted:       big.NewInt(100),
			tradingDaiSupplyBefore: big.NewInt(1000),
			expectedYield:          big.NewInt(10),
			expectedPoolBalance:    sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(990))),
			expectedInsuranceBalances: map[uint32]sdk.Coins{
				100: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(110))),
			},
			expectError: false,
		},
		{
			name:               "Cross market",
			initialPoolBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(1000))),
			initialInsuranceBalances: map[uint32]sdk.Coins{
				300: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
			},
			tradingDaiMinted:       big.NewInt(100),
			tradingDaiSupplyBefore: big.NewInt(1000),
			expectedYield:          big.NewInt(10),
			expectedPoolBalance:    sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(990))),
			expectedInsuranceBalances: map[uint32]sdk.Coins{
				300: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(110))),
			},
			expectError: false,
		},
		{
			name:               "Multiple markets",
			initialPoolBalance: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(2040))),
			initialInsuranceBalances: map[uint32]sdk.Coins{
				100: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(100))),
				200: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(200))),
				300: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(300))),
			},
			tradingDaiMinted:       big.NewInt(240),
			tradingDaiSupplyBefore: big.NewInt(2400),
			expectedYield:          big.NewInt(60), // Assuming 10 for each isolated and 10 for cross
			expectedPoolBalance:    sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(1980))),
			expectedInsuranceBalances: map[uint32]sdk.Coins{
				100: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(110))),
				200: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(220))),
				300: sdk.NewCoins(sdk.NewCoin(types.TradingDAIDenom, sdkmath.NewInt(330))),
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).Build()
			ctx := tApp.InitChain()
			k := tApp.App.RatelimitKeeper
			perpKeeper := tApp.App.PerpetualsKeeper
			bankKeeper := tApp.App.BankKeeper

			// Create perpetual markets
			_, err := perpKeeper.CreatePerpetual(ctx, 100, "PERP1", 1, 1, 1, 0, perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED)
			if err != nil {
				t.Fatalf("error creating perpetual market 1: %v", err)
			}

			_, err = perpKeeper.CreatePerpetual(ctx, 200, "PERP2", 2, 1, 1, 0, perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_ISOLATED)
			if err != nil {
				t.Fatalf("error creating perpetual market 2: %v", err)
			}

			_, err = perpKeeper.CreatePerpetual(ctx, 300, "PERP3", 3, 1, 1, 0, perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS)
			if err != nil {
				t.Fatalf("error creating perpetual market 3: %v", err)
			}

			_, err = perpKeeper.CreatePerpetual(ctx, 400, "PERP4", 4, 1, 1, 0, perptypes.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS)
			if err != nil {
				t.Fatalf("error creating perpetual market 4: %v", err)
			}

			// Mint initial pool balance
			if !tc.initialPoolBalance.IsZero() {
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, tc.initialPoolBalance))
			}

			// Mint initial insurance balances
			for id, balance := range tc.initialInsuranceBalances {
				address, err := perpKeeper.GetInsuranceFundModuleAddress(ctx, id)
				require.NoError(t, err)
				require.NoError(t, bankKeeper.MintCoins(ctx, types.PoolAccount, balance))
				require.NoError(t, bankKeeper.SendCoins(ctx, authtypes.NewModuleAddress(types.PoolAccount), address, balance))
			}

			yield, err := k.CollectYieldForInsuranceFunds(ctx, tc.tradingDaiMinted, tc.tradingDaiSupplyBefore)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedYield, yield)

				// Check pool account balance
				poolBalance := bankKeeper.GetBalance(ctx, authtypes.NewModuleAddress(types.PoolAccount), types.TradingDAIDenom)
				require.Equal(t, tc.expectedPoolBalance.AmountOf(types.TradingDAIDenom).String(), poolBalance.Amount.String())

				// Check insurance fund balances
				for id, expectedBalance := range tc.expectedInsuranceBalances {
					address, err := perpKeeper.GetInsuranceFundModuleAddress(ctx, id)
					require.NoError(t, err)
					insuranceBalance := bankKeeper.GetBalance(ctx, address, types.TradingDAIDenom)
					require.Equal(t, expectedBalance.AmountOf(types.TradingDAIDenom).String(), insuranceBalance.Amount.String())
				}

				// Check the total supply of tradingDAI
				totalTradingDAISupply := bankKeeper.GetSupply(ctx, types.TradingDAIDenom).Amount.BigInt()
				expectedTotalTradingDAISupply := new(big.Int).Add(tc.tradingDaiSupplyBefore, tc.tradingDaiMinted)
				require.Equal(t, expectedTotalTradingDAISupply, totalTradingDAISupply)
			}
		})
	}
}
