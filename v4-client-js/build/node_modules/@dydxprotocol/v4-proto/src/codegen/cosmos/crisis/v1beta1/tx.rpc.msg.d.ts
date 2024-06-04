import { Rpc } from "../../../helpers";
import { MsgVerifyInvariant, MsgVerifyInvariantResponse, MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the bank Msg service. */
export interface Msg {
    /** VerifyInvariant defines a method to verify a particular invariant. */
    verifyInvariant(request: MsgVerifyInvariant): Promise<MsgVerifyInvariantResponse>;
    /**
     * UpdateParams defines a governance operation for updating the x/crisis module
     * parameters. The authority is defined in the keeper.
     *
     * Since: cosmos-sdk 0.47
     */
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    verifyInvariant(request: MsgVerifyInvariant): Promise<MsgVerifyInvariantResponse>;
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
