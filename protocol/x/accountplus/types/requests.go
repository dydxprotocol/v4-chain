package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// These structs define the data structure for authentication, used with AuthenticationRequest struct.
//

// SignModeData represents the signing modes with direct bytes and textual representation.
type SignModeData struct {
	Direct  []byte `json:"sign_mode_direct"`
	Textual string `json:"sign_mode_textual"`
}

// SimplifiedSignatureData contains lists of signers and their corresponding signatures.
type SimplifiedSignatureData struct {
	Signers    []sdk.AccAddress `json:"signers"`
	Signatures [][]byte         `json:"signatures"`
}

// ExplicitTxData encapsulates key transaction data like chain ID, account info, and messages.
type ExplicitTxData struct {
	ChainID         string    `json:"chain_id"`
	AccountNumber   uint64    `json:"account_number"`
	AccountSequence uint64    `json:"sequence"`
	TimeoutHeight   uint64    `json:"timeout_height"`
	Msgs            []sdk.Msg `json:"msgs"`
	Memo            string    `json:"memo"`
}

type TrackRequest struct {
	AuthenticatorId     string         `json:"authenticator_id"`
	Account             sdk.AccAddress `json:"account"`
	FeePayer            sdk.AccAddress `json:"fee_payer"`
	FeeGranter          sdk.AccAddress `json:"fee_granter,omitempty"`
	Fee                 sdk.Coins      `json:"fee"`
	Msg                 sdk.Msg        `json:"msg"`
	MsgIndex            uint64         `json:"msg_index"`
	AuthenticatorParams []byte         `json:"authenticator_params,omitempty"`
}

type ConfirmExecutionRequest struct {
	AuthenticatorId     string         `json:"authenticator_id"`
	Account             sdk.AccAddress `json:"account"`
	FeePayer            sdk.AccAddress `json:"fee_payer"`
	FeeGranter          sdk.AccAddress `json:"fee_granter,omitempty"`
	Fee                 sdk.Coins      `json:"fee"`
	Msg                 sdk.Msg        `json:"msg"`
	MsgIndex            uint64         `json:"msg_index"`
	AuthenticatorParams []byte         `json:"authenticator_params,omitempty"`
}

type AuthenticationRequest struct {
	AuthenticatorId string         `json:"authenticator_id"`
	Account         sdk.AccAddress `json:"account"`
	FeePayer        sdk.AccAddress `json:"fee_payer"`
	FeeGranter      sdk.AccAddress `json:"fee_granter,omitempty"`
	Fee             sdk.Coins      `json:"fee"`
	Msg             sdk.Msg        `json:"msg"`

	// Since array size is int, and size depends on the system architecture,
	// we use uint64 to cover all available architectures.
	// It is unsigned, so at this point, it can't be negative.
	MsgIndex uint64 `json:"msg_index"`

	// Only allowing messages with a single signer, so the signature can be a single byte array.
	Signature           []byte                  `json:"signature"`
	SignModeTxData      SignModeData            `json:"sign_mode_tx_data"`
	TxData              ExplicitTxData          `json:"tx_data"`
	SignatureData       SimplifiedSignatureData `json:"signature_data"`
	Simulate            bool                    `json:"simulate"`
	AuthenticatorParams []byte                  `json:"authenticator_params,omitempty"`
}
