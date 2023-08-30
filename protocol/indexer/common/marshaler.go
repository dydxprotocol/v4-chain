package common

import (
	"github.com/cosmos/gogoproto/proto"
)

// MarshalerImpl is the struct that implements the `Marshaler` interface.
type MarshalerImpl struct{}

// Ensure the `MarshalerImpl` struct is implemented at compile time.
var _ Marshaler = (*MarshalerImpl)(nil)

// Marshaler is an interface that encapsulates the gogo protobuf function `Marshal`.
type Marshaler interface {
	Marshal(pb proto.Message) ([]byte, error)
}

// Marshal wraps `proto.Marshal` function which marshals a proto message into a byte array.
func (m *MarshalerImpl) Marshal(pb proto.Message) ([]byte, error) {
	return proto.Marshal(pb)
}
