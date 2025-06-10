package authenticator_test

import (
	"os"
	"testing"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"

	"github.com/stretchr/testify/suite"
)

type SubaccountFilterTest struct {
	BaseAuthenticatorSuite

	SubaccountFilter authenticator.SubaccountFilter
}

func TestSubaccountFilterTest(t *testing.T) {
	suite.Run(t, new(SubaccountFilterTest))
}

func (s *SubaccountFilterTest) SetupTest() {
	s.SetupKeys()
	s.SubaccountFilter = authenticator.NewSubaccountFilter()
}

func (s *SubaccountFilterTest) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

// TestFilter tests the SubaccountFilter with multiple clob messages
func (s *SubaccountFilterTest) TestFilter() {
	tests := map[string]struct {
		whitelist string
		msg       sdk.Msg

		match bool
	}{
		"order place": {
			whitelist: "0",
			msg:       constants.Msg_PlaceOrder_LongTerm,
			match:     true,
		},
		"order cancel": {
			whitelist: "0",
			msg:       constants.Msg_CancelOrder_LongTerm,
			match:     true,
		},
		"order batch cancel": {
			whitelist: "0",
			msg:       constants.Msg_BatchCancel,
			match:     true,
		},
		"order place - fail": {
			whitelist: "1",
			msg:       constants.Msg_PlaceOrder_LongTerm,
			match:     false,
		},
		"order cancel - fail": {
			whitelist: "1",
			msg:       constants.Msg_CancelOrder_LongTerm,
			match:     false,
		},
		"order batch cancel - fail": {
			whitelist: "1",
			msg:       constants.Msg_BatchCancel,
			match:     false,
		},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			requireSigVerification, err := s.SubaccountFilter.OnAuthenticatorAdded(
				s.Ctx,
				sdk.AccAddress{},
				[]byte(tt.whitelist),
				"1",
			)
			s.Require().False(requireSigVerification)
			s.Require().NoError(err)
			filter, err := s.SubaccountFilter.Initialize([]byte(tt.whitelist))
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
				constants.AliceAccAddress,
				constants.AliceAccAddress,
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
				s.Require().ErrorIs(err, types.ErrSubaccountVerification)
			}
		})
	}
}
