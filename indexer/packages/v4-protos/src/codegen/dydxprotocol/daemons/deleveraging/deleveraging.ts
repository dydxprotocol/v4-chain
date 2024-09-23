import { SubaccountOpenPositionInfo, SubaccountOpenPositionInfoSDKType } from "../../clob/liquidations";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * UpdateSubaccountsListForDeleveragingDaemonRequest is a request message that contains a list of
 * perpetuals with the associated subaccounts that have open long and short positions
 */

export interface UpdateSubaccountsListForDeleveragingDaemonRequest {
  subaccountOpenPositionInfo: SubaccountOpenPositionInfo[];
}
/**
 * UpdateSubaccountsListForDeleveragingDaemonRequest is a request message that contains a list of
 * perpetuals with the associated subaccounts that have open long and short positions
 */

export interface UpdateSubaccountsListForDeleveragingDaemonRequestSDKType {
  subaccount_open_position_info: SubaccountOpenPositionInfoSDKType[];
}
/**
 * UpdateSubaccountsListForDeleveragingDaemonResponse is a response message for
 * UpdateSubaccountsListForDeleveragingDaemonRequest.
 */

export interface UpdateSubaccountsListForDeleveragingDaemonResponse {}
/**
 * UpdateSubaccountsListForDeleveragingDaemonResponse is a response message for
 * UpdateSubaccountsListForDeleveragingDaemonRequest.
 */

export interface UpdateSubaccountsListForDeleveragingDaemonResponseSDKType {}

function createBaseUpdateSubaccountsListForDeleveragingDaemonRequest(): UpdateSubaccountsListForDeleveragingDaemonRequest {
  return {
    subaccountOpenPositionInfo: []
  };
}

export const UpdateSubaccountsListForDeleveragingDaemonRequest = {
  encode(message: UpdateSubaccountsListForDeleveragingDaemonRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccountOpenPositionInfo) {
      SubaccountOpenPositionInfo.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSubaccountsListForDeleveragingDaemonRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSubaccountsListForDeleveragingDaemonRequest();

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

  fromPartial(object: DeepPartial<UpdateSubaccountsListForDeleveragingDaemonRequest>): UpdateSubaccountsListForDeleveragingDaemonRequest {
    const message = createBaseUpdateSubaccountsListForDeleveragingDaemonRequest();
    message.subaccountOpenPositionInfo = object.subaccountOpenPositionInfo?.map(e => SubaccountOpenPositionInfo.fromPartial(e)) || [];
    return message;
  }

};

function createBaseUpdateSubaccountsListForDeleveragingDaemonResponse(): UpdateSubaccountsListForDeleveragingDaemonResponse {
  return {};
}

export const UpdateSubaccountsListForDeleveragingDaemonResponse = {
  encode(_: UpdateSubaccountsListForDeleveragingDaemonResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSubaccountsListForDeleveragingDaemonResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSubaccountsListForDeleveragingDaemonResponse();

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

  fromPartial(_: DeepPartial<UpdateSubaccountsListForDeleveragingDaemonResponse>): UpdateSubaccountsListForDeleveragingDaemonResponse {
    const message = createBaseUpdateSubaccountsListForDeleveragingDaemonResponse();
    return message;
  }

};