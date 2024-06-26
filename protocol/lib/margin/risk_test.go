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

func TestRisk_Cmp(t *testing.T) {
	tests := map[string]struct {
		firstNC  *big.Int
		firstMMR *big.Int

		secondNC  *big.Int
		secondMMR *big.Int

		expected int
	}{
		// Normal cases: different ratios.
		"first is less risky than second": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"first is less risky than second - second has zero TNC": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(0),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"first is less risky than second - second has negative TNC": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(-100),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"first is less risky than second - both have negative TNC": {
			firstNC:   big.NewInt(-100),
			firstMMR:  big.NewInt(150),
			secondNC:  big.NewInt(-100),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"first is more risky than second": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(150),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  1,
		},
		"first is more risky than second - first has zero TNC": {
			firstNC:   big.NewInt(0),
			firstMMR:  big.NewInt(150),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  1,
		},
		"first is more risky than second - first has negative TNC": {
			firstNC:   big.NewInt(-100),
			firstMMR:  big.NewInt(150),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  1,
		},
		"first is more risky than second - both hahave negative TNC": {
			firstNC:   big.NewInt(-100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(-100),
			secondMMR: big.NewInt(100),
			expected:  1,
		},
		"first is less risky than second - first has zero MMR": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(0),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"first is more risky than second - second has zero MMR": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(100),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(0),
			expected:  1,
		},
		// special cases: ratio is the same
		"special case: equally risky": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(100),
			secondNC:  big.NewInt(100),
			secondMMR: big.NewInt(100),
			expected:  0,
		},
		"special case: same ratio, tie break by MMR": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(200),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"special case: first with zero NC and second with zero NC, tie break by MMR": {
			firstNC:   big.NewInt(0),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(0),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"special case: first with zero NC and MMR, tie break by MMR": {
			firstNC:   big.NewInt(0),
			firstMMR:  big.NewInt(0),
			secondNC:  big.NewInt(200),
			secondMMR: big.NewInt(100),
			expected:  -1,
		},
		"special case: second with zero NC and MMR, tie break by MMR": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(50),
			secondNC:  big.NewInt(0),
			secondMMR: big.NewInt(0),
			expected:  1,
		},
		"special case: first with zero MMR and second with zero MMR, tie break by NC": {
			firstNC:   big.NewInt(100),
			firstMMR:  big.NewInt(0),
			secondNC:  big.NewInt(200),
			secondMMR: big.NewInt(0),
			expected:  1,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			firstRisk := margin.Risk{
				MMR: tc.firstMMR,
				NC:  tc.firstNC,
			}
			secondRisk := margin.Risk{
				MMR: tc.secondMMR,
				NC:  tc.secondNC,
			}
			require.Equal(t, tc.expected, firstRisk.Cmp(secondRisk))
		})
	}
}
