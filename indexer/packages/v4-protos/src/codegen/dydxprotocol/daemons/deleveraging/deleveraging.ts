import { SubaccountOpenPositionInfo, SubaccountOpenPositionInfoSDKType } from "../../clob/liquidations";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * DeleveragingSubaccountsRequest is a request message that contains a list of
 * perpetuals with the associated subaccounts that have open long and short positions
 */

export interface DeleveragingSubaccountsRequest {
  subaccountOpenPositionInfo: SubaccountOpenPositionInfo[];
}
/**
 * DeleveragingSubaccountsRequest is a request message that contains a list of
 * perpetuals with the associated subaccounts that have open long and short positions
 */

export interface DeleveragingSubaccountsRequestSDKType {
  subaccount_open_position_info: SubaccountOpenPositionInfoSDKType[];
}
/**
 * DeleveragingSubaccountsResponse is a response message for
 * DeleverageSubaccountsRequest.
 */

export interface DeleveragingSubaccountsResponse {}
/**
 * DeleveragingSubaccountsResponse is a response message for
 * DeleverageSubaccountsRequest.
 */

export interface DeleveragingSubaccountsResponseSDKType {}

function createBaseDeleveragingSubaccountsRequest(): DeleveragingSubaccountsRequest {
  return {
    subaccountOpenPositionInfo: []
  };
}

export const DeleveragingSubaccountsRequest = {
  encode(message: DeleveragingSubaccountsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccountOpenPositionInfo) {
      SubaccountOpenPositionInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleveragingSubaccountsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleveragingSubaccountsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountOpenPositionInfo.push(SubaccountOpenPositionInfo.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<DeleveragingSubaccountsRequest>): DeleveragingSubaccountsRequest {
    const message = createBaseDeleveragingSubaccountsRequest();
    message.subaccountOpenPositionInfo = object.subaccountOpenPositionInfo?.map(e => SubaccountOpenPositionInfo.fromPartial(e)) || [];
    return message;
  }

};

function createBaseDeleveragingSubaccountsResponse(): DeleveragingSubaccountsResponse {
  return {};
}

export const DeleveragingSubaccountsResponse = {
  encode(_: DeleveragingSubaccountsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleveragingSubaccountsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleveragingSubaccountsResponse();

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

  fromPartial(_: DeepPartial<DeleveragingSubaccountsResponse>): DeleveragingSubaccountsResponse {
    const message = createBaseDeleveragingSubaccountsResponse();
    return message;
  }

};