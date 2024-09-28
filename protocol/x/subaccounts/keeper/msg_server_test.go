package keeper_test

import (
	"context"
	"math/big"
	"testing"

	sdkmath "cosmossdk.io/math"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	bank_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/keeper"
	asstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	ratelimitkeeper "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/keeper"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t *testing.T) (keeper.Keeper, types.MsgServer, context.Context) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.SubaccountsKeeper

	return k, keeper.NewMsgServerImpl(k), ctx
}

func TestMsgServer(t *testing.T) {
	k, ms, ctx := setupMsgServer(t)
	require.NotNil(t, k)
	require.NotNil(t, ms)
	require.NotNil(t, ctx)
}

func TestClaimYieldForSubaccount(t *testing.T) {
	// default subaccount id, the first subaccount id generated when calling createNSubaccount
	defaultSubaccountId := types.SubaccountId{
		Owner:  "0",
		Number: 0,
	}

	tests := map[string]struct {
		// state
		perpetuals []perptypes.Perpetual
		assets     []*asstypes.Asset
		// Only set when specified. Defaults to 0/1.
		// Set perpYieldIndex in the perpetuals state.
		globalAssetYieldIndex *big.Rat
		fundsInTDaiPool       *big.Int

		// subaccount state
		perpetualPositions        []*types.PerpetualPosition
		assetPositions            []*types.AssetPosition
		subaccountAssetYieldIndex string

		// collateral pool state
		collateralPoolTDaiBalances map[string]int64

		// extra test state
		msgClaimYieldForSubaccount types.MsgClaimYieldForSubaccount

		// expectations
		expectedCollateralPoolTDaiBalances map[string]int64
		expectedPerpetualPositions         []*types.PerpetualPosition
		expectedAssetPositions             []*types.AssetPosition
		expectedTDaiYieldPoolBalance       *big.Int
		expectedErr                        error
		expectedAssetYieldIndex            string
	}{
		"Successfully claims yield for tDai asset position and no other position exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(2, 1).String(),
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset position when perp with no yield exists": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(2, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(2, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_000_000_000,
			},
		},
		"Successfully claims yield for one perp with no asset positions existing before yield claim": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(1_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(199_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 101_000_000_000,
			},
		},
		"Successfully claims yield for one perp with asset position existing but not claiming yield": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(1_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(101_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(199_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 101_000_000_000,
			},
		},
		"Successfully claims yield for tDai asset and one perp": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 2).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(200_100_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(99_900_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 200_100_000_000,
			},
		},
		"Successfully claims yield when multiple perp positions are open and tDai position open": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(50_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(13, 11).String(),
			globalAssetYieldIndex:     big.NewRat(26, 11),
			fundsInTDaiPool:           big.NewInt(222_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 50_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					Params:       constants.EthUsd_NoMarginRequirement.Params,
					FundingIndex: constants.EthUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.EthUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 9).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(26, 11).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(11, 3).String(),
				},
				{
					PerpetualId:  uint32(1),
					Quantums:     dtypes.NewInt(-2_000_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 3).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_094_444_443), // Total Yield: 50_094_444_443
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(171_905_555_557),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_094_444_443,
			},
		},
		"Successfully claims all yield in tDaiPool": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(3, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(3, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(300_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 300_000_000_000,
			},
		},
		"Successfully claims yield for isolated market": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(5, 4),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.IsoUsd_IsolatedMarket.Params,
					FundingIndex: constants.IsoUsd_IsolatedMarket.FundingIndex,
					OpenInterest: constants.IsoUsd_IsolatedMarket.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(5, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(3),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(124_920_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(175_080_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				authtypes.NewModuleAddress(
					types.ModuleName + ":" + lib.UintToString(constants.PerpetualPosition_OneISOLong.PerpetualId),
				).String(): 124_920_000_000,
				types.ModuleAddress.String(): 0,
			},
		},
		"Successfully does not claim yield when asset yield index is already updated": {
			globalAssetYieldIndex: big.NewRat(5, 4),
			fundsInTDaiPool:       big.NewInt(1_200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
			subaccountAssetYieldIndex: big.NewRat(5, 4).String(),
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(5, 4).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(4, 5).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000), // Yield Collected: 0 tDAI
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(1_200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 1_000_000_000_000,
			},
		},
		"Succesfully does not claim yield when negative positions cancel out positive position yield claims": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)),
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1000, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Negative general asset yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(-1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex:    big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Asset yield index in account higher than in general ": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: big.NewRat(1, 1).String(),
			globalAssetYieldIndex:     big.NewRat(1, 2),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				constants.BtcUsd_NoMarginRequirement,
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex:    big.NewRat(1, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Negative general perp yield index": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(-1, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                types.ErrGlobalYieldIndexNegative,
			expectedAssetYieldIndex:    big.NewRat(-1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Perp yield index in subaccount higher than in general": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                types.ErrGeneralYieldIndexSmallerThanYieldIndexInSubaccount,
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: Perp yield index in subaccount badly initialized": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                types.ErrYieldIndexUninitialized,
			expectedAssetYieldIndex:    big.NewRat(0, 2).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   "",
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Succesfull yield claim: not enough yield in tdai pool so we take available": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(1),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_001),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(0),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_001,
			},
		},
		"Fails yield claim: no open positions": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			msgClaimYieldForSubaccount:   types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedErr:                  types.ErrNoYieldToClaim,
			expectedAssetYieldIndex:      big.NewRat(1, 1).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: subaccountId is nil": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			msgClaimYieldForSubaccount:   types.MsgClaimYieldForSubaccount{},
			expectedErr:                  types.ErrSubaccountIdIsNil,
			expectedAssetYieldIndex:      big.NewRat(1, 1).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Fails yield claim: subaccount with yield claim exists, but different one with no positions is passed in": {
			globalAssetYieldIndex: big.NewRat(1, 1),
			fundsInTDaiPool:       big.NewInt(100_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(1, 2).String(),
				},
			},
			msgClaimYieldForSubaccount:   types.MsgClaimYieldForSubaccount{Id: &types.SubaccountId{Owner: "0", Number: 1}},
			expectedErr:                  types.ErrNoYieldToClaim,
			expectedAssetYieldIndex:      big.NewRat(1, 1).String(),
			expectedTDaiYieldPoolBalance: big.NewInt(100_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
		"Successfully claims 0 yield when subaccount's yield is negative": {
			assetPositions:            testutil.CreateTDaiAssetPosition(big.NewInt(100_000_000_000)), // $100,000
			subaccountAssetYieldIndex: constants.AssetYieldIndex_Zero,
			globalAssetYieldIndex:     big.NewRat(1, 1),
			fundsInTDaiPool:           big.NewInt(200_000_000_000),
			collateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
			perpetuals: []perptypes.Perpetual{
				{
					Params:       constants.BtcUsd_NoMarginRequirement.Params,
					FundingIndex: constants.BtcUsd_NoMarginRequirement.FundingIndex,
					OpenInterest: constants.BtcUsd_NoMarginRequirement.OpenInterest,
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			perpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(0, 1).String(),
				},
			},
			msgClaimYieldForSubaccount: types.MsgClaimYieldForSubaccount{Id: &defaultSubaccountId},
			expectedAssetYieldIndex:    big.NewRat(1, 1).String(),
			expectedPerpetualPositions: []*types.PerpetualPosition{
				{
					PerpetualId:  uint32(0),
					Quantums:     dtypes.NewInt(-100_000_000),
					FundingIndex: dtypes.NewInt(0),
					YieldIndex:   big.NewRat(10_000, 1).String(),
				},
			},
			expectedAssetPositions: []*types.AssetPosition{
				{
					AssetId:  uint32(0),
					Quantums: dtypes.NewInt(100_000_000_000),
				},
			},
			expectedTDaiYieldPoolBalance: big.NewInt(200_000_000_000),
			expectedCollateralPoolTDaiBalances: map[string]int64{
				types.ModuleAddress.String(): 100_000_000_000,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, subaccountsKeeper, pricesKeeper, perpetualsKeeper, accountKeeper, bankKeeper, assetsKeeper, rateLimitKeeper, _, _ := testutil.SubaccountsKeepers(
				t,
				true,
			)
			ctx = ctx.WithTxBytes(constants.TestTxBytes)
			testutil.CreateTestMarkets(t, ctx, pricesKeeper)
			testutil.CreateTestLiquidityTiers(t, ctx, perpetualsKeeper)

			// Set up initial sdai price
			rateString := sdaiservertypes.TestSDAIEventRequest.ConversionRate
			rate, conversionErr := ratelimitkeeper.ConvertStringToBigInt(rateString)
			require.NoError(t, conversionErr)

			rateLimitKeeper.SetSDAIPrice(ctx, rate)
			globalAssetYieldIndex := big.NewRat(1, 1)
			if tc.globalAssetYieldIndex != nil {
				globalAssetYieldIndex = tc.globalAssetYieldIndex
			}
			rateLimitKeeper.SetAssetYieldIndex(ctx, globalAssetYieldIndex)

			// Always creates TDai asset first
			require.NoError(t, testutil.CreateTDaiAsset(ctx, assetsKeeper))
			for _, a := range tc.assets {
				_, err := assetsKeeper.CreateAsset(
					ctx,
					a.Id,
					a.Symbol,
					a.Denom,
					a.DenomExponent,
					a.HasMarket,
					a.MarketId,
					a.AtomicResolution,
					a.AssetYieldIndex,
				)
				require.NoError(t, err)
			}

			for _, p := range tc.perpetuals {
				perpetualsKeeper.SetPerpetualForTest(
					ctx,
					p,
				)
			}

			for collateralPoolAddr, TDaiBal := range tc.collateralPoolTDaiBalances {
				err := bank_testutil.FundAccount(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(TDaiBal)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			if tc.fundsInTDaiPool != nil {
				err := bank_testutil.FundModuleAccount(
					ctx,
					ratelimittypes.TDaiPoolAccount,
					sdk.Coins{
						sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.fundsInTDaiPool)),
					},
					*bankKeeper,
				)
				require.NoError(t, err)
			}

			subaccount := createNSubaccount(subaccountsKeeper, ctx, 1, big.NewInt(1_000))[0]
			subaccount.PerpetualPositions = tc.perpetualPositions
			subaccount.AssetPositions = tc.assetPositions
			subaccountYieldIndex := constants.AssetYieldIndex_Zero
			if tc.subaccountAssetYieldIndex != "" {
				subaccountYieldIndex = tc.subaccountAssetYieldIndex
			}
			subaccount.AssetYieldIndex = subaccountYieldIndex
			subaccountsKeeper.SetSubaccount(ctx, subaccount)
			subaccountId := *subaccount.Id

			msgServer := keeper.NewMsgServerImpl(*subaccountsKeeper)

			resp, err := msgServer.ClaimYieldForSubaccount(ctx, &tc.msgClaimYieldForSubaccount)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.Empty(t, resp)
			}
			newSubaccount := subaccountsKeeper.GetSubaccount(ctx, subaccountId)
			require.Equal(t, len(newSubaccount.PerpetualPositions), len(tc.expectedPerpetualPositions))
			for i, ep := range tc.expectedPerpetualPositions {
				require.Equal(t, *ep, *newSubaccount.PerpetualPositions[i])
			}
			require.Equal(t, len(newSubaccount.AssetPositions), len(tc.expectedAssetPositions))
			for i, ep := range tc.expectedAssetPositions {
				require.Equal(t, *ep, *newSubaccount.AssetPositions[i])
			}
			if tc.expectedErr == nil {
				require.Equal(t, 0, globalAssetYieldIndex.Cmp(ratelimitkeeper.ConvertStringToBigRatWithPanicOnErr(newSubaccount.AssetYieldIndex)),
					"Expected AssetYieldIndex %v. Got %v.", globalAssetYieldIndex, newSubaccount.AssetYieldIndex,
				)
			}

			for collateralPoolAddr, expectedTDaiBal := range tc.expectedCollateralPoolTDaiBalances {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					sdk.MustAccAddressFromBech32(collateralPoolAddr),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewInt(expectedTDaiBal)),
					TDaiBal,
				)
			}

			if tc.expectedTDaiYieldPoolBalance != nil {
				TDaiBal := bankKeeper.GetBalance(
					ctx,
					accountKeeper.GetModuleAddress(ratelimittypes.TDaiPoolAccount),
					asstypes.AssetTDai.Denom,
				)
				require.Equal(t,
					sdk.NewCoin(asstypes.AssetTDai.Denom, sdkmath.NewIntFromBigInt(tc.expectedTDaiYieldPoolBalance)),
					TDaiBal,
				)

			}
		})
	}
}
