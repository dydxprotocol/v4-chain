package authenticator_test

import (
	"os"
	"testing"

	sdkmath "cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"

	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
	sendingtypes "github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/suite"
)

type ClobPairIdFilterTest struct {
	BaseAuthenticatorSuite

	ClobPairIdFilter authenticator.ClobPairIdFilter
}

func TestClobPairIdFilterTest(t *testing.T) {
	suite.Run(t, new(ClobPairIdFilterTest))
}

func (s *ClobPairIdFilterTest) SetupTest() {
	s.SetupKeys()
	s.ClobPairIdFilter = authenticator.NewClobPairIdFilter()
}

func (s *ClobPairIdFilterTest) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

// TestFilter tests the ClobPairIdFilter with multiple clob messages
func (s *ClobPairIdFilterTest) TestFilter() {
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
		// Regression: non-CLOB messages must be rejected, not silently allowed.
		// See git history for the original fail-open vulnerability.
		"withdraw from subaccount - rejected": {
			whitelist: "0",
			msg: &sendingtypes.MsgWithdrawFromSubaccount{
				Sender:    satypes.SubaccountId{Owner: constants.AliceAccAddress.String(), Number: 0},
				Recipient: constants.BobAccAddress.String(),
				AssetId:   0,
				Quantums:  1_000_000_000,
			},
			match: false,
		},
		"create transfer - rejected": {
			whitelist: "0",
			msg: &sendingtypes.MsgCreateTransfer{
				Transfer: &sendingtypes.Transfer{
					Sender:    satypes.SubaccountId{Owner: constants.AliceAccAddress.String(), Number: 0},
					Recipient: satypes.SubaccountId{Owner: constants.BobAccAddress.String(), Number: 0},
					AssetId:   0,
					Amount:    1_000_000,
				},
			},
			match: false,
		},
		"bank send - rejected": {
			whitelist: "0",
			msg: &banktypes.MsgSend{
				FromAddress: constants.AliceAccAddress.String(),
				ToAddress:   constants.BobAccAddress.String(),
				Amount:      sdk.NewCoins(sdk.NewCoin("uusdc", sdkmath.NewInt(1_000_000))),
			},
			match: false,
		},
	}

	for name, tt := range tests {
		s.Run(name, func() {
			requireSigVerification, err := s.ClobPairIdFilter.OnAuthenticatorAdded(
				s.Ctx,
				sdk.AccAddress{},
				[]byte(tt.whitelist),
				"1",
			)
			s.Require().NoError(err)
			s.Require().False(requireSigVerification)
			filter, err := s.ClobPairIdFilter.Initialize([]byte(tt.whitelist))
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
				s.Require().ErrorIs(err, types.ErrClobPairIdVerification)
			}
		})
	}
}
