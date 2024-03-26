//go:build all || container_test

package v_5_0_0_test

import (
	"testing"

	"github.com/cosmos/gogoproto/proto"
	v_5_0_0 "github.com/dydxprotocol/v4-chain/protocol/app/upgrades/v5.0.0"
	containertest "github.com/dydxprotocol/v4-chain/protocol/testing/containertest"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	perpetuals "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStateUpgrade(t *testing.T) {
	testnet, err := containertest.NewTestnetWithPreupgradeGenesis()
	require.NoError(t, err, "failed to create testnet - is docker daemon running?")
	err = testnet.Start()
	require.NoError(t, err)
	defer testnet.MustCleanUp()
	node := testnet.Nodes["alice"]
	nodeAddress := constants.AliceAccAddress.String()

	preUpgradeChecks(node, t)

	err = containertest.UpgradeTestnet(nodeAddress, t, node, v_5_0_0.UpgradeName)
	require.NoError(t, err)

	postUpgradeChecks(node, t)
}

func preUpgradeChecks(node *containertest.Node, t *testing.T) {
	preUpgradeCheckPerpetualMarketType(node, t)
	// Add test for your upgrade handler logic below
}

func postUpgradeChecks(node *containertest.Node, t *testing.T) {
	postUpgradecheckPerpetualMarketType(node, t)
	postUpgradeCheckLiquidityTiers(node, t)
	// Add test for your upgrade handler logic below
}

func preUpgradeCheckPerpetualMarketType(node *containertest.Node, t *testing.T) {
	perpetualsList := &perpetuals.QueryAllPerpetualsResponse{}
	resp, err := containertest.Query(
		node,
		perpetuals.NewQueryClient,
		perpetuals.QueryClient.AllPerpetuals,
		&perpetuals.QueryAllPerpetualsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), perpetualsList)
	require.NoError(t, err)
	for _, perpetual := range perpetualsList.Perpetual {
		assert.Equal(t, perpetuals.PerpetualMarketType_PERPETUAL_MARKET_TYPE_UNSPECIFIED, perpetual.Params.MarketType)
	}
}

func postUpgradecheckPerpetualMarketType(node *containertest.Node, t *testing.T) {
	perpetualsList := &perpetuals.QueryAllPerpetualsResponse{}
	resp, err := containertest.Query(
		node,
		perpetuals.NewQueryClient,
		perpetuals.QueryClient.AllPerpetuals,
		&perpetuals.QueryAllPerpetualsRequest{},
	)
	require.NoError(t, err)
	err = proto.UnmarshalText(resp.String(), perpetualsList)
	require.NoError(t, err)
	for _, perpetual := range perpetualsList.Perpetual {
		assert.Equal(t, perpetuals.PerpetualMarketType_PERPETUAL_MARKET_TYPE_CROSS, perpetual.Params.MarketType)
	}
}

func postUpgradeCheckLiquidityTiers(node *containertest.Node, t *testing.T) {
	resp, err := containertest.Query(
		node,
		perpetuals.NewQueryClient,
		perpetuals.QueryClient.AllLiquidityTiers,
		&perpetuals.QueryAllLiquidityTiersRequest{},
	)
	require.NoError(t, err)

	liquidityTiersResponse := &perpetuals.QueryAllLiquidityTiersResponse{}
	err = proto.UnmarshalText(resp.String(), liquidityTiersResponse)
	require.NoError(t, err)

	assert.Equal(t, 4, len(liquidityTiersResponse.LiquidityTiers))
	assert.Equal(t, perpetuals.LiquidityTier{
		Id:                     0,
		Name:                   "Large-Cap",
		InitialMarginPpm:       50_000,
		MaintenanceFractionPpm: 600_000,
		ImpactNotional:         10_000_000_000,
		OpenInterestLowerCap:   uint64(0),
		OpenInterestUpperCap:   uint64(0),
	}, liquidityTiersResponse.LiquidityTiers[0])

	assert.Equal(t, perpetuals.LiquidityTier{
		Id:                     1,
		Name:                   "Mid-Cap",
		InitialMarginPpm:       100_000,
		MaintenanceFractionPpm: 500_000,
		ImpactNotional:         5_000_000_000,
		OpenInterestLowerCap:   uint64(25_000_000_000_000),
		OpenInterestUpperCap:   uint64(50_000_000_000_000),
	}, liquidityTiersResponse.LiquidityTiers[1])

	assert.Equal(t, perpetuals.LiquidityTier{
		Id:                     2,
		Name:                   "Long-Tail",
		InitialMarginPpm:       200_000,
		MaintenanceFractionPpm: 500_000,
		ImpactNotional:         2_500_000_000,
		OpenInterestLowerCap:   uint64(10_000_000_000_000),
		OpenInterestUpperCap:   uint64(20_000_000_000_000),
	}, liquidityTiersResponse.LiquidityTiers[2])

	assert.Equal(t, perpetuals.LiquidityTier{
		Id:                     3,
		Name:                   "Safety",
		InitialMarginPpm:       1_000_000,
		MaintenanceFractionPpm: 200_000,
		ImpactNotional:         2_500_000_000,
		OpenInterestLowerCap:   uint64(500_000_000_000),
		OpenInterestUpperCap:   uint64(1_000_000_000_000),
	}, liquidityTiersResponse.LiquidityTiers[3])
}
