package margin_test

import (
	"math/big"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/stretchr/testify/require"
)

func TestRisk_AddInPlace(t *testing.T) {
	tests := map[string]struct {
		a        margin.Risk
		b        margin.Risk
		expected margin.Risk
	}{
		"zero + zero": {
			a:        margin.Risk{},
			b:        margin.Risk{},
			expected: margin.ZeroRisk(),
		},
		"zero + non-zero": {
			a: margin.Risk{},
			b: margin.Risk{
				MMR: big.NewInt(100),
				IMR: big.NewInt(200),
				NC:  big.NewInt(300),
			},
			expected: margin.Risk{
				MMR: big.NewInt(100),
				IMR: big.NewInt(200),
				NC:  big.NewInt(300),
			},
		},
		"non-zero + zero": {
			a: margin.Risk{
				MMR: big.NewInt(100),
				IMR: big.NewInt(200),
				NC:  big.NewInt(300),
			},
			b: margin.Risk{},
			expected: margin.Risk{
				MMR: big.NewInt(100),
				IMR: big.NewInt(200),
				NC:  big.NewInt(300),
			},
		},
		"non-zero + non-zero": {
			a: margin.Risk{
				MMR: big.NewInt(100),
				IMR: big.NewInt(200),
				NC:  big.NewInt(300),
			},
			b: margin.Risk{
				MMR: big.NewInt(50),
				IMR: big.NewInt(100),
				NC:  big.NewInt(150),
			},
			expected: margin.Risk{
				MMR: big.NewInt(150),
				IMR: big.NewInt(300),
				NC:  big.NewInt(450),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.a.AddInPlace(tc.b)
			require.Equal(t, tc.expected, tc.a)
		})
	}
}

func TestRisk_IsInitialCollateralized(t *testing.T) {
	tests := map[string]struct {
		NC       *big.Int
		IMR      *big.Int
		expected bool
	}{
		"NC > IMR": {
			NC:       big.NewInt(200),
			IMR:      big.NewInt(100),
			expected: true,
		},
		"NC = IMR": {
			NC:       big.NewInt(100),
			IMR:      big.NewInt(100),
			expected: true,
		},
		"NC < IMR": {
			NC:       big.NewInt(50),
			IMR:      big.NewInt(100),
			expected: false,
		},
		"NC = 0, IMR = 0": {
			NC:       big.NewInt(0),
			IMR:      big.NewInt(0),
			expected: true,
		},
		"NC < 0, IMR = 0": {
			NC:       big.NewInt(-100),
			IMR:      big.NewInt(0),
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := margin.Risk{
				IMR: tc.IMR,
				NC:  tc.NC,
			}
			result := r.IsInitialCollateralized()
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRisk_IsMaintenanceCollateralized(t *testing.T) {
	tests := map[string]struct {
		NC       *big.Int
		MMR      *big.Int
		expected bool
	}{
		"NC > MMR": {
			NC:       big.NewInt(200),
			MMR:      big.NewInt(100),
			expected: true,
		},
		"NC = MMR": {
			NC:       big.NewInt(100),
			MMR:      big.NewInt(100),
			expected: true,
		},
		"NC < MMR": {
			NC:       big.NewInt(50),
			MMR:      big.NewInt(100),
			expected: false,
		},
		"NC = 0, MMR = 0": {
			NC:       big.NewInt(0),
			MMR:      big.NewInt(0),
			expected: true,
		},
		"NC < 0, MMR = 0": {
			NC:       big.NewInt(-100),
			MMR:      big.NewInt(0),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := margin.Risk{
				MMR: tc.MMR,
				NC:  tc.NC,
			}
			result := r.IsMaintenanceCollateralized()
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestRisk_IsLiquidatable(t *testing.T) {
	tests := map[string]struct {
		NC       *big.Int
		MMR      *big.Int
		expected bool
	}{
		"NC > 0, MMR = 0": {
			NC:       big.NewInt(100),
			MMR:      big.NewInt(0),
			expected: false,
		},
		"NC = 0, MMR > 0": {
			NC:       big.NewInt(0),
			MMR:      big.NewInt(100),
			expected: true,
		},
		"NC = 0, MMR = 0": {
			NC:       big.NewInt(0),
			MMR:      big.NewInt(0),
			expected: false,
		},
		"NC < 0, MMR > 0": {
			NC:       big.NewInt(-100),
			MMR:      big.NewInt(100),
			expected: true,
		},
		"NC < 0, MMR = 0": {
			NC:       big.NewInt(-100),
			MMR:      big.NewInt(0),
			expected: false,
		},
		"NC < MMR": {
			NC:       big.NewInt(75),
			MMR:      big.NewInt(100),
			expected: true,
		},
		"NC = MMR": {
			NC:       big.NewInt(100),
			MMR:      big.NewInt(100),
			expected: false,
		},
		"NC > MMR": {
			NC:       big.NewInt(125),
			MMR:      big.NewInt(100),
			expected: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := margin.Risk{
				MMR: tc.MMR,
				NC:  tc.NC,
			}
			result := r.IsLiquidatable()
			require.Equal(t, tc.expected, result)
		})
	}
}
