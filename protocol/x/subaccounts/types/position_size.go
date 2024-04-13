package types

import (
	errorsmod "cosmossdk.io/errors"

	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
)

const (
	AssetProductType     = "asset"
	PerpetualProductType = "perpetual"
	UnknownProductTYpe   = "unknown"
)

// PositionSize is an interface for expressing the size of a position
type PositionSize interface {
	// Returns true if and only if the position size is positive.
	GetIsLong() bool
	// Returns the signed position size in in256.Int.
	GetQuantums() *int256.Int
	GetId() uint32
	GetProductType() string
}

type PositionUpdate struct {
	Id       uint32
	Quantums *int256.Int
}

func NewPositionUpdate(id uint32) PositionUpdate {
	return PositionUpdate{
		Id:       id,
		Quantums: int256.NewInt(0),
	}
}

// Both updates and positions should conform to this interface
var _ PositionSize = AssetUpdate{}
var _ PositionSize = PerpetualUpdate{}
var _ PositionSize = PositionUpdate{}

// AssetPositions and PerpetualPositions use pointer receivers
// due to the way proto-gen generates them.
var _ PositionSize = &AssetPosition{}
var _ PositionSize = &PerpetualPosition{}

func (m *AssetPosition) GetId() uint32 {
	return m.GetAssetId()
}

// Get the asset position quantum size in int256.Int. Panics if the size is zero.
func (m *AssetPosition) GetQuantums() *int256.Int {
	if m == nil {
		return new(int256.Int)
	}

	if m.Quantums.BigInt().Sign() == 0 {
		panic(errorsmod.Wrapf(
			ErrAssetPositionZeroQuantum,
			"asset position (asset Id: %v) has zero quantum",
			m.AssetId,
		))
	}

	return int256.MustFromBig(m.Quantums.BigInt())
}

func (m *AssetPosition) GetIsLong() bool {
	if m == nil {
		return false
	}
	return m.GetQuantums().Sign() > 0
}

func (m *AssetPosition) GetProductType() string {
	return AssetProductType
}

func (m *PerpetualPosition) GetId() uint32 {
	return m.GetPerpetualId()
}

func (m *PerpetualPosition) SetQuantums(sizeQuantums int64) {
	m.Quantums = dtypes.NewInt(sizeQuantums)
}

// Get the perpetual position quantum size in int256.Int. Panics if the size is zero.
func (m *PerpetualPosition) GetQuantums() *int256.Int {
	if m == nil {
		return new(int256.Int)
	}

	if m.Quantums.BigInt().Sign() == 0 {
		panic(errorsmod.Wrapf(
			ErrPerpPositionZeroQuantum,
			"perpetual position (perpetual Id: %v) has zero quantum",
			m.PerpetualId,
		))
	}

	return int256.MustFromBig(m.Quantums.BigInt())
}

func (m *PerpetualPosition) GetIsLong() bool {
	if m == nil {
		return false
	}
	return m.GetQuantums().Sign() > 0
}

func (m *PerpetualPosition) GetProductType() string {
	return PerpetualProductType
}

func (au AssetUpdate) GetIsLong() bool {
	return au.GetQuantums().Sign() > 0
}

func (au AssetUpdate) GetQuantums() *int256.Int {
	return au.QuantumsDelta
}

func (au AssetUpdate) GetId() uint32 {
	return au.AssetId
}

func (au AssetUpdate) GetProductType() string {
	return AssetProductType
}

func (pu PerpetualUpdate) GetQuantums() *int256.Int {
	return pu.QuantumsDelta
}

func (pu PerpetualUpdate) GetId() uint32 {
	return pu.PerpetualId
}

func (pu PerpetualUpdate) GetIsLong() bool {
	return pu.GetQuantums().Sign() > 0
}

func (pu PerpetualUpdate) GetProductType() string {
	return PerpetualProductType
}

func (pu PositionUpdate) GetId() uint32 {
	return pu.Id
}

func (pu PositionUpdate) GetIsLong() bool {
	return pu.Quantums.Sign() > 0
}

func (pu PositionUpdate) SetQuantums(Quantums *int256.Int) {
	pu.Quantums.Set(Quantums)
}

func (pu PositionUpdate) GetQuantums() *int256.Int {
	return pu.Quantums
}
func (pu PositionUpdate) GetProductType() string {
	// PositionUpdate is generic and doesn't have a product type.
	return UnknownProductTYpe
}
