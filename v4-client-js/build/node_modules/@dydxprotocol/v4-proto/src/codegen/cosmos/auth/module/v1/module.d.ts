import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
/** Module is the config object for the auth module. */
export interface Module {
    /** bech32_prefix is the bech32 account prefix for the app. */
    bech32Prefix: string;
    /** module_account_permissions are module account permissions. */
    moduleAccountPermissions: ModuleAccountPermission[];
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
/** Module is the config object for the auth module. */
export interface ModuleSDKType {
    bech32_prefix: string;
    module_account_permissions: ModuleAccountPermissionSDKType[];
    authority: string;
}
/** ModuleAccountPermission represents permissions for a module account. */
export interface ModuleAccountPermission {
    /** account is the name of the module. */
    account: string;
    /**
     * permissions are the permissions this module has. Currently recognized
     * values are minter, burner and staking.
     */
    permissions: string[];
}
/** ModuleAccountPermission represents permissions for a module account. */
export interface ModuleAccountPermissionSDKType {
    account: string;
    permissions: string[];
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
export declare const ModuleAccountPermission: {
    encode(message: ModuleAccountPermission, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): ModuleAccountPermission;
    fromPartial(object: DeepPartial<ModuleAccountPermission>): ModuleAccountPermission;
};
