//go:build all || integration_test

package cli_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	appconstants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/network"
	sa_testutil "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/client/testutil"
	satypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/subaccounts/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	// // Deterministic Mnemonic.
	validatorMnemonic := constants.AliceMnenomic

	// Generated from the above Mnemonic.
	s.validatorAddress = constants.AliceAccAddress
	fmt.Println("Validator address", s.validatorAddress)

	// Configure test network.
	s.cfg = network.DefaultConfig(nil)

	s.cfg.Mnemonics = append(s.cfg.Mnemonics, validatorMnemonic)
	s.cfg.ChainID = appconstants.AppName

	// Set min gas prices to zero so that we can submit transactions with zero gas price.
	s.cfg.MinGasPrices = fmt.Sprintf("0%s", sdk.DefaultBondDenom)

	genesisChanges := "\".app_state.subaccounts.subaccounts = [{\\\"id\\\": {\\\"number\\\": \\\"0\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"500000000\\\"}], \\\"perpetual_positions\\\": []}, {\\\"id\\\": {\\\"number\\\": \\\"1\\\", \\\"owner\\\": \\\"dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6\\\"}, \\\"asset_positions\\\": [{\\\"asset_id\\\": \\\"0\\\", \\\"quantums\\\": \\\"500000000\\\"}], \\\"perpetual_positions\\\": []}] | .app_state.epochs.epoch_info_list = [{\\\"name\\\": \\\"funding-sample\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"funding-tick\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}, {\\\"name\\\": \\\"stats-epoch\\\", \\\"next_tick\\\": \\\"1747543084\\\", \\\"duration\\\": \\\"31536000\\\", \\\"current_epoch\\\": \\\"0\\\", \\\"current_epoch_start_block\\\": \\\"0\\\", \\\"fast_forward_next_tick\\\": false}]\" \"\""
	network.DeployCustomNetwork(genesisChanges)

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
	network.CleanupCustomNetwork()
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
	network.CleanupCustomNetwork()
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
	network.CleanupCustomNetwork()
}

func (s *SendingIntegrationTestSuite) sendTransferAndVerifyBalance(
	senderSubaccountNumber uint32,
	recipientSubaccountNumber uint32,
	amount uint64,
	expectedSenderQuoteBalance *big.Int,
	expectedRecipientQuoteBalance *big.Int,
) {

	cfg := network.DefaultConfig(nil)
	transferTx := fmt.Sprintf("docker exec interchain-security-instance interchain-security-cd tx sending create-transfer dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 %d dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 %d %d --from dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6 --chain-id consu --home /consu/validatoralice --node tcp://7.7.8.4:26658 --keyring-backend test -y -o json", senderSubaccountNumber, recipientSubaccountNumber, amount)
	_, _, err := network.QueryCustomNetwork(transferTx)
	if err != nil {
		s.T().Fatalf("failed to send transfer: %v", err)
	}
	time.Sleep(5 * time.Second)

	// Query both subaccounts.
	resp, err := sa_testutil.MsgQuerySubaccountExec("dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6", senderSubaccountNumber)
	s.Require().NoError(err)

	var subaccountResp satypes.QuerySubaccountResponse
	s.Require().NoError(cfg.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
	sender := subaccountResp.Subaccount

	resp, err = sa_testutil.MsgQuerySubaccountExec("dydx1eeeggku6dzk3mv7wph3zq035rhtd890smfq5z6", recipientSubaccountNumber)
	s.Require().NoError(err)

	s.Require().NoError(cfg.Codec.UnmarshalJSON(resp.Bytes(), &subaccountResp))
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
