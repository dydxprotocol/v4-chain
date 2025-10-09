package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/testutil/constants"
	"github.com/dydxprotocol/v4-chain/protocol/x/feetiers/types"
	"github.com/stretchr/testify/require"
)

func TestMsgSetStakingTiers_ValidateBasic(t *testing.T) {
	tests := map[string]struct {
		msg           *types.MsgSetStakingTiers
		expectedError string
	}{
		"success - empty": {
			msg: &types.MsgSetStakingTiers{
				Authority:    constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{},
			},
			expectedError: "",
		},
		"success - single tier": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100000000000000000000",
								FeeDiscountPpm:      50000,
							},
						},
					},
				},
			},
			expectedError: "",
		},
		"success - multiple tiers and levels": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100",
								FeeDiscountPpm:      10000,
							},
							{
								MinStakedBaseTokens: "1000",
								FeeDiscountPpm:      20000,
							},
						},
					},
					{
						FeeTierName: "2",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "500",
								FeeDiscountPpm:      15000,
							},
						},
					},
				},
			},
			expectedError: "",
		},
		"success - zero discount": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100",
								FeeDiscountPpm:      0,
							},
						},
					},
				},
			},
			expectedError: "",
		},
		"success - 100% discount": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100",
								FeeDiscountPpm:      1_000_000, // 100%
							},
						},
					},
				},
			},
			expectedError: "",
		},
		"success - valid message with empty tiers": {
			msg: &types.MsgSetStakingTiers{
				Authority:    constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{},
			},
			expectedError: "",
		},
		"error - invalid authority address": {
			msg: &types.MsgSetStakingTiers{
				Authority:    "invalid-address",
				StakingTiers: []*types.StakingTier{},
			},
			expectedError: "invalid authority address",
		},
		"error - empty fee tier name": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100",
								FeeDiscountPpm:      10000,
							},
						},
					},
				},
			},
			expectedError: "fee tier name cannot be empty",
		},
		"error - duplicate fee tier names": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels:      []*types.StakingLevel{},
					},
					{
						FeeTierName: "1",
						Levels:      []*types.StakingLevel{},
					},
				},
			},
			expectedError: "duplicate staking tier for fee tier: 1",
		},
		"error - invalid min staked tokens": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "not-a-number",
								FeeDiscountPpm:      10000,
							},
						},
					},
				},
			},
			expectedError: "invalid min staked tokens for tier 1 level 0",
		},
		"error - negative min staked tokens": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "-100",
								FeeDiscountPpm:      10000,
							},
						},
					},
				},
			},
			expectedError: "min staked tokens cannot be negative for tier 1 level 0",
		},
		"error - levels in decreasing order of staked amount": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "1000",
								FeeDiscountPpm:      10000,
							},
							{
								MinStakedBaseTokens: "999", // Less than previous
								FeeDiscountPpm:      20000,
							},
						},
					},
				},
			},
			expectedError: "staking levels must be in increasing order for tier 1",
		},
		"error - levels with equal staked amounts": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "1000",
								FeeDiscountPpm:      10000,
							},
							{
								MinStakedBaseTokens: "1000", // Equal to previous
								FeeDiscountPpm:      20000,
							},
						},
					},
				},
			},
			expectedError: "staking levels must be in increasing order for tier 1",
		},
		"error - discount exceeds 100%": {
			msg: &types.MsgSetStakingTiers{
				Authority: constants.AliceAccAddress.String(),
				StakingTiers: []*types.StakingTier{
					{
						FeeTierName: "1",
						Levels: []*types.StakingLevel{
							{
								MinStakedBaseTokens: "100",
								FeeDiscountPpm:      1_000_001, // > 100%
							},
						},
					},
				},
			},
			expectedError: "fee discount cannot exceed 100% for tier 1 level 0",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.msg.ValidateBasic()
			if tc.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
