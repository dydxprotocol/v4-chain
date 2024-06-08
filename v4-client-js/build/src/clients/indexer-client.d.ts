import { IndexerConfig } from './constants';
import AccountClient from './modules/account';
import MarketsClient from './modules/markets';
import UtilityClient from './modules/utility';
/**
 * @description Client for Indexer
 */
export declare class IndexerClient {
    readonly config: IndexerConfig;
    readonly apiTimeout: number;
    readonly _markets: MarketsClient;
    readonly _account: AccountClient;
    readonly _utility: UtilityClient;
    constructor(config: IndexerConfig, apiTimeout?: number);
    /**
     * @description Get the public module, used for interacting with public endpoints.
     *
     * @returns The public module
     */
    get markets(): MarketsClient;
    /**
     * @description Get the private module, used for interacting with private endpoints.
     *
     * @returns The private module
     */
    get account(): AccountClient;
    /**
     * @description Get the utility module, used for interacting with non-market public endpoints.
     */
    get utility(): UtilityClient;
}
