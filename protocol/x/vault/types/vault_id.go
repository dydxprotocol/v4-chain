package types

import (
	fmt "fmt"
	"strconv"
	"strings"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/dydxprotocol/v4-chain/protocol/lib/metrics"
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
	satypes "github.com/dydxprotocol/v4-chain/protocol/x/subaccounts/types"
)

// ToString returns the string representation of a vault ID.
func (id *VaultId) ToString() string {
	return fmt.Sprintf("%s-%d", id.Type, id.Number)
}

// ToStateKey returns the state key for the vault ID.
func (id *VaultId) ToStateKey() []byte {
	return []byte(id.ToString())
}

// ToStateKeyPrefix returns the state key prefix for the vault ID.
func (id *VaultId) ToStateKeyPrefix() []byte {
	return []byte(fmt.Sprintf("%s:", id.ToString()))
}

// GetVaultIdFromStateKey returns a vault ID from a given state key.
func GetVaultIdFromStateKey(stateKey []byte) (*VaultId, error) {
	stateKeyStr := string(stateKey)

	// Split state key string into type and number.
	split := strings.Split(stateKeyStr, "-")
	if len(split) != 2 {
		return nil, fmt.Errorf("stateKey in string must follow format <type>-<number> but got %s", stateKeyStr)
	}

	// Parse vault type.
	vaultTypeInt, exists := VaultType_value[split[0]]
	if !exists {
		return nil, fmt.Errorf("unknown vault type: %s", split[0])
	}

	// Parse vault number.
	number, err := strconv.ParseUint(split[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse number: %s", err.Error())
	}

	return &VaultId{
		Type:   VaultType(vaultTypeInt),
		Number: uint32(number),
	}, nil
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

// GetClobOrderId returns a vault CLOB order ID given a client ID.
func (id *VaultId) GetClobOrderId(clientId uint32) *clobtypes.OrderId {
	return &clobtypes.OrderId{
		SubaccountId: *id.ToSubaccountId(),
		ClientId:     clientId,
		OrderFlags:   clobtypes.OrderIdFlags_LongTerm,
		ClobPairId:   uint32(id.Number),
	}
}

// IncrCounterWithLabels increments counter with labels with added vault ID labels.
func (id *VaultId) IncrCounterWithLabels(metricName string, labels ...metrics.Label) {
	// Append vault labels.
	labels = id.addLabels(labels...)

	metrics.IncrCounterWithLabels(
		metricName,
		1,
		labels...,
	)
}

// IncrCounterWithLabels sets gauge with labels with added vault ID labels.
func (id *VaultId) SetGaugeWithLabels(
	metricName string,
	value float32,
	labels ...metrics.Label,
) {
	// Append vault labels.
	labels = id.addLabels(labels...)

	metrics.SetGaugeWithLabels(
		metricName,
		value,
		labels...,
	)
}

func (id *VaultId) addLabels(labels ...metrics.Label) []metrics.Label {
	return append(
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
}
