package types

import (
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	"github.com/dydxprotocol/v4-chain/protocol/lib/int256"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"
)

const (
	MaxSubaccountIdNumber = 128_000 // 0 ... 128,000 are valid numbers.
)

// BaseQuantums is used to represent an amount in base quantums.
type BaseQuantums uint64

// Get the BaseQuantum value in *int256.Int.
func (bq BaseQuantums) ToInt256() *int256.Int {
	return int256.NewUnsignedInt(uint64(bq))
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
func (m *Subaccount) GetUsdcPosition() *int256.Int {
	usdcAssetPosition := m.getUsdcAssetPosition()
	if usdcAssetPosition == nil {
		return new(int256.Int)
	}
	return usdcAssetPosition.GetQuantums()
}

// SetUsdcAssetPosition sets the balance of the USDC asset position to `newUsdcPosition`.
func (m *Subaccount) SetUsdcAssetPosition(newUsdcPosition *int256.Int) {
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
		usdcAssetPosition.Quantums = dtypes.NewIntFromBigInt(newUsdcPosition.ToBig())
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
