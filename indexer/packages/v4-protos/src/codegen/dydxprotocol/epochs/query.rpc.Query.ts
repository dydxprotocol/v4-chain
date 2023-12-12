import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetEpochInfoRequest, QueryEpochInfoResponse, QueryAllEpochInfoRequest, QueryEpochInfoAllResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a EpochInfo by name. */
  epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse>;
  /** Queries a list of EpochInfo items. */

  epochInfoAll(request?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.epochInfo = this.epochInfo.bind(this);
    this.epochInfoAll = this.epochInfoAll.bind(this);
  }

  epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse> {
    const data = QueryGetEpochInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.epochs.Query", "EpochInfo", data);
    return promise.then(data => QueryEpochInfoResponse.decode(new _m0.Reader(data)));
  }

  epochInfoAll(request: QueryAllEpochInfoRequest = {
    pagination: undefined
  }): Promise<QueryEpochInfoAllResponse> {
    const data = QueryAllEpochInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.epochs.Query", "EpochInfoAll", data);
    return promise.then(data => QueryEpochInfoAllResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    epochInfo(request: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponse> {
      return queryService.epochInfo(request);
    },

    epochInfoAll(request?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponse> {
      return queryService.epochInfoAll(request);
    }

  };
};