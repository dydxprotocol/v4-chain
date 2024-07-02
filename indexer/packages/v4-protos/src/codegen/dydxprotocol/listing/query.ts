import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Queries for the hard cap on listed markets */

export interface QueryMarketsHardCap {}
/** Queries for the hard cap on listed markets */

export interface QueryMarketsHardCapSDKType {}
/** Response type indicating the hard cap on listed markets */

export interface QueryMarketsHardCapResponse {
  /** Response type indicating the hard cap on listed markets */
  hardCap: number;
}
/** Response type indicating the hard cap on listed markets */

export interface QueryMarketsHardCapResponseSDKType {
  /** Response type indicating the hard cap on listed markets */
  hard_cap: number;
}

function createBaseQueryMarketsHardCap(): QueryMarketsHardCap {
  return {};
}

export const QueryMarketsHardCap = {
  encode(_: QueryMarketsHardCap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketsHardCap {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketsHardCap();

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

  fromPartial(_: DeepPartial<QueryMarketsHardCap>): QueryMarketsHardCap {
    const message = createBaseQueryMarketsHardCap();
    return message;
  }

};

function createBaseQueryMarketsHardCapResponse(): QueryMarketsHardCapResponse {
  return {
    hardCap: 0
  };
}

export const QueryMarketsHardCapResponse = {
  encode(message: QueryMarketsHardCapResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.hardCap !== 0) {
      writer.uint32(8).uint32(message.hardCap);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryMarketsHardCapResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryMarketsHardCapResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.hardCap = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryMarketsHardCapResponse>): QueryMarketsHardCapResponse {
    const message = createBaseQueryMarketsHardCapResponse();
    message.hardCap = object.hardCap ?? 0;
    return message;
  }

};