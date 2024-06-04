import { BridgeEvent, BridgeEventSDKType } from "../../bridge/bridge_event";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * AddBridgeEventsRequest is a request message that contains a list of new
 * bridge events. The events should be contiguous and sorted by (unique) id.
 */

export interface AddBridgeEventsRequest {
  bridgeEvents: BridgeEvent[];
}
/**
 * AddBridgeEventsRequest is a request message that contains a list of new
 * bridge events. The events should be contiguous and sorted by (unique) id.
 */

export interface AddBridgeEventsRequestSDKType {
  bridge_events: BridgeEventSDKType[];
}
/** AddBridgeEventsResponse is a response message for BridgeEventRequest. */

export interface AddBridgeEventsResponse {}
/** AddBridgeEventsResponse is a response message for BridgeEventRequest. */

export interface AddBridgeEventsResponseSDKType {}

function createBaseAddBridgeEventsRequest(): AddBridgeEventsRequest {
  return {
    bridgeEvents: []
  };
}

export const AddBridgeEventsRequest = {
  encode(message: AddBridgeEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.bridgeEvents) {
      BridgeEvent.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddBridgeEventsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddBridgeEventsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.bridgeEvents.push(BridgeEvent.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AddBridgeEventsRequest>): AddBridgeEventsRequest {
    const message = createBaseAddBridgeEventsRequest();
    message.bridgeEvents = object.bridgeEvents?.map(e => BridgeEvent.fromPartial(e)) || [];
    return message;
  }

};

function createBaseAddBridgeEventsResponse(): AddBridgeEventsResponse {
  return {};
}

export const AddBridgeEventsResponse = {
  encode(_: AddBridgeEventsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddBridgeEventsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddBridgeEventsResponse();

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

  fromPartial(_: DeepPartial<AddBridgeEventsResponse>): AddBridgeEventsResponse {
    const message = createBaseAddBridgeEventsResponse();
    return message;
  }

};