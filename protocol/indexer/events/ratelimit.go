package events

// NewUpdatePerpetualEventV1 creates a UpdatePerpetualEventV1 representing
// update of a perpetual.
func NewUpdateYieldParamsEventV1(
	sdaiPrice string,
	assetYieldIndex string,
) *UpdateYieldParamsEventV1 {
	return &UpdateYieldParamsEventV1{
		SdaiPrice:       sdaiPrice,
		AssetYieldIndex: assetYieldIndex,
	}
}
