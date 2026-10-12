package txdecoder

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"google.golang.org/protobuf/encoding/protowire"
)

// NewTxConfig returns a TxConfig whose proto tx decoder rejects a TxRaw that repeats its
// body_bytes or auth_info_bytes field, in addition to the standard decoding. Every consumer of the
// returned config's TxDecoder() (baseapp, the clob/price decoders, the proposal handlers) picks up
// the check.
func NewTxConfig(cdc codec.Codec) client.TxConfig {
	cfg, err := authtx.NewTxConfigWithOptions(cdc, authtx.ConfigOptions{
		EnabledSignModes: authtx.DefaultSignModes,
		ProtoDecoder:     WrapTxDecoder(authtx.DefaultTxDecoder(cdc)),
	})
	if err != nil {
		panic(err)
	}
	return cfg
}

// WrapTxDecoder returns a decoder that rejects duplicate TxRaw fields before delegating to next.
func WrapTxDecoder(next sdk.TxDecoder) sdk.TxDecoder {
	return func(b []byte) (sdk.Tx, error) {
		if err := RejectDuplicateTxRawFields(b); err != nil {
			return nil, errorsmod.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}
		return next(b)
	}
}

// RejectDuplicateTxRawFields returns an error if the encoded TxRaw repeats field 1 (body_bytes) or
// field 2 (auth_info_bytes). The SDK's ADR-027 check only rejects fields encoded in decreasing tag
// order, so bytes like 1,1,2,3 pass it; gogoproto then decodes such a duplicate last-wins, but the
// cometbft mempool classifier merges the two bodies, so the same bytes mean two different txs on
// different code paths. Field 3 (signatures) is legitimately repeated (multisig) and is left alone.
func RejectDuplicateTxRawFields(b []byte) error {
	var seen [3]bool
	for len(b) > 0 {
		num, wt, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]

		if num == 1 || num == 2 {
			if seen[num] {
				return fmt.Errorf("TxRaw field %d must not be repeated", num)
			}
			seen[num] = true
		}

		if n = protowire.ConsumeFieldValue(num, wt, b); n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
	}
	return nil
}
