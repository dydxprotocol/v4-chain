import { Rpc } from "../../helpers";
import { MsgAcknowledgeBridges, MsgAcknowledgeBridgesResponse, MsgCompleteBridge, MsgCompleteBridgeResponse, MsgUpdateEventParams, MsgUpdateEventParamsResponse, MsgUpdateProposeParams, MsgUpdateProposeParamsResponse, MsgUpdateSafetyParams, MsgUpdateSafetyParamsResponse } from "./tx";
/** Msg defines the Msg service. */
export interface Msg {
    /**
     * AcknowledgeBridges acknowledges bridges and sets them to complete at a
     * later block.
     */
    acknowledgeBridges(request: MsgAcknowledgeBridges): Promise<MsgAcknowledgeBridgesResponse>;
    /** CompleteBridge finalizes a bridge by minting coins to an address. */
    completeBridge(request: MsgCompleteBridge): Promise<MsgCompleteBridgeResponse>;
    /** UpdateEventParams updates the EventParams in state. */
    updateEventParams(request: MsgUpdateEventParams): Promise<MsgUpdateEventParamsResponse>;
    /** UpdateProposeParams updates the ProposeParams in state. */
    updateProposeParams(request: MsgUpdateProposeParams): Promise<MsgUpdateProposeParamsResponse>;
    /** UpdateSafetyParams updates the SafetyParams in state. */
    updateSafetyParams(request: MsgUpdateSafetyParams): Promise<MsgUpdateSafetyParamsResponse>;
}
export declare class MsgClientImpl implements Msg {
    private readonly rpc;
    constructor(rpc: Rpc);
    acknowledgeBridges(request: MsgAcknowledgeBridges): Promise<MsgAcknowledgeBridgesResponse>;
    completeBridge(request: MsgCompleteBridge): Promise<MsgCompleteBridgeResponse>;
    updateEventParams(request: MsgUpdateEventParams): Promise<MsgUpdateEventParamsResponse>;
    updateProposeParams(request: MsgUpdateProposeParams): Promise<MsgUpdateProposeParamsResponse>;
    updateSafetyParams(request: MsgUpdateSafetyParams): Promise<MsgUpdateSafetyParamsResponse>;
}
