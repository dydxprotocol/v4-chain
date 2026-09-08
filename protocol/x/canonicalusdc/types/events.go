package types

const (
	EventTypeCanonicalTransfer = "canonical_usdc_transfer"
	EventTypeSettlement        = "canonical_usdc_settlement"
	EventTypeBackingFunding    = "canonical_usdc_backing_funding"

	AttributeKeyRoute          = "route"
	AttributeKeyPhysicalDenom  = "physical_denom"
	AttributeKeyLogicalDenom   = "logical_denom"
	AttributeKeyPending        = "pending"
	AttributeKeyOutcome        = "outcome"
	AttributeKeySequence       = "sequence"
	AttributeKeySourceChannel  = "source_channel"
	AttributeKeyController     = "controller"
	AttributeKeyNobleRecipient = "noble_recipient"
	AttributeValueSuccess      = "success"
	AttributeValueError        = "error"
	AttributeValueTimeout      = "timeout"
)
