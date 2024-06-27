import { LCDClient } from "@osmonauts/lcd";
import { QueryPerpetualFeeParamsRequest, QueryPerpetualFeeParamsResponseSDKType, QueryUserFeeTierRequest, QueryUserFeeTierResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    perpetualFeeParams(_params?: QueryPerpetualFeeParamsRequest): Promise<QueryPerpetualFeeParamsResponseSDKType>;
    userFeeTier(params: QueryUserFeeTierRequest): Promise<QueryUserFeeTierResponseSDKType>;
}
