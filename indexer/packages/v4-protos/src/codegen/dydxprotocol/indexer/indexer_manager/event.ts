import { Timestamp } from "../../../google/protobuf/timestamp";
import { BinaryReader, BinaryWriter } from "../../../binary";
import { bytesFromBase64, base64FromBytes, toTimestamp, fromTimestamp } from "../../../helpers";
/** enum to specify that the IndexerTendermintEvent is a block event. */
export enum IndexerTendermintEvent_BlockEvent {
  /** BLOCK_EVENT_UNSPECIFIED - Default value. This value is invalid and unused. */
  BLOCK_EVENT_UNSPECIFIED = 0,
  /**
   * BLOCK_EVENT_BEGIN_BLOCK - BLOCK_EVENT_BEGIN_BLOCK indicates that the event was generated during
   * BeginBlock.
   */
  BLOCK_EVENT_BEGIN_BLOCK = 1,
  /**
   * BLOCK_EVENT_END_BLOCK - BLOCK_EVENT_END_BLOCK indicates that the event was generated during
   * EndBlock.
   */
  BLOCK_EVENT_END_BLOCK = 2,
  UNRECOGNIZED = -1,
}
export const IndexerTendermintEvent_BlockEventSDKType = IndexerTendermintEvent_BlockEvent;
export const IndexerTendermintEvent_BlockEventAmino = IndexerTendermintEvent_BlockEvent;
export function indexerTendermintEvent_BlockEventFromJSON(object: any): IndexerTendermintEvent_BlockEvent {
  switch (object) {
    case 0:
    case "BLOCK_EVENT_UNSPECIFIED":
      return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_UNSPECIFIED;
    case 1:
    case "BLOCK_EVENT_BEGIN_BLOCK":
      return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK;
    case 2:
    case "BLOCK_EVENT_END_BLOCK":
      return IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK;
    case -1:
    case "UNRECOGNIZED":
    default:
      return IndexerTendermintEvent_BlockEvent.UNRECOGNIZED;
  }
}
export function indexerTendermintEvent_BlockEventToJSON(object: IndexerTendermintEvent_BlockEvent): string {
  switch (object) {
    case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_UNSPECIFIED:
      return "BLOCK_EVENT_UNSPECIFIED";
    case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK:
      return "BLOCK_EVENT_BEGIN_BLOCK";
    case IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK:
      return "BLOCK_EVENT_END_BLOCK";
    case IndexerTendermintEvent_BlockEvent.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/**
 * IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
 * with an additional txn_hash field.
 */
export interface IndexerTendermintEventWrapper {
  event?: IndexerTendermintEvent;
  txnHash: string;
}
export interface IndexerTendermintEventWrapperProtoMsg {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEventWrapper";
  value: Uint8Array;
}
/**
 * IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
 * with an additional txn_hash field.
 */
export interface IndexerTendermintEventWrapperAmino {
  event?: IndexerTendermintEventAmino;
  txn_hash?: string;
}
export interface IndexerTendermintEventWrapperAminoMsg {
  type: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEventWrapper";
  value: IndexerTendermintEventWrapperAmino;
}
/**
 * IndexerTendermintEventWrapper is a wrapper around IndexerTendermintEvent,
 * with an additional txn_hash field.
 */
export interface IndexerTendermintEventWrapperSDKType {
  event?: IndexerTendermintEventSDKType;
  txn_hash: string;
}
/**
 * IndexerEventsStoreValue represents the type of the value of the
 * `IndexerEventsStore` in state.
 */
export interface IndexerEventsStoreValue {
  events: IndexerTendermintEventWrapper[];
}
export interface IndexerEventsStoreValueProtoMsg {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerEventsStoreValue";
  value: Uint8Array;
}
/**
 * IndexerEventsStoreValue represents the type of the value of the
 * `IndexerEventsStore` in state.
 */
export interface IndexerEventsStoreValueAmino {
  events?: IndexerTendermintEventWrapperAmino[];
}
export interface IndexerEventsStoreValueAminoMsg {
  type: "/dydxprotocol.indexer.indexer_manager.IndexerEventsStoreValue";
  value: IndexerEventsStoreValueAmino;
}
/**
 * IndexerEventsStoreValue represents the type of the value of the
 * `IndexerEventsStore` in state.
 */
export interface IndexerEventsStoreValueSDKType {
  events: IndexerTendermintEventWrapperSDKType[];
}
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */
export interface IndexerTendermintEvent {
  /** Subtype of the event e.g. "order_fill", "subaccount_update", etc. */
  subtype: string;
  transactionIndex?: number;
  blockEvent?: IndexerTendermintEvent_BlockEvent;
  /**
   * Index of the event within the list of events that happened either during a
   * transaction or during processing of a block.
   * TODO(DEC-537): Deprecate this field because events are already ordered.
   */
  eventIndex: number;
  /** Version of the event. */
  version: number;
  /** Tendermint event bytes. */
  dataBytes: Uint8Array;
}
export interface IndexerTendermintEventProtoMsg {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent";
  value: Uint8Array;
}
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */
export interface IndexerTendermintEventAmino {
  /** Subtype of the event e.g. "order_fill", "subaccount_update", etc. */
  subtype?: string;
  transaction_index?: number;
  block_event?: IndexerTendermintEvent_BlockEvent;
  /**
   * Index of the event within the list of events that happened either during a
   * transaction or during processing of a block.
   * TODO(DEC-537): Deprecate this field because events are already ordered.
   */
  event_index?: number;
  /** Version of the event. */
  version?: number;
  /** Tendermint event bytes. */
  data_bytes?: string;
}
export interface IndexerTendermintEventAminoMsg {
  type: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent";
  value: IndexerTendermintEventAmino;
}
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */
export interface IndexerTendermintEventSDKType {
  subtype: string;
  transaction_index?: number;
  block_event?: IndexerTendermintEvent_BlockEvent;
  event_index: number;
  version: number;
  data_bytes: Uint8Array;
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */
export interface IndexerTendermintBlock {
  height: number;
  time: Date;
  events: IndexerTendermintEvent[];
  txHashes: string[];
}
export interface IndexerTendermintBlockProtoMsg {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintBlock";
  value: Uint8Array;
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */
export interface IndexerTendermintBlockAmino {
  height?: number;
  time?: string;
  events?: IndexerTendermintEventAmino[];
  tx_hashes?: string[];
}
export interface IndexerTendermintBlockAminoMsg {
  type: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintBlock";
  value: IndexerTendermintBlockAmino;
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */
export interface IndexerTendermintBlockSDKType {
  height: number;
  time: Date;
  events: IndexerTendermintEventSDKType[];
  tx_hashes: string[];
}
function createBaseIndexerTendermintEventWrapper(): IndexerTendermintEventWrapper {
  return {
    event: undefined,
    txnHash: ""
  };
}
export const IndexerTendermintEventWrapper = {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEventWrapper",
  encode(message: IndexerTendermintEventWrapper, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.event !== undefined) {
      IndexerTendermintEvent.encode(message.event, writer.uint32(10).fork()).ldelim();
    }
    if (message.txnHash !== "") {
      writer.uint32(18).string(message.txnHash);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerTendermintEventWrapper {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerTendermintEventWrapper();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.event = IndexerTendermintEvent.decode(reader, reader.uint32());
          break;
        case 2:
          message.txnHash = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerTendermintEventWrapper>): IndexerTendermintEventWrapper {
    const message = createBaseIndexerTendermintEventWrapper();
    message.event = object.event !== undefined && object.event !== null ? IndexerTendermintEvent.fromPartial(object.event) : undefined;
    message.txnHash = object.txnHash ?? "";
    return message;
  },
  fromAmino(object: IndexerTendermintEventWrapperAmino): IndexerTendermintEventWrapper {
    const message = createBaseIndexerTendermintEventWrapper();
    if (object.event !== undefined && object.event !== null) {
      message.event = IndexerTendermintEvent.fromAmino(object.event);
    }
    if (object.txn_hash !== undefined && object.txn_hash !== null) {
      message.txnHash = object.txn_hash;
    }
    return message;
  },
  toAmino(message: IndexerTendermintEventWrapper): IndexerTendermintEventWrapperAmino {
    const obj: any = {};
    obj.event = message.event ? IndexerTendermintEvent.toAmino(message.event) : undefined;
    obj.txn_hash = message.txnHash;
    return obj;
  },
  fromAminoMsg(object: IndexerTendermintEventWrapperAminoMsg): IndexerTendermintEventWrapper {
    return IndexerTendermintEventWrapper.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerTendermintEventWrapperProtoMsg): IndexerTendermintEventWrapper {
    return IndexerTendermintEventWrapper.decode(message.value);
  },
  toProto(message: IndexerTendermintEventWrapper): Uint8Array {
    return IndexerTendermintEventWrapper.encode(message).finish();
  },
  toProtoMsg(message: IndexerTendermintEventWrapper): IndexerTendermintEventWrapperProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEventWrapper",
      value: IndexerTendermintEventWrapper.encode(message).finish()
    };
  }
};
function createBaseIndexerEventsStoreValue(): IndexerEventsStoreValue {
  return {
    events: []
  };
}
export const IndexerEventsStoreValue = {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerEventsStoreValue",
  encode(message: IndexerEventsStoreValue, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.events) {
      IndexerTendermintEventWrapper.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerEventsStoreValue {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerEventsStoreValue();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.events.push(IndexerTendermintEventWrapper.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerEventsStoreValue>): IndexerEventsStoreValue {
    const message = createBaseIndexerEventsStoreValue();
    message.events = object.events?.map(e => IndexerTendermintEventWrapper.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: IndexerEventsStoreValueAmino): IndexerEventsStoreValue {
    const message = createBaseIndexerEventsStoreValue();
    message.events = object.events?.map(e => IndexerTendermintEventWrapper.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: IndexerEventsStoreValue): IndexerEventsStoreValueAmino {
    const obj: any = {};
    if (message.events) {
      obj.events = message.events.map(e => e ? IndexerTendermintEventWrapper.toAmino(e) : undefined);
    } else {
      obj.events = [];
    }
    return obj;
  },
  fromAminoMsg(object: IndexerEventsStoreValueAminoMsg): IndexerEventsStoreValue {
    return IndexerEventsStoreValue.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerEventsStoreValueProtoMsg): IndexerEventsStoreValue {
    return IndexerEventsStoreValue.decode(message.value);
  },
  toProto(message: IndexerEventsStoreValue): Uint8Array {
    return IndexerEventsStoreValue.encode(message).finish();
  },
  toProtoMsg(message: IndexerEventsStoreValue): IndexerEventsStoreValueProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerEventsStoreValue",
      value: IndexerEventsStoreValue.encode(message).finish()
    };
  }
};
function createBaseIndexerTendermintEvent(): IndexerTendermintEvent {
  return {
    subtype: "",
    transactionIndex: undefined,
    blockEvent: undefined,
    eventIndex: 0,
    version: 0,
    dataBytes: new Uint8Array()
  };
}
export const IndexerTendermintEvent = {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent",
  encode(message: IndexerTendermintEvent, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.subtype !== "") {
      writer.uint32(10).string(message.subtype);
    }
    if (message.transactionIndex !== undefined) {
      writer.uint32(24).uint32(message.transactionIndex);
    }
    if (message.blockEvent !== undefined) {
      writer.uint32(32).int32(message.blockEvent);
    }
    if (message.eventIndex !== 0) {
      writer.uint32(40).uint32(message.eventIndex);
    }
    if (message.version !== 0) {
      writer.uint32(48).uint32(message.version);
    }
    if (message.dataBytes.length !== 0) {
      writer.uint32(58).bytes(message.dataBytes);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerTendermintEvent {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerTendermintEvent();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.subtype = reader.string();
          break;
        case 3:
          message.transactionIndex = reader.uint32();
          break;
        case 4:
          message.blockEvent = (reader.int32() as any);
          break;
        case 5:
          message.eventIndex = reader.uint32();
          break;
        case 6:
          message.version = reader.uint32();
          break;
        case 7:
          message.dataBytes = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerTendermintEvent>): IndexerTendermintEvent {
    const message = createBaseIndexerTendermintEvent();
    message.subtype = object.subtype ?? "";
    message.transactionIndex = object.transactionIndex ?? undefined;
    message.blockEvent = object.blockEvent ?? undefined;
    message.eventIndex = object.eventIndex ?? 0;
    message.version = object.version ?? 0;
    message.dataBytes = object.dataBytes ?? new Uint8Array();
    return message;
  },
  fromAmino(object: IndexerTendermintEventAmino): IndexerTendermintEvent {
    const message = createBaseIndexerTendermintEvent();
    if (object.subtype !== undefined && object.subtype !== null) {
      message.subtype = object.subtype;
    }
    if (object.transaction_index !== undefined && object.transaction_index !== null) {
      message.transactionIndex = object.transaction_index;
    }
    if (object.block_event !== undefined && object.block_event !== null) {
      message.blockEvent = indexerTendermintEvent_BlockEventFromJSON(object.block_event);
    }
    if (object.event_index !== undefined && object.event_index !== null) {
      message.eventIndex = object.event_index;
    }
    if (object.version !== undefined && object.version !== null) {
      message.version = object.version;
    }
    if (object.data_bytes !== undefined && object.data_bytes !== null) {
      message.dataBytes = bytesFromBase64(object.data_bytes);
    }
    return message;
  },
  toAmino(message: IndexerTendermintEvent): IndexerTendermintEventAmino {
    const obj: any = {};
    obj.subtype = message.subtype;
    obj.transaction_index = message.transactionIndex;
    obj.block_event = indexerTendermintEvent_BlockEventToJSON(message.blockEvent);
    obj.event_index = message.eventIndex;
    obj.version = message.version;
    obj.data_bytes = message.dataBytes ? base64FromBytes(message.dataBytes) : undefined;
    return obj;
  },
  fromAminoMsg(object: IndexerTendermintEventAminoMsg): IndexerTendermintEvent {
    return IndexerTendermintEvent.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerTendermintEventProtoMsg): IndexerTendermintEvent {
    return IndexerTendermintEvent.decode(message.value);
  },
  toProto(message: IndexerTendermintEvent): Uint8Array {
    return IndexerTendermintEvent.encode(message).finish();
  },
  toProtoMsg(message: IndexerTendermintEvent): IndexerTendermintEventProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintEvent",
      value: IndexerTendermintEvent.encode(message).finish()
    };
  }
};
function createBaseIndexerTendermintBlock(): IndexerTendermintBlock {
  return {
    height: 0,
    time: new Date(),
    events: [],
    txHashes: []
  };
}
export const IndexerTendermintBlock = {
  typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintBlock",
  encode(message: IndexerTendermintBlock, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.height !== 0) {
      writer.uint32(8).uint32(message.height);
    }
    if (message.time !== undefined) {
      Timestamp.encode(toTimestamp(message.time), writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.events) {
      IndexerTendermintEvent.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    for (const v of message.txHashes) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): IndexerTendermintBlock {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIndexerTendermintBlock();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = reader.uint32();
          break;
        case 2:
          message.time = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          break;
        case 3:
          message.events.push(IndexerTendermintEvent.decode(reader, reader.uint32()));
          break;
        case 4:
          message.txHashes.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<IndexerTendermintBlock>): IndexerTendermintBlock {
    const message = createBaseIndexerTendermintBlock();
    message.height = object.height ?? 0;
    message.time = object.time ?? undefined;
    message.events = object.events?.map(e => IndexerTendermintEvent.fromPartial(e)) || [];
    message.txHashes = object.txHashes?.map(e => e) || [];
    return message;
  },
  fromAmino(object: IndexerTendermintBlockAmino): IndexerTendermintBlock {
    const message = createBaseIndexerTendermintBlock();
    if (object.height !== undefined && object.height !== null) {
      message.height = object.height;
    }
    if (object.time !== undefined && object.time !== null) {
      message.time = fromTimestamp(Timestamp.fromAmino(object.time));
    }
    message.events = object.events?.map(e => IndexerTendermintEvent.fromAmino(e)) || [];
    message.txHashes = object.tx_hashes?.map(e => e) || [];
    return message;
  },
  toAmino(message: IndexerTendermintBlock): IndexerTendermintBlockAmino {
    const obj: any = {};
    obj.height = message.height;
    obj.time = message.time ? Timestamp.toAmino(toTimestamp(message.time)) : undefined;
    if (message.events) {
      obj.events = message.events.map(e => e ? IndexerTendermintEvent.toAmino(e) : undefined);
    } else {
      obj.events = [];
    }
    if (message.txHashes) {
      obj.tx_hashes = message.txHashes.map(e => e);
    } else {
      obj.tx_hashes = [];
    }
    return obj;
  },
  fromAminoMsg(object: IndexerTendermintBlockAminoMsg): IndexerTendermintBlock {
    return IndexerTendermintBlock.fromAmino(object.value);
  },
  fromProtoMsg(message: IndexerTendermintBlockProtoMsg): IndexerTendermintBlock {
    return IndexerTendermintBlock.decode(message.value);
  },
  toProto(message: IndexerTendermintBlock): Uint8Array {
    return IndexerTendermintBlock.encode(message).finish();
  },
  toProtoMsg(message: IndexerTendermintBlock): IndexerTendermintBlockProtoMsg {
    return {
      typeUrl: "/dydxprotocol.indexer.indexer_manager.IndexerTendermintBlock",
      value: IndexerTendermintBlock.encode(message).finish()
    };
  }
};