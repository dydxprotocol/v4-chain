import { DowntimeParams, DowntimeParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */
export interface MsgUpdateDowntimeParams {
    authority: string;
    /** Defines the parameters to update. All parameters must be supplied. */
    params?: DowntimeParams;
}
/** MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type. */
export interface MsgUpdateDowntimeParamsSDKType {
    authority: string;
    params?: DowntimeParamsSDKType;
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */
export interface MsgUpdateDowntimeParamsResponse {
}
/**
 * MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
 * type.
 */
export interface MsgUpdateDowntimeParamsResponseSDKType {
}
export declare const MsgUpdateDowntimeParams: {
    encode(message: MsgUpdateDowntimeParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDowntimeParams;
    fromPartial(object: DeepPartial<MsgUpdateDowntimeParams>): MsgUpdateDowntimeParams;
};
export declare const MsgUpdateDowntimeParamsResponse: {
    encode(_: MsgUpdateDowntimeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateDowntimeParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateDowntimeParamsResponse>): MsgUpdateDowntimeParamsResponse;
};
