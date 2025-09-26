import { setPaginationParams } from "../../helpers";
import { LCDClient } from "@osmonauts/lcd";
import { QueryGetClobPairRequest, QueryClobPairResponseSDKType, QueryAllClobPairRequest, QueryClobPairAllResponseSDKType, QueryEquityTierLimitConfigurationRequest, QueryEquityTierLimitConfigurationResponseSDKType, QueryBlockRateLimitConfigurationRequest, QueryBlockRateLimitConfigurationResponseSDKType, QueryLiquidationsConfigurationRequest, QueryLiquidationsConfigurationResponseSDKType, QueryNextClobPairIdRequest, QueryNextClobPairIdResponseSDKType, QueryLeverageRequest, QueryLeverageResponseSDKType } from "./query";
export class LCDQueryClient {
  req: LCDClient;

  constructor({
    requestClient
  }: {
    requestClient: LCDClient;
  }) {
    this.req = requestClient;
    this.clobPair = this.clobPair.bind(this);
    this.clobPairAll = this.clobPairAll.bind(this);
    this.equityTierLimitConfiguration = this.equityTierLimitConfiguration.bind(this);
    this.blockRateLimitConfiguration = this.blockRateLimitConfiguration.bind(this);
    this.liquidationsConfiguration = this.liquidationsConfiguration.bind(this);
    this.nextClobPairId = this.nextClobPairId.bind(this);
    this.leverage = this.leverage.bind(this);
  }
  /* Queries a ClobPair by id. */


  async clobPair(params: QueryGetClobPairRequest): Promise<QueryClobPairResponseSDKType> {
    const endpoint = `dydxprotocol/clob/clob_pair/${params.id}`;
    return await this.req.get<QueryClobPairResponseSDKType>(endpoint);
  }
  /* Queries a list of ClobPair items. */


  async clobPairAll(params: QueryAllClobPairRequest = {
    pagination: undefined
  }): Promise<QueryClobPairAllResponseSDKType> {
    const options: any = {
      params: {}
    };

    if (typeof params?.pagination !== "undefined") {
      setPaginationParams(options, params.pagination);
    }

    const endpoint = `dydxprotocol/clob/clob_pair`;
    return await this.req.get<QueryClobPairAllResponseSDKType>(endpoint, options);
  }
  /* Queries EquityTierLimitConfiguration. */


  async equityTierLimitConfiguration(_params: QueryEquityTierLimitConfigurationRequest = {}): Promise<QueryEquityTierLimitConfigurationResponseSDKType> {
    const endpoint = `dydxprotocol/clob/equity_tier`;
    return await this.req.get<QueryEquityTierLimitConfigurationResponseSDKType>(endpoint);
  }
  /* Queries BlockRateLimitConfiguration. */


  async blockRateLimitConfiguration(_params: QueryBlockRateLimitConfigurationRequest = {}): Promise<QueryBlockRateLimitConfigurationResponseSDKType> {
    const endpoint = `dydxprotocol/clob/block_rate`;
    return await this.req.get<QueryBlockRateLimitConfigurationResponseSDKType>(endpoint);
  }
  /* Queries LiquidationsConfiguration. */


  async liquidationsConfiguration(_params: QueryLiquidationsConfigurationRequest = {}): Promise<QueryLiquidationsConfigurationResponseSDKType> {
    const endpoint = `dydxprotocol/clob/liquidations_config`;
    return await this.req.get<QueryLiquidationsConfigurationResponseSDKType>(endpoint);
  }
  /* Queries the next clob pair id. */


  async nextClobPairId(_params: QueryNextClobPairIdRequest = {}): Promise<QueryNextClobPairIdResponseSDKType> {
    const endpoint = `dydxprotocol/clob/next_clob_pair_id`;
    return await this.req.get<QueryNextClobPairIdResponseSDKType>(endpoint);
  }
  /* Queries leverage for a subaccount. */


  async leverage(params: QueryLeverageRequest): Promise<QueryLeverageResponseSDKType> {
    const endpoint = `dydxprotocol/clob/leverage/${params.owner}/${params.number}`;
    return await this.req.get<QueryLeverageResponseSDKType>(endpoint);
  }

}