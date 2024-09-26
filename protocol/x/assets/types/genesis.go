package types

import "github.com/StreamFinance-Protocol/stream-chain/protocol/lib"

const (
	// TDai denom is defined as equal to the Uusdc denom.
	// utdai stands for micro TDAI, where 1 TDAI = 10^6 utdai.
	TDaiDenom         = "utdai"
	TDaiDenomExponent = -6
)

var (
	AssetTDai Asset = Asset{
		Id:               0,
		Symbol:           "TDAI",
		DenomExponent:    TDaiDenomExponent,
		Denom:            TDaiDenom,
		HasMarket:        false,
		AtomicResolution: lib.QuoteCurrencyAtomicResolution,
		AssetYieldIndex:  "1/1",
	}
)

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Assets: []Asset{
			AssetTDai,
		},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Genesis state should contain at least one asset.
	if len(gs.Assets) == 0 {
		return ErrNoAssetInGenesis
	}

	// The first asset should always be TDai.
	if gs.Assets[0] != AssetTDai {
		return ErrTDaiMustBeAssetZero
	}

	// Provided assets should not contain duplicated asset ids, and denoms.
	// Asset ids should be sequential.
	// MarketId should be 0 if HasMarket is false.
	assetIdSet := make(map[uint32]struct{})
	denomSet := make(map[string]struct{})
	expectedId := uint32(0)

	for _, asset := range gs.Assets {
		if _, exists := assetIdSet[asset.Id]; exists {
			return ErrAssetIdAlreadyExists
		}
		if _, exists := denomSet[asset.Denom]; exists {
			return ErrAssetDenomAlreadyExists
		}
		if asset.Id != expectedId {
			return ErrGapFoundInAssetId
		}
		if !asset.HasMarket && asset.MarketId > 0 {
			return ErrInvalidMarketId
		}
		assetIdSet[asset.Id] = struct{}{}
		denomSet[asset.Denom] = struct{}{}
		expectedId = expectedId + 1
	}
	return nil
}
