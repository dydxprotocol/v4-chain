package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	"github.com/dydxprotocol/v4-chain/protocol/dtypes"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultControls() Controls {
	return Controls{
		Mode:              Mode_MODE_DISABLED,
		LogicalDenom:      assettypes.UusdcDenom,
		MaxTransferAmount: dtypes.ZeroInt(),
		MigrationCeiling:  dtypes.ZeroInt(),
	}
}

func DefaultLedger() Ledger {
	return Ledger{
		NobleBacking:     dtypes.ZeroInt(),
		InjectiveBacking: dtypes.ZeroInt(),
		PendingInjective: dtypes.ZeroInt(),
		PendingNoble:     dtypes.ZeroInt(),
	}
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{Controls: DefaultControls(), Ledger: DefaultLedger()}
}

// AmountOf converts a stored amount to math.Int. An unset value is zero.
// Negative values and values wider than math.Int are rejected.
func AmountOf(value dtypes.SerializableInt) (math.Int, error) {
	if value.IsNil() {
		return math.ZeroInt(), nil
	}
	bigValue := value.BigInt()
	if bigValue.Sign() < 0 || bigValue.BitLen() > math.MaxBitLen {
		return math.Int{}, ErrInvalidAmount
	}
	return math.NewIntFromBigInt(bigValue), nil
}

// PositiveAmountOf is AmountOf restricted to strictly positive values.
func PositiveAmountOf(value dtypes.SerializableInt) (math.Int, error) {
	amount, err := AmountOf(value)
	if err != nil || !amount.IsPositive() {
		return math.Int{}, ErrInvalidAmount
	}
	return amount, nil
}

// NewAmount converts a math.Int to its stored form.
func NewAmount(amount math.Int) dtypes.SerializableInt {
	return dtypes.NewIntFromBigInt(amount.BigInt())
}

// NobleTrace is the full denom path that ICS20 writes into packets carrying
// Noble-origin USDC out of this chain.
func (controls Controls) NobleTrace() string {
	return transfertypes.GetDenomPrefix(transfertypes.PortID, controls.NobleChannel) + controls.NoblePacketDenom
}

// InjectiveTrace is the full denom path that ICS20 writes into packets carrying
// Injective-origin USDC out of this chain.
func (controls Controls) InjectiveTrace() string {
	return transfertypes.GetDenomPrefix(transfertypes.PortID, controls.InjectiveChannel) +
		controls.InjectivePacketDenom
}

// SameRoute reports whether both controls bind the same channels, clients,
// connections, and denoms.
func (controls Controls) SameRoute(other Controls) bool {
	return controls.LogicalDenom == other.LogicalDenom &&
		controls.NobleChannel == other.NobleChannel &&
		controls.NobleClient == other.NobleClient &&
		controls.NobleConnection == other.NobleConnection &&
		controls.NoblePacketDenom == other.NoblePacketDenom &&
		controls.InjectiveChannel == other.InjectiveChannel &&
		controls.InjectiveClient == other.InjectiveClient &&
		controls.InjectiveConnection == other.InjectiveConnection &&
		controls.InjectivePacketDenom == other.InjectivePacketDenom &&
		controls.InjectiveDenom == other.InjectiveDenom
}

func (controls Controls) Validate() error {
	if controls.Mode < Mode_MODE_DISABLED || controls.Mode > Mode_MODE_PAUSED {
		return fmt.Errorf("%w: unsupported mode %s", ErrInvalidControls, controls.Mode)
	}
	if err := sdk.ValidateDenom(controls.LogicalDenom); err != nil {
		return fmt.Errorf("%w: logical denom: %v", ErrInvalidControls, err)
	}
	if controls.LogicalDenom != assettypes.UusdcDenom {
		return fmt.Errorf("%w: logical denom must remain %s", ErrInvalidControls, assettypes.UusdcDenom)
	}
	if controls.Mode == Mode_MODE_DISABLED {
		return nil
	}
	if err := controls.validateRoutes(); err != nil {
		return err
	}
	return controls.validateLimits()
}

func (controls Controls) validateRoutes() error {
	if controls.NobleChannel == controls.InjectiveChannel {
		return fmt.Errorf("%w: route channels must differ", ErrInvalidControls)
	}
	if err := host.ChannelIdentifierValidator(controls.NobleChannel); err != nil {
		return fmt.Errorf("%w: Noble channel: %v", ErrInvalidControls, err)
	}
	if err := host.ChannelIdentifierValidator(controls.InjectiveChannel); err != nil {
		return fmt.Errorf("%w: Injective channel: %v", ErrInvalidControls, err)
	}
	if err := host.ClientIdentifierValidator(controls.NobleClient); err != nil {
		return fmt.Errorf("%w: Noble client: %v", ErrInvalidControls, err)
	}
	if err := host.ClientIdentifierValidator(controls.InjectiveClient); err != nil {
		return fmt.Errorf("%w: Injective client: %v", ErrInvalidControls, err)
	}
	if err := host.ConnectionIdentifierValidator(controls.NobleConnection); err != nil {
		return fmt.Errorf("%w: Noble connection: %v", ErrInvalidControls, err)
	}
	if err := host.ConnectionIdentifierValidator(controls.InjectiveConnection); err != nil {
		return fmt.Errorf("%w: Injective connection: %v", ErrInvalidControls, err)
	}
	if err := sdk.ValidateDenom(controls.InjectiveDenom); err != nil {
		return fmt.Errorf("%w: Injective denom: %v", ErrInvalidControls, err)
	}
	if controls.InjectivePacketDenom == "" || controls.NoblePacketDenom == "" {
		return fmt.Errorf("%w: packet denoms must be set", ErrInvalidControls)
	}
	if strings.ContainsRune(controls.InjectivePacketDenom, '\x00') ||
		strings.ContainsRune(controls.NoblePacketDenom, '\x00') {
		return fmt.Errorf("%w: packet denoms contain NUL", ErrInvalidControls)
	}
	if nobleDenom := transfertypes.ParseDenomTrace(controls.NobleTrace()).IBCDenom(); nobleDenom != controls.LogicalDenom {
		return fmt.Errorf("%w: Noble trace resolves to %s, not logical denom", ErrInvalidControls, nobleDenom)
	}
	injectiveDenom := transfertypes.ParseDenomTrace(controls.InjectiveTrace()).IBCDenom()
	if injectiveDenom != controls.InjectiveDenom {
		return fmt.Errorf("%w: Injective trace resolves to %s, not physical denom", ErrInvalidControls, injectiveDenom)
	}
	return nil
}

func (controls Controls) validateLimits() error {
	maxTransfer, err := PositiveAmountOf(controls.MaxTransferAmount)
	if err != nil {
		return fmt.Errorf("%w: max transfer amount", ErrInvalidControls)
	}
	ceiling, err := PositiveAmountOf(controls.MigrationCeiling)
	if err != nil || ceiling.LT(maxTransfer) {
		return fmt.Errorf("%w: migration ceiling", ErrInvalidControls)
	}
	if controls.MaxPendingSettlements == 0 || controls.MaxPendingSettlements > HardMaxPending {
		return fmt.Errorf("%w: max pending settlements", ErrInvalidControls)
	}
	return nil
}

func (ledger Ledger) Validate() error {
	values := []dtypes.SerializableInt{
		ledger.NobleBacking, ledger.InjectiveBacking, ledger.PendingInjective, ledger.PendingNoble,
	}
	for _, value := range values {
		if _, err := AmountOf(value); err != nil {
			return ErrInvalidLedger
		}
	}
	return nil
}

func (genesis GenesisState) Validate() error {
	if err := genesis.Controls.Validate(); err != nil {
		return err
	}
	if err := genesis.Ledger.Validate(); err != nil {
		return err
	}
	if len(genesis.PendingSettlements) > int(genesis.Controls.MaxPendingSettlements) {
		return ErrPendingLimit
	}
	if genesis.Controls.Mode == Mode_MODE_DISABLED && len(genesis.PendingSettlements) != 0 {
		return fmt.Errorf("%w: disabled genesis cannot contain pending settlements", ErrInvalidControls)
	}
	for _, pending := range genesis.PendingSettlements {
		if err := pending.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (pending PendingSettlement) Validate() error {
	if pending.Route < Route_ROUTE_INJECTIVE || pending.Route > Route_ROUTE_NOBLE {
		return fmt.Errorf("invalid pending route")
	}
	if err := host.ChannelIdentifierValidator(pending.SourceChannel); err != nil {
		return err
	}
	if pending.Sequence == 0 {
		return fmt.Errorf("pending sequence must be positive")
	}
	_, err := PositiveAmountOf(pending.Amount)
	return err
}
