package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateUpdateLeverageMsg(ctx sdk.Context, msg *MsgUpdateLeverage, clobKeeper ClobKeeper) error {
	if msg.SubaccountId == nil {
		return errorsmod.Wrap(ErrInvalidAddress, "subaccount ID cannot be nil")
	}

	if err := msg.SubaccountId.Validate(); err != nil {
		return err
	}

	// Validate that leverage entries are not empty
	if len(msg.ClobPairLeverage) == 0 {
		return errorsmod.Wrap(ErrInvalidLeverage, "clob pair leverage entries cannot be empty")
	}

	// Validate leverage values are positive and clob pair IDs are unique
	clobPairIds := make(map[uint32]bool)
	for _, entry := range msg.ClobPairLeverage {
		if entry == nil {
			return errorsmod.Wrap(ErrInvalidLeverage, "leverage entry cannot be nil")
		}

		if entry.CustomImfPpm == 0 || entry.CustomImfPpm > 1_000_000 {
			return errorsmod.Wrap(
				ErrInvalidLeverage,
				fmt.Sprintf("imf ppm for clob pair %d must be between (0, 1,000,000]", entry.ClobPairId),
			)
		}

		if clobPairIds[entry.ClobPairId] {
			return errorsmod.Wrap(
				ErrInvalidLeverage,
				fmt.Sprintf("duplicate clob pair ID %d", entry.ClobPairId),
			)
		}

		// Validate that the clob pair ID is a valid clob pair ID
		if _, found := clobKeeper.GetClobPair(ctx, ClobPairId(entry.ClobPairId)); !found {
			return errorsmod.Wrap(
				ErrInvalidClob,
				fmt.Sprintf("clob pair ID %d does not exist", entry.ClobPairId),
			)
		}
		clobPairIds[entry.ClobPairId] = true
	}

	return nil
}

func ValidateAndConstructPerpetualLeverageMap(
	ctx sdk.Context,
	msg *MsgUpdateLeverage,
	clobKeeper ClobKeeper,
) (map[uint32]uint32, error) {
	if err := ValidateUpdateLeverageMsg(ctx, msg, clobKeeper); err != nil {
		return nil, err
	}

	perpetualLeverageMap := make(map[uint32]uint32)
	for _, entry := range msg.ClobPairLeverage {
		clob, _ := clobKeeper.GetClobPair(ctx, ClobPairId(entry.ClobPairId))
		perpetualId := clob.MustGetPerpetualId()
		perpetualLeverageMap[perpetualId] = entry.CustomImfPpm
	}

	return perpetualLeverageMap, nil
}
