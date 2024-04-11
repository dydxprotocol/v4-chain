package int256

import (
	"math/big"

	"github.com/holiman/uint256"
)

// Signed wrapper for github.com/holiman/uint256.
//
// WARNING: Do not write to pointers until reading from all pointers is complete. This will cause incorrect
// behavior if the same pointer is passed in through multiple arguments.

type Int uint256.Int

var (
	OneInt256        = NewInt(1)
	TenInt256        = NewInt(10)
	OneMillionInt256 = NewInt(1_000_000)

	exp10Lookup = createExp10Lookup()
)

func (z *Int) String() string {
	if z.Sign() >= 0 {
		return (*uint256.Int)(z).Dec()
	}
	return "-" + (*uint256.Int)(new(uint256.Int).Neg((*uint256.Int)(z))).Dec()
}

// NewInt creates a new Int from an int64.
func NewInt(val int64) *Int {
	if val < 0 {
		u := uint256.NewInt(uint64(-val))
		return (*Int)(u.Neg(u))
	}
	return (*Int)(uint256.NewInt(uint64(val)))
}

// NewUnsignedInt creates a new Int from a uint64.
func NewUnsignedInt(val uint64) *Int {
	return (*Int)(uint256.NewInt(uint64(val)))
}

// MustFromBig creates a new Int from a big.Int. Panics on failure.
func MustFromBig(b *big.Int) *Int {
	return (*Int)(uint256.MustFromBig(b))
}

// ToBig converts z to a big.Int.
func (z *Int) ToBig() *big.Int {
	if z.Sign() >= 0 {
		return (*uint256.Int)(z).ToBig()
	}
	r := new(uint256.Int).Neg((*uint256.Int)(z)).ToBig()
	return r.Neg(r)
}

// Set sets z to the value of x.
func (z *Int) Set(x *Int) *Int {
	return (*Int)((*uint256.Int)(z).Set((*uint256.Int)(x)))
}

// Set sets z to the value of a uint256.
func (z *Int) SetUint64(x uint64) *Int {
	return (*Int)((*uint256.Int)(z).SetUint64(x))
}

// Sign returns -1 if z is negative, 0 if z is zero, and 1 if z is positive.
func (z *Int) Sign() int {
	return (*uint256.Int)(z).Sign()
}

// IsZero returns true iff z is equal to 0.
func (z *Int) IsZero() bool {
	return (*uint256.Int)(z).IsZero()
}

// Eq returns z == x.
func (z *Int) Eq(x *Int) bool {
	return (*uint256.Int)(z).Eq((*uint256.Int)(x))
}

// Cmp returns -1 if z < x, 0 if z == x, and +1 if z > x.
func (z *Int) Cmp(x *Int) (r int) {
	if z.Sign() >= 0 {
		if x.Sign() >= 0 {
			return (*uint256.Int)(z).Cmp((*uint256.Int)(x))
		} else {
			return 1
		}
	}
	if x.Sign() >= 0 {
		return -1
	}
	return (*uint256.Int)(z).Cmp((*uint256.Int)(x))
}

// Neg sets z to -x and returns z.
func (z *Int) Neg(x *Int) *Int {
	return (*Int)((*uint256.Int)(z).Neg((*uint256.Int)(x)))
}

// Abs sets z to the absolute value of x and returns z.
func (z *Int) Abs(x *Int) *Int {
	return (*Int)((*uint256.Int)(z).Abs((*uint256.Int)(x)))
}

// Add sets z = x + y and returns z.
func (z *Int) Add(x, y *Int) *Int {
	return (*Int)((*uint256.Int)(z).Add((*uint256.Int)(x), (*uint256.Int)(y)))
}

// Sub sets z = x - y and returns z.
func (z *Int) Sub(x, y *Int) *Int {
	return (*Int)((*uint256.Int)(z).Sub((*uint256.Int)(x), (*uint256.Int)(y)))
}

// Mul sets z = x * y and returns z.
func (z *Int) Mul(x, y *Int) *Int {
	if x.Sign() > 0 {
		if y.Sign() > 0 {
			return (*Int)(
				(*uint256.Int)(z).Mul(
					(*uint256.Int)(x),
					(*uint256.Int)(y),
				),
			)
		} else {
			return z.Neg((*Int)(
				(*uint256.Int)(z).Mul(
					(*uint256.Int)(x),
					new(uint256.Int).Neg((*uint256.Int)(y)),
				),
			))
		}
	}
	if y.Sign() > 0 {
		return z.Neg((*Int)(
			(*uint256.Int)(z).Mul(
				new(uint256.Int).Neg((*uint256.Int)(x)),
				(*uint256.Int)(y),
			),
		))
	}
	return (*Int)((*uint256.Int)(z).Mul(
		new(uint256.Int).Neg((*uint256.Int)(x)),
		new(uint256.Int).Neg((*uint256.Int)(y)),
	))
}

// Div sets z = x / y and returns z. If y is 0, z is set to 0.
func (z *Int) Div(x, y *Int) *Int {
	return (*Int)((*uint256.Int)(z).SDiv((*uint256.Int)(x), (*uint256.Int)(y)))
}

func createExp10Lookup() map[uint64]uint256.Int {
	lookup := make(map[uint64]uint256.Int, 100)
	value := uint256.NewInt(1)
	for i := 0; i < 100; i++ {
		lookup[uint64(i)] = *new(uint256.Int).Set(value)
		value.Mul(value, (*uint256.Int)(TenInt256))
	}
	return lookup
}

func mulExp10(z *uint256.Int, x *uint256.Int, y int64) *uint256.Int {
	var abs uint64
	if y < 0 {
		abs = uint64(-y)
	} else {
		abs = uint64(y)
	}
	lookup, ok := exp10Lookup[abs]
	var exp10 *uint256.Int
	if ok {
		exp10 = &lookup
	} else {
		exp10.Exp((*uint256.Int)(TenInt256), uint256.NewInt(abs))
	}
	if y < 0 {
		return z.Div(x, exp10)
	} else {
		return z.Mul(x, exp10)
	}
}

// MulExp10 sets z = x * 10^y and returns z.
func (z *Int) MulExp10(x *Int, y int64) *Int {
	if x.Sign() >= 0 {
		return (*Int)(mulExp10((*uint256.Int)(z), (*uint256.Int)(x), y))
	}
	return z.Neg((*Int)(mulExp10((*uint256.Int)(z), (*uint256.Int)(z.Neg(x)), y)))
}
