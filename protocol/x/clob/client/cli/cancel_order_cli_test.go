//go:build all || integration_test

package cli_test

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"testing"

	appconstants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/appoptions"
	testutil_bank "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	cli_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/client/testutil"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	epochstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs/types"
	feetierstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	ratelimitcli "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/client/cli"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	sa_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/client/testutil"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	"github.com/cosmos/cosmos-sdk/client"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/suite"
)

const (
	cancelsInitialQuoteBalance               = int64(1_000_000_000)  // $1,000.
	cancelsInitialSubaccountModuleAccBalance = int64(10_000_000_000) // $10,000.
	initialSDAIBalance                       = int64(10_000_000_000) // 100 sDAI
	cancelsSubaccountNumberZero              = uint32(0)
	cancelsSubaccountNumberOne               = uint32(1)
)

type CancelOrderIntegrationTestSuite struct {
	suite.Suite
	validatorAddress sdk.AccAddress
	cfg              network.Config
	network          *network.Network
}

func GetBalanceAfterYield(clientCtx client.Context, initialBalance *big.Int) (balance int64, err error) {

	args := []string{}
	data, err := clitestutil.ExecTestCLICmd(clientCtx, ratelimitcli.CmdGetSDAIPriceQuery(), args)
	if err != nil {
		return 0, err
	}

	var resp ratelimittypes.GetSDAIPriceQueryResponse
	err = json.Unmarshal(data.Bytes(), &resp)
	if err != nil {
		return 0, err
	}
	fmt.Println("SDAI: ", resp)

	priceFloat, success := new(big.Float).SetString(resp.Price)
	if !success {
		return 0, fmt.Errorf("failed to parse price as big.Float")
	}

	precision := new(big.Int).Exp(
		big.NewInt(ratelimittypes.BASE_10),
		big.NewInt(ratelimittypes.SDAI_DECIMALS),
		nil,
	)

	// Convert initialBalance to big.Float
	initialBalanceFloat := new(big.Float).SetInt(initialBalance)

	// Multiply price by initialBalance
	result := new(big.Float).Mul(priceFloat, initialBalanceFloat)

	// Convert precision to big.Float
	precisionFloat := new(big.Float).SetInt(precision)

	// Divide the result by precision
	result.Quo(result, precisionFloat)

	// Convert result to int64
	balanceFloat, _ := result.Float64()
	balance = int64(balanceFloat)

	return balance, nil
}

func TestCancelOrderIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &CancelOrderIntegrationTestSuite{})
}

func (s *CancelOrderIntegrationTestSuite) SetupTest() {
	s.T().Log("setting up cancel order integration test")

	// // Deterministic Mnemonic.
	validatorMnemonic := constants.AliceMnenomic

	// Generated from the above Mnemonic.
	s.validatorAddress = constants.AliceAccAddress

	// // Configure test network.
	appOptions := appoptions.NewFakeAppOptions()
	s.cfg = network.DefaultConfig(&network.NetworkConfigOptions{
		AppOptions: appOptions,
		OnNewApp: func(val networktestutil.ValidatorI) {
			testval, ok := val.(networktestutil.Validator)
			if !ok {
				panic("incorrect validator type")
			}

			// Disable the Price daemon in the integration tests.
			appOptions.Set(daemonflags.FlagPriceDaemonEnabled, false)
			appOptions.Set(daemonflags.FlagSDAIDaemonMockEnabled, true)

			// Effectively disable the health monitor panic timeout for these tests. This is necessary
			// because all clob cli tests are running in the same process and the total time to run is >> 5 minutes
			// on CI, causing the panic to trigger for liquidations daemon go routines that haven't been properly
			// cleaned up after a test run.
			// TODO(CORE-29): Remove this once the liquidations daemon is refactored to be stoppable.
			appOptions.Set(daemonflags.FlagMaxDaemonUnhealthySeconds, math.MaxUint32)

			// Make sure the daemon is using the correct GRPC address.
			appOptions.Set(appflags.GrpcAddress, testval.AppConfig.GRPC.Address)
		},
	})

	s.cfg.Mnemonics = append(s.cfg.Mnemonics, validatorMnemonic)
	s.cfg.ChainID = appconstants.AppName

	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	clobPair := constants.ClobPair_Btc
	state := types.GenesisState{}

	state.ClobPairs = append(state.ClobPairs, clobPair)
	state.LiquidationsConfig = types.LiquidationsConfig_NoFee

	perpstate := perptypes.GenesisState{}
	perpstate.LiquidityTiers = constants.LiquidityTiers
	perpstate.Params = constants.PerpetualsGenesisParams
	perpetual := constants.BtcUsd_50PercentInitial_40PercentMaintenance
	perpstate.Perpetuals = append(perpstate.Perpetuals, perpetual)

	pricesstate := constants.Prices_DefaultGenesisState

	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf

	s.cfg.GenesisState[banktypes.ModuleName] = cli_testutil.CreateBankGenesisState(
		s.T(),
		s.cfg,
	)

	sastate := satypes.GenesisState{}
	sastate.Subaccounts = append(
		sastate.Subaccounts,
		satypes.Subaccount{
			Id: &satypes.SubaccountId{
				Owner:  s.validatorAddress.String(),
				Number: 0,
			},
			AssetPositions: []*satypes.AssetPosition{
				{
					AssetId:  0,
					Quantums: dtypes.NewInt(1000000000),
				},
			},
			PerpetualPositions: []*satypes.PerpetualPosition{},
		},
		satypes.Subaccount{
			Id: &satypes.SubaccountId{
				Owner:  s.validatorAddress.String(),
				Number: 1,
			},
			AssetPositions: []*satypes.AssetPosition{
				{
					AssetId:  0,
					Quantums: dtypes.NewInt(1000000000),
				},
			},
			PerpetualPositions: []*satypes.PerpetualPosition{},
		},
	)

	epstate := constants.GenerateEpochGenesisStateWithoutFunding()

	feeTiersState := feetierstypes.GenesisState{}
	feeTiersState.Params = constants.PerpetualFeeParams

	epbuf, err := s.cfg.Codec.MarshalJSON(&epstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	sabuf, err := s.cfg.Codec.MarshalJSON(&sastate)
	s.Require().NoError(err)
	s.cfg.GenesisState[satypes.ModuleName] = sabuf

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

// TestCLICancelPendingOrder places then cancels an order from a subaccount, and then places an order from
// a different subaccount (with the same owner and different numbers).
// The orders placed are expected to match, but should not due to the first order being canceled.
// Afterwards, an additional cancel of an unknown is order is made (expected to be a no-op).
// The subaccounts are then queried and assertions are performed on their QuoteBalance and PerpetualPositions.
// The account which places the orders is also the validator's AccAddress.
func (s *CancelOrderIntegrationTestSuite) TestCLICancelPendingOrder() {

	val := s.network.Validators[0]
	ctx := val.ClientCtx

	currentHeight, err := s.network.LatestHeight()
	s.Require().NoError(err)

	goodTilBlock := uint32(currentHeight) + types.ShortBlockWindow
	clientId := uint64(1)
	quantums := satypes.BaseQuantums(1_000)
	subticks := types.Subticks(50_000_000_000)

	// Place the first order.
	_, err = cli_testutil.MsgPlaceOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
		clientId,
		constants.ClobPair_Btc.Id,
		types.Order_SIDE_BUY,
		quantums,
		subticks.ToUint64(),
		goodTilBlock,
	)
	s.Require().NoError(err)

	// Cancel the first order.
	_, err = cli_testutil.MsgCancelOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
		clientId,
		constants.ClobPair_Btc.Id,
		goodTilBlock,
	)
	s.Require().NoError(err)

	// Place the second order.
	_, err = cli_testutil.MsgPlaceOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberOne,
		clientId,
		constants.ClobPair_Btc.Id,
		types.Order_SIDE_SELL,
		quantums,
		subticks.ToUint64(),
		goodTilBlock,
	)
	s.Require().NoError(err)

	// Cancel an unknown order.
	unknownClientId := uint64(10)
	_, err = cli_testutil.MsgCancelOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
		unknownClientId,
		constants.ClobPair_Btc.Id,
		goodTilBlock,
	)
	s.Require().NoError(err)

	currentHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)

	// Wait for a few blocks.
	_, err = s.network.WaitForHeight(currentHeight + 3)
	s.Require().NoError(err)

	// Check that subaccounts balance have not changed, and no positions were opened.
	for _, subaccountNumber := range []uint32{cancelsSubaccountNumberZero, cancelsSubaccountNumberOne} {
		resp, err := sa_testutil.MsgQuerySubaccountExec(
			ctx,
			s.validatorAddress,
			subaccountNumber,
		)
		s.Require().NoError(err)

		var subaccountResp satypes.QuerySubaccountResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
		subaccount := subaccountResp.Subaccount

		s.Require().Equal(
			new(big.Int).SetInt64(cancelsInitialQuoteBalance),
			subaccount.GetTDaiPosition(),
		)
		s.Require().Len(subaccount.PerpetualPositions, 0)

		s.Require().Equal(
			new(big.Int).SetInt64(cancelsInitialQuoteBalance),
			subaccount.GetTDaiPosition())
		s.Require().Len(subaccount.PerpetualPositions, 0)
	}

	// Check that the `subaccounts` module account balance has not changed.
	saModuleTDaiBalance, err := testutil_bank.GetModuleAccTDaiBalance(
		val,
		s.network.Config.Codec,
		satypes.ModuleName,
	)
	s.Require().NoError(err)
	s.Require().Equal(
		cancelsInitialSubaccountModuleAccBalance,
		saModuleTDaiBalance,
	)

	distrModuleTDaiBalance, err := testutil_bank.GetModuleAccTDaiBalance(
		val,
		s.network.Config.Codec,
		distrtypes.ModuleName,
	)

	s.Require().NoError(err)
	s.Require().Equal(int64(0), distrModuleTDaiBalance)

}

// TestCLICancelMatchingOrders places two matching orders from two different subaccounts (with the
// same owner and different numbers), then cancels the first matching order a few blocks later.
// The matching orders should not be canceled.
// The subaccounts are then queried and assertions are performed on their QuoteBalance and PerpetualPositions.
// The account which places the orders is also the validator's AccAddress.
func (s *CancelOrderIntegrationTestSuite) TestCLICancelMatchingOrders() {

	val := s.network.Validators[0]
	ctx := val.ClientCtx

	currentHeight, err := s.network.LatestHeight()
	s.Require().NoError(err)

	goodTilBlock := uint32(currentHeight) + types.ShortBlockWindow
	clientId := uint64(2)
	quantums := satypes.BaseQuantums(1_000)
	subticks := types.Subticks(50_000_000_000)

	// Place the first order.
	_, err = cli_testutil.MsgPlaceOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
		clientId,
		constants.ClobPair_Btc.Id,
		types.Order_SIDE_BUY,
		quantums,
		subticks.ToUint64(),
		goodTilBlock,
	)
	s.Require().NoError(err)

	// Place the second order.
	_, err = cli_testutil.MsgPlaceOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberOne,
		clientId,
		constants.ClobPair_Btc.Id,
		types.Order_SIDE_SELL,
		quantums,
		subticks.ToUint64(),
		goodTilBlock,
	)
	s.Require().NoError(err)

	currentHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)

	// Wait for a few blocks.
	_, err = s.network.WaitForHeight(currentHeight + 3)
	s.Require().NoError(err)

	// Cancel the first order.
	_, err = cli_testutil.MsgCancelOrderExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
		clientId,
		constants.ClobPair_Btc.Id,
		goodTilBlock,
	)
	s.Require().NoError(err)

	currentHeight, err = s.network.LatestHeight()
	s.Require().NoError(err)

	// Wait for a few blocks.
	_, err = s.network.WaitForHeight(currentHeight + 3)
	s.Require().NoError(err)

	// Query both subaccounts.
	accResp, accErr := sa_testutil.MsgQuerySubaccountExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberZero,
	)
	s.Require().NoError(accErr)

	var subaccountResp satypes.QuerySubaccountResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(accResp.Bytes(), &subaccountResp))
	subaccountZero := subaccountResp.Subaccount

	accResp, accErr = sa_testutil.MsgQuerySubaccountExec(
		ctx,
		s.validatorAddress,
		cancelsSubaccountNumberOne,
	)
	s.Require().NoError(accErr)

	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(accResp.Bytes(), &subaccountResp))
	subaccountOne := subaccountResp.Subaccount

	// Compute the fill price so as to know how much QuoteBalance should be remaining.
	fillSizeQuoteQuantums := types.FillAmountToQuoteQuantums(
		subticks,
		quantums,
		constants.ClobPair_Btc.QuantumConversionExponent,
	).Int64()

	// cancelsInitialQuoteBalanceAfterYield, err := GetBalanceAfterYield(ctx, new(big.Int).SetInt64(cancelsInitialQuoteBalance))
	// s.Require().NoError(err)

	// Assert that both Subaccounts have the appropriate state.
	// Order could be maker or taker after Uncross, so assert that account could have been either.
	takerFee := fillSizeQuoteQuantums *
		int64(constants.PerpetualFeeParams.Tiers[0].TakerFeePpm) /
		int64(lib.OneMillion)
	makerFee := fillSizeQuoteQuantums *
		int64(constants.PerpetualFeeParams.Tiers[0].MakerFeePpm) /
		int64(lib.OneMillion)

	fmt.Println(subaccountZero.GetTDaiPosition())
	fmt.Println(subaccountOne.GetTDaiPosition())
	fmt.Println(cancelsInitialQuoteBalance)
	fmt.Println(fillSizeQuoteQuantums)
	fmt.Println(takerFee)
	fmt.Println(makerFee)

	s.Require().Contains(
		[]*big.Int{
			new(big.Int).SetInt64(cancelsInitialQuoteBalance - fillSizeQuoteQuantums - takerFee),
			new(big.Int).SetInt64(cancelsInitialQuoteBalance - fillSizeQuoteQuantums - makerFee),
		},
		subaccountZero.GetTDaiPosition(),
	)
	s.Require().Len(subaccountZero.PerpetualPositions, 1)
	s.Require().Equal(quantums.ToBigInt(), subaccountZero.PerpetualPositions[0].GetBigQuantums())

	s.Require().Contains(
		[]*big.Int{
			new(big.Int).SetInt64(cancelsInitialQuoteBalance + fillSizeQuoteQuantums - takerFee),
			new(big.Int).SetInt64(cancelsInitialQuoteBalance + fillSizeQuoteQuantums - makerFee),
		},
		subaccountOne.GetTDaiPosition(),
	)
	s.Require().Len(subaccountOne.PerpetualPositions, 1)
	s.Require().Equal(new(big.Int).Neg(
		quantums.ToBigInt()),
		subaccountOne.PerpetualPositions[0].GetBigQuantums(),
	)

	// Check that the `subaccounts` module account has expected remaining TDAI balance.
	saModuleTDaiBalance, err := testutil_bank.GetModuleAccTDaiBalance(
		val,
		s.network.Config.Codec,
		satypes.ModuleName,
	)
	s.Require().NoError(err)
	s.Require().Equal(
		int64(cancelsInitialSubaccountModuleAccBalance-takerFee-makerFee),
		saModuleTDaiBalance,
	)

	distrModuleTDaiBalance, err := testutil_bank.GetModuleAccTDaiBalance(
		val,
		s.network.Config.Codec,
		distrtypes.ModuleName,
	)

	s.Require().NoError(err)
	s.Require().Equal(makerFee+takerFee, distrModuleTDaiBalance)

	// test the sdai - we comment this out because we test it above as we include yield in the subaccount calc

	// cfg := network.DefaultConfig(nil)

	// chi := "1006681181716810314385961731"

	// time.Sleep(15 * time.Second)

	// rateQuery := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" query ratelimit get-sdai-price "
	// data, _, err = network.QueryCustomNetwork(rateQuery)

	// require.NoError(s.T(), err)
	// var respSdai ratelimittypes.GetSDAIPriceQueryResponse
	// require.NoError(s.T(), cfg.Codec.UnmarshalJSON(data, &respSdai))

	// chiFloat, success := new(big.Float).SetString(chi)
	// require.True(s.T(), success, "Failed to parse chi as big.Float")

	// priceFloat, success := new(big.Float).SetString(respSdai.Price)
	// require.True(s.T(), success, "Failed to parse price as big.Float")

	// // Compare the big.Float values directly
	// comparison := new(big.Float).Quo(priceFloat, chiFloat)

	// minThreshold := big.NewFloat(0.99)
	// maxThreshold := big.NewFloat(1.16)

	// require.True(s.T(), comparison.Cmp(minThreshold) >= 0, "Price should be at least 99% of chi")
	// require.True(s.T(), comparison.Cmp(maxThreshold) <= 0, "Price should be at most 116% of chi")

	// network.CleanupCustomNetwork()
}
