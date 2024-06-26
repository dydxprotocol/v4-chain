package aggregator

import (
	"log"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/app/ve/codec"
	pk "github.com/StreamFinance-Protocol/stream-chain/protocol/x/prices/keeper"
)

// PriceWriter is an interface that defines the methods required to aggregate and apply prices from VE's
type PriceWriter struct {
	// va is a VoteAggregator that is used to aggregate votes into prices.
	va VoteAggregator

	// pk is the prices keeper that is used to write prices to state.
	pk pk.Keeper

	// logger
	logger log.Logger

	// codecs
	voteExtensionCodec  codec.VoteExtensionCodec
	extendedCommitCodec codec.ExtendedCommitCodec
}
