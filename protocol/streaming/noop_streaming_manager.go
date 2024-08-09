package streaming

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/streaming/types"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

var _ types.FullNodeStreamingManager = (*NoopGrpcStreamingManager)(nil)

type NoopGrpcStreamingManager struct{}

func NewNoopGrpcStreamingManager() *NoopGrpcStreamingManager {
	return &NoopGrpcStreamingManager{}
}

func (sm *NoopGrpcStreamingManager) Enabled() bool {
	return false
}

func (sm *NoopGrpcStreamingManager) Subscribe(
	_ []uint32,
	_ []*satypes.SubaccountId,
	_ types.OutgoingMessageSender,
) (
	err error,
) {
	return types.ErrNotImplemented
}

func (sm *NoopGrpcStreamingManager) SendOrderbookUpdates(
	updates *clobtypes.OffchainUpdates,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) SendOrderbookFillUpdates(
	orderbookFills []clobtypes.StreamOrderbookFill,
	blockHeight uint32,
	execMode sdk.ExecMode,
	perpetualIdToClobPairId map[uint32][]clobtypes.ClobPairId,
) {
}

func (sm *NoopGrpcStreamingManager) SendSubaccountUpdates(
	subaccountUpdates []satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) SendTakerOrderStatus(
	takerOrder clobtypes.StreamTakerOrder,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) InitializeNewStreams(
	getOrderbookSnapshot func(clobPairId clobtypes.ClobPairId) *clobtypes.OffchainUpdates,
	getSubaccountSnapshot func(subaccountId satypes.SubaccountId) *satypes.StreamSubaccountUpdate,
	blockHeight uint32,
	execMode sdk.ExecMode,
) {
}

func (sm *NoopGrpcStreamingManager) Stop() {
}
