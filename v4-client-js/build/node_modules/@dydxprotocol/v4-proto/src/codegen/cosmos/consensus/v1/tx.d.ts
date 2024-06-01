import { BlockParams, BlockParamsSDKType, EvidenceParams, EvidenceParamsSDKType, ValidatorParams, ValidatorParamsSDKType, ABCIParams, ABCIParamsSDKType } from "../../../tendermint/types/params";
import * as _m0 from "protobufjs/minimal";
import { DeepPartial } from "../../../helpers";
/** MsgUpdateParams is the Msg/UpdateParams request type. */
export interface MsgUpdateParams {
    /** authority is the address that controls the module (defaults to x/gov unless overwritten). */
    authority: string;
    /**
     * params defines the x/consensus parameters to update.
     * VersionsParams is not included in this Msg because it is tracked
     * separarately in x/upgrade.
     *
     * NOTE: All parameters must be supplied.
     */
    block?: BlockParams;
    evidence?: EvidenceParams;
    validator?: ValidatorParams;
    /** Since: cosmos-sdk 0.50 */
    abci?: ABCIParams;
}
/** MsgUpdateParams is the Msg/UpdateParams request type. */
export interface MsgUpdateParamsSDKType {
    authority: string;
    block?: BlockParamsSDKType;
    evidence?: EvidenceParamsSDKType;
    validator?: ValidatorParamsSDKType;
    abci?: ABCIParamsSDKType;
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 */
export interface MsgUpdateParamsResponse {
}
/**
 * MsgUpdateParamsResponse defines the response structure for executing a
 * MsgUpdateParams message.
 */
export interface MsgUpdateParamsResponseSDKType {
}
export declare const MsgUpdateParams: {
    encode(message: MsgUpdateParams, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParams;
    fromPartial(object: DeepPartial<MsgUpdateParams>): MsgUpdateParams;
};
export declare const MsgUpdateParamsResponse: {
    encode(_: MsgUpdateParamsResponse, writer?: _m0.Writer): _m0.Writer;
    decode(input: _m0.Reader | Uint8Array, length?: number): MsgUpdateParamsResponse;
    fromPartial(_: DeepPartial<MsgUpdateParamsResponse>): MsgUpdateParamsResponse;
};
