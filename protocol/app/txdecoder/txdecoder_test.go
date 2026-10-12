package txdecoder_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/dydxprotocol/v4-chain/protocol/app/txdecoder"
	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

// dupField prepends an extra field/value pair with the given tag to valid, producing a TxRaw
// that repeats that field.
func dupField(tag uint32, valid []byte) []byte {
	dup := protowire.AppendTag(nil, protowire.Number(tag), protowire.BytesType)
	dup = protowire.AppendBytes(dup, []byte("x"))
	return append(dup, valid...)
}

func TestRejectDuplicateTxRawFields(t *testing.T) {
	valid := constants.Msg_Send_TxBytes

	truncated := protowire.AppendTag(nil, 1, protowire.BytesType)

	tests := map[string]struct {
		txBytes []byte
		wantErr bool
	}{
		"valid tx: no error": {
			txBytes: valid,
			wantErr: false,
		},
		"duplicate body_bytes (field 1): error": {
			txBytes: dupField(1, valid),
			wantErr: true,
		},
		"duplicate auth_info_bytes (field 2): error": {
			txBytes: dupField(2, valid),
			wantErr: true,
		},
		"extra signatures (field 3) is allowed: no error": {
			txBytes: dupField(3, valid),
			wantErr: false,
		},
		"truncated bytes: error": {
			txBytes: truncated,
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := txdecoder.RejectDuplicateTxRawFields(tc.txBytes)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestWrappedDecoder_ReturnsErrTxDecode(t *testing.T) {
	decoder := txdecoder.WrapTxDecoder(tx.DefaultTxDecoder(constants.TestEncodingCfg.Codec))

	valid := constants.Msg_Send_TxBytes
	_, err := decoder(valid)
	require.NoError(t, err)

	dup := dupField(1, valid)
	_, err = decoder(dup)
	require.Error(t, err)
	require.ErrorIs(t, err, sdkerrors.ErrTxDecode)
}
