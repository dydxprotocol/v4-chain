package ante

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
)

func IsUnsupportedMsg(msg sdk.Msg) bool {
	switch msg.(type) {
	case
		// ICA Controller messages
		*icacontrollertypes.MsgUpdateParams,
		*icacontrollertypes.MsgSendTx,
		*icacontrollertypes.MsgRegisterInterchainAccount,
		// ------- CosmosSDK default modules
		// gov
		*govbeta.MsgSubmitProposal,
		*gov.MsgCancelProposal,
		// ------- CosmWasm
		*wasmtypes.MsgClearAdmin,
		*wasmtypes.MsgIBCCloseChannel,
		*wasmtypes.MsgIBCSend,
		*wasmtypes.MsgPinCodes,
		*wasmtypes.MsgRemoveCodeUploadParamsAddresses,
		*wasmtypes.MsgSudoContract,
		*wasmtypes.MsgUnpinCodes,
		*wasmtypes.MsgUpdateAdmin,
		*wasmtypes.MsgUpdateContractLabel,
		*wasmtypes.MsgUpdateInstantiateConfig,
		*wasmtypes.MsgMigrateContract:
		return true
	}
	return false
}
