import { Rpc } from "../../helpers";
import { MsgUpdateDowntimeParams, MsgUpdateDowntimeParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /** UpdateDowntimeParams updates the DowntimeParams in state. */
    updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    updateDowntimeParams(request: MsgUpdateDowntimeParams): Promise<MsgUpdateDowntimeParamsResponse>;
}
