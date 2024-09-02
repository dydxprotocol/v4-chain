package ratelimit_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
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
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sDAIOracle"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

var (
	globalStartTime              = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	chainIDPrefix                = "localdydxprotocol"
	sDaiPoolAccountAddressString = "dydx1r3fsd6humm0ghyq0te5jf8eumklmclya37zle0"
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

func setupChainForIBC(
	t *testing.T,
	coord *ibctesting.Coordinator,
	chainID string,
	accountCoinDenom string,
	accountCoinAmount *big.Int,
) *ibctesting.TestChain {
	if accountCoinDenom == ratelimittypes.SDaiDenom {
		panic("Cannot use sDAI denom as the coin denom. Cannot deposit sDAI directly into user account")
	}

	t.Helper()
	r := testutil.NewRand()
	simAccounts := simtypes.RandomAccounts(r, 10)

	// sdai_amount / tdai_amount
	sDaiToTDaiConversionRate := sdaiservertypes.TestSDAIEventRequests[0].ConversionRate
	sDaiToTDaiConversionRateAsBigInt, found := new(big.Int).SetString(sDaiToTDaiConversionRate, 10)
	if !found {
		panic("Could not convert sdai to tdai conversion rate to big.Int")
	}

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
				// deposit equivalent sdai amount into the sdai pool account
				if accountCoinDenom == assettypes.TDaiDenom {
					totalUTDaiAmount := big.NewInt(0).Mul(accountCoinAmount, big.NewInt(int64(len(simAccounts))))
					tenScaledBySDaiDecimals := new(big.Int).Exp(
						big.NewInt(ratelimittypes.BASE_10),
						big.NewInt(ratelimittypes.SDAI_DECIMALS),
						nil,
					)
					scaledSDaiAmount := big.NewInt(0).Mul(totalUTDaiAmount, tenScaledBySDaiDecimals)
					sDaiAmount := scaledSDaiAmount.Div(scaledSDaiAmount, sDaiToTDaiConversionRateAsBigInt)

					// Perform denom exponent conversion
					conversionDecimals := new(big.Int).Abs(
						big.NewInt(ratelimittypes.SDaiDenomExponent - assettypes.TDaiDenomExponent),
					)
					tenScaledByConversionDecimals := new(big.Int).Exp(
						big.NewInt(10),
						conversionDecimals,
						nil,
					)
					gSDaiAmount := new(big.Int).Mul(sDaiAmount, tenScaledByConversionDecimals)
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: sDaiPoolAccountAddressString,
						Coins: sdktypes.NewCoins(sdktypes.NewCoin(
							ratelimittypes.SDaiDenom,
							math.NewIntFromBigInt(gSDaiAmount),
						)),
					})

				}

				for _, simAccount := range simAccounts {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						Coins: sdktypes.NewCoins(sdktypes.NewCoin(
							accountCoinDenom,
							math.NewIntFromBigInt(accountCoinAmount),
						)),
					})
				}
			},
		)
		return genesis
	}).WithNonDeterminismChecksEnabled(false).Build()
	ctx := tApp.InitChain()

	if accountCoinDenom == assettypes.TDaiDenom {
		k := tApp.App.RatelimitKeeper
		k.SetSDAIPrice(ctx, sDaiToTDaiConversionRateAsBigInt)
		denomTrace := transfertypes.DenomTrace{
			Path:      ratelimittypes.SDaiBaseDenomPathPrefix,
			BaseDenom: ratelimittypes.SDaiBaseDenom,
		}
		tApp.App.TransferKeeper.SetDenomTrace(
			ctx,
			denomTrace,
		)
	}

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

func NewCoordinator(t *testing.T, n int, accountCoinDenom string, accountCoinAmount *big.Int) *ibctesting.Coordinator {
	t.Helper()
	chains := make(map[string]*ibctesting.TestChain)
	coord := &ibctesting.Coordinator{
		T:           t,
		CurrentTime: globalStartTime,
	}

	for i := 1; i <= n; i++ {
		chainID := GetChainID(i)
		chains[chainID] = setupChainForIBC(t, coord, chainID, accountCoinDenom, accountCoinAmount)
	}
	coord.Chains = chains

	return coord
}

func GetChainID(i int) string {
	return chainIDPrefix + "-" + strconv.Itoa(i)
}

func (suite *KeeperTestSuite) SetupTest(accountCoinDenom string, accountCoinAmount *big.Int) {
	suite.coordinator = NewCoordinator(suite.T(), 3, accountCoinDenom, accountCoinAmount)
	suite.chainA = suite.coordinator.GetChain(GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(GetChainID(3))

	queryHelper := baseapp.NewQueryServerTestHelper(suite.chainA.GetContext(), suite.chainA.App.(*app.App).InterfaceRegistry())
	transfertypes.RegisterQueryServer(queryHelper, suite.chainA.App.(*app.App).TransferKeeper)
}

func (suite *KeeperTestSuite) setupSDaiDenomTrace() {
	chains := []*ibctesting.TestChain{suite.chainA, suite.chainB, suite.chainC}
	sDaiToTDaiConversionRate := sdaiservertypes.TestSDAIEventRequests[0].ConversionRate
	sDaiToTDaiConversionRateAsBigInt, found := new(big.Int).SetString(sDaiToTDaiConversionRate, 10)
	if !found {
		panic("Could not convert sdai to tdai conversion rate to big.Int")
	}

	for _, chain := range chains {
		chainApp := chain.App.(*app.App)
		chainApp.TransferKeeper.SetDenomTrace(
			chain.GetContext(),
			transfertypes.DenomTrace{
				Path:      ratelimittypes.SDaiBaseDenomPathPrefix,
				BaseDenom: ratelimittypes.SDaiBaseDenom,
			},
		)

		chainApp.RatelimitKeeper.SetSDAIPrice(
			chain.GetContext(),
			sDaiToTDaiConversionRateAsBigInt,
		)
	}
}

func (suite *KeeperTestSuite) TestSendTransfer() {
	var (
		coin                  sdk.Coin
		path                  *ibctesting.Path
		sender                sdk.AccAddress
		timeoutHeight         clienttypes.Height
		isEscrow              bool     // if false, then we expect token to be burned
		accountCoinAmountSent *big.Int // amount of coins in account that were sent
		memo                  string
	)

	testCases := []struct {
		name              string
		accountCoinDenom  string
		accountCoinAmount *big.Int
		sendCoinDenom     string
		sendCoinAmount    math.Int
		additionalSetup   func()
		malleate          func()
		expPass           bool
		expEarlyErr       bool
	}{
		{
			name:              "successful transfer with native token",
			accountCoinDenom:  sdk.DefaultBondDenom,
			accountCoinAmount: constants.TDai_Asset_500_000.Quantums.BigInt(),
			sendCoinDenom:     sdk.DefaultBondDenom,
			sendCoinAmount:    sdkmath.NewInt(100),
			additionalSetup:   func() {},
			malleate: func() {
				accountCoinAmountSent = big.NewInt(100)
				isEscrow = true
			},
			expPass:     true,
			expEarlyErr: false,
		},
		{
			name:              "successful transfer with tDAI and sDAI: basic",
			accountCoinDenom:  assettypes.TDaiDenom,
			accountCoinAmount: constants.TDai_Asset_500_000.Quantums.BigInt(), // Note: represents amount of utdai
			sendCoinDenom:     ratelimittypes.SDaiDenom,
			sendCoinAmount:    sdkmath.NewInt(100), // Note: represents amount of gsdai
			additionalSetup:   suite.setupSDaiDenomTrace,
			malleate: func() {
				accountCoinAmountSent = big.NewInt(1)
				isEscrow = false
			},
			expPass:     true,
			expEarlyErr: false,
		},
		{
			name:              "successful transfer with tDAI and sDAI: real scenario",
			accountCoinDenom:  assettypes.TDaiDenom,
			accountCoinAmount: constants.TDai_Asset_500_000.Quantums.BigInt(), // Note: represents amount of utdai
			sendCoinDenom:     ratelimittypes.SDaiDenom,
			sendCoinAmount: func() sdkmath.Int {
				bi, _ := new(big.Int).SetString("496681580107906596", 10)
				return sdkmath.NewIntFromBigInt(bi)
			}(), // Note: represents amount of gsdai
			additionalSetup: suite.setupSDaiDenomTrace,
			malleate: func() {
				accountCoinAmountSent = big.NewInt(500000)
				isEscrow = false
			},
			expPass:     true,
			expEarlyErr: false,
		},
		{
			name:              "successful transfer with tDAI and sDAI: account sends its entire tDAI balance",
			accountCoinDenom:  assettypes.TDaiDenom,
			accountCoinAmount: big.NewInt(500000), // Note: represents amount of utdai
			sendCoinDenom:     ratelimittypes.SDaiDenom,
			sendCoinAmount: func() sdkmath.Int {
				bi, _ := new(big.Int).SetString("496681580107906596", 10)
				return sdkmath.NewIntFromBigInt(bi)
			}(), // Note: represents amount of gsdai
			additionalSetup: suite.setupSDaiDenomTrace,
			malleate: func() {
				accountCoinAmountSent = big.NewInt(500000)
				isEscrow = false
			},
			expPass:     true,
			expEarlyErr: false,
		},
		{
			name:              "failed transfer with tDAI and sDAI: utdai balance too low",
			accountCoinDenom:  assettypes.TDaiDenom,
			accountCoinAmount: big.NewInt(1), // Note: represents amount of utdai
			sendCoinDenom:     ratelimittypes.SDaiDenom,
			sendCoinAmount: func() sdkmath.Int {
				bi, _ := new(big.Int).SetString("496681580107906596", 10)
				return sdkmath.NewIntFromBigInt(bi)
			}(), // Note: represents amount of gsdai
			additionalSetup: suite.setupSDaiDenomTrace,
			malleate:        func() {},
			expPass:         false,
			expEarlyErr:     true,
		},
		{
			name:              "failed transfer with tDAI and sDAI: sending 0 amount",
			accountCoinDenom:  assettypes.TDaiDenom,
			accountCoinAmount: big.NewInt(10000), // Note: represents amount of utdai
			sendCoinDenom:     ratelimittypes.SDaiDenom,
			sendCoinAmount:    sdkmath.NewInt(0), // Note: represents amount of gsdai
			additionalSetup:   suite.setupSDaiDenomTrace,
			malleate:          func() {},
			expPass:           false,
			expEarlyErr:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest(tc.accountCoinDenom, tc.accountCoinAmount)

			tc.additionalSetup()

			path = ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			coin = sdk.NewCoin(tc.sendCoinDenom, tc.sendCoinAmount)
			sender = suite.chainA.SenderAccount.GetAddress()
			memo = ""
			timeoutHeight = suite.chainB.GetTimeoutHeight()

			//create IBC token on chainA
			transferMsg := transfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, coin, suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainA.GetTimeoutHeight(), 0, "")
			result, err := suite.chainB.SendMsgs(transferMsg)
			if tc.expEarlyErr {
				fmt.Println("ERR IS ", err)
				suite.Require().Error(err)
				return
			}
			suite.Require().NoError(err) // message committed

			packet, err := ibctesting.ParsePacketFromEvents(result.Events)
			suite.Require().NoError(err)

			err = path.RelayPacket(packet)
			suite.Require().NoError(err)

			tc.malleate()

			initialSupply := suite.chainA.App.(*app.App).BankKeeper.GetSupply(suite.chainA.GetContext(), coin.GetDenom()).Amount

			msg := transfertypes.NewMsgTransfer(
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				coin, sender.String(), suite.chainB.SenderAccount.GetAddress().String(),
				timeoutHeight, 0, // only use timeout height
				memo,
			)

			res, err := suite.chainA.App.(*app.App).TransferKeeper.Transfer(suite.chainA.GetContext(), msg)

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			} else {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				supplyRemaining := suite.chainA.App.(*app.App).BankKeeper.GetSupply(suite.chainA.GetContext(), coin.GetDenom())

				// Check whether appropriate amount of coins has been escrowed / burned
				if isEscrow {
					// When a token is escrowed, we still expect it to be part of the total denom supply
					suite.Require().Equal(initialSupply, supplyRemaining.Amount)
					amount := suite.chainA.App.(*app.App).TransferKeeper.GetTotalEscrowForDenom(suite.chainA.GetContext(), coin.GetDenom())
					suite.Require().Equal(tc.sendCoinAmount, amount.Amount)
				} else {
					deltaAmount := initialSupply.Sub(supplyRemaining.Amount)
					suite.Require().Equal(tc.sendCoinAmount, deltaAmount)
				}

				// Check that account does not hold the sent amount anymore
				accountBalance := suite.chainA.App.(*app.App).BankKeeper.GetBalance(suite.chainA.GetContext(), sender, tc.accountCoinDenom)
				actualAccountCoinAmountSent := tc.accountCoinAmount.Sub(tc.accountCoinAmount, accountBalance.Amount.BigInt())
				suite.Require().Equal(accountCoinAmountSent, actualAccountCoinAmountSent)
			}

			// if tc.expPass {
			// 	suite.Require().NoError(err)
			// 	suite.Require().NotNil(res)
			// } else {
			// 	suite.Require().Error(err)
			// 	suite.Require().Nil(res)
			// // }

			// supplyRemaining := suite.chainA.App.(*app.App).BankKeeper.GetSupply(suite.chainA.GetContext(), coin.GetDenom())

			// // Check whether appropriate amount of coins has been escrowed / burned
			// if isEscrow {
			// 	// When a token is escrowed, we still expect it to be part of the total denom supply
			// 	suite.Require().Equal(initialSupply, supplyRemaining.Amount)
			// 	amount := suite.chainA.App.(*app.App).TransferKeeper.GetTotalEscrowForDenom(suite.chainA.GetContext(), coin.GetDenom())
			// 	suite.Require().Equal(tc.sendCoinAmount, amount.Amount)
			// } else {
			// 	deltaAmount := initialSupply.Sub(supplyRemaining.Amount)
			// 	suite.Require().Equal(tc.sendCoinAmount, deltaAmount)
			// }

			// // Check that account does not hold the sent amount anymore
			// accountBalance := suite.chainA.App.(*app.App).BankKeeper.GetBalance(suite.chainA.GetContext(), sender, tc.accountCoinDenom)
			// actualAccountCoinAmountSent := tc.accountCoinAmount.Sub(tc.accountCoinAmount, accountBalance.Amount.BigInt())
			// suite.Require().Equal(accountCoinAmountSent, actualAccountCoinAmountSent)

			// if tc.expPass {
			// 	suite.Require().NoError(err)
			// 	suite.Require().NotNil(res)
			// } else {
			// 	suite.Require().Error(err)
			// 	suite.Require().Nil(res)
			// }
		})
	}
}
