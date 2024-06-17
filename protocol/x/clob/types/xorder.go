package types

import (
	"encoding/binary"
	"math/bits"

	"github.com/cespare/xxhash"
	"github.com/dydxprotocol/v4-chain/protocol/lib"
)

// Place Order Flags

type Order_ReplaceFlags uint32

const (
	Order_REPLACE_FLAGS_NEW_ONLY Order_ReplaceFlags = 0
	Order_REPLACE_FLAGS_UPSERT   Order_ReplaceFlags = 1
	Order_REPLACE_FLAGS_INC_SIZE Order_ReplaceFlags = 2
	Order_REPLACE_FLAGS_DEC_SIZE Order_ReplaceFlags = 3
)

func GetReplaceFlags(placeFlags uint32) Order_ReplaceFlags {
	return Order_ReplaceFlags(placeFlags & 0b11)
}

// Self-Trade Prevention

type Order_STP uint32

const (
	Order_STP_EXPIRE_MAKER Order_STP = 0
	Order_STP_EXPIRE_TAKER Order_STP = 1
	Order_STP_EXPIRE_BOTH  Order_STP = 2
)

func FormOrder(
	uid XUID,
	base XBase,
) XOrder {
	return XOrder{
		Uid:  uid,
		Base: base,
	}
}

// Subaccount ID

func SplitSID(sid uint64) (
	ownerAccountNumber uint64,
	subaccountIdNumber uint32,
) {
	// ownerAccountNumber is the last 44 bits, big-endian.
	ownerAccountNumber = sid & 0x0000ffffffffffff

	// subaccountIdNumber is the first 20 bits, little-endian.
	subaccountIdNumber = uint32(bits.Reverse64(sid)) & uint32(0x0000ffff)

	return ownerAccountNumber, subaccountIdNumber
}

// Byte Keys

func (uid *XUID) ToBytes() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b, uid.Sid)
	binary.BigEndian.PutUint32(b[8:], uid.Iid.ClobId)
	binary.BigEndian.PutUint32(b[12:], uid.Iid.ClientId)
	return b
}

func (order *XOrder) Validate() error {
	if order.Base.HasOcoClientId() && order.Base.MustGetOcoClientId() == order.Uid.Iid.ClientId {
		return ErrOcoIsSameAsOrder
	}
	return nil
}

func (order *XOrder) ToOrdersKey() []byte {
	return order.Uid.ToBytes()
}

func (order *XOrder) ToStopKey(priority uint64) []byte {
	return ToStopKey(
		order.Uid.Iid.ClobId,
		order.Base.GetStopSide() == Order_SIDE_BUY,
		order.Base.MustGetStopSubticks(),
		priority,
	)
}

func ToStopKey(
	clobId uint32,
	isStopSideBuy bool,
	stopSubticks uint64,
	priority uint64,
) []byte {
	b := make([]byte, 21)
	binary.BigEndian.PutUint32(b, clobId)
	b[4] = lib.BoolToByte(isStopSideBuy)
	if isStopSideBuy {
		stopSubticks = ^stopSubticks
	}
	binary.BigEndian.PutUint64(b[5:], stopSubticks)
	binary.BigEndian.PutUint64(b[13:], priority)
	return b
}

func (order *XOrder) ToRestingKey(priority uint64) []byte {
	return ToRestingKey(
		order.Uid.Iid.ClobId,
		order.Base.Side() == Order_SIDE_BUY,
		order.Base.Subticks,
		priority,
	)
}

func ToRestingKey(
	clobId uint32,
	isBuy bool,
	subticks uint64,
	priority uint64,
) []byte {
	b := make([]byte, 21)
	binary.BigEndian.PutUint32(b, clobId)
	b[4] = lib.BoolToByte(isBuy)
	if isBuy {
		subticks = ^subticks
	}
	binary.BigEndian.PutUint64(b[5:], subticks)
	binary.BigEndian.PutUint64(b[13:], priority)
	return b
}

func (order *XOrder) ToExpiryKeyOrNil(priority uint64) []byte {
	if !order.Base.HasGoodTilTime() {
		return nil
	}
	b := make([]byte, 12)
	binary.BigEndian.PutUint32(b, order.Base.MustGetGoodTilTime())
	binary.BigEndian.PutUint64(b[4:], priority)
	return b
}

// Priority

// GetPriority returns a deterministic priority for the order.
// The priority is:
// - used to determine the order in which orders of the same price are matched or triggered
// - intended to be "fair" (averaging over a long period) between any two accounts on any orderbook
// - intended to be minimally-gameable by the client (i.e. not allow a better priority to be "mined" via clientId)
//   - the priority can be influenced by the clientId only in the least-significant 32 bits
//
// - must be deterministic so that keys that include the priority can be found when canceling the order
// - must be difficult to predict so that it is hard to front-run an order using an order of the same priority
func (order *XOrder) GetPriority() uint64 {
	b := make([]byte, 20)
	binary.BigEndian.PutUint64(b, order.Uid.Sid)
	binary.BigEndian.PutUint32(b[8:], order.Uid.Iid.ClobId)
	binary.BigEndian.PutUint64(b[12:], order.Base.Subticks)
	return xxhash.Sum64(b) + uint64(order.Uid.Iid.ClientId)
}

// Side

func (base *XBase) Side() Order_Side {
	if base.Flags&uint32(0b1) == uint32(1) {
		return Order_SIDE_BUY
	} else {
		return Order_SIDE_SELL
	}
}

// TIF (Time in Force)

func (base *XBase) GetTif() Order_TimeInForce {
	tbits := uint32(base.Flags>>1) & uint32(0b11)
	if tbits == uint32(0) {
		return Order_TIME_IN_FORCE_UNSPECIFIED
	}
	if tbits == uint32(1) {
		return Order_TIME_IN_FORCE_POST_ONLY
	}

	// 2 is IOC
	// 3 is IOC and reduce-only
	return Order_TIME_IN_FORCE_IOC
}

func (base *XBase) IsReduceOnly() bool {
	tbits := uint32(base.Flags>>1) & uint32(0b11)
	return tbits == uint32(3)
}

// Self-Trade Prevention

func (base *XBase) GetStp() Order_STP {
	tbits := uint32(base.Flags>>4) & uint32(0b11)
	return Order_STP(tbits)
}

// GTT (Good Til Time)

func (base *XBase) HasGoodTilTime() bool {
	return base.XGoodTilTime != nil
}

func (base *XBase) MustGetGoodTilTime() uint32 {
	if !base.HasGoodTilTime() {
		panic("MustGetGoodTilTime: order does not have a good til time")
	}
	return base.GetGoodTilTime()
}

// Stop Orders

func (base *XBase) IsStop() bool {
	return base.XStopSubticks != nil
}

func (base *XBase) MustGetStopSubticks() uint64 {
	if !base.IsStop() {
		panic("MustGetStopSubticks: order is not stop order")
	}
	return base.GetStopSubticks()
}

func (base *XBase) GetStopSide() Order_Side {
	tbits := uint32(base.Flags>>3) & uint32(0b1)
	if tbits == uint32(1) {
		return Order_SIDE_BUY
	} else {
		return Order_SIDE_SELL
	}
}

// OCO (One Cancels Other)

func (base *XBase) HasOcoClientId() bool {
	return base.XOcoClientId != nil
}

func (base *XBase) MustGetOcoClientId() uint32 {
	if !base.HasOcoClientId() {
		panic("MustGetOcoClientId: order does not have an OCO")
	}
	return base.GetOcoClientId()
}
