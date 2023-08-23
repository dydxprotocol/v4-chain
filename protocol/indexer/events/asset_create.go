package events

// NewAssetCreateEvent creates a AssetCreateEvent representing creation of an asset.
func NewAssetCreateEvent(
	id uint32,
	symbol string,
	hasMarket bool,
	marketId uint32,
	atomicResolution int32,
) *AssetCreateEventV1 {
	return &AssetCreateEventV1{
		Id:               id,
		Symbol:           symbol,
		HasMarket:        hasMarket,
		MarketId:         marketId,
		AtomicResolution: atomicResolution,
	}
}
