package types

type ClobMetadata struct {
	ClobPair    ClobPair
	MidPrice    Subticks
	OraclePrice Subticks
	BestBid     Order
	BestAsk     Order
}

type MevTelemetryConfig struct {
	Enabled    bool
	Hosts      []string
	Identifier string
}
