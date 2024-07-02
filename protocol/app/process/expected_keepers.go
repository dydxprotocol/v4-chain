package process

import (
	"context"

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
	GetValidatorByConsAddr(ctx context.Context, consAddr sdk.ConsAddress) (validator stakingtypes.Validator, err error)
}

// ProcessPerpetualKeeper defines the expected perpetual keeper used for `ProcessProposal`.
type ProcessPerpetualKeeper interface {
	MaybeProcessNewFundingTickEpoch(ctx sdk.Context)
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
