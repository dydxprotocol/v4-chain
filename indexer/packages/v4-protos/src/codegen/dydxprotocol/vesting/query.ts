import { VestingEntry, VestingEntrySDKType } from "./vesting_entry";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryVestingEntryRequest is a request type for the VestingEntry RPC method. */

export interface QueryVestingEntryRequest {
  /** QueryVestingEntryRequest is a request type for the VestingEntry RPC method. */
  vesterAccount: string;
}
/** QueryVestingEntryRequest is a request type for the VestingEntry RPC method. */

export interface QueryVestingEntryRequestSDKType {
  /** QueryVestingEntryRequest is a request type for the VestingEntry RPC method. */
  vester_account: string;
}
/** QueryVestingEntryResponse is a response type for the VestingEntry RPC method. */

export interface QueryVestingEntryResponse {
  /** QueryVestingEntryResponse is a response type for the VestingEntry RPC method. */
  entry?: VestingEntry;
}
/** QueryVestingEntryResponse is a response type for the VestingEntry RPC method. */

export interface QueryVestingEntryResponseSDKType {
  /** QueryVestingEntryResponse is a response type for the VestingEntry RPC method. */
  entry?: VestingEntrySDKType;
}

function createBaseQueryVestingEntryRequest(): QueryVestingEntryRequest {
  return {
    vesterAccount: ""
  };
}

export const QueryVestingEntryRequest = {
  encode(message: QueryVestingEntryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.vesterAccount !== "") {
      writer.uint32(10).string(message.vesterAccount);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVestingEntryRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVestingEntryRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.vesterAccount = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVestingEntryRequest>): QueryVestingEntryRequest {
    const message = createBaseQueryVestingEntryRequest();
    message.vesterAccount = object.vesterAccount ?? "";
    return message;
  }

};

function createBaseQueryVestingEntryResponse(): QueryVestingEntryResponse {
  return {
    entry: undefined
  };
}

export const QueryVestingEntryResponse = {
  encode(message: QueryVestingEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      VestingEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryVestingEntryResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryVestingEntryResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.entry = VestingEntry.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryVestingEntryResponse>): QueryVestingEntryResponse {
    const message = createBaseQueryVestingEntryResponse();
    message.entry = object.entry !== undefined && object.entry !== null ? VestingEntry.fromPartial(object.entry) : undefined;
    return message;
  }

};