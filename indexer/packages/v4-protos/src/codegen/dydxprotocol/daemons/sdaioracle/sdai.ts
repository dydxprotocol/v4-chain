import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/**
 * AddsDAIEventsRequest is a request message that contains a new
 * sDAI conversion rate.
 */

export interface AddsDAIEventsRequest {
  conversionRate: string;
}
/**
 * AddsDAIEventsRequest is a request message that contains a new
 * sDAI conversion rate.
 */

export interface AddsDAIEventsRequestSDKType {
  conversion_rate: string;
}
/** AddsDAIEventsResponse is a response message for AddsDAIEventsRequest. */

export interface AddsDAIEventsResponse {}
/** AddsDAIEventsResponse is a response message for AddsDAIEventsRequest. */

export interface AddsDAIEventsResponseSDKType {}

function createBaseAddsDAIEventsRequest(): AddsDAIEventsRequest {
  return {
    conversionRate: ""
  };
}

export const AddsDAIEventsRequest = {
  encode(message: AddsDAIEventsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.conversionRate !== "") {
      writer.uint32(10).string(message.conversionRate);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddsDAIEventsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddsDAIEventsRequest();

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

  fromPartial(object: DeepPartial<AddsDAIEventsRequest>): AddsDAIEventsRequest {
    const message = createBaseAddsDAIEventsRequest();
    message.conversionRate = object.conversionRate ?? "";
    return message;
  }

};

function createBaseAddsDAIEventsResponse(): AddsDAIEventsResponse {
  return {};
}

export const AddsDAIEventsResponse = {
  encode(_: AddsDAIEventsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddsDAIEventsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddsDAIEventsResponse();

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

  fromPartial(_: DeepPartial<AddsDAIEventsResponse>): AddsDAIEventsResponse {
    const message = createBaseAddsDAIEventsResponse();
    return message;
  }

};