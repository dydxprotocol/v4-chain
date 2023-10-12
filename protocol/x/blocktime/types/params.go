package types

func (m *DowntimeParams) Validate() error {
	if m.Durations != nil {
		for i := 0; i < len(m.Durations); i++ {
			if m.Durations[i] <= 0 {
				return ErrNonpositiveDuration
			}
		}

		for i := 0; i < len(m.Durations)-1; i++ {
			if m.Durations[i] >= m.Durations[i+1] {
				return ErrUnorderedDurations
			}
		}
	}
	return nil
}
