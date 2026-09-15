package types

const (
	EventTypeCanonicalTransfer = "canonical_usdc_transfer"
	EventTypeSettlement        = "canonical_usdc_settlement"
	EventTypeInbound           = "canonical_usdc_inbound"

	AttributeKeyRoute         = "route"
	AttributeKeyPhysicalDenom = "physical_denom"
	AttributeKeyLogicalDenom  = "logical_denom"
	AttributeKeyPending       = "pending"
	AttributeKeyOutcome       = "outcome"
	AttributeKeySequence      = "sequence"
	AttributeKeySourceChannel = "source_channel"
	AttributeValueSuccess     = "success"
	AttributeValueError       = "error"
	AttributeValueTimeout     = "timeout"
)
