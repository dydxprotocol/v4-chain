package tx

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/dydxprotocol/v4-chain/protocol/app/module"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
)

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

// CreateTestTx is a helper function to create a tx given multiple inputs.
func CreateTestTx(
	ctx sdk.Context,
	msgs []sdk.Msg,
	privs []cryptotypes.PrivKey,
	accNums, accSeqs []uint64,
	chainID string, signMode signing.SignMode, txConfig client.TxConfig,
	timeoutHeight uint64,
) (xauthsigning.Tx, error) {
	txBuilder := txConfig.NewTxBuilder()
	txBuilder.SetTimeoutHeight(timeoutHeight)
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		panic(err)
	}

	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signMode,
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	// Second round: all signer infos are set, so each signer can sign.
	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			Address:       sdk.AccAddress(priv.PubKey().Address()).String(),
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
			PubKey:        priv.PubKey(),
		}
		sigV2, err := tx.SignWithPrivKey(
			ctx, signMode, signerData,
			txBuilder, priv, txConfig, accSeqs[i])
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err = txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return nil, err
	}

	return txBuilder.GetTx(), nil
}
