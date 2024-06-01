import { Rpc } from "../../helpers";
import { MsgUpdatePerpetualFeeParams, MsgUpdatePerpetualFeeParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /** UpdatePerpetualFeeParams updates the PerpetualFeeParams in state. */
    updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    updatePerpetualFeeParams(request: MsgUpdatePerpetualFeeParams): Promise<MsgUpdatePerpetualFeeParamsResponse>;
}
