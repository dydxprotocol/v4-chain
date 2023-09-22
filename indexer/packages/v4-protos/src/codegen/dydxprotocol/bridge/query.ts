import { EventParams, EventParamsSDKType, ProposeParams, ProposeParamsSDKType, SafetyParams, SafetyParamsSDKType } from "./params";
import { BridgeEventInfo, BridgeEventInfoSDKType } from "./bridge_event_info";
import { MsgCompleteBridge, MsgCompleteBridgeSDKType } from "./tx";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */

export interface QueryEventParamsRequest {}
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */

export interface QueryEventParamsRequestSDKType {}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */

export interface QueryEventParamsResponse {
  params?: EventParams;
}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */

export interface QueryEventParamsResponseSDKType {
  params?: EventParamsSDKType;
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */

export interface QueryProposeParamsRequest {}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */

export interface QueryProposeParamsRequestSDKType {}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */

export interface QueryProposeParamsResponse {
  params?: ProposeParams;
}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */

export interface QueryProposeParamsResponseSDKType {
  params?: ProposeParamsSDKType;
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */

export interface QuerySafetyParamsRequest {}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */

export interface QuerySafetyParamsRequestSDKType {}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */

export interface QuerySafetyParamsResponse {
  params?: SafetyParams;
}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */

export interface QuerySafetyParamsResponseSDKType {
  params?: SafetyParamsSDKType;
}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */

export interface QueryAcknowledgedEventInfoRequest {}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */

export interface QueryAcknowledgedEventInfoRequestSDKType {}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */

export interface QueryAcknowledgedEventInfoResponse {
  info?: BridgeEventInfo;
}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */

export interface QueryAcknowledgedEventInfoResponseSDKType {
  info?: BridgeEventInfoSDKType;
}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */

export interface QueryRecognizedEventInfoRequest {}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */

export interface QueryRecognizedEventInfoRequestSDKType {}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */

export interface QueryRecognizedEventInfoResponse {
  info?: BridgeEventInfo;
}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */

export interface QueryRecognizedEventInfoResponseSDKType {
  info?: BridgeEventInfoSDKType;
}
/**
 * QueryInFlightCompleteBridgeMessagesRequest is a request type for the
 * InFlightCompleteBridgeMessages RPC method.
 */

export interface QueryInFlightCompleteBridgeMessagesRequest {
  /**
   * QueryInFlightCompleteBridgeMessagesRequest is a request type for the
   * InFlightCompleteBridgeMessages RPC method.
   */
  address: string;
}
/**
 * QueryInFlightCompleteBridgeMessagesRequest is a request type for the
 * InFlightCompleteBridgeMessages RPC method.
 */

export interface QueryInFlightCompleteBridgeMessagesRequestSDKType {
  /**
   * QueryInFlightCompleteBridgeMessagesRequest is a request type for the
   * InFlightCompleteBridgeMessages RPC method.
   */
  address: string;
}
/**
 * QueryInFlightCompleteBridgeMessagesResponse is a response type for the
 * InFlightCompleteBridgeMessages RPC method.
 */

export interface QueryInFlightCompleteBridgeMessagesResponse {
  messages: InFlightCompleteBridgeMessage[];
}
/**
 * QueryInFlightCompleteBridgeMessagesResponse is a response type for the
 * InFlightCompleteBridgeMessages RPC method.
 */

export interface QueryInFlightCompleteBridgeMessagesResponseSDKType {
  messages: InFlightCompleteBridgeMessageSDKType[];
}
/**
 * InFlightCompleteBridgeMessage is a message type for the response of
 * InFlightCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */

export interface InFlightCompleteBridgeMessage {
  message?: MsgCompleteBridge;
  blockHeight: Long;
}
/**
 * InFlightCompleteBridgeMessage is a message type for the response of
 * InFlightCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */

export interface InFlightCompleteBridgeMessageSDKType {
  message?: MsgCompleteBridgeSDKType;
  block_height: Long;
}

function createBaseQueryEventParamsRequest(): QueryEventParamsRequest {
  return {};
}

export const QueryEventParamsRequest = {
  encode(_: QueryEventParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEventParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEventParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryEventParamsRequest>): QueryEventParamsRequest {
    const message = createBaseQueryEventParamsRequest();
    return message;
  }

};

function createBaseQueryEventParamsResponse(): QueryEventParamsResponse {
  return {
    params: undefined
  };
}

export const QueryEventParamsResponse = {
  encode(message: QueryEventParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      EventParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryEventParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryEventParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = EventParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryEventParamsResponse>): QueryEventParamsResponse {
    const message = createBaseQueryEventParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? EventParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryProposeParamsRequest(): QueryProposeParamsRequest {
  return {};
}

export const QueryProposeParamsRequest = {
  encode(_: QueryProposeParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProposeParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProposeParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryProposeParamsRequest>): QueryProposeParamsRequest {
    const message = createBaseQueryProposeParamsRequest();
    return message;
  }

};

function createBaseQueryProposeParamsResponse(): QueryProposeParamsResponse {
  return {
    params: undefined
  };
}

export const QueryProposeParamsResponse = {
  encode(message: QueryProposeParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      ProposeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryProposeParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryProposeParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = ProposeParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryProposeParamsResponse>): QueryProposeParamsResponse {
    const message = createBaseQueryProposeParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? ProposeParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQuerySafetyParamsRequest(): QuerySafetyParamsRequest {
  return {};
}

export const QuerySafetyParamsRequest = {
  encode(_: QuerySafetyParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySafetyParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySafetyParamsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QuerySafetyParamsRequest>): QuerySafetyParamsRequest {
    const message = createBaseQuerySafetyParamsRequest();
    return message;
  }

};

function createBaseQuerySafetyParamsResponse(): QuerySafetyParamsResponse {
  return {
    params: undefined
  };
}

export const QuerySafetyParamsResponse = {
  encode(message: QuerySafetyParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      SafetyParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QuerySafetyParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQuerySafetyParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = SafetyParams.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QuerySafetyParamsResponse>): QuerySafetyParamsResponse {
    const message = createBaseQuerySafetyParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? SafetyParams.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryAcknowledgedEventInfoRequest(): QueryAcknowledgedEventInfoRequest {
  return {};
}

export const QueryAcknowledgedEventInfoRequest = {
  encode(_: QueryAcknowledgedEventInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAcknowledgedEventInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAcknowledgedEventInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryAcknowledgedEventInfoRequest>): QueryAcknowledgedEventInfoRequest {
    const message = createBaseQueryAcknowledgedEventInfoRequest();
    return message;
  }

};

function createBaseQueryAcknowledgedEventInfoResponse(): QueryAcknowledgedEventInfoResponse {
  return {
    info: undefined
  };
}

export const QueryAcknowledgedEventInfoResponse = {
  encode(message: QueryAcknowledgedEventInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      BridgeEventInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAcknowledgedEventInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAcknowledgedEventInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.info = BridgeEventInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAcknowledgedEventInfoResponse>): QueryAcknowledgedEventInfoResponse {
    const message = createBaseQueryAcknowledgedEventInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? BridgeEventInfo.fromPartial(object.info) : undefined;
    return message;
  }

};

function createBaseQueryRecognizedEventInfoRequest(): QueryRecognizedEventInfoRequest {
  return {};
}

export const QueryRecognizedEventInfoRequest = {
  encode(_: QueryRecognizedEventInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecognizedEventInfoRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRecognizedEventInfoRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(_: DeepPartial<QueryRecognizedEventInfoRequest>): QueryRecognizedEventInfoRequest {
    const message = createBaseQueryRecognizedEventInfoRequest();
    return message;
  }

};

function createBaseQueryRecognizedEventInfoResponse(): QueryRecognizedEventInfoResponse {
  return {
    info: undefined
  };
}

export const QueryRecognizedEventInfoResponse = {
  encode(message: QueryRecognizedEventInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.info !== undefined) {
      BridgeEventInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecognizedEventInfoResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRecognizedEventInfoResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.info = BridgeEventInfo.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryRecognizedEventInfoResponse>): QueryRecognizedEventInfoResponse {
    const message = createBaseQueryRecognizedEventInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? BridgeEventInfo.fromPartial(object.info) : undefined;
    return message;
  }

};

function createBaseQueryInFlightCompleteBridgeMessagesRequest(): QueryInFlightCompleteBridgeMessagesRequest {
  return {
    address: ""
  };
}

export const QueryInFlightCompleteBridgeMessagesRequest = {
  encode(message: QueryInFlightCompleteBridgeMessagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryInFlightCompleteBridgeMessagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryInFlightCompleteBridgeMessagesRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryInFlightCompleteBridgeMessagesRequest>): QueryInFlightCompleteBridgeMessagesRequest {
    const message = createBaseQueryInFlightCompleteBridgeMessagesRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseQueryInFlightCompleteBridgeMessagesResponse(): QueryInFlightCompleteBridgeMessagesResponse {
  return {
    messages: []
  };
}

export const QueryInFlightCompleteBridgeMessagesResponse = {
  encode(message: QueryInFlightCompleteBridgeMessagesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.messages) {
      InFlightCompleteBridgeMessage.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryInFlightCompleteBridgeMessagesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryInFlightCompleteBridgeMessagesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.messages.push(InFlightCompleteBridgeMessage.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryInFlightCompleteBridgeMessagesResponse>): QueryInFlightCompleteBridgeMessagesResponse {
    const message = createBaseQueryInFlightCompleteBridgeMessagesResponse();
    message.messages = object.messages?.map(e => InFlightCompleteBridgeMessage.fromPartial(e)) || [];
    return message;
  }

};

function createBaseInFlightCompleteBridgeMessage(): InFlightCompleteBridgeMessage {
  return {
    message: undefined,
    blockHeight: Long.ZERO
  };
}

export const InFlightCompleteBridgeMessage = {
  encode(message: InFlightCompleteBridgeMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.message !== undefined) {
      MsgCompleteBridge.encode(message.message, writer.uint32(10).fork()).ldelim();
    }

    if (!message.blockHeight.isZero()) {
      writer.uint32(16).int64(message.blockHeight);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): InFlightCompleteBridgeMessage {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInFlightCompleteBridgeMessage();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.message = MsgCompleteBridge.decode(reader, reader.uint32());
          break;

        case 2:
          message.blockHeight = (reader.int64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<InFlightCompleteBridgeMessage>): InFlightCompleteBridgeMessage {
    const message = createBaseInFlightCompleteBridgeMessage();
    message.message = object.message !== undefined && object.message !== null ? MsgCompleteBridge.fromPartial(object.message) : undefined;
    message.blockHeight = object.blockHeight !== undefined && object.blockHeight !== null ? Long.fromValue(object.blockHeight) : Long.ZERO;
    return message;
  }

};