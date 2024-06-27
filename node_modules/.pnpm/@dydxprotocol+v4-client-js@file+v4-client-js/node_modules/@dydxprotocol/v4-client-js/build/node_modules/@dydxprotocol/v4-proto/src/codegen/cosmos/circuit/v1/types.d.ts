import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** Level is the permission level. */
export declare enum Permissions_Level {
    /**
     * LEVEL_NONE_UNSPECIFIED - LEVEL_NONE_UNSPECIFIED indicates that the account will have no circuit
     * breaker permissions.
     */
    LEVEL_NONE_UNSPECIFIED = 0,
    /**
     * LEVEL_SOME_MSGS - LEVEL_SOME_MSGS indicates that the account will have permission to
     * trip or reset the circuit breaker for some Msg type URLs. If this level
     * is chosen, a non-empty list of Msg type URLs must be provided in
     * limit_type_urls.
     */
    LEVEL_SOME_MSGS = 1,
    /**
     * LEVEL_ALL_MSGS - LEVEL_ALL_MSGS indicates that the account can trip or reset the circuit
     * breaker for Msg's of all type URLs.
     */
    LEVEL_ALL_MSGS = 2,
    /**
     * LEVEL_SUPER_ADMIN - LEVEL_SUPER_ADMIN indicates that the account can take all circuit breaker
     * actions and can grant permissions to other accounts.
     */
    LEVEL_SUPER_ADMIN = 3,
    UNRECOGNIZED = -1
}
export declare const Permissions_LevelSDKType: typeof Permissions_Level;
export declare function permissions_LevelFromJSON(object: any): Permissions_Level;
export declare function permissions_LevelToJSON(object: Permissions_Level): string;
/**
 * Permissions are the permissions that an account has to trip
 * or reset the circuit breaker.
 */
export interface Permissions {
    /** level is the level of permissions granted to this account. */
    level: Permissions_Level;
    /**
     * limit_type_urls is used with LEVEL_SOME_MSGS to limit the lists of Msg type
     * URLs that the account can trip. It is an error to use limit_type_urls with
     * a level other than LEVEL_SOME_MSGS.
     */
    limitTypeUrls: string[];
}
/**
 * Permissions are the permissions that an account has to trip
 * or reset the circuit breaker.
 */
export interface PermissionsSDKType {
    level: Permissions_Level;
    limit_type_urls: string[];
}
/** GenesisAccountPermissions is the account permissions for the circuit breaker in genesis */
export interface GenesisAccountPermissions {
    address: string;
    permissions?: Permissions;
}
/** GenesisAccountPermissions is the account permissions for the circuit breaker in genesis */
export interface GenesisAccountPermissionsSDKType {
    address: string;
    permissions?: PermissionsSDKType;
}
/** GenesisState is the state that must be provided at genesis. */
export interface GenesisState {
    accountPermissions: GenesisAccountPermissions[];
    disabledTypeUrls: string[];
}
/** GenesisState is the state that must be provided at genesis. */
export interface GenesisStateSDKType {
    account_permissions: GenesisAccountPermissionsSDKType[];
    disabled_type_urls: string[];
}
export declare const Permissions: {
    encode(message: Permissions, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Permissions;
    fromPartial(object: DeepPartial<Permissions>): Permissions;
};
export declare const GenesisAccountPermissions: {
    encode(message: GenesisAccountPermissions, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisAccountPermissions;
    fromPartial(object: DeepPartial<GenesisAccountPermissions>): GenesisAccountPermissions;
};
export declare const GenesisState: {
    encode(message: GenesisState, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState;
    fromPartial(object: DeepPartial<GenesisState>): GenesisState;
};
