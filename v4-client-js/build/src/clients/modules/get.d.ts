import { Coin } from '@cosmjs/proto-signing';
import { Account, Block, QueryClient as StargateQueryClient, TxExtension } from '@cosmjs/stargate';
import Long from 'long';
import { BridgeModule, ClobModule, FeeTierModule, PerpetualsModule, PricesModule, RewardsModule, StakingModule, SubaccountsModule } from './proto-includes';
import { TendermintClient } from './tendermintClient';
export declare class Get {
    readonly tendermintClient: TendermintClient;
    readonly stargateQueryClient: (StargateQueryClient & TxExtension);
    constructor(tendermintClient: TendermintClient, stargateQueryClient: (StargateQueryClient & TxExtension));
    /**
     * @description Get latest block
     *
     * @returns last block structure
     */
    latestBlock(): Promise<Block>;
    /**
     * @description Get latest block height
     *
     * @returns last height
     */
    latestBlockHeight(): Promise<number>;
    /**
     * @description Get all fee tier params.
     *
     * @returns All fee tier params.
     */
    getFeeTiers(): Promise<FeeTierModule.QueryPerpetualFeeParamsResponse>;
    /**
     * @description Get fee tier the user belongs to
     *
     * @returns the fee tier user belongs to.
     */
    getUserFeeTier(address: string): Promise<FeeTierModule.QueryUserFeeTierResponse>;
    /**
     * @description Get get trading stats
     *
     * @returns return the user's taker and maker volume
     */
    getUserStats(address: string): Promise<{
        takerNotional: Long;
        makerNotional: Long;
    } | undefined>;
    /**
     * @description Get all balances for an account.
     *
     * @returns Array of Coin balances for all tokens held by an account.
     */
    getAccountBalances(address: string): Promise<Coin[]>;
    /**
     * @description Get balances of one denom for an account.
     *
     * @returns Coin balance for denom tokens held by an account.
     */
    getAccountBalance(address: string, denom: string): Promise<Coin | undefined>;
    /**
     * @description Get all subaccounts
     *
     * @returns All subaccounts
     */
    getSubaccounts(): Promise<SubaccountsModule.QuerySubaccountAllResponse>;
    /**
     * @description Get a specific subaccount for an account.
     *
     * @returns Subaccount for account with given accountNumber or default subaccount if none exists.
     */
    getSubaccount(address: string, accountNumber: number): Promise<SubaccountsModule.QuerySubaccountResponse>;
    /**
     * @description Get the params for the rewards module.
     *
     * @returns Params for the rewards module.
     */
    getRewardsParams(): Promise<RewardsModule.QueryParamsResponse>;
    /**
     * @description Get all Clob Pairs.
     *
     * @returns Information on all Clob Pairs.
     */
    getAllClobPairs(): Promise<ClobModule.QueryClobPairAllResponse>;
    /**
     * @description Get Clob Pair for an Id or the promise is rejected if no pair exists.
     *
     * @returns Clob Pair for a given Clob Pair Id.
     */
    getClobPair(pairId: number): Promise<ClobModule.QueryClobPairResponse>;
    /**
     * @description Get all Prices across markets.
     *
     * @returns Prices across all markets.
     */
    getAllPrices(): Promise<PricesModule.QueryAllMarketPricesResponse>;
    /**
     * @description Get Price for a clob Id or the promise is rejected if none exists.
     *
     * @returns Price for a given Market Id.
     */
    getPrice(marketId: number): Promise<PricesModule.QueryMarketPriceResponse>;
    /**
     * @description Get all Perpetuals.
     *
     * @returns Information on all Perpetual pairs.
     */
    getAllPerpetuals(): Promise<PerpetualsModule.QueryAllPerpetualsResponse>;
    /**
     * @description Get Perpetual for an Id or the promise is rejected if none exists.
     *
     * @returns The Perpetual for a given Perpetual Id.
     */
    getPerpetual(perpetualId: number): Promise<PerpetualsModule.QueryPerpetualResponse>;
    /**
     * @description Get Account for an address or the promise is rejected if the account
     * does not exist on-chain.
     *
     * @throws UnexpectedClientError if a malformed response is returned with no GRPC error.
     * @returns An account for a given address.
     */
    getAccount(address: string): Promise<Account>;
    /**
     * @description Get equity tier limit configuration.
     *
     * @returns Information on all equity tiers that are configured.
     */
    getEquityTierLimitConfiguration(): Promise<ClobModule.QueryEquityTierLimitConfigurationResponse>;
    /**
     *
     * @description Get all delegations from a delegator.
     *
     * @returns All delegations from a delegator.
     */
    getDelegatorDelegations(delegatorAddr: string): Promise<StakingModule.QueryDelegatorDelegationsResponse>;
    /**
     *
     * @description Get all unbonding delegations from a delegator.
     *
     * @returns All unbonding delegations from a delegator.
     */
    getDelegatorUnbondingDelegations(delegatorAddr: string): Promise<StakingModule.QueryDelegatorUnbondingDelegationsResponse>;
    /**
     * @description Get all delayed complete bridge messages, optionally filtered by address.
     *
     * @returns Information on all delayed complete bridge messages.
     */
    getDelayedCompleteBridgeMessages(address?: string): Promise<BridgeModule.QueryDelayedCompleteBridgeMessagesResponse>;
    /**
     * @description Get all validators of a status.
     *
     * @returns all validators of a status.
     */
    getAllValidators(status?: string): Promise<StakingModule.QueryValidatorsResponse>;
    private sendQuery;
}
