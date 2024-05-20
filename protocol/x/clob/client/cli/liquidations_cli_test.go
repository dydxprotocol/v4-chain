//go:build all || integration_test

package cli_test

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"math/big"
	"testing"

	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	blocktypes "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/lib"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	sa_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/client/testutil"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	appconstants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/appoptions"
	testutil_bank "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/bank"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/stretchr/testify/require"
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

			// Disable the Price daemon in the integration tests.
			appOptions.Set(daemonflags.FlagPriceDaemonEnabled, false)

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

	// clobPair := constants.ClobPair_Btc
	// state := types.GenesisState{}

	// state.ClobPairs = append(state.ClobPairs, clobPair)
	// state.LiquidationsConfig = types.LiquidationsConfig_Default

	// perpstate := perptypes.GenesisState{}
	// perpstate.LiquidityTiers = constants.LiquidityTiers
	// perpstate.Params = constants.PerpetualsGenesisParams
	// perpetual := constants.BtcUsd_20PercentInitial_10PercentMaintenance
	// perpstate.Perpetuals = append(perpstate.Perpetuals, perpetual)

	// pricesstate := constants.Prices_DefaultGenesisState

	// buf, err := s.cfg.Codec.MarshalJSON(&state)
	// s.NoError(err)
	// s.cfg.GenesisState[types.ModuleName] = buf

	// // Set the balances in the genesis state.
	// s.cfg.GenesisState[banktypes.ModuleName] = testutil.CreateBankGenesisState(
	// 	s.T(),
	// 	s.cfg,
	// 	liqTestInitialSubaccountModuleAccBalance,
	// )

	// sastate := satypes.GenesisState{}
	// sastate.Subaccounts = append(
	// 	sastate.Subaccounts,
	// 	satypes.Subaccount{
	// 		Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: liqTestSubaccountNumberZero},
	// 		AssetPositions: []*satypes.AssetPosition{
	// 			&constants.Usdc_Asset_100_000,
	// 		},
	// 		PerpetualPositions: []*satypes.PerpetualPosition{},
	// 	},
	// 	satypes.Subaccount{
	// 		Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: liqTestSubaccountNumberOne},
	// 		AssetPositions: []*satypes.AssetPosition{
	// 			{
	// 				AssetId:  assettypes.AssetUsdc.Id,
	// 				Quantums: dtypes.NewInt(-45_001_000_000), // -$45,001
	// 			},
	// 		},
	// 		PerpetualPositions: []*satypes.PerpetualPosition{
	// 			{
	// 				PerpetualId: 0,
	// 				Quantums:    dtypes.NewInt(100_000_000), // 1 BTC
	// 			},
	// 		},
	// 	},
	// )

	// subbaccountGenesis := "\".app_state.subaccounts.subaccounts = [{\\\"id\\\": {\\\"number\\\": \\\"0\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"100000000000\\\"}], \\\"perpetual_positions\\\": []}, {\\\"id\\\": {\\\"number\\\": \\\"1\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"-45001000000\\\"}], \\\"perpetual_positions\\\": [{\\\"perpetual_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"100000000\\\", \\\"funding_index\\\": \\\"0\\\"}]}]\" \"\""

	// sabuf, err := s.cfg.Codec.MarshalJSON(&sastate)
	// s.Require().NoError(err)
	// s.cfg.GenesisState[satypes.ModuleName] = sabuf

	// // Ensure that no funding payments will occur during this test.
	// epstate := constants.GenerateEpochGenesisStateWithoutFunding()

	// feeTiersState := feetierstypes.GenesisState{}
	// feeTiersState.Params = constants.PerpetualFeeParams

	// epbuf, err := s.cfg.Codec.MarshalJSON(&epstate)
	// s.Require().NoError(err)
	// s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	// perpbuf, err := s.cfg.Codec.MarshalJSON(&perpstate)
	// s.Require().NoError(err)
	// s.cfg.GenesisState[perptypes.ModuleName] = perpbuf

	// pricesbuf, err := s.cfg.Codec.MarshalJSON(&pricesstate)
	// s.Require().NoError(err)
	// s.cfg.GenesisState[pricestypes.ModuleName] = pricesbuf

	// feeTiersBuf, err := s.cfg.Codec.MarshalJSON(&feeTiersState)
	// s.Require().NoError(err)
	// s.cfg.GenesisState[feetierstypes.ModuleName] = feeTiersBuf
	// perpGenesis := "\".app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"1\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"750000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"2\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"3\\\", \\\"name\\\": \\\"3\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"4\\\", \\\"name\\\": \\\"4\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"800000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"5\\\", \\\"name\\\": \\\"5\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"6\\\", \\\"name\\\": \\\"6\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"900000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"7\\\", \\\"name\\\": \\\"7\\\", \\\"initial_margin_ppm\\\": \\\"0\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"8\\\", \\\"name\\\": \\\"8\\\", \\\"initial_margin_ppm\\\": \\\"9910\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"50454000000\\\"}, {\\\"id\\\": \\\"101\\\", \\\"name\\\": \\\"101\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}] | .app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"} | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"-8\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"3\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"BTC-USD 20/10 margin requirements\\\"}, \\\"funding_index\\\": \\\"0\\\"}]\" \"\""
	// clobGenesis := "\".app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"-8\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}]\""
	// priceGenesis := "\".app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}]\" \"\""
	// liquidationGenesis := "\".app_state.clob.liquidations_config = {\\\"max_liquidation_fee_ppm\\\": \\\"5000\\\", \\\"fillable_price_config\\\": {\\\"bankruptcy_adjustment_ppm\\\": \\\"1000000\\\", \\\"spread_to_maintenance_margin_ratio_ppm\\\": \\\"100000\\\"}, \\\"position_block_limits\\\": {\\\"max_position_portion_liquidated_ppm\\\": \\\"1000000\\\", \\\"min_position_notional_liquidated\\\": \\\"1000\\\"}, \\\"subaccount_block_limits\\\": {\\\"max_notional_liquidated\\\": \\\"100000000000000\\\", \\\"max_quantums_insurance_lost\\\": \\\"100000000000000\\\"}}\" \"\""
	// subbaccountGenesis := "\".app_state.subaccounts.subaccounts = [{\\\"id\\\": {\\\"number\\\": \\\"0\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"1000000000\\\"}], \\\"perpetual_positions\\\": []}, {\\\"id\\\": {\\\"number\\\": \\\"1\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"1000000000\\\"}], \\\"perpetual_positions\\\": []}]\" \"\""
	// bankGenesis := "\".app_state.bank.balances = [{\\\"address\\\": \\\"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6\\\", \\\"coins\\\": [{\\\"denom\\\": \\\"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5\\\", \\\"amount\\\": \\\"10000000000\\\"}]}] | .app_state.bank.supply = [{\\\"denom\\\": \\\"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5\\\", \\\"amount\\\": \\\"10000000000\\\"}]\" \"\""

	genesis := "\".app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}] | .app_state.feetiers.params = {\\\"tiers\\\": [{\\\"name\\\": \\\"1\\\", \\\"maker_fee_ppm\\\": \\\"200\\\", \\\"taker_fee_ppm\\\": \\\"500\\\"}]} | .app_state.subaccounts.subaccounts = [{\\\"id\\\": {\\\"number\\\": \\\"0\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"100000000000\\\"}], \\\"perpetual_positions\\\": []}, {\\\"id\\\": {\\\"number\\\": \\\"1\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"-45001000000\\\"}], \\\"perpetual_positions\\\": [{\\\"perpetual_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"100000000\\\", \\\"funding_index\\\": \\\"0\\\"}]}] | .app_state.clob.liquidations_config = {\\\"max_liquidation_fee_ppm\\\": \\\"5000\\\", \\\"fillable_price_config\\\": {\\\"bankruptcy_adjustment_ppm\\\": \\\"1000000\\\", \\\"spread_to_maintenance_margin_ratio_ppm\\\": \\\"100000\\\"}, \\\"position_block_limits\\\": {\\\"max_position_portion_liquidated_ppm\\\": \\\"1000000\\\", \\\"min_position_notional_liquidated\\\": \\\"1000\\\"}, \\\"subaccount_block_limits\\\": {\\\"max_notional_liquidated\\\": \\\"100000000000000\\\", \\\"max_quantums_insurance_lost\\\": \\\"100000000000000\\\"}} | .app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}] | .app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"-8\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}] | .app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"0\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"1\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"750000\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"2\\\", \\\"initial_margin_ppm\\\": \\\"1000000\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"500000000\\\"}, {\\\"id\\\": \\\"3\\\", \\\"name\\\": \\\"3\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"4\\\", \\\"name\\\": \\\"4\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"800000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"5\\\", \\\"name\\\": \\\"5\\\", \\\"initial_margin_ppm\\\": \\\"500000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"6\\\", \\\"name\\\": \\\"6\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"900000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"7\\\", \\\"name\\\": \\\"7\\\", \\\"initial_margin_ppm\\\": \\\"0\\\", \\\"maintenance_fraction_ppm\\\": \\\"0\\\", \\\"impact_notional\\\": \\\"1000000000\\\"}, {\\\"id\\\": \\\"8\\\", \\\"name\\\": \\\"8\\\", \\\"initial_margin_ppm\\\": \\\"9910\\\", \\\"maintenance_fraction_ppm\\\": \\\"1000000\\\", \\\"impact_notional\\\": \\\"50454000000\\\"}, {\\\"id\\\": \\\"101\\\", \\\"name\\\": \\\"101\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}] | .app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"} | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"-8\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"3\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"BTC-USD 20/10 margin requirements\\\"}, \\\"funding_index\\\": \\\"0\\\"}] | .app_state.bank.balances = [{\\\"address\\\": \\\"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6\\\", \\\"coins\\\": [{\\\"denom\\\": \\\"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5\\\", \\\"amount\\\": \\\"10000000000\\\"}]}] | .app_state.bank.supply = [{\\\"denom\\\": \\\"ibc/8E27BA2D5493AF5636760E354E46004562C46AB7EC0CC4C1CA14E9E20E2545B5\\\", \\\"amount\\\": \\\"10000000000\\\"}, {\\\"denom\\\": \\\"stake\\\", \\\"amount\\\": \\\"30000000000\\\"}]\" \"--price-daemon-enabled=false --max-daemon-unhealthy-seconds=4294967295 --unix-socket-address=/tmp/liquidations_cli_test.sock\""
	network.DeployCustomNetwork(genesis)
	// s.network = network.New(s.T(), s.cfg)

	// _, err = s.network.WaitForHeight(1)
	// s.Require().NoError(err)
}

// TestCLILiquidations creates two subaccounts (where one is undercollateralized), and places a
// maker order from the well-collateralized subaccount that should match with the liquidation order.
// The account which is liquidated and places the maker orders is also the validator's AccAddress.
// After the matching, the subaccounts are queried and assertions are performed on their
// QuoteBalance and PerpetualPositions, along with the balances of the Subaccounts module,
// Distribution module, and insurance fund.
func (s *LiquidationsIntegrationTestSuite) TestCLILiquidations() {
	// val := s.network.Validators[0]
	// ctx := val.ClientCtx

	// currentHeight, err := s.network.LatestHeight()
	// s.Require().NoError(err)

	// goodTilBlock := uint32(currentHeight) + types.ShortBlockWindow
	// clientId := uint64(1)
	// subticks := types.Subticks(50_000_000_000)
	goodTilBlock := uint32(0)
	// quantums := satypes.BaseQuantums(1_000)
	subticks := types.Subticks(50_000_000_000)

	blockHeightQuery := "docker exec interchain-security-instance interchain-security-cd query block --type=height 0 --node tcp://7.7.8.4:26658 -o json"
	data, _, err := network.QueryCustomNetwork(blockHeightQuery)
	var resp blocktypes.Block
	require.NoError(s.T(), s.cfg.Codec.UnmarshalJSON(data, &resp))
	blockHeight := resp.LastCommit.Height

	goodTilBlock = uint32(blockHeight) + types.ShortBlockWindow
	goodTilBlockStr := strconv.Itoa(int(goodTilBlock))

	buyTx := "docker exec interchain-security-instance interchain-security-cd tx clob place-order dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 0 1 0 1 100000000 50000000000 " + goodTilBlockStr + " --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 --chain-id consu --home /consu/validatoralice --keyring-backend test -y --node tcp://7.7.8.4:26658"
	_, _, err = network.QueryCustomNetwork(buyTx)
	s.Require().NoError(err)
	// // Place the maker order that should be filled by the liquidation order.
	// _, err = testutil.MsgPlaceOrderExec(
	// 	ctx,
	// 	s.validatorAddress,
	// 	liqTestSubaccountNumberZero,
	// 	clientId,
	// 	constants.ClobPair_Btc.Id,
	// 	types.Order_SIDE_BUY,
	// 	liqTestMakerOrderQuantums,
	// 	subticks.ToUint64(),
	// 	goodTilBlock,
	// )
	// s.Require().NoError(err)

	// currentHeight, err = s.network.LatestHeight()
	// s.Require().NoError(err)

	// // Wait for a few blocks to ensure the liquidation order was placed, matched, and included
	// // in a block.
	// _, err = s.network.WaitForHeight(currentHeight + 3)
	// s.Require().NoError(err)
	time.Sleep(5 * time.Second)

	// Query both subaccounts.
	accResp, accErr := sa_testutil.MsgQuerySubaccountExec("dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6", liqTestSubaccountNumberZero)
	s.Require().NoError(accErr)

	var subaccountResp satypes.QuerySubaccountResponse
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(accResp.Bytes(), &subaccountResp))
	subaccountZero := subaccountResp.Subaccount

	accResp, err = sa_testutil.MsgQuerySubaccountExec("dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6", liqTestSubaccountNumberOne)
	s.Require().NoError(accErr)

	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(accResp.Bytes(), &subaccountResp))
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
		"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6",
		s.cfg.Codec,
		satypes.ModuleName,
	)
	s.Require().NoError(err)
	s.Require().Equal(
		initialSubaccountModuleAccBalance-makerFee-liquidationFee,
		saModuleUSDCBalance,
	)

	// Check that the insurance fund has expected USDC balance.
	insuranceFundBalance, err := testutil_bank.GetModuleAccUsdcBalance(
		"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6",
		s.cfg.Codec,
		types.InsuranceFundName,
	)

	s.Require().NoError(err)
	s.Require().Equal(liquidationFee, insuranceFundBalance)
	network.CleanupCustomNetwork()
	// Check that the `distribution` module account has expected remaining USDC balance.
	// During `BeginBlock()`, the `fee-collector` module account will send all fees
	// to the `distribution` module account, and the fees will stay in `distribution`
	// until withdrawn. More details at:
	// https://docs.cosmos.network/v0.45/modules/distribution/03_begin_block.html#the-distribution-scheme
	// distrModuleUSDCBalance, err := testutil_bank.GetModuleAccUsdcBalance(
	// 	val,
	// 	s.network.Config.Codec,
	// 	distrtypes.ModuleName,
	// )

	// s.Require().NoError(err)
	// s.Require().Equal(makerFee, distrModuleUSDCBalance)
}
