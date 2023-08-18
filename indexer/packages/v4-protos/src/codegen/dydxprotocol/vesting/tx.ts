import { VestingEntry, VestingEntrySDKType } from "./vesting_entry";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgDeleteVestingEntry is the Msg/DeleteVestingEntry request type. */

export interface MsgDeleteVestingEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vester account of the vesting entry to delete. */

  vesterAccount: string;
}
/** MsgDeleteVestingEntry is the Msg/DeleteVestingEntry request type. */

export interface MsgDeleteVestingEntrySDKType {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vester account of the vesting entry to delete. */

  vester_account: string;
}
/** MsgDeleteVestingEntryResponse is the Msg/DeleteVestingEntry response type. */

export interface MsgDeleteVestingEntryResponse {}
/** MsgDeleteVestingEntryResponse is the Msg/DeleteVestingEntry response type. */

export interface MsgDeleteVestingEntryResponseSDKType {}
/** MsgSetVestingEntry is the Msg/SetVestingEntry request type. */

export interface MsgSetVestingEntry {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vesting entry to set. */

  entry?: VestingEntry;
}
/** MsgSetVestingEntry is the Msg/SetVestingEntry request type. */

export interface MsgSetVestingEntrySDKType {
  /** authority is the address that controls the module. */
  authority: string;
  /** The vesting entry to set. */

  entry?: VestingEntrySDKType;
}
/** MsgSetVestingEntryResponse is the Msg/SetVestingEntry response type. */

export interface MsgSetVestingEntryResponse {}
/** MsgSetVestingEntryResponse is the Msg/SetVestingEntry response type. */

export interface MsgSetVestingEntryResponseSDKType {}

function createBaseMsgDeleteVestingEntry(): MsgDeleteVestingEntry {
  return {
    authority: "",
    vesterAccount: ""
  };
}

export const MsgDeleteVestingEntry = {
  encode(message: MsgDeleteVestingEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.vesterAccount !== "") {
      writer.uint32(18).string(message.vesterAccount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestingEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteVestingEntry();

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

  fromPartial(object: DeepPartial<MsgDeleteVestingEntry>): MsgDeleteVestingEntry {
    const message = createBaseMsgDeleteVestingEntry();
    message.authority = object.authority ?? "";
    message.vesterAccount = object.vesterAccount ?? "";
    return message;
  }

};

function createBaseMsgDeleteVestingEntryResponse(): MsgDeleteVestingEntryResponse {
  return {};
}

export const MsgDeleteVestingEntryResponse = {
  encode(_: MsgDeleteVestingEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestingEntryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgDeleteVestingEntryResponse();

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

  fromPartial(_: DeepPartial<MsgDeleteVestingEntryResponse>): MsgDeleteVestingEntryResponse {
    const message = createBaseMsgDeleteVestingEntryResponse();
    return message;
  }

};

function createBaseMsgSetVestingEntry(): MsgSetVestingEntry {
  return {
    authority: "",
    entry: undefined
  };
}

export const MsgSetVestingEntry = {
  encode(message: MsgSetVestingEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.entry !== undefined) {
      VestingEntry.encode(message.entry, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestingEntry {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVestingEntry();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.entry = VestingEntry.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetVestingEntry>): MsgSetVestingEntry {
    const message = createBaseMsgSetVestingEntry();
    message.authority = object.authority ?? "";
    message.entry = object.entry !== undefined && object.entry !== null ? VestingEntry.fromPartial(object.entry) : undefined;
    return message;
  }

};

function createBaseMsgSetVestingEntryResponse(): MsgSetVestingEntryResponse {
  return {};
}

export const MsgSetVestingEntryResponse = {
  encode(_: MsgSetVestingEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestingEntryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetVestingEntryResponse();

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

  fromPartial(_: DeepPartial<MsgSetVestingEntryResponse>): MsgSetVestingEntryResponse {
    const message = createBaseMsgSetVestingEntryResponse();
    return message;
  }

};