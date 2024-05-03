package encoding

import (
	"testing"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/ante"

	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	custommodule "github.com/StreamFinance-Protocol/stream-chain/protocol/app/module"
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/feetiers"
	perpetualtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/types"
	sendingtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/sending/types"
	subaccountsmodule "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/ibc-go/modules/capability"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	"github.com/stretchr/testify/require"
)

// GetTestEncodingCfg returns an encoding config for testing purposes.
func GetTestEncodingCfg() testutil.TestEncodingConfig {
	encodingCfg := ante.MakeTestEncodingConfig(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		custommodule.SlashingModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		feetiers.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctm.AppModuleBasic{},
		ica.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		transfer.AppModuleBasic{},
		consensus.AppModuleBasic{},

		// Custom modules
		subaccountsmodule.AppModuleBasic{})

	msgInterfacesToRegister := []sdk.Msg{
		// Clob.
		&clobtypes.MsgProposedOperations{},
		&clobtypes.MsgPlaceOrder{},
		&clobtypes.MsgCancelOrder{},

		// Perpetuals.
		&perpetualtypes.MsgAddPremiumVotes{},

		// Prices.
		&pricestypes.MsgUpdateMarketPrices{},

		// Sending.
		&sendingtypes.MsgCreateTransfer{},
		&sendingtypes.MsgDepositToSubaccount{},
		&sendingtypes.MsgWithdrawFromSubaccount{},
	}

	for _, msg := range msgInterfacesToRegister {
		encodingCfg.InterfaceRegistry.RegisterInterface(
			"/"+proto.MessageName(msg),
			(*sdk.Msg)(nil),
			msg,
		)
	}

	return encodingCfg
}

// EncodeMessageToAny converts a message to an Any object for protobuf encoding.
func EncodeMessageToAny(t *testing.T, msg sdk.Msg) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	return any
}
