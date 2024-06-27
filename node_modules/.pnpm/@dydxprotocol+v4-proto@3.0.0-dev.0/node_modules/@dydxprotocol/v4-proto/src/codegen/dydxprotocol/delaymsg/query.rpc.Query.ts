import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
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
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.nextDelayedMessageId = this.nextDelayedMessageId.bind(this);
    this.message = this.message.bind(this);
    this.blockMessageIds = this.blockMessageIds.bind(this);
  }

  nextDelayedMessageId(request: QueryNextDelayedMessageIdRequest = {}): Promise<QueryNextDelayedMessageIdResponse> {
    const data = QueryNextDelayedMessageIdRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Query", "NextDelayedMessageId", data);
    return promise.then(data => QueryNextDelayedMessageIdResponse.decode(new _m0.Reader(data)));
  }

  message(request: QueryMessageRequest): Promise<QueryMessageResponse> {
    const data = QueryMessageRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Query", "Message", data);
    return promise.then(data => QueryMessageResponse.decode(new _m0.Reader(data)));
  }

  blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse> {
    const data = QueryBlockMessageIdsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Query", "BlockMessageIds", data);
    return promise.then(data => QueryBlockMessageIdsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    nextDelayedMessageId(request?: QueryNextDelayedMessageIdRequest): Promise<QueryNextDelayedMessageIdResponse> {
      return queryService.nextDelayedMessageId(request);
    },

    message(request: QueryMessageRequest): Promise<QueryMessageResponse> {
      return queryService.message(request);
    },

    blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse> {
      return queryService.blockMessageIds(request);
    }

  };
};