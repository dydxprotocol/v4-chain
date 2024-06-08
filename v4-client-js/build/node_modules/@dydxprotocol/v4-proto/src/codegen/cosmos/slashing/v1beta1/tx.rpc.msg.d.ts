import { Rpc } from "../../../helpers";
import { MsgUnjail, MsgUnjailResponse, MsgUpdateParams, MsgUpdateParamsResponse } from "./tx";
/** Msg defines the slashing Msg service. */
export interface Msg {
    /**
     * Unjail defines a method for unjailing a jailed validator, thus returning
     * them into the bonded validator set, so they can begin receiving provisions
     * and rewards again.
     */
    unjail(request: MsgUnjail): Promise<MsgUnjailResponse>;
    /**
     * UpdateParams defines a governance operation for updating the x/slashing module
     * parameters. The authority defaults to the x/gov module account.
     *
     * Since: cosmos-sdk 0.47
     */
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    unjail(request: MsgUnjail): Promise<MsgUnjailResponse>;
    updateParams(request: MsgUpdateParams): Promise<MsgUpdateParamsResponse>;
}
