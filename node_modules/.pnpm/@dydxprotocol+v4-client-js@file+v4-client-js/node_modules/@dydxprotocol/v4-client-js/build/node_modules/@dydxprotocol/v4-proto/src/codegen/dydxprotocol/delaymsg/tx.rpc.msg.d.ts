import { Rpc } from "../../helpers";
import { MsgDelayMessage, MsgDelayMessageResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /**
     * DelayMessage delays the execution of a message for a given number of
     * blocks.
     */
    delayMessage(request: MsgDelayMessage): Promise<MsgDelayMessageResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    delayMessage(request: MsgDelayMessage): Promise<MsgDelayMessageResponse>;
}
