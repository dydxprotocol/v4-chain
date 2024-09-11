//go:build all || integration_test

package cli_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	appconstants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	appflags "github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/configs"
	daemonflags "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/flags"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/pricefeed/client"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/appoptions"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/client/testutil"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

var (
	genesisState = constants.Prices_MultiExchangeMarketGenesisState

	medianUpdatedMarket1price = uint64(9_001_000_000)
	medianUpdatedMarket0Price = uint64(10_100_000)

	// expectedPricesWithNoUpdates is the set of genesis prices.
	expectedPricesWithNoUpdates = map[uint32]ve.VEPricePair{
		0: {
			SpotPrice: genesisState.MarketPrices[0].SpotPrice,
			PnlPrice:  genesisState.MarketPrices[0].PnlPrice,
		},
		1: {
			SpotPrice: genesisState.MarketPrices[1].SpotPrice,
			PnlPrice:  genesisState.MarketPrices[1].PnlPrice,
		},
	}

	// expectedPricesWithPartialUpdate is the expected prices after updating prices with the partial update.
	expectedPricesWithPartialUpdate = map[uint32]ve.VEPricePair{
		// No price update; 2 error and 1 valid responses. However, min req for valid exchange prices is 2.
		0: {
			SpotPrice: genesisState.MarketPrices[0].SpotPrice,
			PnlPrice:  genesisState.MarketPrices[0].PnlPrice,
		},
		// Valid price update; 1 error and 2 valid responses. Min req for valid exchange prices is 2.
		// Median of 9_000 and 9_002.
		1: {
			SpotPrice: medianUpdatedMarket1price,
			PnlPrice:  medianUpdatedMarket1price,
		},
	}

	// expectedPricesWithFullUpdate is the expected prices after updating all prices.
	expectedPricesWithFullUpdate = map[uint32]ve.VEPricePair{
		0: {
			SpotPrice: medianUpdatedMarket0Price,
			PnlPrice:  medianUpdatedMarket0Price,
		},
		1: {
			SpotPrice: medianUpdatedMarket1price,
			PnlPrice:  medianUpdatedMarket1price,
		},
	}
)

type PricesIntegrationTestSuite struct {
	suite.Suite

	validatorAddress sdk.AccAddress
	cfg              network.Config
}

func TestPricesIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &PricesIntegrationTestSuite{})
}

func (s *PricesIntegrationTestSuite) SetupTest() {
	s.T().Log("setting up prices integration test")

	// Deterministic Mnemonic.
	validatorMnemonic := constants.AliceMnenomic

	// Generated from the above Mnemonic.
	s.validatorAddress = constants.AliceAccAddress

	appOptions := appoptions.NewFakeAppOptions()

	// Configure test network.
	s.cfg = network.DefaultConfig(&network.NetworkConfigOptions{
		AppOptions: appOptions,
		OnNewApp: func(val networktestutil.ValidatorI) {
			testval, ok := val.(networktestutil.Validator)
			if !ok {
				panic("incorrect validator type")
			}

			// Disable the Liquidations daemon.
			appOptions.Set(daemonflags.FlagDeleveragingDaemonEnabled, false)

			// Enable the Price daemon.
			appOptions.Set(daemonflags.FlagPriceDaemonEnabled, true)
			appOptions.Set(daemonflags.FlagPriceDaemonLoopDelayMs, 1_000)

			homeDir := filepath.Join(testval.Dir, "simd")
			configs.WriteDefaultPricefeedExchangeToml(homeDir) // must manually create config file.

			// Make sure the daemon is using the correct GRPC address.
			appOptions.Set(appflags.GrpcAddress, testval.AppConfig.GRPC.Address)
		},
	})

	s.cfg.Mnemonics = append(s.cfg.Mnemonics, validatorMnemonic)
	s.cfg.ChainID = appconstants.AppName

	// Set min gas prices to zero so that we can submit transactions with zero gas price.
	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	// // Gock setup.
	defer gock.Off()         // Flush pending mocks after test execution.
	gock.DisableNetworking() // Disables real networking.
	gock.InterceptClient(&client.HttpClient)

	// Starting the network is delayed on purpose.
	// `gock` HTTP mocking must be setup before the network starts.
}

// expectMarketPricesWithTimeout waits for the specified timeout for the market prices to be updated with the
// expected values. If the prices are not updated to match the expected prices within the timeout, the test fails.
func (s *PricesIntegrationTestSuite) expectMarketPricesWithTimeout(prices map[uint32]ve.VEPricePair, timeout time.Duration) {
	start := time.Now()

	for {
		if time.Since(start) > timeout {
			s.Require().Fail("timed out waiting for market prices")
		}

		time.Sleep(100 * time.Millisecond)

		resp, err := testutil.MsgQueryAllMarketPriceExec()
		s.Require().NoError(err)

		var allMarketPricesQueryResponse types.QueryAllMarketPricesResponse
		s.Require().NoError(s.cfg.Codec.UnmarshalJSON(resp, &allMarketPricesQueryResponse))

		if len(allMarketPricesQueryResponse.MarketPrices) != len(prices) {
			continue
		}

		// Compare for equality. If prices are not equal, continue waiting.
		actualPrices := make(map[uint32]ve.VEPricePair, len(allMarketPricesQueryResponse.MarketPrices))
		for _, actualPrice := range allMarketPricesQueryResponse.MarketPrices {
			actualPrices[actualPrice.Id] = ve.VEPricePair{
				SpotPrice: actualPrice.SpotPrice,
				PnlPrice:  actualPrice.PnlPrice,
			}
		}

		for marketId, expectedPrice := range prices {
			actualPrice, ok := actualPrices[marketId]
			if !ok {
				continue
			}
			if actualPrice.SpotPrice != expectedPrice.SpotPrice || actualPrice.PnlPrice != expectedPrice.PnlPrice {
				continue
			}
		}

		// All prices match - return.
		return
	}
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_AllEmptyResponses_NoPriceUpdate() {
	// Setup.
	ts := s.T()

	testutil.SetupExchangeResponses(ts, testutil.EmptyResponses_AllExchanges)

	// // Run.
	// s.network = network.New(ts, s.cfg)
	genesis := "\".app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}] | .app_state.feetiers.params = {\\\"tiers\\\": [{\\\"name\\\": \\\"1\\\", \\\"maker_fee_ppm\\\": \\\"200\\\", \\\"taker_fee_ppm\\\": \\\"500\\\"}]}\" \"--deleveraging-daemon-enabled=false --price-daemon-enabled=true --price-daemon-loop-delay-ms=1000\""
	network.DeployCustomNetwork(genesis)

	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithNoUpdates, 30*time.Second)

	network.CleanupCustomNetwork()
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_PartialResponses_PartialPriceUpdate() {
	// Setup.
	ts := s.T()

	// Add logging to see what's going on in circleCI.
	testutil.SetupExchangeResponses(ts, testutil.PartialResponses_AllExchanges_Eth9001)

	genesis := "\".app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}] | .app_state.feetiers.params = {\\\"tiers\\\": [{\\\"name\\\": \\\"1\\\", \\\"maker_fee_ppm\\\": \\\"200\\\", \\\"taker_fee_ppm\\\": \\\"500\\\"}]}\" \"--deleveraging-daemon-enabled=false --price-daemon-enabled=true --price-daemon-loop-delay-ms=1000\""
	network.DeployCustomNetwork(genesis)

	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithPartialUpdate, 30*time.Second)
	network.CleanupCustomNetwork()
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_AllValidResponses_ValidPriceUpdate() {
	// Setup.
	ts := s.T()
	testutil.SetupExchangeResponses(ts, testutil.FullResponses_AllExchanges_Btc101_Eth9001)
	genesis := "\".app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}] | .app_state.feetiers.params = {\\\"tiers\\\": [{\\\"name\\\": \\\"1\\\", \\\"maker_fee_ppm\\\": \\\"200\\\", \\\"taker_fee_ppm\\\": \\\"500\\\"}]}\" \"--deleveraging-daemon-enabled=false --price-daemon-enabled=true --price-daemon-loop-delay-ms=1000\""
	network.DeployCustomNetwork(genesis)
	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithFullUpdate, 30*time.Second)
	network.CleanupCustomNetwork()
}
