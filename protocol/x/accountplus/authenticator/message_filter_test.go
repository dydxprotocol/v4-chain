package authenticator_test

import (
	"os"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"

	"github.com/stretchr/testify/suite"
)

type MessageFilterTest struct {
	BaseAuthenticatorSuite

	MessageFilter authenticator.MessageFilter
}

func TestMessageFilterTest(t *testing.T) {
	suite.Run(t, new(MessageFilterTest))
}

func (s *MessageFilterTest) SetupTest() {
	s.SetupKeys()
	s.MessageFilter = authenticator.NewMessageFilter()
}

func (s *MessageFilterTest) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

// TestBankSend tests the MessageFilter with multiple bank send messages
func (s *MessageFilterTest) TestBankSend() {
	tests := map[string]struct {
		msgType string
		msg     sdk.Msg

		match bool
	}{
		"bank send": {
			msgType: "/cosmos.bank.v1beta1.MsgSend",
			msg: &bank.MsgSend{
				FromAddress: s.TestAccAddress[0].String(),
				ToAddress:   "to",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin("foo", 100)),
			},
			match: true,
		},
		"bank send - multiple types": {
			msgType: "/cosmos.bank.v1beta1.MsgMultiSend,/cosmos.bank.v1beta1.MsgSend",
			msg: &bank.MsgSend{
				FromAddress: s.TestAccAddress[0].String(),
				ToAddress:   "to",
				Amount:      sdk.NewCoins(sdk.NewInt64Coin("foo", 100)),
			},
			match: true,
		},
		"bank send. fail on different message type": {
			msgType: "/cosmos.bank.v1beta1.MsgSend",
			msg: &clobtypes.MsgPlaceOrder{
				Order: clobtypes.Order{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner: s.TestAccAddress[0].String(),
						},
					},
				},
			},
			match: false,
		},
		"bank send - multiple types. fail on different message type": {
			msgType: "/cosmos.bank.v1beta1.MsgMultiSend,/cosmos.bank.v1beta1.MsgSend",
			msg: &clobtypes.MsgPlaceOrder{
				Order: clobtypes.Order{
					OrderId: clobtypes.OrderId{
						SubaccountId: satypes.SubaccountId{
							Owner: s.TestAccAddress[0].String(),
						},
					},
				},
			},
			match: false,
		},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			requireSigVerification, err := s.MessageFilter.OnAuthenticatorAdded(s.Ctx, sdk.AccAddress{}, []byte(tt.msgType), "1")
			s.Require().False(requireSigVerification)
			s.Require().NoError(err)
			filter, err := s.MessageFilter.Initialize([]byte(tt.msgType))
			s.Require().NoError(err)

			ak := s.tApp.App.AccountKeeper
			sigModeHandler := s.EncodingConfig.TxConfig.SignModeHandler()
			tx, err := s.GenSimpleTx([]sdk.Msg{tt.msg}, []cryptotypes.PrivKey{s.TestPrivKeys[0]})
			s.Require().NoError(err)
			request, err := lib.GenerateAuthenticationRequest(
				s.Ctx,
				s.tApp.App.AppCodec(),
				ak,
				sigModeHandler,
				s.TestAccAddress[0],
				s.TestAccAddress[0],
				nil,
				sdk.NewCoins(),
				tt.msg,
				tx,
				0,
				false,
			)
			s.Require().NoError(err)

			err = filter.Authenticate(s.Ctx, request)
			if tt.match {
				s.Require().NoError(err)
			} else {
				s.Require().ErrorIs(err, types.ErrMessageTypeVerification)
			}
		})
	}
}
