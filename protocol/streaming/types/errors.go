package types

import errorsmod "cosmossdk.io/errors"

const (
	ModuleName = "full_node_streaming"
)

var (
	ErrNotImplemented = errorsmod.Register(ModuleName, 1, "Not implemented")
)
