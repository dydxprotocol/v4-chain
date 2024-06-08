/// <reference types="long" />
import { Long, DeepPartial } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** Module is the config object of the gov module. */
export interface Module {
    /**
     * max_metadata_len defines the maximum proposal metadata length.
     * Defaults to 255 if not explicitly set.
     */
    maxMetadataLen: Long;
    /** authority defines the custom module authority. If not set, defaults to the governance module. */
    authority: string;
}
/** Module is the config object of the gov module. */
export interface ModuleSDKType {
    max_metadata_len: Long;
    authority: string;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
