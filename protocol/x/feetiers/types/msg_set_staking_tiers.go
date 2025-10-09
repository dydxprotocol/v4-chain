package types

import (
	"fmt"
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidateBasic performs basic validation on MsgSetStakingTiers
func (msg *MsgSetStakingTiers) ValidateBasic() error {
	// Validate authority address is valid
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrap(
			err,
			"invalid authority address",
		)
	}

	// Validate staking tiers where
	// - No duplicate fee tier names
	// - Fee tier name is not empty
	// - Staking levels are valid:
	//   - Min staked tokens is a valid non-negative number
	//   - Levels are in strictly increasing order of min staked tokens
	//   - Discount is not more than 100%
	seenTiers := make(map[string]bool)
	for _, tier := range msg.StakingTiers {
		// Validate fee tier name is not empty
		if tier.FeeTierName == "" {
			return fmt.Errorf("fee tier name cannot be empty")
		}

		if seenTiers[tier.FeeTierName] {
			return fmt.Errorf("duplicate staking tier for fee tier: %s", tier.FeeTierName)
		}
		seenTiers[tier.FeeTierName] = true

		// Validate staking levels
		var prevMinStaked *big.Int
		for i, level := range tier.Levels {
			// Validate min staked tokens is a valid number
			minStaked := new(big.Int)
			if _, ok := minStaked.SetString(level.MinStakedBaseTokens, 10); !ok {
				return fmt.Errorf("invalid min staked tokens for tier %s level %d: %s",
					tier.FeeTierName, i, level.MinStakedBaseTokens)
			}

			// Check that min staked is non-negative
			if minStaked.Sign() < 0 {
				return fmt.Errorf("min staked tokens cannot be negative for tier %s level %d",
					tier.FeeTierName, i)
			}

			// Check that levels are in increasing order
			if prevMinStaked != nil && minStaked.Cmp(prevMinStaked) <= 0 {
				return fmt.Errorf("staking levels must be in increasing order for tier %s",
					tier.FeeTierName)
			}
			prevMinStaked = minStaked

			// Validate discount is not more than 100%
			if level.FeeDiscountPpm > 1_000_000 {
				return fmt.Errorf("fee discount cannot exceed 100%% for tier %s level %d",
					tier.FeeTierName, i)
			}
		}
	}

	return nil
}
