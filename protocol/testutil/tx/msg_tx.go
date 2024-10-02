package tx

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/accountplus/types"
)

func MakeTxBuilder(
	ctx sdk.Context,
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums, accSeqs []uint64,
	signers, signatures []cryptotypes.PrivKey,
	selectedAuthenticators []uint64,
) (client.TxBuilder, error) {
	sigs := make([]signing.SignatureV2, len(signers))

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

	baseTxBuilder := gen.NewTxBuilder()

	txBuilder, ok := baseTxBuilder.(authtx.ExtensionOptionsTxBuilder)
	if !ok {
		return nil, fmt.Errorf("expected authtx.ExtensionOptionsTxBuilder, got %T", baseTxBuilder)
	}
	if len(selectedAuthenticators) > 0 {
		value, err := codectypes.NewAnyWithValue(&types.TxExtension{
			SelectedAuthenticators: selectedAuthenticators,
		})
		if err != nil {
			return nil, err
		}
		txBuilder.SetNonCriticalExtensionOptions(value)
	}

	err = txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	err = txBuilder.SetSignatures(sigs...)
	if err != nil {
		return nil, err
	}
	txBuilder.SetMemo(memo)
	txBuilder.SetFeeAmount(feeAmt)
	txBuilder.SetGasLimit(gas)

	// 2nd round: once all signer infos are set, every signer can sign.
	for i, p := range signatures {
		signerData := authsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		signBytes, err := authsigning.GetSignBytesAdapter(
			ctx, gen.SignModeHandler(), signMode, signerData, txBuilder.GetTx())
		if err != nil {
			panic(err)
		}

		sig, err := p.Sign(signBytes)
		if err != nil {
			panic(err)
		}
		sigs[i].Data.(*signing.SingleSignatureData).Signature = sig
		err = txBuilder.SetSignatures(sigs...)
		if err != nil {
			panic(err)
		}
	}

	return txBuilder, nil
}

// GenTx generates a signed mock transaction.
func GenTx(
	ctx sdk.Context,
	gen client.TxConfig,
	msgs []sdk.Msg,
	feeAmt sdk.Coins,
	gas uint64,
	chainID string,
	accNums, accSeqs []uint64,
	signers, signatures []cryptotypes.PrivKey,
	selectedAuthenticators []uint64,
) (sdk.Tx, error) {
	builder, err := MakeTxBuilder(
		ctx,
		gen,
		msgs,
		feeAmt,
		gas,
		chainID,
		accNums,
		accSeqs,
		signers,
		signatures,
		selectedAuthenticators,
	)
	if err != nil {
		return nil, err
	}
	return builder.GetTx(), nil
}

// Returns the encoded msg as transaction. Will panic if encoding fails.
func MustGetTxBytes(msgs ...sdk.Msg) []byte {
	tx := constants.TestEncodingCfg.TxConfig.NewTxBuilder()
	err := tx.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}
	bz, err := constants.TestEncodingCfg.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		panic(err)
	}
	return bz
}

// Returns the account address that should sign the msg. Will panic if it is an unsupported message type.
func MustGetOnlySignerAddress(cdc codec.Codec, msg sdk.Msg) string {
	signers, _, err := cdc.GetMsgV1Signers(msg)
	if err != nil {
		panic(err)
	}
	if len(signers) == 0 {
		panic(fmt.Errorf("msg does not have designated signer: %T", msg))
	}
	if len(signers) > 1 {
		panic(fmt.Errorf("not supported - msg has multiple signers: %T", msg))
	}
	signer, err := module.InterfaceRegistry.SigningContext().AddressCodec().BytesToString(signers[0])
	if err != nil {
		panic(err)
	}
	return signer
}
