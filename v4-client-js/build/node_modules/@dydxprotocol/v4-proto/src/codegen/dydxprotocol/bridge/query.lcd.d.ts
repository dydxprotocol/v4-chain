import { LCDClient } from "@osmonauts/lcd";
import { QueryEventParamsRequest, QueryEventParamsResponseSDKType, QueryProposeParamsRequest, QueryProposeParamsResponseSDKType, QuerySafetyParamsRequest, QuerySafetyParamsResponseSDKType, QueryAcknowledgedEventInfoRequest, QueryAcknowledgedEventInfoResponseSDKType, QueryRecognizedEventInfoRequest, QueryRecognizedEventInfoResponseSDKType, QueryDelayedCompleteBridgeMessagesRequest, QueryDelayedCompleteBridgeMessagesResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    eventParams(_params?: QueryEventParamsRequest): Promise<QueryEventParamsResponseSDKType>;
    proposeParams(_params?: QueryProposeParamsRequest): Promise<QueryProposeParamsResponseSDKType>;
    safetyParams(_params?: QuerySafetyParamsRequest): Promise<QuerySafetyParamsResponseSDKType>;
    acknowledgedEventInfo(_params?: QueryAcknowledgedEventInfoRequest): Promise<QueryAcknowledgedEventInfoResponseSDKType>;
    recognizedEventInfo(_params?: QueryRecognizedEventInfoRequest): Promise<QueryRecognizedEventInfoResponseSDKType>;
    delayedCompleteBridgeMessages(params: QueryDelayedCompleteBridgeMessagesRequest): Promise<QueryDelayedCompleteBridgeMessagesResponseSDKType>;
}
