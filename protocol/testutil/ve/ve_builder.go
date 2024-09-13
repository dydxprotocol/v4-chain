package ve_testutils

import (
	"bytes"
	"fmt"
	"math/big"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"
	ccvkeeper "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/keeper"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SignedVEInfo struct {
	Val     sdk.ConsAddress
	Power   int64
	Prices  []vetypes.PricePair
	Height  int64
	Round   int64
	ChainId string
}

func NewDefaultSignedVeInfo(
	val sdk.ConsAddress,
	prices []vetypes.PricePair,
) SignedVEInfo {
	return SignedVEInfo{
		Val:     val,
		Power:   500,
		Prices:  prices,
		Height:  3,
		Round:   0,
		ChainId: "localdydxprotocol",
	}
}

var (
	voteCodec = vecodec.NewDefaultVoteExtensionCodec()
	extCodec  = vecodec.NewDefaultExtendedCommitCodec()
)

func CreateSignedExtendedCommitInfo(
	veInfo []SignedVEInfo,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	var votes []cometabci.ExtendedVoteInfo
	for _, info := range veInfo {
		extVoteInfo, err := CreateSignedExtendedVoteInfo(info)
		if err != nil {
			continue
		}
		votes = append(votes, extVoteInfo)
	}

	return CreateExtendedCommitInfo(votes)
}

func GetEmptyLocalLastCommit(
	validators []ccvtypes.CrossChainValidator,
	height int64,
	round int64,
	chainId string,
) cometabci.ExtendedCommitInfo {
	var votes []cometabci.ExtendedVoteInfo
	for _, validator := range validators {
		ve, err := CreateSignedExtendedVoteInfo(
			SignedVEInfo{
				Val:     sdk.ConsAddress(validator.Address),
				Power:   validator.GetPower(),
				Prices:  []vetypes.PricePair{},
				Height:  height,
				Round:   round,
				ChainId: chainId,
			},
		)

		if err != nil {
			panic(err)
		}
		votes = append(votes, ve)
	}
	extCommitInfo, _, _ := CreateExtendedCommitInfo(votes)
	return extCommitInfo
}

// CreateExtendedCommitInfo creates an extended commit info with the given commit info.
func CreateExtendedCommitInfo(
	commitInfo []cometabci.ExtendedVoteInfo,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	extendedCommitInfo := cometabci.ExtendedCommitInfo{
		Votes: commitInfo,
	}

	bz, err := extCodec.Encode(extendedCommitInfo)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, err
	}

	return extendedCommitInfo, bz, nil
}

// CreateExtendedVoteInfoWithPower CreateExtendedVoteInfo creates an extended vote info
// with the given power, prices, timestamp and height.
func CreateSignedExtendedVoteInfo(veInfo SignedVEInfo) (cometabci.ExtendedVoteInfo, error) {
	ve, err := CreateVoteExtensionBytes(veInfo.Prices)
	if err != nil {
		return cometabci.ExtendedVoteInfo{}, err
	}

	sig, err := signVoteExtension(veInfo, ve)
	if err != nil {
		return cometabci.ExtendedVoteInfo{}, err
	}

	voteInfo := cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: veInfo.Val,
			Power:   veInfo.Power,
		},
		VoteExtension:      ve,
		BlockIdFlag:        cometproto.BlockIDFlagCommit,
		ExtensionSignature: sig,
	}

	return voteInfo, nil
}

// CreateVoteExtensionBytes creates a vote extension bytes with the given prices, timestamp and height.
func CreateVoteExtensionBytes(
	prices []vetypes.PricePair,
) ([]byte, error) {
	voteExtension := CreateVoteExtension(prices)
	voteExtensionBz, err := voteCodec.Encode(voteExtension)
	if err != nil {
		return nil, err
	}

	return voteExtensionBz, nil
}

// CreateVoteExtension creates a vote extension with the given prices, timestamp and height.
func CreateVoteExtension(
	prices []vetypes.PricePair,
) vetypes.DaemonVoteExtension {
	return vetypes.DaemonVoteExtension{
		Prices: prices,
	}
}

func GetVeEnabledCtx(ctx sdk.Context, blockHeight int64) sdk.Context {
	ctx = ctx.WithConsensusParams(
		cometproto.ConsensusParams{
			Abci: &cometproto.ABCIParams{
				VoteExtensionsEnableHeight: 2,
			},
		},
	).WithBlockHeight(blockHeight)
	return ctx
}

func GetEmptyProposedLastCommit() cometabci.CommitInfo {
	return cometabci.CommitInfo{
		Round: 0,
		Votes: []cometabci.VoteInfo{
			{
				Validator: cometabci.Validator{
					Address: constants.CarlEthosConsAddress,
					Power:   500,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
			{
				Validator: cometabci.Validator{
					Address: constants.AliceEthosConsAddress,
					Power:   500,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
			{
				Validator: cometabci.Validator{
					Address: constants.BobEthosConsAddress,
					Power:   500,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
		},
	}
}

func GetDaemonPriceCacheEncodedPrice(price *big.Int) ([]byte, error) {
	if price.Sign() < 0 {
		return nil, fmt.Errorf("price must be non-negative %v", price.String())
	}

	return price.GobEncode()
}

func GetDaemonPriceCacheDecodedPrice(priceBz []byte) (*big.Int, error) {
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

func CreateSingleValidatorExtendedCommitInfo(
	consAddr sdk.ConsAddress,
	prices []vetypes.PricePair,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	voteInfo, err := CreateSignedExtendedVoteInfo(
		NewDefaultSignedVeInfo(
			consAddr,
			prices,
		),
	)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, err
	}

	return CreateExtendedCommitInfo([]cometabci.ExtendedVoteInfo{voteInfo})
}

func signVoteExtension(
	veInfo SignedVEInfo,
	voteExtension []byte,
) ([]byte, error) {
	privKey := constants.GetPrivKeyFromConsAddress(veInfo.Val)

	cve := cometproto.CanonicalVoteExtension{
		Height:    veInfo.Height,
		Round:     veInfo.Round,
		ChainId:   veInfo.ChainId,
		Extension: voteExtension,
	}
	extSignBytes, err := marshalDelimited(&cve)

	if err != nil {
		return nil, err
	}

	return privKey.Sign(extSignBytes)
}

func marshalDelimited(msg proto.Message) ([]byte, error) {
	var buf bytes.Buffer
	if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GetInjectedExtendedCommitInfoForTestApp(
	consumerKeeper *ccvkeeper.Keeper,
	ctx sdk.Context,
	prices map[uint32]ve.VEPricePair,
	height int64,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	var pricesBz = make([]vetypes.PricePair, len(prices))
	for marketId, price := range prices {
		encodedSpotPrice, err := GetDaemonPriceCacheEncodedPrice(new(big.Int).SetUint64(price.SpotPrice))
		if err != nil {
			return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to encode price: %w", err)
		}

		encodedPnlPrice, err := GetDaemonPriceCacheEncodedPrice(new(big.Int).SetUint64(price.PnlPrice))
		if err != nil {
			return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to encode price: %w", err)
		}

		pricesBz[marketId] = vetypes.PricePair{
			MarketId:  marketId,
			SpotPrice: encodedSpotPrice,
			PnlPrice:  encodedPnlPrice,
		}
	}

	validators := consumerKeeper.GetAllCCValidator(ctx)

	var veSignedInfos []SignedVEInfo
	for _, v := range validators {
		veSignedInfos = append(veSignedInfos, SignedVEInfo{
			Val:     sdk.ConsAddress(v.Address),
			Power:   v.GetPower(),
			Prices:  pricesBz,
			Height:  height,
			Round:   0,
			ChainId: "localdydxprotocol",
		})
	}

	extCommitInfo, extCommitBz, err := CreateSignedExtendedCommitInfo(veSignedInfos)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to create signed extended commit info: %w", err)
	}
	return extCommitInfo, extCommitBz, nil
}
