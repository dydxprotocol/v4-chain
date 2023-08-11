package types_test

import (
	"errors"
	"testing"

	"github.com/dydxprotocol/v4/lib"
	"github.com/dydxprotocol/v4/testutil/constants"
	"github.com/dydxprotocol/v4/x/clob/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	tests := map[string]struct {
		genState      *types.GenesisState
		expectedError error
	}{
		"default is valid": {
			genState:      types.DefaultGenesis(),
			expectedError: nil,
		},
		"valid genesis state": {
			genState: &types.GenesisState{

				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(1),
					},
				},
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion + 1,
						SpreadToMaintenanceMarginRatioPpm: 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: nil,
		},
		"duplicated clobPair": {
			genState: &types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(0),
					},
				},
			},
			expectedError: errors.New("duplicated id for clobPair"),
		},
		"gap in clobPair": {
			genState: &types.GenesisState{
				ClobPairs: []types.ClobPair{
					{
						Id: uint32(0),
					},
					{
						Id: uint32(2),
					},
				},
			},
			expectedError: errors.New("found gap in clobPair id"),
		},
		"spread to maintenance margin ratio of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: 0,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid SpreadToMaintenanceMarginRatioPpm: Proposed LiquidationsConfig is invalid"),
		},
		"spread to maintenance margin ratio of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion + 1,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid SpreadToMaintenanceMarginRatioPpm: Proposed LiquidationsConfig is invalid"),
		},
		"bankruptcy adjustment ppm of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           0,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"bankruptcy adjustment ppm of less than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: 100_00,
					FillablePriceConfig: types.FillablePriceConfig{
						BankruptcyAdjustmentPpm:           lib.OneMillion - 1,
						SpreadToMaintenanceMarginRatioPpm: lib.OneMillion,
					},
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"999999 is not a valid BankruptcyAdjustmentPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of zero is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  0,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxLiquidationFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max liquidation fee ppm of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm:  lib.OneMillion + 1,
					FillablePriceConfig:   constants.FillablePriceConfig_Default,
					PositionBlockLimits:   constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid MaxLiquidationFeePpm: Proposed LiquidationsConfig is invalid"),
		},
		"max position portion liquidated ppm of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   1_000,
						MaxPositionPortionLiquidatedPpm: 0,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxPositionPortionLiquidatedPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max position portion liquidated ppm of greater than one million is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits: types.PositionBlockLimits{
						MinPositionNotionalLiquidated:   1_000,
						MaxPositionPortionLiquidatedPpm: lib.OneMillion + 1,
					},
					SubaccountBlockLimits: constants.SubaccountBlockLimits_Default,
				},
			},
			expectedError: errors.New(
				"1000001 is not a valid MaxPositionPortionLiquidatedPpm: Proposed LiquidationsConfig is invalid"),
		},
		"max notional liquidated of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    0,
						MaxQuantumsInsuranceLost: 100_000_000_000_000,
					},
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxNotionalLiquidated: Proposed LiquidationsConfig is invalid"),
		},
		"max quantums insurance lost of 0 is invalid": {
			genState: &types.GenesisState{
				LiquidationsConfig: types.LiquidationsConfig{
					MaxLiquidationFeePpm: lib.OneMillion,
					FillablePriceConfig:  constants.FillablePriceConfig_Default,
					PositionBlockLimits:  constants.PositionBlockLimits_Default,
					SubaccountBlockLimits: types.SubaccountBlockLimits{
						MaxNotionalLiquidated:    100_000_000_000_000,
						MaxQuantumsInsuranceLost: 0,
					},
				},
			},
			expectedError: errors.New(
				"0 is not a valid MaxQuantumsInsuranceLost: Proposed LiquidationsConfig is invalid"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.expectedError == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedError.Error())
			}
		})
	}
}
