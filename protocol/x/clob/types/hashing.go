package types

import (
	"crypto/sha256"
)

// ProtobufHashable represents a struct that is generated
// from protobuf. We use the two provided methods to hash
// the struct and compare equality.
type ProtobufHashable interface {
	Marshal() ([]byte, error)
	Size() int
}

// ProtoHash is the result of hashing a ProtobufHashable.
type ProtoHash [32]byte

var _ ProtobufHashable = &MatchPerpetualLiquidation{}
var _ ProtobufHashable = &MatchOrders{}
var _ ProtobufHashable = &MsgCancelOrder{}
var _ ProtobufHashable = &MsgPlaceOrder{}
var _ ProtobufHashable = &MatchPerpetualDeleveraging{}
var _ ProtobufHashable = &Operation{}
var _ ProtobufHashable = &Order{}

// GetHash returns the 32-byte sha256 hash of a hashable object.
func GetHash(m ProtobufHashable) ProtoHash {
	bytes, err := m.Marshal()
	if err != nil {
		panic(err)
	}
	return sha256.Sum256(bytes)
}

// IsEqual returns if 2 hashable objects are equal by sha.
// Be wary of golang type casting in nil cases.
func IsEqual(x ProtobufHashable, y ProtobufHashable) bool {
	var tValuedNil ProtobufHashable // T-typed nil value
	if x == tValuedNil || y == tValuedNil {
		panic("IsEqual cannot compare a nil object")
	}
	if x.Size() != y.Size() {
		return false
	}
	return GetHash(x) == GetHash(y)
}
