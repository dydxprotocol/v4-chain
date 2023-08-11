package types_test

import (
	"testing"

	"github.com/dydxprotocol/v4/app/config"
	"github.com/dydxprotocol/v4/x/perpetuals/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	testAddress = "dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf"
)

func TestMsgAddPremiumVotes(t *testing.T) {
	sample := types.NewFundingPremium(uint32(0), int32(10))
	samples := []types.FundingPremium{*sample}
	msg := types.NewMsgAddPremiumVotes(testAddress, samples)

	require.Equal(t, uint32(0), sample.PerpetualId)
	require.Equal(t, int32(10), sample.PremiumPpm)
	require.Equal(t, samples, msg.Votes)
}

func TestMsgAddPremiumVotes_GetSigners_Success(t *testing.T) {
	// This package does not contain the `app/config` package in its import chain, and therefore needs to call
	// SetAddressPrefixes() explicitly in order to set the `dydx` address prefixes.
	config.SetAddressPrefixes()

	sample := types.NewFundingPremium(uint32(0), int32(10))
	samples := []types.FundingPremium{*sample}
	msg := types.NewMsgAddPremiumVotes(testAddress, samples)

	expectedAddress, err := sdk.AccAddressFromBech32(testAddress)
	require.NoError(t, err)

	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, expectedAddress, signers[0])
}

func TestMsgUpdateMarketPrices_GetSigners_Panics(t *testing.T) {
	sample := types.NewFundingPremium(uint32(0), int32(10))
	samples := []types.FundingPremium{*sample}
	msg := types.NewMsgAddPremiumVotes("invalid", samples)

	require.PanicsWithError(
		t,
		"decoding bech32 failed: invalid bech32 string length 7",
		func() { msg.GetSigners() })
}

func TestValidateBasic(t *testing.T) {
	errStr := "premium votes must be sorted by perpetual id in ascending order and cannot " +
		"contain duplicates: MsgAddPremiumVotes is invalid"

	tests := map[string]struct {
		samples []types.FundingPremium

		expectedErr bool
	}{
		"Error: duplicate perpetual ids": {
			samples: []types.FundingPremium{
				{PerpetualId: 1},
				{PerpetualId: 1},
			},
			expectedErr: true,
		},
		"Error: desending perpetual ids": {
			samples: []types.FundingPremium{
				{PerpetualId: 2},
				{PerpetualId: 3},
				{PerpetualId: 4},
				{PerpetualId: 1},
			},
			expectedErr: true,
		},
		"No error: empty samples": {
			samples: []types.FundingPremium{},
		},
		"No error: valid ordering": {
			samples: []types.FundingPremium{
				{PerpetualId: 1},
				{PerpetualId: 3},
				{PerpetualId: 99},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			msg := types.NewMsgAddPremiumVotes("", tc.samples)
			err := msg.ValidateBasic()
			if tc.expectedErr {
				require.ErrorContains(t, err, errStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
