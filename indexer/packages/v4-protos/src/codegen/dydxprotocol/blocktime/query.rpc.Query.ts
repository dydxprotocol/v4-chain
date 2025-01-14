import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryDowntimeParamsRequest, QueryDowntimeParamsResponse, QueryPreviousBlockInfoRequest, QueryPreviousBlockInfoResponse, QueryAllDowntimeInfoRequest, QueryAllDowntimeInfoResponse, QuerySynchronyParamsRequest, QuerySynchronyParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the DowntimeParams. */
  downtimeParams(request?: QueryDowntimeParamsRequest): Promise<QueryDowntimeParamsResponse>;
  /** Queries the information of the previous block */

  previousBlockInfo(request?: QueryPreviousBlockInfoRequest): Promise<QueryPreviousBlockInfoResponse>;
  /** Queries all recorded downtime info. */

  allDowntimeInfo(request?: QueryAllDowntimeInfoRequest): Promise<QueryAllDowntimeInfoResponse>;
  /** Queries the SynchronyParams. */

  synchronyParams(request?: QuerySynchronyParamsRequest): Promise<QuerySynchronyParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.downtimeParams = this.downtimeParams.bind(this);
    this.previousBlockInfo = this.previousBlockInfo.bind(this);
    this.allDowntimeInfo = this.allDowntimeInfo.bind(this);
    this.synchronyParams = this.synchronyParams.bind(this);
  }

  downtimeParams(request: QueryDowntimeParamsRequest = {}): Promise<QueryDowntimeParamsResponse> {
    const data = QueryDowntimeParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Query", "DowntimeParams", data);
    return promise.then(data => QueryDowntimeParamsResponse.decode(new _m0.Reader(data)));
  }

  previousBlockInfo(request: QueryPreviousBlockInfoRequest = {}): Promise<QueryPreviousBlockInfoResponse> {
    const data = QueryPreviousBlockInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Query", "PreviousBlockInfo", data);
    return promise.then(data => QueryPreviousBlockInfoResponse.decode(new _m0.Reader(data)));
  }

  allDowntimeInfo(request: QueryAllDowntimeInfoRequest = {}): Promise<QueryAllDowntimeInfoResponse> {
    const data = QueryAllDowntimeInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Query", "AllDowntimeInfo", data);
    return promise.then(data => QueryAllDowntimeInfoResponse.decode(new _m0.Reader(data)));
  }

  synchronyParams(request: QuerySynchronyParamsRequest = {}): Promise<QuerySynchronyParamsResponse> {
    const data = QuerySynchronyParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.blocktime.Query", "SynchronyParams", data);
    return promise.then(data => QuerySynchronyParamsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    downtimeParams(request?: QueryDowntimeParamsRequest): Promise<QueryDowntimeParamsResponse> {
      return queryService.downtimeParams(request);
    },

    previousBlockInfo(request?: QueryPreviousBlockInfoRequest): Promise<QueryPreviousBlockInfoResponse> {
      return queryService.previousBlockInfo(request);
    },

    allDowntimeInfo(request?: QueryAllDowntimeInfoRequest): Promise<QueryAllDowntimeInfoResponse> {
      return queryService.allDowntimeInfo(request);
    },

    synchronyParams(request?: QuerySynchronyParamsRequest): Promise<QuerySynchronyParamsResponse> {
      return queryService.synchronyParams(request);
    }

  };
};