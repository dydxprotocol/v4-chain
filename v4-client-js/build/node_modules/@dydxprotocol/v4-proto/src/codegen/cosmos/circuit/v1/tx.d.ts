import { Permissions, PermissionsSDKType } from "./types";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** MsgAuthorizeCircuitBreaker defines the Msg/AuthorizeCircuitBreaker request type. */
export interface MsgAuthorizeCircuitBreaker {
    /**
     * granter is the granter of the circuit breaker permissions and must have
     * LEVEL_SUPER_ADMIN.
     */
    granter: string;
    /** grantee is the account authorized with the provided permissions. */
    grantee: string;
    /**
     * permissions are the circuit breaker permissions that the grantee receives.
     * These will overwrite any existing permissions. LEVEL_NONE_UNSPECIFIED can
     * be specified to revoke all permissions.
     */
    permissions?: Permissions;
}
/** MsgAuthorizeCircuitBreaker defines the Msg/AuthorizeCircuitBreaker request type. */
export interface MsgAuthorizeCircuitBreakerSDKType {
    granter: string;
    grantee: string;
    permissions?: PermissionsSDKType;
}
/** MsgAuthorizeCircuitBreakerResponse defines the Msg/AuthorizeCircuitBreaker response type. */
export interface MsgAuthorizeCircuitBreakerResponse {
    success: boolean;
}
/** MsgAuthorizeCircuitBreakerResponse defines the Msg/AuthorizeCircuitBreaker response type. */
export interface MsgAuthorizeCircuitBreakerResponseSDKType {
    success: boolean;
}
/** MsgTripCircuitBreaker defines the Msg/TripCircuitBreaker request type. */
export interface MsgTripCircuitBreaker {
    /** authority is the account authorized to trip the circuit breaker. */
    authority: string;
    /**
     * msg_type_urls specifies a list of type URLs to immediately stop processing.
     * IF IT IS LEFT EMPTY, ALL MSG PROCESSING WILL STOP IMMEDIATELY.
     * This value is validated against the authority's permissions and if the
     * authority does not have permissions to trip the specified msg type URLs
     * (or all URLs), the operation will fail.
     */
    msgTypeUrls: string[];
}
/** MsgTripCircuitBreaker defines the Msg/TripCircuitBreaker request type. */
export interface MsgTripCircuitBreakerSDKType {
    authority: string;
    msg_type_urls: string[];
}
/** MsgTripCircuitBreakerResponse defines the Msg/TripCircuitBreaker response type. */
export interface MsgTripCircuitBreakerResponse {
    success: boolean;
}
/** MsgTripCircuitBreakerResponse defines the Msg/TripCircuitBreaker response type. */
export interface MsgTripCircuitBreakerResponseSDKType {
    success: boolean;
}
/** MsgResetCircuitBreaker defines the Msg/ResetCircuitBreaker request type. */
export interface MsgResetCircuitBreaker {
    /** authority is the account authorized to trip or reset the circuit breaker. */
    authority: string;
    /**
     * msg_type_urls specifies a list of Msg type URLs to resume processing. If
     * it is left empty all Msg processing for type URLs that the account is
     * authorized to trip will resume.
     */
    msgTypeUrls: string[];
}
/** MsgResetCircuitBreaker defines the Msg/ResetCircuitBreaker request type. */
export interface MsgResetCircuitBreakerSDKType {
    authority: string;
    msg_type_urls: string[];
}
/** MsgResetCircuitBreakerResponse defines the Msg/ResetCircuitBreaker response type. */
export interface MsgResetCircuitBreakerResponse {
    success: boolean;
}
/** MsgResetCircuitBreakerResponse defines the Msg/ResetCircuitBreaker response type. */
export interface MsgResetCircuitBreakerResponseSDKType {
    success: boolean;
}
export declare const MsgAuthorizeCircuitBreaker: {
    encode(message: MsgAuthorizeCircuitBreaker, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAuthorizeCircuitBreaker;
    fromPartial(object: DeepPartial<MsgAuthorizeCircuitBreaker>): MsgAuthorizeCircuitBreaker;
};
export declare const MsgAuthorizeCircuitBreakerResponse: {
    encode(message: MsgAuthorizeCircuitBreakerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgAuthorizeCircuitBreakerResponse;
    fromPartial(object: DeepPartial<MsgAuthorizeCircuitBreakerResponse>): MsgAuthorizeCircuitBreakerResponse;
};
export declare const MsgTripCircuitBreaker: {
    encode(message: MsgTripCircuitBreaker, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgTripCircuitBreaker;
    fromPartial(object: DeepPartial<MsgTripCircuitBreaker>): MsgTripCircuitBreaker;
};
export declare const MsgTripCircuitBreakerResponse: {
    encode(message: MsgTripCircuitBreakerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgTripCircuitBreakerResponse;
    fromPartial(object: DeepPartial<MsgTripCircuitBreakerResponse>): MsgTripCircuitBreakerResponse;
};
export declare const MsgResetCircuitBreaker: {
    encode(message: MsgResetCircuitBreaker, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgResetCircuitBreaker;
    fromPartial(object: DeepPartial<MsgResetCircuitBreaker>): MsgResetCircuitBreaker;
};
export declare const MsgResetCircuitBreakerResponse: {
    encode(message: MsgResetCircuitBreakerResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgResetCircuitBreakerResponse;
    fromPartial(object: DeepPartial<MsgResetCircuitBreakerResponse>): MsgResetCircuitBreakerResponse;
};
