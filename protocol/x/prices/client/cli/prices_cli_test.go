//go:build integration_test

package cli_test

import (
	"fmt"

	"path/filepath"
	"testing"

	networktestutil "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/daemons/configs"
	"github.com/dydxprotocol/v4/daemons/pricefeed"
	"github.com/dydxprotocol/v4/daemons/pricefeed/client"
	"github.com/dydxprotocol/v4/testutil/appoptions"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/testutil/network"
	epochstypes "github.com/dydxprotocol/v4/x/epochs/types"
	"github.com/dydxprotocol/v4/x/prices/client/testutil"
	"github.com/dydxprotocol/v4/x/prices/types"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/suite"
)

var (
	genesisState = constants.Prices_MultiExchangeMarketGenesisState
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

			// Enable the PriceFeed daemon in the integration tests.
			appOptions.Set(pricefeed.FlagPriceFeedEnabled, true)
			homeDir := filepath.Join(testval.Dir, "simd")
			configs.WriteDefaultPricefeedExchangeToml(homeDir) // must manually create config file.
			appOptions.Set(pricefeed.FlagPriceFeedPriceUpdaterLoopDelayMs, 1000)

			// Enable the common gRPC daemon server.
			appOptions.Set(pricefeed.GrpcAddress, testval.AppConfig.GRPC.Address)
		},
	})

	s.cfg.Mnemonics = append(s.cfg.Mnemonics, validatorMnemonic)
	s.cfg.ChainID = app.AppName

	// Set min gas prices to zero so that we can submit transactions with zero gas price.
	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	// Setting genesis state for Prices.
	state := genesisState

	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf

	// Ensure that no funding-related epochs will occur during this test.
	epstate := constants.GenerateEpochGenesisStateWithoutFunding()

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

func (s *PricesIntegrationTestSuite) TestCLIPrices_AllErrorResponses_NoPriceUpdate() {
	// Setup.
	ts := s.T()
	testutil.SetupExchangeResponses(ts, testutil.AllErrorResponses, genesisState)

	// Run.
	s.network = network.New(ts, s.cfg)

	_, err := s.network.WaitForHeight(5)
	s.Require().NoError(err)

	// Verify.
	val := s.network.Validators[0]
	ctx := val.ClientCtx
	resp, err := testutil.MsgQueryAllMarketExec(ctx)
	s.Require().NoError(err)

	var allMarketQueryResponse types.QueryAllMarketsResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &allMarketQueryResponse))
	s.Require().Len(allMarketQueryResponse.Market, 2)

	s.Require().Equal(uint32(0), allMarketQueryResponse.Market[0].Id)
	s.Require().Equal(uint64(0), allMarketQueryResponse.Market[0].Price)

	s.Require().Equal(uint32(1), allMarketQueryResponse.Market[1].Id)
	s.Require().Equal(uint64(0), allMarketQueryResponse.Market[1].Price)
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_PartialValidResponses_PartialPriceUpdate() {
	// Setup.
	ts := s.T()
	testutil.SetupExchangeResponses(ts, testutil.MixedResponses, genesisState)

	// Run.
	s.network = network.New(ts, s.cfg)

	_, err := s.network.WaitForHeight(5)
	s.Require().NoError(err)

	// Verify.
	val := s.network.Validators[0]
	ctx := val.ClientCtx
	resp, err := testutil.MsgQueryAllMarketExec(ctx)
	s.Require().NoError(err)

	var allMarketQueryResponse types.QueryAllMarketsResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &allMarketQueryResponse))
	s.Require().Len(allMarketQueryResponse.Market, 2)

	s.Require().Equal(uint32(0), allMarketQueryResponse.Market[0].Id)
	// No price update; 2 error and 1 valid responses. However, min req for valid exchange prices is 2.
	s.Require().Equal(uint64(0), allMarketQueryResponse.Market[0].Price)

	s.Require().Equal(uint32(1), allMarketQueryResponse.Market[1].Id)
	// Valid price update; 1 error and 2 valid responses. Min req for valid exchange prices is 2.
	// Median of 9_000 and 9_002.
	s.Require().Equal(uint64(9_001_000_000), allMarketQueryResponse.Market[1].Price)
}

func (s *PricesIntegrationTestSuite) TestCLIPrices_AllValidResponses_ValidPriceUpdate() {
	// Setup.
	ts := s.T()
	testutil.SetupExchangeResponses(ts, testutil.AllValidResponses, genesisState)

	// Run.
	s.network = network.New(ts, s.cfg)

	_, err := s.network.WaitForHeight(5)
	s.Require().NoError(err)

	// Verify.
	val := s.network.Validators[0]
	ctx := val.ClientCtx
	resp, err := testutil.MsgQueryAllMarketExec(ctx)
	s.Require().NoError(err)

	var allMarketQueryResponse types.QueryAllMarketsResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &allMarketQueryResponse))
	s.Require().Len(allMarketQueryResponse.Market, 2)

	s.Require().Equal(uint32(0), allMarketQueryResponse.Market[0].Id)
	// Median of 100, 101, 102.
	s.Require().Equal(uint64(10_100_000), allMarketQueryResponse.Market[0].Price)

	s.Require().Equal(uint32(1), allMarketQueryResponse.Market[1].Id)
	// Median of 9_000, 9_001, 9_002.
	s.Require().Equal(uint64(9_001_000_000), allMarketQueryResponse.Market[1].Price)
}
