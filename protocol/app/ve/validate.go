package ve

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"cosmossdk.io/core/comet"
	constants "github.com/StreamFinance-Protocol/stream-chain/protocol/app/constants"
	codec "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	vetypes "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/types"
	veutils "github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/utils"
	ratelimittypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/ratelimit/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtprotocrypto "github.com/cometbft/cometbft/proto/tendermint/crypto"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
// function is not explicitly used to validate the oracle data but rather that
// the signed vote extensions included in the proposal are valid and provide
// a super-majority of vote extensions for the current block. This method is
// expected to be used in PrepareProposal and ProcessProposal.
type ValidateVEConsensusInfoFn func(
	ctx sdk.Context,
	extInfo cometabci.ExtendedCommitInfo,
) error

// ValidatorStore defines the interface contract require for verifying vote
// extension signatures. Typically, this will be implemented by the x/staking
// module, which has knowledge of the CometBFT public key.
type ValidatorStore interface {
	GetCCValidator(ctx sdk.Context, addr []byte) (ccvtypes.CrossChainValidator, bool)
}

// ---------------------------- VE VALIDATION ----------------------------

func CleanAndValidateExtCommitInfo(
	ctx sdk.Context,
	extCommitInfo cometabci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
	pricesKeeper PreBlockExecPricesKeeper,
	ratelimitKeeper VoteExtensionRateLimitKeeper,
) (cometabci.ExtendedCommitInfo, error) {
	for i, vote := range extCommitInfo.Votes {
		if err := validateIndividualVoteExtension(ctx, vote, veCodec, pricesKeeper, ratelimitKeeper); err != nil {
			ctx.Logger().Info(
				"failed to validate vote extension - pruning vote",
				"err", err,
				"validator", vote.Validator.Address,
			)

			// failed to validate this vote-extension, mark it as absent in the original commit
			pruneVoteFromExtCommitInfo(&vote, &extCommitInfo, i)
		}
	}

	return extCommitInfo, nil
}

func ValidateExtendedCommitInfo(
	ctx sdk.Context,
	height int64,
	extCommitInfo cometabci.ExtendedCommitInfo,
	veCodec codec.VoteExtensionCodec,
	pricesKeeper PreBlockExecPricesKeeper,
	ratelimitKeeper VoteExtensionRateLimitKeeper,
	validateVEConsensusInfo ValidateVEConsensusInfoFn,
) error {
	if err := validateVEConsensusInfo(ctx, extCommitInfo); err != nil {
		ctx.Logger().Error(
			"failed to validate vote extension",
			"height", height,
			"err", err,
		)
		return err
	}

	for _, vote := range extCommitInfo.Votes {
		addr := sdk.ConsAddress(vote.Validator.Address)

		if err := validateIndividualVoteExtension(ctx, vote, veCodec, pricesKeeper, ratelimitKeeper); err != nil {
			ctx.Logger().Error(
				"failed to validate vote extension",
				"height", height,
				"validator", addr,
				"err", err,
			)
			return err
		}
	}
	return nil
}

func validateIndividualVoteExtension(
	ctx sdk.Context,
	vote cometabci.ExtendedVoteInfo,
	voteCodec codec.VoteExtensionCodec,
	pricesKeeper PreBlockExecPricesKeeper,
	ratelimitKeeper VoteExtensionRateLimitKeeper,

) error {
	if vote.VoteExtension == nil && vote.ExtensionSignature == nil {
		return nil
	}

	if err := ValidateVEMarketsAndPrices(ctx, pricesKeeper, vote.VoteExtension, voteCodec); err != nil {
		return err
	}

	if err := ValidateVeSDaiConversionRate(ctx, ratelimitKeeper, vote.VoteExtension, voteCodec); err != nil {
		return err
	}

	return nil
}

func ValidateVEMarketsAndPrices(
	ctx sdk.Context,
	pricesKeeper PreBlockExecPricesKeeper,
	veBytes []byte,
	voteCodec codec.VoteExtensionCodec,
) error {
	ve, err := voteCodec.Decode(veBytes)

	if err != nil {
		return err
	}

	if err := ValidateMarketCountInVE(ctx, ve, pricesKeeper); err != nil {
		return err
	}

	if err := ValidatePricesBytesSizeInVE(ctx, ve); err != nil {
		return err
	}

	return nil
}

func ValidateVeSDaiConversionRate(
	ctx sdk.Context,
	ratelimitKeeper VoteExtensionRateLimitKeeper,
	veBytes []byte,
	voteCodec codec.VoteExtensionCodec,
) error {
	ve, err := voteCodec.Decode(veBytes)

	if err != nil {
		return err
	}

	if ve.SDaiConversionRate == "" {
		return nil
	}

	if err := ValidateSDaiConversionRateHeightInVE(ctx, ve, ratelimitKeeper); err != nil {
		return err
	}

	if err := ValidateSDaiConversionRateSizeInVE(ctx, ve); err != nil {
		return err
	}

	if err := ValidateSDaiConversionRateValueInVE(ctx, ve, ratelimitKeeper); err != nil {
		return err
	}

	return nil
}

func ValidateMarketCountInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
	pricesKeeper PreBlockExecPricesKeeper,
) error {
	maxPairs := GetMaxMarketPairs(ctx, pricesKeeper)
	if uint32(len(ve.Prices)) > maxPairs {
		return fmt.Errorf(
			"number of oracle vote extension pairs of %d greater than maximum expected pairs of %d",
			uint64(len(ve.Prices)),
			uint64(maxPairs),
		)
	}
	return nil
}

func ValidatePricesBytesSizeInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
) error {
	for _, pricePair := range ve.Prices {
		if len(pricePair.SpotPrice) > constants.MaximumPriceSizeInBytes {
			return fmt.Errorf("spot price bytes are too long: %d", len(pricePair.SpotPrice))
		}

		if len(pricePair.PnlPrice) > constants.MaximumPriceSizeInBytes {
			return fmt.Errorf("pnl price bytes are too long: %d", len(pricePair.PnlPrice))
		}
	}
	return nil
}

func ValidateSDaiConversionRateSizeInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
) error {
	if len(ve.SDaiConversionRate) > constants.MaxSDaiConversionRateLengthCharacters {
		return fmt.Errorf("sDai conversion rate length (%d) exceeds maximum allowed length (%d)", len(ve.SDaiConversionRate), constants.MaxSDaiConversionRateLengthCharacters)
	}
	return nil
}

func ValidateSDaiConversionRateHeightInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
	ratelimitKeeper VoteExtensionRateLimitKeeper,
) error {

	lastBlockUpdated, found := ratelimitKeeper.GetSDAILastBlockUpdated(ctx)
	if found {
		if ctx.BlockHeight()-lastBlockUpdated.Int64() < ratelimittypes.SDAI_UPDATE_BLOCK_DELAY {
			return fmt.Errorf("sDai conversion rate height is not within the allowed delay of %d blocks", ratelimittypes.SDAI_UPDATE_BLOCK_DELAY)
		}
	}
	return nil
}

func ValidateSDaiConversionRateValueInVE(
	ctx sdk.Context,
	ve vetypes.DaemonVoteExtension,
	ratelimitKeeper VoteExtensionRateLimitKeeper,
) error {
	sDaiConversionRate, ok := new(big.Int).SetString(ve.SDaiConversionRate, 10)
	if !ok {
		return fmt.Errorf("failed to convert sDai conversion rate to big.Int: %s", ve.SDaiConversionRate)
	}

	// TODO: Left in to exit early if the rate is not positive. Could remove this given below check.
	if sDaiConversionRate.Sign() <= 0 {
		return fmt.Errorf("sDai conversion rate must be positive: %s", ve.SDaiConversionRate)
	}

	prevRate, found := ratelimitKeeper.GetSDAIPrice(ctx)
	if found && sDaiConversionRate.Cmp(prevRate) <= 0 {
		return fmt.Errorf("new sDai conversion rate (%s) is not greater than the previous rate (%s)", ve.SDaiConversionRate, prevRate.String())
	}

	return nil
}

func pruneVoteFromExtCommitInfo(
	vote *cometabci.ExtendedVoteInfo,
	extCommitInfo *cometabci.ExtendedCommitInfo,
	index int,
) {
	vote.BlockIdFlag = cmtproto.BlockIDFlagAbsent
	vote.ExtensionSignature = nil
	vote.VoteExtension = nil
	extCommitInfo.Votes[index] = *vote
}

// ---------------------------- CONSENSUS VALIDATION ----------------------------

// NewDefaultValidateVoteExtensionsFn returns a new DefaultValidateVoteExtensionsFn.
func NewValidateVEConsensusInfo(validatorStore ValidatorStore) ValidateVEConsensusInfoFn {
	return func(ctx sdk.Context, info cometabci.ExtendedCommitInfo) error {
		return ValidateVEConsensusInfo(ctx, validatorStore, info)
	}
}

// ValidateVoteExtensions defines a helper function for verifying vote extension
// signatures that may be passed or manually injected into a block proposal from
// a proposer in PrepareProposal. It returns an error if any signature is invalid
// or if unexpected vote extensions and/or signatures are found or less than 2/3
// power is received.
func ValidateVEConsensusInfo(
	ctx sdk.Context,
	valStore ValidatorStore,
	extCommit cometabci.ExtendedCommitInfo,
) error {
	currentHeight := ctx.HeaderInfo().Height
	chainID := ctx.HeaderInfo().ChainID
	commitInfo := ctx.CometInfo().GetLastCommit()

	if err := ValidateExtendedCommitAgainstLastCommit(extCommit, commitInfo); err != nil {
		return err
	}

	// Start checking vote extensions only **after** the vote extensions enable
	// height, because when `currentHeight == VoteExtensionsEnableHeight`
	// PrepareProposal doesn't get any vote extensions in its request.

	var (
		// Total voting power of all vote extensions.
		totalVP int64
		// Total voting power of all validators that submitted valid vote extensions.
		sumVP int64
	)

	for _, vote := range extCommit.Votes {
		totalVP += vote.Validator.Power

		if err := validateVoteSignatureExistence(vote); err != nil {
			return err
		}

		// Only check + include power if the vote is a commit vote. There must be super-majority, otherwise the
		// previous block (the block vote is for) could not have been committed.
		if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit {
			continue
		}

		// If the validator does not have a valid public key, we skip the signature verification logic but still include
		// the validator's voting power in the total voting power. The app may have pruned the validator's public key
		// from the store, but comet considered the validator as active and included them in the commit since there
		// is a 1 block delay between the validator set update on the app and comet.
		sumVP += vote.Validator.Power
		cmtPubKey, err := veutils.GetValCmtPubKeyFromVote(ctx, vote, valStore)
		if err != nil {
			var notFoundErr *veutils.ValidatorNotFoundError
			if errors.As(err, &notFoundErr) {
				continue
			} else {
				return fmt.Errorf("failed to convert validator: %w", err)
			}
		}
		cve := cmtproto.CanonicalVoteExtension{
			Extension: vote.VoteExtension,
			Height:    currentHeight - 1, // the vote extension was signed in the previous height
			Round:     int64(extCommit.Round),
			ChainId:   chainID,
		}

		extSignBytes, err := veutils.MarshalDelimited(&cve)
		if err != nil {
			return fmt.Errorf("failed to encode CanonicalVoteExtension: %w", err)
		}

		if !cmtPubKey.VerifySignature(extSignBytes, vote.ExtensionSignature) {
			return fmt.Errorf("failed to verify validator %X vote extension signature", cmtPubKey.Address().String())
		}
	}

	// This check is probably unnecessary, but better safe than sorry.
	if totalVP <= 0 {
		return fmt.Errorf("total voting power must be positive, got: %d", totalVP)
	}

	// If the sum of the voting power has not reached (2/3 + 1) we need to error.
	if requiredVP := getRequiredVotingPower(totalVP); sumVP < requiredVP {
		return fmt.Errorf(
			"insufficient cumulative voting power received to verify vote extensions; got: %d, expected: >=%d",
			sumVP, requiredVP,
		)
	}

	return nil
}

// ValidateExtendedCommitAgainstLastCommit validates an ExtendedCommitInfo against a LastCommit. Specifically,
// it checks that the ExtendedCommit + LastCommit (for the same height), are consistent with each other + that
// they are ordered correctly (by voting power) in accordance with
// [comet](https://github.com/cometbft/cometbft/blob/4ce0277b35f31985bbf2c25d3806a184a4510010/types/validator_set.go#L784).
func ValidateExtendedCommitAgainstLastCommit(extCommitInfo cometabci.ExtendedCommitInfo, cmtLastCommit comet.CommitInfo) error {
	// check that the rounds are the same
	if err := validateExtCommitRound(extCommitInfo, cmtLastCommit); err != nil {
		return err
	}

	if err := validateExtCommitVoteCount(extCommitInfo, cmtLastCommit); err != nil {
		return err
	}

	if err := validateVotesSignerInfo(extCommitInfo, cmtLastCommit); err != nil {
		return err
	}

	return nil
}

func validateExtCommitRound(valExtCommitInfo cometabci.ExtendedCommitInfo, cmtLastCommit comet.CommitInfo) error {
	if valExtCommitInfo.Round != cmtLastCommit.Round() {
		return fmt.Errorf(
			"extended commit round %d does not match last commit round %d",
			valExtCommitInfo.Round,
			cmtLastCommit.Round(),
		)
	}
	return nil
}

func getRequiredVotingPower(totalVP int64) int64 {
	return ((totalVP * 2) / 3) + 1
}

func validateVoteSignatureExistence(vote cometabci.ExtendedVoteInfo) error {
	if vote.BlockIdFlag == cmtproto.BlockIDFlagCommit && len(vote.ExtensionSignature) == 0 {
		return fmt.Errorf("vote extension signature is missing; validator addr %s",
			vote.Validator.String(),
		)
	}
	if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit && len(vote.VoteExtension) != 0 {
		return fmt.Errorf("non-commit vote extension present; validator addr %s",
			vote.Validator.String(),
		)
	}
	if vote.BlockIdFlag != cmtproto.BlockIDFlagCommit && len(vote.ExtensionSignature) != 0 {
		return fmt.Errorf("non-commit vote extension signature present; validator addr %s",
			vote.Validator.String(),
		)
	}
	return nil
}

func validateExtCommitVoteCount(
	valExtCommitInfo cometabci.ExtendedCommitInfo,
	cmtLastCommit comet.CommitInfo,
) error {
	if len(valExtCommitInfo.Votes) != cmtLastCommit.Votes().Len() {
		return fmt.Errorf(
			"extended commit votes length %d does not match last commit votes length %d",
			len(valExtCommitInfo.Votes),
			cmtLastCommit.Votes().Len(),
		)
	}
	return nil
}

// GetPubKeyByConsAddr returns the public key of a validator given the consensus addr.
func GetPubKeyByConsAddr(ccvalidator ccvtypes.CrossChainValidator) (cmtprotocrypto.PublicKey, error) {
	consPubKey, err := ccvalidator.ConsPubKey()
	if err != nil {
		return cmtprotocrypto.PublicKey{}, fmt.Errorf("could not get pubkey for val %s: %w", ccvalidator.String(), err)
	}
	tmPubKey, err := cryptocodec.ToCmtProtoPublicKey(consPubKey)
	if err != nil {
		return cmtprotocrypto.PublicKey{}, err
	}

	return tmPubKey, nil
}

func validateVotesSignerInfo(valExtCommitInfo cometabci.ExtendedCommitInfo, cmtLastCommit comet.CommitInfo) error {
	addressCache := make(map[string]struct{}, len(valExtCommitInfo.Votes))
	for i, vote := range valExtCommitInfo.Votes {
		if err := validateVoteAddress(
			vote,
			cmtLastCommit.Votes().Get(i),
			addressCache,
		); err != nil {
			return err
		}
		if err := validateVoteBlockIdFlag(vote, cmtLastCommit.Votes().Get(i)); err != nil {
			return err
		}
	}
	return nil
}

func validateVoteBlockIdFlag(vote cometabci.ExtendedVoteInfo, cmtLastCommitVote comet.VoteInfo) error {
	if !(vote.BlockIdFlag == cmtproto.BlockIDFlagAbsent && len(vote.VoteExtension) == 0 && len(vote.ExtensionSignature) == 0) {
		if int32(vote.BlockIdFlag) != int32(cmtLastCommitVote.GetBlockIDFlag()) {
			return fmt.Errorf(
				"mismatched block ID flag between extended commit vote %d and last proposed commit %d",
				int32(vote.BlockIdFlag),
				int32(cmtLastCommitVote.GetBlockIDFlag()),
			)
		}
	}
	return nil
}

func validateVoteAddress(
	vote cometabci.ExtendedVoteInfo,
	cmtLastCommitVote comet.VoteInfo,
	addressCache map[string]struct{},
) error {
	if _, ok := addressCache[string(vote.Validator.Address)]; ok {
		return fmt.Errorf("extended commit vote address %X is duplicated", vote.Validator.Address)
	}
	addressCache[string(vote.Validator.Address)] = struct{}{}

	if !bytes.Equal(vote.Validator.Address, cmtLastCommitVote.Validator().Address()) {
		return fmt.Errorf(
			"extended commit vote address %X does not match last commit vote address %X",
			vote.Validator.Address,
			cmtLastCommitVote.Validator().Address(),
		)
	}
	if vote.Validator.Power != cmtLastCommitVote.Validator().Power() {
		return fmt.Errorf(
			"extended commit vote power %d does not match last commit vote power %d",
			vote.Validator.Power,
			cmtLastCommitVote.Validator().Power(),
		)
	}
	return nil
}

func GetMaxMarketPairs(ctx sdk.Context, pricesKeeper PreBlockExecPricesKeeper) uint32 {
	markets := pricesKeeper.GetAllMarketParams(ctx)
	return uint32(len(markets))
}
