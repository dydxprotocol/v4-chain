import { Rpc } from "../../../helpers";
import { MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the consensus Msg service. */
export interface Msg {
    /**
     * UpdateParams defines a governance operation for updating the x/consensus module parameters.
     * The authority is defined in the keeper.
     *
     * Since: cosmos-sdk 0.47
     */
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
