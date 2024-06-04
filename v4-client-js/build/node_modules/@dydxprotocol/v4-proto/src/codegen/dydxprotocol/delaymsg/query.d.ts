import { DelayedMessage, DelayedMessageSDKType } from "./delayed_message";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdRequest {
}
/**
 * QueryNextDelayedMessageIdRequest is the request type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdRequestSDKType {
}
/**
 * QueryNextDelayedMessageIdResponse is the response type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdResponse {
    nextDelayedMessageId: number;
}
/**
 * QueryNextDelayedMessageIdResponse is the response type for the
 * NextDelayedMessageId RPC method.
 */
export interface QueryNextDelayedMessageIdResponseSDKType {
    next_delayed_message_id: number;
}
/** QueryMessageRequest is the request type for the Message RPC method. */
export interface QueryMessageRequest {
    /** QueryMessageRequest is the request type for the Message RPC method. */
    id: number;
}
/** QueryMessageRequest is the request type for the Message RPC method. */
export interface QueryMessageRequestSDKType {
    id: number;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */
export interface QueryMessageResponse {
    /** QueryGetMessageResponse is the response type for the Message RPC method. */
    message?: DelayedMessage;
}
/** QueryGetMessageResponse is the response type for the Message RPC method. */
export interface QueryMessageResponseSDKType {
    message?: DelayedMessageSDKType;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsRequest {
    /**
     * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
     * RPC method.
     */
    blockHeight: number;
}
/**
 * QueryBlockMessageIdsRequest is the request type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsRequestSDKType {
    block_height: number;
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsResponse {
    /**
     * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
     * RPC method.
     */
    messageIds: number[];
}
/**
 * QueryGetBlockMessageIdsResponse is the response type for the BlockMessageIds
 * RPC method.
 */
export interface QueryBlockMessageIdsResponseSDKType {
    message_ids: number[];
}
export declare const QueryNextDelayedMessageIdRequest: {
    encode(_: QueryNextDelayedMessageIdRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextDelayedMessageIdRequest;
    fromPartial(_: DeepPartial<QueryNextDelayedMessageIdRequest>): QueryNextDelayedMessageIdRequest;
};
export declare const QueryNextDelayedMessageIdResponse: {
    encode(message: QueryNextDelayedMessageIdResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryNextDelayedMessageIdResponse;
    fromPartial(object: DeepPartial<QueryNextDelayedMessageIdResponse>): QueryNextDelayedMessageIdResponse;
};
export declare const QueryMessageRequest: {
    encode(message: QueryMessageRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMessageRequest;
    fromPartial(object: DeepPartial<QueryMessageRequest>): QueryMessageRequest;
};
export declare const QueryMessageResponse: {
    encode(message: QueryMessageResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryMessageResponse;
    fromPartial(object: DeepPartial<QueryMessageResponse>): QueryMessageResponse;
};
export declare const QueryBlockMessageIdsRequest: {
    encode(message: QueryBlockMessageIdsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockMessageIdsRequest;
    fromPartial(object: DeepPartial<QueryBlockMessageIdsRequest>): QueryBlockMessageIdsRequest;
};
export declare const QueryBlockMessageIdsResponse: {
    encode(message: QueryBlockMessageIdsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryBlockMessageIdsResponse;
    fromPartial(object: DeepPartial<QueryBlockMessageIdsResponse>): QueryBlockMessageIdsResponse;
};
