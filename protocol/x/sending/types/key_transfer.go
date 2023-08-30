package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// TransferToGoodTilBlockKeyPrefix is the prefix to retrieve all transfer hash to `GoodTilBlock` mappings.
	TransferToGoodTilBlockKeyPrefix = "TransferToGoodTilBlock/value/"
	// BlockExpirationForTransfersKeyPrefix is the prefix to retrieve all expiring transfers given a block height.
	BlockExpirationForTransfersKeyPrefix = "BlockExpirationForTransfers/value/"
)
