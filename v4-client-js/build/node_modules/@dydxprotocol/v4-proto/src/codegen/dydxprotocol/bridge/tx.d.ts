import { BridgeEvent, BridgeEventSDKType } from "./bridge_event";
import { EventParams, EventParamsSDKType, ProposeParams, ProposeParamsSDKType, SafetyParams, SafetyParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgAcknowledgeBridges is the Msg/AcknowledgeBridges request type. */
export interface MsgAcknowledgeBridges {
    /** The events to acknowledge. */
    events: BridgeEvent[];
}
/** MsgAcknowledgeBridges is the Msg/AcknowledgeBridges request type. */
export interface MsgAcknowledgeBridgesSDKType {
    events: BridgeEventSDKType[];
}
/**
 * MsgAcknowledgeBridgesResponse is the Msg/AcknowledgeBridgesResponse response
 * type.
 */
export interface MsgAcknowledgeBridgesResponse {
}
/**
 * MsgAcknowledgeBridgesResponse is the Msg/AcknowledgeBridgesResponse response
 * type.
 */
export interface MsgAcknowledgeBridgesResponseSDKType {
}
/** MsgCompleteBridge is the Msg/CompleteBridgeResponse request type. */
export interface MsgCompleteBridge {
    authority: string;
    /** The event to complete. */
    event?: BridgeEvent;
}
/** MsgCompleteBridge is the Msg/CompleteBridgeResponse request type. */
export interface MsgCompleteBridgeSDKType {
    authority: string;
    event?: BridgeEventSDKType;
}
/** MsgCompleteBridgeResponse is the Msg/CompleteBridgeResponse response type. */
export interface MsgCompleteBridgeResponse {
}
/** MsgCompleteBridgeResponse is the Msg/CompleteBridgeResponse response type. */
export interface MsgCompleteBridgeResponseSDKType {
}
/** MsgUpdateEventParams is the Msg/UpdateEventParams request type. */
export interface MsgUpdateEventParams {
    authority: string;
    /** The parameters to update. Each field must be set. */
    params?: EventParams;
}
/** MsgUpdateEventParams is the Msg/UpdateEventParams request type. */
export interface MsgUpdateEventParamsSDKType {
    authority: string;
    params?: EventParamsSDKType;
}
/** MsgUpdateEventParamsResponse is the Msg/UpdateEventParams response type. */
export interface MsgUpdateEventParamsResponse {
}
/** MsgUpdateEventParamsResponse is the Msg/UpdateEventParams response type. */
export interface MsgUpdateEventParamsResponseSDKType {
}
/** MsgUpdateProposeParams is the Msg/UpdateProposeParams request type. */
export interface MsgUpdateProposeParams {
    authority: string;
    /** The parameters to update. Each field must be set. */
    params?: ProposeParams;
}
/** MsgUpdateProposeParams is the Msg/UpdateProposeParams request type. */
export interface MsgUpdateProposeParamsSDKType {
    authority: string;
    params?: ProposeParamsSDKType;
}
/** MsgUpdateProposeParamsResponse is the Msg/UpdateProposeParams response type. */
export interface MsgUpdateProposeParamsResponse {
}
/** MsgUpdateProposeParamsResponse is the Msg/UpdateProposeParams response type. */
export interface MsgUpdateProposeParamsResponseSDKType {
}
/** MsgUpdateSafetyParams is the Msg/UpdateSafetyParams request type. */
export interface MsgUpdateSafetyParams {
    authority: string;
    /** The parameters to update. Each field must be set. */
    params?: SafetyParams;
}
/** MsgUpdateSafetyParams is the Msg/UpdateSafetyParams request type. */
export interface MsgUpdateSafetyParamsSDKType {
    authority: string;
    params?: SafetyParamsSDKType;
}
/** MsgUpdateSafetyParamsResponse is the Msg/UpdateSafetyParams response type. */
export interface MsgUpdateSafetyParamsResponse {
}
/** MsgUpdateSafetyParamsResponse is the Msg/UpdateSafetyParams response type. */
export interface MsgUpdateSafetyParamsResponseSDKType {
}
export declare const MsgAcknowledgeBridges: {
    encode(message: MsgAcknowledgeBridges, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcknowledgeBridges;
    fromPartial(object: DeepPartial<MsgAcknowledgeBridges>): MsgAcknowledgeBridges;
};
export declare const MsgAcknowledgeBridgesResponse: {
    encode(_: MsgAcknowledgeBridgesResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAcknowledgeBridgesResponse;
    fromPartial(_: DeepPartial<MsgAcknowledgeBridgesResponse>): MsgAcknowledgeBridgesResponse;
};
export declare const MsgCompleteBridge: {
    encode(message: MsgCompleteBridge, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCompleteBridge;
    fromPartial(object: DeepPartial<MsgCompleteBridge>): MsgCompleteBridge;
};
export declare const MsgCompleteBridgeResponse: {
    encode(_: MsgCompleteBridgeResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgCompleteBridgeResponse;
    fromPartial(_: DeepPartial<MsgCompleteBridgeResponse>): MsgCompleteBridgeResponse;
};
export declare const MsgUpdateEventParams: {
    encode(message: MsgUpdateEventParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateEventParams;
    fromPartial(object: DeepPartial<MsgUpdateEventParams>): MsgUpdateEventParams;
};
export declare const MsgUpdateEventParamsResponse: {
    encode(_: MsgUpdateEventParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateEventParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateEventParamsResponse>): MsgUpdateEventParamsResponse;
};
export declare const MsgUpdateProposeParams: {
    encode(message: MsgUpdateProposeParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateProposeParams;
    fromPartial(object: DeepPartial<MsgUpdateProposeParams>): MsgUpdateProposeParams;
};
export declare const MsgUpdateProposeParamsResponse: {
    encode(_: MsgUpdateProposeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateProposeParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateProposeParamsResponse>): MsgUpdateProposeParamsResponse;
};
export declare const MsgUpdateSafetyParams: {
    encode(message: MsgUpdateSafetyParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSafetyParams;
    fromPartial(object: DeepPartial<MsgUpdateSafetyParams>): MsgUpdateSafetyParams;
};
export declare const MsgUpdateSafetyParamsResponse: {
    encode(_: MsgUpdateSafetyParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateSafetyParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateSafetyParamsResponse>): MsgUpdateSafetyParamsResponse;
};
