package int256

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInt(t *testing.T) {
	require.Equal(t, "10", NewInt(10).String())
	require.Equal(t, "-10", NewInt(-10).String())
	require.Equal(t, "0", NewInt(0).String())
}

func TestNewUnsignedInt(t *testing.T) {
	require.Equal(t, NewInt(1), NewUnsignedInt(1))
	require.Equal(t, NewInt(0), NewUnsignedInt(0))
	require.Equal(t, NewInt(123), NewUnsignedInt(123))
}

func TestMustFromBig(t *testing.T) {
	require.Equal(t, NewInt(0), MustFromBig(big.NewInt(0)))
	require.Equal(t, NewInt(123132451), MustFromBig(big.NewInt(123132451)))
	require.Equal(t, NewInt(-123132451), MustFromBig(big.NewInt(-123132451)))
}

func TestToBig(t *testing.T) {
	require.True(t, big.NewInt(0).Cmp(NewInt(0).ToBig()) == 0)
	require.True(t, big.NewInt(10000000).Cmp(NewInt(10000000).ToBig()) == 0)
	require.True(t, big.NewInt(-10000000).Cmp(NewInt(-10000000).ToBig()) == 0)
}

func TestSet(t *testing.T) {
	tests := map[string]struct {
		x *Int
	}{
		"positive": {
			x: NewInt(123),
		},
		"zero": {
			x: NewInt(0),
		},
		"negative": {
			x: NewInt(-456),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			ret := zNew.Set(tc.x)
			require.Equal(t, tc.x, zNew)
			require.Equal(t, tc.x, ret)
			require.Equal(t, xOrig, tc.x)
		})
	}
}

func TestSetUint64(t *testing.T) {
	tests := map[string]struct {
		x uint64
	}{
		"positive": {
			x: 123,
		},
		"zero": {
			x: 0,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			require.Equal(t, zNew.SetUint64(tc.x), NewUnsignedInt(tc.x))
		})
	}
}

func TestSign(t *testing.T) {
	tests := map[string]struct {
		z *Int
		r int
	}{
		"positive": {
			z: NewInt(-12323),
			r: -1,
		},
		"zero": {
			z: NewInt(0),
			r: 0,
		},
		"negative": {
			z: NewInt(451236),
			r: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.r, tc.z.Sign())
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := map[string]struct {
		z *Int
		r bool
	}{
		"positive": {
			z: NewInt(-12323),
			r: false,
		},
		"zero": {
			z: NewInt(0),
			r: true,
		},
		"negative": {
			z: NewInt(451236),
			r: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.r, tc.z.IsZero())
		})
	}
}

func TestEq(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		r bool
	}{
		"equal": {
			z: NewInt(-12323),
			x: NewInt(-12323),
			r: true,
		},
		"not equal": {
			z: NewInt(-12323),
			x: NewInt(12323),
			r: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			require.Equal(t, tc.r, tc.z.Eq(tc.x))
		})
	}
}

func TestCmp(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		r int
	}{
		"pos pos -1": {
			z: NewInt(100),
			x: NewInt(200),
			r: -1,
		},
		"pos pos 0": {
			z: NewInt(100),
			x: NewInt(100),
			r: 0,
		},
		"pos pos 1": {
			z: NewInt(200),
			x: NewInt(100),
			r: 1,
		},
		"zero pos -1": {
			z: NewInt(0),
			x: NewInt(100),
			r: -1,
		},
		"pos zero 1": {
			z: NewInt(100),
			x: NewInt(0),
			r: 1,
		},
		"zero neg 1": {
			z: NewInt(0),
			x: NewInt(-100),
			r: 1,
		},
		"neg zero -1": {
			z: NewInt(-200),
			x: NewInt(0),
			r: -1,
		},
		"pos neg 1": {
			z: NewInt(200),
			x: NewInt(-100),
			r: 1,
		},
		"neg pos -1": {
			z: NewInt(-200),
			x: NewInt(100),
			r: -1,
		},
		"neg neg -1": {
			z: NewInt(-200),
			x: NewInt(-100),
			r: -1,
		},
		"neg neg 0": {
			z: NewInt(-200),
			x: NewInt(-200),
			r: 0,
		},
		"neg neg 1": {
			z: NewInt(-100),
			x: NewInt(-200),
			r: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zOrig := new(Int).Set(tc.z)
			xOrig := new(Int).Set(tc.x)
			require.Equal(t, tc.r, tc.z.Cmp(tc.x))
			require.Equal(t, zOrig, tc.z)
			require.Equal(t, xOrig, tc.x)
		})
	}
}

func TestNeg(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
	}{
		"positive": {
			z: NewInt(-123),
			x: NewInt(123),
		},
		"zero": {
			z: NewInt(0),
			x: NewInt(0),
		},
		"negative": {
			z: NewInt(456),
			x: NewInt(-456),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			ret := zNew.Neg(tc.x)
			require.Equal(t, tc.z, ret)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, xOrig, tc.x)
		})
	}
}

func TestAbs(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
	}{
		"positive": {
			z: NewInt(123),
			x: NewInt(123),
		},
		"zero": {
			z: NewInt(0),
			x: NewInt(0),
		},
		"negative": {
			z: NewInt(456),
			x: NewInt(-456),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			ret := zNew.Abs(tc.x)
			require.Equal(t, tc.z, ret)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, xOrig, tc.x)
		})
	}
}

func TestAdd(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		y *Int
	}{
		"pos neg equals pos": {
			z: NewInt(23),
			x: NewInt(-100),
			y: NewInt(123),
		},
		"pos neg equals neg": {
			z: NewInt(-50),
			x: NewInt(-100),
			y: NewInt(50),
		},
		"neg neg": {
			z: NewInt(-150),
			x: NewInt(-100),
			y: NewInt(-50),
		},
		"pos pos": {
			z: NewInt(150),
			x: NewInt(100),
			y: NewInt(50),
		},
		"zero": {
			z: NewInt(50),
			x: NewInt(0),
			y: NewInt(50),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			yOrig := new(Int).Set(tc.y)
			ret := zNew.Add(tc.x, tc.y)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, tc.z, ret)
			require.Equal(t, xOrig, tc.x)
			require.Equal(t, yOrig, tc.y)
		})
	}
}

func TestSub(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		y *Int
	}{
		"pos neg": {
			z: NewInt(223),
			x: NewInt(100),
			y: NewInt(-123),
		},
		"neg pos": {
			z: NewInt(-150),
			x: NewInt(-100),
			y: NewInt(50),
		},
		"neg neg equals pos": {
			z: NewInt(400),
			x: NewInt(-100),
			y: NewInt(-500),
		},
		"neg neg equals neg": {
			z: NewInt(-50),
			x: NewInt(-100),
			y: NewInt(-50),
		},
		"zero pos": {
			z: NewInt(-50),
			x: NewInt(0),
			y: NewInt(50),
		},
		"pos zero": {
			z: NewInt(50),
			x: NewInt(50),
			y: NewInt(0),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			yOrig := new(Int).Set(tc.y)
			ret := zNew.Sub(tc.x, tc.y)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, tc.z, ret)
			require.Equal(t, xOrig, tc.x)
			require.Equal(t, yOrig, tc.y)
		})
	}
}

func TestMul(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		y *Int
	}{
		"zero zero": {
			z: NewInt(0),
			x: NewInt(0),
			y: NewInt(0),
		},
		"zero pos": {
			z: NewInt(0),
			x: NewInt(0),
			y: NewInt(100),
		},
		"pos zero": {
			z: NewInt(0),
			x: NewInt(123),
			y: NewInt(0),
		},
		"zero neg": {
			z: NewInt(0),
			x: NewInt(0),
			y: NewInt(-123),
		},
		"neg zero": {
			z: NewInt(0),
			x: NewInt(-4634562343265),
			y: NewInt(0),
		},
		"pos neg": {
			z: NewInt(-600),
			x: NewInt(200),
			y: NewInt(-3),
		},
		"neg pos": {
			z: NewInt(-39),
			x: NewInt(-13),
			y: NewInt(3),
		},
		"pos pos": {
			z: NewInt(50),
			x: NewInt(5),
			y: NewInt(10),
		},
		"neg neg": {
			z: NewInt(50),
			x: NewInt(-10),
			y: NewInt(-5),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			yOrig := new(Int).Set(tc.y)
			ret := zNew.Mul(tc.x, tc.y)
			require.Equal(t, tc.z, ret)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, xOrig, tc.x)
			require.Equal(t, yOrig, tc.y)
		})
	}
}

func TestDiv(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		y *Int
	}{
		"div zero": {
			z: NewInt(0),
			x: NewInt(10),
			y: NewInt(0),
		},
		"pos pos remainder": {
			z: NewInt(2),
			x: NewInt(11),
			y: NewInt(5),
		},
		"neg neg remainder": {
			z: NewInt(5),
			x: NewInt(-110),
			y: NewInt(-20),
		},
		"neg pos remainder": {
			z: NewInt(-5),
			x: NewInt(-110),
			y: NewInt(20),
		},
		"pos neg remainder": {
			z: NewInt(-4),
			x: NewInt(150),
			y: NewInt(-32),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			yOrig := new(Int).Set(tc.y)
			ret := zNew.Div(tc.x, tc.y)
			require.Equal(t, tc.z, ret)
			require.Equal(t, tc.z, zNew)
			require.Equal(t, xOrig, tc.x)
			require.Equal(t, yOrig, tc.y)
		})
	}
}

func TestMulExp10(t *testing.T) {
	tests := map[string]struct {
		z *Int
		x *Int
		y int64
	}{
		"zero with zero exponent": {
			z: NewInt(0),
			x: NewInt(0),
			y: 0,
		},
		"zero with neg exponent": {
			z: NewInt(0),
			x: NewInt(0),
			y: 10,
		},
		"zero with pos exponent": {
			z: NewInt(0),
			x: NewInt(0),
			y: -5,
		},
		"pos with zero exponent": {
			z: NewInt(31),
			x: NewInt(31),
			y: 0,
		},
		"pos with neg exponent": {
			z: NewInt(2431),
			x: NewInt(243100),
			y: -2,
		},
		"pos with pos exponent": {
			z: NewInt(234000000),
			x: NewInt(2340),
			y: 5,
		},
		"neg with zero exponent": {
			z: NewInt(-50),
			x: NewInt(-50),
			y: 0,
		},
		"neg with neg exponent": {
			z: NewInt(-3310),
			x: NewInt(-3310000),
			y: -3,
		},
		"neg with pos exponent": {
			z: NewInt(-90000000),
			x: NewInt(-9),
			y: 7,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			zNew := NewInt(99999)
			xOrig := new(Int).Set(tc.x)
			ret := zNew.MulExp10(tc.x, tc.y)
			require.Equal(t, tc.z.String(), ret.String())
			require.Equal(t, tc.z, zNew)
			require.Equal(t, xOrig, tc.x)
		})
	}
}
