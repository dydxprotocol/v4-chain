import { Timestamp } from "../../../google/protobuf/timestamp";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, toTimestamp, fromTimestamp } from "../../../helpers";
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
/** enum to specify that the IndexerTendermintEvent is a block event. */

export enum IndexerTendermintEvent_BlockEventSDKType {
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
/**
 * IndexerTendermintEvent contains the base64 encoded event proto emitted from
 * the dYdX application as well as additional metadata to determine the ordering
 * of the event within the block and the subtype of the event.
 */

export interface IndexerTendermintEventSDKType {
  /** Subtype of the event e.g. "order_fill", "subaccount_update", etc. */
  subtype: string;
  transaction_index?: number;
  block_event?: IndexerTendermintEvent_BlockEventSDKType;
  /**
   * Index of the event within the list of events that happened either during a
   * transaction or during processing of a block.
   * TODO(DEC-537): Deprecate this field because events are already ordered.
   */

  event_index: number;
  /** Version of the event. */

  version: number;
  /** Tendermint event bytes. */

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
  time?: Date;
  events: IndexerTendermintEvent[];
  txHashes: string[];
}
/**
 * IndexerTendermintBlock contains all the events for the block along with
 * metadata for the block height, timestamp of the block and a list of all the
 * hashes of the transactions within the block. The transaction hashes follow
 * the ordering of the transactions as they appear within the block.
 */

export interface IndexerTendermintBlockSDKType {
  height: number;
  time?: Date;
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
  encode(message: IndexerTendermintEventWrapper, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.event !== undefined) {
      IndexerTendermintEvent.encode(message.event, writer.uint32(10).fork()).ldelim();
    }

    if (message.txnHash !== "") {
      writer.uint32(18).string(message.txnHash);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintEventWrapper {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerTendermintEventWrapper>): IndexerTendermintEventWrapper {
    const message = createBaseIndexerTendermintEventWrapper();
    message.event = object.event !== undefined && object.event !== null ? IndexerTendermintEvent.fromPartial(object.event) : undefined;
    message.txnHash = object.txnHash ?? "";
    return message;
  }

};

function createBaseIndexerEventsStoreValue(): IndexerEventsStoreValue {
  return {
    events: []
  };
}

export const IndexerEventsStoreValue = {
  encode(message: IndexerEventsStoreValue, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.events) {
      IndexerTendermintEventWrapper.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerEventsStoreValue {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerEventsStoreValue>): IndexerEventsStoreValue {
    const message = createBaseIndexerEventsStoreValue();
    message.events = object.events?.map(e => IndexerTendermintEventWrapper.fromPartial(e)) || [];
    return message;
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
  encode(message: IndexerTendermintEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerTendermintEvent>): IndexerTendermintEvent {
    const message = createBaseIndexerTendermintEvent();
    message.subtype = object.subtype ?? "";
    message.transactionIndex = object.transactionIndex ?? undefined;
    message.blockEvent = object.blockEvent ?? undefined;
    message.eventIndex = object.eventIndex ?? 0;
    message.version = object.version ?? 0;
    message.dataBytes = object.dataBytes ?? new Uint8Array();
    return message;
  }

};

function createBaseIndexerTendermintBlock(): IndexerTendermintBlock {
  return {
    height: 0,
    time: undefined,
    events: [],
    txHashes: []
  };
}

export const IndexerTendermintBlock = {
  encode(message: IndexerTendermintBlock, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): IndexerTendermintBlock {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<IndexerTendermintBlock>): IndexerTendermintBlock {
    const message = createBaseIndexerTendermintBlock();
    message.height = object.height ?? 0;
    message.time = object.time ?? undefined;
    message.events = object.events?.map(e => IndexerTendermintEvent.fromPartial(e)) || [];
    message.txHashes = object.txHashes?.map(e => e) || [];
    return message;
  }

};