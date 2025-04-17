package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	marketmapmoduletypes "github.com/dydxprotocol/slinky/x/marketmap/types"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	bridgemoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	perpetualsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	rewardsmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/rewards/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	vaultmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	vestmoduletypes "github.com/dydxprotocol/v4-chain/protocol/x/vest/types"
)

func TestModuleAccountsToAddresses(t *testing.T) {
	expectedModuleAccToAddresses := map[string]string{
		authtypes.FeeCollectorName:                   "dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2",
		bridgemoduletypes.ModuleName:                 "dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv",
		distrtypes.ModuleName:                        "dydx1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wx2cfg",
		stakingtypes.BondedPoolName:                  "dydx1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uz8teq",
		stakingtypes.NotBondedPoolName:               "dydx1tygms3xhhs3yv487phx3dw4a95jn7t7lgzm605",
		govtypes.ModuleName:                          "dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky",
		ibctransfertypes.ModuleName:                  "dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5",
		satypes.ModuleName:                           "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6",
		perpetualsmoduletypes.InsuranceFundName:      "dydx1c7ptc87hkd54e3r7zjy92q29xkq7t79w64slrq",
		rewardsmoduletypes.TreasuryAccountName:       "dydx16wrau2x4tsg033xfrrdpae6kxfn9kyuerr5jjp",
		rewardsmoduletypes.VesterAccountName:         "dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp",
		vestmoduletypes.CommunityTreasuryAccountName: "dydx15ztc7xy42tn2ukkc0qjthkucw9ac63pgp70urn",
		vestmoduletypes.CommunityVesterAccountName:   "dydx1wxje320an3karyc6mjw4zghs300dmrjkwn7xtk",
		icatypes.ModuleName:                          "dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8",
		marketmapmoduletypes.ModuleName:              "dydx16j3d86dww8p2rzdlqsv7wle98cxzjxw6gjjyzn",
		vaultmoduletypes.MegavaultAccountName:        "dydx18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxqje4r",
	}

	require.True(t, len(expectedModuleAccToAddresses) == len(app.GetMaccPerms()),
		"expected %d, got %d", len(expectedModuleAccToAddresses), len(app.GetMaccPerms()))
	for acc, address := range expectedModuleAccToAddresses {
		expectedAddr := authtypes.NewModuleAddress(acc).String()
		require.Equal(t, address, expectedAddr, "module (%v) should have address (%s)", acc, expectedAddr)
	}
}

func TestBlockedAddresses(t *testing.T) {
	expectedBlockedAddresses := map[string]bool{
		"dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2": true,
		"dydx1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wx2cfg": true,
		"dydx1tygms3xhhs3yv487phx3dw4a95jn7t7lgzm605": true,
		"dydx1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uz8teq": true,
		"dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5": true,
		"dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8": true,
	}
	require.Equal(t, expectedBlockedAddresses, app.BlockedAddresses())
}

func TestMaccPerms(t *testing.T) {
	maccPerms := app.GetMaccPerms()
	expectedMaccPerms := map[string][]string{
		"bonded_tokens_pool":     {"burner", "staking"},
		"bridge":                 {"minter"},
		"distribution":           nil,
		"fee_collector":          nil,
		"gov":                    {"burner"},
		"insurance_fund":         nil,
		"not_bonded_tokens_pool": {"burner", "staking"},
		"subaccounts":            nil,
		"transfer":               {"minter", "burner"},
		"interchainaccounts":     nil,
		"rewards_treasury":       nil,
		"rewards_vester":         nil,
		"community_treasury":     nil,
		"community_vester":       nil,
		"marketmap":              nil,
		"megavault":              nil,
	}
	require.Equal(t, expectedMaccPerms, maccPerms, "default macc perms list does not match expected")
}

func TestModuleAccountAddrs(t *testing.T) {
	expectedModuleAccAddresses := map[string]bool{
		"dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2": true, // x/auth.FeeCollector
		"dydx1zlefkpe3g0vvm9a4h0jf9000lmqutlh9jwjnsv": true, // x/bridge
		"dydx1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wx2cfg": true, // x/distribution
		"dydx1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uz8teq": true, // x/staking.bondedPool
		"dydx1tygms3xhhs3yv487phx3dw4a95jn7t7lgzm605": true, // x/staking.notBondedPool
		"dydx10d07y265gmmuvt4z0w9aw880jnsr700jnmapky": true, // x/ gov
		"dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5": true, // ibc transfer
		"dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8": true, // interchainaccounts
		"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6": true, // x/subaccount
		"dydx1c7ptc87hkd54e3r7zjy92q29xkq7t79w64slrq": true, // x/clob.insuranceFund
		"dydx16wrau2x4tsg033xfrrdpae6kxfn9kyuerr5jjp": true, // x/rewards.treasury
		"dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp": true, // x/rewards.vester
		"dydx15ztc7xy42tn2ukkc0qjthkucw9ac63pgp70urn": true, // x/vest.communityTreasury
		"dydx1wxje320an3karyc6mjw4zghs300dmrjkwn7xtk": true, // x/vest.communityVester
		"dydx16j3d86dww8p2rzdlqsv7wle98cxzjxw6gjjyzn": true, // x/marketmap
		"dydx18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxqje4r": true, // x/vault.megavault
	}

	require.Equal(t, expectedModuleAccAddresses, app.ModuleAccountAddrs())
}
