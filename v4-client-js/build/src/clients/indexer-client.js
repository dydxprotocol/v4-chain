"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.IndexerClient = void 0;
const constants_1 = require("./constants");
const account_1 = __importDefault(require("./modules/account"));
const markets_1 = __importDefault(require("./modules/markets"));
const utility_1 = __importDefault(require("./modules/utility"));
/**
 * @description Client for Indexer
 */
class IndexerClient {
    constructor(config, apiTimeout) {
        this.config = config;
        this.apiTimeout = apiTimeout !== null && apiTimeout !== void 0 ? apiTimeout : constants_1.DEFAULT_API_TIMEOUT;
        this._markets = new markets_1.default(config.restEndpoint);
        this._account = new account_1.default(config.restEndpoint);
        this._utility = new utility_1.default(config.restEndpoint);
    }
    /**
     * @description Get the public module, used for interacting with public endpoints.
     *
     * @returns The public module
     */
    get markets() {
        return this._markets;
    }
    /**
     * @description Get the private module, used for interacting with private endpoints.
     *
     * @returns The private module
     */
    get account() {
        return this._account;
    }
    /**
     * @description Get the utility module, used for interacting with non-market public endpoints.
     */
    get utility() {
        return this._utility;
    }
}
exports.IndexerClient = IndexerClient;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW5kZXhlci1jbGllbnQuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi9zcmMvY2xpZW50cy9pbmRleGVyLWNsaWVudC50cyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiOzs7Ozs7QUFBQSwyQ0FBaUU7QUFDakUsZ0VBQThDO0FBQzlDLGdFQUE4QztBQUM5QyxnRUFBOEM7QUFFOUM7O0dBRUc7QUFDSCxNQUFhLGFBQWE7SUFPdEIsWUFBWSxNQUFxQixFQUFFLFVBQW1CO1FBQ3BELElBQUksQ0FBQyxNQUFNLEdBQUcsTUFBTSxDQUFDO1FBQ3JCLElBQUksQ0FBQyxVQUFVLEdBQUcsVUFBVSxhQUFWLFVBQVUsY0FBVixVQUFVLEdBQUksK0JBQW1CLENBQUM7UUFFcEQsSUFBSSxDQUFDLFFBQVEsR0FBRyxJQUFJLGlCQUFhLENBQUMsTUFBTSxDQUFDLFlBQVksQ0FBQyxDQUFDO1FBQ3ZELElBQUksQ0FBQyxRQUFRLEdBQUcsSUFBSSxpQkFBYSxDQUFDLE1BQU0sQ0FBQyxZQUFZLENBQUMsQ0FBQztRQUN2RCxJQUFJLENBQUMsUUFBUSxHQUFHLElBQUksaUJBQWEsQ0FBQyxNQUFNLENBQUMsWUFBWSxDQUFDLENBQUM7SUFDekQsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxJQUFJLE9BQU87UUFDVCxPQUFPLElBQUksQ0FBQyxRQUFRLENBQUM7SUFDdkIsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxJQUFJLE9BQU87UUFDVCxPQUFPLElBQUksQ0FBQyxRQUFRLENBQUM7SUFDdkIsQ0FBQztJQUVEOztPQUVHO0lBQ0gsSUFBSSxPQUFPO1FBQ1QsT0FBTyxJQUFJLENBQUMsUUFBUSxDQUFDO0lBQ3ZCLENBQUM7Q0FDSjtBQXhDRCxzQ0F3Q0MifQ==