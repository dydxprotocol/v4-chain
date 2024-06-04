import { IndexerConfig, DEFAULT_API_TIMEOUT } from './constants';
import AccountClient from './modules/account';
import MarketsClient from './modules/markets';
import UtilityClient from './modules/utility';

/**
 * @description Client for Indexer
 */
export class IndexerClient {
    public readonly config: IndexerConfig;
    readonly apiTimeout: number;
    readonly _markets: MarketsClient;
    readonly _account: AccountClient;
    readonly _utility: UtilityClient;

    constructor(config: IndexerConfig, apiTimeout?: number) {
      this.config = config;
      this.apiTimeout = apiTimeout ?? DEFAULT_API_TIMEOUT;

      this._markets = new MarketsClient(config.restEndpoint);
      this._account = new AccountClient(config.restEndpoint);
      this._utility = new UtilityClient(config.restEndpoint);
    }

    /**
     * @description Get the public module, used for interacting with public endpoints.
     *
     * @returns The public module
     */
    get markets(): MarketsClient {
      return this._markets;
    }

    /**
     * @description Get the private module, used for interacting with private endpoints.
     *
     * @returns The private module
     */
    get account(): AccountClient {
      return this._account;
    }

    /**
     * @description Get the utility module, used for interacting with non-market public endpoints.
     */
    get utility(): UtilityClient {
      return this._utility;
    }
}
