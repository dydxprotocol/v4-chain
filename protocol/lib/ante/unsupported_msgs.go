package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govbeta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
	vaulttypes "github.com/dydxprotocol/v4-chain/protocol/x/vault/types"
)

// IsUnsupportedMsg returns true if the msg is unsupported by the app.
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
		// ------- dYdX custom modules
		// vault
		// nolint:staticcheck
		*vaulttypes.MsgSetVaultQuotingParams,
		// nolint:staticcheck
		*vaulttypes.MsgUpdateParams,
		*crisis.MsgVerifyInvariant:
		return true
	}
	return false
}
