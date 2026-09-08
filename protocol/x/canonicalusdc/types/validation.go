package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"
	assettypes "github.com/dydxprotocol/v4-chain/protocol/x/assets/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func DefaultControls() Controls {
	return Controls{
		Mode:         Mode_MODE_DISABLED,
		LogicalDenom: assettypes.UusdcDenom,
		MemoVersion:  SupportedMemoVersion,
	}
}

func DefaultLedger() Ledger {
	return Ledger{
		NobleBacking:       "0",
		InjectiveBacking:   "0",
		LegacyDownstream:   "0",
		RestrictedFunding:  "0",
		PendingInjective:   "0",
		PendingNoble:       "0",
		PendingBackingSwap: "0",
	}
}

func DefaultGenesis() *GenesisState {
	return &GenesisState{Controls: DefaultControls(), Ledger: DefaultLedger()}
}

func ParseAmount(value string) (math.Int, error) {
	if len(value) == 0 || len(value) > 78 {
		return math.Int{}, ErrInvalidAmount
	}
	amount, ok := math.NewIntFromString(value)
	if !ok || amount.IsNegative() {
		return math.Int{}, ErrInvalidAmount
	}
	return amount, nil
}

func ParsePositiveAmount(value string) (math.Int, error) {
	amount, err := ParseAmount(value)
	if err != nil || !amount.IsPositive() {
		return math.Int{}, ErrInvalidAmount
	}
	return amount, nil
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
	if controls.MemoVersion != SupportedMemoVersion {
		return fmt.Errorf("%w: memo version must be %d", ErrInvalidControls, SupportedMemoVersion)
	}
	if controls.Mode == Mode_MODE_DISABLED {
		return nil
	}
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
	nobleDenom := transfertypes.ParseDenomTrace(
		transfertypes.GetDenomPrefix(transfertypes.PortID, controls.NobleChannel) + controls.NoblePacketDenom,
	).IBCDenom()
	if nobleDenom != controls.LogicalDenom {
		return fmt.Errorf("%w: Noble trace resolves to %s, not logical denom", ErrInvalidControls, nobleDenom)
	}
	injectiveDenom := transfertypes.ParseDenomTrace(
		transfertypes.GetDenomPrefix(transfertypes.PortID, controls.InjectiveChannel) + controls.InjectivePacketDenom,
	).IBCDenom()
	if injectiveDenom != controls.InjectiveDenom {
		return fmt.Errorf("%w: Injective trace resolves to %s, not physical denom", ErrInvalidControls, injectiveDenom)
	}
	if strings.ContainsRune(controls.InjectivePacketDenom, '\x00') ||
		strings.ContainsRune(controls.NoblePacketDenom, '\x00') {
		return fmt.Errorf("%w: packet denoms contain NUL", ErrInvalidControls)
	}
	maxTransfer, err := ParsePositiveAmount(controls.MaxTransferAmount)
	if err != nil {
		return fmt.Errorf("%w: max transfer amount", ErrInvalidControls)
	}
	ceiling, err := ParsePositiveAmount(controls.MigrationCeiling)
	if err != nil || ceiling.LT(maxTransfer) {
		return fmt.Errorf("%w: migration ceiling", ErrInvalidControls)
	}
	if controls.MaxPendingSettlements == 0 || controls.MaxPendingSettlements > HardMaxPending {
		return fmt.Errorf("%w: max pending settlements", ErrInvalidControls)
	}
	if controls.NobleWithdrawalCutoffTimestamp < 0 {
		return fmt.Errorf("%w: negative Noble cutoff", ErrInvalidControls)
	}
	return nil
}

func (ledger Ledger) Validate() error {
	values := []string{
		ledger.NobleBacking, ledger.InjectiveBacking, ledger.LegacyDownstream,
		ledger.RestrictedFunding, ledger.PendingInjective, ledger.PendingNoble,
		ledger.PendingBackingSwap,
	}
	for _, value := range values {
		if _, err := ParseAmount(value); err != nil {
			return ErrInvalidLedger
		}
	}
	return nil
}

func (participant Participant) Validate() error {
	if _, err := sdk.AccAddressFromBech32(participant.Controller); err != nil {
		return fmt.Errorf("%w: controller: %v", ErrInvalidParticipant, err)
	}
	if participant.NobleRecipient == "" || len(participant.NobleRecipient) > 256 {
		return fmt.Errorf("%w: Noble recipient", ErrInvalidParticipant)
	}
	maxRelease, err := ParsePositiveAmount(participant.MaxRelease)
	if err != nil {
		return fmt.Errorf("%w: max release", ErrInvalidParticipant)
	}
	released, err := ParseAmount(participant.Released)
	if err != nil || released.GT(maxRelease) {
		return fmt.Errorf("%w: released amount", ErrInvalidParticipant)
	}
	if _, err := ParseAmount(participant.Funded); err != nil {
		return fmt.Errorf("%w: funded amount", ErrInvalidParticipant)
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
	if len(genesis.Participants) > MaxParticipants {
		return fmt.Errorf("%w: too many participants", ErrInvalidParticipant)
	}
	seen := make(map[string]struct{}, len(genesis.Participants))
	for _, participant := range genesis.Participants {
		if err := participant.Validate(); err != nil {
			return err
		}
		if _, ok := seen[participant.Controller]; ok {
			return fmt.Errorf("%w: duplicate controller", ErrInvalidParticipant)
		}
		seen[participant.Controller] = struct{}{}
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
	if pending.Route < Route_ROUTE_INJECTIVE || pending.Route > Route_ROUTE_BACKING_SWAP {
		return fmt.Errorf("invalid pending route")
	}
	if err := host.ChannelIdentifierValidator(pending.SourceChannel); err != nil {
		return err
	}
	if pending.Sequence == 0 {
		return fmt.Errorf("pending sequence must be positive")
	}
	_, err := ParsePositiveAmount(pending.Amount)
	return err
}
