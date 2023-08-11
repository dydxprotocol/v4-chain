package types

func (m *Params) Validate() error {
	if m.WindowDuration <= 0 {
		return ErrNonpositiveDuration
	}
	return nil
}
