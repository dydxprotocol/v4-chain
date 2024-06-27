import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
/** Module is the config object of the distribution module. */
export interface Module {
    feeCollectorName: string;
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
/** Module is the config object of the distribution module. */
export interface ModuleSDKType {
    fee_collector_name: string;
    authority: string;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
