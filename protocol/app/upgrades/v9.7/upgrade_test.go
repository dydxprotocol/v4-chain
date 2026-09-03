package v_9_7_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v_9_7 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v9.7"
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	keepertest "github.com/dydxprotocol/v4-chain/protocol/testutil/keeper"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/memclob"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Only the isolated market in final settlement is swept; the bank mock fails the test on any
// call for the cross final-settlement market or the active isolated market.
func TestSweepIsolatedInsuranceFunds(t *testing.T) {
	memClob := memclob.NewMemClobPriceTimePriority(false)
	bankMock := &mocks.BankKeeper{}
	mockIndexerEventManager := &mocks.IndexerEventManager{}
	mockIndexerEventManager.On("AddTxnEvent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return()
	ks := keepertest.NewClobKeepersTestContext(t, memClob, bankMock, mockIndexerEventManager)

	require.NoError(t, keepertest.CreateUsdcAsset(ks.Ctx, ks.AssetsKeeper))
	keepertest.CreateTestMarkets(t, ks.Ctx, ks.PricesKeeper)
	keepertest.CreateTestLiquidityTiers(t, ks.Ctx, ks.PerpetualsKeeper)
	keepertest.CreateTestPerpetuals(t, ks.Ctx, ks.PerpetualsKeeper)
	keepertest.CreateTestClobPairs(t, ks.Ctx, ks.ClobKeeper, []clobtypes.ClobPair{
		constants.ClobPair_Btc_Final_Settlement,
		constants.ClobPair_3_Iso_Final_Settlement,
		constants.ClobPair_4_Iso2,
	})

	// The upgrade module's PreBlocker runs before the clob module's PreBlocker, so the in-memory
	// perpetual-to-clob-pair mapping is not hydrated when the upgrade handler executes.
	ks.ClobKeeper.PerpetualIdToClobPairId = make(map[uint32][]clobtypes.ClobPairId)

	isolatedInsuranceFundAddr, err := ks.PerpetualsKeeper.GetInsuranceFundModuleAddress(
		ks.Ctx,
		constants.IsoUsd_IsolatedMarket.Params.Id,
	)
	require.NoError(t, err)
	require.NotEqual(t, perptypes.InsuranceFundModuleAddress, isolatedInsuranceFundAddr)

	isolatedFundBalance := sdk.NewCoin(constants.Usdc.Denom, sdkmath.NewInt(1_500_000_000))
	bankMock.On(
		"GetBalance",
		mock.Anything,
		isolatedInsuranceFundAddr,
		constants.Usdc.Denom,
	).Return(isolatedFundBalance)
	bankMock.On(
		"SendCoins",
		mock.Anything,
		isolatedInsuranceFundAddr,
		perptypes.InsuranceFundModuleAddress,
		sdk.Coins{isolatedFundBalance},
	).Return(nil)

	require.NoError(t, v_9_7.SweepIsolatedInsuranceFunds(ks.Ctx, ks.ClobKeeper, ks.SubaccountsKeeper))

	bankMock.AssertExpectations(t)
}
