package encoding

import (
	"github.com/cosmos/gogoproto/proto"

	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	upgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	clobtypes "github.com/dydxprotocol/v4/x/clob/types"
	perpetualtypes "github.com/dydxprotocol/v4/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4/x/prices/types"
	sendingtypes "github.com/dydxprotocol/v4/x/sending/types"
)

// GetTestEncodingCfg returns an encoding config for testing purposes.
func GetTestEncodingCfg() simappparams.EncodingConfig {
	encodingCfg := simappparams.MakeTestEncodingConfig()

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

		// Bank.
		&banktypes.MsgSend{},

		// Gov.
		&gov.MsgSubmitProposal{},
		&govbeta.MsgSubmitProposal{},

		// Upgrade.
		&upgrade.MsgSoftwareUpgrade{},
		&upgrade.MsgCancelUpgrade{},
	}

	for _, msg := range msgInterfacesToRegister {
		encodingCfg.InterfaceRegistry.RegisterInterface(
			"/"+proto.MessageName(msg),
			(*sdk.Msg)(nil),
			msg,
		)
	}

	std.RegisterInterfaces(encodingCfg.InterfaceRegistry)

	return encodingCfg
}
