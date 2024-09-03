package app_test

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app"
	perpetualsmoduletypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	consumertypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
	"github.com/stretchr/testify/require"
)

func TestModuleAccountsToAddresses(t *testing.T) {
	expectedModuleAccToAddresses := map[string]string{
		authtypes.FeeCollectorName: "dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2",

		ibctransfertypes.ModuleName:                "dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5",
		satypes.ModuleName:                         "dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6",
		perpetualsmoduletypes.InsuranceFundName:    "dydx1c7ptc87hkd54e3r7zjy92q29xkq7t79w64slrq",
		icatypes.ModuleName:                        "dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8",
		consumertypes.ConsumerRedistributeName:     "dydx1x69dz0c0emw8m2c6kp5v6c08kgjxmu30yn6p5y",
		consumertypes.ConsumerToSendToProviderName: "dydx1ywtansy6ss0jtq8ckrcv6jzkps8yh8mf37gcch",
		satypes.LiquidityFeeModuleAddress:          "dydx1l4fct6xefgds6tsslrluwy2juuyaet369u29e7",
	}

	require.True(t, len(expectedModuleAccToAddresses) == len(app.GetMaccPerms()))
	for acc, address := range expectedModuleAccToAddresses {
		expectedAddr := authtypes.NewModuleAddress(acc).String()
		require.Equal(t, address, expectedAddr, "module (%v) should have address (%s)", acc, expectedAddr)
	}
}

func TestBlockedAddresses(t *testing.T) {
	expectedBlockedAddresses := map[string]bool{
		"dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2": true,
		"dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5": true,
		"dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8": true,
		"dydx1x69dz0c0emw8m2c6kp5v6c08kgjxmu30yn6p5y": true,
		"dydx1ywtansy6ss0jtq8ckrcv6jzkps8yh8mf37gcch": true,
	}
	require.Equal(t, expectedBlockedAddresses, app.BlockedAddresses())
}

func TestMaccPerms(t *testing.T) {
	maccPerms := app.GetMaccPerms()
	expectedMaccPerms := map[string][]string{

		"fee_collector":            nil,
		"insurance_fund":           nil,
		"subaccounts":              nil,
		"transfer":                 {"minter", "burner"},
		"interchainaccounts":       nil,
		"cons_redistribute":        nil,
		"cons_to_send_to_provider": nil,
		"liquidity_module":         nil,
	}
	require.Equal(t, expectedMaccPerms, maccPerms, "default macc perms list does not match expected")
}

func TestModuleAccountAddrs(t *testing.T) {
	expectedModuleAccAddresses := map[string]bool{
		"dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2": true, // x/auth.FeeCollector
		"dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5": true, // ibc transfer
		"dydx1vlthgax23ca9syk7xgaz347xmf4nunefw3cnv8": true, // interchainaccounts
		"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6": true, // x/subaccount
		"dydx1c7ptc87hkd54e3r7zjy92q29xkq7t79w64slrq": true, // x/clob.insuranceFund
		"dydx1x69dz0c0emw8m2c6kp5v6c08kgjxmu30yn6p5y": true, // x/ccvconsumer.ConsumerRedistribute
		"dydx1ywtansy6ss0jtq8ckrcv6jzkps8yh8mf37gcch": true, // x/ccvconsumer.ConsumerToSendToProvider
		"dydx1l4fct6xefgds6tsslrluwy2juuyaet369u29e7": true, // x/subaccount.LiquidityFeeModuleAddress
	}

	require.Equal(t, expectedModuleAccAddresses, app.ModuleAccountAddrs())
}
