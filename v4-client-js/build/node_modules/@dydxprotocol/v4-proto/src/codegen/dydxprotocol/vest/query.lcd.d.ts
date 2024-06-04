import { LCDClient } from "@osmonauts/lcd";
import { QueryVestEntryRequest, QueryVestEntryResponseSDKType } from "./query";
export declare class LCDQueryClient {
    req: LCDClient;
    constructor({ requestClient }: {
        requestClient: LCDClient;
    });
    vestEntry(params: QueryVestEntryRequest): Promise<QueryVestEntryResponseSDKType>;
}
