package tx

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/dydxprotocol/v4-chain/protocol/app"
)

// CreateTestTx is a helper function to create a tx given multiple inputs.
func CreateTestTx(
	privs []cryptotypes.PrivKey,
	accNums []uint64,
	accSeqs []uint64,
	chainID string,
	msgs []sdk.Msg,
	timeoutHeight uint64,
) (xauthsigning.Tx, error) {
	encodingConfig := app.GetEncodingConfig()
	clientCtx := client.Context{}.WithTxConfig(encodingConfig.TxConfig)
	txBuilder := clientCtx.TxConfig.NewTxBuilder()
	txBuilder.SetTimeoutHeight(timeoutHeight)
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}

	// First round: we gather all the signer infos. We use the "set empty
	// signature" hack to do that.
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  signing.SignMode_SIGN_MODE_DIRECT,
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
			context.Background(),
			signing.SignMode_SIGN_MODE_DIRECT,
			signerData,
			txBuilder,
			priv,
			clientCtx.TxConfig,
			accSeqs[i],
		)
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
