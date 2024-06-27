import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryGetSubaccountRequest, QuerySubaccountResponseSDKType, QueryAllSubaccountRequest, QuerySubaccountAllResponseSDKType } from "./query";
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

}