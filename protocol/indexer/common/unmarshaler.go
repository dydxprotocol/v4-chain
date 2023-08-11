package common

import (
	"github.com/cosmos/gogoproto/proto"
)

// UnmarshalerImpl is the struct that implements the `Unmarshaler` interface.
type UnmarshalerImpl struct{}

// Ensure the `UnmarshalerImpl` struct is implemented at compile time.
var _ Unmarshaler = (*UnmarshalerImpl)(nil)

// Unmarshaler is an interface that encapsulates the gogo protobuf function `Unmarshal`.
type Unmarshaler interface {
	Unmarshal(bytes []byte, pb proto.Message) error
}

// Unmarshal wraps `proto.Unmarshal` function which unmarshals a byte array into a proto message.
func (m *UnmarshalerImpl) Unmarshal(bytes []byte, pb proto.Message) error {
	return proto.Unmarshal(bytes, pb)
}
