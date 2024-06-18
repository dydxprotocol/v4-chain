import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * MsgSetMarketsHardCap is used to set a hard cap on the number of markets
 * listed
 */

export interface MsgSetMarketsHardCap {
  authority: string;
  /** Hard cap for the total number of markets listed */

  hardCapForMarkets: number;
}
/**
 * MsgSetMarketsHardCap is used to set a hard cap on the number of markets
 * listed
 */

export interface MsgSetMarketsHardCapSDKType {
  authority: string;
  /** Hard cap for the total number of markets listed */

  hard_cap_for_markets: number;
}
/** MsgSetMarketsHardCapResponse defines the MsgSetMarketsHardCap response */

export interface MsgSetMarketsHardCapResponse {}
/** MsgSetMarketsHardCapResponse defines the MsgSetMarketsHardCap response */

export interface MsgSetMarketsHardCapResponseSDKType {}

function createBaseMsgSetMarketsHardCap(): MsgSetMarketsHardCap {
  return {
    authority: "",
    hardCapForMarkets: 0
  };
}

export const MsgSetMarketsHardCap = {
  encode(message: MsgSetMarketsHardCap, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.hardCapForMarkets !== 0) {
      writer.uint32(16).uint32(message.hardCapForMarkets);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketsHardCap {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketsHardCap();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.hardCapForMarkets = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgSetMarketsHardCap>): MsgSetMarketsHardCap {
    const message = createBaseMsgSetMarketsHardCap();
    message.authority = object.authority ?? "";
    message.hardCapForMarkets = object.hardCapForMarkets ?? 0;
    return message;
  }

};

function createBaseMsgSetMarketsHardCapResponse(): MsgSetMarketsHardCapResponse {
  return {};
}

export const MsgSetMarketsHardCapResponse = {
  encode(_: MsgSetMarketsHardCapResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetMarketsHardCapResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgSetMarketsHardCapResponse();

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

  fromPartial(_: DeepPartial<MsgSetMarketsHardCapResponse>): MsgSetMarketsHardCapResponse {
    const message = createBaseMsgSetMarketsHardCapResponse();
    return message;
  }

};