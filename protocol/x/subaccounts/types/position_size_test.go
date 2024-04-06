package types_test

import (
	"testing"

	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestGetQuantums(t *testing.T) {
	p := types.PerpetualPosition{
		Quantums: dtypes.NewInt(42),
	}

	require.Equal(t, int256.NewInt(42), p.GetQuantums())

	p = types.PerpetualPosition{
		Quantums: dtypes.NewInt(-42),
	}

	require.Equal(t, int256.NewInt(-42), p.GetQuantums())

	nilPosition := (*types.PerpetualPosition)(nil)
	require.Equal(t, int256.NewInt(0), nilPosition.GetQuantums())
}

func TestPerpetualPosition_GetIsLong(t *testing.T) {
	longPosition := types.PerpetualPosition{
		Quantums: dtypes.NewInt(1000),
	}
	zeroPosition := types.PerpetualPosition{
		Quantums: dtypes.NewInt(0),
	}
	shortPosition := types.PerpetualPosition{
		Quantums: dtypes.NewInt(-10000000),
	}
	nilPosition := (*types.PerpetualPosition)(nil)

	require.True(t,
		longPosition.GetIsLong(),
	)
	require.PanicsWithError(t,
		errorsmod.Wrapf(
			types.ErrPerpPositionZeroQuantum,
			"perpetual position (perpetual Id: 0) has zero quantum",
		).Error(),
		func() {
			zeroPosition.GetIsLong()
		},
	)
	require.False(t,
		shortPosition.GetIsLong(),
	)
	require.False(t,
		nilPosition.GetIsLong(),
	)
}

func TestAssetPosition_GetIsLong(t *testing.T) {
	longPosition := types.AssetPosition{
		Quantums: dtypes.NewInt(1000),
	}
	zeroPosition := types.AssetPosition{
		Quantums: dtypes.NewInt(0),
	}
	shortPosition := types.AssetPosition{
		Quantums: dtypes.NewInt(-10000000),
	}
	nilPosition := (*types.AssetPosition)(nil)

	require.True(t,
		longPosition.GetIsLong(),
	)
	require.PanicsWithError(t,
		errorsmod.Wrapf(
			types.ErrAssetPositionZeroQuantum,
			"asset position (asset Id: 0) has zero quantum",
		).Error(),
		func() {
			zeroPosition.GetIsLong()
		},
	)
	require.False(t,
		shortPosition.GetIsLong(),
	)
	require.False(t,
		nilPosition.GetIsLong(),
	)
}

func TestAssetUpdate_GetIsLong(t *testing.T) {
	longUpdate := types.AssetUpdate{
		QuantumsDelta: int256.NewInt(1000),
	}
	zeroUpdate := types.AssetUpdate{
		QuantumsDelta: int256.NewInt(0),
	}
	shortUpdate := types.AssetUpdate{
		QuantumsDelta: int256.NewInt(-10000000),
	}

	require.True(t,
		longUpdate.GetIsLong(),
	)
	require.False(t,
		zeroUpdate.GetIsLong(),
	)
	require.False(t,
		shortUpdate.GetIsLong(),
	)
}

func TestPerpetualUpdate_GetIsLong(t *testing.T) {
	longUpdate := types.PerpetualUpdate{
		QuantumsDelta: int256.NewInt(1000),
	}
	zeroUpdate := types.PerpetualUpdate{
		QuantumsDelta: int256.NewInt(0),
	}
	shortUpdate := types.PerpetualUpdate{
		QuantumsDelta: int256.NewInt(-10000000),
	}

	require.True(t,
		longUpdate.GetIsLong(),
	)
	require.False(t,
		zeroUpdate.GetIsLong(),
	)
	require.False(t,
		shortUpdate.GetIsLong(),
	)
}
