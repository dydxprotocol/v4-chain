package ve_testutils

import (
	"bytes"
	"fmt"
	"math/big"

	ve "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve"
	vecodec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	voteweighted "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/math"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	"github.com/StreamFinance-Protocol/stream-chain/protocol/testutil/constants"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type SignedVEInfo struct {
	Val                sdk.ConsAddress
	Power              int64
	Prices             []vetypes.PricePair
	SDaiConversionRate string
	Height             int64
	Round              int64
	ChainId            string
}

func NewDefaultSignedVeInfo(
	val sdk.ConsAddress,
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) SignedVEInfo {
	return SignedVEInfo{
		Val:                val,
		Power:              500,
		Prices:             prices,
		SDaiConversionRate: sdaiConversionRate,
		Height:             3,
		Round:              0,
		ChainId:            "localdydxprotocol",
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
	validators []stakingtypes.Validator,
	height int64,
	round int64,
	chainId string,
) cometabci.ExtendedCommitInfo {
	var votes []cometabci.ExtendedVoteInfo
	fmt.Println("GETEMPTYLOCALLASTCOMMIT STARTED")
	for _, validator := range validators {
		valConsAddr := constants.GetConsAddressFromStringValidatorAddress(validator.OperatorAddress)
		fmt.Println("SUCCESSFULLY GOT VALIDATOR ADDRESS IN GET EMPTY LOCAL LAST COMMIT", valConsAddr)
		votingPower := voteweighted.GetPowerFromBondedTokens(validator.Tokens)
		fmt.Println("SUCCESSFULLY GOT VOTING POWER IN GET EMPTY LOCAL LAST COMMIT", votingPower)

		ve, err := CreateSignedExtendedVoteInfo(
			SignedVEInfo{
				Val:                valConsAddr,
				Power:              votingPower,
				Prices:             []vetypes.PricePair{},
				SDaiConversionRate: "",
				Height:             height,
				Round:              round,
				ChainId:            chainId,
			},
		)

		if err != nil {
			panic(err)
		}
		votes = append(votes, ve)
	}
	fmt.Println("VOTES", votes)
	extCommitInfo, _, _ := CreateExtendedCommitInfo(votes)
	fmt.Println("EXTENDED COMMIT INFO ", extCommitInfo)
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
	ve, err := CreateVoteExtensionBytes(veInfo.Prices, veInfo.SDaiConversionRate)
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

func CreateNilVoteExtensionInfo(consAddress sdk.ConsAddress, power int64) (cometabci.ExtendedVoteInfo, error) {
	voteInfo := cometabci.ExtendedVoteInfo{
		Validator: cometabci.Validator{
			Address: consAddress,
			Power:   power,
		},
		VoteExtension:      nil,
		BlockIdFlag:        cometproto.BlockIDFlagAbsent,
		ExtensionSignature: nil,
	}

	return voteInfo, nil
}

// CreateVoteExtensionBytes creates a vote extension bytes with the given prices, timestamp and height.
func CreateVoteExtensionBytes(
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) ([]byte, error) {
	voteExtension := CreateVoteExtension(prices, sdaiConversionRate)
	voteExtensionBz, err := voteCodec.Encode(voteExtension)
	if err != nil {
		return nil, err
	}

	return voteExtensionBz, nil
}

// CreateVoteExtension creates a vote extension with the given prices, timestamp and height.
func CreateVoteExtension(
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) vetypes.DaemonVoteExtension {
	return vetypes.DaemonVoteExtension{
		Prices:             prices,
		SDaiConversionRate: sdaiConversionRate,
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
					Address: constants.AliceConsAddress,
					Power:   500000,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
			{
				Validator: cometabci.Validator{
					Address: constants.CarlConsAddress,
					Power:   500000,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
			{
				Validator: cometabci.Validator{
					Address: constants.DaveConsAddress,
					Power:   500000,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
			{
				Validator: cometabci.Validator{
					Address: constants.BobConsAddress,
					Power:   500000,
				},
				BlockIdFlag: cometproto.BlockIDFlagCommit,
			},
		},
	}
}

func GetVECacheEncodedPrice(price *big.Int) ([]byte, error) {
	if price.Sign() < 0 {
		return nil, fmt.Errorf("price must be non-negative %v", price.String())
	}

	return price.GobEncode()
}

func CreateSingleValidatorExtendedCommitInfo(
	consAddr sdk.ConsAddress,
	prices []vetypes.PricePair,
	sdaiConversionRate string,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	voteInfo, err := CreateSignedExtendedVoteInfo(
		NewDefaultSignedVeInfo(
			consAddr,
			prices,
			sdaiConversionRate,
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
	stakingKeeper *stakingkeeper.Keeper,
	ctx sdk.Context,
	prices map[uint32]ve.VEPricePair,
	sdaiConversionRate string,
	height int64,
) (cometabci.ExtendedCommitInfo, []byte, error) {
	var pricesBz = make([]vetypes.PricePair, len(prices))
	for marketId, price := range prices {
		encodedSpotPrice, err := GetVECacheEncodedPrice(new(big.Int).SetUint64(price.SpotPrice))
		if err != nil {
			return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to encode price: %w", err)
		}

		encodedPnlPrice, err := GetVECacheEncodedPrice(new(big.Int).SetUint64(price.PnlPrice))
		if err != nil {
			return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to encode price: %w", err)
		}

		pricesBz[marketId] = vetypes.PricePair{
			MarketId:  marketId,
			SpotPrice: encodedSpotPrice,
			PnlPrice:  encodedPnlPrice,
		}
	}

	validators, err := stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		fmt.Println("MASSIVE FAILURE IN GETTING BONDED VALIDATORS")
		return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to get bonded validators: %w", err)
	}

	fmt.Println("MASSIVEVALIDATORS", validators)

	var veSignedInfos []SignedVEInfo
	for _, validator := range validators {
		valConsAddr := constants.GetConsAddressFromStringValidatorAddress(validator.OperatorAddress)

		fmt.Println("IN LOOP VALIDATOR", validator.OperatorAddress)
		fmt.Println("VALCONSADDR", valConsAddr)
		veSignedInfos = append(veSignedInfos, SignedVEInfo{
			Val:                valConsAddr,
			Power:              voteweighted.GetPowerFromBondedTokens(validator.Tokens),
			Prices:             pricesBz,
			SDaiConversionRate: sdaiConversionRate,
			Height:             height,
			Round:              0,
			ChainId:            "localdydxprotocol",
		})
	}

	extCommitInfo, extCommitBz, err := CreateSignedExtendedCommitInfo(veSignedInfos)
	if err != nil {
		return cometabci.ExtendedCommitInfo{}, nil, fmt.Errorf("failed to create signed extended commit info: %w", err)
	}
	return extCommitInfo, extCommitBz, nil
}
