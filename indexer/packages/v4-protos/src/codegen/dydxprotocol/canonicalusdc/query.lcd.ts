import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryStateRequest, QueryStateResponseSDKType, QueryPendingSettlementsRequest, QueryPendingSettlementsResponseSDKType, QueryPendingSettlementRequest, QueryPendingSettlementResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.state = this.state.bind(this);
    this.pendingSettlements = this.pendingSettlements.bind(this);
    this.pendingSettlement = this.pendingSettlement.bind(this);
  }
  /* State returns controls and aggregate backing state. */


  async state(_params: QueryStateRequest = {}): Promise<QueryStateResponseSDKType> {
    const endpoint = `dydxprotocol/v4/canonicalusdc/state`;
    return await this.req.get<QueryStateResponseSDKType>(endpoint);
  }
  /* PendingSettlements returns a page of in-flight canonical settlements. */


  async pendingSettlements(params: QueryPendingSettlementsRequest = {
    pagination: undefined
  }): Promise<QueryPendingSettlementsResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/v4/canonicalusdc/pending_settlements`;
    return await this.req.get<QueryPendingSettlementsResponseSDKType>(endpoint, options);
  }
  /* PendingSettlement returns the in-flight canonical settlement for an IBC
   packet identified by source channel and sequence. */


  async pendingSettlement(params: QueryPendingSettlementRequest): Promise<QueryPendingSettlementResponseSDKType> {
    const endpoint = `dydxprotocol/v4/canonicalusdc/pending_settlements/${params.sourceChannel}/${params.sequence}`;
    return await this.req.get<QueryPendingSettlementResponseSDKType>(endpoint);
  }

}