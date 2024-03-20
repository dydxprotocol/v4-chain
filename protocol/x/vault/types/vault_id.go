package types

import (
	fmt "fmt"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// ToStateKey returns the state key for the vault ID.
func (id *VaultId) ToStateKey() []byte {
	b, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	return b
}

// ToModuleAccountAddress returns the module account address for the vault ID
// (generated from string "vault-<type>-<number>").
func (id *VaultId) ToModuleAccountAddress() string {
	return authtypes.NewModuleAddress(
		fmt.Sprintf("vault-%s-%d", id.Type, id.Number),
	).String()
}

// ToSubaccountId returns the subaccount ID for the vault ID (owner is
// corresponding module account address and number is 0).
func (id *VaultId) ToSubaccountId() *satypes.SubaccountId {
	return &satypes.SubaccountId{
		Owner:  id.ToModuleAccountAddress(),
		Number: 0,
	}
}

// IncrCounterWithLabels increments counter with labels with added vault ID labels.
func (id *VaultId) IncrCounterWithLabels(metricName string, labels ...metrics.Label) {
	// Append vault type and number to labels.
	labels = append(
		labels,
		metrics.GetLabelForIntValue(
			metrics.VaultType,
			int(id.Type),
		),
		metrics.GetLabelForIntValue(
			metrics.VaultId,
			int(id.Number),
		),
	)

	metrics.IncrCounterWithLabels(
		metricName,
		1,
		labels...,
	)
}
