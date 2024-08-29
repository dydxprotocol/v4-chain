package app_test

import (
	"reflect"
	"strings"
	"testing"
	"time"

	evidencemodule "cosmossdk.io/x/evidence"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/flags"
	custommodule "github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/mocks"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	assetsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets"
	blocktimemodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/blocktime"
	clobmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob"
	delaymsgmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/delaymsg"
	epochsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/epochs"
	feetiersmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers"
	perpetualsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals"
	pricesmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices"
	ratelimitmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit"
	sendingmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending"
	statsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/stats"
	subaccountsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/ibc-go/modules/capability"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	consumer "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer"
	"github.com/stretchr/testify/require"
	"gopkg.in/typ.v4/slices"
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

			expectedUninitializedFields := []string{
				// Note that the daemon clients are currently hard coded as disabled in GetDefaultTestAppOptions.
				// Normally they would be only disabled for non-validating full nodes or for nodes where any
				// daemon is explicitly disabled.
				"PriceFeedClient",
				"DeleveragingClient",

				// Any default constructed type can be considered initialized if the default is what is
				// expected. getUninitializedStructFields relies on fields being the non-default and
				// reports the following as uninitialized.
				"event",
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
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		custommodule.SlashingModuleBasic{},
		evidencemodule.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ica.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		transfer.AppModuleBasic{},
		consensus.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		consumer.AppModuleBasic{},

		// Custom modules
		pricesmodule.AppModuleBasic{},
		assetsmodule.AppModuleBasic{},
		blocktimemodule.AppModuleBasic{},
		feetiersmodule.AppModuleBasic{},
		perpetualsmodule.AppModuleBasic{},
		statsmodule.AppModuleBasic{},
		subaccountsmodule.AppModuleBasic{},
		clobmodule.AppModuleBasic{},
		sendingmodule.AppModuleBasic{},
		delaymsgmodule.AppModuleBasic{},
		epochsmodule.AppModuleBasic{},
		ratelimitmodule.AppModuleBasic{},
	)

	app := testapp.DefaultTestApp(nil)

	expectedFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(defaultAppModuleBasics))
	actualFieldTypes := getMapFieldsAndTypes(reflect.ValueOf(app.ModuleBasics))
	require.Equal(t, expectedFieldTypes, actualFieldTypes, "Module basics does not match expected")
}

func TestRegisterDaemonWithHealthMonitor_Panics(t *testing.T) {
	app := testapp.DefaultTestApp(nil)
	hc := &mocks.HealthCheckable{}
	hc.On("ServiceName").Return("test-service")
	hc.On("HealthCheck").Return(nil)

	app.RegisterDaemonWithHealthMonitor(hc, 5*time.Minute)
	// The second registration should fail, causing a panic.
	require.PanicsWithError(
		t,
		"service test-service already registered",
		func() { app.RegisterDaemonWithHealthMonitor(hc, 5*time.Minute) },
	)
}
