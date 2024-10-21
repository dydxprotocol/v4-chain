import { EventParams, EventParamsSDKType, ProposeParams, ProposeParamsSDKType, SafetyParams, SafetyParamsSDKType } from "./params";
import { BridgeEventInfo, BridgeEventInfoSDKType } from "./bridge_event_info";
import { MsgCompleteBridge, MsgCompleteBridgeSDKType } from "./tx";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */
export interface QueryEventParamsRequest {
}
/** QueryEventParamsRequest is a request type for the EventParams RPC method. */
export interface QueryEventParamsRequestSDKType {
}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */
export interface QueryEventParamsResponse {
    params?: EventParams;
}
/** QueryEventParamsResponse is a response type for the EventParams RPC method. */
export interface QueryEventParamsResponseSDKType {
    params?: EventParamsSDKType;
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */
export interface QueryProposeParamsRequest {
}
/** QueryProposeParamsRequest is a request type for the ProposeParams RPC method. */
export interface QueryProposeParamsRequestSDKType {
}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */
export interface QueryProposeParamsResponse {
    params?: ProposeParams;
}
/**
 * QueryProposeParamsResponse is a response type for the ProposeParams RPC
 * method.
 */
export interface QueryProposeParamsResponseSDKType {
    params?: ProposeParamsSDKType;
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */
export interface QuerySafetyParamsRequest {
}
/** QuerySafetyParamsRequest is a request type for the SafetyParams RPC method. */
export interface QuerySafetyParamsRequestSDKType {
}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */
export interface QuerySafetyParamsResponse {
    params?: SafetyParams;
}
/** QuerySafetyParamsResponse is a response type for the SafetyParams RPC method. */
export interface QuerySafetyParamsResponseSDKType {
    params?: SafetyParamsSDKType;
}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoRequest {
}
/**
 * QueryAcknowledgedEventInfoRequest is a request type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoRequestSDKType {
}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoResponse {
    info?: BridgeEventInfo;
}
/**
 * QueryAcknowledgedEventInfoResponse is a response type for the
 * AcknowledgedEventInfo RPC method.
 */
export interface QueryAcknowledgedEventInfoResponseSDKType {
    info?: BridgeEventInfoSDKType;
}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoRequest {
}
/**
 * QueryRecognizedEventInfoRequest is a request type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoRequestSDKType {
}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoResponse {
    info?: BridgeEventInfo;
}
/**
 * QueryRecognizedEventInfoResponse is a response type for the
 * RecognizedEventInfo RPC method.
 */
export interface QueryRecognizedEventInfoResponseSDKType {
    info?: BridgeEventInfoSDKType;
}
/**
 * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesRequest {
    /**
     * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
     * DelayedCompleteBridgeMessages RPC method.
     */
    address: string;
}
/**
 * QueryDelayedCompleteBridgeMessagesRequest is a request type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesRequestSDKType {
    address: string;
}
/**
 * QueryDelayedCompleteBridgeMessagesResponse is a response type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesResponse {
    messages: DelayedCompleteBridgeMessage[];
}
/**
 * QueryDelayedCompleteBridgeMessagesResponse is a response type for the
 * DelayedCompleteBridgeMessages RPC method.
 */
export interface QueryDelayedCompleteBridgeMessagesResponseSDKType {
    messages: DelayedCompleteBridgeMessageSDKType[];
}
/**
 * DelayedCompleteBridgeMessage is a message type for the response of
 * DelayedCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */
export interface DelayedCompleteBridgeMessage {
    message?: MsgCompleteBridge;
    blockHeight: number;
}
/**
 * DelayedCompleteBridgeMessage is a message type for the response of
 * DelayedCompleteBridgeMessages RPC method. It contains the message
 * and the block height at which it will execute.
 */
export interface DelayedCompleteBridgeMessageSDKType {
    message?: MsgCompleteBridgeSDKType;
    block_height: number;
}
export declare const QueryEventParamsRequest: {
    encode(_: QueryEventParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryEventParamsRequest;
    fromPartial(_: DeepPartial<QueryEventParamsRequest>): QueryEventParamsRequest;
};
export declare const QueryEventParamsResponse: {
    encode(message: QueryEventParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryEventParamsResponse;
    fromPartial(object: DeepPartial<QueryEventParamsResponse>): QueryEventParamsResponse;
};
export declare const QueryProposeParamsRequest: {
    encode(_: QueryProposeParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryProposeParamsRequest;
    fromPartial(_: DeepPartial<QueryProposeParamsRequest>): QueryProposeParamsRequest;
};
export declare const QueryProposeParamsResponse: {
    encode(message: QueryProposeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryProposeParamsResponse;
    fromPartial(object: DeepPartial<QueryProposeParamsResponse>): QueryProposeParamsResponse;
};
export declare const QuerySafetyParamsRequest: {
    encode(_: QuerySafetyParamsRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QuerySafetyParamsRequest;
    fromPartial(_: DeepPartial<QuerySafetyParamsRequest>): QuerySafetyParamsRequest;
};
export declare const QuerySafetyParamsResponse: {
    encode(message: QuerySafetyParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QuerySafetyParamsResponse;
    fromPartial(object: DeepPartial<QuerySafetyParamsResponse>): QuerySafetyParamsResponse;
};
export declare const QueryAcknowledgedEventInfoRequest: {
    encode(_: QueryAcknowledgedEventInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAcknowledgedEventInfoRequest;
    fromPartial(_: DeepPartial<QueryAcknowledgedEventInfoRequest>): QueryAcknowledgedEventInfoRequest;
};
export declare const QueryAcknowledgedEventInfoResponse: {
    encode(message: QueryAcknowledgedEventInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryAcknowledgedEventInfoResponse;
    fromPartial(object: DeepPartial<QueryAcknowledgedEventInfoResponse>): QueryAcknowledgedEventInfoResponse;
};
export declare const QueryRecognizedEventInfoRequest: {
    encode(_: QueryRecognizedEventInfoRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecognizedEventInfoRequest;
    fromPartial(_: DeepPartial<QueryRecognizedEventInfoRequest>): QueryRecognizedEventInfoRequest;
};
export declare const QueryRecognizedEventInfoResponse: {
    encode(message: QueryRecognizedEventInfoResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryRecognizedEventInfoResponse;
    fromPartial(object: DeepPartial<QueryRecognizedEventInfoResponse>): QueryRecognizedEventInfoResponse;
};
export declare const QueryDelayedCompleteBridgeMessagesRequest: {
    encode(message: QueryDelayedCompleteBridgeMessagesRequest, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelayedCompleteBridgeMessagesRequest;
    fromPartial(object: DeepPartial<QueryDelayedCompleteBridgeMessagesRequest>): QueryDelayedCompleteBridgeMessagesRequest;
};
export declare const QueryDelayedCompleteBridgeMessagesResponse: {
    encode(message: QueryDelayedCompleteBridgeMessagesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): QueryDelayedCompleteBridgeMessagesResponse;
    fromPartial(object: DeepPartial<QueryDelayedCompleteBridgeMessagesResponse>): QueryDelayedCompleteBridgeMessagesResponse;
};
export declare const DelayedCompleteBridgeMessage: {
    encode(message: DelayedCompleteBridgeMessage, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): DelayedCompleteBridgeMessage;
    fromPartial(object: DeepPartial<DelayedCompleteBridgeMessage>): DelayedCompleteBridgeMessage;
};
