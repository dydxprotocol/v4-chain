// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: dydxprotocol/indexer/indexer_manager/event.proto

package indexer_manager

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/cosmos/gogoproto/types"
	github_com_cosmos_gogoproto_types "github.com/cosmos/gogoproto/types"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// enum to specify that the IndexerTendermintEvent is a block event.
type IndexerTendermintEvent_BlockEvent int32

const (
	// Default value. This value is invalid and unused.
	IndexerTendermintEvent_BLOCK_EVENT_UNSPECIFIED IndexerTendermintEvent_BlockEvent = 0
	// BLOCK_EVENT_BEGIN_BLOCK indicates that the event was generated during
	// BeginBlock.
	IndexerTendermintEvent_BLOCK_EVENT_BEGIN_BLOCK IndexerTendermintEvent_BlockEvent = 1
	// BLOCK_EVENT_END_BLOCK indicates that the event was generated during
	// EndBlock.
	IndexerTendermintEvent_BLOCK_EVENT_END_BLOCK IndexerTendermintEvent_BlockEvent = 2
)

var IndexerTendermintEvent_BlockEvent_name = map[int32]string{
	0: "BLOCK_EVENT_UNSPECIFIED",
	1: "BLOCK_EVENT_BEGIN_BLOCK",
	2: "BLOCK_EVENT_END_BLOCK",
}

var IndexerTendermintEvent_BlockEvent_value = map[string]int32{
	"BLOCK_EVENT_UNSPECIFIED": 0,
	"BLOCK_EVENT_BEGIN_BLOCK": 1,
	"BLOCK_EVENT_END_BLOCK":   2,
}

func (x IndexerTendermintEvent_BlockEvent) String() string {
	return proto.EnumName(IndexerTendermintEvent_BlockEvent_name, int32(x))
}

func (IndexerTendermintEvent_BlockEvent) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_18a6a94c31da6b1f, []int{2, 0}
}

// IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
// with an additional txn_hash field.
type IndexerTendermintEventWrapper struct {
	Event   *IndexerTendermintEvent `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	TxnHash string                  `protobuf:"bytes,2,opt,name=txn_hash,json=txnHash,proto3" json:"txn_hash,omitempty"`
}

func (m *IndexerTendermintEventWrapper) Reset()         { *m = IndexerTendermintEventWrapper{} }
func (m *IndexerTendermintEventWrapper) String() string { return proto.CompactTextString(m) }
func (*IndexerTendermintEventWrapper) ProtoMessage()    {}
func (*IndexerTendermintEventWrapper) Descriptor() ([]byte, []int) {
	return fileDescriptor_18a6a94c31da6b1f, []int{0}
}
func (m *IndexerTendermintEventWrapper) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerTendermintEventWrapper) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerTendermintEventWrapper.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerTendermintEventWrapper) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerTendermintEventWrapper.Merge(m, src)
}
func (m *IndexerTendermintEventWrapper) XXX_Size() int {
	return m.Size()
}
func (m *IndexerTendermintEventWrapper) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerTendermintEventWrapper.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerTendermintEventWrapper proto.InternalMessageInfo

func (m *IndexerTendermintEventWrapper) GetEvent() *IndexerTendermintEvent {
	if m != nil {
		return m.Event
	}
	return nil
}

func (m *IndexerTendermintEventWrapper) GetTxnHash() string {
	if m != nil {
		return m.TxnHash
	}
	return ""
}

// IndexerEventsStoreValue represents the type of the value of the
// `IndexerEventsStore` in state.
type IndexerEventsStoreValue struct {
	Events []*IndexerTendermintEventWrapper `protobuf:"bytes,1,rep,name=events,proto3" json:"events,omitempty"`
}

func (m *IndexerEventsStoreValue) Reset()         { *m = IndexerEventsStoreValue{} }
func (m *IndexerEventsStoreValue) String() string { return proto.CompactTextString(m) }
func (*IndexerEventsStoreValue) ProtoMessage()    {}
func (*IndexerEventsStoreValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_18a6a94c31da6b1f, []int{1}
}
func (m *IndexerEventsStoreValue) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerEventsStoreValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerEventsStoreValue.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerEventsStoreValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerEventsStoreValue.Merge(m, src)
}
func (m *IndexerEventsStoreValue) XXX_Size() int {
	return m.Size()
}
func (m *IndexerEventsStoreValue) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerEventsStoreValue.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerEventsStoreValue proto.InternalMessageInfo

func (m *IndexerEventsStoreValue) GetEvents() []*IndexerTendermintEventWrapper {
	if m != nil {
		return m.Events
	}
	return nil
}

// IndexerTendermintEvent contains the base64 encoded event proto emitted from
// the V4 application as well as additional metadata to determine the ordering
// of the event within the block and the subtype of the event.
type IndexerTendermintEvent struct {
	// Subtype of the event e.g. "order_fill", "subaccount_update", etc.
	Subtype string `protobuf:"bytes,1,opt,name=subtype,proto3" json:"subtype,omitempty"`
	// Base64 encoded proto from the Tendermint event.
	// TODO(DEC-1720): Change to bytes post-migration.
	Data string `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// ordering_within_block is either the transaction index or a boolean
	// indicating the event was generated during processing the block rather than
	// any specific transaction e.g. during FinalizeBlock.
	//
	// Types that are valid to be assigned to OrderingWithinBlock:
	//
	//	*IndexerTendermintEvent_TransactionIndex
	//	*IndexerTendermintEvent_BlockEvent_
	OrderingWithinBlock isIndexerTendermintEvent_OrderingWithinBlock `protobuf_oneof:"ordering_within_block"`
	// Index of the event within the list of events that happened either during a
	// transaction or during processing of a block.
	// TODO(DEC-537): Deprecate this field because events are already ordered.
	EventIndex uint32 `protobuf:"varint,5,opt,name=event_index,json=eventIndex,proto3" json:"event_index,omitempty"`
	// Version of the event.
	Version uint32 `protobuf:"varint,6,opt,name=version,proto3" json:"version,omitempty"`
	// Tendermint event bytes.
	DataBytes []byte `protobuf:"bytes,7,opt,name=data_bytes,json=dataBytes,proto3" json:"data_bytes,omitempty"`
}

func (m *IndexerTendermintEvent) Reset()         { *m = IndexerTendermintEvent{} }
func (m *IndexerTendermintEvent) String() string { return proto.CompactTextString(m) }
func (*IndexerTendermintEvent) ProtoMessage()    {}
func (*IndexerTendermintEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_18a6a94c31da6b1f, []int{2}
}
func (m *IndexerTendermintEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerTendermintEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerTendermintEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerTendermintEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerTendermintEvent.Merge(m, src)
}
func (m *IndexerTendermintEvent) XXX_Size() int {
	return m.Size()
}
func (m *IndexerTendermintEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerTendermintEvent.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerTendermintEvent proto.InternalMessageInfo

type isIndexerTendermintEvent_OrderingWithinBlock interface {
	isIndexerTendermintEvent_OrderingWithinBlock()
	MarshalTo([]byte) (int, error)
	Size() int
}

type IndexerTendermintEvent_TransactionIndex struct {
	TransactionIndex uint32 `protobuf:"varint,3,opt,name=transaction_index,json=transactionIndex,proto3,oneof" json:"transaction_index,omitempty"`
}
type IndexerTendermintEvent_BlockEvent_ struct {
	BlockEvent IndexerTendermintEvent_BlockEvent `protobuf:"varint,4,opt,name=block_event,json=blockEvent,proto3,enum=dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent_BlockEvent,oneof" json:"block_event,omitempty"`
}

func (*IndexerTendermintEvent_TransactionIndex) isIndexerTendermintEvent_OrderingWithinBlock() {}
func (*IndexerTendermintEvent_BlockEvent_) isIndexerTendermintEvent_OrderingWithinBlock()      {}

func (m *IndexerTendermintEvent) GetOrderingWithinBlock() isIndexerTendermintEvent_OrderingWithinBlock {
	if m != nil {
		return m.OrderingWithinBlock
	}
	return nil
}

func (m *IndexerTendermintEvent) GetSubtype() string {
	if m != nil {
		return m.Subtype
	}
	return ""
}

func (m *IndexerTendermintEvent) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *IndexerTendermintEvent) GetTransactionIndex() uint32 {
	if x, ok := m.GetOrderingWithinBlock().(*IndexerTendermintEvent_TransactionIndex); ok {
		return x.TransactionIndex
	}
	return 0
}

func (m *IndexerTendermintEvent) GetBlockEvent() IndexerTendermintEvent_BlockEvent {
	if x, ok := m.GetOrderingWithinBlock().(*IndexerTendermintEvent_BlockEvent_); ok {
		return x.BlockEvent
	}
	return IndexerTendermintEvent_BLOCK_EVENT_UNSPECIFIED
}

func (m *IndexerTendermintEvent) GetEventIndex() uint32 {
	if m != nil {
		return m.EventIndex
	}
	return 0
}

func (m *IndexerTendermintEvent) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *IndexerTendermintEvent) GetDataBytes() []byte {
	if m != nil {
		return m.DataBytes
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*IndexerTendermintEvent) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*IndexerTendermintEvent_TransactionIndex)(nil),
		(*IndexerTendermintEvent_BlockEvent_)(nil),
	}
}

// IndexerTendermintBlock contains all the events for the block along with
// metadata for the block height, timestamp of the block and a list of all the
// hashes of the transactions within the block. The transaction hashes follow
// the ordering of the transactions as they appear within the block.
type IndexerTendermintBlock struct {
	Height   uint32                    `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	Time     time.Time                 `protobuf:"bytes,2,opt,name=time,proto3,stdtime" json:"time"`
	Events   []*IndexerTendermintEvent `protobuf:"bytes,3,rep,name=events,proto3" json:"events,omitempty"`
	TxHashes []string                  `protobuf:"bytes,4,rep,name=tx_hashes,json=txHashes,proto3" json:"tx_hashes,omitempty"`
}

func (m *IndexerTendermintBlock) Reset()         { *m = IndexerTendermintBlock{} }
func (m *IndexerTendermintBlock) String() string { return proto.CompactTextString(m) }
func (*IndexerTendermintBlock) ProtoMessage()    {}
func (*IndexerTendermintBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_18a6a94c31da6b1f, []int{3}
}
func (m *IndexerTendermintBlock) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *IndexerTendermintBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_IndexerTendermintBlock.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *IndexerTendermintBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IndexerTendermintBlock.Merge(m, src)
}
func (m *IndexerTendermintBlock) XXX_Size() int {
	return m.Size()
}
func (m *IndexerTendermintBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_IndexerTendermintBlock.DiscardUnknown(m)
}

var xxx_messageInfo_IndexerTendermintBlock proto.InternalMessageInfo

func (m *IndexerTendermintBlock) GetHeight() uint32 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *IndexerTendermintBlock) GetTime() time.Time {
	if m != nil {
		return m.Time
	}
	return time.Time{}
}

func (m *IndexerTendermintBlock) GetEvents() []*IndexerTendermintEvent {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *IndexerTendermintBlock) GetTxHashes() []string {
	if m != nil {
		return m.TxHashes
	}
	return nil
}

func init() {
	proto.RegisterEnum("dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent_BlockEvent", IndexerTendermintEvent_BlockEvent_name, IndexerTendermintEvent_BlockEvent_value)
	proto.RegisterType((*IndexerTendermintEventWrapper)(nil), "dydxprotocol.indexer.indexer_manager.IndexerTendermintEventWrapper")
	proto.RegisterType((*IndexerEventsStoreValue)(nil), "dydxprotocol.indexer.indexer_manager.IndexerEventsStoreValue")
	proto.RegisterType((*IndexerTendermintEvent)(nil), "dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent")
	proto.RegisterType((*IndexerTendermintBlock)(nil), "dydxprotocol.indexer.indexer_manager.IndexerTendermintBlock")
}

func init() {
	proto.RegisterFile("dydxprotocol/indexer/indexer_manager/event.proto", fileDescriptor_18a6a94c31da6b1f)
}

var fileDescriptor_18a6a94c31da6b1f = []byte{
	// 597 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcd, 0x6e, 0xd3, 0x4c,
	0x14, 0xf5, 0x34, 0x69, 0xd2, 0xdc, 0x7c, 0xfd, 0x54, 0x46, 0xb4, 0x75, 0x5b, 0x35, 0xb1, 0x22,
	0x16, 0xde, 0xd4, 0x46, 0x81, 0x05, 0x0b, 0x36, 0xb8, 0x35, 0x4d, 0x04, 0x0a, 0x68, 0x1a, 0x8a,
	0x44, 0x17, 0xd6, 0x38, 0x19, 0x6c, 0x43, 0x32, 0x13, 0xd9, 0x93, 0xe0, 0xbe, 0x04, 0xea, 0xdb,
	0xf0, 0x0a, 0x5d, 0x76, 0xc9, 0x86, 0x1f, 0xb5, 0x2f, 0x82, 0x3c, 0x76, 0x4a, 0xa9, 0x2a, 0x81,
	0x94, 0x95, 0xef, 0xef, 0xb9, 0xe7, 0x7a, 0xe6, 0x0c, 0x3c, 0x1c, 0x9e, 0x0e, 0xd3, 0x49, 0x2c,
	0xa4, 0x18, 0x88, 0x91, 0x1d, 0xf1, 0x21, 0x4b, 0x59, 0x3c, 0xff, 0x7a, 0x63, 0xca, 0x69, 0xc0,
	0x62, 0x9b, 0xcd, 0x18, 0x97, 0x96, 0x2a, 0xc3, 0x0f, 0x6e, 0x76, 0x58, 0x45, 0xa5, 0x75, 0xab,
	0x63, 0xbb, 0x19, 0x08, 0x11, 0x8c, 0x98, 0xad, 0x0a, 0xfd, 0xe9, 0x7b, 0x5b, 0x46, 0x63, 0x96,
	0x48, 0x3a, 0x9e, 0xe4, 0x30, 0xdb, 0xf7, 0x03, 0x11, 0x08, 0x65, 0xda, 0x99, 0x95, 0x47, 0x5b,
	0x9f, 0x11, 0xec, 0x76, 0x73, 0xa8, 0x3e, 0xe3, 0x43, 0x16, 0x8f, 0x23, 0x2e, 0xdd, 0x6c, 0xfa,
	0xdb, 0x98, 0x4e, 0x26, 0x2c, 0xc6, 0x04, 0x96, 0x15, 0x1b, 0x1d, 0x19, 0xc8, 0xac, 0xb7, 0x9f,
	0x5a, 0xff, 0x42, 0xc7, 0xba, 0x1b, 0x93, 0xe4, 0x50, 0x78, 0x0b, 0x56, 0x64, 0xca, 0xbd, 0x90,
	0x26, 0xa1, 0xbe, 0x64, 0x20, 0xb3, 0x46, 0xaa, 0x32, 0xe5, 0x1d, 0x9a, 0x84, 0xad, 0x19, 0x6c,
	0x16, 0xbd, 0xaa, 0x23, 0x39, 0x92, 0x22, 0x66, 0xc7, 0x74, 0x34, 0x65, 0xf8, 0x04, 0x2a, 0xaa,
	0x3d, 0xd1, 0x91, 0x51, 0x32, 0xeb, 0xed, 0xfd, 0x45, 0xa8, 0x14, 0xeb, 0x91, 0x02, 0xb2, 0xf5,
	0xa5, 0x04, 0x1b, 0x77, 0x57, 0x62, 0x1d, 0xaa, 0xc9, 0xd4, 0x97, 0xa7, 0x13, 0xa6, 0xfe, 0x41,
	0x8d, 0xcc, 0x5d, 0x8c, 0xa1, 0x3c, 0xa4, 0x92, 0x16, 0x3b, 0x28, 0x1b, 0xef, 0xc1, 0x3d, 0x19,
	0x53, 0x9e, 0xd0, 0x81, 0x8c, 0x04, 0xf7, 0x14, 0x1b, 0xbd, 0x64, 0x20, 0x73, 0xb5, 0xa3, 0x91,
	0xb5, 0x1b, 0x29, 0x35, 0x0d, 0x7f, 0x80, 0xba, 0x3f, 0x12, 0x83, 0x8f, 0x5e, 0xfe, 0x93, 0xcb,
	0x06, 0x32, 0xff, 0x6f, 0x1f, 0x2e, 0xb2, 0x99, 0xe5, 0x64, 0x78, 0xca, 0xec, 0x68, 0x04, 0xfc,
	0x6b, 0x0f, 0x37, 0xa1, 0xae, 0xa6, 0x14, 0xa4, 0x96, 0x33, 0x52, 0x04, 0x54, 0x28, 0x27, 0xa3,
	0x43, 0x75, 0xc6, 0xe2, 0x24, 0x12, 0x5c, 0xaf, 0xa8, 0xe4, 0xdc, 0xc5, 0xbb, 0x00, 0xd9, 0x76,
	0x9e, 0x7f, 0x2a, 0x59, 0xa2, 0x57, 0x0d, 0x64, 0xfe, 0x47, 0x6a, 0x59, 0xc4, 0xc9, 0x02, 0x2d,
	0x0a, 0xf0, 0x7b, 0x2a, 0xde, 0x81, 0x4d, 0xe7, 0xe5, 0xab, 0xfd, 0x17, 0x9e, 0x7b, 0xec, 0xf6,
	0xfa, 0xde, 0x9b, 0xde, 0xd1, 0x6b, 0x77, 0xbf, 0xfb, 0xbc, 0xeb, 0x1e, 0xac, 0x69, 0xb7, 0x93,
	0x8e, 0x7b, 0xd8, 0xed, 0x79, 0x2a, 0xb2, 0x86, 0xf0, 0x16, 0xac, 0xdf, 0x4c, 0xba, 0xbd, 0x83,
	0x22, 0xb5, 0xe4, 0x6c, 0xc2, 0xba, 0x88, 0x87, 0x2c, 0x8e, 0x78, 0xe0, 0x7d, 0x8a, 0x64, 0x18,
	0x71, 0x4f, 0xad, 0xd6, 0xfa, 0x86, 0xee, 0x38, 0x39, 0xc5, 0x06, 0x6f, 0x40, 0x25, 0x64, 0x51,
	0x10, 0xe6, 0x97, 0x77, 0x95, 0x14, 0x1e, 0x7e, 0x02, 0xe5, 0x4c, 0x1e, 0xea, 0xdc, 0xea, 0xed,
	0x6d, 0x2b, 0xd7, 0x8e, 0x35, 0xd7, 0x8e, 0xd5, 0x9f, 0x6b, 0xc7, 0x59, 0x39, 0xff, 0xde, 0xd4,
	0xce, 0x7e, 0x34, 0x11, 0x51, 0x1d, 0xb8, 0x7f, 0x7d, 0x07, 0x4b, 0xea, 0x0e, 0x2e, 0x26, 0x87,
	0x02, 0x0b, 0xef, 0x40, 0x4d, 0xa6, 0x4a, 0x0e, 0x2c, 0xd1, 0xcb, 0x46, 0xc9, 0xac, 0x91, 0x15,
	0x99, 0x76, 0x94, 0xef, 0x9c, 0x9c, 0x5f, 0x36, 0xd0, 0xc5, 0x65, 0x03, 0xfd, 0xbc, 0x6c, 0xa0,
	0xb3, 0xab, 0x86, 0x76, 0x71, 0xd5, 0xd0, 0xbe, 0x5e, 0x35, 0xb4, 0x77, 0xcf, 0x82, 0x48, 0x86,
	0x53, 0xdf, 0x1a, 0x88, 0xb1, 0xfd, 0xc7, 0xb3, 0x32, 0x7b, 0xbc, 0x37, 0x08, 0x69, 0xc4, 0xed,
	0xbf, 0x3d, 0x34, 0x7e, 0x45, 0x55, 0x3c, 0xfa, 0x15, 0x00, 0x00, 0xff, 0xff, 0x59, 0xcc, 0xe8,
	0x07, 0x97, 0x04, 0x00, 0x00,
}

func (m *IndexerTendermintEventWrapper) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerTendermintEventWrapper) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerTendermintEventWrapper) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxnHash) > 0 {
		i -= len(m.TxnHash)
		copy(dAtA[i:], m.TxnHash)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.TxnHash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Event != nil {
		{
			size, err := m.Event.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintEvent(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IndexerEventsStoreValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerEventsStoreValue) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerEventsStoreValue) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintEvent(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *IndexerTendermintEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerTendermintEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerTendermintEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DataBytes) > 0 {
		i -= len(m.DataBytes)
		copy(dAtA[i:], m.DataBytes)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.DataBytes)))
		i--
		dAtA[i] = 0x3a
	}
	if m.Version != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x30
	}
	if m.EventIndex != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.EventIndex))
		i--
		dAtA[i] = 0x28
	}
	if m.OrderingWithinBlock != nil {
		{
			size := m.OrderingWithinBlock.Size()
			i -= size
			if _, err := m.OrderingWithinBlock.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Subtype) > 0 {
		i -= len(m.Subtype)
		copy(dAtA[i:], m.Subtype)
		i = encodeVarintEvent(dAtA, i, uint64(len(m.Subtype)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *IndexerTendermintEvent_TransactionIndex) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerTendermintEvent_TransactionIndex) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintEvent(dAtA, i, uint64(m.TransactionIndex))
	i--
	dAtA[i] = 0x18
	return len(dAtA) - i, nil
}
func (m *IndexerTendermintEvent_BlockEvent_) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerTendermintEvent_BlockEvent_) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintEvent(dAtA, i, uint64(m.BlockEvent))
	i--
	dAtA[i] = 0x20
	return len(dAtA) - i, nil
}
func (m *IndexerTendermintBlock) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *IndexerTendermintBlock) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *IndexerTendermintBlock) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.TxHashes) > 0 {
		for iNdEx := len(m.TxHashes) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.TxHashes[iNdEx])
			copy(dAtA[i:], m.TxHashes[iNdEx])
			i = encodeVarintEvent(dAtA, i, uint64(len(m.TxHashes[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Events) > 0 {
		for iNdEx := len(m.Events) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Events[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintEvent(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	n2, err2 := github_com_cosmos_gogoproto_types.StdTimeMarshalTo(m.Time, dAtA[i-github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintEvent(dAtA, i, uint64(n2))
	i--
	dAtA[i] = 0x12
	if m.Height != 0 {
		i = encodeVarintEvent(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintEvent(dAtA []byte, offset int, v uint64) int {
	offset -= sovEvent(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *IndexerTendermintEventWrapper) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Event != nil {
		l = m.Event.Size()
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.TxnHash)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func (m *IndexerEventsStoreValue) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovEvent(uint64(l))
		}
	}
	return n
}

func (m *IndexerTendermintEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Subtype)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	if m.OrderingWithinBlock != nil {
		n += m.OrderingWithinBlock.Size()
	}
	if m.EventIndex != 0 {
		n += 1 + sovEvent(uint64(m.EventIndex))
	}
	if m.Version != 0 {
		n += 1 + sovEvent(uint64(m.Version))
	}
	l = len(m.DataBytes)
	if l > 0 {
		n += 1 + l + sovEvent(uint64(l))
	}
	return n
}

func (m *IndexerTendermintEvent_TransactionIndex) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovEvent(uint64(m.TransactionIndex))
	return n
}
func (m *IndexerTendermintEvent_BlockEvent_) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovEvent(uint64(m.BlockEvent))
	return n
}
func (m *IndexerTendermintBlock) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovEvent(uint64(m.Height))
	}
	l = github_com_cosmos_gogoproto_types.SizeOfStdTime(m.Time)
	n += 1 + l + sovEvent(uint64(l))
	if len(m.Events) > 0 {
		for _, e := range m.Events {
			l = e.Size()
			n += 1 + l + sovEvent(uint64(l))
		}
	}
	if len(m.TxHashes) > 0 {
		for _, s := range m.TxHashes {
			l = len(s)
			n += 1 + l + sovEvent(uint64(l))
		}
	}
	return n
}

func sovEvent(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozEvent(x uint64) (n int) {
	return sovEvent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *IndexerTendermintEventWrapper) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IndexerTendermintEventWrapper: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerTendermintEventWrapper: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Event", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Event == nil {
				m.Event = &IndexerTendermintEvent{}
			}
			if err := m.Event.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxnHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxnHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IndexerEventsStoreValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IndexerEventsStoreValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerEventsStoreValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, &IndexerTendermintEventWrapper{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IndexerTendermintEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IndexerTendermintEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerTendermintEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subtype", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subtype = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TransactionIndex", wireType)
			}
			var v uint32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.OrderingWithinBlock = &IndexerTendermintEvent_TransactionIndex{v}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockEvent", wireType)
			}
			var v IndexerTendermintEvent_BlockEvent
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= IndexerTendermintEvent_BlockEvent(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.OrderingWithinBlock = &IndexerTendermintEvent_BlockEvent_{v}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EventIndex", wireType)
			}
			m.EventIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EventIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DataBytes", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DataBytes = append(m.DataBytes[:0], dAtA[iNdEx:postIndex]...)
			if m.DataBytes == nil {
				m.DataBytes = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *IndexerTendermintBlock) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: IndexerTendermintBlock: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: IndexerTendermintBlock: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_cosmos_gogoproto_types.StdTimeUnmarshal(&m.Time, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Events", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Events = append(m.Events, &IndexerTendermintEvent{})
			if err := m.Events[len(m.Events)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxHashes", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthEvent
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthEvent
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TxHashes = append(m.TxHashes, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipEvent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthEvent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipEvent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowEvent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowEvent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthEvent
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupEvent
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthEvent
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthEvent        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowEvent          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupEvent = fmt.Errorf("proto: unexpected end of group")
)
