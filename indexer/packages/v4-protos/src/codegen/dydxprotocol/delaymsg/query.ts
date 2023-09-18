import { DelayedMessage, DelayedMessageSDKType } from "./delayed_message";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial, Long } from "../../helpers";
/**
 * QueryNumMessagesRequest is the request type for the NumMessages RPC
 * method.
 */

export interface QueryNumMessagesRequest {}
/**
 * QueryNumMessagesRequest is the request type for the NumMessages RPC
 * method.
 */

export interface QueryNumMessagesRequestSDKType {}
/**
 * QueryGetNumMessagesResponse is the response type for the NumMessages RPC
 * method.
 */

export interface QueryNumMessagesResponse {
  /**
   * QueryGetNumMessagesResponse is the response type for the NumMessages RPC
   * method.
   */
  numMessages: number;
}
/**
 * QueryGetNumMessagesResponse is the response type for the NumMessages RPC
 * method.
 */

export interface QueryNumMessagesResponseSDKType {
  /**
   * QueryGetNumMessagesResponse is the response type for the NumMessages RPC
   * method.
   */
  num_messages: number;
}
/** QueryMessageRequest is the request type for the Message RPC method. */

export interface QueryMessageRequest {
  /** QueryMessageRequest is the request type for the Message RPC method. */
  id: number;
}
/** QueryMessageRequest is the request type for the Message RPC method. */

export interface QueryMessageRequestSDKType {
  /** QueryMessageRequest is the request type for the Message RPC method. */
  id: number;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */

export interface QueryMessageResponse {
  /** QueryGetMessageResponse is the response type for the Message RPC method. */
  message?: DelayedMessage;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */

export interface QueryMessageResponseSDKType {
  /** QueryGetMessageResponse is the response type for the Message RPC method. */
  message?: DelayedMessageSDKType;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */

export interface QueryBlockMessageIdsRequest {
  /**
   * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
   * RPC method.
   */
  blockHeight: Long;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */

export interface QueryBlockMessageIdsRequestSDKType {
  /**
   * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
   * RPC method.
   */
  block_height: Long;
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */

export interface QueryBlockMessageIdsResponse {
  /**
   * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
   * RPC method.
   */
  messageIds: number[];
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */

export interface QueryBlockMessageIdsResponseSDKType {
  /**
   * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
   * RPC method.
   */
  message_ids: number[];
}
/** QueryAllMessagesRequest is the request type for the AllMessages RPC method. */

export interface QueryAllMessagesRequest {}
/** QueryAllMessagesRequest is the request type for the AllMessages RPC method. */

export interface QueryAllMessagesRequestSDKType {}
/** QueryAllMessagesResponse is the response type for the AllMessages RPC method. */

export interface QueryAllMessagesResponse {
  /** QueryAllMessagesResponse is the response type for the AllMessages RPC method. */
  messages: DelayedMessage[];
}
/** QueryAllMessagesResponse is the response type for the AllMessages RPC method. */

export interface QueryAllMessagesResponseSDKType {
  /** QueryAllMessagesResponse is the response type for the AllMessages RPC method. */
  messages: DelayedMessageSDKType[];
}

function createBaseQueryNumMessagesRequest(): QueryNumMessagesRequest {
  return {};
}

export const QueryNumMessagesRequest = {
  encode(_: QueryNumMessagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNumMessagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNumMessagesRequest();

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

  fromPartial(_: DeepPartial<QueryNumMessagesRequest>): QueryNumMessagesRequest {
    const message = createBaseQueryNumMessagesRequest();
    return message;
  }

};

function createBaseQueryNumMessagesResponse(): QueryNumMessagesResponse {
  return {
    numMessages: 0
  };
}

export const QueryNumMessagesResponse = {
  encode(message: QueryNumMessagesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.numMessages !== 0) {
      writer.uint32(8).uint32(message.numMessages);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNumMessagesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNumMessagesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.numMessages = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryNumMessagesResponse>): QueryNumMessagesResponse {
    const message = createBaseQueryNumMessagesResponse();
    message.numMessages = object.numMessages ?? 0;
    return message;
  }

};

function createBaseQueryMessageRequest(): QueryMessageRequest {
  return {
    id: 0
  };
}

export const QueryMessageRequest = {
  encode(message: QueryMessageRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMessageRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMessageRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMessageRequest>): QueryMessageRequest {
    const message = createBaseQueryMessageRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryMessageResponse(): QueryMessageResponse {
  return {
    message: undefined
  };
}

export const QueryMessageResponse = {
  encode(message: QueryMessageResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.message !== undefined) {
      DelayedMessage.encode(message.message, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMessageResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMessageResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.message = DelayedMessage.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMessageResponse>): QueryMessageResponse {
    const message = createBaseQueryMessageResponse();
    message.message = object.message !== undefined && object.message !== null ? DelayedMessage.fromPartial(object.message) : undefined;
    return message;
  }

};

function createBaseQueryBlockMessageIdsRequest(): QueryBlockMessageIdsRequest {
  return {
    blockHeight: Long.ZERO
  };
}

export const QueryBlockMessageIdsRequest = {
  encode(message: QueryBlockMessageIdsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (!message.blockHeight.isZero()) {
      writer.uint32(8).sint64(message.blockHeight);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockMessageIdsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlockMessageIdsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.blockHeight = (reader.sint64() as Long);
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryBlockMessageIdsRequest>): QueryBlockMessageIdsRequest {
    const message = createBaseQueryBlockMessageIdsRequest();
    message.blockHeight = object.blockHeight !== undefined && object.blockHeight !== null ? Long.fromValue(object.blockHeight) : Long.ZERO;
    return message;
  }

};

function createBaseQueryBlockMessageIdsResponse(): QueryBlockMessageIdsResponse {
  return {
    messageIds: []
  };
}

export const QueryBlockMessageIdsResponse = {
  encode(message: QueryBlockMessageIdsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();

    for (const v of message.messageIds) {
      writer.uint32(v);
    }

    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockMessageIdsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryBlockMessageIdsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;

            while (reader.pos < end2) {
              message.messageIds.push(reader.uint32());
            }
          } else {
            message.messageIds.push(reader.uint32());
          }

          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryBlockMessageIdsResponse>): QueryBlockMessageIdsResponse {
    const message = createBaseQueryBlockMessageIdsResponse();
    message.messageIds = object.messageIds?.map(e => e) || [];
    return message;
  }

};

function createBaseQueryAllMessagesRequest(): QueryAllMessagesRequest {
  return {};
}

export const QueryAllMessagesRequest = {
  encode(_: QueryAllMessagesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMessagesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMessagesRequest();

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

  fromPartial(_: DeepPartial<QueryAllMessagesRequest>): QueryAllMessagesRequest {
    const message = createBaseQueryAllMessagesRequest();
    return message;
  }

};

function createBaseQueryAllMessagesResponse(): QueryAllMessagesResponse {
  return {
    messages: []
  };
}

export const QueryAllMessagesResponse = {
  encode(message: QueryAllMessagesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.messages) {
      DelayedMessage.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllMessagesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllMessagesResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.messages.push(DelayedMessage.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllMessagesResponse>): QueryAllMessagesResponse {
    const message = createBaseQueryAllMessagesResponse();
    message.messages = object.messages?.map(e => DelayedMessage.fromPartial(e)) || [];
    return message;
  }

};