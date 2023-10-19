import { Params, ParamsSDKType } from "./params";
import { RewardShare, RewardShareSDKType } from "./reward_share";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequest {}
/** QueryParamsRequest is a request type for the Params RPC method. */

export interface QueryParamsRequestSDKType {}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponse {
  params?: Params;
}
/** QueryParamsResponse is a response type for the Params RPC method. */

export interface QueryParamsResponseSDKType {
  params?: ParamsSDKType;
}
/** QueryRewardShareRequest is a request type for the RewardShare RPC method. */

export interface QueryRewardShareRequest {
  address: string;
}
/** QueryRewardShareRequest is a request type for the RewardShare RPC method. */

export interface QueryRewardShareRequestSDKType {
  address: string;
}
/** QueryRewardShareResponse is a response type for the RewardsShare RPC method. */

export interface QueryRewardShareResponse {
  rewardShare?: RewardShare;
}
/** QueryRewardShareResponse is a response type for the RewardsShare RPC method. */

export interface QueryRewardShareResponseSDKType {
  reward_share?: RewardShareSDKType;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();

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

  fromPartial(_: DeepPartial<QueryParamsRequest>): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  }

};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return {
    params: undefined
  };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryParamsResponse>): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = object.params !== undefined && object.params !== null ? Params.fromPartial(object.params) : undefined;
    return message;
  }

};

function createBaseQueryRewardShareRequest(): QueryRewardShareRequest {
  return {
    address: ""
  };
}

export const QueryRewardShareRequest = {
  encode(message: QueryRewardShareRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.address !== "") {
      writer.uint32(10).string(message.address);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRewardShareRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRewardShareRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.address = reader.string();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryRewardShareRequest>): QueryRewardShareRequest {
    const message = createBaseQueryRewardShareRequest();
    message.address = object.address ?? "";
    return message;
  }

};

function createBaseQueryRewardShareResponse(): QueryRewardShareResponse {
  return {
    rewardShare: undefined
  };
}

export const QueryRewardShareResponse = {
  encode(message: QueryRewardShareResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.rewardShare !== undefined) {
      RewardShare.encode(message.rewardShare, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryRewardShareResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryRewardShareResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.rewardShare = RewardShare.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryRewardShareResponse>): QueryRewardShareResponse {
    const message = createBaseQueryRewardShareResponse();
    message.rewardShare = object.rewardShare !== undefined && object.rewardShare !== null ? RewardShare.fromPartial(object.rewardShare) : undefined;
    return message;
  }

};