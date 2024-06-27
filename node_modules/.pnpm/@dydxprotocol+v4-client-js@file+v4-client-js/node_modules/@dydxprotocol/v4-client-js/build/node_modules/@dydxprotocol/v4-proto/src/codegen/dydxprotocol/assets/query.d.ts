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
    id: number;
}
/** QueryAssetResponse is response type for the Asset RPC method. */
export interface QueryAssetResponse {
    /** QueryAssetResponse is response type for the Asset RPC method. */
    asset?: Asset;
}
/** QueryAssetResponse is response type for the Asset RPC method. */
export interface QueryAssetResponseSDKType {
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
export declare const QueryAssetRequest: {
    encode(message: QueryAssetRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAssetRequest;
    fromPartial(object: DeepPartial<QueryAssetRequest>): QueryAssetRequest;
};
export declare const QueryAssetResponse: {
    encode(message: QueryAssetResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAssetResponse;
    fromPartial(object: DeepPartial<QueryAssetResponse>): QueryAssetResponse;
};
export declare const QueryAllAssetsRequest: {
    encode(message: QueryAllAssetsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllAssetsRequest;
    fromPartial(object: DeepPartial<QueryAllAssetsRequest>): QueryAllAssetsRequest;
};
export declare const QueryAllAssetsResponse: {
    encode(message: QueryAllAssetsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllAssetsResponse;
    fromPartial(object: DeepPartial<QueryAllAssetsResponse>): QueryAllAssetsResponse;
};
