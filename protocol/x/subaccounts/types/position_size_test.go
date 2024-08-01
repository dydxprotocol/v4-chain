package types_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"

	"github.com/stretchr/testify/require"
)

func TestGetBigQuantums(t *testing.T) {
	p := types.PerpetualPosition{
		Quantums: dtypes.NewInt(42),
	}

	require.Equal(t, big.NewInt(42), p.GetBigQuantums())

	p = types.PerpetualPosition{
		Quantums: dtypes.NewInt(-42),
	}

	require.Equal(t, big.NewInt(-42), p.GetBigQuantums())

	nilPosition := (*types.PerpetualPosition)(nil)
	require.Equal(t, big.NewInt(0), nilPosition.GetBigQuantums())
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
	require.False(t,
		zeroPosition.GetIsLong(),
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
	require.False(t,
		zeroPosition.GetIsLong(),
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
		BigQuantumsDelta: big.NewInt(1000),
	}
	zeroUpdate := types.AssetUpdate{
		BigQuantumsDelta: big.NewInt(0),
	}
	shortUpdate := types.AssetUpdate{
		BigQuantumsDelta: big.NewInt(-10000000),
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
		BigQuantumsDelta: big.NewInt(1000),
	}
	zeroUpdate := types.PerpetualUpdate{
		BigQuantumsDelta: big.NewInt(0),
	}
	shortUpdate := types.PerpetualUpdate{
		BigQuantumsDelta: big.NewInt(-10000000),
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
