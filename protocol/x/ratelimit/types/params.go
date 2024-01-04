package types

// Validate validates the set of params
func (p *LimitParams) Validate() error {
	// TODO(CORE-824): implement keepers. Check that `BaselineMinimum` and `BaselineTvlPpm` are both positive.
	return nil
}
