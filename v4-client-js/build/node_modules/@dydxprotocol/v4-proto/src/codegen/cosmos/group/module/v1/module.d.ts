/// <reference types="long" />
import { Duration, DurationSDKType } from "../../../../google/protobuf/duration";
import { Long, DeepPartial } from "../../../../helpers";
import * as _m0 from "protobufjs/minimal";
/** Module is the config object of the group module. */
export interface Module {
    /**
     * max_execution_period defines the max duration after a proposal's voting period ends that members can send a MsgExec
     * to execute the proposal.
     */
    maxExecutionPeriod?: Duration;
    /**
     * max_metadata_len defines the max length of the metadata bytes field for various entities within the group module.
     * Defaults to 255 if not explicitly set.
     */
    maxMetadataLen: Long;
}
/** Module is the config object of the group module. */
export interface ModuleSDKType {
    max_execution_period?: DurationSDKType;
    max_metadata_len: Long;
}
export declare const Module: {
    encode(message: Module, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): Module;
    fromPartial(object: DeepPartial<Module>): Module;
};
