//go:build all || container_test

package containertest

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"
	"github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/client/types"
	"github.com/dydxprotocol/v4-chain/protocol/testing/version"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/daemons/pricefeed/exchange_config"
	assets "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
	clob "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	prices "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const expectDirName = "expect"

var acceptFlag = flag.Bool("accept", false, "Accept new values for expect files")

// Compare a message against an expected output. Use flag `-accept` to write or modify expected output.
// Expected output will read/written from `expect/{testName}_{tag}.expect`.
func assertExpectedProto(t *testing.T, tag string, message proto.Message) {
	expectFilePath := filepath.Join(expectDirName, fmt.Sprintf("%s_%s.expect", t.Name(), tag))
	marshaler := &jsonpb.Marshaler{
		Indent: "  ",
	}
	actual, err := marshaler.MarshalToString(message)
	assert.NoError(t, err)

	if *acceptFlag {
		err = os.WriteFile(expectFilePath, []byte(actual), 0644)
		assert.NoError(t, err)
	} else {
		expected, err := os.ReadFile(expectFilePath)
		assert.NoError(t, err)
		assert.JSONEq(t, string(expected), actual, "rerun with -accept to accept all new output")
	}
}

// expectProto is like assertExpectedProto, but returns a bool instead of calling assert.
// It returns true if the proto serialized into a JSON file specified by the tag matches the message.
// It does not write to the expect file. If the accept flag is set, it will always return true.
func expectProto(t *testing.T, tag string, message proto.Message) bool {
	expectFilePath := filepath.Join(expectDirName, fmt.Sprintf("%s_%s.expect", t.Name(), tag))
	marshaler := &jsonpb.Marshaler{
		Indent: "  ",
	}
	actual, err := marshaler.MarshalToString(message)
	assert.NoError(t, err)

	if *acceptFlag {
		return true
	} else {
		expected, err := os.ReadFile(expectFilePath)
		assert.NoError(t, err)
		var expectedJSONAsInterface, actualJSONAsInterface interface{}

		err = json.Unmarshal([]byte(expected), &expectedJSONAsInterface)
		require.NoError(t, err)

		err = json.Unmarshal([]byte(actual), &actualJSONAsInterface)
		require.NoError(t, err)

		return assert.ObjectsAreEqual(expectedJSONAsInterface, actualJSONAsInterface)
	}
}

func TestPlaceOrder(t *testing.T) {
	// TODO(DEC-2198): Reenable these tests after fixing flakiness on CI.
	// Seems to occur only because multiple container tests run.
	if os.Getenv("SKIP_DISABLED") != "" {
		t.Skip("Skipping disabled test")
	}
	testnet, err := NewTestnet()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]

	assert.NoError(t, BroadcastTx(
		node,
		&clob.MsgPlaceOrder{
			Order: clob.Order{
				OrderId: clob.OrderId{
					ClientId: 0,
					SubaccountId: satypes.SubaccountId{
						Owner:  constants.AliceAccAddress.String(),
						Number: 0,
					},
					ClobPairId: 0,
				},
				Side:     clob.Order_SIDE_BUY,
				Quantums: 10_000_000,
				Subticks: 1_000_000,
				GoodTilOneof: &clob.Order_GoodTilBlock{
					GoodTilBlock: 20,
				},
			},
		},
		constants.AliceAccAddress.String(),
	))
	// TODO(CLOB-905): place another matching order, and verify that the trade is executed.
}

func TestBankSend(t *testing.T) {
	if os.Getenv("SKIP_DISABLED") != "" {
		t.Skip("Skipping disabled test")
	}
	testnet, err := NewTestnet()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]

	resp, err := Query(
		node,
		bank.NewQueryClient,
		bank.QueryClient.AllBalances,
		&bank.QueryAllBalancesRequest{
			Address: constants.AliceAccAddress.String(),
		},
	)
	assert.NoError(t, err)
	assertExpectedProto(t, "aliceInitialBalances", resp)

	resp, err = Query(
		node,
		bank.NewQueryClient,
		bank.QueryClient.AllBalances,
		&bank.QueryAllBalancesRequest{
			Address: constants.BobAccAddress.String(),
		},
	)
	assert.NoError(t, err)
	assertExpectedProto(t, "bobInitialBalances", resp)

	assert.NoError(t, BroadcastTx(
		node,
		&bank.MsgSend{
			FromAddress: constants.BobAccAddress.String(),
			ToAddress:   constants.AliceAccAddress.String(),
			Amount: []sdk.Coin{
				sdk.NewCoin(assets.AssetUsdc.Denom, sdkmath.NewInt(1)),
			},
		},
		constants.BobAccAddress.String(),
	))
	err = node.Wait(2)
	assert.NoError(t, err)

	resp, err = Query(
		node,
		bank.NewQueryClient,
		bank.QueryClient.AllBalances,
		&bank.QueryAllBalancesRequest{
			Address: constants.AliceAccAddress.String(),
		},
	)
	assert.NoError(t, err)
	assertExpectedProto(t, "aliceFinalBalances", resp)

	resp, err = Query(
		node,
		bank.NewQueryClient,
		bank.QueryClient.AllBalances,
		&bank.QueryAllBalancesRequest{
			Address: constants.BobAccAddress.String(),
		},
	)
	assert.NoError(t, err)
	assertExpectedProto(t, "bobFinalBalances", resp)
}

// assertExpectedPrices compares a message against an expected output. This method utilized the write
// functionality of `assertExpectedProto` and is useful to run in order to ensure that the test output
// is written to appropriate files.
func assertExpectedPrices(t *testing.T, node *Node, marketTags map[types.MarketId]string) {
	for marketId, tag := range marketTags {
		resp, err := Query(
			node,
			prices.NewQueryClient,
			prices.QueryClient.MarketPrice,
			&prices.QueryMarketPriceRequest{
				Id: marketId,
			},
		)
		require.NoError(t, err)
		assertExpectedProto(t, tag, resp)
	}
}

// expectPrices evaluates if the current market prices, when individually queried, match the expected prices.
func expectPrices(t *testing.T, node *Node, marketTags map[types.MarketId]string) bool {
	for marketId, tag := range marketTags {
		resp, err := Query(
			node,
			prices.NewQueryClient,
			prices.QueryClient.MarketPrice,
			&prices.QueryMarketPriceRequest{
				Id: marketId,
			},
		)
		require.NoError(t, err)
		if !expectProto(t, tag, resp) {
			return false
		}
	}
	return true
}

// assertPricesWithTimeout polls the node for the expected prices until the timeout is reached. If the
// accept flag is set, it will wait the full duration and then write the current prices to the appropriate files
// based on the contents of marketTags.
func assertPricesWithTimeout(t *testing.T, node *Node, marketTags map[types.MarketId]string, timeout time.Duration) {
	start := time.Now()
	for {
		// If we're not accepting, return as soon as we see the expected prices. Use short circuit evaluation
		// to skip price comparison when the accept flag is set.
		if !*acceptFlag && expectPrices(t, node, marketTags) {
			return
		}

		// When we see the timeout, we should either fail or write the expected prices.
		if time.Since(start) > timeout {
			if *acceptFlag {
				// Write prices!
				assertExpectedPrices(t, node, marketTags)
			} else {
				require.Fail(t, "timed out waiting for expected prices")
			}
		}

		// Sleep for the poll interval.
		time.Sleep(100 * time.Millisecond)
	}
}

func TestMarketPrices(t *testing.T) {
	if os.Getenv("SKIP_DISABLED") != "" {
		t.Skip("Skipping disabled test")
	}
	testnet, err := NewTestnet()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	testnet.setPrice(exchange_config.MARKET_BTC_USD, 50001)
	testnet.setPrice(exchange_config.MARKET_ETH_USD, 55002)
	testnet.setPrice(exchange_config.MARKET_LINK_USD, 55003)

	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]

	expectedPrices := map[types.MarketId]string{
		exchange_config.MARKET_BTC_USD:  "initialBTCPrice",
		exchange_config.MARKET_ETH_USD:  "initialETHPrice",
		exchange_config.MARKET_LINK_USD: "initialLINKPrice",
	}
	assertPricesWithTimeout(t, node, expectedPrices, 30*time.Second)
}

func TestUpgrade(t *testing.T) {
	testnet, err := NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]

	err = UpgradeTestnet(constants.AliceAccAddress.String(), t, node, version.CurrentVersion)
	require.NoError(t, err)
}
