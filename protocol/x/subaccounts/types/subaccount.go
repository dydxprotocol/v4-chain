package types

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

const (
	MaxSubaccountIdNumber = 128_000 // 0 ... 128,000 are valid numbers.
)

// BaseQuantums is used to represent an amount in base quantums.
type BaseQuantums uint64

type SafetyHeapPositionSide uint

const (
	Long SafetyHeapPositionSide = iota
	Short
)

// Get the BaseQuantum value in *big.Int.
func (bq BaseQuantums) ToBigInt() *big.Int {
	return new(big.Int).SetUint64(bq.ToUint64())
}

// Get the BaseQuantum value in uint64.
func (bq BaseQuantums) ToUint64() uint64 {
	return uint64(bq)
}

func (m *SubaccountId) Validate() error {
	if _, err := sdk.AccAddressFromBech32(m.Owner); err != nil {
		return errorsmod.Wrapf(ErrInvalidSubaccountIdOwner,
			"invalid SubaccountId Owner address (%s). Error: (%s)", m.Owner, err)
	}

	if m.Number > MaxSubaccountIdNumber {
		return ErrInvalidSubaccountIdNumber
	}

	return nil
}

func (m *SubaccountId) MustGetAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(m.Owner)
}

func (m *Subaccount) DeepCopy() Subaccount {
	b, err := m.Marshal()
	if err != nil {
		panic(err)
	}
	newSubaccount := Subaccount{}
	if err := newSubaccount.Unmarshal(b); err != nil {
		panic(err)
	}
	return newSubaccount
}

// GetPerpetualPositionForId returns the perpetual position with the given
// perpetual id. Returns nil if subaccount does not have an open position
// for the perpetual.
func (m *Subaccount) GetPerpetualPositionForId(
	perpetualId uint32,
) (
	perpetualPosition *PerpetualPosition,
	exists bool,
) {
	if m != nil {
		for _, position := range m.PerpetualPositions {
			if position.PerpetualId == perpetualId {
				return position, true
			}
		}
	}
	return nil, false
}

// GetUsdcPosition returns the balance of the USDC asset position.
func (m *Subaccount) GetUsdcPosition() *big.Int {
	usdcAssetPosition := m.getUsdcAssetPosition()
	if usdcAssetPosition == nil {
		return new(big.Int)
	}
	return usdcAssetPosition.GetBigQuantums()
}

// SetUsdcAssetPosition sets the balance of the USDC asset position to `newUsdcPosition`.
func (m *Subaccount) SetUsdcAssetPosition(newUsdcPosition *big.Int) {
	if m == nil {
		return
	}

	usdcAssetPosition := m.getUsdcAssetPosition()
	if newUsdcPosition == nil || newUsdcPosition.Sign() == 0 {
		if usdcAssetPosition != nil {
			m.AssetPositions = m.AssetPositions[1:]
		}
	} else {
		if usdcAssetPosition == nil {
			usdcAssetPosition = &AssetPosition{
				AssetId: assettypes.AssetUsdc.Id,
			}
			m.AssetPositions = append([]*AssetPosition{usdcAssetPosition}, m.AssetPositions...)
		}
		usdcAssetPosition.Quantums = dtypes.NewIntFromBigInt(newUsdcPosition)
	}
}

func (m *Subaccount) getUsdcAssetPosition() *AssetPosition {
	if m == nil || len(m.AssetPositions) == 0 {
		return nil
	}

	firstAsset := m.AssetPositions[0]
	if firstAsset.AssetId != assettypes.AssetUsdc.Id {
		return nil
	}
	return firstAsset
}
