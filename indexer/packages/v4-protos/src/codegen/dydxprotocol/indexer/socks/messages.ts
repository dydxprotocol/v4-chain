import { IndexerSubaccountId, IndexerSubaccountIdSDKType } from "../protocol/v1/subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** TODO(IND-210): Make this proto conform and update downstream indexer logic */

export enum CandleMessage_Resolution {
  /**
   * ONE_MINUTE - buf:lint:ignore ENUM_VALUE_PREFIX
   * buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
   */
  ONE_MINUTE = 0,

  /** FIVE_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  FIVE_MINUTES = 1,

  /** FIFTEEN_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  FIFTEEN_MINUTES = 2,

  /** THIRTY_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  THIRTY_MINUTES = 3,

  /** ONE_HOUR - buf:lint:ignore ENUM_VALUE_PREFIX */
  ONE_HOUR = 4,

  /** FOUR_HOURS - buf:lint:ignore ENUM_VALUE_PREFIX */
  FOUR_HOURS = 5,

  /** ONE_DAY - buf:lint:ignore ENUM_VALUE_PREFIX */
  ONE_DAY = 6,
  UNRECOGNIZED = -1,
}
/** TODO(IND-210): Make this proto conform and update downstream indexer logic */

export enum CandleMessage_ResolutionSDKType {
  /**
   * ONE_MINUTE - buf:lint:ignore ENUM_VALUE_PREFIX
   * buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
   */
  ONE_MINUTE = 0,

  /** FIVE_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  FIVE_MINUTES = 1,

  /** FIFTEEN_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  FIFTEEN_MINUTES = 2,

  /** THIRTY_MINUTES - buf:lint:ignore ENUM_VALUE_PREFIX */
  THIRTY_MINUTES = 3,

  /** ONE_HOUR - buf:lint:ignore ENUM_VALUE_PREFIX */
  ONE_HOUR = 4,

  /** FOUR_HOURS - buf:lint:ignore ENUM_VALUE_PREFIX */
  FOUR_HOURS = 5,

  /** ONE_DAY - buf:lint:ignore ENUM_VALUE_PREFIX */
  ONE_DAY = 6,
  UNRECOGNIZED = -1,
}
export function candleMessage_ResolutionFromJSON(object: any): CandleMessage_Resolution {
  switch (object) {
    case 0:
    case "ONE_MINUTE":
      return CandleMessage_Resolution.ONE_MINUTE;

    case 1:
    case "FIVE_MINUTES":
      return CandleMessage_Resolution.FIVE_MINUTES;

    case 2:
    case "FIFTEEN_MINUTES":
      return CandleMessage_Resolution.FIFTEEN_MINUTES;

    case 3:
    case "THIRTY_MINUTES":
      return CandleMessage_Resolution.THIRTY_MINUTES;

    case 4:
    case "ONE_HOUR":
      return CandleMessage_Resolution.ONE_HOUR;

    case 5:
    case "FOUR_HOURS":
      return CandleMessage_Resolution.FOUR_HOURS;

    case 6:
    case "ONE_DAY":
      return CandleMessage_Resolution.ONE_DAY;

    case -1:
    case "UNRECOGNIZED":
    default:
      return CandleMessage_Resolution.UNRECOGNIZED;
  }
}
export function candleMessage_ResolutionToJSON(object: CandleMessage_Resolution): string {
  switch (object) {
    case CandleMessage_Resolution.ONE_MINUTE:
      return "ONE_MINUTE";

    case CandleMessage_Resolution.FIVE_MINUTES:
      return "FIVE_MINUTES";

    case CandleMessage_Resolution.FIFTEEN_MINUTES:
      return "FIFTEEN_MINUTES";

    case CandleMessage_Resolution.THIRTY_MINUTES:
      return "THIRTY_MINUTES";

    case CandleMessage_Resolution.ONE_HOUR:
      return "ONE_HOUR";

    case CandleMessage_Resolution.FOUR_HOURS:
      return "FOUR_HOURS";

    case CandleMessage_Resolution.ONE_DAY:
      return "ONE_DAY";

    case CandleMessage_Resolution.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
/** Message to be sent through the 'to-websockets-orderbooks` kafka topic. */

export interface OrderbookMessage {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Clob pair id of the Orderbook message. */

  clobPairId: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-orderbooks` kafka topic. */

export interface OrderbookMessageSDKType {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Clob pair id of the Orderbook message. */

  clob_pair_id: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-subaccounts` kafka topic. */

export interface SubaccountMessage {
  /** Block height where the contents occur. */
  blockHeight: string;
  /** Transaction index where the contents occur. */

  transactionIndex: number;
  /** Event index where the contents occur. */

  eventIndex: number;
  /** Stringified JSON object of all events to be streamed. */

  contents: string;
  /** Subaccount id that the content corresponds to. */

  subaccountId?: IndexerSubaccountId;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-subaccounts` kafka topic. */

export interface SubaccountMessageSDKType {
  /** Block height where the contents occur. */
  block_height: string;
  /** Transaction index where the contents occur. */

  transaction_index: number;
  /** Event index where the contents occur. */

  event_index: number;
  /** Stringified JSON object of all events to be streamed. */

  contents: string;
  /** Subaccount id that the content corresponds to. */

  subaccount_id?: IndexerSubaccountIdSDKType;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-trades` kafka topic. */

export interface TradeMessage {
  /** Block height where the contents occur. */
  blockHeight: string;
  /** Stringified JSON object of all events to be streamed. */

  contents: string;
  /** Clob pair id of the Trade message. */

  clobPairId: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-trades` kafka topic. */

export interface TradeMessageSDKType {
  /** Block height where the contents occur. */
  block_height: string;
  /** Stringified JSON object of all events to be streamed. */

  contents: string;
  /** Clob pair id of the Trade message. */

  clob_pair_id: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-markets` kafka topic. */

export interface MarketMessage {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-markets` kafka topic. */

export interface MarketMessageSDKType {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-candles` kafka topic. */

export interface CandleMessage {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Clob pair id of the Candle message. */

  clobPairId: string;
  /** Resolution of the candle update. */

  resolution: CandleMessage_Resolution;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-candles` kafka topic. */

export interface CandleMessageSDKType {
  /** Stringified JSON object of all events to be streamed. */
  contents: string;
  /** Clob pair id of the Candle message. */

  clob_pair_id: string;
  /** Resolution of the candle update. */

  resolution: CandleMessage_ResolutionSDKType;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-block-height` kafka topic. */

export interface BlockHeightMessage {
  /** Block height where the contents occur. */
  blockHeight: string;
  /** ISO formatted time of the block height. */

  time: string;
  /** Version of the websocket message. */

  version: string;
}
/** Message to be sent through the 'to-websockets-block-height` kafka topic. */

export interface BlockHeightMessageSDKType {
  /** Block height where the contents occur. */
  block_height: string;
  /** ISO formatted time of the block height. */

  time: string;
  /** Version of the websocket message. */

  version: string;
}

function createBaseOrderbookMessage(): OrderbookMessage {
  return {
    contents: "",
    clobPairId: "",
    version: ""
  };
}

export const OrderbookMessage = {
  encode(message: OrderbookMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contents !== "") {
      writer.uint32(10).string(message.contents);
    }

    if (message.clobPairId !== "") {
      writer.uint32(18).string(message.clobPairId);
    }

    if (message.version !== "") {
      writer.uint32(26).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OrderbookMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOrderbookMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.contents = reader.string();
          break;

        case 2:
          message.clobPairId = reader.string();
          break;

        case 3:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<OrderbookMessage>): OrderbookMessage {
    const message = createBaseOrderbookMessage();
    message.contents = object.contents ?? "";
    message.clobPairId = object.clobPairId ?? "";
    message.version = object.version ?? "";
    return message;
  }

};

function createBaseSubaccountMessage(): SubaccountMessage {
  return {
    blockHeight: "",
    transactionIndex: 0,
    eventIndex: 0,
    contents: "",
    subaccountId: undefined,
    version: ""
  };
}

export const SubaccountMessage = {
  encode(message: SubaccountMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== "") {
      writer.uint32(10).string(message.blockHeight);
    }

    if (message.transactionIndex !== 0) {
      writer.uint32(16).int32(message.transactionIndex);
    }

    if (message.eventIndex !== 0) {
      writer.uint32(24).uint32(message.eventIndex);
    }

    if (message.contents !== "") {
      writer.uint32(34).string(message.contents);
    }

    if (message.subaccountId !== undefined) {
      IndexerSubaccountId.encode(message.subaccountId, writer.uint32(42).fork()).ldelim();
    }

    if (message.version !== "") {
      writer.uint32(50).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SubaccountMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSubaccountMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = reader.string();
          break;

        case 2:
          message.transactionIndex = reader.int32();
          break;

        case 3:
          message.eventIndex = reader.uint32();
          break;

        case 4:
          message.contents = reader.string();
          break;

        case 5:
          message.subaccountId = IndexerSubaccountId.decode(reader, reader.uint32());
          break;

        case 6:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<SubaccountMessage>): SubaccountMessage {
    const message = createBaseSubaccountMessage();
    message.blockHeight = object.blockHeight ?? "";
    message.transactionIndex = object.transactionIndex ?? 0;
    message.eventIndex = object.eventIndex ?? 0;
    message.contents = object.contents ?? "";
    message.subaccountId = object.subaccountId !== undefined && object.subaccountId !== null ? IndexerSubaccountId.fromPartial(object.subaccountId) : undefined;
    message.version = object.version ?? "";
    return message;
  }

};

function createBaseTradeMessage(): TradeMessage {
  return {
    blockHeight: "",
    contents: "",
    clobPairId: "",
    version: ""
  };
}

export const TradeMessage = {
  encode(message: TradeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== "") {
      writer.uint32(10).string(message.blockHeight);
    }

    if (message.contents !== "") {
      writer.uint32(34).string(message.contents);
    }

    if (message.clobPairId !== "") {
      writer.uint32(42).string(message.clobPairId);
    }

    if (message.version !== "") {
      writer.uint32(50).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TradeMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTradeMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = reader.string();
          break;

        case 4:
          message.contents = reader.string();
          break;

        case 5:
          message.clobPairId = reader.string();
          break;

        case 6:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<TradeMessage>): TradeMessage {
    const message = createBaseTradeMessage();
    message.blockHeight = object.blockHeight ?? "";
    message.contents = object.contents ?? "";
    message.clobPairId = object.clobPairId ?? "";
    message.version = object.version ?? "";
    return message;
  }

};

function createBaseMarketMessage(): MarketMessage {
  return {
    contents: "",
    version: ""
  };
}

export const MarketMessage = {
  encode(message: MarketMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contents !== "") {
      writer.uint32(10).string(message.contents);
    }

    if (message.version !== "") {
      writer.uint32(18).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MarketMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMarketMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.contents = reader.string();
          break;

        case 2:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MarketMessage>): MarketMessage {
    const message = createBaseMarketMessage();
    message.contents = object.contents ?? "";
    message.version = object.version ?? "";
    return message;
  }

};

function createBaseCandleMessage(): CandleMessage {
  return {
    contents: "",
    clobPairId: "",
    resolution: 0,
    version: ""
  };
}

export const CandleMessage = {
  encode(message: CandleMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contents !== "") {
      writer.uint32(10).string(message.contents);
    }

    if (message.clobPairId !== "") {
      writer.uint32(18).string(message.clobPairId);
    }

    if (message.resolution !== 0) {
      writer.uint32(24).int32(message.resolution);
    }

    if (message.version !== "") {
      writer.uint32(34).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CandleMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCandleMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.contents = reader.string();
          break;

        case 2:
          message.clobPairId = reader.string();
          break;

        case 3:
          message.resolution = (reader.int32() as any);
          break;

        case 4:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<CandleMessage>): CandleMessage {
    const message = createBaseCandleMessage();
    message.contents = object.contents ?? "";
    message.clobPairId = object.clobPairId ?? "";
    message.resolution = object.resolution ?? 0;
    message.version = object.version ?? "";
    return message;
  }

};

function createBaseBlockHeightMessage(): BlockHeightMessage {
  return {
    blockHeight: "",
    time: "",
    version: ""
  };
}

export const BlockHeightMessage = {
  encode(message: BlockHeightMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== "") {
      writer.uint32(10).string(message.blockHeight);
    }

    if (message.time !== "") {
      writer.uint32(18).string(message.time);
    }

    if (message.version !== "") {
      writer.uint32(26).string(message.version);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BlockHeightMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBlockHeightMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = reader.string();
          break;

        case 2:
          message.time = reader.string();
          break;

        case 3:
          message.version = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<BlockHeightMessage>): BlockHeightMessage {
    const message = createBaseBlockHeightMessage();
    message.blockHeight = object.blockHeight ?? "";
    message.time = object.time ?? "";
    message.version = object.version ?? "";
    return message;
  }

};