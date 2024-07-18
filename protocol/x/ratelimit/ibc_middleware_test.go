package ratelimit_test

import (
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	"github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cometbft/cometbft/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb/testutil"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

var (
	globalStartTime = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	chainIDPrefix   = "localdydxprotocol"
)

type KeeperTestSuite struct {
	testifysuite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
	chainC *ibctesting.TestChain
}

func TestKeeperTestSuite(t *testing.T) {
	testifysuite.Run(t, new(KeeperTestSuite))
}

func createSignersByAddress(t *testing.T, val *ccvtypes.CrossChainValidator) (string, cmttypes.PrivValidator, *cmttypes.Validator) {

	consAddress := sdk.ConsAddress(val.Address)
	privKey := constants.GetPrivKeyFromConsAddress(consAddress)
	edPriv := ed25519.PrivKey(privKey.Bytes())
	priv := cmttypes.MockPV{
		PrivKey: edPriv,
	}

	pubKey, err := priv.GetPubKey()
	require.NoError(t, err)
	validator := cmttypes.NewValidator(pubKey, 500)

	return pubKey.Address().String(), priv, validator
}

func convertSimAccountsToSenderAccounts(simAccounts []simtypes.Account) []ibctesting.SenderAccount {
	senderAccounts := make([]ibctesting.SenderAccount, len(simAccounts))
	for i, simAccount := range simAccounts {
		baseAccount := auth.NewBaseAccount(simAccount.Address, simAccount.PubKey, uint64(i+5), 0)
		senderAccounts[i] = ibctesting.SenderAccount{
			SenderPrivKey: simAccount.PrivKey,
			SenderAccount: baseAccount,
		}
	}
	return senderAccounts
}

func setupChainForIBC(t *testing.T, coord *ibctesting.Coordinator, chainID string) *ibctesting.TestChain {
	t.Helper()
	r := testutil.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 10)
	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		genesis.ChainID = chainID // Update chain_id to chainID
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *auth.GenesisState) {
				for _, simAccount := range simAccounts {
					acct := &auth.BaseAccount{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						PubKey:  codectypes.UnsafePackAny(simAccount.PubKey),
					}
					genesisState.Accounts = append(genesisState.Accounts, codectypes.UnsafePackAny(acct))
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *banktypes.GenesisState) {
				for _, simAccount := range simAccounts {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						Coins: sdktypes.NewCoins(sdktypes.NewInt64Coin(
							sdk.DefaultBondDenom,
							constants.Usdc_Asset_500_000.Quantums.BigInt().Int64(),
						)),
					})
				}
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()
	ctx := tApp.InitChain()

	// create current header and call begin block
	header := cmtproto.Header{
		ChainID: chainID,
		Height:  2,
		Time:    coord.CurrentTime.UTC(),
	}

	txConfig := tApp.App.GetTxConfig()

	// convert ccv validators to standard validators
	vals := tApp.App.ConsumerKeeper.GetAllCCValidator(ctx)
	validators := make([]*cmttypes.Validator, len(vals))
	signers := make(map[string]cmttypes.PrivValidator, len(validators))
	for i, val := range vals {
		address, priv, validator := createSignersByAddress(t, &val)
		validators[i] = validator
		signers[address] = priv
	}
	valSet := cmttypes.NewValidatorSet(validators)

	senderAccounts := convertSimAccountsToSenderAccounts(simAccounts)

	// create an account to send transactions from
	chain := &ibctesting.TestChain{
		TB:             t,
		Coordinator:    coord,
		ChainID:        chainID,
		App:            tApp.App,
		CurrentHeader:  header,
		QueryServer:    tApp.App.GetIBCKeeper(),
		TxConfig:       txConfig,
		Codec:          tApp.App.AppCodec(),
		Vals:           valSet,
		NextVals:       valSet,
		Signers:        signers,
		SenderPrivKey:  senderAccounts[0].SenderPrivKey,
		SenderAccount:  senderAccounts[0].SenderAccount,
		SenderAccounts: senderAccounts,
	}

	// commit genesis block
	chain.NextBlock()

	return chain

}

func NewCoordinator(t *testing.T, n int) *ibctesting.Coordinator {
	t.Helper()
	chains := make(map[string]*ibctesting.TestChain)
	coord := &ibctesting.Coordinator{
		T:           t,
		CurrentTime: globalStartTime,
	}

	for i := 1; i <= n; i++ {
		chainID := GetChainID(i)
		chains[chainID] = setupChainForIBC(t, coord, chainID)
	}
	coord.Chains = chains

	return coord
}

func GetChainID(i int) string {
	return chainIDPrefix + "-" + strconv.Itoa(i)
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.coordinator = NewCoordinator(suite.T(), 3)
	suite.chainA = suite.coordinator.GetChain(GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.App.(*app.App).InterfaceRegistry())
	transfertypes.RegisterQueryServer(queryHelper, suite.chainA.App.(*app.App).TransferKeeper)
}

func (suite *KeeperTestSuite) TestSendTransfer() {
	var (
		coin            sdk.Coin
		path            *ibctesting.Path
		sender          sdk.AccAddress
		timeoutHeight   clienttypes.Height
		memo            string
		expEscrowAmount sdkmath.Int // total amount in escrow for denom on receiving chain
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"successful transfer with native token",
			func() {
				expEscrowAmount = sdkmath.NewInt(100)
			}, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path = ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			coin = sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100))
			sender = suite.chainA.SenderAccount.GetAddress()
			memo = ""
			timeoutHeight = suite.chainB.GetTimeoutHeight()
			expEscrowAmount = sdkmath.ZeroInt()

			//create IBC token on chainA
			transferMsg := transfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, coin, suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainA.GetTimeoutHeight(), 0, "")
			result, err := suite.chainB.SendMsgs(transferMsg)
			suite.Require().NoError(err) // message committed

			packet, err := ibctesting.ParsePacketFromEvents(result.Events)
			suite.Require().NoError(err)

			err = path.RelayPacket(packet)
			suite.Require().NoError(err)

			tc.malleate()

			msg := transfertypes.NewMsgTransfer(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				coin, sender.String(), suite.chainB.SenderAccount.GetAddress().String(),
				timeoutHeight, 0, // only use timeout height
				memo,
			)

			res, err := suite.chainA.App.(*app.App).TransferKeeper.Transfer(suite.chainA.GetContext(), msg)

			// check total amount in escrow of sent token denom on sending chain
			amount := suite.chainA.App.(*app.App).TransferKeeper.GetTotalEscrowForDenom(suite.chainA.GetContext(), coin.GetDenom())
			suite.Require().Equal(expEscrowAmount, amount.Amount)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}
