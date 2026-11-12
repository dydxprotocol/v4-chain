package types

// Indexers should read canonical USDC activity from the events in this file,
// not from the ICS20 events emitted alongside them.
//
// On the Injective route the physical packet is sent and received by the
// module account. The ICS20 `ibc_transfer` event for that packet names the
// module as sender, and the ICS20 `fungible_token_packet` event for a deposit
// names the module as receiver. A withdrawal on that route also carries a
// second `ibc_transfer` event that names the user and the logical denom, so
// counting `ibc_transfer` events counts an Injective withdrawal twice.
//
// A withdrawal is one canonical_usdc_transfer followed by exactly one
// canonical_usdc_settlement with the same route, source_channel and sequence.
// A settlement with outcome legacy_refund is the exception: its packet was sent
// while DISABLED, so it has no canonical_usdc_transfer. A deposit is one
// canonical_usdc_inbound.
const (
	// EventTypeCanonicalTransfer is emitted when a withdrawal is accepted and
	// its packet is sent. Attributes: route, logical_denom, physical_denom,
	// source_channel, sequence, amount, sender, receiver, pending.
	EventTypeCanonicalTransfer = "canonical_usdc_transfer"
	// EventTypeSettlement is emitted when a withdrawal's acknowledgement or
	// timeout is processed. Attributes: route, logical_denom, physical_denom,
	// source_channel, sequence, amount, sender, receiver, outcome.
	EventTypeSettlement = "canonical_usdc_settlement"
	// EventTypeInbound is emitted when an Injective deposit is credited.
	// Attributes: physical_denom, logical_denom, amount, receiver (the dYdX
	// account credited), sender (the Injective sender), source_channel,
	// sequence.
	EventTypeInbound = "canonical_usdc_inbound"
	// EventTypeBackingBreach is emitted when the EndBlocker pauses routing.
	// Attributes: reason.
	EventTypeBackingBreach = "canonical_usdc_backing_breach"

	AttributeKeyRoute         = "route"
	AttributeKeyPhysicalDenom = "physical_denom"
	AttributeKeyLogicalDenom  = "logical_denom"
	AttributeKeyPending       = "pending"
	AttributeKeyOutcome       = "outcome"
	AttributeKeySequence      = "sequence"
	AttributeKeySourceChannel = "source_channel"
	AttributeKeyReason        = "reason"
	AttributeValueSuccess     = "success"
	AttributeValueError       = "error"
	AttributeValueTimeout     = "timeout"
	// AttributeValueLegacyRefund marks a refund of a Noble packet sent while DISABLED.
	AttributeValueLegacyRefund = "legacy_refund"
)
