import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * AddsDAIEventsRequest is a request message that contains a new
 * sDAI conversion rate.
 */

export interface AddsDAIEventRequest {
  /**
   * AddsDAIEventsRequest is a request message that contains a new
   * sDAI conversion rate.
   */
  conversionRate: string;
}
/**
 * AddsDAIEventsRequest is a request message that contains a new
 * sDAI conversion rate.
 */

export interface AddsDAIEventRequestSDKType {
  /**
   * AddsDAIEventsRequest is a request message that contains a new
   * sDAI conversion rate.
   */
  conversion_rate: string;
}
/** AddsDAIEventsResponse is a response message for AddsDAIEventsRequest. */

export interface AddsDAIEventResponse {}
/** AddsDAIEventsResponse is a response message for AddsDAIEventsRequest. */

export interface AddsDAIEventResponseSDKType {}

function createBaseAddsDAIEventRequest(): AddsDAIEventRequest {
  return {
    conversionRate: ""
  };
}

export const AddsDAIEventRequest = {
  encode(message: AddsDAIEventRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.conversionRate !== "") {
      writer.uint32(10).string(message.conversionRate);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddsDAIEventRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddsDAIEventRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.conversionRate = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<AddsDAIEventRequest>): AddsDAIEventRequest {
    const message = createBaseAddsDAIEventRequest();
    message.conversionRate = object.conversionRate ?? "";
    return message;
  }

};

function createBaseAddsDAIEventResponse(): AddsDAIEventResponse {
  return {};
}

export const AddsDAIEventResponse = {
  encode(_: AddsDAIEventResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddsDAIEventResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddsDAIEventResponse();

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

  fromPartial(_: DeepPartial<AddsDAIEventResponse>): AddsDAIEventResponse {
    const message = createBaseAddsDAIEventResponse();
    return message;
  }

};