package authenticator_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/stretchr/testify/suite"

	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/app/config"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/authenticator"
)

type SigVerifyAuthenticationSuite struct {
	BaseAuthenticatorSuite

	SigVerificationAuthenticator authenticator.SignatureVerification
}

func TestSigVerifyAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(SigVerifyAuthenticationSuite))
}

func (s *SigVerifyAuthenticationSuite) SetupTest() {
	s.SetupKeys()

	s.EncodingConfig = app.GetEncodingConfig()
	ak := s.tApp.App.AccountKeeper

	// Create a new Secp256k1SignatureAuthenticator for testing
	s.SigVerificationAuthenticator = authenticator.NewSignatureVerification(
		ak,
	)
}

func (s *SigVerifyAuthenticationSuite) TearDownTest() {
	os.RemoveAll(s.HomeDir)
}

type SignatureVerificationTestData struct {
	Msgs                               []sdk.Msg
	AccNums                            []uint64
	AccSeqs                            []uint64
	Signers                            []cryptotypes.PrivKey
	Signatures                         []cryptotypes.PrivKey
	NumberOfExpectedSigners            int
	NumberOfExpectedSignatures         int
	ShouldSucceedGettingData           bool
	ShouldSucceedSignatureVerification bool
}

type SignatureVerificationTest struct {
	Description string
	TestData    SignatureVerificationTestData
}

// TestSignatureAuthenticator test a non-smart account signature verification
func (s *SigVerifyAuthenticationSuite) TestSignatureAuthenticator() {
	bech32Prefix := config.Bech32PrefixAccAddr
	coins := sdk.Coins{sdk.NewInt64Coin(constants.TestNativeTokenDenom, 2500)}

	// Create a test messages for signing
	testMsg1 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	testMsg2 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	testMsg3 := &banktypes.MsgSend{
		FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[2]),
		ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
		Amount:      coins,
	}
	//testMsg4 := &banktypes.MsgSend{
	//	FromAddress: sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[0]),
	//	ToAddress:   sdk.MustBech32ifyAddressBytes(bech32Prefix, s.TestAccAddress[1]),
	//	Amount:      coins,
	//}
	feeCoins := constants.TestFeeCoins_5Cents

	tests := []SignatureVerificationTest{
		{
			Description: "Test: successfully verified authenticator with one signer: base case: PASS",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
				},
				[]uint64{5},
				[]uint64{0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
				},
				1,
				1,
				true,
				true,
			},
		},

		{
			Description: "Test: successfully verified authenticator: multiple signers: PASS",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
					testMsg2,
					testMsg3,
				},
				[]uint64{5, 5, 5, 5},
				[]uint64{0, 0, 0, 0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
					s.TestPrivKeys[2],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
					s.TestPrivKeys[2],
				},
				3,
				3,
				true,
				true,
			},
		},

		{
			// This test case tests if there are two messages with the same signer
			// with two successful signatures.
			Description: "Test: verified authenticator with 2 messages signed correctly with the same address: PASS",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
					testMsg2,
					testMsg2,
				},
				[]uint64{5, 5},
				[]uint64{0, 0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
				},
				2,
				2,
				true,
				true,
			},
		},

		{
			// This test case tests if there are two messages with the same signer
			// with two successful signatures.
			Description: "Test: verified authenticator with 2 messages but only first signed signed correctly: Fail",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
					testMsg2,
					testMsg2,
				},
				[]uint64{5, 5},
				[]uint64{0, 0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[0],
				},
				2,
				2,
				true,
				false,
			},
		},

		{
			// This test case tests if there are two messages with the same signer
			// with two successful signatures.
			Description: "Test: verified authenticator with 2 messages but only second signed signed correctly: Fail",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
					testMsg2,
					testMsg2,
				},
				[]uint64{5, 5},
				[]uint64{0, 0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[0],
					s.TestPrivKeys[1],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[1],
					s.TestPrivKeys[1],
				},
				2,
				2,
				true,
				false,
			},
		},

		{
			Description: "Test: unsuccessful signature authentication invalid signatures: FAIL",
			TestData: SignatureVerificationTestData{
				[]sdk.Msg{
					testMsg1,
					testMsg2,
				},
				[]uint64{5, 5},
				[]uint64{0, 0},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[1],
					s.TestPrivKeys[0],
				},
				[]cryptotypes.PrivKey{
					s.TestPrivKeys[2],
					s.TestPrivKeys[0],
				},
				2,
				2,
				false,
				false,
			},
		},
	}

	for _, tc := range tests {
		s.Run(tc.Description, func() {
			// Generate a transaction based on the test cases
			tx, _ := GenTx(
				s.Ctx,
				s.EncodingConfig.TxConfig,
				tc.TestData.Msgs,
				feeCoins,
				300000,
				s.Ctx.ChainID(),
				tc.TestData.AccNums,
				tc.TestData.AccSeqs,
				tc.TestData.Signers,
				tc.TestData.Signatures,
			)
			ak := s.tApp.App.AccountKeeper
			sigModeHandler := s.EncodingConfig.TxConfig.SignModeHandler()

			// Only the first message is tested for authenticate
			addr := sdk.AccAddress(tc.TestData.Signers[0].PubKey().Address())

			if tc.TestData.ShouldSucceedGettingData {
				// request for the first message
				request, err := authenticator.GenerateAuthenticationRequest(
					s.Ctx,
					s.tApp.App.AppCodec(),
					ak,
					sigModeHandler,
					addr,
					addr,
					nil,
					sdk.NewCoins(),
					tc.TestData.Msgs[0],
					tx,
					0,
					false,
				)
				s.Require().NoError(err)

				// Test Authenticate method
				if tc.TestData.ShouldSucceedSignatureVerification {
					initialized, err := s.SigVerificationAuthenticator.Initialize(tc.TestData.Signers[0].PubKey().Bytes())
					s.Require().NoError(err)
					err = initialized.Authenticate(s.Ctx, request)
					s.Require().NoError(err)
				} else {
					err = s.SigVerificationAuthenticator.Authenticate(s.Ctx, request)
					s.Require().Error(err)
				}
			} else {
				_, err := authenticator.GenerateAuthenticationRequest(
					s.Ctx,
					s.tApp.App.AppCodec(),
					ak,
					sigModeHandler,
					addr,
					addr,
					nil,
					sdk.NewCoins(),
					tc.TestData.Msgs[0],
					tx,
					0,
					false,
				)
				s.Require().Error(err)
			}
		})
	}
}

func MakeTxBuilder(ctx sdk.Context,
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums,
	accSeqs []uint64,
	signers []cryptotypes.PrivKey,
	signatures []cryptotypes.PrivKey,
) (client.TxBuilder, error) {
	sigs := make([]signing.SignatureV2, len(signatures))

	// create a random length memo
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	memo := simulation.RandStringOfLength(r, simulation.RandIntBetween(r, 0, 100))
	signMode, err := authsigning.APISignModeToInternal(gen.SignModeHandler().DefaultMode())
	if err != nil {
		return nil, err
	}

	// 1st round: set SignatureV2 with empty signatures, to set correct
	// signer infos.
	for i, p := range signers {
		sigs[i] = signing.SignatureV2{
			PubKey: p.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode: signMode,
			},
			Sequence: accSeqs[i],
		}
	}

	tx := gen.NewTxBuilder()
	err = tx.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = tx.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	tx.SetMemo(memo)
	tx.SetFeeAmount(feeAmt)
	tx.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range signatures {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := authsigning.GetSignBytesAdapter(
			ctx, gen.SignModeHandler(), signMode, signerData, tx.GetTx())
		if err != nil {
			panic(err)
		}
		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
	}

	err = tx.SetSignatures(sigs...)
	if err != nil {
		panic(err)
	}
	return tx, nil
}

// GenTx generates a signed mock transaction.
func GenTx(
	ctx sdk.Context,
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums,
	accSeqs []uint64,
	signers []cryptotypes.PrivKey,
	signatures []cryptotypes.PrivKey,
) (sdk.Tx, error) {
	tx, err := MakeTxBuilder(ctx, gen, msgs, feeAmt, gas, chainID, accNums, accSeqs, signers, signatures)
	if err != nil {
		return nil, err
	}
	return tx.GetTx(), nil
}

// func generatePubKeysForMultiSig(
// 	priv ...cryptotypes.PrivKey,
// ) (pubkeys []cryptotypes.PubKey) {
// 	pubkeys = make([]cryptotypes.PubKey, len(priv))
// 	for i, p := range priv {
// 		pubkeys[i] = p.PubKey()
// 	}
// 	return
// }

// func generateSignaturesForMultiSig(
// 	msg []byte,
// 	priv ...cryptotypes.PrivKey,
// ) (signatures [][]byte) {
// 	signatures = make([][]byte, len(priv))
// 	for i, p := range priv {
// 		var err error
// 		signatures[i], err = p.Sign(msg)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// 	return
// }
