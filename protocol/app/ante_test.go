package app_test

import (
	"cosmossdk.io/store/rootmulti"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	testApp "github.com/dydxprotocol/v4-chain/protocol/testutil/app"
	"github.com/stretchr/testify/require"
)

func newHandlerOptions() app.HandlerOptions {
	encodingConfig := app.GetEncodingConfig()
	dydxApp := testApp.DefaultTestApp(nil)
	return app.HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   dydxApp.AccountKeeper,
			BankKeeper:      dydxApp.BankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  dydxApp.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		ClobKeeper:       dydxApp.ClobKeeper,
		Codec:            encodingConfig.Codec,
		AuthStoreKey:     dydxApp.CommitMultiStore().(*rootmulti.Store).StoreKeysByName()[authtypes.StoreKey],
		PerpetualsKeeper: dydxApp.PerpetualsKeeper,
	}
}

func TestNewAnteHandler(t *testing.T) {
	handlerOptions := newHandlerOptions()
	anteHandler, err := app.NewAnteHandler(handlerOptions)
	require.NoError(t, err, "NewAnteHandler call failed")
	require.NotNil(t, anteHandler, "expected non-nil AnteHandler function")
}

func TestNewAnteHandler_Error(t *testing.T) {
	tests := map[string]struct {
		handlerMutation func(*app.HandlerOptions)
		errorMsg        string
	}{
		"nil handlerOptions.AccountKeeper": {
			handlerMutation: func(options *app.HandlerOptions) { options.AccountKeeper = nil },
			errorMsg:        "account keeper is required for ante builder",
		},
		"nil handlerOptions.BankKeeper": {
			handlerMutation: func(options *app.HandlerOptions) { options.BankKeeper = nil },
			errorMsg:        "bank keeper is required for ante builder",
		},
		"nil PerpetualsKeeper": {
			handlerMutation: func(options *app.HandlerOptions) { options.PerpetualsKeeper = nil },
			errorMsg:        "perpetuals keeper is required for ante builder",
		},
		"nil handlerOptions.SignModeHandler": {
			handlerMutation: func(options *app.HandlerOptions) { options.SignModeHandler = nil },
			errorMsg:        "sign mode handler is required for ante builder",
		},
		"nil ClobKeeper": {
			handlerMutation: func(options *app.HandlerOptions) { options.ClobKeeper = nil },
			errorMsg:        "clob keeper is required for ante builder",
		},
		"nil Codec": {
			handlerMutation: func(options *app.HandlerOptions) { options.Codec = nil },
			errorMsg:        "codec is required for ante builder",
		},
		"nil AuthStoreKey": {
			handlerMutation: func(options *app.HandlerOptions) { options.AuthStoreKey = nil },
			errorMsg:        "auth store key is required for ante builder",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handlerOptions := newHandlerOptions()
			tc.handlerMutation(&handlerOptions)

			anteHandler, err := app.NewAnteHandler(handlerOptions)
			require.Nil(t, anteHandler, "Expected Ante Handler creation to error")
			require.Errorf(t, err, tc.errorMsg)
		})
	}
}
