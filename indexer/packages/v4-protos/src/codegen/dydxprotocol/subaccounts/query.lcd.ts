import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryGetSubaccountRequest, QuerySubaccountResponseSDKType, QueryAllSubaccountRequest, QuerySubaccountAllResponseSDKType, QueryGetWithdrawalAndTransfersBlockedInfoRequest, QueryGetWithdrawalAndTransfersBlockedInfoResponseSDKType, QueryCollateralPoolAddressRequest, QueryCollateralPoolAddressResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.subaccount = this.subaccount.bind(this);
    this.subaccountAll = this.subaccountAll.bind(this);
    this.getWithdrawalAndTransfersBlockedInfo = this.getWithdrawalAndTransfersBlockedInfo.bind(this);
    this.collateralPoolAddress = this.collateralPoolAddress.bind(this);
  }
  /* Queries a Subaccount by id */


  async subaccount(params: QueryGetSubaccountRequest): Promise<QuerySubaccountResponseSDKType> {
    const endpoint = `dydxprotocol/subaccounts/subaccount/${params.owner}/${params.number}`;
    return await this.req.get<QuerySubaccountResponseSDKType>(endpoint);
  }
  /* Queries a list of Subaccount items. */


  async subaccountAll(params: QueryAllSubaccountRequest = {
    pagination: undefined
  }): Promise<QuerySubaccountAllResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/subaccounts/subaccount`;
    return await this.req.get<QuerySubaccountAllResponseSDKType>(endpoint, options);
  }
  /* Queries information about whether withdrawal and transfers are blocked, and
   if so which block they are re-enabled on. */


  async getWithdrawalAndTransfersBlockedInfo(params: QueryGetWithdrawalAndTransfersBlockedInfoRequest): Promise<QueryGetWithdrawalAndTransfersBlockedInfoResponseSDKType> {
    const endpoint = `dydxprotocol/subaccounts/withdrawals_and_transfers_blocked_info/${params.perpetualId}`;
    return await this.req.get<QueryGetWithdrawalAndTransfersBlockedInfoResponseSDKType>(endpoint);
  }
  /* Queries the collateral pool account address for a perpetual id. */


  async collateralPoolAddress(params: QueryCollateralPoolAddressRequest): Promise<QueryCollateralPoolAddressResponseSDKType> {
    const endpoint = `dydxprotocol/subaccounts/collateral_pool_address/${params.perpetualId}`;
    return await this.req.get<QueryCollateralPoolAddressResponseSDKType>(endpoint);
  }

}