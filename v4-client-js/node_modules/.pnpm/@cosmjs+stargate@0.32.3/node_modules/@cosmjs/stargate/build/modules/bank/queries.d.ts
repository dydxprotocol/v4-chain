import { Metadata } from "cosmjs-types/cosmos/bank/v1beta1/bank";
import { QueryTotalSupplyResponse } from "cosmjs-types/cosmos/bank/v1beta1/query";
import { Coin } from "cosmjs-types/cosmos/base/v1beta1/coin";
import { QueryClient } from "../../queryclient";
export interface BankExtension {
    readonly bank: {
        readonly balance: (address: string, denom: string) => Promise<Coin>;
        readonly allBalances: (address: string) => Promise<Coin[]>;
        readonly totalSupply: (paginationKey?: Uint8Array) => Promise<QueryTotalSupplyResponse>;
        readonly supplyOf: (denom: string) => Promise<Coin>;
        readonly denomMetadata: (denom: string) => Promise<Metadata>;
        readonly denomsMetadata: () => Promise<Metadata[]>;
    };
}
export declare function setupBankExtension(base: QueryClient): BankExtension;
