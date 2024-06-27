import { Rpc } from "../../helpers";
import { QueryClient } from "@cosmjs/stargate";
import { QueryNextDelayedMessageIdRequest, QueryNextDelayedMessageIdResponse, QueryMessageRequest, QueryMessageResponse, QueryBlockMessageIdsRequest, QueryBlockMessageIdsResponse } from "./query";
/** Query defines the gRPC querier service. */
export interface Query {
    /** Queries the next DelayedMessage's id. */
    nextDelayedMessageId(request?: QueryNextDelayedMessageIdRequest): Promise<QueryNextDelayedMessageIdResponse>;
    /** Queries the DelayedMessage by id. */
    message(request: QueryMessageRequest): Promise<QueryMessageResponse>;
    /** Queries the DelayedMessages at a given block height. */
    blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse>;
}
export declare class QueryClientImpl implements Query {
    private readonly rpc;
    constructor(rpc: Rpc);
    nextDelayedMessageId(request?: QueryNextDelayedMessageIdRequest): Promise<QueryNextDelayedMessageIdResponse>;
    message(request: QueryMessageRequest): Promise<QueryMessageResponse>;
    blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse>;
}
export declare const createRpcQueryExtension: (base: QueryClient) => {
    nextDelayedMessageId(request?: QueryNextDelayedMessageIdRequest): Promise<QueryNextDelayedMessageIdResponse>;
    message(request: QueryMessageRequest): Promise<QueryMessageResponse>;
    blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse>;
};
