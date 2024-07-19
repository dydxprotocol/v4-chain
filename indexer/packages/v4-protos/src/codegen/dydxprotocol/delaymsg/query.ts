import { DelayedMessage, DelayedMessageSDKType } from "./delayed_message";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */

export interface QueryNextDelayedMessageIdRequest {}
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */

export interface QueryNextDelayedMessageIdRequestSDKType {}
/**
 * QueryNextDelayedMessageIdResponse is the response type for the
 * NextDelayedMessageId RPC method.
 */

export interface QueryNextDelayedMessageIdResponse {
  nextDelayedMessageId: number;
}
/**
 * QueryNextDelayedMessageIdResponse is the response type for the
 * NextDelayedMessageId RPC method.
 */

export interface QueryNextDelayedMessageIdResponseSDKType {
  next_delayed_message_id: number;
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
  blockHeight: number;
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
  block_height: number;
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

function createBaseQueryNextDelayedMessageIdRequest(): QueryNextDelayedMessageIdRequest {
  return {};
}

export const QueryNextDelayedMessageIdRequest = {
  encode(_: QueryNextDelayedMessageIdRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextDelayedMessageIdRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextDelayedMessageIdRequest();

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

  fromPartial(_: DeepPartial<QueryNextDelayedMessageIdRequest>): QueryNextDelayedMessageIdRequest {
    const message = createBaseQueryNextDelayedMessageIdRequest();
    return message;
  }

};

function createBaseQueryNextDelayedMessageIdResponse(): QueryNextDelayedMessageIdResponse {
  return {
    nextDelayedMessageId: 0
  };
}

export const QueryNextDelayedMessageIdResponse = {
  encode(message: QueryNextDelayedMessageIdResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nextDelayedMessageId !== 0) {
      writer.uint32(8).uint32(message.nextDelayedMessageId);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextDelayedMessageIdResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryNextDelayedMessageIdResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.nextDelayedMessageId = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryNextDelayedMessageIdResponse>): QueryNextDelayedMessageIdResponse {
    const message = createBaseQueryNextDelayedMessageIdResponse();
    message.nextDelayedMessageId = object.nextDelayedMessageId ?? 0;
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
    blockHeight: 0
  };
}

export const QueryBlockMessageIdsRequest = {
  encode(message: QueryBlockMessageIdsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.blockHeight !== 0) {
      writer.uint32(8).uint32(message.blockHeight);
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
          message.blockHeight = reader.uint32();
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
    message.blockHeight = object.blockHeight ?? 0;
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