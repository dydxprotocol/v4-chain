import { PageRequest, PageRequestAmino, PageRequestSDKType, PageResponse, PageResponseAmino, PageResponseSDKType } from "../../cosmos/base/query/v1beta1/pagination";
import { Asset, AssetAmino, AssetSDKType } from "./asset";
import { BinaryReader, BinaryWriter } from "../../binary";
/** Queries an Asset by id. */
export interface QueryAssetRequest {
  /** Queries an Asset by id. */
  id: number;
}
export interface QueryAssetRequestProtoMsg {
  typeUrl: "/dydxprotocol.assets.QueryAssetRequest";
  value: Uint8Array;
}
/** Queries an Asset by id. */
export interface QueryAssetRequestAmino {
  /** Queries an Asset by id. */
  id?: number;
}
export interface QueryAssetRequestAminoMsg {
  type: "/dydxprotocol.assets.QueryAssetRequest";
  value: QueryAssetRequestAmino;
}
/** Queries an Asset by id. */
export interface QueryAssetRequestSDKType {
  id: number;
}
/** QueryAssetResponse is response type for the Asset RPC method. */
export interface QueryAssetResponse {
  /** QueryAssetResponse is response type for the Asset RPC method. */
  asset: Asset;
}
export interface QueryAssetResponseProtoMsg {
  typeUrl: "/dydxprotocol.assets.QueryAssetResponse";
  value: Uint8Array;
}
/** QueryAssetResponse is response type for the Asset RPC method. */
export interface QueryAssetResponseAmino {
  /** QueryAssetResponse is response type for the Asset RPC method. */
  asset?: AssetAmino;
}
export interface QueryAssetResponseAminoMsg {
  type: "/dydxprotocol.assets.QueryAssetResponse";
  value: QueryAssetResponseAmino;
}
/** QueryAssetResponse is response type for the Asset RPC method. */
export interface QueryAssetResponseSDKType {
  asset: AssetSDKType;
}
/** Queries a list of Asset items. */
export interface QueryAllAssetsRequest {
  pagination?: PageRequest;
}
export interface QueryAllAssetsRequestProtoMsg {
  typeUrl: "/dydxprotocol.assets.QueryAllAssetsRequest";
  value: Uint8Array;
}
/** Queries a list of Asset items. */
export interface QueryAllAssetsRequestAmino {
  pagination?: PageRequestAmino;
}
export interface QueryAllAssetsRequestAminoMsg {
  type: "/dydxprotocol.assets.QueryAllAssetsRequest";
  value: QueryAllAssetsRequestAmino;
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
export interface QueryAllAssetsResponseProtoMsg {
  typeUrl: "/dydxprotocol.assets.QueryAllAssetsResponse";
  value: Uint8Array;
}
/** QueryAllAssetsResponse is response type for the AllAssets RPC method. */
export interface QueryAllAssetsResponseAmino {
  asset?: AssetAmino[];
  pagination?: PageResponseAmino;
}
export interface QueryAllAssetsResponseAminoMsg {
  type: "/dydxprotocol.assets.QueryAllAssetsResponse";
  value: QueryAllAssetsResponseAmino;
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
  typeUrl: "/dydxprotocol.assets.QueryAssetRequest",
  encode(message: QueryAssetRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.id !== 0) {
      writer.uint32(8).uint32(message.id);
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAssetRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAssetRequest>): QueryAssetRequest {
    const message = createBaseQueryAssetRequest();
    message.id = object.id ?? 0;
    return message;
  },
  fromAmino(object: QueryAssetRequestAmino): QueryAssetRequest {
    const message = createBaseQueryAssetRequest();
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    }
    return message;
  },
  toAmino(message: QueryAssetRequest): QueryAssetRequestAmino {
    const obj: any = {};
    obj.id = message.id;
    return obj;
  },
  fromAminoMsg(object: QueryAssetRequestAminoMsg): QueryAssetRequest {
    return QueryAssetRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAssetRequestProtoMsg): QueryAssetRequest {
    return QueryAssetRequest.decode(message.value);
  },
  toProto(message: QueryAssetRequest): Uint8Array {
    return QueryAssetRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAssetRequest): QueryAssetRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.assets.QueryAssetRequest",
      value: QueryAssetRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAssetResponse(): QueryAssetResponse {
  return {
    asset: Asset.fromPartial({})
  };
}
export const QueryAssetResponse = {
  typeUrl: "/dydxprotocol.assets.QueryAssetResponse",
  encode(message: QueryAssetResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.asset !== undefined) {
      Asset.encode(message.asset, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAssetResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAssetResponse>): QueryAssetResponse {
    const message = createBaseQueryAssetResponse();
    message.asset = object.asset !== undefined && object.asset !== null ? Asset.fromPartial(object.asset) : undefined;
    return message;
  },
  fromAmino(object: QueryAssetResponseAmino): QueryAssetResponse {
    const message = createBaseQueryAssetResponse();
    if (object.asset !== undefined && object.asset !== null) {
      message.asset = Asset.fromAmino(object.asset);
    }
    return message;
  },
  toAmino(message: QueryAssetResponse): QueryAssetResponseAmino {
    const obj: any = {};
    obj.asset = message.asset ? Asset.toAmino(message.asset) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAssetResponseAminoMsg): QueryAssetResponse {
    return QueryAssetResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAssetResponseProtoMsg): QueryAssetResponse {
    return QueryAssetResponse.decode(message.value);
  },
  toProto(message: QueryAssetResponse): Uint8Array {
    return QueryAssetResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAssetResponse): QueryAssetResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.assets.QueryAssetResponse",
      value: QueryAssetResponse.encode(message).finish()
    };
  }
};
function createBaseQueryAllAssetsRequest(): QueryAllAssetsRequest {
  return {
    pagination: undefined
  };
}
export const QueryAllAssetsRequest = {
  typeUrl: "/dydxprotocol.assets.QueryAllAssetsRequest",
  encode(message: QueryAllAssetsRequest, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllAssetsRequest {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllAssetsRequest>): QueryAllAssetsRequest {
    const message = createBaseQueryAllAssetsRequest();
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageRequest.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllAssetsRequestAmino): QueryAllAssetsRequest {
    const message = createBaseQueryAllAssetsRequest();
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageRequest.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllAssetsRequest): QueryAllAssetsRequestAmino {
    const obj: any = {};
    obj.pagination = message.pagination ? PageRequest.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllAssetsRequestAminoMsg): QueryAllAssetsRequest {
    return QueryAllAssetsRequest.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllAssetsRequestProtoMsg): QueryAllAssetsRequest {
    return QueryAllAssetsRequest.decode(message.value);
  },
  toProto(message: QueryAllAssetsRequest): Uint8Array {
    return QueryAllAssetsRequest.encode(message).finish();
  },
  toProtoMsg(message: QueryAllAssetsRequest): QueryAllAssetsRequestProtoMsg {
    return {
      typeUrl: "/dydxprotocol.assets.QueryAllAssetsRequest",
      value: QueryAllAssetsRequest.encode(message).finish()
    };
  }
};
function createBaseQueryAllAssetsResponse(): QueryAllAssetsResponse {
  return {
    asset: [],
    pagination: undefined
  };
}
export const QueryAllAssetsResponse = {
  typeUrl: "/dydxprotocol.assets.QueryAllAssetsResponse",
  encode(message: QueryAllAssetsResponse, writer: BinaryWriter = BinaryWriter.create()): BinaryWriter {
    for (const v of message.asset) {
      Asset.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: BinaryReader | Uint8Array, length?: number): QueryAllAssetsResponse {
    const reader = input instanceof BinaryReader ? input : new BinaryReader(input);
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
  fromPartial(object: Partial<QueryAllAssetsResponse>): QueryAllAssetsResponse {
    const message = createBaseQueryAllAssetsResponse();
    message.asset = object.asset?.map(e => Asset.fromPartial(e)) || [];
    message.pagination = object.pagination !== undefined && object.pagination !== null ? PageResponse.fromPartial(object.pagination) : undefined;
    return message;
  },
  fromAmino(object: QueryAllAssetsResponseAmino): QueryAllAssetsResponse {
    const message = createBaseQueryAllAssetsResponse();
    message.asset = object.asset?.map(e => Asset.fromAmino(e)) || [];
    if (object.pagination !== undefined && object.pagination !== null) {
      message.pagination = PageResponse.fromAmino(object.pagination);
    }
    return message;
  },
  toAmino(message: QueryAllAssetsResponse): QueryAllAssetsResponseAmino {
    const obj: any = {};
    if (message.asset) {
      obj.asset = message.asset.map(e => e ? Asset.toAmino(e) : undefined);
    } else {
      obj.asset = [];
    }
    obj.pagination = message.pagination ? PageResponse.toAmino(message.pagination) : undefined;
    return obj;
  },
  fromAminoMsg(object: QueryAllAssetsResponseAminoMsg): QueryAllAssetsResponse {
    return QueryAllAssetsResponse.fromAmino(object.value);
  },
  fromProtoMsg(message: QueryAllAssetsResponseProtoMsg): QueryAllAssetsResponse {
    return QueryAllAssetsResponse.decode(message.value);
  },
  toProto(message: QueryAllAssetsResponse): Uint8Array {
    return QueryAllAssetsResponse.encode(message).finish();
  },
  toProtoMsg(message: QueryAllAssetsResponse): QueryAllAssetsResponseProtoMsg {
    return {
      typeUrl: "/dydxprotocol.assets.QueryAllAssetsResponse",
      value: QueryAllAssetsResponse.encode(message).finish()
    };
  }
};