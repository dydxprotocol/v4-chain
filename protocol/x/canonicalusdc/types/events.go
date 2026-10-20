package types

const (
	EventTypeCanonicalTransfer = "canonical_usdc_transfer"
	EventTypeSettlement        = "canonical_usdc_settlement"
	EventTypeInbound           = "canonical_usdc_inbound"
	EventTypeBackingBreach     = "canonical_usdc_backing_breach"

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
