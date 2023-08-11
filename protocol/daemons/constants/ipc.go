package constants

const (
	// UnixProtocol is the network for gRPC protocol used by the price daemon and server.
	UnixProtocol = "unix"

	// UmaskUserReadWriteOnly is the Inverse Unix Permission code for only the user to read or write.
	UmaskUserReadWriteOnly = 0177

	// DefaultGrpcEndpoint is the default grpc endpoint for Cosmos query services.
	DefaultGrpcEndpoint = "localhost:9090"
)
