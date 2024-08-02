//go:build all || integration_test

package cli_test

import (
	"fmt"
	"math"

	"math/big"
	"testing"

	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"

	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
	testutil "github.com/dydxprotocol/v4-chain/protocol/testutil/util"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clob_testutil "github.com/dydxprotocol/v4-chain/protocol/x/clob/client/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sa_testutil "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/client/testutil"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	appconstants "github.com/dydxprotocol/v4-chain/protocol/app/constants"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	testutil_bank "github.com/dydxprotocol/v4-chain/protocol/testutil/bank"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	"github.com/stretchr/testify/suite"
)

var (
	liqTestMakerOrderQuantums                = satypes.BaseQuantums(100_000_000) // 1 BTC.
	liqTestInitialSubaccountModuleAccBalance = int64(
		10_000 * constants.QuoteBalance_OneDollar, // $10,000.
	)
	liqTestSubaccountNumberZero = uint32(0)
	liqTestSubaccountNumberOne  = uint32(1)
	liqTestUnixSocketAddress    = "/tmp/liquidations_cli_test.sock"
)

type LiquidationsIntegrationTestSuite struct {
	suite.Suite

	validatorAddress sdk.AccAddress
	cfg              network.Config
	network          *network.Network
}

func TestLiquidationOrderIntegrationTestSuite(t *testing.T) {
	// Deterministic Mnemonic.
	validatorMnemonic := constants.AliceMnenomic

	// Generated from the above Mnemonic.
	validatorAddress := constants.AliceAccAddress

	appOptions := appoptions.NewFakeAppOptions()

	// Configure test network.
	cfg := network.DefaultConfig(&network.NetworkConfigOptions{
		AppOptions: appOptions,
		OnNewApp: func(val networktestutil.ValidatorI) {
			testval, ok := val.(networktestutil.Validator)
			if !ok {
				panic("incorrect validator type")
			}

			// Disable the Bridge and Price daemons in the integration tests.
			appOptions.Set(daemonflags.FlagPriceDaemonEnabled, false)
			appOptions.Set(daemonflags.FlagBridgeDaemonEnabled, false)
			appOptions.Set(daemonflags.FlagOracleEnabled, false)

			// Effectively disable the health monitor panic timeout for these tests. This is necessary
			// because all clob cli tests are running in the same process and the total time to run is >> 5 minutes
			// on CI, causing the panic to trigger for liquidations daemon go routines that haven't been properly
			// cleaned up after a test run.
			// TODO(CORE-29): Remove this once the liquidations daemon is refactored to be stoppable.
			appOptions.Set(daemonflags.FlagMaxDaemonUnhealthySeconds, math.MaxUint32)

			// Make sure the daemon is using the correct GRPC address.
			appOptions.Set(appflags.GrpcAddress, testval.AppConfig.GRPC.Address)

			// Enable the liquidations daemon in the integration tests.
			appOptions.Set(daemonflags.FlagUnixSocketAddress, liqTestUnixSocketAddress)
		},
	})

	cfg.Mnemonics = append(cfg.Mnemonics, validatorMnemonic)
	cfg.ChainID = appconstants.AppName

	suite.Run(t, NewLiquidationsIntegrationTestSuite(cfg, validatorAddress))
}

func NewLiquidationsIntegrationTestSuite(
	cfg network.Config,
	validatorAddress sdk.AccAddress,
) *LiquidationsIntegrationTestSuite {
	return &LiquidationsIntegrationTestSuite{cfg: cfg, validatorAddress: validatorAddress}
}

func (s *LiquidationsIntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up liquidations integration test suite")

	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	clobPair := constants.ClobPair_Btc
	state := types.GenesisState{}

	state.ClobPairs = append(state.ClobPairs, clobPair)
	state.LiquidationsConfig = types.LiquidationsConfig_Default

	perpstate := perptypes.GenesisState{}
	perpstate.LiquidityTiers = constants.LiquidityTiers
	perpstate.Params = constants.PerpetualsGenesisParams
	perpetual := constants.BtcUsd_20PercentInitial_10PercentMaintenance
	perpstate.Perpetuals = append(perpstate.Perpetuals, perpetual)

	pricesstate := constants.Prices_DefaultGenesisState

	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf

	// Set the balances in the genesis state.
	s.cfg.GenesisState[banktypes.ModuleName] = clob_testutil.CreateBankGenesisState(
		s.T(),
		s.cfg,
		liqTestInitialSubaccountModuleAccBalance,
	)

	sastate := satypes.GenesisState{}
	sastate.Subaccounts = append(
		sastate.Subaccounts,
		satypes.Subaccount{
			Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: liqTestSubaccountNumberZero},
			AssetPositions: []*satypes.AssetPosition{
				&constants.Usdc_Asset_100_000,
			},
			PerpetualPositions: []*satypes.PerpetualPosition{},
		},
		satypes.Subaccount{
			Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: liqTestSubaccountNumberOne},
			AssetPositions: []*satypes.AssetPosition{
				testutil.CreateSingleAssetPosition(
					assettypes.AssetUsdc.Id,
					big.NewInt(-45_001_000_000), // -$45,001
				),
			},
			PerpetualPositions: []*satypes.PerpetualPosition{
				testutil.CreateSinglePerpetualPosition(
					0,
					big.NewInt(100_000_000), // 1 BTC
					big.NewInt(0),
					big.NewInt(0),
				),
			},
		},
	)

	sabuf, err := s.cfg.Codec.MarshalJSON(&sastate)
	s.Require().NoError(err)
	s.cfg.GenesisState[satypes.ModuleName] = sabuf

	// Ensure that no funding payments will occur during this test.
	epstate := constants.GenerateEpochGenesisStateWithoutFunding()

	feeTiersState := feetierstypes.GenesisState{}
	feeTiersState.Params = constants.PerpetualFeeParams

	epbuf, err := s.cfg.Codec.MarshalJSON(&epstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	perpbuf, err := s.cfg.Codec.MarshalJSON(&perpstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[perptypes.ModuleName] = perpbuf

	pricesbuf, err := s.cfg.Codec.MarshalJSON(&pricesstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[pricestypes.ModuleName] = pricesbuf

	feeTiersBuf, err := s.cfg.Codec.MarshalJSON(&feeTiersState)
	s.Require().NoError(err)
	s.cfg.GenesisState[feetierstypes.ModuleName] = feeTiersBuf

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TestCLILiquidations creates two subaccounts (where one is undercollateralized), and places a
// maker order from the well-collateralized subaccount that should match with the liquidation order.
// The account which is liquidated and places the maker orders is also the validator's AccAddress.
// After the matching, the subaccounts are queried and assertions are performed on their
// QuoteBalance and PerpetualPositions, along with the balances of the Subaccounts module,
// Distribution module, and insurance fund.
func (s *LiquidationsIntegrationTestSuite) TestCLILiquidations() {
	val := s.network.Validators[0]
	ctx := val.ClientCtx

	currentHeight, err := s.network.LatestHeight()
	s.Require().NoError(err)

	goodTilBlock := uint32(currentHeight) + types.ShortBlockWindow
	clientId := uint64(1)
	subticks := types.Subticks(50_000_000_000)

	// Place the maker order that should be filled by the liquidation order.
	_, err = clob_testutil.MsgPlaceOrderExec(
		ctx,
		s.validatorAddress,
		liqTestSubaccountNumberZero,
		clientId,
		constants.ClobPair_Btc.Id,
		types.Order_SIDE_BUY,
		liqTestMakerOrderQuantums,
		subticks.ToUint64(),
		goodTilBlock,
	)
	s.Require().NoError(err)

	currentHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)

	// Wait for a few blocks to ensure the liquidation order was placed, matched, and included
	// in a block.
	_, err = s.network.WaitForHeight(currentHeight + 3)
	s.Require().NoError(err)

	// Query both subaccounts.
	resp, err := sa_testutil.MsgQuerySubaccountExec(ctx, s.validatorAddress, liqTestSubaccountNumberZero)
	s.Require().NoError(err)

	var subaccountResp satypes.QuerySubaccountResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
	subaccountZero := subaccountResp.Subaccount

	resp, err = sa_testutil.MsgQuerySubaccountExec(ctx, s.validatorAddress, liqTestSubaccountNumberOne)
	s.Require().NoError(err)

	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
	subaccountOne := subaccountResp.Subaccount

	// Compute the fill price so as to know how much QuoteBalance should be remaining.
	fillSizeQuoteQuantums := types.FillAmountToQuoteQuantums(
		subticks,
		liqTestMakerOrderQuantums,
		constants.ClobPair_Btc.QuantumConversionExponent,
	).Int64()

	// Assert that both Subaccounts have the appropriate state.
	takerFee := fillSizeQuoteQuantums * int64(constants.PerpetualFeeParams.Tiers[0].TakerFeePpm) / int64(lib.OneMillion)
	makerFee := fillSizeQuoteQuantums * int64(constants.PerpetualFeeParams.Tiers[0].MakerFeePpm) / int64(lib.OneMillion)
	subaccountZeroInitialQuoteBalance := constants.Usdc_Asset_100_000.GetBigQuantums().Int64()
	s.Require().Contains(
		[]*big.Int{
			new(big.Int).SetInt64(subaccountZeroInitialQuoteBalance - fillSizeQuoteQuantums - takerFee),
			new(big.Int).SetInt64(subaccountZeroInitialQuoteBalance - fillSizeQuoteQuantums - makerFee),
		},
		subaccountZero.GetUsdcPosition(),
	)
	s.Require().Len(subaccountZero.PerpetualPositions, 1)
	s.Require().Equal(liqTestMakerOrderQuantums.ToBigInt(), subaccountZero.PerpetualPositions[0].GetBigQuantums())

	subaccountOneInitialQuoteBalance := int64(-45_001_000_000)
	liquidationFee := fillSizeQuoteQuantums *
		int64(types.LiquidationsConfig_Default.MaxLiquidationFeePpm) /
		int64(lib.OneMillion)
	s.Require().Equal(
		new(big.Int).SetInt64(subaccountOneInitialQuoteBalance+fillSizeQuoteQuantums-liquidationFee),
		subaccountOne.GetUsdcPosition(),
	)
	s.Require().Empty(subaccountOne.PerpetualPositions)

	// Check that the `subaccounts` module account has expected remaining USDC balance.
	saModuleUSDCBalance, err := testutil_bank.GetModuleAccUsdcBalance(
		val,
		s.network.Config.Codec,
		satypes.ModuleName,
	)
	s.Require().NoError(err)
	s.Require().Equal(
		initialSubaccountModuleAccBalance-makerFee-liquidationFee,
		saModuleUSDCBalance,
	)

	// Check that the insurance fund has expected USDC balance.
	insuranceFundBalance, err := testutil_bank.GetModuleAccUsdcBalance(
		val,
		s.network.Config.Codec,
		perptypes.InsuranceFundName,
	)

	s.Require().NoError(err)
	s.Require().Equal(liquidationFee, insuranceFundBalance)

	// Check that the `distribution` module account has expected remaining USDC balance.
	// During `BeginBlock()`, the `fee-collector` module account will send all fees
	// to the `distribution` module account, and the fees will stay in `distribution`
	// until withdrawn. More details at:
	// https://docs.cosmos.network/v0.45/modules/distribution/03_begin_block.html#the-distribution-scheme
	distrModuleUSDCBalance, err := testutil_bank.GetModuleAccUsdcBalance(
		val,
		s.network.Config.Codec,
		distrtypes.ModuleName,
	)

	s.Require().NoError(err)
	s.Require().Equal(makerFee, distrModuleUSDCBalance)
}
