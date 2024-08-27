package prepare

import (
	clobtypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/clob/types"
	perpstypes "github.com/StreamFinance-Protocol/stream-chain/protocol/x/perpetuals/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ccvtypes "github.com/ethos-works/ethos/ethos-chain/x/ccv/consumer/types"
)

// PrepareClobKeeper defines the expected CLOB keeper used for `PrepareProposal`.
type PrepareClobKeeper interface {
	GetOperations(ctx sdk.Context) *clobtypes.MsgProposedOperations
}

// PreparePerpetualsKeeper defines the expected Perpetuals keeper used for `PrepareProposal`.
type PreparePerpetualsKeeper interface {
	GetAddPremiumVotes(ctx sdk.Context) *perpstypes.MsgAddPremiumVotes
}

type PrepareConsumerKeeper interface {
	GetCCValidator(ctx sdk.Context, addr []byte) (ccvtypes.CrossChainValidator, bool)
}
