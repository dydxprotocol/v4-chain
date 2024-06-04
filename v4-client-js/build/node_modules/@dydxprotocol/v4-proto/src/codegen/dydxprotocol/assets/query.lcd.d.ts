import { LCDClient } from "@osmonauts/lcd";
import { QueryAssetRequest, QueryAssetResponseSDKType, QueryAllAssetsRequest, QueryAllAssetsResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    asset(params: QueryAssetRequest): Promise<QueryAssetResponseSDKType>;
    allAssets(params?: QueryAllAssetsRequest): Promise<QueryAllAssetsResponseSDKType>;
}
