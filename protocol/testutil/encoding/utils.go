package encoding

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/ante"

	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/gogoproto/proto"
	"github.com/cosmos/ibc-go/modules/capability"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
	custommodule "github.com/dydxprotocol/v4-chain/protocol/app/module"
	bridgemodule "github.com/dydxprotocol/v4-chain/protocol/x/bridge"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers"
	perpetualtypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	subaccountsmodule "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts"
	vaultmodule "github.com/dydxprotocol/v4-chain/protocol/x/vault"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

// GetTestEncodingCfg returns an encoding config for testing purposes.
func GetTestEncodingCfg() testutil.TestEncodingConfig {
	encodingCfg := ante.MakeTestEncodingConfig(
		auth.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(
			[]govclient.ProposalHandler{
				paramsclient.ProposalHandler,
			},
		),
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
		bridgemodule.AppModuleBasic{},
		subaccountsmodule.AppModuleBasic{},
		vaultmodule.AppModuleBasic{},
	)

	msgInterfacesToRegister := []sdk.Msg{
		// Clob.
		&clobtypes.MsgProposedOperations{},
		&clobtypes.MsgPlaceOrder{},
		&clobtypes.MsgCancelOrder{},
		&clobtypes.MsgBatchCancel{},

		// Perpetuals.
		&perpetualtypes.MsgAddPremiumVotes{},

		// Prices.
		&pricestypes.MsgUpdateMarketPrices{},

		// Sending.
		&sendingtypes.MsgCreateTransfer{},
		&sendingtypes.MsgDepositToSubaccount{},
		&sendingtypes.MsgWithdrawFromSubaccount{},

		// Vault.
		&vaulttypes.MsgDepositToMegavault{},
		&vaulttypes.MsgWithdrawFromMegavault{},
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
func EncodeMessageToAny(t testing.TB, msg sdk.Msg) *codectypes.Any {
	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)
	return any
}
