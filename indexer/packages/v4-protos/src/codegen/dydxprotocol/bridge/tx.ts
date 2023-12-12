import { BridgeEvent, BridgeEventAmino, BridgeEventSDKType } from "./bridge_event";
import { EventParams, EventParamsAmino, EventParamsSDKType, ProposeParams, ProposeParamsAmino, ProposeParamsSDKType, SafetyParams, SafetyParamsAmino, SafetyParamsSDKType } from "./params";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgAcknowledgeBridges is the Msg/AcknowledgeBridges request type. */
export interface MsgAcknowledgeBridges {
  /** The events to acknowledge. */
  events: BridgeEvent[];
}
export interface MsgAcknowledgeBridgesProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridges";
  value: Uint8Array;
}
/** MsgAcknowledgeBridges is the Msg/AcknowledgeBridges request type. */
export interface MsgAcknowledgeBridgesAmino {
  /** The events to acknowledge. */
  events?: BridgeEventAmino[];
}
export interface MsgAcknowledgeBridgesAminoMsg {
  type: "/dydxprotocol.bridge.MsgAcknowledgeBridges";
  value: MsgAcknowledgeBridgesAmino;
}
/** MsgAcknowledgeBridges is the Msg/AcknowledgeBridges request type. */
export interface MsgAcknowledgeBridgesSDKType {
  events: BridgeEventSDKType[];
}
/**
 * MsgAcknowledgeBridgesResponse is the Msg/AcknowledgeBridgesResponse response
 * type.
 */
export interface MsgAcknowledgeBridgesResponse {}
export interface MsgAcknowledgeBridgesResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridgesResponse";
  value: Uint8Array;
}
/**
 * MsgAcknowledgeBridgesResponse is the Msg/AcknowledgeBridgesResponse response
 * type.
 */
export interface MsgAcknowledgeBridgesResponseAmino {}
export interface MsgAcknowledgeBridgesResponseAminoMsg {
  type: "/dydxprotocol.bridge.MsgAcknowledgeBridgesResponse";
  value: MsgAcknowledgeBridgesResponseAmino;
}
/**
 * MsgAcknowledgeBridgesResponse is the Msg/AcknowledgeBridgesResponse response
 * type.
 */
export interface MsgAcknowledgeBridgesResponseSDKType {}
/** MsgCompleteBridge is the Msg/CompleteBridgeResponse request type. */
export interface MsgCompleteBridge {
  authority: string;
  /** The event to complete. */
  event: BridgeEvent;
}
export interface MsgCompleteBridgeProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgCompleteBridge";
  value: Uint8Array;
}
/** MsgCompleteBridge is the Msg/CompleteBridgeResponse request type. */
export interface MsgCompleteBridgeAmino {
  authority?: string;
  /** The event to complete. */
  event?: BridgeEventAmino;
}
export interface MsgCompleteBridgeAminoMsg {
  type: "/dydxprotocol.bridge.MsgCompleteBridge";
  value: MsgCompleteBridgeAmino;
}
/** MsgCompleteBridge is the Msg/CompleteBridgeResponse request type. */
export interface MsgCompleteBridgeSDKType {
  authority: string;
  event: BridgeEventSDKType;
}
/** MsgCompleteBridgeResponse is the Msg/CompleteBridgeResponse response type. */
export interface MsgCompleteBridgeResponse {}
export interface MsgCompleteBridgeResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgCompleteBridgeResponse";
  value: Uint8Array;
}
/** MsgCompleteBridgeResponse is the Msg/CompleteBridgeResponse response type. */
export interface MsgCompleteBridgeResponseAmino {}
export interface MsgCompleteBridgeResponseAminoMsg {
  type: "/dydxprotocol.bridge.MsgCompleteBridgeResponse";
  value: MsgCompleteBridgeResponseAmino;
}
/** MsgCompleteBridgeResponse is the Msg/CompleteBridgeResponse response type. */
export interface MsgCompleteBridgeResponseSDKType {}
/** MsgUpdateEventParams is the Msg/UpdateEventParams request type. */
export interface MsgUpdateEventParams {
  authority: string;
  /** The parameters to update. Each field must be set. */
  params: EventParams;
}
export interface MsgUpdateEventParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParams";
  value: Uint8Array;
}
/** MsgUpdateEventParams is the Msg/UpdateEventParams request type. */
export interface MsgUpdateEventParamsAmino {
  authority?: string;
  /** The parameters to update. Each field must be set. */
  params?: EventParamsAmino;
}
export interface MsgUpdateEventParamsAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateEventParams";
  value: MsgUpdateEventParamsAmino;
}
/** MsgUpdateEventParams is the Msg/UpdateEventParams request type. */
export interface MsgUpdateEventParamsSDKType {
  authority: string;
  params: EventParamsSDKType;
}
/** MsgUpdateEventParamsResponse is the Msg/UpdateEventParams response type. */
export interface MsgUpdateEventParamsResponse {}
export interface MsgUpdateEventParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParamsResponse";
  value: Uint8Array;
}
/** MsgUpdateEventParamsResponse is the Msg/UpdateEventParams response type. */
export interface MsgUpdateEventParamsResponseAmino {}
export interface MsgUpdateEventParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateEventParamsResponse";
  value: MsgUpdateEventParamsResponseAmino;
}
/** MsgUpdateEventParamsResponse is the Msg/UpdateEventParams response type. */
export interface MsgUpdateEventParamsResponseSDKType {}
/** MsgUpdateProposeParams is the Msg/UpdateProposeParams request type. */
export interface MsgUpdateProposeParams {
  authority: string;
  /** The parameters to update. Each field must be set. */
  params: ProposeParams;
}
export interface MsgUpdateProposeParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParams";
  value: Uint8Array;
}
/** MsgUpdateProposeParams is the Msg/UpdateProposeParams request type. */
export interface MsgUpdateProposeParamsAmino {
  authority?: string;
  /** The parameters to update. Each field must be set. */
  params?: ProposeParamsAmino;
}
export interface MsgUpdateProposeParamsAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateProposeParams";
  value: MsgUpdateProposeParamsAmino;
}
/** MsgUpdateProposeParams is the Msg/UpdateProposeParams request type. */
export interface MsgUpdateProposeParamsSDKType {
  authority: string;
  params: ProposeParamsSDKType;
}
/** MsgUpdateProposeParamsResponse is the Msg/UpdateProposeParams response type. */
export interface MsgUpdateProposeParamsResponse {}
export interface MsgUpdateProposeParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParamsResponse";
  value: Uint8Array;
}
/** MsgUpdateProposeParamsResponse is the Msg/UpdateProposeParams response type. */
export interface MsgUpdateProposeParamsResponseAmino {}
export interface MsgUpdateProposeParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateProposeParamsResponse";
  value: MsgUpdateProposeParamsResponseAmino;
}
/** MsgUpdateProposeParamsResponse is the Msg/UpdateProposeParams response type. */
export interface MsgUpdateProposeParamsResponseSDKType {}
/** MsgUpdateSafetyParams is the Msg/UpdateSafetyParams request type. */
export interface MsgUpdateSafetyParams {
  authority: string;
  /** The parameters to update. Each field must be set. */
  params: SafetyParams;
}
export interface MsgUpdateSafetyParamsProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParams";
  value: Uint8Array;
}
/** MsgUpdateSafetyParams is the Msg/UpdateSafetyParams request type. */
export interface MsgUpdateSafetyParamsAmino {
  authority?: string;
  /** The parameters to update. Each field must be set. */
  params?: SafetyParamsAmino;
}
export interface MsgUpdateSafetyParamsAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateSafetyParams";
  value: MsgUpdateSafetyParamsAmino;
}
/** MsgUpdateSafetyParams is the Msg/UpdateSafetyParams request type. */
export interface MsgUpdateSafetyParamsSDKType {
  authority: string;
  params: SafetyParamsSDKType;
}
/** MsgUpdateSafetyParamsResponse is the Msg/UpdateSafetyParams response type. */
export interface MsgUpdateSafetyParamsResponse {}
export interface MsgUpdateSafetyParamsResponseProtoMsg {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse";
  value: Uint8Array;
}
/** MsgUpdateSafetyParamsResponse is the Msg/UpdateSafetyParams response type. */
export interface MsgUpdateSafetyParamsResponseAmino {}
export interface MsgUpdateSafetyParamsResponseAminoMsg {
  type: "/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse";
  value: MsgUpdateSafetyParamsResponseAmino;
}
/** MsgUpdateSafetyParamsResponse is the Msg/UpdateSafetyParams response type. */
export interface MsgUpdateSafetyParamsResponseSDKType {}
function createBaseMsgAcknowledgeBridges(): MsgAcknowledgeBridges {
  return {
    events: []
  };
}
export const MsgAcknowledgeBridges = {
  typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridges",
  encode(message: MsgAcknowledgeBridges, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.events) {
      BridgeEvent.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAcknowledgeBridges {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAcknowledgeBridges();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.events.push(BridgeEvent.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgAcknowledgeBridges>): MsgAcknowledgeBridges {
    const message = createBaseMsgAcknowledgeBridges();
    message.events = object.events?.map(e => BridgeEvent.fromPartial(e)) || [];
    return message;
  },
  fromAmino(object: MsgAcknowledgeBridgesAmino): MsgAcknowledgeBridges {
    const message = createBaseMsgAcknowledgeBridges();
    message.events = object.events?.map(e => BridgeEvent.fromAmino(e)) || [];
    return message;
  },
  toAmino(message: MsgAcknowledgeBridges): MsgAcknowledgeBridgesAmino {
    const obj: any = {};
    if (message.events) {
      obj.events = message.events.map(e => e ? BridgeEvent.toAmino(e) : undefined);
    } else {
      obj.events = [];
    }
    return obj;
  },
  fromAminoMsg(object: MsgAcknowledgeBridgesAminoMsg): MsgAcknowledgeBridges {
    return MsgAcknowledgeBridges.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAcknowledgeBridgesProtoMsg): MsgAcknowledgeBridges {
    return MsgAcknowledgeBridges.decode(message.value);
  },
  toProto(message: MsgAcknowledgeBridges): Uint8Array {
    return MsgAcknowledgeBridges.encode(message).finish();
  },
  toProtoMsg(message: MsgAcknowledgeBridges): MsgAcknowledgeBridgesProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridges",
      value: MsgAcknowledgeBridges.encode(message).finish()
    };
  }
};
function createBaseMsgAcknowledgeBridgesResponse(): MsgAcknowledgeBridgesResponse {
  return {};
}
export const MsgAcknowledgeBridgesResponse = {
  typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridgesResponse",
  encode(_: MsgAcknowledgeBridgesResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgAcknowledgeBridgesResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgAcknowledgeBridgesResponse();
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
  fromPartial(_: Partial<MsgAcknowledgeBridgesResponse>): MsgAcknowledgeBridgesResponse {
    const message = createBaseMsgAcknowledgeBridgesResponse();
    return message;
  },
  fromAmino(_: MsgAcknowledgeBridgesResponseAmino): MsgAcknowledgeBridgesResponse {
    const message = createBaseMsgAcknowledgeBridgesResponse();
    return message;
  },
  toAmino(_: MsgAcknowledgeBridgesResponse): MsgAcknowledgeBridgesResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgAcknowledgeBridgesResponseAminoMsg): MsgAcknowledgeBridgesResponse {
    return MsgAcknowledgeBridgesResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgAcknowledgeBridgesResponseProtoMsg): MsgAcknowledgeBridgesResponse {
    return MsgAcknowledgeBridgesResponse.decode(message.value);
  },
  toProto(message: MsgAcknowledgeBridgesResponse): Uint8Array {
    return MsgAcknowledgeBridgesResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgAcknowledgeBridgesResponse): MsgAcknowledgeBridgesResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgAcknowledgeBridgesResponse",
      value: MsgAcknowledgeBridgesResponse.encode(message).finish()
    };
  }
};
function createBaseMsgCompleteBridge(): MsgCompleteBridge {
  return {
    authority: "",
    event: BridgeEvent.fromPartial({})
  };
}
export const MsgCompleteBridge = {
  typeUrl: "/dydxprotocol.bridge.MsgCompleteBridge",
  encode(message: MsgCompleteBridge, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.event !== undefined) {
      BridgeEvent.encode(message.event, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCompleteBridge {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCompleteBridge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.event = BridgeEvent.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgCompleteBridge>): MsgCompleteBridge {
    const message = createBaseMsgCompleteBridge();
    message.authority = object.authority ?? "";
    message.event = object.event !== undefined && object.event !== null ? BridgeEvent.fromPartial(object.event) : undefined;
    return message;
  },
  fromAmino(object: MsgCompleteBridgeAmino): MsgCompleteBridge {
    const message = createBaseMsgCompleteBridge();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.event !== undefined && object.event !== null) {
      message.event = BridgeEvent.fromAmino(object.event);
    }
    return message;
  },
  toAmino(message: MsgCompleteBridge): MsgCompleteBridgeAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.event = message.event ? BridgeEvent.toAmino(message.event) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgCompleteBridgeAminoMsg): MsgCompleteBridge {
    return MsgCompleteBridge.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCompleteBridgeProtoMsg): MsgCompleteBridge {
    return MsgCompleteBridge.decode(message.value);
  },
  toProto(message: MsgCompleteBridge): Uint8Array {
    return MsgCompleteBridge.encode(message).finish();
  },
  toProtoMsg(message: MsgCompleteBridge): MsgCompleteBridgeProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgCompleteBridge",
      value: MsgCompleteBridge.encode(message).finish()
    };
  }
};
function createBaseMsgCompleteBridgeResponse(): MsgCompleteBridgeResponse {
  return {};
}
export const MsgCompleteBridgeResponse = {
  typeUrl: "/dydxprotocol.bridge.MsgCompleteBridgeResponse",
  encode(_: MsgCompleteBridgeResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgCompleteBridgeResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCompleteBridgeResponse();
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
  fromPartial(_: Partial<MsgCompleteBridgeResponse>): MsgCompleteBridgeResponse {
    const message = createBaseMsgCompleteBridgeResponse();
    return message;
  },
  fromAmino(_: MsgCompleteBridgeResponseAmino): MsgCompleteBridgeResponse {
    const message = createBaseMsgCompleteBridgeResponse();
    return message;
  },
  toAmino(_: MsgCompleteBridgeResponse): MsgCompleteBridgeResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgCompleteBridgeResponseAminoMsg): MsgCompleteBridgeResponse {
    return MsgCompleteBridgeResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgCompleteBridgeResponseProtoMsg): MsgCompleteBridgeResponse {
    return MsgCompleteBridgeResponse.decode(message.value);
  },
  toProto(message: MsgCompleteBridgeResponse): Uint8Array {
    return MsgCompleteBridgeResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgCompleteBridgeResponse): MsgCompleteBridgeResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgCompleteBridgeResponse",
      value: MsgCompleteBridgeResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateEventParams(): MsgUpdateEventParams {
  return {
    authority: "",
    params: EventParams.fromPartial({})
  };
}
export const MsgUpdateEventParams = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParams",
  encode(message: MsgUpdateEventParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      EventParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateEventParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateEventParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.params = EventParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateEventParams>): MsgUpdateEventParams {
    const message = createBaseMsgUpdateEventParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? EventParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateEventParamsAmino): MsgUpdateEventParams {
    const message = createBaseMsgUpdateEventParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = EventParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgUpdateEventParams): MsgUpdateEventParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? EventParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateEventParamsAminoMsg): MsgUpdateEventParams {
    return MsgUpdateEventParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateEventParamsProtoMsg): MsgUpdateEventParams {
    return MsgUpdateEventParams.decode(message.value);
  },
  toProto(message: MsgUpdateEventParams): Uint8Array {
    return MsgUpdateEventParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateEventParams): MsgUpdateEventParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParams",
      value: MsgUpdateEventParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateEventParamsResponse(): MsgUpdateEventParamsResponse {
  return {};
}
export const MsgUpdateEventParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParamsResponse",
  encode(_: MsgUpdateEventParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateEventParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateEventParamsResponse();
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
  fromPartial(_: Partial<MsgUpdateEventParamsResponse>): MsgUpdateEventParamsResponse {
    const message = createBaseMsgUpdateEventParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateEventParamsResponseAmino): MsgUpdateEventParamsResponse {
    const message = createBaseMsgUpdateEventParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateEventParamsResponse): MsgUpdateEventParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateEventParamsResponseAminoMsg): MsgUpdateEventParamsResponse {
    return MsgUpdateEventParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateEventParamsResponseProtoMsg): MsgUpdateEventParamsResponse {
    return MsgUpdateEventParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateEventParamsResponse): Uint8Array {
    return MsgUpdateEventParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateEventParamsResponse): MsgUpdateEventParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateEventParamsResponse",
      value: MsgUpdateEventParamsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateProposeParams(): MsgUpdateProposeParams {
  return {
    authority: "",
    params: ProposeParams.fromPartial({})
  };
}
export const MsgUpdateProposeParams = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParams",
  encode(message: MsgUpdateProposeParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      ProposeParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateProposeParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateProposeParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.params = ProposeParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateProposeParams>): MsgUpdateProposeParams {
    const message = createBaseMsgUpdateProposeParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? ProposeParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateProposeParamsAmino): MsgUpdateProposeParams {
    const message = createBaseMsgUpdateProposeParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = ProposeParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgUpdateProposeParams): MsgUpdateProposeParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? ProposeParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateProposeParamsAminoMsg): MsgUpdateProposeParams {
    return MsgUpdateProposeParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateProposeParamsProtoMsg): MsgUpdateProposeParams {
    return MsgUpdateProposeParams.decode(message.value);
  },
  toProto(message: MsgUpdateProposeParams): Uint8Array {
    return MsgUpdateProposeParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateProposeParams): MsgUpdateProposeParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParams",
      value: MsgUpdateProposeParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateProposeParamsResponse(): MsgUpdateProposeParamsResponse {
  return {};
}
export const MsgUpdateProposeParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParamsResponse",
  encode(_: MsgUpdateProposeParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateProposeParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateProposeParamsResponse();
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
  fromPartial(_: Partial<MsgUpdateProposeParamsResponse>): MsgUpdateProposeParamsResponse {
    const message = createBaseMsgUpdateProposeParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateProposeParamsResponseAmino): MsgUpdateProposeParamsResponse {
    const message = createBaseMsgUpdateProposeParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateProposeParamsResponse): MsgUpdateProposeParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateProposeParamsResponseAminoMsg): MsgUpdateProposeParamsResponse {
    return MsgUpdateProposeParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateProposeParamsResponseProtoMsg): MsgUpdateProposeParamsResponse {
    return MsgUpdateProposeParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateProposeParamsResponse): Uint8Array {
    return MsgUpdateProposeParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateProposeParamsResponse): MsgUpdateProposeParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateProposeParamsResponse",
      value: MsgUpdateProposeParamsResponse.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateSafetyParams(): MsgUpdateSafetyParams {
  return {
    authority: "",
    params: SafetyParams.fromPartial({})
  };
}
export const MsgUpdateSafetyParams = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParams",
  encode(message: MsgUpdateSafetyParams, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.params !== undefined) {
      SafetyParams.encode(message.params, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateSafetyParams {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSafetyParams();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.params = SafetyParams.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgUpdateSafetyParams>): MsgUpdateSafetyParams {
    const message = createBaseMsgUpdateSafetyParams();
    message.authority = object.authority ?? "";
    message.params = object.params !== undefined && object.params !== null ? SafetyParams.fromPartial(object.params) : undefined;
    return message;
  },
  fromAmino(object: MsgUpdateSafetyParamsAmino): MsgUpdateSafetyParams {
    const message = createBaseMsgUpdateSafetyParams();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.params !== undefined && object.params !== null) {
      message.params = SafetyParams.fromAmino(object.params);
    }
    return message;
  },
  toAmino(message: MsgUpdateSafetyParams): MsgUpdateSafetyParamsAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.params = message.params ? SafetyParams.toAmino(message.params) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgUpdateSafetyParamsAminoMsg): MsgUpdateSafetyParams {
    return MsgUpdateSafetyParams.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateSafetyParamsProtoMsg): MsgUpdateSafetyParams {
    return MsgUpdateSafetyParams.decode(message.value);
  },
  toProto(message: MsgUpdateSafetyParams): Uint8Array {
    return MsgUpdateSafetyParams.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateSafetyParams): MsgUpdateSafetyParamsProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParams",
      value: MsgUpdateSafetyParams.encode(message).finish()
    };
  }
};
function createBaseMsgUpdateSafetyParamsResponse(): MsgUpdateSafetyParamsResponse {
  return {};
}
export const MsgUpdateSafetyParamsResponse = {
  typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse",
  encode(_: MsgUpdateSafetyParamsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgUpdateSafetyParamsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgUpdateSafetyParamsResponse();
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
  fromPartial(_: Partial<MsgUpdateSafetyParamsResponse>): MsgUpdateSafetyParamsResponse {
    const message = createBaseMsgUpdateSafetyParamsResponse();
    return message;
  },
  fromAmino(_: MsgUpdateSafetyParamsResponseAmino): MsgUpdateSafetyParamsResponse {
    const message = createBaseMsgUpdateSafetyParamsResponse();
    return message;
  },
  toAmino(_: MsgUpdateSafetyParamsResponse): MsgUpdateSafetyParamsResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgUpdateSafetyParamsResponseAminoMsg): MsgUpdateSafetyParamsResponse {
    return MsgUpdateSafetyParamsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgUpdateSafetyParamsResponseProtoMsg): MsgUpdateSafetyParamsResponse {
    return MsgUpdateSafetyParamsResponse.decode(message.value);
  },
  toProto(message: MsgUpdateSafetyParamsResponse): Uint8Array {
    return MsgUpdateSafetyParamsResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgUpdateSafetyParamsResponse): MsgUpdateSafetyParamsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.bridge.MsgUpdateSafetyParamsResponse",
      value: MsgUpdateSafetyParamsResponse.encode(message).finish()
    };
  }
};