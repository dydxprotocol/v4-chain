import { DelayedMessage, DelayedMessageAmino, DelayedMessageSDKType } from "./delayed_message";
import { BinaryReader, BinaryWriter } from "../../binary";
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdRequest {}
export interface QueryNextDelayedMessageIdRequestProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdRequest";
  value: Uint8Array;
}
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdRequestAmino {}
export interface QueryNextDelayedMessageIdRequestAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdRequest";
  value: QueryNextDelayedMessageIdRequestAmino;
}
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
export interface QueryNextDelayedMessageIdResponseProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdResponse";
  value: Uint8Array;
}
/**
 * QueryNextDelayedMessageIdResponse is the response type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdResponseAmino {
  next_delayed_message_id?: number;
}
export interface QueryNextDelayedMessageIdResponseAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdResponse";
  value: QueryNextDelayedMessageIdResponseAmino;
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
export interface QueryMessageRequestProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryMessageRequest";
  value: Uint8Array;
}
/** QueryMessageRequest is the request type for the Message RPC method. */
export interface QueryMessageRequestAmino {
  /** QueryMessageRequest is the request type for the Message RPC method. */
  id?: number;
}
export interface QueryMessageRequestAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryMessageRequest";
  value: QueryMessageRequestAmino;
}
/** QueryMessageRequest is the request type for the Message RPC method. */
export interface QueryMessageRequestSDKType {
  id: number;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */
export interface QueryMessageResponse {
  /** QueryGetMessageResponse is the response type for the Message RPC method. */
  message?: DelayedMessage;
}
export interface QueryMessageResponseProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryMessageResponse";
  value: Uint8Array;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */
export interface QueryMessageResponseAmino {
  /** QueryGetMessageResponse is the response type for the Message RPC method. */
  message?: DelayedMessageAmino;
}
export interface QueryMessageResponseAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryMessageResponse";
  value: QueryMessageResponseAmino;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */
export interface QueryMessageResponseSDKType {
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
export interface QueryBlockMessageIdsRequestProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsRequest";
  value: Uint8Array;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsRequestAmino {
  /**
   * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
   * RPC method.
   */
  block_height?: number;
}
export interface QueryBlockMessageIdsRequestAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryBlockMessageIdsRequest";
  value: QueryBlockMessageIdsRequestAmino;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsRequestSDKType {
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
export interface QueryBlockMessageIdsResponseProtoMsg {
  typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsResponse";
  value: Uint8Array;
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsResponseAmino {
  /**
   * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
   * RPC method.
   */
  message_ids?: number[];
}
export interface QueryBlockMessageIdsResponseAminoMsg {
  type: "/dydxprotocol.delaymsg.QueryBlockMessageIdsResponse";
  value: QueryBlockMessageIdsResponseAmino;
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsResponseSDKType {
  message_ids: number[];
}
function createBaseQueryNextDelayedMessageIdRequest(): QueryNextDelayedMessageIdRequest {
  return {};
}
export const QueryNextDelayedMessageIdRequest = {
  typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdRequest",
  encode(_: QueryNextDelayedMessageIdRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryNextDelayedMessageIdRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QueryNextDelayedMessageIdRequest>): QueryNextDelayedMessageIdRequest {
    const message = createBaseQueryNextDelayedMessageIdRequest();
    return message;
  },
  fromAmino(_: QueryNextDelayedMessageIdRequestAmino): QueryNextDelayedMessageIdRequest {
    const message = createBaseQueryNextDelayedMessageIdRequest();
    return message;
  },
  toAmino(_: QueryNextDelayedMessageIdRequest): QueryNextDelayedMessageIdRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryNextDelayedMessageIdRequestAminoMsg): QueryNextDelayedMessageIdRequest {
    return QueryNextDelayedMessageIdRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryNextDelayedMessageIdRequestProtoMsg): QueryNextDelayedMessageIdRequest {
    return QueryNextDelayedMessageIdRequest.decode(message.value);
  },
  toProto(message: QueryNextDelayedMessageIdRequest): Uint8Array {
    return QueryNextDelayedMessageIdRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryNextDelayedMessageIdRequest): QueryNextDelayedMessageIdRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdRequest",
      value: QueryNextDelayedMessageIdRequest.encode(message).finish()
    };
  }
};
function createBaseQueryNextDelayedMessageIdResponse(): QueryNextDelayedMessageIdResponse {
  return {
    nextDelayedMessageId: 0
  };
}
export const QueryNextDelayedMessageIdResponse = {
  typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdResponse",
  encode(message: QueryNextDelayedMessageIdResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.nextDelayedMessageId !== 0) {
      writer.uint32(8).uint32(message.nextDelayedMessageId);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryNextDelayedMessageIdResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryNextDelayedMessageIdResponse>): QueryNextDelayedMessageIdResponse {
    const message = createBaseQueryNextDelayedMessageIdResponse();
    message.nextDelayedMessageId = object.nextDelayedMessageId ?? 0;
    return message;
  },
  fromAmino(object: QueryNextDelayedMessageIdResponseAmino): QueryNextDelayedMessageIdResponse {
    const message = createBaseQueryNextDelayedMessageIdResponse();
    if (object.next_delayed_message_id !== undefined && object.next_delayed_message_id !== null) {
      message.nextDelayedMessageId = object.next_delayed_message_id;
    }
    return message;
  },
  toAmino(message: QueryNextDelayedMessageIdResponse): QueryNextDelayedMessageIdResponseAmino {
    const obj: any = {};
    obj.next_delayed_message_id = message.nextDelayedMessageId;
    return obj;
  },
  fromAminoMsg(object: QueryNextDelayedMessageIdResponseAminoMsg): QueryNextDelayedMessageIdResponse {
    return QueryNextDelayedMessageIdResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryNextDelayedMessageIdResponseProtoMsg): QueryNextDelayedMessageIdResponse {
    return QueryNextDelayedMessageIdResponse.decode(message.value);
  },
  toProto(message: QueryNextDelayedMessageIdResponse): Uint8Array {
    return QueryNextDelayedMessageIdResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryNextDelayedMessageIdResponse): QueryNextDelayedMessageIdResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryNextDelayedMessageIdResponse",
      value: QueryNextDelayedMessageIdResponse.encode(message).finish()
    };
  }
};
function createBaseQueryMessageRequest(): QueryMessageRequest {
  return {
    id: 0
  };
}
export const QueryMessageRequest = {
  typeUrl: "/dydxprotocol.delaymsg.QueryMessageRequest",
  encode(message: QueryMessageRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMessageRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMessageRequest>): QueryMessageRequest {
    const message = createBaseQueryMessageRequest();
    message.id = object.id ?? 0;
    return message;
  },
  fromAmino(object: QueryMessageRequestAmino): QueryMessageRequest {
    const message = createBaseQueryMessageRequest();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    return message;
  },
  toAmino(message: QueryMessageRequest): QueryMessageRequestAmino {
    const obj: any = {};
    obj.id = message.id;
    return obj;
  },
  fromAminoMsg(object: QueryMessageRequestAminoMsg): QueryMessageRequest {
    return QueryMessageRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMessageRequestProtoMsg): QueryMessageRequest {
    return QueryMessageRequest.decode(message.value);
  },
  toProto(message: QueryMessageRequest): Uint8Array {
    return QueryMessageRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryMessageRequest): QueryMessageRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryMessageRequest",
      value: QueryMessageRequest.encode(message).finish()
    };
  }
};
function createBaseQueryMessageResponse(): QueryMessageResponse {
  return {
    message: undefined
  };
}
export const QueryMessageResponse = {
  typeUrl: "/dydxprotocol.delaymsg.QueryMessageResponse",
  encode(message: QueryMessageResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.message !== undefined) {
      DelayedMessage.encode(message.message, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryMessageResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryMessageResponse>): QueryMessageResponse {
    const message = createBaseQueryMessageResponse();
    message.message = object.message !== undefined && object.message !== null ? DelayedMessage.fromPartial(object.message) : undefined;
    return message;
  },
  fromAmino(object: QueryMessageResponseAmino): QueryMessageResponse {
    const message = createBaseQueryMessageResponse();
    if (object.message !== undefined && object.message !== null) {
      message.message = DelayedMessage.fromAmino(object.message);
    }
    return message;
  },
  toAmino(message: QueryMessageResponse): QueryMessageResponseAmino {
    const obj: any = {};
    obj.message = message.message ? DelayedMessage.toAmino(message.message) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryMessageResponseAminoMsg): QueryMessageResponse {
    return QueryMessageResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryMessageResponseProtoMsg): QueryMessageResponse {
    return QueryMessageResponse.decode(message.value);
  },
  toProto(message: QueryMessageResponse): Uint8Array {
    return QueryMessageResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryMessageResponse): QueryMessageResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryMessageResponse",
      value: QueryMessageResponse.encode(message).finish()
    };
  }
};
function createBaseQueryBlockMessageIdsRequest(): QueryBlockMessageIdsRequest {
  return {
    blockHeight: 0
  };
}
export const QueryBlockMessageIdsRequest = {
  typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsRequest",
  encode(message: QueryBlockMessageIdsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.blockHeight !== 0) {
      writer.uint32(8).uint32(message.blockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlockMessageIdsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryBlockMessageIdsRequest>): QueryBlockMessageIdsRequest {
    const message = createBaseQueryBlockMessageIdsRequest();
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  },
  fromAmino(object: QueryBlockMessageIdsRequestAmino): QueryBlockMessageIdsRequest {
    const message = createBaseQueryBlockMessageIdsRequest();
    if (object.block_height !== undefined && object.block_height !== null) {
      message.blockHeight = object.block_height;
    }
    return message;
  },
  toAmino(message: QueryBlockMessageIdsRequest): QueryBlockMessageIdsRequestAmino {
    const obj: any = {};
    obj.block_height = message.blockHeight;
    return obj;
  },
  fromAminoMsg(object: QueryBlockMessageIdsRequestAminoMsg): QueryBlockMessageIdsRequest {
    return QueryBlockMessageIdsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlockMessageIdsRequestProtoMsg): QueryBlockMessageIdsRequest {
    return QueryBlockMessageIdsRequest.decode(message.value);
  },
  toProto(message: QueryBlockMessageIdsRequest): Uint8Array {
    return QueryBlockMessageIdsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryBlockMessageIdsRequest): QueryBlockMessageIdsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsRequest",
      value: QueryBlockMessageIdsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryBlockMessageIdsResponse(): QueryBlockMessageIdsResponse {
  return {
    messageIds: []
  };
}
export const QueryBlockMessageIdsResponse = {
  typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsResponse",
  encode(message: QueryBlockMessageIdsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    writer.uint32(10).fork();
    for (const v of message.messageIds) {
      writer.uint32(v);
    }
    writer.ldelim();
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryBlockMessageIdsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryBlockMessageIdsResponse>): QueryBlockMessageIdsResponse {
    const message = createBaseQueryBlockMessageIdsResponse();
    message.messageIds = object.messageIds?.map(e => e) || [];
    return message;
  },
  fromAmino(object: QueryBlockMessageIdsResponseAmino): QueryBlockMessageIdsResponse {
    const message = createBaseQueryBlockMessageIdsResponse();
    message.messageIds = object.message_ids?.map(e => e) || [];
    return message;
  },
  toAmino(message: QueryBlockMessageIdsResponse): QueryBlockMessageIdsResponseAmino {
    const obj: any = {};
    if (message.messageIds) {
      obj.message_ids = message.messageIds.map(e => e);
    } else {
      obj.message_ids = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryBlockMessageIdsResponseAminoMsg): QueryBlockMessageIdsResponse {
    return QueryBlockMessageIdsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryBlockMessageIdsResponseProtoMsg): QueryBlockMessageIdsResponse {
    return QueryBlockMessageIdsResponse.decode(message.value);
  },
  toProto(message: QueryBlockMessageIdsResponse): Uint8Array {
    return QueryBlockMessageIdsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryBlockMessageIdsResponse): QueryBlockMessageIdsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.delaymsg.QueryBlockMessageIdsResponse",
      value: QueryBlockMessageIdsResponse.encode(message).finish()
    };
  }
};