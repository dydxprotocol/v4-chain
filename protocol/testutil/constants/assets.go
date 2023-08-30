package constants

import (
	"math"
	"math/big"

	asstypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

const (
	// TestFee is the gas fee offered for test transactions.
	TestFee = "50000" + asstypes.UusdcDenom // 5 cents
	// TestGasLimit is the gas limit used for test transactions.
	// It's set to a larger amount such that the transaction never runs out of gas.
	TestGasLimit = 1_000_000
)

// BigNegMaxUint64 returns a `big.Int` that is set to -math.MaxUint64.
func BigNegMaxUint64() *big.Int {
	return new(big.Int).Neg(
		new(big.Int).SetUint64(math.MaxUint64),
	)
}

var (
	BtcUsd = &asstypes.Asset{
		Id:               1,
		Symbol:           "BTC",
		Denom:            "btc-denom",
		DenomExponent:    int32(-8),
		HasMarket:        true,
		MarketId:         uint32(0),
		AtomicResolution: int32(-8),
	}

	Usdc = &asstypes.Asset{
		Id:               0,
		Symbol:           "USDC",
		Denom:            asstypes.AssetUsdc.Denom,
		DenomExponent:    int32(-6),
		HasMarket:        false,
		MarketId:         uint32(0),
		AtomicResolution: int32(-6),
	}
)
