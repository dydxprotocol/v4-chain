package types

import time "time"

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

func (s SynchronyParams) Validate() error {
	if s.NextBlockDelay < 0 {
		return ErrNegativeNextBlockDelay
	}
	return nil
}

func DefaultSynchronyParams() SynchronyParams {
	return SynchronyParams{
		// CometBFT defaults back to `timeout_commit` if application sends over
		// `NextBlockDelay` of 0.
		NextBlockDelay: 0 * time.Second,
	}
}
