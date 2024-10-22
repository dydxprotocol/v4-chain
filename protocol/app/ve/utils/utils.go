package ve_utils

import (
	"bytes"
	"context"
	"fmt"
	"math/big"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto"
	cryptoenc "github.com/cometbft/cometbft/crypto/encoding"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"
)

type ValidatorNotFoundError struct {
	Address []byte
}

func (e *ValidatorNotFoundError) Error() string {
	return fmt.Sprintf("validator %X not found", e.Address)
}

type ValidatorPubKeyStore interface {
	GetPubKeyByConsAddr(context.Context, sdk.ConsAddress) (cmtprotocrypto.PublicKey, error)
}

func AreVEEnabled(ctx sdk.Context) bool {
	cp := ctx.ConsensusParams()
	if cp.Abci == nil || cp.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	if ctx.BlockHeight() <= 1 {
		return false
	}

	return cp.Abci.VoteExtensionsEnableHeight < ctx.BlockHeight()
}

func GetPriceFromBytes(
	id uint32,
	bz []byte,
) (*big.Int, error) {
	price, err := GetVEDecodedPrice(bz)

	if err != nil {
		return nil, err
	}

	return price, nil
}

func GetVEDecodedPrice(
	priceBz []byte,
) (*big.Int, error) {
	var price big.Int
	err := price.GobDecode(priceBz)
	if err != nil {
		return nil, err
	}

	if price.Sign() < 0 {
		return nil, fmt.Errorf("price must be non-negative %v", price.String())
	}

	return &price, nil
}

func GetVEEncodedPrice(
	price *big.Int,
) ([]byte, error) {
	if price.Sign() < 0 {
		return nil, fmt.Errorf("price must be non-negative %v", price.String())
	}

	return price.GobEncode()
}

// marshalDelimited serializes a proto.Message into a delimited byte slice.
func MarshalDelimited(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetValPubKeyFromVote(
	ctx sdk.Context,
	vote cometabci.ExtendedVoteInfo,
	validatorStore ValidatorPubKeyStore,
) (crypto.PubKey, error) {
	valConsAddr := sdk.ConsAddress(vote.Validator.Address)

	pubKeyProto, err := validatorStore.GetPubKeyByConsAddr(ctx, valConsAddr)
	if err != nil {
		return nil, &ValidatorNotFoundError{Address: valConsAddr}
	}

	cmtPubKey, err := cryptoenc.PubKeyFromProto(pubKeyProto)
	if err != nil {
		return nil, fmt.Errorf("failed to convert validator %X public key: %w", valConsAddr, err)
	}

	return cmtPubKey, nil
}
