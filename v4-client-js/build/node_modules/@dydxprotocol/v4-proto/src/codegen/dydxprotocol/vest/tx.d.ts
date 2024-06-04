import { VestEntry, VestEntrySDKType } from "./vest_entry";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */
export interface MsgDeleteVestEntry {
    /** authority is the address that controls the module. */
    authority: string;
    /** The vester account of the vest entry to delete. */
    vesterAccount: string;
}
/** MsgDeleteVestEntry is the Msg/DeleteVestEntry request type. */
export interface MsgDeleteVestEntrySDKType {
    authority: string;
    vester_account: string;
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */
export interface MsgDeleteVestEntryResponse {
}
/** MsgDeleteVestEntryResponse is the Msg/DeleteVestEntry response type. */
export interface MsgDeleteVestEntryResponseSDKType {
}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */
export interface MsgSetVestEntry {
    /** authority is the address that controls the module. */
    authority: string;
    /** The vest entry to set. */
    entry?: VestEntry;
}
/** MsgSetVestEntry is the Msg/SetVestEntry request type. */
export interface MsgSetVestEntrySDKType {
    authority: string;
    entry?: VestEntrySDKType;
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */
export interface MsgSetVestEntryResponse {
}
/** MsgSetVestEntryResponse is the Msg/SetVestEntry response type. */
export interface MsgSetVestEntryResponseSDKType {
}
export declare const MsgDeleteVestEntry: {
    encode(message: MsgDeleteVestEntry, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestEntry;
    fromPartial(object: DeepPartial<MsgDeleteVestEntry>): MsgDeleteVestEntry;
};
export declare const MsgDeleteVestEntryResponse: {
    encode(_: MsgDeleteVestEntryResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgDeleteVestEntryResponse;
    fromPartial(_: DeepPartial<MsgDeleteVestEntryResponse>): MsgDeleteVestEntryResponse;
};
export declare const MsgSetVestEntry: {
    encode(message: MsgSetVestEntry, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestEntry;
    fromPartial(object: DeepPartial<MsgSetVestEntry>): MsgSetVestEntry;
};
export declare const MsgSetVestEntryResponse: {
    encode(_: MsgSetVestEntryResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgSetVestEntryResponse;
    fromPartial(_: DeepPartial<MsgSetVestEntryResponse>): MsgSetVestEntryResponse;
};
