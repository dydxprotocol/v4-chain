import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../../helpers";
/** Module is the config object of the authz module. */
export interface Module {
}
/** Module is the config object of the authz module. */
export interface ModuleSDKType {
}
export declare const Module: {
    encode(_: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(_: DeepPartial<Module>): Module;
};
