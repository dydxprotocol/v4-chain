import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgEnablePermissionlessMarketListing is used to enable/disable permissionless // market listing */

export interface MsgEnablePermissionlessMarketListing {
  authority: string;
  /** boolean flag to enable/disable permissionless market listing */

  enablePermissionlessMarketListing: boolean;
}
/** MsgEnablePermissionlessMarketListing is used to enable/disable permissionless // market listing */

export interface MsgEnablePermissionlessMarketListingSDKType {
  authority: string;
  /** boolean flag to enable/disable permissionless market listing */

  enablePermissionlessMarketListing: boolean;
}
export interface MsgEnablePermissionlessMarketListingResponse {}
export interface MsgEnablePermissionlessMarketListingResponseSDKType {}

function createBaseMsgEnablePermissionlessMarketListing(): MsgEnablePermissionlessMarketListing {
  return {
    authority: "",
    enablePermissionlessMarketListing: false
  };
}

export const MsgEnablePermissionlessMarketListing = {
  encode(message: MsgEnablePermissionlessMarketListing, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authority !== "") {
      writer.uint32(10).string(message.authority);
    }

    if (message.enablePermissionlessMarketListing === true) {
      writer.uint32(16).bool(message.enablePermissionlessMarketListing);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnablePermissionlessMarketListing {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgEnablePermissionlessMarketListing();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.authority = reader.string();
          break;

        case 2:
          message.enablePermissionlessMarketListing = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<MsgEnablePermissionlessMarketListing>): MsgEnablePermissionlessMarketListing {
    const message = createBaseMsgEnablePermissionlessMarketListing();
    message.authority = object.authority ?? "";
    message.enablePermissionlessMarketListing = object.enablePermissionlessMarketListing ?? false;
    return message;
  }

};

function createBaseMsgEnablePermissionlessMarketListingResponse(): MsgEnablePermissionlessMarketListingResponse {
  return {};
}

export const MsgEnablePermissionlessMarketListingResponse = {
  encode(_: MsgEnablePermissionlessMarketListingResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgEnablePermissionlessMarketListingResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgEnablePermissionlessMarketListingResponse();

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

  fromPartial(_: DeepPartial<MsgEnablePermissionlessMarketListingResponse>): MsgEnablePermissionlessMarketListingResponse {
    const message = createBaseMsgEnablePermissionlessMarketListingResponse();
    return message;
  }

};