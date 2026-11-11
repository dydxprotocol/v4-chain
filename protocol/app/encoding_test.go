package app_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/dydxprotocol/v4-chain/protocol/app"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

func TestAppTxDecoder_RejectsDuplicateBody(t *testing.T) {
	txConfig := app.GetEncodingConfig().TxConfig

	txBuilder := txConfig.NewTxBuilder()
	require.NoError(t, txBuilder.SetMsgs(constants.Msg_Send))
	valid, err := txConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err)

	_, err = txConfig.TxDecoder()(valid)
	require.NoError(t, err)

	dup := protowire.AppendTag(nil, 1, protowire.BytesType)
	dup = protowire.AppendBytes(dup, []byte("x"))
	dup = append(dup, valid...)

	_, err = txConfig.TxDecoder()(dup)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrTxDecode)
}
