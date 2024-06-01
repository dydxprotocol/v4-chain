import { PerpetualFeeParams, PerpetualFeeParamsSDKType } from "./params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../helpers";
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */
export interface MsgUpdatePerpetualFeeParams {
    authority: string;
    /** Defines the parameters to update. All parameters must be supplied. */
    params?: PerpetualFeeParams;
}
/** MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type. */
export interface MsgUpdatePerpetualFeeParamsSDKType {
    authority: string;
    params?: PerpetualFeeParamsSDKType;
}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */
export interface MsgUpdatePerpetualFeeParamsResponse {
}
/**
 * MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
 * response type.
 */
export interface MsgUpdatePerpetualFeeParamsResponseSDKType {
}
export declare const MsgUpdatePerpetualFeeParams: {
    encode(message: MsgUpdatePerpetualFeeParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParams;
    fromPartial(object: DeepPartial<MsgUpdatePerpetualFeeParams>): MsgUpdatePerpetualFeeParams;
};
export declare const MsgUpdatePerpetualFeeParamsResponse: {
    encode(_: MsgUpdatePerpetualFeeParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdatePerpetualFeeParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdatePerpetualFeeParamsResponse>): MsgUpdatePerpetualFeeParamsResponse;
};
