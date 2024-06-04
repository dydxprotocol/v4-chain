import { LCDClient } from "@osmonauts/lcd";
import { QueryNextDelayedMessageIdRequest, QueryNextDelayedMessageIdResponseSDKType, QueryMessageRequest, QueryMessageResponseSDKType, QueryBlockMessageIdsRequest, QueryBlockMessageIdsResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    nextDelayedMessageId(_params?: QueryNextDelayedMessageIdRequest): Promise<QueryNextDelayedMessageIdResponseSDKType>;
    message(params: QueryMessageRequest): Promise<QueryMessageResponseSDKType>;
    blockMessageIds(params: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponseSDKType>;
}
