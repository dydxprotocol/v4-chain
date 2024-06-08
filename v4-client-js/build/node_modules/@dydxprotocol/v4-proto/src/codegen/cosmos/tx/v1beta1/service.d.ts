/// <reference types="long" />
import { Tx, TxSDKType } from "./tx";
import { PageRequest, PageRequestSDKType, PageResponse, PageResponseSDKType } from "../../base/query/v1beta1/pagination";
import { TxResponse, TxResponseSDKType, GasInfo, GasInfoSDKType, Result, ResultSDKType } from "../../base/abci/v1beta1/abci";
import { BlockID, BlockIDSDKType } from "../../../tendermint/types/types";
import { Block, BlockSDKType } from "../../../tendermint/types/block";
import { Long, DeepPartial } from "../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** OrderBy defines the sorting order */
export declare enum OrderBy {
    /**
     * ORDER_BY_UNSPECIFIED - ORDER_BY_UNSPECIFIED specifies an unknown sorting order. OrderBy defaults
     * to ASC in this case.
     */
    ORDER_BY_UNSPECIFIED = 0,
    /** ORDER_BY_ASC - ORDER_BY_ASC defines ascending order */
    ORDER_BY_ASC = 1,
    /** ORDER_BY_DESC - ORDER_BY_DESC defines descending order */
    ORDER_BY_DESC = 2,
    UNRECOGNIZED = -1
}
export declare const OrderBySDKType: typeof OrderBy;
export declare function orderByFromJSON(object: any): OrderBy;
export declare function orderByToJSON(object: OrderBy): string;
/**
 * BroadcastMode specifies the broadcast mode for the TxService.Broadcast RPC
 * method.
 */
export declare enum BroadcastMode {
    /** BROADCAST_MODE_UNSPECIFIED - zero-value for mode ordering */
    BROADCAST_MODE_UNSPECIFIED = 0,
    /**
     * BROADCAST_MODE_BLOCK - DEPRECATED: use BROADCAST_MODE_SYNC instead,
     * BROADCAST_MODE_BLOCK is not supported by the SDK from v0.47.x onwards.
     */
    BROADCAST_MODE_BLOCK = 1,
    /**
     * BROADCAST_MODE_SYNC - BROADCAST_MODE_SYNC defines a tx broadcasting mode where the client waits
     * for a CheckTx execution response only.
     */
    BROADCAST_MODE_SYNC = 2,
    /**
     * BROADCAST_MODE_ASYNC - BROADCAST_MODE_ASYNC defines a tx broadcasting mode where the client
     * returns immediately.
     */
    BROADCAST_MODE_ASYNC = 3,
    UNRECOGNIZED = -1
}
export declare const BroadcastModeSDKType: typeof BroadcastMode;
export declare function broadcastModeFromJSON(object: any): BroadcastMode;
export declare function broadcastModeToJSON(object: BroadcastMode): string;
/**
 * GetTxsEventRequest is the request type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventRequest {
    /**
     * events is the list of transaction event type.
     * Deprecated post v0.47.x: use query instead, which should contain a valid
     * events query.
     */
    /** @deprecated */
    events: string[];
    /**
     * pagination defines a pagination for the request.
     * Deprecated post v0.46.x: use page and limit instead.
     */
    /** @deprecated */
    pagination?: PageRequest;
    orderBy: OrderBy;
    /**
     * page is the page number to query, starts at 1. If not provided, will
     * default to first page.
     */
    page: Long;
    /**
     * limit is the total number of results to be returned in the result page.
     * If left empty it will default to a value to be set by each app.
     */
    limit: Long;
    /**
     * query defines the transaction event query that is proxied to Tendermint's
     * TxSearch RPC method. The query must be valid.
     *
     * Since cosmos-sdk 0.50
     */
    query: string;
}
/**
 * GetTxsEventRequest is the request type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventRequestSDKType {
    /** @deprecated */
    events: string[];
    /** @deprecated */
    pagination?: PageRequestSDKType;
    order_by: OrderBy;
    page: Long;
    limit: Long;
    query: string;
}
/**
 * GetTxsEventResponse is the response type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventResponse {
    /** txs is the list of queried transactions. */
    txs: Tx[];
    /** tx_responses is the list of queried TxResponses. */
    txResponses: TxResponse[];
    /**
     * pagination defines a pagination for the response.
     * Deprecated post v0.46.x: use total instead.
     */
    /** @deprecated */
    pagination?: PageResponse;
    /** total is total number of results available */
    total: Long;
}
/**
 * GetTxsEventResponse is the response type for the Service.TxsByEvents
 * RPC method.
 */
export interface GetTxsEventResponseSDKType {
    txs: TxSDKType[];
    tx_responses: TxResponseSDKType[];
    /** @deprecated */
    pagination?: PageResponseSDKType;
    total: Long;
}
/**
 * BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
 * RPC method.
 */
export interface BroadcastTxRequest {
    /** tx_bytes is the raw transaction. */
    txBytes: Uint8Array;
    mode: BroadcastMode;
}
/**
 * BroadcastTxRequest is the request type for the Service.BroadcastTxRequest
 * RPC method.
 */
export interface BroadcastTxRequestSDKType {
    tx_bytes: Uint8Array;
    mode: BroadcastMode;
}
/**
 * BroadcastTxResponse is the response type for the
 * Service.BroadcastTx method.
 */
export interface BroadcastTxResponse {
    /** tx_response is the queried TxResponses. */
    txResponse?: TxResponse;
}
/**
 * BroadcastTxResponse is the response type for the
 * Service.BroadcastTx method.
 */
export interface BroadcastTxResponseSDKType {
    tx_response?: TxResponseSDKType;
}
/**
 * SimulateRequest is the request type for the Service.Simulate
 * RPC method.
 */
export interface SimulateRequest {
    /**
     * tx is the transaction to simulate.
     * Deprecated. Send raw tx bytes instead.
     */
    /** @deprecated */
    tx?: Tx;
    /**
     * tx_bytes is the raw transaction.
     *
     * Since: cosmos-sdk 0.43
     */
    txBytes: Uint8Array;
}
/**
 * SimulateRequest is the request type for the Service.Simulate
 * RPC method.
 */
export interface SimulateRequestSDKType {
    /** @deprecated */
    tx?: TxSDKType;
    tx_bytes: Uint8Array;
}
/**
 * SimulateResponse is the response type for the
 * Service.SimulateRPC method.
 */
export interface SimulateResponse {
    /** gas_info is the information about gas used in the simulation. */
    gasInfo?: GasInfo;
    /** result is the result of the simulation. */
    result?: Result;
}
/**
 * SimulateResponse is the response type for the
 * Service.SimulateRPC method.
 */
export interface SimulateResponseSDKType {
    gas_info?: GasInfoSDKType;
    result?: ResultSDKType;
}
/**
 * GetTxRequest is the request type for the Service.GetTx
 * RPC method.
 */
export interface GetTxRequest {
    /** hash is the tx hash to query, encoded as a hex string. */
    hash: string;
}
/**
 * GetTxRequest is the request type for the Service.GetTx
 * RPC method.
 */
export interface GetTxRequestSDKType {
    hash: string;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxResponse {
    /** tx is the queried transaction. */
    tx?: Tx;
    /** tx_response is the queried TxResponses. */
    txResponse?: TxResponse;
}
/** GetTxResponse is the response type for the Service.GetTx method. */
export interface GetTxResponseSDKType {
    tx?: TxSDKType;
    tx_response?: TxResponseSDKType;
}
/**
 * GetBlockWithTxsRequest is the request type for the Service.GetBlockWithTxs
 * RPC method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsRequest {
    /** height is the height of the block to query. */
    height: Long;
    /** pagination defines a pagination for the request. */
    pagination?: PageRequest;
}
/**
 * GetBlockWithTxsRequest is the request type for the Service.GetBlockWithTxs
 * RPC method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsRequestSDKType {
    height: Long;
    pagination?: PageRequestSDKType;
}
/**
 * GetBlockWithTxsResponse is the response type for the Service.GetBlockWithTxs
 * method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsResponse {
    /** txs are the transactions in the block. */
    txs: Tx[];
    blockId?: BlockID;
    block?: Block;
    /** pagination defines a pagination for the response. */
    pagination?: PageResponse;
}
/**
 * GetBlockWithTxsResponse is the response type for the Service.GetBlockWithTxs
 * method.
 *
 * Since: cosmos-sdk 0.45.2
 */
export interface GetBlockWithTxsResponseSDKType {
    txs: TxSDKType[];
    block_id?: BlockIDSDKType;
    block?: BlockSDKType;
    pagination?: PageResponseSDKType;
}
/**
 * TxDecodeRequest is the request type for the Service.TxDecode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeRequest {
    /** tx_bytes is the raw transaction. */
    txBytes: Uint8Array;
}
/**
 * TxDecodeRequest is the request type for the Service.TxDecode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeRequestSDKType {
    tx_bytes: Uint8Array;
}
/**
 * TxDecodeResponse is the response type for the
 * Service.TxDecode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeResponse {
    /** tx is the decoded transaction. */
    tx?: Tx;
}
/**
 * TxDecodeResponse is the response type for the
 * Service.TxDecode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeResponseSDKType {
    tx?: TxSDKType;
}
/**
 * TxEncodeRequest is the request type for the Service.TxEncode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeRequest {
    /** tx is the transaction to encode. */
    tx?: Tx;
}
/**
 * TxEncodeRequest is the request type for the Service.TxEncode
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeRequestSDKType {
    tx?: TxSDKType;
}
/**
 * TxEncodeResponse is the response type for the
 * Service.TxEncode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeResponse {
    /** tx_bytes is the encoded transaction bytes. */
    txBytes: Uint8Array;
}
/**
 * TxEncodeResponse is the response type for the
 * Service.TxEncode method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeResponseSDKType {
    tx_bytes: Uint8Array;
}
/**
 * TxEncodeAminoRequest is the request type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoRequest {
    aminoJson: string;
}
/**
 * TxEncodeAminoRequest is the request type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoRequestSDKType {
    amino_json: string;
}
/**
 * TxEncodeAminoResponse is the response type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoResponse {
    aminoBinary: Uint8Array;
}
/**
 * TxEncodeAminoResponse is the response type for the Service.TxEncodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxEncodeAminoResponseSDKType {
    amino_binary: Uint8Array;
}
/**
 * TxDecodeAminoRequest is the request type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoRequest {
    aminoBinary: Uint8Array;
}
/**
 * TxDecodeAminoRequest is the request type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoRequestSDKType {
    amino_binary: Uint8Array;
}
/**
 * TxDecodeAminoResponse is the response type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoResponse {
    aminoJson: string;
}
/**
 * TxDecodeAminoResponse is the response type for the Service.TxDecodeAmino
 * RPC method.
 *
 * Since: cosmos-sdk 0.47
 */
export interface TxDecodeAminoResponseSDKType {
    amino_json: string;
}
export declare const GetTxsEventRequest: {
    encode(message: GetTxsEventRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetTxsEventRequest;
    fromPartial(object: DeepPartial<GetTxsEventRequest>): GetTxsEventRequest;
};
export declare const GetTxsEventResponse: {
    encode(message: GetTxsEventResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetTxsEventResponse;
    fromPartial(object: DeepPartial<GetTxsEventResponse>): GetTxsEventResponse;
};
export declare const BroadcastTxRequest: {
    encode(message: BroadcastTxRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BroadcastTxRequest;
    fromPartial(object: DeepPartial<BroadcastTxRequest>): BroadcastTxRequest;
};
export declare const BroadcastTxResponse: {
    encode(message: BroadcastTxResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): BroadcastTxResponse;
    fromPartial(object: DeepPartial<BroadcastTxResponse>): BroadcastTxResponse;
};
export declare const SimulateRequest: {
    encode(message: SimulateRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SimulateRequest;
    fromPartial(object: DeepPartial<SimulateRequest>): SimulateRequest;
};
export declare const SimulateResponse: {
    encode(message: SimulateResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): SimulateResponse;
    fromPartial(object: DeepPartial<SimulateResponse>): SimulateResponse;
};
export declare const GetTxRequest: {
    encode(message: GetTxRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetTxRequest;
    fromPartial(object: DeepPartial<GetTxRequest>): GetTxRequest;
};
export declare const GetTxResponse: {
    encode(message: GetTxResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetTxResponse;
    fromPartial(object: DeepPartial<GetTxResponse>): GetTxResponse;
};
export declare const GetBlockWithTxsRequest: {
    encode(message: GetBlockWithTxsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockWithTxsRequest;
    fromPartial(object: DeepPartial<GetBlockWithTxsRequest>): GetBlockWithTxsRequest;
};
export declare const GetBlockWithTxsResponse: {
    encode(message: GetBlockWithTxsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GetBlockWithTxsResponse;
    fromPartial(object: DeepPartial<GetBlockWithTxsResponse>): GetBlockWithTxsResponse;
};
export declare const TxDecodeRequest: {
    encode(message: TxDecodeRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxDecodeRequest;
    fromPartial(object: DeepPartial<TxDecodeRequest>): TxDecodeRequest;
};
export declare const TxDecodeResponse: {
    encode(message: TxDecodeResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxDecodeResponse;
    fromPartial(object: DeepPartial<TxDecodeResponse>): TxDecodeResponse;
};
export declare const TxEncodeRequest: {
    encode(message: TxEncodeRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxEncodeRequest;
    fromPartial(object: DeepPartial<TxEncodeRequest>): TxEncodeRequest;
};
export declare const TxEncodeResponse: {
    encode(message: TxEncodeResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxEncodeResponse;
    fromPartial(object: DeepPartial<TxEncodeResponse>): TxEncodeResponse;
};
export declare const TxEncodeAminoRequest: {
    encode(message: TxEncodeAminoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxEncodeAminoRequest;
    fromPartial(object: DeepPartial<TxEncodeAminoRequest>): TxEncodeAminoRequest;
};
export declare const TxEncodeAminoResponse: {
    encode(message: TxEncodeAminoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxEncodeAminoResponse;
    fromPartial(object: DeepPartial<TxEncodeAminoResponse>): TxEncodeAminoResponse;
};
export declare const TxDecodeAminoRequest: {
    encode(message: TxDecodeAminoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxDecodeAminoRequest;
    fromPartial(object: DeepPartial<TxDecodeAminoRequest>): TxDecodeAminoRequest;
};
export declare const TxDecodeAminoResponse: {
    encode(message: TxDecodeAminoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): TxDecodeAminoResponse;
    fromPartial(object: DeepPartial<TxDecodeAminoResponse>): TxDecodeAminoResponse;
};
