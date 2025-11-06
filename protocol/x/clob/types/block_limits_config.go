package types

// Validate checks that the BlockLimitsConfig is valid.
// Note: MaxStatefulOrderRemovalsPerBlock can be 0, which means no cap (process all expired orders).
func (config *BlockLimitsConfig) Validate() error {
	// No validation needed - 0 is a valid value meaning "no cap"
	return nil
}
