import { QueryAllowanceResponse, QueryAllowancesResponse } from "cosmjs-types/cosmos/feegrant/v1beta1/query";
import { QueryClient } from "../../queryclient";
export interface FeegrantExtension {
    readonly feegrant: {
        readonly allowance: (granter: string, grantee: string) => Promise<QueryAllowanceResponse>;
        readonly allowances: (grantee: string, paginationKey?: Uint8Array) => Promise<QueryAllowancesResponse>;
    };
}
export declare function setupFeegrantExtension(base: QueryClient): FeegrantExtension;
