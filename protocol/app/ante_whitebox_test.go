package app

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/dydxprotocol/v4/app/basic_manager"
	libante "github.com/dydxprotocol/v4/lib/ante"
	"github.com/dydxprotocol/v4/lib/encoding"
	clobante "github.com/dydxprotocol/v4/x/clob/ante"
	clobmodulekeeper "github.com/dydxprotocol/v4/x/clob/keeper"
	clobmodulememclob "github.com/dydxprotocol/v4/x/clob/memclob"
	"github.com/stretchr/testify/require"
)

func newTestHandlerOptions() HandlerOptions {
	encodingConfig := encoding.MakeEncodingConfig(basic_manager.ModuleBasics)
	appCodec := encodingConfig.Codec
	txConfig := encodingConfig.TxConfig

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		nil,
		authtypes.ProtoBaseAccount,
		nil,
		sdk.Bech32MainPrefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		nil,
		accountKeeper,
		BlockedAddresses(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	feeGrantKeeper := feegrantkeeper.NewKeeper(appCodec, nil, accountKeeper)

	memClob := clobmodulememclob.NewMemClobPriceTimePriority(false)
	clobKeeper := clobmodulekeeper.NewKeeper(
		appCodec,
		nil,
		nil,
		nil,
		memClob,
		nil,
		nil,
		bankKeeper,
		nil,
		nil,
	)
	return HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   accountKeeper,
			BankKeeper:      bankKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			FeegrantKeeper:  feeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		ClobKeeper: clobKeeper,
	}
}

type WrappedHandlerType reflect.Type

func wrapDecoratorStr(decoratorStr string) string {
	return "(" + decoratorStr + ")"
}

func humanReadableDecoratorTypes(decoratorChain []sdk.AnteDecorator) []string {
	var dTypes []string

	for _, decorator := range decoratorChain {
		decoratorType := reflect.TypeOf(decorator).String()
		switch decorator := decorator.(type) {
		case libante.AppInjectedMsgAnteWrapper:
			switch nestedDecorator := decorator.GetAnteHandler().(type) {
			// The OffChainSingleMsgClobTxnAnteWrapper has an additional layer of nesting
			case clobante.OffChainSingleMsgClobTxAnteWrapper:
				dTypes = append(dTypes, decoratorType+
					wrapDecoratorStr(
						reflect.TypeOf(nestedDecorator).String()+
							wrapDecoratorStr(
								reflect.TypeOf(nestedDecorator.GetAnteHandler()).String(),
							),
					),
				)
			default:
				dTypes = append(dTypes, decoratorType+wrapDecoratorStr(reflect.TypeOf(nestedDecorator).String()))
			}
		default:
			dTypes = append(dTypes, decoratorType)
		}
	}
	return dTypes
}

func TestAnteHandlerChainOrder_Valid(t *testing.T) {
	handlerOptions := newTestHandlerOptions()
	decoratorChain := newAnteDecoratorChain(handlerOptions)
	decoratorTypes := humanReadableDecoratorTypes(decoratorChain)

	expectedDecoratorTypes := []string{
		"ante.AppInjectedMsgAnteWrapper(ante.OffChainSingleMsgClobTxAnteWrapper(ante.SetUpContextDecorator))",
		"ante.FreeInfiniteGasDecorator",
		"ante.RejectExtensionOptionsDecorator",
		"ante.ValidateMsgTypeDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.ValidateBasicDecorator)",
		"ante.TxTimeoutHeightDecorator",
		"ante.ValidateMemoDecorator",
		"ante.ConsumeTxSizeGasDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.OffChainSingleMsgClobTxAnteWrapper(ante.DeductFeeDecorator))",
		"ante.AppInjectedMsgAnteWrapper(ante.SetPubKeyDecorator)",
		"ante.ValidateSigCountDecorator",
		"ante.AppInjectedMsgAnteWrapper(ante.SigGasConsumeDecorator)",
		"ante.AppInjectedMsgAnteWrapper(ante.SigVerificationDecorator)",
		"ante.AppInjectedMsgAnteWrapper(ante.IncrementSequenceDecorator)",
		"ante.ClobDecorator",
	}

	require.Equal(t, expectedDecoratorTypes, decoratorTypes, "Decorator order does not match expected")
}
