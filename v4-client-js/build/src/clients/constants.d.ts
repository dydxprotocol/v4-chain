import { PageRequest } from '@dydxprotocol/v4-proto/src/codegen/cosmos/base/query/v1beta1/pagination';
import { BroadcastOptions, DenomConfig } from './types';
export * from '../lib/constants';
export declare const DEV_CHAIN_ID = "dydxprotocol-testnet";
export declare const STAGING_CHAIN_ID = "dydxprotocol-testnet";
export declare const TESTNET_CHAIN_ID = "dydx-testnet-4";
export declare const LOCAL_CHAIN_ID = "consu";
export declare enum IndexerApiHost {
    TESTNET = "https://dydx-testnet.imperator.co",
    LOCAL = "http://localhost:3002"
}
export declare enum IndexerWSHost {
    TESTNET = "wss://dydx-testnet.imperator.co/v4/ws",
    LOCAL = "ws://localhost:3003"
}
export declare enum FaucetApiHost {
    TESTNET = "https://faucet.v4testnet.dydx.exchange"
}
export declare enum ValidatorApiHost {
    TESTNET = "https://test-dydx.kingnodes.com",
    LOCAL = "http://localhost:26657"
}
export declare enum NetworkId {
    TESTNET = "dydx-testnet-4"
}
export declare const NETWORK_ID_MAINNET: string | null;
export declare const NETWORK_ID_TESTNET: string;
export declare enum MarketStatisticDay {
    ONE = "1",
    SEVEN = "7",
    THIRTY = "30"
}
export declare enum OrderType {
    LIMIT = "LIMIT",
    MARKET = "MARKET",
    STOP_LIMIT = "STOP_LIMIT",
    TAKE_PROFIT_LIMIT = "TAKE_PROFIT",
    STOP_MARKET = "STOP_MARKET",
    TAKE_PROFIT_MARKET = "TAKE_PROFIT_MARKET"
}
export declare enum OrderSide {
    BUY = "BUY",
    SELL = "SELL"
}
export declare enum OrderTimeInForce {
    GTT = "GTT",
    IOC = "IOC",
    FOK = "FOK"
}
export declare enum OrderExecution {
    DEFAULT = "DEFAULT",
    IOC = "IOC",
    FOK = "FOK",
    POST_ONLY = "POST_ONLY"
}
export declare enum OrderStatus {
    BEST_EFFORT_OPENED = "BEST_EFFORT_OPENED",
    OPEN = "OPEN",
    FILLED = "FILLED",
    BEST_EFFORT_CANCELED = "BEST_EFFORT_CANCELED",
    CANCELED = "CANCELED"
}
export declare enum TickerType {
    PERPETUAL = "PERPETUAL"
}
export declare enum PositionStatus {
    OPEN = "OPEN",
    CLOSED = "CLOSED",
    LIQUIDATED = "LIQUIDATED"
}
export declare enum TimePeriod {
    ONE_DAY = "ONE_DAY",
    SEVEN_DAYS = "SEVEN_DAYS"
}
export declare const DEFAULT_API_TIMEOUT: number;
export declare const MAX_MEMO_CHARACTERS: number;
export declare const SHORT_BLOCK_WINDOW: number;
export declare const SHORT_BLOCK_FORWARD: number;
export declare const PAGE_REQUEST: PageRequest;
export declare class IndexerConfig {
    restEndpoint: string;
    websocketEndpoint: string;
    constructor(restEndpoint: string, websocketEndpoint: string);
}
export declare class ValidatorConfig {
    restEndpoint: string;
    chainId: string;
    denoms: DenomConfig;
    broadcastOptions?: BroadcastOptions;
    constructor(restEndpoint: string, chainId: string, denoms: DenomConfig, broadcastOptions?: BroadcastOptions);
}
export declare class Network {
    env: string;
    indexerConfig: IndexerConfig;
    validatorConfig: ValidatorConfig;
    constructor(env: string, indexerConfig: IndexerConfig, validatorConfig: ValidatorConfig);
    static testnet(): Network;
    static local(): Network;
    getString(): string;
}
