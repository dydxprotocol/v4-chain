import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Queries if permissionless listings are enabled */

export interface QueryPermissionlessMarketListingStatus {}
/** Queries if permissionless listings are enabled */

export interface QueryPermissionlessMarketListingStatusSDKType {}
/** Response type indicating if permissionless listings are enabled */

export interface QueryPermissionlessMarketListingStatusResponse {
  /** Response type indicating if permissionless listings are enabled */
  enabled: boolean;
}
/** Response type indicating if permissionless listings are enabled */

export interface QueryPermissionlessMarketListingStatusResponseSDKType {
  /** Response type indicating if permissionless listings are enabled */
  enabled: boolean;
}

function createBaseQueryPermissionlessMarketListingStatus(): QueryPermissionlessMarketListingStatus {
  return {};
}

export const QueryPermissionlessMarketListingStatus = {
  encode(_: QueryPermissionlessMarketListingStatus, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPermissionlessMarketListingStatus {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPermissionlessMarketListingStatus();

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

  fromPartial(_: DeepPartial<QueryPermissionlessMarketListingStatus>): QueryPermissionlessMarketListingStatus {
    const message = createBaseQueryPermissionlessMarketListingStatus();
    return message;
  }

};

function createBaseQueryPermissionlessMarketListingStatusResponse(): QueryPermissionlessMarketListingStatusResponse {
  return {
    enabled: false
  };
}

export const QueryPermissionlessMarketListingStatusResponse = {
  encode(message: QueryPermissionlessMarketListingStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryPermissionlessMarketListingStatusResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryPermissionlessMarketListingStatusResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.enabled = reader.bool();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryPermissionlessMarketListingStatusResponse>): QueryPermissionlessMarketListingStatusResponse {
    const message = createBaseQueryPermissionlessMarketListingStatusResponse();
    message.enabled = object.enabled ?? false;
    return message;
  }

};