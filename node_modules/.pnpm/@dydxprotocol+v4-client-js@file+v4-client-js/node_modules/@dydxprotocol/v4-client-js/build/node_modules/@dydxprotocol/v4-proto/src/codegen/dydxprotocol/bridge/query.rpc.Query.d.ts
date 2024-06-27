import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryEventParamsRequest, QueryEventParamsResponse, QueryProposeParamsRequest, QueryProposeParamsResponse, QuerySafetyParamsRequest, QuerySafetyParamsResponse, QueryAcknowledgedEventInfoRequest, QueryAcknowledgedEventInfoResponse, QueryRecognizedEventInfoRequest, QueryRecognizedEventInfoResponse, QueryDelayedCompleteBridgeMessagesRequest, QueryDelayedCompleteBridgeMessagesResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries the EventParams. */
    eventParams(request?: QueryEventParamsRequest): Promise<QueryEventParamsResponse>;
    /** Queries the ProposeParams. */
    proposeParams(request?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponse>;
    /** Queries the SafetyParams. */
    safetyParams(request?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponse>;
    /**
     * Queries the AcknowledgedEventInfo.
     * An "acknowledged" event is one that is in-consensus and has been stored
     * in-state.
     */
    acknowledgedEventInfo(request?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponse>;
    /**
     * Queries the RecognizedEventInfo.
     * A "recognized" event is one that is finalized on the Ethereum blockchain
     * and has been identified by the queried node. It is not yet in-consensus.
     */
    recognizedEventInfo(request?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponse>;
    /**
     * Queries all `MsgCompleteBridge` messages that are delayed (not yet
     * executed) and corresponding block heights at which they will execute.
     */
    delayedCompleteBridgeMessages(request: QueryDelayedCompleteBridgeMessagesRequest): Promise<QueryDelayedCompleteBridgeMessagesResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    eventParams(request?: QueryEventParamsRequest): Promise<QueryEventParamsResponse>;
    proposeParams(request?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponse>;
    safetyParams(request?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponse>;
    acknowledgedEventInfo(request?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponse>;
    recognizedEventInfo(request?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponse>;
    delayedCompleteBridgeMessages(request: QueryDelayedCompleteBridgeMessagesRequest): Promise<QueryDelayedCompleteBridgeMessagesResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    eventParams(request?: QueryEventParamsRequest): Promise<QueryEventParamsResponse>;
    proposeParams(request?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponse>;
    safetyParams(request?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponse>;
    acknowledgedEventInfo(request?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponse>;
    recognizedEventInfo(request?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponse>;
    delayedCompleteBridgeMessages(request: QueryDelayedCompleteBridgeMessagesRequest): Promise<QueryDelayedCompleteBridgeMessagesResponse>;
};
