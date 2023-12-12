import { LimitParams, LimitParamsSDKType } from "./limit_params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** ListLimitParamsRequest is a request type of the ListLimitParams RPC method. */

export interface ListLimitParamsRequest {}
/** ListLimitParamsRequest is a request type of the ListLimitParams RPC method. */

export interface ListLimitParamsRequestSDKType {}
/** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */

export interface ListLimitParamsResponse {
  /** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */
  limitParamsList: LimitParams[];
}
/** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */

export interface ListLimitParamsResponseSDKType {
  /** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */
  limit_params_list: LimitParamsSDKType[];
}
/**
 * QueryCapacityByDenomRequest is a request type for the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomRequest {
  /**
   * QueryCapacityByDenomRequest is a request type for the CapacityByDenom RPC
   * method.
   */
  denom: string;
}
/**
 * QueryCapacityByDenomRequest is a request type for the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomRequestSDKType {
  /**
   * QueryCapacityByDenomRequest is a request type for the CapacityByDenom RPC
   * method.
   */
  denom: string;
}
/** CapacityResult is a specific rate limit for a denom. */

export interface CapacityResult {
  periodSec: number;
  capacity: Uint8Array;
}
/** CapacityResult is a specific rate limit for a denom. */

export interface CapacityResultSDKType {
  period_sec: number;
  capacity: Uint8Array;
}
/**
 * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomResponse {
  /**
   * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
   * method.
   */
  results: CapacityResult[];
}
/**
 * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomResponseSDKType {
  /**
   * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
   * method.
   */
  results: CapacityResultSDKType[];
}

function createBaseListLimitParamsRequest(): ListLimitParamsRequest {
  return {};
}

export const ListLimitParamsRequest = {
  encode(_: ListLimitParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLimitParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLimitParamsRequest();

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

  fromPartial(_: DeepPartial<ListLimitParamsRequest>): ListLimitParamsRequest {
    const message = createBaseListLimitParamsRequest();
    return message;
  }

};

function createBaseListLimitParamsResponse(): ListLimitParamsResponse {
  return {
    limitParamsList: []
  };
}

export const ListLimitParamsResponse = {
  encode(message: ListLimitParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.limitParamsList) {
      LimitParams.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListLimitParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListLimitParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.limitParamsList.push(LimitParams.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<ListLimitParamsResponse>): ListLimitParamsResponse {
    const message = createBaseListLimitParamsResponse();
    message.limitParamsList = object.limitParamsList?.map(e => LimitParams.fromPartial(e)) || [];
    return message;
  }

};

function createBaseQueryCapacityByDenomRequest(): QueryCapacityByDenomRequest {
  return {
    denom: ""
  };
}

export const QueryCapacityByDenomRequest = {
  encode(message: QueryCapacityByDenomRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.denom !== "") {
      writer.uint32(10).string(message.denom);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCapacityByDenomRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCapacityByDenomRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.denom = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryCapacityByDenomRequest>): QueryCapacityByDenomRequest {
    const message = createBaseQueryCapacityByDenomRequest();
    message.denom = object.denom ?? "";
    return message;
  }

};

function createBaseCapacityResult(): CapacityResult {
  return {
    periodSec: 0,
    capacity: new Uint8Array()
  };
}

export const CapacityResult = {
  encode(message: CapacityResult, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.periodSec !== 0) {
      writer.uint32(8).uint32(message.periodSec);
    }

    if (message.capacity.length !== 0) {
      writer.uint32(18).bytes(message.capacity);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CapacityResult {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCapacityResult();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.periodSec = reader.uint32();
          break;

        case 2:
          message.capacity = reader.bytes();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<CapacityResult>): CapacityResult {
    const message = createBaseCapacityResult();
    message.periodSec = object.periodSec ?? 0;
    message.capacity = object.capacity ?? new Uint8Array();
    return message;
  }

};

function createBaseQueryCapacityByDenomResponse(): QueryCapacityByDenomResponse {
  return {
    results: []
  };
}

export const QueryCapacityByDenomResponse = {
  encode(message: QueryCapacityByDenomResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.results) {
      CapacityResult.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryCapacityByDenomResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryCapacityByDenomResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.results.push(CapacityResult.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryCapacityByDenomResponse>): QueryCapacityByDenomResponse {
    const message = createBaseQueryCapacityByDenomResponse();
    message.results = object.results?.map(e => CapacityResult.fromPartial(e)) || [];
    return message;
  }

};