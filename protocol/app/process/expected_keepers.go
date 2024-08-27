package process

import (
	"context"
	"math/big"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perptypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

// ProcessClobKeeper defines the expected clob keeper used for `ProcessProposal`.
type ProcessClobKeeper interface {
	RecordMevMetricsIsEnabled() bool
	RecordMevMetrics(
		ctx sdk.Context,
		stakingKeeper ProcessStakingKeeper,
		perpetualKeeper ProcessPerpetualKeeper,
		msgProposedOperations *types.MsgProposedOperations,
	)
}

type ProcessProposalPriceApplier interface {
	ApplyPricesFromVE(ctx sdk.Context, req *abci.RequestFinalizeBlock, writeToCache bool) error
}

// ProcessStakingKeeper defines the expected staking keeper used for `ProcessProposal`.
type ProcessStakingKeeper interface {
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, err error)
}

// ProcessPerpetualKeeper defines the expected perpetual keeper used for `ProcessProposal`.
type ProcessPerpetualKeeper interface {
	MaybeProcessNewFundingTickEpoch(ctx sdk.Context)
	GetSettlementPpm(
		ctx sdk.Context,
		perpetualId uint32,
		quantums *big.Int,
		index *big.Int,
	) (
		bigNetSettlementPpm *big.Int,
		newFundingIndex *big.Int,
		err error,
	)
	GetPerpetual(ctx sdk.Context, id uint32) (val perptypes.Perpetual, err error)
}

type ProcessConsumerKeeper interface {
	GetCCValidator(ctx sdk.Context, addr []byte) (ccvtypes.CrossChainValidator, bool)
}
