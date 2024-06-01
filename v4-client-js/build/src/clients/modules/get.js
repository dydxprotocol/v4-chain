"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Get = void 0;
const stargate_1 = require("@cosmjs/stargate");
const AuthModule = __importStar(require("cosmjs-types/cosmos/auth/v1beta1/query"));
const BankModule = __importStar(require("cosmjs-types/cosmos/bank/v1beta1/query"));
const long_1 = __importDefault(require("long"));
const protobufjs_1 = __importDefault(require("protobufjs"));
const constants_1 = require("../constants");
const errors_1 = require("../lib/errors");
const proto_includes_1 = require("./proto-includes");
// Required for encoding and decoding queries that are of type Long.
// Must be done once but since the individal modules should be usable without
// dYdXClient - must be set in each module that encounters encoding/decoding Longs.
// Reference: https://github.com/protobufjs/protobuf.js/issues/921
protobufjs_1.default.util.Long = long_1.default;
protobufjs_1.default.configure();
class Get {
    constructor(tendermintClient, stargateQueryClient) {
        this.tendermintClient = tendermintClient;
        this.stargateQueryClient = stargateQueryClient;
    }
    /**
     * @description Get latest block
     *
     * @returns last block structure
     */
    async latestBlock() {
        return this.tendermintClient.getBlock();
    }
    /**
     * @description Get latest block height
     *
     * @returns last height
     */
    async latestBlockHeight() {
        const block = await this.latestBlock();
        return block.header.height;
    }
    /**
     * @description Get all fee tier params.
     *
     * @returns All fee tier params.
     */
    async getFeeTiers() {
        const requestData = Uint8Array.from(proto_includes_1.FeeTierModule.QueryPerpetualFeeParamsRequest.encode({})
            .finish());
        const data = await this.sendQuery('/dydxprotocol.feetiers.Query/PerpetualFeeParams', requestData);
        return proto_includes_1.FeeTierModule.QueryPerpetualFeeParamsResponse.decode(data);
    }
    /**
     * @description Get fee tier the user belongs to
     *
     * @returns the fee tier user belongs to.
     */
    async getUserFeeTier(address) {
        const requestData = Uint8Array.from(proto_includes_1.FeeTierModule.QueryUserFeeTierRequest.encode({ user: address })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.feetiers.Query/UserFeeTier', requestData);
        return proto_includes_1.FeeTierModule.QueryUserFeeTierResponse.decode(data);
    }
    /**
     * @description Get get trading stats
     *
     * @returns return the user's taker and maker volume
     */
    async getUserStats(address) {
        const requestData = Uint8Array.from(proto_includes_1.StatsModule.QueryUserStatsRequest.encode({ user: address })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.stats.Query/UserStats', requestData);
        return proto_includes_1.StatsModule.QueryUserStatsResponse.decode(data).stats;
    }
    /**
     * @description Get all balances for an account.
     *
     * @returns Array of Coin balances for all tokens held by an account.
     */
    async getAccountBalances(address) {
        const requestData = Uint8Array.from(BankModule.QueryAllBalancesRequest.encode({ address })
            .finish());
        const data = await this.sendQuery('/cosmos.bank.v1beta1.Query/AllBalances', requestData);
        return BankModule.QueryAllBalancesResponse.decode(data).balances;
    }
    /**
     * @description Get balances of one denom for an account.
     *
     * @returns Coin balance for denom tokens held by an account.
     */
    async getAccountBalance(address, denom) {
        const requestData = Uint8Array.from(BankModule.QueryBalanceRequest.encode({
            address,
            denom,
        })
            .finish());
        const data = await this.sendQuery('/cosmos.bank.v1beta1.Query/Balance', requestData);
        const coin = BankModule.QueryBalanceResponse.decode(data).balance;
        return coin;
    }
    /**
     * @description Get all subaccounts
     *
     * @returns All subaccounts
     */
    async getSubaccounts() {
        const requestData = Uint8Array.from(proto_includes_1.SubaccountsModule.QueryAllSubaccountRequest.encode({})
            .finish());
        const data = await this.sendQuery('/dydxprotocol.subaccounts.Query/SubaccountAll', requestData);
        return proto_includes_1.SubaccountsModule.QuerySubaccountAllResponse.decode(data);
    }
    /**
     * @description Get a specific subaccount for an account.
     *
     * @returns Subaccount for account with given accountNumber or default subaccount if none exists.
     */
    async getSubaccount(address, accountNumber) {
        const requestData = Uint8Array.from(proto_includes_1.SubaccountsModule.QueryGetSubaccountRequest.encode({
            owner: address,
            number: accountNumber,
        })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.subaccounts.Query/Subaccount', requestData);
        return proto_includes_1.SubaccountsModule.QuerySubaccountResponse.decode(data);
    }
    /**
     * @description Get the params for the rewards module.
     *
     * @returns Params for the rewards module.
     */
    async getRewardsParams() {
        const requestData = Uint8Array.from(proto_includes_1.RewardsModule.QueryParamsRequest.encode({})
            .finish());
        const data = await this.sendQuery('/dydxprotocol.rewards.Query/Params', requestData);
        return proto_includes_1.RewardsModule.QueryParamsResponse.decode(data);
    }
    /**
     * @description Get all Clob Pairs.
     *
     * @returns Information on all Clob Pairs.
     */
    async getAllClobPairs() {
        const requestData = Uint8Array.from(proto_includes_1.ClobModule.QueryAllClobPairRequest.encode({ pagination: constants_1.PAGE_REQUEST })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.clob.Query/ClobPairAll', requestData);
        return proto_includes_1.ClobModule.QueryClobPairAllResponse.decode(data);
    }
    /**
     * @description Get Clob Pair for an Id or the promise is rejected if no pair exists.
     *
     * @returns Clob Pair for a given Clob Pair Id.
     */
    async getClobPair(pairId) {
        const requestData = Uint8Array.from(proto_includes_1.ClobModule.QueryGetClobPairRequest.encode({ id: pairId })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.clob.Query/ClobPair', requestData);
        return proto_includes_1.ClobModule.QueryClobPairResponse.decode(data);
    }
    /**
     * @description Get all Prices across markets.
     *
     * @returns Prices across all markets.
     */
    async getAllPrices() {
        const requestData = Uint8Array.from(proto_includes_1.PricesModule.QueryAllMarketPricesRequest.encode({ pagination: constants_1.PAGE_REQUEST })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.prices.Query/AllMarketPrices', requestData);
        return proto_includes_1.PricesModule.QueryAllMarketPricesResponse.decode(data);
    }
    /**
     * @description Get Price for a clob Id or the promise is rejected if none exists.
     *
     * @returns Price for a given Market Id.
     */
    async getPrice(marketId) {
        const requestData = Uint8Array.from(proto_includes_1.PricesModule.QueryMarketPriceRequest.encode({ id: marketId })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.prices.Query/MarketPrice', requestData);
        return proto_includes_1.PricesModule.QueryMarketPriceResponse.decode(data);
    }
    /**
     * @description Get all Perpetuals.
     *
     * @returns Information on all Perpetual pairs.
     */
    async getAllPerpetuals() {
        const requestData = Uint8Array.from(proto_includes_1.PerpetualsModule.QueryAllPerpetualsRequest.encode({ pagination: constants_1.PAGE_REQUEST })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.perpetuals.Query/AllPerpetuals', requestData);
        return proto_includes_1.PerpetualsModule.QueryAllPerpetualsResponse.decode(data);
    }
    /**
     * @description Get Perpetual for an Id or the promise is rejected if none exists.
     *
     * @returns The Perpetual for a given Perpetual Id.
     */
    async getPerpetual(perpetualId) {
        const requestData = Uint8Array.from(proto_includes_1.PerpetualsModule.QueryPerpetualRequest.encode({ id: perpetualId })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.perpetuals.Query/Perpetual', requestData);
        return proto_includes_1.PerpetualsModule.QueryPerpetualResponse.decode(data);
    }
    /**
     * @description Get Account for an address or the promise is rejected if the account
     * does not exist on-chain.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error.
     * @returns An account for a given address.
     */
    async getAccount(address) {
        const requestData = Uint8Array.from(AuthModule.QueryAccountRequest.encode({ address })
            .finish());
        const data = await this.sendQuery('/cosmos.auth.v1beta1.Query/Account', requestData);
        const rawAccount = AuthModule.QueryAccountResponse.decode(data).account;
        // The promise should have been rejected if the rawAccount was undefined.
        if (rawAccount === undefined) {
            throw new errors_1.UnexpectedClientError();
        }
        return (0, stargate_1.accountFromAny)(rawAccount);
    }
    /**
     * @description Get equity tier limit configuration.
     *
     * @returns Information on all equity tiers that are configured.
     */
    async getEquityTierLimitConfiguration() {
        const requestData = Uint8Array.from(proto_includes_1.ClobModule.QueryEquityTierLimitConfigurationRequest.encode({})
            .finish());
        const data = await this.sendQuery('/dydxprotocol.clob.Query/EquityTierLimitConfiguration', requestData);
        return proto_includes_1.ClobModule.QueryEquityTierLimitConfigurationResponse.decode(data);
    }
    /**
     *
     * @description Get all delegations from a delegator.
     *
     * @returns All delegations from a delegator.
     */
    async getDelegatorDelegations(delegatorAddr) {
        const requestData = Uint8Array.from(proto_includes_1.StakingModule.QueryDelegatorDelegationsRequest.encode({
            delegatorAddr,
            pagination: constants_1.PAGE_REQUEST,
        })
            .finish());
        const data = await this.sendQuery('/cosmos.staking.v1beta1.Query/DelegatorDelegations', requestData);
        return proto_includes_1.StakingModule.QueryDelegatorDelegationsResponse.decode(data);
    }
    /**
     *
     * @description Get all unbonding delegations from a delegator.
     *
     * @returns All unbonding delegations from a delegator.
     */
    async getDelegatorUnbondingDelegations(delegatorAddr) {
        const requestData = Uint8Array.from(proto_includes_1.StakingModule.QueryDelegatorUnbondingDelegationsRequest.encode({
            delegatorAddr,
            pagination: constants_1.PAGE_REQUEST,
        })
            .finish());
        const data = await this.sendQuery('/cosmos.staking.v1beta1.Query/DelegatorUnbondingDelegations', requestData);
        return proto_includes_1.StakingModule.QueryDelegatorUnbondingDelegationsResponse.decode(data);
    }
    /**
     * @description Get all delayed complete bridge messages, optionally filtered by address.
     *
     * @returns Information on all delayed complete bridge messages.
     */
    async getDelayedCompleteBridgeMessages(address = '') {
        const requestData = Uint8Array.from(proto_includes_1.BridgeModule.QueryDelayedCompleteBridgeMessagesRequest.encode({ address })
            .finish());
        const data = await this.sendQuery('/dydxprotocol.bridge.Query/DelayedCompleteBridgeMessages', requestData);
        return proto_includes_1.BridgeModule.QueryDelayedCompleteBridgeMessagesResponse.decode(data);
    }
    /**
     * @description Get all validators of a status.
     *
     * @returns all validators of a status.
     */
    async getAllValidators(status = '') {
        const requestData = Uint8Array.from(proto_includes_1.StakingModule.QueryValidatorsRequest
            .encode({
            status,
            pagination: constants_1.PAGE_REQUEST,
        })
            .finish());
        const data = await this.sendQuery('/cosmos.staking.v1beta1.Query/Validators', requestData);
        return proto_includes_1.StakingModule.QueryValidatorsResponse.decode(data);
    }
    async sendQuery(requestUrl, requestData) {
        const resp = await this.stargateQueryClient.queryAbci(requestUrl, requestData);
        return resp.value;
    }
}
exports.Get = Get;
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiZ2V0LmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vc3JjL2NsaWVudHMvbW9kdWxlcy9nZXQudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7Ozs7QUFDQSwrQ0FPMEI7QUFDMUIsbUZBQXFFO0FBQ3JFLG1GQUFxRTtBQUVyRSxnREFBd0I7QUFDeEIsNERBQWtDO0FBRWxDLDRDQUE0QztBQUM1QywwQ0FBc0Q7QUFDdEQscURBVTBCO0FBRzFCLG9FQUFvRTtBQUNwRSw2RUFBNkU7QUFDN0UsbUZBQW1GO0FBQ25GLGtFQUFrRTtBQUNsRSxvQkFBUSxDQUFDLElBQUksQ0FBQyxJQUFJLEdBQUcsY0FBSSxDQUFDO0FBQzFCLG9CQUFRLENBQUMsU0FBUyxFQUFFLENBQUM7QUFFckIsTUFBYSxHQUFHO0lBSWQsWUFDRSxnQkFBa0MsRUFDbEMsbUJBQXdEO1FBRXhELElBQUksQ0FBQyxnQkFBZ0IsR0FBRyxnQkFBZ0IsQ0FBQztRQUN6QyxJQUFJLENBQUMsbUJBQW1CLEdBQUcsbUJBQW1CLENBQUM7SUFDakQsQ0FBQztJQUVEOzs7O09BSUc7SUFDSCxLQUFLLENBQUMsV0FBVztRQUNmLE9BQU8sSUFBSSxDQUFDLGdCQUFnQixDQUFDLFFBQVEsRUFBRSxDQUFDO0lBQzFDLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsS0FBSyxDQUFDLGlCQUFpQjtRQUNyQixNQUFNLEtBQUssR0FBRyxNQUFNLElBQUksQ0FBQyxXQUFXLEVBQUUsQ0FBQztRQUN2QyxPQUFPLEtBQUssQ0FBQyxNQUFNLENBQUMsTUFBTSxDQUFDO0lBQzdCLENBQUM7SUFFRDs7OztPQUlHO0lBQ0gsS0FBSyxDQUFDLFdBQVc7UUFDZixNQUFNLFdBQVcsR0FBRyxVQUFVLENBQUMsSUFBSSxDQUNqQyw4QkFBYSxDQUFDLDhCQUE4QixDQUFDLE1BQU0sQ0FBQyxFQUFFLENBQUM7YUFDcEQsTUFBTSxFQUFFLENBQ1osQ0FBQztRQUVGLE1B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