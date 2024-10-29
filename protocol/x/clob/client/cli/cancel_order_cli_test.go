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
	"github.com/cosmos/cosmos-sdk/client/flags"
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

	// rateQuery := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" query ratelimit get-sdai-price "
	// data, _, err := network.QueryCustomNetwork(rateQuery)

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, "node0"),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	}
	data, err := clitestutil.ExecTestCLICmd(clientCtx, ratelimitcli.CmdGetSDAIPriceQuery(), args)
	if err != nil {
		return 0, err
	}

	var resp ratelimittypes.GetSDAIPriceQueryResponse
	err = json.Unmarshal(data.Bytes(), &resp)
	if err != nil {
		return 0, err
	}

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

	epbuf, err := s.network.Config.Codec.MarshalJSON(&epstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	sabuf, err := s.network.Config.Codec.MarshalJSON(&sastate)
	s.Require().NoError(err)
	s.cfg.GenesisState[satypes.ModuleName] = sabuf

	perpbuf, err := s.network.Config.Codec.MarshalJSON(&perpstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[perptypes.ModuleName] = perpbuf

	pricesbuf, err := s.network.Config.Codec.MarshalJSON(&pricesstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[pricestypes.ModuleName] = pricesbuf

	feeTiersBuf, err := s.network.Config.Codec.MarshalJSON(&feeTiersState)
	s.Require().NoError(err)
	s.cfg.GenesisState[feetierstypes.ModuleName] = feeTiersBuf

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// fullGenesis := "\".app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"-8\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}] | .app_state.clob.liquidations_config = {\\\"insurance_fund_fee_ppm\\\": \\\"5000\\\", \\\"validator_fee_ppm\\\": \\\"0\\\", \\\"liquidity_fee_ppm\\\": \\\"0\\\", \\\"max_cumulative_insurance_fund_delta\\\": \\\"1000000000000\\\", \\\"fillable_price_config\\\": {\\\"bankruptcy_adjustment_ppm\\\": \\\"1000000\\\", \\\"spread_to_maintenance_margin_ratio_ppm\\\": \\\"100000\\\"}} | .app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"1\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"750000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"2\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"3\\\", \\\"name\\\": \\\"3\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"4\\\", \\\"name\\\": \\\"4\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"800000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"5\\\", \\\"name\\\": \\\"5\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"6\\\", \\\"name\\\": \\\"6\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"900000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"7\\\", \\\"name\\\": \\\"7\\\", \\\"initial_margin_ppm\\\": \\\"0\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"8\\\", \\\"name\\\": \\\"8\\\", \\\"initial_margin_ppm\\\": \\\"9910\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"50454000000\\\"}, {\\\"id\\\": \\\"101\\\", \\\"name\\\": \\\"101\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}] | .app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"} | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"-8\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"4\\\", \\\"market_id\\\": \\\"0\\\", \\\"market_type\\\": \\\"PERPETUAL_MARKET_TYPE_CROSS\\\", \\\"ticker\\\": \\\"BTC-USD 50/40 margin requirements\\\"}, \\\"funding_index\\\": \\\"0\\\", \\\"yield_index\\\": \\\"0/1\\\", \\\"last_funding_rate\\\": \\\"0\\\"}] | .app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"spot_price\\\": \\\"5000000000\\\", \\\"pnl_price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"spot_price\\\": \\\"3000000000\\\", \\\"pnl_price\\\": \\\"3000000000\\\"}] | .app_state.bank.balances = [{\\\"address\\\": \\\"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6\\\", \\\"coins\\\": [{\\\"denom\\\": \\\"utdai\\\", \\\"amount\\\": \\\"10000000000\\\"}]}, {\\\"address\\\": \\\"dydx1r3fsd6humm0ghyq0te5jf8eumklmclya37zle0\\\", \\\"coins\\\": [{\\\"denom\\\": \\\"ibc/DEEFE2DEFDC8EA8879923C4CCA42BB888C3CD03FF7ECFEFB1C2FEC27A732ACC8\\\", \\\"amount\\\": \\\"10000000000000000000000\\\"}]}] | .app_state.bank.supply = [] | .app_state.subaccounts.subaccounts = [{\\\"id\\\": {\\\"number\\\": \\\"0\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"1000000000\\\"}], \\\"perpetual_positions\\\": []}, {\\\"id\\\": {\\\"number\\\": \\\"1\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"1000000000\\\"}], \\\"perpetual_positions\\\": []}] | .app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}] | .app_state.feetiers.params = {\\\"tiers\\\": [{\\\"name\\\": \\\"1\\\", \\\"maker_fee_ppm\\\": \\\"-200\\\", \\\"taker_fee_ppm\\\": \\\"500\\\"}]}\" \"--sDAI-daemon-mock-enabled=true\""
	// network.DeployCustomNetwork(fullGenesis)
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

	// goodTilBlock := uint32(0)
	// query := "docker exec interchain-security-instance interchain-security-cd query block --type=height 0"
	// data, _, _ := network.QueryCustomNetwork(query)
	// var resp blocktypes.Block
	// require.NoError(s.T(), s.cfg.Codec.UnmarshalJSON(data, &resp))
	// blockHeight := resp.LastCommit.Height

	// goodTilBlock = uint32(blockHeight) + types.ShortBlockWindow
	// goodTilBlockStr := strconv.Itoa(int(goodTilBlock))

	// buyTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 0 1 0 1 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err := network.QueryCustomNetwork(buyTx)
	// s.Require().NoError(err)

	// cancelTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob cancel-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 0 1 0 " +
	// 	goodTilBlockStr + " --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(cancelTx)
	// s.Require().NoError(err)

	// sellTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 1 1 0 2 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(sellTx)
	// s.Require().NoError(err)

	// cancelUknownTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob cancel-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 0 10 0 " +
	// 	goodTilBlockStr + " --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(cancelUknownTx)
	// s.Require().NoError(err)

	// time.Sleep(5 * time.Second)

	// goodTilBlock := uint32(0)
	// query := "docker exec interchain-security-instance interchain-security-cd query block --type=height 0"
	// data, _, _ := network.QueryCustomNetwork(query)
	// var resp blocktypes.Block
	// require.NoError(s.T(), s.cfg.Codec.UnmarshalJSON(data, &resp))
	// blockHeight := resp.LastCommit.Height

	// goodTilBlock = uint32(blockHeight) + types.ShortBlockWindow
	// goodTilBlockStr := strconv.Itoa(int(goodTilBlock))

	// buyTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 0 1 0 1 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err := network.QueryCustomNetwork(buyTx)
	// s.Require().NoError(err)

	// cancelTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob cancel-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 0 1 0 " +
	// 	goodTilBlockStr + " --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(cancelTx)
	// s.Require().NoError(err)

	// sellTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 1 1 0 2 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(sellTx)
	// s.Require().NoError(err)

	// cancelUknownTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob cancel-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 0 10 0 " +
	// 	goodTilBlockStr + " --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, err = network.QueryCustomNetwork(cancelUknownTx)
	// s.Require().NoError(err)

	// time.Sleep(5 * time.Second)

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

	// network.CleanupCustomNetwork()
}

// TestCLICancelMatchingOrders places two matching orders from two different subaccounts (with the
// same owner and different numbers), then cancels the first matching order a few blocks later.
// The matching orders should not be canceled.
// The subaccounts are then queried and assertions are performed on their QuoteBalance and PerpetualPositions.
// The account which places the orders is also the validator's AccAddress.
func (s *CancelOrderIntegrationTestSuite) TestCLICancelMatchingOrders() {
	// goodTilBlock := uint32(0)

	// blockHeightQuery := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" query block --type=height 0"
	// data, _, err := network.QueryCustomNetwork(blockHeightQuery)
	// if err != nil {
	// 	s.T().Fatalf("Failed to get block height: %v", err)
	// }
	// var resp blocktypes.Block
	// require.NoError(s.T(), s.cfg.Codec.UnmarshalJSON(data, &resp))
	// blockHeight := resp.LastCommit.Height

	// goodTilBlock = uint32(blockHeight) + types.ShortBlockWindow
	// goodTilBlockStr := strconv.Itoa(int(goodTilBlock))
	// quantums := satypes.BaseQuantums(1_000)
	// subticks := types.Subticks(50_000_000_000)

	// placeBuyTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 0 2 0 1 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, buyerr := network.QueryCustomNetwork(placeBuyTx)
	// if buyerr != nil {
	// 	s.T().Fatalf("Failed to place order: %v", buyerr)
	// }
	// s.Require().NoError(buyerr)

	// placeSellTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 1 2 0 2 1000 50000000000 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, sellerr := network.QueryCustomNetwork(placeSellTx)
	// if sellerr != nil {
	// 	s.T().Fatalf("Failed to place order: %v", sellerr)
	// }
	// s.Require().NoError(sellerr)

	// time.Sleep(5 * time.Second)

	// cancelBuyTx := "docker exec interchain-security-instance interchain-security-cd" +
	// 	" tx clob cancel-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" 0 2 0 " + goodTilBlockStr +
	// 	" --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6" +
	// 	" --chain-id consu --home /consu/validatoralice --keyring-backend test -y"
	// _, _, cancelerr := network.QueryCustomNetwork(cancelBuyTx)

	// if cancelerr != nil {
	// 	s.T().Fatalf("Failed to cancel order: %v", cancelerr)
	// }
	// s.Require().NoError(cancelerr)

	// time.Sleep(5 * time.Second)

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

	cancelsInitialQuoteBalanceAfterYield, err := GetBalanceAfterYield(ctx, new(big.Int).SetInt64(cancelsInitialQuoteBalance))
	s.Require().NoError(err)

	// Assert that both Subaccounts have the appropriate state.
	// Order could be maker or taker after Uncross, so assert that account could have been either.
	takerFee := fillSizeQuoteQuantums *
		int64(constants.PerpetualFeeParamsMakerRebate.Tiers[0].TakerFeePpm) /
		int64(lib.OneMillion)
	makerFee := fillSizeQuoteQuantums *
		int64(constants.PerpetualFeeParamsMakerRebate.Tiers[0].MakerFeePpm) /
		int64(lib.OneMillion)

	s.Require().Contains(
		[]*big.Int{
			new(big.Int).SetInt64(cancelsInitialQuoteBalanceAfterYield - fillSizeQuoteQuantums - takerFee),
			new(big.Int).SetInt64(cancelsInitialQuoteBalanceAfterYield - fillSizeQuoteQuantums - makerFee),
		},
		subaccountZero.GetTDaiPosition(),
	)
	s.Require().Len(subaccountZero.PerpetualPositions, 1)
	s.Require().Equal(quantums.ToBigInt(), subaccountZero.PerpetualPositions[0].GetBigQuantums())

	s.Require().Contains(
		[]*big.Int{
			new(big.Int).SetInt64(cancelsInitialQuoteBalanceAfterYield + fillSizeQuoteQuantums - takerFee),
			new(big.Int).SetInt64(cancelsInitialQuoteBalanceAfterYield + fillSizeQuoteQuantums - makerFee),
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
		int64(10013362212),
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
