import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
/** Module is the config object of the staking module. */
export interface Module {
    /**
     * hooks_order specifies the order of staking hooks and should be a list
     * of module names which provide a staking hooks instance. If no order is
     * provided, then hooks will be applied in alphabetical order of module names.
     */
    hooksOrder: string[];
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
    /** bech32_prefix_validator is the bech32 validator prefix for the app. */
    bech32PrefixValidator: string;
    /** bech32_prefix_consensus is the bech32 consensus node prefix for the app. */
    bech32PrefixConsensus: string;
}
/** Module is the config object of the staking module. */
export interface ModuleSDKType {
    hooks_order: string[];
    authority: string;
    bech32_prefix_validator: string;
    bech32_prefix_consensus: string;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
