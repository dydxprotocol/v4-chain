import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Asset, AssetSDKType } from "./asset";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** Queries an Asset by id. */

export interface QueryAssetRequest {
  /** Queries an Asset by id. */
  id: number;
}
/** Queries an Asset by id. */

export interface QueryAssetRequestSDKType {
  /** Queries an Asset by id. */
  id: number;
}
/** QueryAssetResponse is response type for the Asset RPC method. */

export interface QueryAssetResponse {
  /** QueryAssetResponse is response type for the Asset RPC method. */
  asset?: Asset;
}
/** QueryAssetResponse is response type for the Asset RPC method. */

export interface QueryAssetResponseSDKType {
  /** QueryAssetResponse is response type for the Asset RPC method. */
  asset?: AssetSDKType;
}
/** Queries a list of Asset items. */

export interface QueryAllAssetsRequest {
  pagination?: PageRequest;
}
/** Queries a list of Asset items. */

export interface QueryAllAssetsRequestSDKType {
  pagination?: PageRequestSDKType;
}
/** QueryAllAssetsResponse is response type for the AllAssets RPC method. */

export interface QueryAllAssetsResponse {
  asset: Asset[];
  pagination?: PageResponse;
}
/** QueryAllAssetsResponse is response type for the AllAssets RPC method. */

export interface QueryAllAssetsResponseSDKType {
  asset: AssetSDKType[];
  pagination?: PageResponseSDKType;
}

function createBaseQueryAssetRequest(): QueryAssetRequest {
  return {
    id: 0
  };
}

export const QueryAssetRequest = {
  encode(message: QueryAssetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAssetRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAssetRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.id = reader.uint32();
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAssetRequest>): QueryAssetRequest {
    const message = createBaseQueryAssetRequest();
    message.id = object.id ?? 0;
    return message;
  }

};

function createBaseQueryAssetResponse(): QueryAssetResponse {
  return {
    asset: undefined
  };
}

export const QueryAssetResponse = {
  encode(message: QueryAssetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.asset !== undefined) {
      Asset.encode(message.asset, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAssetResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAssetResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.asset = Asset.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAssetResponse>): QueryAssetResponse {
    const message = createBaseQueryAssetResponse();
    message.asset = object.asset !== undefined && object.asset !== null ? Asset.fromPartial(object.asset) : undefined;
    return message;
  }

};

function createBaseQueryAllAssetsRequest(): QueryAllAssetsRequest {
  return {
    pagination: undefined
  };
}

export const QueryAllAssetsRequest = {
  encode(message: QueryAllAssetsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllAssetsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllAssetsRequest();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllAssetsRequest>): QueryAllAssetsRequest {
    const message = createBaseQueryAllAssetsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  }

};

function createBaseQueryAllAssetsResponse(): QueryAllAssetsResponse {
  return {
    asset: [],
    pagination: undefined
  };
}

export const QueryAllAssetsResponse = {
  encode(message: QueryAllAssetsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.asset) {
      Asset.encode(v!, writer.uint32(10).fork()).ldelim();
    }

    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }

    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllAssetsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllAssetsResponse();

    while (reader.pos < end) {
      const tag = reader.uint32();

      switch (tag >>> 3) {
        case 1:
          message.asset.push(Asset.decode(reader, reader.uint32()));
          break;

        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;

        default:
          reader.skipType(tag & 7);
          break;
      }
    }

    return message;
  },

  fromPartial(object: DeepPartial<QueryAllAssetsResponse>): QueryAllAssetsResponse {
    const message = createBaseQueryAllAssetsResponse();
    message.asset = object.asset?.map(e => Asset.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  }

};