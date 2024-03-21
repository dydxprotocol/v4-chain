import { LimitParams, LimitParamsSDKType } from "./limit_params";
import { LimiterCapacity, LimiterCapacitySDKType } from "./capacity";
import { PendingSendPacket, PendingSendPacketSDKType } from "./pending_send_packet";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** ListLimitParamsRequest is a request type of the ListLimitParams RPC method. */

export interface ListLimitParamsRequest {}
/** ListLimitParamsRequest is a request type of the ListLimitParams RPC method. */

export interface ListLimitParamsRequestSDKType {}
/** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */

export interface ListLimitParamsResponse {
  limitParamsList: LimitParams[];
}
/** ListLimitParamsResponse is a response type of the ListLimitParams RPC method. */

export interface ListLimitParamsResponseSDKType {
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
/**
 * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomResponse {
  limiterCapacityList: LimiterCapacity[];
}
/**
 * QueryCapacityByDenomResponse is a response type of the CapacityByDenom RPC
 * method.
 */

export interface QueryCapacityByDenomResponseSDKType {
  limiter_capacity_list: LimiterCapacitySDKType[];
}
/**
 * QueryAllPendingSendPacketsRequest is a request type for the
 * AllPendingSendPackets RPC
 */

export interface QueryAllPendingSendPacketsRequest {}
/**
 * QueryAllPendingSendPacketsRequest is a request type for the
 * AllPendingSendPackets RPC
 */

export interface QueryAllPendingSendPacketsRequestSDKType {}
/**
 * QueryAllPendingSendPacketsResponse is a response type of the
 * AllPendingSendPackets RPC
 */

export interface QueryAllPendingSendPacketsResponse {
  pendingSendPackets: PendingSendPacket[];
}
/**
 * QueryAllPendingSendPacketsResponse is a response type of the
 * AllPendingSendPackets RPC
 */

export interface QueryAllPendingSendPacketsResponseSDKType {
  pending_send_packets: PendingSendPacketSDKType[];
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

function createBaseQueryCapacityByDenomResponse(): QueryCapacityByDenomResponse {
  return {
    limiterCapacityList: []
  };
}

export const QueryCapacityByDenomResponse = {
  encode(message: QueryCapacityByDenomResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.limiterCapacityList) {
      LimiterCapacity.encode(v!, writer.uint32(10).fork()).ldelim();
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
          message.limiterCapacityList.push(LimiterCapacity.decode(reader, reader.uint32()));
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
    message.limiterCapacityList = object.limiterCapacityList?.map(e => LimiterCapacity.fromPartial(e)) || [];
    return message;
  }

};

function createBaseQueryAllPendingSendPacketsRequest(): QueryAllPendingSendPacketsRequest {
  return {};
}

export const QueryAllPendingSendPacketsRequest = {
  encode(_: QueryAllPendingSendPacketsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPendingSendPacketsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPendingSendPacketsRequest();

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

  fromPartial(_: DeepPartial<QueryAllPendingSendPacketsRequest>): QueryAllPendingSendPacketsRequest {
    const message = createBaseQueryAllPendingSendPacketsRequest();
    return message;
  }

};

function createBaseQueryAllPendingSendPacketsResponse(): QueryAllPendingSendPacketsResponse {
  return {
    pendingSendPackets: []
  };
}

export const QueryAllPendingSendPacketsResponse = {
  encode(message: QueryAllPendingSendPacketsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.pendingSendPackets) {
      PendingSendPacket.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllPendingSendPacketsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllPendingSendPacketsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pendingSendPackets.push(PendingSendPacket.decode(reader, reader.uint32()));
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllPendingSendPacketsResponse>): QueryAllPendingSendPacketsResponse {
    const message = createBaseQueryAllPendingSendPacketsResponse();
    message.pendingSendPackets = object.pendingSendPackets?.map(e => PendingSendPacket.fromPartial(e)) || [];
    return message;
  }

};