import { SubaccountId, SubaccountIdSDKType } from "../../subaccounts/subaccount";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * LiquidateSubaccountsRequest is a request message that contains a list of
 * subaccount ids that potentially need to be liquidated. The list of subaccount
 * ids should not contain duplicates. The application should re-verify these
 * subaccount ids against current state before liquidating their positions.
 */

export interface LiquidateSubaccountsRequest {
  subaccountIds: SubaccountId[];
}
/**
 * LiquidateSubaccountsRequest is a request message that contains a list of
 * subaccount ids that potentially need to be liquidated. The list of subaccount
 * ids should not contain duplicates. The application should re-verify these
 * subaccount ids against current state before liquidating their positions.
 */

export interface LiquidateSubaccountsRequestSDKType {
  subaccount_ids: SubaccountIdSDKType[];
}
/**
 * LiquidateSubaccountsResponse is a response message for
 * LiquidateSubaccountsRequest.
 */

export interface LiquidateSubaccountsResponse {}
/**
 * LiquidateSubaccountsResponse is a response message for
 * LiquidateSubaccountsRequest.
 */

export interface LiquidateSubaccountsResponseSDKType {}

function createBaseLiquidateSubaccountsRequest(): LiquidateSubaccountsRequest {
  return {
    subaccountIds: []
  };
}

export const LiquidateSubaccountsRequest = {
  encode(message: LiquidateSubaccountsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.subaccountIds) {
      SubaccountId.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidateSubaccountsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidateSubaccountsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.subaccountIds.push(SubaccountId.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<LiquidateSubaccountsRequest>): LiquidateSubaccountsRequest {
    const message = createBaseLiquidateSubaccountsRequest();
    message.subaccountIds = object.subaccountIds?.map(e => SubaccountId.fromPartial(e)) || [];
    return message;
  }

};

function createBaseLiquidateSubaccountsResponse(): LiquidateSubaccountsResponse {
  return {};
}

export const LiquidateSubaccountsResponse = {
  encode(_: LiquidateSubaccountsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): LiquidateSubaccountsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLiquidateSubaccountsResponse();

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

  fromPartial(_: DeepPartial<LiquidateSubaccountsResponse>): LiquidateSubaccountsResponse {
    const message = createBaseLiquidateSubaccountsResponse();
    return message;
  }

};