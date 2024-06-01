import { LCDClient } from "@osmonauts/lcd";
import { QueryGetEpochInfoRequest, QueryEpochInfoResponseSDKType, QueryAllEpochInfoRequest, QueryEpochInfoAllResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    epochInfo(params: QueryGetEpochInfoRequest): Promise<QueryEpochInfoResponseSDKType>;
    epochInfoAll(params?: QueryAllEpochInfoRequest): Promise<QueryEpochInfoAllResponseSDKType>;
}
