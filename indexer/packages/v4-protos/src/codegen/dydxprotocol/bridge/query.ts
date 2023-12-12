import { EventParams, EventParamsAmino, EventParamsSDKType, ProposeParams, ProposeParamsAmino, ProposeParamsSDKType, SafetyParams, SafetyParamsAmino, SafetyParamsSDKType } from "./params";
import { BridgeEventInfo, BridgeEventInfoAmino, BridgeEventInfoSDKType } from "./bridge_event_info";
import { MsgCompleteBridge, MsgCompleteBridgeAmino, MsgCompleteBridgeSDKType } from "./tx";
import { BinaryReader, BinaryWriter } from "../../binary";
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */
export interface QueryEventParamsRequest {}
export interface QueryEventParamsRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryEventParamsRequest";
  value: Uint8Array;
}
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */
export interface QueryEventParamsRequestAmino {}
export interface QueryEventParamsRequestAminoMsg {
  type: "/dydxprotocol.bridge.QueryEventParamsRequest";
  value: QueryEventParamsRequestAmino;
}
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */
export interface QueryEventParamsRequestSDKType {}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */
export interface QueryEventParamsResponse {
  params: EventParams;
}
export interface QueryEventParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryEventParamsResponse";
  value: Uint8Array;
}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */
export interface QueryEventParamsResponseAmino {
  params?: EventParamsAmino;
}
export interface QueryEventParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.QueryEventParamsResponse";
  value: QueryEventParamsResponseAmino;
}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */
export interface QueryEventParamsResponseSDKType {
  params: EventParamsSDKType;
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */
export interface QueryProposeParamsRequest {}
export interface QueryProposeParamsRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryProposeParamsRequest";
  value: Uint8Array;
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */
export interface QueryProposeParamsRequestAmino {}
export interface QueryProposeParamsRequestAminoMsg {
  type: "/dydxprotocol.bridge.QueryProposeParamsRequest";
  value: QueryProposeParamsRequestAmino;
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */
export interface QueryProposeParamsRequestSDKType {}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */
export interface QueryProposeParamsResponse {
  params: ProposeParams;
}
export interface QueryProposeParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryProposeParamsResponse";
  value: Uint8Array;
}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */
export interface QueryProposeParamsResponseAmino {
  params?: ProposeParamsAmino;
}
export interface QueryProposeParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.QueryProposeParamsResponse";
  value: QueryProposeParamsResponseAmino;
}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */
export interface QueryProposeParamsResponseSDKType {
  params: ProposeParamsSDKType;
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */
export interface QuerySafetyParamsRequest {}
export interface QuerySafetyParamsRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsRequest";
  value: Uint8Array;
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */
export interface QuerySafetyParamsRequestAmino {}
export interface QuerySafetyParamsRequestAminoMsg {
  type: "/dydxprotocol.bridge.QuerySafetyParamsRequest";
  value: QuerySafetyParamsRequestAmino;
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */
export interface QuerySafetyParamsRequestSDKType {}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */
export interface QuerySafetyParamsResponse {
  params: SafetyParams;
}
export interface QuerySafetyParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsResponse";
  value: Uint8Array;
}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */
export interface QuerySafetyParamsResponseAmino {
  params?: SafetyParamsAmino;
}
export interface QuerySafetyParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.QuerySafetyParamsResponse";
  value: QuerySafetyParamsResponseAmino;
}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */
export interface QuerySafetyParamsResponseSDKType {
  params: SafetyParamsSDKType;
}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoRequest {}
export interface QueryAcknowledgedEventInfoRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoRequest";
  value: Uint8Array;
}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoRequestAmino {}
export interface QueryAcknowledgedEventInfoRequestAminoMsg {
  type: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoRequest";
  value: QueryAcknowledgedEventInfoRequestAmino;
}
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
  info: BridgeEventInfo;
}
export interface QueryAcknowledgedEventInfoResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoResponse";
  value: Uint8Array;
}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoResponseAmino {
  info?: BridgeEventInfoAmino;
}
export interface QueryAcknowledgedEventInfoResponseAminoMsg {
  type: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoResponse";
  value: QueryAcknowledgedEventInfoResponseAmino;
}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoResponseSDKType {
  info: BridgeEventInfoSDKType;
}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoRequest {}
export interface QueryRecognizedEventInfoRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoRequest";
  value: Uint8Array;
}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoRequestAmino {}
export interface QueryRecognizedEventInfoRequestAminoMsg {
  type: "/dydxprotocol.bridge.QueryRecognizedEventInfoRequest";
  value: QueryRecognizedEventInfoRequestAmino;
}
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
  info: BridgeEventInfo;
}
export interface QueryRecognizedEventInfoResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoResponse";
  value: Uint8Array;
}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoResponseAmino {
  info?: BridgeEventInfoAmino;
}
export interface QueryRecognizedEventInfoResponseAminoMsg {
  type: "/dydxprotocol.bridge.QueryRecognizedEventInfoResponse";
  value: QueryRecognizedEventInfoResponseAmino;
}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoResponseSDKType {
  info: BridgeEventInfoSDKType;
}
/**
 * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesRequest {
  /**
   * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
   * DelayedCompleteBridgeMessages RPC method.
   */
  address: string;
}
export interface QueryDelayedCompleteBridgeMessagesRequestProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesRequest";
  value: Uint8Array;
}
/**
 * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesRequestAmino {
  /**
   * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
   * DelayedCompleteBridgeMessages RPC method.
   */
  address?: string;
}
export interface QueryDelayedCompleteBridgeMessagesRequestAminoMsg {
  type: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesRequest";
  value: QueryDelayedCompleteBridgeMessagesRequestAmino;
}
/**
 * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesRequestSDKType {
  address: string;
}
/**
 * QueryDelayedCompleteBridgeMessagesResponse is a response type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesResponse {
  messages: DelayedCompleteBridgeMessage[];
}
export interface QueryDelayedCompleteBridgeMessagesResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesResponse";
  value: Uint8Array;
}
/**
 * QueryDelayedCompleteBridgeMessagesResponse is a response type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesResponseAmino {
  messages?: DelayedCompleteBridgeMessageAmino[];
}
export interface QueryDelayedCompleteBridgeMessagesResponseAminoMsg {
  type: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesResponse";
  value: QueryDelayedCompleteBridgeMessagesResponseAmino;
}
/**
 * QueryDelayedCompleteBridgeMessagesResponse is a response type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesResponseSDKType {
  messages: DelayedCompleteBridgeMessageSDKType[];
}
/**
 * DelayedCompleteBridgeMessage is a message type for the response of
 * DelayedCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */
export interface DelayedCompleteBridgeMessage {
  message: MsgCompleteBridge;
  blockHeight: number;
}
export interface DelayedCompleteBridgeMessageProtoMsg {
  typeUrl: "/dydxprotocol.bridge.DelayedCompleteBridgeMessage";
  value: Uint8Array;
}
/**
 * DelayedCompleteBridgeMessage is a message type for the response of
 * DelayedCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */
export interface DelayedCompleteBridgeMessageAmino {
  message?: MsgCompleteBridgeAmino;
  block_height?: number;
}
export interface DelayedCompleteBridgeMessageAminoMsg {
  type: "/dydxprotocol.bridge.DelayedCompleteBridgeMessage";
  value: DelayedCompleteBridgeMessageAmino;
}
/**
 * DelayedCompleteBridgeMessage is a message type for the response of
 * DelayedCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */
export interface DelayedCompleteBridgeMessageSDKType {
  message: MsgCompleteBridgeSDKType;
  block_height: number;
}
function createBaseQueryEventParamsRequest(): QueryEventParamsRequest {
  return {};
}
export const QueryEventParamsRequest = {
  typeUrl: "/dydxprotocol.bridge.QueryEventParamsRequest",
  encode(_: QueryEventParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryEventParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QueryEventParamsRequest>): QueryEventParamsRequest {
    const message = createBaseQueryEventParamsRequest();
    return message;
  },
  fromAmino(_: QueryEventParamsRequestAmino): QueryEventParamsRequest {
    const message = createBaseQueryEventParamsRequest();
    return message;
  },
  toAmino(_: QueryEventParamsRequest): QueryEventParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryEventParamsRequestAminoMsg): QueryEventParamsRequest {
    return QueryEventParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryEventParamsRequestProtoMsg): QueryEventParamsRequest {
    return QueryEventParamsRequest.decode(message.value);
  },
  toProto(message: QueryEventParamsRequest): Uint8Array {
    return QueryEventParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryEventParamsRequest): QueryEventParamsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryEventParamsRequest",
      value: QueryEventParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryEventParamsResponse(): QueryEventParamsResponse {
  return {
    params: EventParams.fromPartial({})
  };
}
export const QueryEventParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.QueryEventParamsResponse",
  encode(message: QueryEventParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      EventParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryEventParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryEventParamsResponse>): QueryEventParamsResponse {
    const message = createBaseQueryEventParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? EventParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QueryEventParamsResponseAmino): QueryEventParamsResponse {
    const message = createBaseQueryEventParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = EventParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QueryEventParamsResponse): QueryEventParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? EventParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryEventParamsResponseAminoMsg): QueryEventParamsResponse {
    return QueryEventParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryEventParamsResponseProtoMsg): QueryEventParamsResponse {
    return QueryEventParamsResponse.decode(message.value);
  },
  toProto(message: QueryEventParamsResponse): Uint8Array {
    return QueryEventParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryEventParamsResponse): QueryEventParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryEventParamsResponse",
      value: QueryEventParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryProposeParamsRequest(): QueryProposeParamsRequest {
  return {};
}
export const QueryProposeParamsRequest = {
  typeUrl: "/dydxprotocol.bridge.QueryProposeParamsRequest",
  encode(_: QueryProposeParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryProposeParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QueryProposeParamsRequest>): QueryProposeParamsRequest {
    const message = createBaseQueryProposeParamsRequest();
    return message;
  },
  fromAmino(_: QueryProposeParamsRequestAmino): QueryProposeParamsRequest {
    const message = createBaseQueryProposeParamsRequest();
    return message;
  },
  toAmino(_: QueryProposeParamsRequest): QueryProposeParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryProposeParamsRequestAminoMsg): QueryProposeParamsRequest {
    return QueryProposeParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProposeParamsRequestProtoMsg): QueryProposeParamsRequest {
    return QueryProposeParamsRequest.decode(message.value);
  },
  toProto(message: QueryProposeParamsRequest): Uint8Array {
    return QueryProposeParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryProposeParamsRequest): QueryProposeParamsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryProposeParamsRequest",
      value: QueryProposeParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryProposeParamsResponse(): QueryProposeParamsResponse {
  return {
    params: ProposeParams.fromPartial({})
  };
}
export const QueryProposeParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.QueryProposeParamsResponse",
  encode(message: QueryProposeParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      ProposeParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryProposeParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryProposeParamsResponse>): QueryProposeParamsResponse {
    const message = createBaseQueryProposeParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? ProposeParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QueryProposeParamsResponseAmino): QueryProposeParamsResponse {
    const message = createBaseQueryProposeParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = ProposeParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QueryProposeParamsResponse): QueryProposeParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? ProposeParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryProposeParamsResponseAminoMsg): QueryProposeParamsResponse {
    return QueryProposeParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryProposeParamsResponseProtoMsg): QueryProposeParamsResponse {
    return QueryProposeParamsResponse.decode(message.value);
  },
  toProto(message: QueryProposeParamsResponse): Uint8Array {
    return QueryProposeParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryProposeParamsResponse): QueryProposeParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryProposeParamsResponse",
      value: QueryProposeParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQuerySafetyParamsRequest(): QuerySafetyParamsRequest {
  return {};
}
export const QuerySafetyParamsRequest = {
  typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsRequest",
  encode(_: QuerySafetyParamsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySafetyParamsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QuerySafetyParamsRequest>): QuerySafetyParamsRequest {
    const message = createBaseQuerySafetyParamsRequest();
    return message;
  },
  fromAmino(_: QuerySafetyParamsRequestAmino): QuerySafetyParamsRequest {
    const message = createBaseQuerySafetyParamsRequest();
    return message;
  },
  toAmino(_: QuerySafetyParamsRequest): QuerySafetyParamsRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QuerySafetyParamsRequestAminoMsg): QuerySafetyParamsRequest {
    return QuerySafetyParamsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySafetyParamsRequestProtoMsg): QuerySafetyParamsRequest {
    return QuerySafetyParamsRequest.decode(message.value);
  },
  toProto(message: QuerySafetyParamsRequest): Uint8Array {
    return QuerySafetyParamsRequest.encode(message).finish();
  },
  toProtoMsg(message: QuerySafetyParamsRequest): QuerySafetyParamsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsRequest",
      value: QuerySafetyParamsRequest.encode(message).finish()
    };
  }
};
function createBaseQuerySafetyParamsResponse(): QuerySafetyParamsResponse {
  return {
    params: SafetyParams.fromPartial({})
  };
}
export const QuerySafetyParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsResponse",
  encode(message: QuerySafetyParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.params !== undefined) {
      SafetyParams.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QuerySafetyParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QuerySafetyParamsResponse>): QuerySafetyParamsResponse {
    const message = createBaseQuerySafetyParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? SafetyParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: QuerySafetyParamsResponseAmino): QuerySafetyParamsResponse {
    const message = createBaseQuerySafetyParamsResponse();
    if (object.params !== undefined && object.params !== null) {
      message.params = SafetyParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: QuerySafetyParamsResponse): QuerySafetyParamsResponseAmino {
    const obj: any = {};
    obj.params = message.params ? SafetyParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: QuerySafetyParamsResponseAminoMsg): QuerySafetyParamsResponse {
    return QuerySafetyParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QuerySafetyParamsResponseProtoMsg): QuerySafetyParamsResponse {
    return QuerySafetyParamsResponse.decode(message.value);
  },
  toProto(message: QuerySafetyParamsResponse): Uint8Array {
    return QuerySafetyParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: QuerySafetyParamsResponse): QuerySafetyParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QuerySafetyParamsResponse",
      value: QuerySafetyParamsResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAcknowledgedEventInfoRequest(): QueryAcknowledgedEventInfoRequest {
  return {};
}
export const QueryAcknowledgedEventInfoRequest = {
  typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoRequest",
  encode(_: QueryAcknowledgedEventInfoRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAcknowledgedEventInfoRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QueryAcknowledgedEventInfoRequest>): QueryAcknowledgedEventInfoRequest {
    const message = createBaseQueryAcknowledgedEventInfoRequest();
    return message;
  },
  fromAmino(_: QueryAcknowledgedEventInfoRequestAmino): QueryAcknowledgedEventInfoRequest {
    const message = createBaseQueryAcknowledgedEventInfoRequest();
    return message;
  },
  toAmino(_: QueryAcknowledgedEventInfoRequest): QueryAcknowledgedEventInfoRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryAcknowledgedEventInfoRequestAminoMsg): QueryAcknowledgedEventInfoRequest {
    return QueryAcknowledgedEventInfoRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAcknowledgedEventInfoRequestProtoMsg): QueryAcknowledgedEventInfoRequest {
    return QueryAcknowledgedEventInfoRequest.decode(message.value);
  },
  toProto(message: QueryAcknowledgedEventInfoRequest): Uint8Array {
    return QueryAcknowledgedEventInfoRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAcknowledgedEventInfoRequest): QueryAcknowledgedEventInfoRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoRequest",
      value: QueryAcknowledgedEventInfoRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAcknowledgedEventInfoResponse(): QueryAcknowledgedEventInfoResponse {
  return {
    info: BridgeEventInfo.fromPartial({})
  };
}
export const QueryAcknowledgedEventInfoResponse = {
  typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoResponse",
  encode(message: QueryAcknowledgedEventInfoResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.info !== undefined) {
      BridgeEventInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAcknowledgedEventInfoResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAcknowledgedEventInfoResponse>): QueryAcknowledgedEventInfoResponse {
    const message = createBaseQueryAcknowledgedEventInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? BridgeEventInfo.fromPartial(object.info) : undefined;
    return message;
  },
  fromAmino(object: QueryAcknowledgedEventInfoResponseAmino): QueryAcknowledgedEventInfoResponse {
    const message = createBaseQueryAcknowledgedEventInfoResponse();
    if (object.info !== undefined && object.info !== null) {
      message.info = BridgeEventInfo.fromAmino(object.info);
    }
    return message;
  },
  toAmino(message: QueryAcknowledgedEventInfoResponse): QueryAcknowledgedEventInfoResponseAmino {
    const obj: any = {};
    obj.info = message.info ? BridgeEventInfo.toAmino(message.info) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAcknowledgedEventInfoResponseAminoMsg): QueryAcknowledgedEventInfoResponse {
    return QueryAcknowledgedEventInfoResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAcknowledgedEventInfoResponseProtoMsg): QueryAcknowledgedEventInfoResponse {
    return QueryAcknowledgedEventInfoResponse.decode(message.value);
  },
  toProto(message: QueryAcknowledgedEventInfoResponse): Uint8Array {
    return QueryAcknowledgedEventInfoResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAcknowledgedEventInfoResponse): QueryAcknowledgedEventInfoResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryAcknowledgedEventInfoResponse",
      value: QueryAcknowledgedEventInfoResponse.encode(message).finish()
    };
  }
};
function createBaseQueryRecognizedEventInfoRequest(): QueryRecognizedEventInfoRequest {
  return {};
}
export const QueryRecognizedEventInfoRequest = {
  typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoRequest",
  encode(_: QueryRecognizedEventInfoRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryRecognizedEventInfoRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(_: Partial<QueryRecognizedEventInfoRequest>): QueryRecognizedEventInfoRequest {
    const message = createBaseQueryRecognizedEventInfoRequest();
    return message;
  },
  fromAmino(_: QueryRecognizedEventInfoRequestAmino): QueryRecognizedEventInfoRequest {
    const message = createBaseQueryRecognizedEventInfoRequest();
    return message;
  },
  toAmino(_: QueryRecognizedEventInfoRequest): QueryRecognizedEventInfoRequestAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: QueryRecognizedEventInfoRequestAminoMsg): QueryRecognizedEventInfoRequest {
    return QueryRecognizedEventInfoRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRecognizedEventInfoRequestProtoMsg): QueryRecognizedEventInfoRequest {
    return QueryRecognizedEventInfoRequest.decode(message.value);
  },
  toProto(message: QueryRecognizedEventInfoRequest): Uint8Array {
    return QueryRecognizedEventInfoRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryRecognizedEventInfoRequest): QueryRecognizedEventInfoRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoRequest",
      value: QueryRecognizedEventInfoRequest.encode(message).finish()
    };
  }
};
function createBaseQueryRecognizedEventInfoResponse(): QueryRecognizedEventInfoResponse {
  return {
    info: BridgeEventInfo.fromPartial({})
  };
}
export const QueryRecognizedEventInfoResponse = {
  typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoResponse",
  encode(message: QueryRecognizedEventInfoResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.info !== undefined) {
      BridgeEventInfo.encode(message.info, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryRecognizedEventInfoResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryRecognizedEventInfoResponse>): QueryRecognizedEventInfoResponse {
    const message = createBaseQueryRecognizedEventInfoResponse();
    message.info = object.info !== undefined && object.info !== null ? BridgeEventInfo.fromPartial(object.info) : undefined;
    return message;
  },
  fromAmino(object: QueryRecognizedEventInfoResponseAmino): QueryRecognizedEventInfoResponse {
    const message = createBaseQueryRecognizedEventInfoResponse();
    if (object.info !== undefined && object.info !== null) {
      message.info = BridgeEventInfo.fromAmino(object.info);
    }
    return message;
  },
  toAmino(message: QueryRecognizedEventInfoResponse): QueryRecognizedEventInfoResponseAmino {
    const obj: any = {};
    obj.info = message.info ? BridgeEventInfo.toAmino(message.info) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryRecognizedEventInfoResponseAminoMsg): QueryRecognizedEventInfoResponse {
    return QueryRecognizedEventInfoResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryRecognizedEventInfoResponseProtoMsg): QueryRecognizedEventInfoResponse {
    return QueryRecognizedEventInfoResponse.decode(message.value);
  },
  toProto(message: QueryRecognizedEventInfoResponse): Uint8Array {
    return QueryRecognizedEventInfoResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryRecognizedEventInfoResponse): QueryRecognizedEventInfoResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryRecognizedEventInfoResponse",
      value: QueryRecognizedEventInfoResponse.encode(message).finish()
    };
  }
};
function createBaseQueryDelayedCompleteBridgeMessagesRequest(): QueryDelayedCompleteBridgeMessagesRequest {
  return {
    address: ""
  };
}
export const QueryDelayedCompleteBridgeMessagesRequest = {
  typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesRequest",
  encode(message: QueryDelayedCompleteBridgeMessagesRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryDelayedCompleteBridgeMessagesRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelayedCompleteBridgeMessagesRequest();
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
  fromPartial(object: Partial<QueryDelayedCompleteBridgeMessagesRequest>): QueryDelayedCompleteBridgeMessagesRequest {
    const message = createBaseQueryDelayedCompleteBridgeMessagesRequest();
    message.address = object.address ?? "";
    return message;
  },
  fromAmino(object: QueryDelayedCompleteBridgeMessagesRequestAmino): QueryDelayedCompleteBridgeMessagesRequest {
    const message = createBaseQueryDelayedCompleteBridgeMessagesRequest();
    if (object.address !== undefined && object.address !== null) {
      message.address = object.address;
    }
    return message;
  },
  toAmino(message: QueryDelayedCompleteBridgeMessagesRequest): QueryDelayedCompleteBridgeMessagesRequestAmino {
    const obj: any = {};
    obj.address = message.address;
    return obj;
  },
  fromAminoMsg(object: QueryDelayedCompleteBridgeMessagesRequestAminoMsg): QueryDelayedCompleteBridgeMessagesRequest {
    return QueryDelayedCompleteBridgeMessagesRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelayedCompleteBridgeMessagesRequestProtoMsg): QueryDelayedCompleteBridgeMessagesRequest {
    return QueryDelayedCompleteBridgeMessagesRequest.decode(message.value);
  },
  toProto(message: QueryDelayedCompleteBridgeMessagesRequest): Uint8Array {
    return QueryDelayedCompleteBridgeMessagesRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryDelayedCompleteBridgeMessagesRequest): QueryDelayedCompleteBridgeMessagesRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesRequest",
      value: QueryDelayedCompleteBridgeMessagesRequest.encode(message).finish()
    };
  }
};
function createBaseQueryDelayedCompleteBridgeMessagesResponse(): QueryDelayedCompleteBridgeMessagesResponse {
  return {
    messages: []
  };
}
export const QueryDelayedCompleteBridgeMessagesResponse = {
  typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesResponse",
  encode(message: QueryDelayedCompleteBridgeMessagesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.messages) {
      DelayedCompleteBridgeMessage.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryDelayedCompleteBridgeMessagesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryDelayedCompleteBridgeMessagesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.messages.push(DelayedCompleteBridgeMessage.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<QueryDelayedCompleteBridgeMessagesResponse>): QueryDelayedCompleteBridgeMessagesResponse {
    const message = createBaseQueryDelayedCompleteBridgeMessagesResponse();
    message.messages = object.messages?.map(e => DelayedCompleteBridgeMessage.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: QueryDelayedCompleteBridgeMessagesResponseAmino): QueryDelayedCompleteBridgeMessagesResponse {
    const message = createBaseQueryDelayedCompleteBridgeMessagesResponse();
    message.messages = object.messages?.map(e => DelayedCompleteBridgeMessage.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: QueryDelayedCompleteBridgeMessagesResponse): QueryDelayedCompleteBridgeMessagesResponseAmino {
    const obj: any = {};
    if (message.messages) {
      obj.messages = message.messages.map(e => e ? DelayedCompleteBridgeMessage.toAmino(e) : undefined);
    } else {
      obj.messages = [];
    }
    return obj;
  },
  fromAminoMsg(object: QueryDelayedCompleteBridgeMessagesResponseAminoMsg): QueryDelayedCompleteBridgeMessagesResponse {
    return QueryDelayedCompleteBridgeMessagesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryDelayedCompleteBridgeMessagesResponseProtoMsg): QueryDelayedCompleteBridgeMessagesResponse {
    return QueryDelayedCompleteBridgeMessagesResponse.decode(message.value);
  },
  toProto(message: QueryDelayedCompleteBridgeMessagesResponse): Uint8Array {
    return QueryDelayedCompleteBridgeMessagesResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryDelayedCompleteBridgeMessagesResponse): QueryDelayedCompleteBridgeMessagesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.QueryDelayedCompleteBridgeMessagesResponse",
      value: QueryDelayedCompleteBridgeMessagesResponse.encode(message).finish()
    };
  }
};
function createBaseDelayedCompleteBridgeMessage(): DelayedCompleteBridgeMessage {
  return {
    message: MsgCompleteBridge.fromPartial({}),
    blockHeight: 0
  };
}
export const DelayedCompleteBridgeMessage = {
  typeUrl: "/dydxprotocol.bridge.DelayedCompleteBridgeMessage",
  encode(message: DelayedCompleteBridgeMessage, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.message !== undefined) {
      MsgCompleteBridge.encode(message.message, writer.uint32(10).fork()).ldelim();
    }
    if (message.blockHeight !== 0) {
      writer.uint32(16).uint32(message.blockHeight);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): DelayedCompleteBridgeMessage {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDelayedCompleteBridgeMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.message = MsgCompleteBridge.decode(reader, reader.uint32());
          break;
        case 2:
          message.blockHeight = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<DelayedCompleteBridgeMessage>): DelayedCompleteBridgeMessage {
    const message = createBaseDelayedCompleteBridgeMessage();
    message.message = object.message !== undefined && object.message !== null ? MsgCompleteBridge.fromPartial(object.message) : undefined;
    message.blockHeight = object.blockHeight ?? 0;
    return message;
  },
  fromAmino(object: DelayedCompleteBridgeMessageAmino): DelayedCompleteBridgeMessage {
    const message = createBaseDelayedCompleteBridgeMessage();
    if (object.message !== undefined && object.message !== null) {
      message.message = MsgCompleteBridge.fromAmino(object.message);
    }
    if (object.block_height !== undefined && object.block_height !== null) {
      message.blockHeight = object.block_height;
    }
    return message;
  },
  toAmino(message: DelayedCompleteBridgeMessage): DelayedCompleteBridgeMessageAmino {
    const obj: any = {};
    obj.message = message.message ? MsgCompleteBridge.toAmino(message.message) : undefined;
    obj.block_height = message.blockHeight;
    return obj;
  },
  fromAminoMsg(object: DelayedCompleteBridgeMessageAminoMsg): DelayedCompleteBridgeMessage {
    return DelayedCompleteBridgeMessage.fromAmino(object.value);
  },
  fromProtoMsg(message: DelayedCompleteBridgeMessageProtoMsg): DelayedCompleteBridgeMessage {
    return DelayedCompleteBridgeMessage.decode(message.value);
  },
  toProto(message: DelayedCompleteBridgeMessage): Uint8Array {
    return DelayedCompleteBridgeMessage.encode(message).finish();
  },
  toProtoMsg(message: DelayedCompleteBridgeMessage): DelayedCompleteBridgeMessageProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.DelayedCompleteBridgeMessage",
      value: DelayedCompleteBridgeMessage.encode(message).finish()
    };
  }
};