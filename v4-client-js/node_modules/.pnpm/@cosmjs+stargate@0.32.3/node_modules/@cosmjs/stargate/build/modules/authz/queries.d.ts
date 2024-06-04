import { QueryGranteeGrantsResponse, QueryGranterGrantsResponse, QueryGrantsResponse } from "cosmjs-types/cosmos/authz/v1beta1/query";
import { QueryClient } from "../../queryclient";
export interface AuthzExtension {
    readonly authz: {
        readonly grants: (granter: string, grantee: string, msgTypeUrl: string, paginationKey?: Uint8Array) => Promise<QueryGrantsResponse>;
        readonly granteeGrants: (grantee: string, paginationKey?: Uint8Array) => Promise<QueryGranteeGrantsResponse>;
        readonly granterGrants: (granter: string, paginationKey?: Uint8Array) => Promise<QueryGranterGrantsResponse>;
    };
}
export declare function setupAuthzExtension(base: QueryClient): AuthzExtension;
