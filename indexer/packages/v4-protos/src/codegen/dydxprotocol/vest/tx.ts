import { VestEntry, VestEntryAmino, VestEntrySDKType } from "./vest_entry";
import { BinaryReader, BinaryWriter } from "../../binary";
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */
export interface MsgDeleteVestEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vester account of the vest entry to delete. */
  vesterAccount: string;
}
export interface MsgDeleteVestEntryProtoMsg {
  typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntry";
  value: Uint8Array;
}
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */
export interface MsgDeleteVestEntryAmino {
  /** authority is the address that controls the module. */
  authority?: string;
  /** The vester account of the vest entry to delete. */
  vester_account?: string;
}
export interface MsgDeleteVestEntryAminoMsg {
  type: "/dydxprotocol.vest.MsgDeleteVestEntry";
  value: MsgDeleteVestEntryAmino;
}
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */
export interface MsgDeleteVestEntrySDKType {
  authority: string;
  vester_account: string;
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */
export interface MsgDeleteVestEntryResponse {}
export interface MsgDeleteVestEntryResponseProtoMsg {
  typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntryResponse";
  value: Uint8Array;
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */
export interface MsgDeleteVestEntryResponseAmino {}
export interface MsgDeleteVestEntryResponseAminoMsg {
  type: "/dydxprotocol.vest.MsgDeleteVestEntryResponse";
  value: MsgDeleteVestEntryResponseAmino;
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */
export interface MsgDeleteVestEntryResponseSDKType {}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */
export interface MsgSetVestEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vest entry to set. */
  entry: VestEntry;
}
export interface MsgSetVestEntryProtoMsg {
  typeUrl: "/dydxprotocol.vest.MsgSetVestEntry";
  value: Uint8Array;
}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */
export interface MsgSetVestEntryAmino {
  /** authority is the address that controls the module. */
  authority?: string;
  /** The vest entry to set. */
  entry?: VestEntryAmino;
}
export interface MsgSetVestEntryAminoMsg {
  type: "/dydxprotocol.vest.MsgSetVestEntry";
  value: MsgSetVestEntryAmino;
}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */
export interface MsgSetVestEntrySDKType {
  authority: string;
  entry: VestEntrySDKType;
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */
export interface MsgSetVestEntryResponse {}
export interface MsgSetVestEntryResponseProtoMsg {
  typeUrl: "/dydxprotocol.vest.MsgSetVestEntryResponse";
  value: Uint8Array;
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */
export interface MsgSetVestEntryResponseAmino {}
export interface MsgSetVestEntryResponseAminoMsg {
  type: "/dydxprotocol.vest.MsgSetVestEntryResponse";
  value: MsgSetVestEntryResponseAmino;
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */
export interface MsgSetVestEntryResponseSDKType {}
function createBaseMsgDeleteVestEntry(): MsgDeleteVestEntry {
  return {
    authority: "",
    vesterAccount: ""
  };
}
export const MsgDeleteVestEntry = {
  typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntry",
  encode(message: MsgDeleteVestEntry, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.vesterAccount !== "") {
      writer.uint32(18).string(message.vesterAccount);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeleteVestEntry {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteVestEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.vesterAccount = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgDeleteVestEntry>): MsgDeleteVestEntry {
    const message = createBaseMsgDeleteVestEntry();
    message.authority = object.authority ?? "";
    message.vesterAccount = object.vesterAccount ?? "";
    return message;
  },
  fromAmino(object: MsgDeleteVestEntryAmino): MsgDeleteVestEntry {
    const message = createBaseMsgDeleteVestEntry();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.vester_account !== undefined && object.vester_account !== null) {
      message.vesterAccount = object.vester_account;
    }
    return message;
  },
  toAmino(message: MsgDeleteVestEntry): MsgDeleteVestEntryAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.vester_account = message.vesterAccount;
    return obj;
  },
  fromAminoMsg(object: MsgDeleteVestEntryAminoMsg): MsgDeleteVestEntry {
    return MsgDeleteVestEntry.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeleteVestEntryProtoMsg): MsgDeleteVestEntry {
    return MsgDeleteVestEntry.decode(message.value);
  },
  toProto(message: MsgDeleteVestEntry): Uint8Array {
    return MsgDeleteVestEntry.encode(message).finish();
  },
  toProtoMsg(message: MsgDeleteVestEntry): MsgDeleteVestEntryProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntry",
      value: MsgDeleteVestEntry.encode(message).finish()
    };
  }
};
function createBaseMsgDeleteVestEntryResponse(): MsgDeleteVestEntryResponse {
  return {};
}
export const MsgDeleteVestEntryResponse = {
  typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntryResponse",
  encode(_: MsgDeleteVestEntryResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgDeleteVestEntryResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteVestEntryResponse();
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
  fromPartial(_: Partial<MsgDeleteVestEntryResponse>): MsgDeleteVestEntryResponse {
    const message = createBaseMsgDeleteVestEntryResponse();
    return message;
  },
  fromAmino(_: MsgDeleteVestEntryResponseAmino): MsgDeleteVestEntryResponse {
    const message = createBaseMsgDeleteVestEntryResponse();
    return message;
  },
  toAmino(_: MsgDeleteVestEntryResponse): MsgDeleteVestEntryResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgDeleteVestEntryResponseAminoMsg): MsgDeleteVestEntryResponse {
    return MsgDeleteVestEntryResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgDeleteVestEntryResponseProtoMsg): MsgDeleteVestEntryResponse {
    return MsgDeleteVestEntryResponse.decode(message.value);
  },
  toProto(message: MsgDeleteVestEntryResponse): Uint8Array {
    return MsgDeleteVestEntryResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgDeleteVestEntryResponse): MsgDeleteVestEntryResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.MsgDeleteVestEntryResponse",
      value: MsgDeleteVestEntryResponse.encode(message).finish()
    };
  }
};
function createBaseMsgSetVestEntry(): MsgSetVestEntry {
  return {
    authority: "",
    entry: VestEntry.fromPartial({})
  };
}
export const MsgSetVestEntry = {
  typeUrl: "/dydxprotocol.vest.MsgSetVestEntry",
  encode(message: MsgSetVestEntry, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }
    if (message.entry !== undefined) {
      VestEntry.encode(message.entry, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetVestEntry {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVestEntry();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;
        case 2:
          message.entry = VestEntry.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromPartial(object: Partial<MsgSetVestEntry>): MsgSetVestEntry {
    const message = createBaseMsgSetVestEntry();
    message.authority = object.authority ?? "";
    message.entry = object.entry !== undefined && object.entry !== null ? VestEntry.fromPartial(object.entry) : undefined;
    return message;
  },
  fromAmino(object: MsgSetVestEntryAmino): MsgSetVestEntry {
    const message = createBaseMsgSetVestEntry();
    if (object.authority !== undefined && object.authority !== null) {
      message.authority = object.authority;
    }
    if (object.entry !== undefined && object.entry !== null) {
      message.entry = VestEntry.fromAmino(object.entry);
    }
    return message;
  },
  toAmino(message: MsgSetVestEntry): MsgSetVestEntryAmino {
    const obj: any = {};
    obj.authority = message.authority;
    obj.entry = message.entry ? VestEntry.toAmino(message.entry) : undefined;
    return obj;
  },
  fromAminoMsg(object: MsgSetVestEntryAminoMsg): MsgSetVestEntry {
    return MsgSetVestEntry.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetVestEntryProtoMsg): MsgSetVestEntry {
    return MsgSetVestEntry.decode(message.value);
  },
  toProto(message: MsgSetVestEntry): Uint8Array {
    return MsgSetVestEntry.encode(message).finish();
  },
  toProtoMsg(message: MsgSetVestEntry): MsgSetVestEntryProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.MsgSetVestEntry",
      value: MsgSetVestEntry.encode(message).finish()
    };
  }
};
function createBaseMsgSetVestEntryResponse(): MsgSetVestEntryResponse {
  return {};
}
export const MsgSetVestEntryResponse = {
  typeUrl: "/dydxprotocol.vest.MsgSetVestEntryResponse",
  encode(_: MsgSetVestEntryResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): MsgSetVestEntryResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVestEntryResponse();
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
  fromPartial(_: Partial<MsgSetVestEntryResponse>): MsgSetVestEntryResponse {
    const message = createBaseMsgSetVestEntryResponse();
    return message;
  },
  fromAmino(_: MsgSetVestEntryResponseAmino): MsgSetVestEntryResponse {
    const message = createBaseMsgSetVestEntryResponse();
    return message;
  },
  toAmino(_: MsgSetVestEntryResponse): MsgSetVestEntryResponseAmino {
    const obj: any = {};
    return obj;
  },
  fromAminoMsg(object: MsgSetVestEntryResponseAminoMsg): MsgSetVestEntryResponse {
    return MsgSetVestEntryResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: MsgSetVestEntryResponseProtoMsg): MsgSetVestEntryResponse {
    return MsgSetVestEntryResponse.decode(message.value);
  },
  toProto(message: MsgSetVestEntryResponse): Uint8Array {
    return MsgSetVestEntryResponse.encode(message).finish();
  },
  toProtoMsg(message: MsgSetVestEntryResponse): MsgSetVestEntryResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.vest.MsgSetVestEntryResponse",
      value: MsgSetVestEntryResponse.encode(message).finish()
    };
  }
};