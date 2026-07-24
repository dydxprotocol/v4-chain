package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Configurable upper bounds for the governance-tunable conditional-order trigger config.
//
// These bound the values a single governance proposal can set: without an upper bound the setter
// would only normalize zero values, so a proposal could set a per-block budget of MaxUint32 —
// erasing the O(budget) EndBlocker work bound that this whole feature exists to provide. The
// bounds are generous relative to observed legitimate peaks (tens of triggers per block) while
// keeping worst-case per-block work firmly bounded.
const (
	// MaxConfigurableTriggersPerBlock bounds MaxTriggersPerBlock. Even at the ceiling this keeps
	// per-block trigger work bounded and far below the attack magnitudes (~250k–1.2M) that motivated
	// the mitigation.
	MaxConfigurableTriggersPerBlock uint32 = 100_000
	// MaxConfigurableRemovalsPerBlock bounds MaxRemovalsPerBlock (expiry drain per block).
	MaxConfigurableRemovalsPerBlock uint32 = 100_000
	// MaxConfigurableUntriggeredGlobal bounds the global admission cap. Steady-state admission is
	// incremental (per-placement), so this can be generous; the one-shot enable-time backfill is
	// bounded separately by the enable ceiling in the msg handler.
	MaxConfigurableUntriggeredGlobal uint32 = 10_000_000
	// MaxConfigurableUntriggeredPerSubaccount bounds the per-subaccount admission cap.
	MaxConfigurableUntriggeredPerSubaccount uint32 = 1_000_000
)

// ValidateBasic performs stateless validation of a conditional-order trigger config update.
// It rejects an invalid authority and any budget / admission cap outside its configurable bound.
// A zero value for any numeric field is permitted: the keeper setter normalizes zeros to sane
// package defaults, so zero means "use the default", not "unbounded".
//
// Stateful checks (e.g. refusing to enable while the resting untriggered set is too large to
// backfill safely) live in the message handler, which has access to consensus state.
func (msg *MsgUpdateConditionalOrderTriggerConfig) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(
			ErrInvalidAuthority,
			"authority '%s' must be a valid bech32 address, but got error '%v'",
			msg.Authority,
			err.Error(),
		)
	}

	if msg.MaxTriggersPerBlock > MaxConfigurableTriggersPerBlock {
		return errorsmod.Wrapf(
			ErrInvalidConditionalOrderTriggerConfig,
			"maxTriggersPerBlock %d exceeds the configurable upper bound %d",
			msg.MaxTriggersPerBlock,
			MaxConfigurableTriggersPerBlock,
		)
	}
	if msg.MaxRemovalsPerBlock > MaxConfigurableRemovalsPerBlock {
		return errorsmod.Wrapf(
			ErrInvalidConditionalOrderTriggerConfig,
			"maxRemovalsPerBlock %d exceeds the configurable upper bound %d",
			msg.MaxRemovalsPerBlock,
			MaxConfigurableRemovalsPerBlock,
		)
	}
	if msg.MaxUntriggeredConditionalOrdersGlobal > MaxConfigurableUntriggeredGlobal {
		return errorsmod.Wrapf(
			ErrInvalidConditionalOrderTriggerConfig,
			"maxUntriggeredConditionalOrdersGlobal %d exceeds the configurable upper bound %d",
			msg.MaxUntriggeredConditionalOrdersGlobal,
			MaxConfigurableUntriggeredGlobal,
		)
	}
	if msg.MaxUntriggeredConditionalOrdersPerSubaccount > MaxConfigurableUntriggeredPerSubaccount {
		return errorsmod.Wrapf(
			ErrInvalidConditionalOrderTriggerConfig,
			"maxUntriggeredConditionalOrdersPerSubaccount %d exceeds the configurable upper bound %d",
			msg.MaxUntriggeredConditionalOrdersPerSubaccount,
			MaxConfigurableUntriggeredPerSubaccount,
		)
	}

	// When both caps are explicitly set (non-zero), the per-subaccount cap must not exceed the
	// global cap — otherwise the per-subaccount cap is meaningless.
	if msg.MaxUntriggeredConditionalOrdersGlobal != 0 &&
		msg.MaxUntriggeredConditionalOrdersPerSubaccount != 0 &&
		msg.MaxUntriggeredConditionalOrdersPerSubaccount > msg.MaxUntriggeredConditionalOrdersGlobal {
		return errorsmod.Wrap(
			ErrInvalidConditionalOrderTriggerConfig,
			fmt.Sprintf(
				"per-subaccount cap %d exceeds global cap %d",
				msg.MaxUntriggeredConditionalOrdersPerSubaccount,
				msg.MaxUntriggeredConditionalOrdersGlobal,
			),
		)
	}

	return nil
}
