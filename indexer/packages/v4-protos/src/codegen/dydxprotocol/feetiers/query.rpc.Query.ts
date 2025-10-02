import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponse, QueryUserFeeTierRequest, QueryUserFeeTierResponse, QueryFeeHolidayParamsRequest, QueryFeeHolidayParamsResponse, QueryAllFeeHolidayParamsRequest, QueryAllFeeHolidayParamsResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries the PerpetualFeeParams. */
  perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse>;
  /** Queries a user's fee tier */

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse>;
  /** Queries the FeeHolidayParams */

  feeHolidayParams(request: QueryFeeHolidayParamsRequest): Promise<QueryFeeHolidayParamsResponse>;
  /** Queries all fee holiday params */

  allFeeHolidays(request?: QueryAllFeeHolidayParamsRequest): Promise<QueryAllFeeHolidayParamsResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.perpetualFeeParams = this.perpetualFeeParams.bind(this);
    this.userFeeTier = this.userFeeTier.bind(this);
    this.feeHolidayParams = this.feeHolidayParams.bind(this);
    this.allFeeHolidays = this.allFeeHolidays.bind(this);
  }

  perpetualFeeParams(request: QueryPerpetualFeeParamsRequest = {}): Promise<QueryPerpetualFeeParamsResponse> {
    const data = QueryPerpetualFeeParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "PerpetualFeeParams", data);
    return promise.then(data => QueryPerpetualFeeParamsResponse.decode(new _m0.Reader(data)));
  }

  userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse> {
    const data = QueryUserFeeTierRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "UserFeeTier", data);
    return promise.then(data => QueryUserFeeTierResponse.decode(new _m0.Reader(data)));
  }

  feeHolidayParams(request: QueryFeeHolidayParamsRequest): Promise<QueryFeeHolidayParamsResponse> {
    const data = QueryFeeHolidayParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "FeeHolidayParams", data);
    return promise.then(data => QueryFeeHolidayParamsResponse.decode(new _m0.Reader(data)));
  }

  allFeeHolidays(request: QueryAllFeeHolidayParamsRequest = {}): Promise<QueryAllFeeHolidayParamsResponse> {
    const data = QueryAllFeeHolidayParamsRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.feetiers.Query", "AllFeeHolidays", data);
    return promise.then(data => QueryAllFeeHolidayParamsResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    perpetualFeeParams(request?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponse> {
      return queryService.perpetualFeeParams(request);
    },

    userFeeTier(request: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponse> {
      return queryService.userFeeTier(request);
    },

    feeHolidayParams(request: QueryFeeHolidayParamsRequest): Promise<QueryFeeHolidayParamsResponse> {
      return queryService.feeHolidayParams(request);
    },

    allFeeHolidays(request?: QueryAllFeeHolidayParamsRequest): Promise<QueryAllFeeHolidayParamsResponse> {
      return queryService.allFeeHolidays(request);
    }

  };
};