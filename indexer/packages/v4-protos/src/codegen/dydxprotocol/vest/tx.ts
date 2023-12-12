import { VestEntry, VestEntrySDKType } from "./vest_entry";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */

export interface MsgDeleteVestEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vester account of the vest entry to delete. */

  vesterAccount: string;
}
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */

export interface MsgDeleteVestEntrySDKType {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vester account of the vest entry to delete. */

  vester_account: string;
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */

export interface MsgDeleteVestEntryResponse {}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */

export interface MsgDeleteVestEntryResponseSDKType {}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */

export interface MsgSetVestEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vest entry to set. */

  entry?: VestEntry;
}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */

export interface MsgSetVestEntrySDKType {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vest entry to set. */

  entry?: VestEntrySDKType;
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */

export interface MsgSetVestEntryResponse {}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */

export interface MsgSetVestEntryResponseSDKType {}

function createBaseMsgDeleteVestEntry(): MsgDeleteVestEntry {
  return {
    authority: "",
    vesterAccount: ""
  };
}

export const MsgDeleteVestEntry = {
  encode(message: MsgDeleteVestEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vesterAccount !== "") {
      writer.uint32(18).string(message.vesterAccount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<MsgDeleteVestEntry>): MsgDeleteVestEntry {
    const message = createBaseMsgDeleteVestEntry();
    message.authority = object.authority ?? "";
    message.vesterAccount = object.vesterAccount ?? "";
    return message;
  }

};

function createBaseMsgDeleteVestEntryResponse(): MsgDeleteVestEntryResponse {
  return {};
}

export const MsgDeleteVestEntryResponse = {
  encode(_: MsgDeleteVestEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestEntryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(_: DeepPartial<MsgDeleteVestEntryResponse>): MsgDeleteVestEntryResponse {
    const message = createBaseMsgDeleteVestEntryResponse();
    return message;
  }

};

function createBaseMsgSetVestEntry(): MsgSetVestEntry {
  return {
    authority: "",
    entry: undefined
  };
}

export const MsgSetVestEntry = {
  encode(message: MsgSetVestEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.entry !== undefined) {
      VestEntry.encode(message.entry, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(object: DeepPartial<MsgSetVestEntry>): MsgSetVestEntry {
    const message = createBaseMsgSetVestEntry();
    message.authority = object.authority ?? "";
    message.entry = object.entry !== undefined && object.entry !== null ? VestEntry.fromPartial(object.entry) : undefined;
    return message;
  }

};

function createBaseMsgSetVestEntryResponse(): MsgSetVestEntryResponse {
  return {};
}

export const MsgSetVestEntryResponse = {
  encode(_: MsgSetVestEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestEntryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
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

  fromPartial(_: DeepPartial<MsgSetVestEntryResponse>): MsgSetVestEntryResponse {
    const message = createBaseMsgSetVestEntryResponse();
    return message;
  }

};