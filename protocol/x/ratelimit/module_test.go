package ratelimit_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	testapp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/dydxprotocol/v4-chain/protocol/x/ratelimit"
	ratelimit_keeper "github.com/dydxprotocol/v4-chain/protocol/x/ratelimit/keeper"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stretchr/testify/require"
)

func createAppModule(t *testing.T) ratelimit.AppModule {
	am, _, _ := createAppModuleWithKeeper(t)
	return am
}

// Returns the keeper and context along with the AppModule.
// This is useful for tests which want to write/read state
// to/from the keeper.
func createAppModuleWithKeeper(t *testing.T) (ratelimit.AppModule, *ratelimit_keeper.Keeper, sdk.Context) {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.RatelimitKeeper

	return ratelimit.NewAppModule(
		appCodec,
		k,
	), &k, ctx
}

func createAppModuleBasic(t *testing.T) ratelimit.AppModuleBasic {
	appCodec := codec.NewProtoCodec(module.InterfaceRegistry)

	appModule := ratelimit.NewAppModuleBasic(appCodec)
	require.NotNil(t, appModule)

	return appModule
}

func TestAppModule_Name(t *testing.T) {
	am := createAppModule(t)
	require.Equal(t, "ratelimit", am.Name())
}

func TestAppModuleBasic_RegisterCodecLegacyAmino(t *testing.T) {
	am := createAppModuleBasic(t)

	cdc := codec.NewLegacyAmino()
	am.RegisterLegacyAminoCodec(cdc)

	var buf bytes.Buffer
	err := cdc.Amino.PrintTypes(&buf)
	require.NoError(t, err)
}

func TestAppModuleBasic_ValidateGenesisErr(t *testing.T) {
	tests := map[string]struct {
		genesisJson string
		expectedErr string
	}{
		"Invalid json": {
			genesisJson: `{"missingClosingQuote: true}`,
			expectedErr: "failed to unmarshal ratelimit genesis state: unexpected EOF",
		},
		"Invalid state": {
			genesisJson: `{"foo":{}}`,
			expectedErr: "failed to unmarshal ratelimit genesis state: unknown field \"foo\" in types.GenesisState",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			am := createAppModuleBasic(t)

			cdc := codec.NewProtoCodec(module.InterfaceRegistry)

			err := am.ValidateGenesis(cdc, nil, json.RawMessage(tc.genesisJson))
			require.EqualError(t, err, tc.expectedErr)
		})
	}
}

func TestAppModuleBasic_RegisterGRPCGatewayRoutes(t *testing.T) {
	am := createAppModuleBasic(t)

	router := runtime.NewServeMux()

	am.RegisterGRPCGatewayRoutes(client.Context{}, router)

	// TODO(CORE-823): implement query for `x/ratelimit`
}

func TestAppModuleBasic_GetTxCmd(t *testing.T) {
	am := createAppModuleBasic(t)

	cmd := am.GetTxCmd()
	require.Equal(t, "ratelimit", cmd.Use)
	require.Equal(t, 0, len(cmd.Commands()))
}

func TestAppModuleBasic_GetQueryCmd(t *testing.T) {
	createAppModuleBasic(t)

	// TODO(CORE-823): implement query for `x/ratelimit`
}
