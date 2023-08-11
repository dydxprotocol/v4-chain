package sim_helpers

const (
	// Used to determine the % frequency that used genesis values are reasonable (e.g. similar to
	// to potential real-world conditions) in simulation tests. Value should be between 0 and 100.
	// For ex. if ReasonableGenesisWeight = 90, 90% of simulation tests will use reasonable values.
	ReasonableGenesisWeight = 90
)
