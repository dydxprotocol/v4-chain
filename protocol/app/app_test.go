package app_test

import (
	"github.com/dydxprotocol/v4-chain/protocol/mocks"
	"gopkg.in/typ.v4/slices"
	"reflect"
	"strings"
	"testing"
	"time"

	delaymsgmodule "github.com/dydxprotocol/v4-chain/protocol/x/delaymsg"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	evidencemodule "github.com/cosmos/cosmos-sdk/x/evidence"
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
	"github.com/dydxprotocol/v4-chain/protocol/app/basic_manager"
	"github.com/dydxprotocol/v4-chain/protocol/app/flags"
	custommodule "github.com/dydxprotocol/v4-chain/protocol/app/module"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	assetsmodule "github.com/dydxprotocol/v4-chain/protocol/x/assets"
	blocktimemodule "github.com/dydxprotocol/v4-chain/protocol/x/blocktime"
	bridgemodule "github.com/dydxprotocol/v4-chain/protocol/x/bridge"
	clobmodule "github.com/dydxprotocol/v4-chain/protocol/x/clob"
	epochsmodule "github.com/dydxprotocol/v4-chain/protocol/x/epochs"
	feetiersmodule "github.com/dydxprotocol/v4-chain/protocol/x/feetiers"
	perpetualsmodule "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals"
	pricesmodule "github.com/dydxprotocol/v4-chain/protocol/x/prices"
	rewardsmodule "github.com/dydxprotocol/v4-chain/protocol/x/rewards"
	sendingmodule "github.com/dydxprotocol/v4-chain/protocol/x/sending"
	statsmodule "github.com/dydxprotocol/v4-chain/protocol/x/stats"
	subaccountsmodule "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts"
	vestmodule "github.com/dydxprotocol/v4-chain/protocol/x/vest"

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
			tApp := testapp.NewTestAppBuilder(t).WithAppOptions(tc.customFlags).Build()
			tApp.InitChain()
			uninitializedFields := getUninitializedStructFields(reflect.ValueOf(*tApp.App))

			// Note that the daemon clients are currently hard coded as disabled in GetDefaultTestAppOptions.
			// Normally they would be only disabled for non-validating full nodes or for nodes where any
			// daemon is explicitly disabled.
			expectedUninitializedFields := []string{
				"PriceFeedClient",
				"LiquidationsClient",
				"BridgeClient",
			}
			for _, field := range expectedUninitializedFields {
				if idx := slices.Index(uninitializedFields, field); idx >= 0 {
					slices.Remove(&uninitializedFields, idx)
				}
			}

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

func TestAppPanicsWithGrpcDisabled(t *testing.T) {
	customFlags := map[string]interface{}{
		flags.GrpcEnable: false,
	}
	require.Panics(t, func() { testapp.DefaultTestApp(customFlags) })
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
		evidencemodule.AppModuleBasic{},
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
		delaymsgmodule.AppModuleBasic{},
		epochsmodule.AppModuleBasic{},
	)

	expectedFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(defaultAppModuleBasics))
	actualFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(basic_manager.ModuleBasics))
	require.Equal(t, expectedFieldTypes, actualFieldTypes, "Module basics does not match expected")
}

func TestMonitorDaemon_Panics(t *testing.T) {
	app := testapp.DefaultTestApp(nil)
	hc := &mocks.HealthCheckable{}
	hc.On("ServiceName").Return("test-service")
	hc.On("HealthCheck").Return(nil)

	app.MonitorDaemon(hc, 5*time.Minute)
	// The second registration should fail, causing a panic.
	require.Panics(t, func() { app.MonitorDaemon(hc, 5*time.Minute) })
}
