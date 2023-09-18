import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryNumMessagesRequest, QueryNumMessagesResponse, QueryMessageRequest, QueryMessageResponse, QueryBlockMessageIdsRequest, QueryBlockMessageIdsResponse, QueryAllMessagesRequest, QueryAllMessagesResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the number of DelayedMessages. */
  numMessages(request?: QueryNumMessagesRequest): Promise<QueryNumMessagesResponse>;
  /** Queries the DelayedMessage by id. */

  message(request: QueryMessageRequest): Promise<QueryMessageResponse>;
  /** Queries the DelayedMessages at a given block height. */

  blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse>;
  /** Queries all DelayedMessages. */

  allMessages(request?: QueryAllMessagesRequest): Promise<QueryAllMessagesResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.numMessages = this.numMessages.bind(this);
    this.message = this.message.bind(this);
    this.blockMessageIds = this.blockMessageIds.bind(this);
    this.allMessages = this.allMessages.bind(this);
  }

  numMessages(request: QueryNumMessagesRequest = {}): Promise<QueryNumMessagesResponse> {
    const data = QueryNumMessagesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Query", "NumMessages", data);
    return promise.then(data => QueryNumMessagesResponse.decode(new _m0.Reader(data)));
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

  allMessages(request: QueryAllMessagesRequest = {}): Promise<QueryAllMessagesResponse> {
    const data = QueryAllMessagesRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.delaymsg.Query", "AllMessages", data);
    return promise.then(data => QueryAllMessagesResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    numMessages(request?: QueryNumMessagesRequest): Promise<QueryNumMessagesResponse> {
      return queryService.numMessages(request);
    },

    message(request: QueryMessageRequest): Promise<QueryMessageResponse> {
      return queryService.message(request);
    },

    blockMessageIds(request: QueryBlockMessageIdsRequest): Promise<QueryBlockMessageIdsResponse> {
      return queryService.blockMessageIds(request);
    },

    allMessages(request?: QueryAllMessagesRequest): Promise<QueryAllMessagesResponse> {
      return queryService.allMessages(request);
    }

  };
};