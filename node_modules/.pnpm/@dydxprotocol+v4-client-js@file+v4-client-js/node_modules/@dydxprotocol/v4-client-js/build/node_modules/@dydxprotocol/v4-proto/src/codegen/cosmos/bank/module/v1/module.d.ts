import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
/** Module is the config object of the bank module. */
export interface Module {
    /**
     * blocked_module_accounts_override configures exceptional module accounts which should be blocked from receiving
     * funds. If left empty it defaults to the list of account names supplied in the auth module configuration as
     * module_account_permissions
     */
    blockedModuleAccountsOverride: string[];
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
/** Module is the config object of the bank module. */
export interface ModuleSDKType {
    blocked_module_accounts_override: string[];
    authority: string;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
