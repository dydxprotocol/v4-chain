package app_test

import (
	"reflect"
	"strings"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclientclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/dydxprotocol/v4/app"
	"github.com/dydxprotocol/v4/app/basic_manager"
	"github.com/dydxprotocol/v4/app/flags"
	custommodule "github.com/dydxprotocol/v4/app/module"
	testapp "github.com/dydxprotocol/v4/testutil/app"
	assetsmodule "github.com/dydxprotocol/v4/x/assets"
	blocktimemodule "github.com/dydxprotocol/v4/x/blocktime"
	bridgemodule "github.com/dydxprotocol/v4/x/bridge"
	clobmodule "github.com/dydxprotocol/v4/x/clob"
	epochsmodule "github.com/dydxprotocol/v4/x/epochs"
	feetiersmodule "github.com/dydxprotocol/v4/x/feetiers"
	perpetualsmodule "github.com/dydxprotocol/v4/x/perpetuals"
	pricesmodule "github.com/dydxprotocol/v4/x/prices"
	rewardsmodule "github.com/dydxprotocol/v4/x/rewards"
	sendingmodule "github.com/dydxprotocol/v4/x/sending"
	statsmodule "github.com/dydxprotocol/v4/x/stats"
	subaccountsmodule "github.com/dydxprotocol/v4/x/subaccounts"
	vestmodule "github.com/dydxprotocol/v4/x/vest"

	"github.com/stretchr/testify/require"
)

func getUninitializedStructFields(reflectedStruct reflect.Value) []string {
	var uninitializedFields []string

	for i := 0; i < reflectedStruct.NumField(); i++ {
		field := reflectedStruct.Field(i)
		if field.IsZero() {
			uninitializedFields = append(uninitializedFields, reflectedStruct.Type().Field(i).Name)
		}
	}
	return uninitializedFields
}

func getMapFieldsAndTypes(reflectedMap reflect.Value) map[string]reflect.Type {
	fieldTypes := map[string]reflect.Type{}
	for _, key := range reflectedMap.MapKeys() {
		keyName := key.String()
		fieldTypes[keyName] = reflectedMap.MapIndex(key).Type()
	}
	return fieldTypes
}

func TestAppIsFullyInitialized(t *testing.T) {
	tests := map[string]struct {
		customFlags map[string]interface{}
	}{
		"default app": {
			customFlags: map[string]interface{}{},
		},
		"nonvalidating node app": {
			customFlags: map[string]interface{}{
				flags.NonValidatingFullNodeFlag: true,
			},
		},
		"validating node app": {
			customFlags: map[string]interface{}{
				flags.NonValidatingFullNodeFlag: false,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dydxApp := testapp.DefaultTestApp(tc.customFlags)
			uninitializedFields := getUninitializedStructFields(reflect.ValueOf(*dydxApp))
			require.Len(
				t,
				uninitializedFields,
				0,
				"The following top-level App fields were unset: %s",
				strings.Join(uninitializedFields, ", "),
			)
		})
	}
}

func TestClobKeeperMemStoreHasBeenInitialized(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	ctx := dydxApp.NewUncachedContext(true, tmproto.Header{})

	// The memstore panics if initialized twice so initializing again outside of application
	// start-up should cause a panic.
	require.Panics(t, func() { dydxApp.ClobKeeper.InitMemStore(ctx) })
}

func TestBaseApp(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.GetBaseApp(), "Expected non-nil BaseApp")
}

func TestLegacyAmino(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.LegacyAmino(), "Expected non-nil LegacyAmino")
}

func TestAppCodec(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.AppCodec(), "Expected non-nil AppCodec")
}

func TestInterfaceRegistry(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.InterfaceRegistry(), "Expected non-nil InterfaceRegistry")
}

func TestTxConfig(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.TxConfig(), "Expected non-nil TxConfig")
}

func TestDefaultGenesis(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.NotNil(t, dydxApp.DefaultGenesis(), "Expected non-nil DefaultGenesis")
}

func TestSimulationManager(t *testing.T) {
	dydxApp := testapp.DefaultTestApp(nil)
	require.Nil(t, dydxApp.SimulationManager(), "Expected nil SimulationManager")
}

func TestUpgrades(t *testing.T) {
	require.Len(t, app.Upgrades, 0, "Expected no Upgrades")
}

func TestForks(t *testing.T) {
	require.Len(t, app.Forks, 0, "Expected no Forks")
}

func TestBlockedAddresses(t *testing.T) {
	blockedAddresses := app.BlockedAddresses()
	expectedBlockedAddresses := map[string]bool{
		"dydx17xpfvakm2amg962yls6f84z3kell8c5leqdyt2": true,
		"dydx1c7ptc87hkd54e3r7zjy92q29xkq7t79w64slrq": true,
		"dydx1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3uz8teq": true,
		"dydx1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wx2cfg": true,
		"dydx1tygms3xhhs3yv487phx3dw4a95jn7t7lgzm605": true,
		"dydx1v88c3xv9xyv3eetdx0tvcmq7ung3dywp5upwc6": true,
		"dydx1yl6hdjhmkf37639730gffanpzndzdpmh8xcdh5": true,
		// `rewards_treasury` module account
		"dydx16wrau2x4tsg033xfrrdpae6kxfn9kyuerr5jjp": true,
		// `vester_treasury` module accoount
		"dydx1ltyc6y4skclzafvpznpt2qjwmfwgsndp458rmp": true,
	}
	require.Equal(t, expectedBlockedAddresses, blockedAddresses, "default blocked address list does not match expected")
}

func TestMaccPerms(t *testing.T) {
	maccPerms := app.GetMaccPerms()
	expectedMaccPerms := map[string][]string{
		"bonded_tokens_pool":     {"burner", "staking"},
		"distribution":           []string(nil),
		"fee_collector":          []string(nil),
		"gov":                    {"burner"},
		"insurance_fund":         []string(nil),
		"not_bonded_tokens_pool": {"burner", "staking"},
		"subaccounts":            []string(nil),
		"transfer":               {"minter", "burner"},
		"rewards_treasury":       nil,
		"rewards_vester":         nil,
	}
	require.Equal(t, expectedMaccPerms, maccPerms, "default macc perms list does not match expected")
}

func TestModuleBasics(t *testing.T) {
	defaultAppModuleBasics := module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
				upgradeclient.LegacyProposalHandler,
				upgradeclient.LegacyCancelProposalHandler,
				ibcclientclient.UpdateClientProposalHandler,
				ibcclientclient.UpgradeProposalHandler,
			},
		),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		custommodule.SlashingModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		transfer.AppModuleBasic{},
		consensus.AppModuleBasic{},

		// Custom modules
		pricesmodule.AppModuleBasic{},
		assetsmodule.AppModuleBasic{},
		blocktimemodule.AppModuleBasic{},
		bridgemodule.AppModuleBasic{},
		feetiersmodule.AppModuleBasic{},
		perpetualsmodule.AppModuleBasic{},
		statsmodule.AppModuleBasic{},
		subaccountsmodule.AppModuleBasic{},
		clobmodule.AppModuleBasic{},
		vestmodule.AppModuleBasic{},
		rewardsmodule.AppModuleBasic{},
		sendingmodule.AppModuleBasic{},
		epochsmodule.AppModuleBasic{},
	)

	expectedFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(defaultAppModuleBasics))
	actualFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(basic_manager.ModuleBasics))
	require.Equal(t, expectedFieldTypes, actualFieldTypes, "Module basics does not match expected")
}
