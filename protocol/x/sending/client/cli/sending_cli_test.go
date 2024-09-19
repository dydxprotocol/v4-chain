package cli_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	appconstants "github.com/dydxprotocol/v4-chain/protocol/app/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/network"
	epochstypes "github.com/dydxprotocol/v4-chain/protocol/x/epochs/types"
	testutil "github.com/dydxprotocol/v4-chain/protocol/x/sending/client/testutil"
	"github.com/dydxprotocol/v4-chain/protocol/x/sending/types"
	sa_testutil "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/client/testutil"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
	"github.com/stretchr/testify/suite"
)

var (
	subaccountNumberZero  = uint32(0)
	subaccountNumberOne   = uint32(1)
	subaccountNonExistent = uint32(127)
)

type SendingIntegrationTestSuite struct {
	suite.Suite

	validatorAddress sdk.AccAddress
	cfg              network.Config
	network          *network.Network
}

func TestSendingIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &SendingIntegrationTestSuite{})
}

func (s *SendingIntegrationTestSuite) SetupTest() {
	s.T().Log("setting up sending integration test")

	// Deterministic Mnemonic.
	validatorMnemonic := constants.AliceMnenomic

	// Generated from the above Mnemonic.
	s.validatorAddress = constants.AliceAccAddress

	// Configure test network.
	s.cfg = network.DefaultConfig(nil)

	s.cfg.Mnemonics = append(s.cfg.Mnemonics, validatorMnemonic)
	s.cfg.ChainID = appconstants.AppName

	// Set min gas prices to zero so that we can submit transactions with zero gas price.
	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	// Setting genesis state for Sending.
	state := types.GenesisState{}

	buf, err := s.cfg.Codec.MarshalJSON(&state)
	s.NoError(err)
	s.cfg.GenesisState[types.ModuleName] = buf

	// Setting genesis state for Subaccounts.
	// Two subaccounts with non-zero USDC balances are added to the genesis state,
	// so that we can initiate transfers from these subaccounts and observe the changes in
	// their USDC positions.
	sastate := satypes.GenesisState{}
	sastate.Subaccounts = append(
		sastate.Subaccounts,
		satypes.Subaccount{
			Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: subaccountNumberZero},
			AssetPositions: []*satypes.AssetPosition{
				&constants.Usdc_Asset_500,
			},
			PerpetualPositions: []*satypes.PerpetualPosition{},
		},
		satypes.Subaccount{
			Id: &satypes.SubaccountId{Owner: s.validatorAddress.String(), Number: subaccountNumberOne},
			AssetPositions: []*satypes.AssetPosition{
				&constants.Usdc_Asset_500,
			},
			PerpetualPositions: []*satypes.PerpetualPosition{},
		},
	)

	sabuf, err := s.cfg.Codec.MarshalJSON(&sastate)
	s.Require().NoError(err)
	s.cfg.GenesisState[satypes.ModuleName] = sabuf

	// Ensure that no funding-related epochs will occur during this test.
	epstate := constants.GenerateEpochGenesisStateWithoutFunding()

	epbuf, err := s.cfg.Codec.MarshalJSON(&epstate)
	s.Require().NoError(err)
	s.cfg.GenesisState[epochstypes.ModuleName] = epbuf

	s.network = network.New(s.T(), s.cfg)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

// TestCLISending_Success sends a transfer from one subaccount to another (with the same owner and different numbers).
// The account which sends the transfer is also the validator's AccAddress.
// The transfer is expected to succeed, and after the transfer, the subaccounts are queried and assertions
// are performed on their new QuoteBalance.
func (s *SendingIntegrationTestSuite) TestCLISending_Success() {
	s.sendTransferAndVerifyBalance(
		subaccountNumberZero,
		subaccountNumberOne,
		uint64(1_000_000),
		new(big.Int).SetUint64(499_000_000),
		new(big.Int).SetUint64(501_000_000),
	)
}

// TestCLISending_InsufficientBalance attempts to send a transfer from one subaccount to
// another (with the same owner and different numbers). The transfer amount is greater than the sender's current
// balance. The transfer is expected to fail, and the subaccounts are expected to have the same QuoteBalance that
// they started with.
func (s *SendingIntegrationTestSuite) TestCLISending_InsufficientBalance() {
	s.sendTransferAndVerifyBalance(
		subaccountNumberZero,
		subaccountNumberOne,
		uint64(501_000_000), // Sender only has $500
		new(big.Int).SetUint64(500_000_000),
		new(big.Int).SetUint64(500_000_000),
	)
}

// TestCLISending_Nonexistent sends a transfer from one subaccount to
// another (with the same owner and different numbers). The recipient subaccount does not exist in state.
// The transfer is expected to succeed, and after the transfer, the subaccounts are queried and assertions
// are performed on their new QuoteBalance.
func (s *SendingIntegrationTestSuite) TestCLISending_Nonexistent() {
	s.sendTransferAndVerifyBalance(
		subaccountNumberZero,
		subaccountNonExistent,
		uint64(1_000_000),
		new(big.Int).SetUint64(499_000_000),
		new(big.Int).SetUint64(1_000_000),
	)
}

func (s *SendingIntegrationTestSuite) sendTransferAndVerifyBalance(
	senderSubaccountNumber uint32,
	recipientSubaccountNumber uint32,
	amount uint64,
	expectedSenderQuoteBalance *big.Int,
	expectedRecipientQuoteBalance *big.Int,
) {
	val := s.network.Validators[0]
	ctx := val.ClientCtx

	// Send the transfer from sender to recipient.
	_, err := testutil.MsgCreateTransferExec(
		ctx,
		s.validatorAddress,
		senderSubaccountNumber,
		s.validatorAddress,
		recipientSubaccountNumber,
		amount,
	)
	s.Require().NoError(err)

	currentHeight, err := s.network.LatestHeight()
	s.Require().NoError(err)

	// Wait for a few blocks to ensure the transfer was complated.
	_, err = s.network.WaitForHeight(currentHeight + 3)
	s.Require().NoError(err)

	// Query both subaccounts.
	resp, err := sa_testutil.MsgQuerySubaccountExec(ctx, s.validatorAddress, senderSubaccountNumber)
	s.Require().NoError(err)

	var subaccountResp satypes.QuerySubaccountResponse
	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
	sender := subaccountResp.Subaccount

	resp, err = sa_testutil.MsgQuerySubaccountExec(ctx, s.validatorAddress, recipientSubaccountNumber)
	s.Require().NoError(err)

	s.Require().NoError(s.network.Config.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
	recipient := subaccountResp.Subaccount

	// Assert that both Subaccounts have the appropriate state.
	s.Require().Equal(
		expectedSenderQuoteBalance,
		sender.GetUsdcPosition(),
	)
	s.Require().Empty(sender.PerpetualPositions)

	s.Require().Equal(
		expectedRecipientQuoteBalance,
		recipient.GetUsdcPosition(),
	)
	s.Require().Empty(recipient.PerpetualPositions)
}
