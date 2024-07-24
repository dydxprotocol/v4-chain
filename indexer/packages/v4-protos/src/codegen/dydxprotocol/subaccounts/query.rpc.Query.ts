import { Rpc } from "../../helpers";
import * as _m0 from "protobufjs/minimal";
import { QueryClient, createProtobufRpcClient } from "@cosmjs/stargate";
import { QueryGetSubaccountRequest, QuerySubaccountResponse, QueryAllSubaccountRequest, QuerySubaccountAllResponse, QueryGetWithdrawalAndTransfersBlockedInfoRequest, QueryGetWithdrawalAndTransfersBlockedInfoResponse, QueryCollateralPoolAddressRequest, QueryCollateralPoolAddressResponse } from "./query";
/** Query defines the gRPC querier service. */

export interface Query {
  /** Queries a Subaccount by id */
  subaccount(request: QueryGetSubaccountRequest): Promise<QuerySubaccountResponse>;
  /** Queries a list of Subaccount items. */

  subaccountAll(request?: QueryAllSubaccountRequest): Promise<QuerySubaccountAllResponse>;
  /**
   * Queries information about whether withdrawal and transfers are blocked, and
   * if so which block they are re-enabled on.
   */

  getWithdrawalAndTransfersBlockedInfo(request: QueryGetWithdrawalAndTransfersBlockedInfoRequest): Promise<QueryGetWithdrawalAndTransfersBlockedInfoResponse>;
  /** Queries the collateral pool account address for a perpetual id. */

  collateralPoolAddress(request: QueryCollateralPoolAddressRequest): Promise<QueryCollateralPoolAddressResponse>;
}
export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.subaccount = this.subaccount.bind(this);
    this.subaccountAll = this.subaccountAll.bind(this);
    this.getWithdrawalAndTransfersBlockedInfo = this.getWithdrawalAndTransfersBlockedInfo.bind(this);
    this.collateralPoolAddress = this.collateralPoolAddress.bind(this);
  }

  subaccount(request: QueryGetSubaccountRequest): Promise<QuerySubaccountResponse> {
    const data = QueryGetSubaccountRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.subaccounts.Query", "Subaccount", data);
    return promise.then(data => QuerySubaccountResponse.decode(new _m0.Reader(data)));
  }

  subaccountAll(request: QueryAllSubaccountRequest = {
    pagination: undefined
  }): Promise<QuerySubaccountAllResponse> {
    const data = QueryAllSubaccountRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.subaccounts.Query", "SubaccountAll", data);
    return promise.then(data => QuerySubaccountAllResponse.decode(new _m0.Reader(data)));
  }

  getWithdrawalAndTransfersBlockedInfo(request: QueryGetWithdrawalAndTransfersBlockedInfoRequest): Promise<QueryGetWithdrawalAndTransfersBlockedInfoResponse> {
    const data = QueryGetWithdrawalAndTransfersBlockedInfoRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.subaccounts.Query", "GetWithdrawalAndTransfersBlockedInfo", data);
    return promise.then(data => QueryGetWithdrawalAndTransfersBlockedInfoResponse.decode(new _m0.Reader(data)));
  }

  collateralPoolAddress(request: QueryCollateralPoolAddressRequest): Promise<QueryCollateralPoolAddressResponse> {
    const data = QueryCollateralPoolAddressRequest.encode(request).finish();
    const promise = this.rpc.request("dydxprotocol.subaccounts.Query", "CollateralPoolAddress", data);
    return promise.then(data => QueryCollateralPoolAddressResponse.decode(new _m0.Reader(data)));
  }

}
export const createRpcQueryExtension = (base: QueryClient) => {
  const rpc = createProtobufRpcClient(base);
  const queryService = new QueryClientImpl(rpc);
  return {
    subaccount(request: QueryGetSubaccountRequest): Promise<QuerySubaccountResponse> {
      return queryService.subaccount(request);
    },

    subaccountAll(request?: QueryAllSubaccountRequest): Promise<QuerySubaccountAllResponse> {
      return queryService.subaccountAll(request);
    },

    getWithdrawalAndTransfersBlockedInfo(request: QueryGetWithdrawalAndTransfersBlockedInfoRequest): Promise<QueryGetWithdrawalAndTransfersBlockedInfoResponse> {
      return queryService.getWithdrawalAndTransfersBlockedInfo(request);
    },

    collateralPoolAddress(request: QueryCollateralPoolAddressRequest): Promise<QueryCollateralPoolAddressResponse> {
      return queryService.collateralPoolAddress(request);
    }

  };
};