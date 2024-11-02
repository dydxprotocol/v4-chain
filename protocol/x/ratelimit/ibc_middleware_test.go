package ratelimit_test

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/dtypes"
	testapp "github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/app"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	cometbfted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	"github.com/stretchr/testify/require"
	testifysuite "github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb/testutil"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app"
	sdaiservertypes "github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/server/types/sdaioracle"
	assettypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/assets/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/keeper"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var (
	globalStartTime              = time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)
	chainIDPrefix                = "localklyraprotocol"
	sDaiPoolAccountAddressString = "klyra1r3fsd6humm0ghyq0te5jf8eumklmclyaw0hs3y"
	exampleAckError              = errors.New("ABCI code: 1: error handling packet: see events for details")
	sequenceNumberOne            = uint64(1)
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

// func createSignersByAddressSecp256k1(t *testing.T, val *stakingtypes.Validator) (string, cmttypes.PrivValidator, *cmttypes.Validator) {
// 	originalPrivKey := constants.GetPrivKeyFromValidatorAddressString(val.OperatorAddress)

// 	// Ensure the original private key is of type Secp256k1
// 	secpPrivKey, ok := originalPrivKey.(*secp256k1.PrivKey)
// 	if !ok {
// 		panic(fmt.Sprintf("Unexpected private key type: %T", originalPrivKey))
// 	}

// 	// Use the Secp256k1 private key directly
// 	priv := cmttypes.MockPV{
// 		PrivKey: secp256k1comet.PrivKey(secpPrivKey.Bytes()),
// 	}

// 	pubKey, err := priv.GetPubKey()
// 	require.NoError(t, err)
// 	validator := cmttypes.NewValidator(pubKey, 500000)

// 	return pubKey.Address().String(), priv, validator
// }

func createSignersByAddressEddsa(t *testing.T, val *stakingtypes.Validator) (string, cmttypes.PrivValidator, *cmttypes.Validator) {

	// Retrieve the EdDSA private key from a map or function
	originalPrivKey := constants.GetEddsaPrivKeyFromValidatorAddressString(val.OperatorAddress)

	// Ensure the public key is of the correct type
	pubKey := originalPrivKey.PubKey()
	ed25519PubKey, ok := pubKey.(cometbfted25519.PubKey)
	if !ok {
		panic(fmt.Sprintf("Unexpected public key type: %T", pubKey))
	}

	// Create a validator using the public key
	validator := cmttypes.NewValidator(ed25519PubKey, 500000)

	// Wrap the private key if necessary
	priv := cmttypes.MockPV{
		PrivKey: originalPrivKey,
	}

	pubKey, err := priv.GetPubKey()
	require.NoError(t, err)

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
	sDaiToTDaiConversionRate := sdaiservertypes.TestSDAIEventRequest.ConversionRate
	sDaiToTDaiConversionRateAsBigInt, found := new(big.Int).SetString(sDaiToTDaiConversionRate, 10)
	if !found {
		panic("Could not convert sdai to tdai conversion rate to big.Int")
	}

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis cmttypes.GenesisDoc) {
		// for IBC testing, disable vote extensions
		if chainID == "localklyraprotocol-1" {
			genesis = testapp.NoVeGenesisIBCTest1()
		} else if chainID == "localklyraprotocol-2" {
			genesis = testapp.NoVeGenesisIBCTest2()
		} else if chainID == "localklyraprotocol-3" {
			genesis = testapp.NoVeGenesisIBCTest3()
		} else {
			genesis = testapp.NoVeGenesis()
		}
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

					var existingBalance *banktypes.Balance
					for i, balance := range genesisState.Balances {
						if balance.Address == sDaiPoolAccountAddressString {
							existingBalance = &genesisState.Balances[i]
							break
						}
					}

					if existingBalance == nil {
						// If no existing balance, append new balance
						genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
							Address: sDaiPoolAccountAddressString,
							Coins: sdktypes.NewCoins(sdktypes.NewCoin(
								ratelimittypes.SDaiDenom,
								sdkmath.NewIntFromBigInt(gSDaiAmount),
							)),
						})
					} else {
						// If existing balance found, increment the amount
						existingCoin := existingBalance.Coins.AmountOf(ratelimittypes.SDaiDenom)
						newAmount := existingCoin.Add(sdkmath.NewIntFromBigInt(gSDaiAmount))
						existingBalance.Coins = sdktypes.NewCoins(sdktypes.NewCoin(
							ratelimittypes.SDaiDenom,
							newAmount,
						))
					}
				}

				for _, simAccount := range simAccounts {
					genesisState.Balances = append(genesisState.Balances, banktypes.Balance{
						Address: sdktypes.AccAddress(simAccount.PubKey.Address()).String(),
						Coins: sdktypes.NewCoins(sdktypes.NewCoin(
							accountCoinDenom,
							sdkmath.NewIntFromBigInt(accountCoinAmount),
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

	// convert cosmos sdk validators to cmt validators
	vals, err := tApp.App.StakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		panic(err)
	}
	validators := make([]*cmttypes.Validator, len(vals))
	signers := make(map[string]cmttypes.PrivValidator, len(validators))
	for i, val := range vals {
		address, priv, validator := createSignersByAddressEddsa(t, &val)
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
	sDaiToTDaiConversionRate := sdaiservertypes.TestSDAIEventRequest.ConversionRate
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
		coin                  sdktypes.Coin
		path                  *ibctesting.Path
		sender                sdktypes.AccAddress
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
		sendCoinAmount    sdkmath.Int
		additionalSetup   func()
		malleate          func()
		expPass           bool
		expEarlyErr       bool
	}{
		{
			name:              "successful transfer with native token",
			accountCoinDenom:  sdktypes.DefaultBondDenom,
			accountCoinAmount: constants.TDai_Asset_500_000.Quantums.BigInt(),
			sendCoinDenom:     sdktypes.DefaultBondDenom,
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

			coin = sdktypes.NewCoin(tc.sendCoinDenom, tc.sendCoinAmount)
			sender = suite.chainA.SenderAccount.GetAddress()
			memo = ""
			timeoutHeight = suite.chainB.GetTimeoutHeight()

			//create IBC token on chainA
			transferMsg := transfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, coin, suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(), suite.chainA.GetTimeoutHeight(), 0, "")
			result, err := suite.chainB.SendMsgs(transferMsg)
			if tc.expEarlyErr {
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
		})
	}
}

func (suite *KeeperTestSuite) TestOnRecvPacket() {
	var (
		path       *ibctesting.Path
		packetData transfertypes.FungibleTokenPacketData
	)

	testCases := []struct {
		name             string
		denom            string
		denomAmountSent  string
		tDaiAmountMinted string
		conversionRate   string
		malleate         func()
		expPass          bool
		expCapacityIncr  bool
		expMintTDAI      bool
	}{
		{
			name:             "successful receive of non-sDAI token",
			denom:            sdktypes.DefaultBondDenom,
			denomAmountSent:  "100",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      false,
		},
		{
			name:             "successful receive of sDAI token",
			denom:            ratelimittypes.SDaiBaseDenom,
			denomAmountSent:  "1000000000000",
			tDaiAmountMinted: "1",
			conversionRate:   "1000000000000000000000000000",
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      true,
		},
		{
			name:             "successful receive of larger amount of sDAI token",
			denom:            ratelimittypes.SDaiBaseDenom,
			denomAmountSent:  "1000000000000000000",
			tDaiAmountMinted: "1000000",
			conversionRate:   "1000000000000000000000000000",
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      true,
		},
		{
			name:             "failed receive due to invalid packet data",
			denom:            sdktypes.DefaultBondDenom,
			denomAmountSent:  "100",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			malleate: func() {
				packetData.Amount = "invalid"
			},
			expPass:         false,
			expCapacityIncr: false,
			expMintTDAI:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest(sdktypes.DefaultBondDenom, big.NewInt(1000000))

			path = ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			sender := suite.chainB.SenderAccount.GetAddress().String()
			receiver := suite.chainA.SenderAccount.GetAddress().String()

			packetData = transfertypes.FungibleTokenPacketData{
				Denom:    tc.denom,
				Amount:   tc.denomAmountSent,
				Sender:   sender,
				Receiver: receiver,
			}

			app := suite.chainA.App.(*app.App)
			ctx := suite.chainA.GetContext()

			app.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr(tc.conversionRate))

			tc.malleate()

			// Record initial capacities and balances
			var initialCapacityList []dtypes.SerializableInt
			var initialBalance sdktypes.Coin
			if tc.expPass {
				initialCapacityList = app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom).CapacityList
				initialBalance = app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Receiver), assettypes.TDaiDenom)
			}

			module, _, err := app.GetIBCKeeper().PortKeeper.LookupModuleByPort(ctx, ibctesting.TransferPort)
			suite.Require().NoError(err)

			cbs, ok := app.GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			packetForTest := channeltypes.Packet{
				SourceChannel:      path.EndpointB.ChannelID,
				SourcePort:         path.EndpointB.ChannelConfig.PortID,
				DestinationChannel: path.EndpointA.ChannelID,
				DestinationPort:    path.EndpointA.ChannelConfig.PortID,
				Data:               packetData.GetBytes(),
			}

			ack := cbs.OnRecvPacket(ctx, packetForTest, suite.chainB.SenderAccount.GetAddress())

			if !tc.expPass {
				suite.Require().False(ack.Success())
			} else {
				suite.Require().True(ack.Success())

				// TODO: Can also include checks for other capacities not changing.
				newCapacity := app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom)
				if tc.expCapacityIncr {
					for i, capacity := range newCapacity.CapacityList {
						expectedCapacity := new(big.Int).Add(initialCapacityList[i].BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(packetData.Amount))
						suite.Require().Equal(expectedCapacity, capacity)
					}
				} else {
					suite.Require().Equal(initialCapacityList, newCapacity.CapacityList)
				}

				newBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Receiver), assettypes.TDaiDenom)
				if tc.expMintTDAI {
					expectedBalance := new(big.Int).Add(initialBalance.Amount.BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(tc.tDaiAmountMinted))
					suite.Require().Equal(expectedBalance, newBalance.Amount.BigInt())
				} else {
					suite.Require().Equal(initialBalance.Amount, newBalance.Amount)
				}

				// User should not have sDAI in their account
				sDaiDenomBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Receiver), ratelimittypes.SDaiDenom)
				suite.Require().Equal(0, sDaiDenomBalance.Amount.BigInt().Cmp(big.NewInt(0)))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnAcknowledgementPacket() {
	var (
		path       *ibctesting.Path
		packetData transfertypes.FungibleTokenPacketData
	)

	testCases := []struct {
		name             string
		denom            string
		denomAmountSent  string
		tDaiAmountMinted string
		conversionRate   string
		ack              channeltypes.Acknowledgement
		malleate         func()
		expPass          bool
		expCapacityIncr  bool
		expMintTDAI      bool
	}{
		{
			name:             "Success: handles ack success in case of non-sDai",
			denom:            sdktypes.DefaultBondDenom,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			ack:              channeltypes.NewResultAcknowledgement([]byte{1}),
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  false,
			expMintTDAI:      false,
		},
		{
			name:             "Success: handles ack failure in case of non-sDai",
			denom:            sdktypes.DefaultBondDenom,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			ack:              channeltypes.NewErrorAcknowledgement(exampleAckError),
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      false,
		},
		{
			name:             "Success: handles ack success in case of sDai",
			denom:            ratelimittypes.SDaiBaseDenomFullPath,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			ack:              channeltypes.NewResultAcknowledgement([]byte{1}),
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  false,
			expMintTDAI:      true,
		},
		{
			name:             "Success: handles ack failure in case of sDai",
			denom:            ratelimittypes.SDaiBaseDenomFullPath,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "10",
			conversionRate:   "1000000000000000000000000000",
			ack:              channeltypes.NewErrorAcknowledgement(exampleAckError),
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest(sdktypes.DefaultBondDenom, big.NewInt(1000000))

			path = ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			sender := suite.chainA.SenderAccount.GetAddress().String()
			receiver := suite.chainB.SenderAccount.GetAddress().String()

			packetData = transfertypes.FungibleTokenPacketData{
				Denom:    tc.denom,
				Amount:   tc.denomAmountSent,
				Sender:   sender,
				Receiver: receiver,
			}

			app := suite.chainA.App.(*app.App)
			ctx := suite.chainA.GetContext()

			app.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr(tc.conversionRate))

			module, _, err := app.GetIBCKeeper().PortKeeper.LookupModuleByPort(ctx, ibctesting.TransferPort)
			suite.Require().NoError(err)

			cbs, ok := app.GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			tc.malleate()

			// Record initial capacities and balances
			var initialCapacityList []dtypes.SerializableInt
			var initialSenderBalance sdktypes.Coin
			if tc.expPass {
				initialCapacityList = app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom).CapacityList
				initialSenderBalance = app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), assettypes.TDaiDenom)
			}

			packetForTest := channeltypes.Packet{
				DestinationChannel: path.EndpointB.ChannelID,
				DestinationPort:    path.EndpointB.ChannelConfig.PortID,
				SourceChannel:      path.EndpointA.ChannelID,
				SourcePort:         path.EndpointA.ChannelConfig.PortID,
				Data:               packetData.GetBytes(),
				Sequence:           sequenceNumberOne,
			}

			// If capacity increases, this means the ack was an error. This means we will unescrow tokens during
			// processing. Simulate previous escrowing here.
			denomToMint := packetData.Denom
			if packetData.Denom == ratelimittypes.SDaiBaseDenomFullPath {
				denomToMint = ratelimittypes.SDaiDenom
			}

			escrowedCoin := sdktypes.NewCoins(sdktypes.NewCoin(denomToMint, sdkmath.NewIntFromBigInt(keeper.ConvertStringToBigIntWithPanicOnErr(tc.denomAmountSent))))
			err = app.BankKeeper.MintCoins(ctx, module, escrowedCoin)
			suite.Require().NoError(err)
			escrowAddress := transfertypes.GetEscrowAddress(packetForTest.SourcePort, packetForTest.SourceChannel)
			err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, module, escrowAddress, escrowedCoin)
			suite.Require().NoError(err)

			app.TransferKeeper.SetTotalEscrowForDenom(ctx, escrowedCoin[0])

			app.RatelimitKeeper.SetPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence)

			// Setup Complete. Run test.
			ackBz, err := transfertypes.ModuleCdc.MarshalJSON(&tc.ack)
			suite.Require().NoError(err, "no error expected when marshalling ack")

			err = cbs.OnAcknowledgementPacket(ctx, packetForTest, ackBz, suite.chainB.SenderAccount.GetAddress())

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().True(app.RatelimitKeeper.HasPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence))
			} else {
				suite.Require().NoError(err)
				suite.Require().False(app.RatelimitKeeper.HasPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence))

				// TODO: Can also include checks for other capacities not changing.
				newCapacity := app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom)
				if tc.expCapacityIncr {
					for i, capacity := range newCapacity.CapacityList {
						expectedCapacity := new(big.Int).Add(initialCapacityList[i].BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(packetData.Amount))
						suite.Require().Equal(expectedCapacity, capacity)
					}
				} else {
					suite.Require().Equal(initialCapacityList, newCapacity.CapacityList)
				}

				newSenderBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), assettypes.TDaiDenom)
				if tc.expMintTDAI {
					expectedBalance := new(big.Int).Add(initialSenderBalance.Amount.BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(tc.tDaiAmountMinted))
					suite.Require().Equal(expectedBalance, newSenderBalance.Amount.BigInt())
				} else {
					suite.Require().Equal(newSenderBalance.Amount, newSenderBalance.Amount)
				}

				// user should not have sDAI at the end
				sDaiDenomBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), ratelimittypes.SDaiDenom)
				suite.Require().Equal(0, sDaiDenomBalance.Amount.BigInt().Cmp(big.NewInt(0)))
			}
		})
	}
}

func (suite *KeeperTestSuite) TestOnTimeoutPacket() {
	var (
		path       *ibctesting.Path
		packetData transfertypes.FungibleTokenPacketData
	)

	testCases := []struct {
		name             string
		denom            string
		denomAmountSent  string
		tDaiAmountMinted string
		conversionRate   string
		malleate         func()
		expPass          bool
		expCapacityIncr  bool
		expMintTDAI      bool
	}{
		{
			name:             "Success: handles non-sDai",
			denom:            sdktypes.DefaultBondDenom,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "0",
			conversionRate:   "1000000000000000000000000000",
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      false,
		},
		{
			name:             "Success: handles sDai",
			denom:            ratelimittypes.SDaiBaseDenomFullPath,
			denomAmountSent:  "10000000000000",
			tDaiAmountMinted: "10",
			conversionRate:   "1000000000000000000000000000",
			malleate:         func() {},
			expPass:          true,
			expCapacityIncr:  true,
			expMintTDAI:      true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest(sdktypes.DefaultBondDenom, big.NewInt(1000000))

			path = ibctesting.NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			sender := suite.chainA.SenderAccount.GetAddress().String()
			receiver := suite.chainB.SenderAccount.GetAddress().String()

			packetData = transfertypes.FungibleTokenPacketData{
				Denom:    tc.denom,
				Amount:   tc.denomAmountSent,
				Sender:   sender,
				Receiver: receiver,
			}

			app := suite.chainA.App.(*app.App)
			ctx := suite.chainA.GetContext()

			app.RatelimitKeeper.SetSDAIPrice(ctx, keeper.ConvertStringToBigIntWithPanicOnErr(tc.conversionRate))

			module, _, err := app.GetIBCKeeper().PortKeeper.LookupModuleByPort(ctx, ibctesting.TransferPort)
			suite.Require().NoError(err)

			cbs, ok := app.GetIBCKeeper().Router.GetRoute(module)
			suite.Require().True(ok)

			tc.malleate()

			// Record initial capacities and balances
			var initialCapacityList []dtypes.SerializableInt
			var initialSenderBalance sdktypes.Coin
			if tc.expPass {
				initialCapacityList = app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom).CapacityList
				initialSenderBalance = app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), assettypes.TDaiDenom)
			}

			packetForTest := channeltypes.Packet{
				DestinationChannel: path.EndpointB.ChannelID,
				DestinationPort:    path.EndpointB.ChannelConfig.PortID,
				SourceChannel:      path.EndpointA.ChannelID,
				SourcePort:         path.EndpointA.ChannelConfig.PortID,
				Data:               packetData.GetBytes(),
				Sequence:           sequenceNumberOne,
			}

			// If capacity increases, this means the ack was an error. This means we will unescrow tokens during
			// processing. Simulate previous escrowing here.
			denomToMint := packetData.Denom
			if packetData.Denom == ratelimittypes.SDaiBaseDenomFullPath {
				denomToMint = ratelimittypes.SDaiDenom
			}

			escrowedCoin := sdktypes.NewCoins(sdktypes.NewCoin(denomToMint, sdkmath.NewIntFromBigInt(keeper.ConvertStringToBigIntWithPanicOnErr(tc.denomAmountSent))))
			err = app.BankKeeper.MintCoins(ctx, module, escrowedCoin)
			suite.Require().NoError(err)
			escrowAddress := transfertypes.GetEscrowAddress(packetForTest.SourcePort, packetForTest.SourceChannel)
			err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, module, escrowAddress, escrowedCoin)
			suite.Require().NoError(err)

			app.TransferKeeper.SetTotalEscrowForDenom(ctx, escrowedCoin[0])

			app.RatelimitKeeper.SetPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence)

			// Setup Complete. Run test.
			err = cbs.OnTimeoutPacket(ctx, packetForTest, suite.chainB.SenderAccount.GetAddress())

			if !tc.expPass {
				suite.Require().Error(err)
				suite.Require().True(app.RatelimitKeeper.HasPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence))
			} else {
				suite.Require().NoError(err)
				suite.Require().False(app.RatelimitKeeper.HasPendingSendPacket(ctx, packetForTest.SourceChannel, packetForTest.Sequence))

				// TODO: Can also include checks for other capacities not changing.
				newCapacity := app.RatelimitKeeper.GetDenomCapacity(ctx, packetData.Denom)
				if tc.expCapacityIncr {
					for i, capacity := range newCapacity.CapacityList {
						expectedCapacity := new(big.Int).Add(initialCapacityList[i].BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(packetData.Amount))
						suite.Require().Equal(expectedCapacity, capacity)
					}
				} else {
					suite.Require().Equal(initialCapacityList, newCapacity.CapacityList)
				}

				newSenderBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), assettypes.TDaiDenom)
				if tc.expMintTDAI {
					expectedBalance := new(big.Int).Add(initialSenderBalance.Amount.BigInt(), keeper.ConvertStringToBigIntWithPanicOnErr(tc.tDaiAmountMinted))
					suite.Require().Equal(expectedBalance, newSenderBalance.Amount.BigInt())
				} else {
					suite.Require().Equal(newSenderBalance.Amount, newSenderBalance.Amount)
				}

				// user should not have sDAI at the end
				sDaiDenomBalance := app.BankKeeper.GetBalance(ctx, sdktypes.MustAccAddressFromBech32(packetData.Sender), ratelimittypes.SDaiDenom)
				suite.Require().Equal(0, sDaiDenomBalance.Amount.BigInt().Cmp(big.NewInt(0)))
			}
		})
	}
}
