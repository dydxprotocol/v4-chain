package types

import errorsmod "cosmossdk.io/errors"

const (
	ModuleName = "full_node_streaming"
)

var (
	ErrNotImplemented          = errorsmod.Register(ModuleName, 1, "Not implemented")
	ErrInvalidStreamingRequest = errorsmod.Register(
		ModuleName,
		2,
		"Invalid full node streaming request",
	)
	ErrInvalidSubaccountFilteringRequest = errorsmod.Register(
		ModuleName,
		3,
		"Invalid subaccount ID filtering request",
	)
)
