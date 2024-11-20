import { StreamOrderbookFill, StreamOrderbookFillSDKType, StreamOrderbookUpdate, StreamOrderbookUpdateSDKType } from "./query";
import { StreamSubaccountUpdate, StreamSubaccountUpdateSDKType } from "../subaccounts/streaming";
import { StreamPriceUpdate, StreamPriceUpdateSDKType } from "../prices/streaming";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** StagedFinalizeBlockEvent is an event staged during `FinalizeBlock`. */

export interface StagedFinalizeBlockEvent {
  orderFill?: StreamOrderbookFill;
  subaccountUpdate?: StreamSubaccountUpdate;
  orderbookUpdate?: StreamOrderbookUpdate;
  priceUpdate?: StreamPriceUpdate;
}
/** StagedFinalizeBlockEvent is an event staged during `FinalizeBlock`. */

export interface StagedFinalizeBlockEventSDKType {
  order_fill?: StreamOrderbookFillSDKType;
  subaccount_update?: StreamSubaccountUpdateSDKType;
  orderbook_update?: StreamOrderbookUpdateSDKType;
  price_update?: StreamPriceUpdateSDKType;
}

function createBaseStagedFinalizeBlockEvent(): StagedFinalizeBlockEvent {
  return {
    orderFill: undefined,
    subaccountUpdate: undefined,
    orderbookUpdate: undefined,
    priceUpdate: undefined
  };
}

export const StagedFinalizeBlockEvent = {
  encode(message: StagedFinalizeBlockEvent, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.orderFill !== undefined) {
      StreamOrderbookFill.encode(message.orderFill, writer.uint32(10).fork()).ldelim();
    }

    if (message.subaccountUpdate !== undefined) {
      StreamSubaccountUpdate.encode(message.subaccountUpdate, writer.uint32(18).fork()).ldelim();
    }

    if (message.orderbookUpdate !== undefined) {
      StreamOrderbookUpdate.encode(message.orderbookUpdate, writer.uint32(26).fork()).ldelim();
    }

    if (message.priceUpdate !== undefined) {
      StreamPriceUpdate.encode(message.priceUpdate, writer.uint32(34).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StagedFinalizeBlockEvent {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStagedFinalizeBlockEvent();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.orderFill = StreamOrderbookFill.decode(reader, reader.uint32());
          break;

        case 2:
          message.subaccountUpdate = StreamSubaccountUpdate.decode(reader, reader.uint32());
          break;

        case 3:
          message.orderbookUpdate = StreamOrderbookUpdate.decode(reader, reader.uint32());
          break;

        case 4:
          message.priceUpdate = StreamPriceUpdate.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<StagedFinalizeBlockEvent>): StagedFinalizeBlockEvent {
    const message = createBaseStagedFinalizeBlockEvent();
    message.orderFill = object.orderFill !== undefined && object.orderFill !== null ? StreamOrderbookFill.fromPartial(object.orderFill) : undefined;
    message.subaccountUpdate = object.subaccountUpdate !== undefined && object.subaccountUpdate !== null ? StreamSubaccountUpdate.fromPartial(object.subaccountUpdate) : undefined;
    message.orderbookUpdate = object.orderbookUpdate !== undefined && object.orderbookUpdate !== null ? StreamOrderbookUpdate.fromPartial(object.orderbookUpdate) : undefined;
    message.priceUpdate = object.priceUpdate !== undefined && object.priceUpdate !== null ? StreamPriceUpdate.fromPartial(object.priceUpdate) : undefined;
    return message;
  }

};