//go:build all || integration_test

package cli_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/nullify"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	tmcli "github.com/cometbft/cometbft/libs/cli"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithClobPairObjects(t *testing.T, n int) (*network.Network, []types.ClobPair) {
	t.Helper()
	// cfg := network.DefaultConfig(nil)

	// Init Prices genesis state.
	// pricesState := constants.Prices_DefaultGenesisState
	// pricesBuf, pricesErr := cfg.Codec.MarshalJSON(&pricesState)
	// require.NoError(t, pricesErr)
	// cfg.GenesisState[pricestypes.ModuleName] = pricesBuf
	// marketParamGenesis := "\".app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}]\""
	// marketPricesGenesis := "\".app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}]\""

	// Init Perpetuals genesis state.
	// Add additional perps for objects exceeding the default perpetual count.
	// ClobPairs and Perpetuals should be one to one.
	// perpetualsState := constants.Perpetuals_DefaultGenesisState
	// for i := 2; i < n; i++ {
	// 	perpetualsState.Perpetuals = append(
	// 		perpetualsState.Perpetuals,
	// 		perpetualstypes.Perpetual{
	// 			Params: perpetualstypes.PerpetualParams{
	// 				Id:            uint32(i),
	// 				Ticker:        fmt.Sprintf("genesis_test_ticker_%d", i),
	// 				LiquidityTier: 0,
	// 			},
	// 			FundingIndex: dtypes.ZeroInt(),
	// 		},
	// 	)
	// }
	// perpetualsBuf, perpetualsErr := cfg.Codec.MarshalJSON(&perpetualsState)
	// require.NoError(t, perpetualsErr)
	// cfg.GenesisState[perpetualstypes.ModuleName] = perpetualsBuf
	// perpstateGenesis := "\".app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"}\""
	// perpetualGenesis := "\".app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_0\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"1\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_1\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"2\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_2\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"3\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_3\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"4\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_4\\\"}, \\\"funding_index\\\": \\\"0\\\"}]\""
	// liquidityTiersGenesis := "\".app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"Large-Cap\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"Mid-Cap\\\", \\\"initial_margin_ppm\\\": \\\"300000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1667000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"Small-Cap\\\", \\\"initial_margin_ppm\\\": \\\"400000\\\", \\\"maintenance_fraction_ppm\\\": \\\"700000\\\", \\\"impact_notional\\\": \\\"1250000000\\\"}]\""
	// Init Clob State.
	state := types.GenesisState{}
	// require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		clobPair := types.ClobPair{
			Id: uint32(i),
			Metadata: &types.ClobPair_PerpetualClobMetadata{
				PerpetualClobMetadata: &types.PerpetualClobMetadata{PerpetualId: uint32(i)},
			},
			SubticksPerTick:  5,
			StepBaseQuantums: 5,
			Status:           types.ClobPair_STATUS_ACTIVE,
		}
		nullify.Fill(&clobPair) //nolint:staticcheck
		state.ClobPairs = append(state.ClobPairs, clobPair)
	}
	// buf, err := cfg.Codec.MarshalJSON(&state)
	// require.NoError(t, err)
	// cfg.GenesisState[types.ModuleName] = buf
	// clobgenesis := "\".app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"1\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"1\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"2\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"2\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"3\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"3\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"4\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"4\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}]\""
	fullGenesisTwo := "\".app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}] | .app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}] | .app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"} | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_0\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"1\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_1\\\"}, \\\"funding_index\\\": \\\"0\\\"}] | .app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"Large-Cap\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"Mid-Cap\\\", \\\"initial_margin_ppm\\\": \\\"300000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1667000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"Small-Cap\\\", \\\"initial_margin_ppm\\\": \\\"400000\\\", \\\"maintenance_fraction_ppm\\\": \\\"700000\\\", \\\"impact_notional\\\": \\\"1250000000\\\"}] | .app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"1\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"1\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}]\" \"\""
	fullGenesisFive := "\".app_state.prices.market_params = [{\\\"id\\\": \\\"0\\\", \\\"pair\\\": \\\"BTC-USD\\\", \\\"exponent\\\": \\\"-5\\\", \\\"min_exchanges\\\": \\\"2\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tBTCUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC/USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTCUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XXBTZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"},  {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"BTC-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}, {\\\"id\\\": \\\"1\\\", \\\"pair\\\": \\\"ETH-USD\\\", \\\"exponent\\\": \\\"-6\\\", \\\"min_exchanges\\\": \\\"1\\\", \\\"min_price_change_ppm\\\": \\\"50\\\", \\\"exchange_config_json\\\": \\\"{\\\\\\\"exchanges\\\\\\\": [{\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Binance\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"BinanceUS\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitfinex\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"tETHUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bitstamp\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH/USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Bybit\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETHUSDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CoinbasePro\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"CryptoCom\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Kraken\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"XETHZUSD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Mexc\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH_USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}, {\\\\\\\"exchangeName\\\\\\\": \\\\\\\"Okx\\\\\\\",\\\\\\\"ticker\\\\\\\": \\\\\\\"ETH-USDT\\\\\\\",\\\\\\\"adjustByMarket\\\\\\\": \\\\\\\"USDT-USD\\\\\\\"}]}\\\"}] | .app_state.prices.market_prices = [{\\\"exponent\\\": \\\"-5\\\", \\\"price\\\": \\\"5000000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"exponent\\\": \\\"-6\\\", \\\"price\\\": \\\"3000000000\\\"}] | .app_state.perpetuals.params = {\\\"funding_rate_clamp_factor_ppm\\\": \\\"6000000\\\", \\\"min_num_votes_per_sample\\\": \\\"15\\\", \\\"premium_vote_clamp_factor_ppm\\\": \\\"60000000\\\"} | .app_state.perpetuals.perpetuals = [{\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"0\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_0\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"1\\\", \\\"liquidity_tier\\\": \\\"1\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_1\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"2\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_2\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"3\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_3\\\"}, \\\"funding_index\\\": \\\"0\\\"}, {\\\"params\\\": {\\\"atomic_resolution\\\": \\\"0\\\", \\\"default_funding_ppm\\\": \\\"0\\\", \\\"id\\\": \\\"4\\\", \\\"liquidity_tier\\\": \\\"0\\\", \\\"market_id\\\": \\\"0\\\", \\\"ticker\\\": \\\"genesis_test_ticker_4\\\"}, \\\"funding_index\\\": \\\"0\\\"}] | .app_state.perpetuals.liquidity_tiers = [{\\\"id\\\": \\\"0\\\", \\\"name\\\": \\\"Large-Cap\\\", \\\"initial_margin_ppm\\\": \\\"200000\\\", \\\"maintenance_fraction_ppm\\\": \\\"500000\\\", \\\"impact_notional\\\": \\\"2500000000\\\"}, {\\\"id\\\": \\\"1\\\", \\\"name\\\": \\\"Mid-Cap\\\", \\\"initial_margin_ppm\\\": \\\"300000\\\", \\\"maintenance_fraction_ppm\\\": \\\"600000\\\", \\\"impact_notional\\\": \\\"1667000000\\\"}, {\\\"id\\\": \\\"2\\\", \\\"name\\\": \\\"Small-Cap\\\", \\\"initial_margin_ppm\\\": \\\"400000\\\", \\\"maintenance_fraction_ppm\\\": \\\"700000\\\", \\\"impact_notional\\\": \\\"1250000000\\\"}] | .app_state.clob.clob_pairs = [{\\\"id\\\": \\\"0\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"0\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"1\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"1\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"2\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"2\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"3\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"3\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}, {\\\"id\\\": \\\"4\\\", \\\"perpetual_clob_metadata\\\": {\\\"perpetual_id\\\": \\\"4\\\"}, \\\"step_base_quantums\\\": \\\"5\\\", \\\"subticks_per_tick\\\": \\\"5\\\", \\\"quantum_conversion_exponent\\\": \\\"0\\\", \\\"status\\\": \\\"STATUS_ACTIVE\\\"}]\" \"\""
	var genesis string
	if n == 2 {
		genesis = fullGenesisTwo
	} else if n == 5 {
		genesis = fullGenesisFive
	}

	setupCmd := exec.Command("bash", "-c", "cd ../../../../ethos/ethos-chain && ./e2e-setup -setup "+genesis)

	fmt.Println("Running setup command", setupCmd.String())
	var out bytes.Buffer
	var stderr bytes.Buffer
	setupCmd.Stdout = &out
	setupCmd.Stderr = &stderr
	err := setupCmd.Run()
	if err != nil {
		t.Fatalf("Failed to set up environment: %v, stdout: %s, stderr: %s", err, out.String(), stderr.String())
	}
	fmt.Println("Setup output:", out.String())
	time.Sleep(5 * time.Second)
	require.NoError(t, err)

	return nil, state.ClobPairs
}

func TestShowClobPair(t *testing.T) {
	_, objs := networkWithClobPairObjects(t, 2)

	// ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	for _, tc := range []struct {
		desc string
		id   uint32

		args []string
		err  string
		obj  types.ClobPair
	}{
		{
			desc: "found",
			id:   objs[0].Id,

			args: common,
			obj:  objs[0],
			err:  "",
		},
		{
			desc: "not found",
			id:   uint32(100000),

			args: common,
			err:  "not found",
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {

			cfg := network.DefaultConfig(nil)
			// args := []string{
			// 	lib.UintToString(tc.id),
			// }
			// args = append(args, tc.args...)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowClobPair(), args)

			clobCmd := exec.Command("bash", "-c", "docker exec interchain-security-instance interchain-security-cd query clob show-clob-pair "+fmt.Sprintf("%d", tc.id)+" --node tcp://7.7.8.4:26658 -o json")
			var clobOut bytes.Buffer
			var clobErr bytes.Buffer
			clobCmd.Stdout = &clobOut
			clobCmd.Stderr = &clobErr
			err := clobCmd.Run()
			data := clobOut.Bytes()
			stderrOutput := clobErr.String()

			if tc.err != "" {
				require.Error(t, err)
				require.Contains(t, stderrOutput, tc.err)
			} else {
				require.NoError(t, err)
				var resp types.QueryClobPairResponse
				require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
				require.NotNil(t, resp.ClobPair)
				require.Equal(t,
					nullify.Fill(&tc.obj),        //nolint:staticcheck
					nullify.Fill(&resp.ClobPair), //nolint:staticcheck
				)
			}
		})
	}
	stopCmd := exec.Command("bash", "-c", "docker stop interchain-security-instance")
	if err := stopCmd.Run(); err != nil {
		t.Fatalf("Failed to stop Docker container: %v", err)
	}
	fmt.Println("Stopped Docker container")
	// Remove the Docker container
	removeCmd := exec.Command("bash", "-c", "docker rm interchain-security-instance")
	if err := removeCmd.Run(); err != nil {
		t.Fatalf("Failed to remove Docker container: %v", err)
	}
	fmt.Println("Removed Docker container")
}

func TestListClobPair(t *testing.T) {
	_, objs := networkWithClobPairObjects(t, 5)

	// ctx := net.Validators[0].ClientCtx
	cfg := network.DefaultConfig(nil)
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			argsString := strings.Join(args, " ")
			commandString := "docker exec interchain-security-instance interchain-security-cd query clob list-clob-pair --node tcp://7.7.8.4:26658 " + argsString
			clobCmd := exec.Command("bash", "-c", commandString)
			var clobOut bytes.Buffer
			clobCmd.Stdout = &clobOut
			err := clobCmd.Run()
			data := clobOut.Bytes()

			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
			require.NoError(t, err)
			var resp types.QueryClobPairAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(objs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			var nextKeyStr string
			if next != nil {
				nextKeyStr = base64.StdEncoding.EncodeToString(next)
			}
			args := request([]byte(nextKeyStr), 0, uint64(step), false)
			// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
			argsString := strings.Join(args, " ")
			commandString := "docker exec interchain-security-instance interchain-security-cd query clob list-clob-pair --node tcp://7.7.8.4:26658 " + argsString
			fmt.Println("Command: ", commandString)
			clobCmd := exec.Command("bash", "-c", commandString)
			var clobOut bytes.Buffer
			clobCmd.Stdout = &clobOut
			err := clobCmd.Run()
			data := clobOut.Bytes()
			require.NoError(t, err)
			var resp types.QueryClobPairAllResponse
			require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
			require.LessOrEqual(t, len(resp.ClobPair), step)
			require.Subset(t,
				nullify.Fill(objs),          //nolint:staticcheck
				nullify.Fill(resp.ClobPair), //nolint:staticcheck
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		argsString := strings.Join(args, " ")
		commandString := "docker exec interchain-security-instance interchain-security-cd query clob list-clob-pair --node tcp://7.7.8.4:26658 " + argsString

		clobCmd := exec.Command("bash", "-c", commandString)
		var clobOut bytes.Buffer
		clobCmd.Stdout = &clobOut
		err := clobCmd.Run()
		data := clobOut.Bytes()
		// out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListClobPair(), args)
		require.NoError(t, err)
		var resp types.QueryClobPairAllResponse
		require.NoError(t, cfg.Codec.UnmarshalJSON(data, &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),          //nolint:staticcheck
			nullify.Fill(resp.ClobPair), //nolint:staticcheck
		)
	})

	stopCmd := exec.Command("bash", "-c", "docker stop interchain-security-instance")
	if err := stopCmd.Run(); err != nil {
		t.Fatalf("Failed to stop Docker container: %v", err)
	}
	fmt.Println("Stopped Docker container")
	// Remove the Docker container
	removeCmd := exec.Command("bash", "-c", "docker rm interchain-security-instance")
	if err := removeCmd.Run(); err != nil {
		t.Fatalf("Failed to remove Docker container: %v", err)
	}
	fmt.Println("Removed Docker container")
}
