package types_test

import (
	fmt "fmt"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAuthenticator_ValidateBasic(t *testing.T) {
	const messageFilterData = "/cosmos.bank.v1beta1.MsgMultiSend,/cosmos.bank.v1beta1.MsgSend"

	tests := map[string]struct {
		msg         types.MsgAddAuthenticator
		expectedErr error
	}{
		"Success": {
			msg: types.MsgAddAuthenticator{
				Sender:            constants.AliceAccAddress.String(),
				AuthenticatorType: "MessageFilter",
				Data:              []byte(messageFilterData),
			},
		},
		"Failure: Not an address": {
			msg: types.MsgAddAuthenticator{
				Sender: "invalid",
			},
			expectedErr: types.ErrInvalidAccountAddress,
		},
		"Failure: Data exceeds max length": {
			msg: types.MsgAddAuthenticator{
				Sender:            constants.AliceAccAddress.String(),
				AuthenticatorType: "AllOf",
				Data: []byte(
					fmt.Sprintf(
						`[
							{"Type":"MessageFilter","Config":"%s"},
							{
								"Type":"AnyOf",
								"Config":"
									[
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"}
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"}
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"},
										{"Type":"SignatureVerification","Config":"%s"}
									]
								"
							},
						]`,
						messageFilterData,
						constants.AlicePrivateKey.PubKey().String(),
						constants.BobPrivateKey.PubKey().String(),
						constants.CarlPrivateKey.PubKey().String(),
						constants.DavePrivateKey.PubKey().String(),
						constants.AlicePrivateKey.PubKey().String(),
						constants.BobPrivateKey.PubKey().String(),
						constants.CarlPrivateKey.PubKey().String(),
						constants.DavePrivateKey.PubKey().String(),
						constants.AlicePrivateKey.PubKey().String(),
						constants.BobPrivateKey.PubKey().String(),
						constants.CarlPrivateKey.PubKey().String(),
						constants.DavePrivateKey.PubKey().String(),
					),
				),
			},
			expectedErr: types.ErrAuthenticatorDataExceedsMaximumLength,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
