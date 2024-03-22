//go:build all || integration_test

package cli_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appconstants "github.com/dydxprotocol/v4-chain/protocol/app/constants"
	appflags "github.com/dydxprotocol/v4-chain/protocol/app/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/configs"
	daemonflags "github.com/dydxprotocol/v4-chain/protocol/daemons/flags"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/appoptions"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	feetierstypes "github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/client/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

var (
	genesisState = constants.Prices_MultiExchangeMarketGenesisState

	medianUpdatedMarket1price = uint64(9_001_000_000)
	medianUpdatedMarket0Price = uint64(10_100_000)

	// expectedPricesWithNoUpdates is the set of genesis prices.
	expectedPricesWithNoUpdates = map[uint32]uint64{
		0: genesisState.MarketPrices[0].Price,
		1: genesisState.MarketPrices[1].Price,
	}

	// expectedPricesWithPartialUpdate is the expected prices after updating prices with the partial update.
	expectedPricesWithPartialUpdate = map[uint32]uint64{
		// No price update; 2 error and 1 valid responses. However, min req for valid exchange prices is 2.
		0: genesisState.MarketPrices[0].Price,
		// Valid price update; 1 error and 2 valid responses. Min req for valid exchange prices is 2.
		// Median of 9_000 and 9_002.
		1: medianUpdatedMarket1price,
	}

	// expectedPricesWithFullUpdate is the expected prices after updating all prices.
	expectedPricesWithFullUpdate = map[uint32]uint64{
		0: medianUpdatedMarket0Price,
		1: medianUpdatedMarket1price,
	}
)

type PricesIntegrationTestSuite struct {
	suite.Suite

	validatorAddress sdk.AccAddress
	cfg              network.Config
	network          *network.Network
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
			appOptions.Set(daemonflags.FlagLiquidationDaemonEnabled, false)

			// Disable the Bridge Daemon.
			appOptions.Set(daemonflags.FlagBridgeDaemonEnabled, false)

			// Enable the Price daemon.
			appOptions.Set(daemonflags.FlagPriceDaemonEnabled, true)
			appOptions.Set(daemonflags.FlagPriceDaemonLoopDelayMs, 1_000)
			appOptions.Set(daemonflags.FlagOracleEnabled, false)

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

	// Setting genesis state for Prices.
	state := genesisState

	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf

	// Ensure that no funding-related epochs will occur during this test.
	epstate := constants.GenerateEpochGenesisStateWithoutFunding()

	feeTiersState := feetierstypes.GenesisState{}
	feeTiersState.Params = constants.PerpetualFeeParams

	feeTiersBuf, err := s.cfg.Codec.MarshalJSON(&feeTiersState)
	s.Require().NoError(err)
	s.cfg.GenesisState[feetierstypes.ModuleName] = feeTiersBuf

	epbuf, err := s.cfg.Codec.MarshalJSON(&epstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	// Gock setup.
	defer gock.Off()         // Flush pending mocks after test execution.
	gock.DisableNetworking() // Disables real networking.
	gock.InterceptClient(&client.HttpClient)

	// Starting the network is delayed on purpose.
	// `gock` HTTP mocking must be setup before the network starts.
}

// expectMarketPricesWithTimeout waits for the specified timeout for the market prices to be updated with the
// expected values. If the prices are not updated to match the expected prices within the timeout, the test fails.
func (s *PricesIntegrationTestSuite) expectMarketPricesWithTimeout(prices map[uint32]uint64, timeout time.Duration) {
	start := time.Now()

	for {
		if time.Since(start) > timeout {
			s.Require().Fail("timed out waiting for market prices")
		}

		time.Sleep(100 * time.Millisecond)

		val := s.network.Validators[0]
		ctx := val.ClientCtx
		resp, err := testutil.MsgQueryAllMarketPriceExec(ctx)
		s.Require().NoError(err)

		var allMarketPricesQueryResponse types.QueryAllMarketPricesResponse
		s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &allMarketPricesQueryResponse))

		if len(allMarketPricesQueryResponse.MarketPrices) != len(prices) {
			continue
		}

		// Compare for equality. If prices are not equal, continue waiting.
		actualPrices := make(map[uint32]uint64, len(allMarketPricesQueryResponse.MarketPrices))
		for _, actualPrice := range allMarketPricesQueryResponse.MarketPrices {
			actualPrices[actualPrice.Id] = actualPrice.Price
		}

		for marketId, expectedPrice := range prices {
			actualPrice, ok := actualPrices[marketId]
			if !ok {
				continue
			}
			if actualPrice != expectedPrice {
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

	// Run.
	s.network = network.New(ts, s.cfg)

	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithNoUpdates, 30*time.Second)
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_PartialResponses_PartialPriceUpdate() {
	// Setup.
	ts := s.T()

	// Add logging to see what's going on in circleCI.
	testutil.SetupExchangeResponses(ts, testutil.PartialResponses_AllExchanges_Eth9001)

	// Run.
	s.network = network.New(ts, s.cfg)

	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithPartialUpdate, 30*time.Second)
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_AllValidResponses_ValidPriceUpdate() {
	// Setup.
	ts := s.T()
	testutil.SetupExchangeResponses(ts, testutil.FullResponses_AllExchanges_Btc101_Eth9001)

	// Run.
	s.network = network.New(ts, s.cfg)

	// Verify.
	s.expectMarketPricesWithTimeout(expectedPricesWithFullUpdate, 30*time.Second)
}
