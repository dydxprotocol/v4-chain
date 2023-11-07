package process

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perptypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// ProcessPricesKeeper defines the expected Prices keeper used for `ProcessProposal`.
type ProcessPricesKeeper interface {
	PerformStatefulPriceUpdateValidation(
		ctx sdk.Context,
		marketPriceUpdates *pricestypes.MsgUpdateMarketPrices,
		performNonDeterministicValidation bool,
	) error

	UpdateSmoothedPrices(
		ctx sdk.Context,
		linearInterpolateFunc func(v0 uint64, v1 uint64, ppm uint32) (uint64, error),
	) error
}

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

// ProcessStakingKeeper defines the expected staking keeper used for `ProcessProposal`.
type ProcessStakingKeeper interface {
	GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, found bool)
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

// ProcessBridgeKeeper defines the expected bridge keeper used for `ProcessProposal`.
type ProcessBridgeKeeper interface {
	GetAcknowledgedEventInfo(
		ctx sdk.Context,
	) (acknowledgedEventInfo bridgetypes.BridgeEventInfo)
	GetRecognizedEventInfo(
		ctx sdk.Context,
	) (recognizedEventInfo bridgetypes.BridgeEventInfo)
	GetBridgeEventFromServer(ctx sdk.Context, id uint32) (event bridgetypes.BridgeEvent, found bool)
	GetSafetyParams(ctx sdk.Context) (safetyParams bridgetypes.SafetyParams)
}
