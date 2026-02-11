package prepare

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bridgetypes "github.com/dydxprotocol/v4-chain/protocol/x/bridge/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	perpstypes "github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// PrepareClobKeeper defines the expected CLOB keeper used for `PrepareProposal`.
type PrepareClobKeeper interface {
	CancelShortTermOrder(ctx sdk.Context, msg *clobtypes.MsgCancelOrder) error
	CancelStatefulOrder(ctx sdk.Context, msg *clobtypes.MsgCancelOrder) error
	BatchCancelShortTermOrder(ctx sdk.Context, msg *clobtypes.MsgBatchCancel) (success []uint32, failure []uint32, err error)
	PlaceShortTermOrder(ctx sdk.Context, msg *clobtypes.MsgPlaceOrder) (satypes.BaseQuantums, clobtypes.OrderStatus, error)
	PlaceStatefulOrder(ctx sdk.Context, msg *clobtypes.MsgPlaceOrder, isInternalOrder bool) error
	GetOperations(ctx sdk.Context) *clobtypes.MsgProposedOperations
}

// PreparePerpetualsKeeper defines the expected Perpetuals keeper used for `PrepareProposal`.
type PreparePerpetualsKeeper interface {
	GetAddPremiumVotes(ctx sdk.Context) *perpstypes.MsgAddPremiumVotes
}

// PrepareBridgeKeeper defines the expected Bridge keeper used for `PrepareProposal`.
type PrepareBridgeKeeper interface {
	GetAcknowledgeBridges(ctx sdk.Context, blockTimestamp time.Time) *bridgetypes.MsgAcknowledgeBridges
}
