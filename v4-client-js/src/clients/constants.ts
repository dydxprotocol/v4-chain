import { PageRequest } from '@dydxprotocol/v4-proto/src/codegen/cosmos/base/query/v1beta1/pagination';
import Long from 'long';

import { BroadcastOptions, DenomConfig } from './types';

export * from '../lib/constants';

// Chain ID
export const DEV_CHAIN_ID = 'dydxprotocol-testnet';
export const STAGING_CHAIN_ID = 'dydxprotocol-testnet';
export const TESTNET_CHAIN_ID = 'dydx-testnet-4';
//export const LOCAL_CHAIN_ID = 'localdydxprotocol';
export const LOCAL_CHAIN_ID = 'consu';

// ------------ API URLs ------------
export enum IndexerApiHost {
  TESTNET = 'https://dydx-testnet.imperator.co',
  LOCAL = 'http://localhost:3002'
  // TODO: Add MAINNET
}

export enum IndexerWSHost {
  TESTNET = 'wss://dydx-testnet.imperator.co/v4/ws',
  // TODO: Add MAINNET
  LOCAL = 'ws://localhost:3003'
}

export enum FaucetApiHost {
  TESTNET = 'https://faucet.v4testnet.dydx.exchange',
}

export enum ValidatorApiHost {
  TESTNET = 'https://test-dydx.kingnodes.com',
  // TODO: Add MAINNET
  LOCAL = 'http://localhost:26657'
}

// ------------ Network IDs ------------

export enum NetworkId {
  TESTNET = 'dydx-testnet-4',
  // TODO: Add MAINNET
}
export const NETWORK_ID_MAINNET: string | null = null;
export const NETWORK_ID_TESTNET: string = 'dydxprotocol-testnet';

// ------------ Market Statistic Day Types ------------
export enum MarketStatisticDay {
  ONE = '1',
  SEVEN = '7',
  THIRTY = '30',
}

// ------------ Order Types ------------
// This should match OrderType in Abacus
export enum OrderType {
  LIMIT = 'LIMIT',
  MARKET = 'MARKET',
  STOP_LIMIT = 'STOP_LIMIT',
  TAKE_PROFIT_LIMIT = 'TAKE_PROFIT',
  STOP_MARKET = 'STOP_MARKET',
  TAKE_PROFIT_MARKET = 'TAKE_PROFIT_MARKET',
}

// ------------ Order Side ------------
// This should match OrderSide in Abacus
export enum OrderSide {
  BUY = 'BUY',
  SELL = 'SELL',
}

// ------------ Order TimeInForce ------------
// This should match OrderTimeInForce in Abacus
export enum OrderTimeInForce {
  GTT = 'GTT',
  IOC = 'IOC',
  FOK = 'FOK',
}

// ------------ Order Execution ------------
// This should match OrderExecution in Abacus
export enum OrderExecution {
  DEFAULT = 'DEFAULT',
  IOC = 'IOC',
  FOK = 'FOK',
  POST_ONLY = 'POST_ONLY',
}

// ------------ Order Status ------------
// This should match OrderStatus in Abacus
export enum OrderStatus {
  BEST_EFFORT_OPENED = 'BEST_EFFORT_OPENED',
  OPEN = 'OPEN',
  FILLED = 'FILLED',
  BEST_EFFORT_CANCELED = 'BEST_EFFORT_CANCELED',
  CANCELED = 'CANCELED',
}

export enum TickerType {
  PERPETUAL = 'PERPETUAL',  // Only PERPETUAL is supported right now
}

export enum PositionStatus {
  OPEN = 'OPEN',
  CLOSED = 'CLOSED',
  LIQUIDATED = 'LIQUIDATED',
}

// ----------- Time Period for Sparklines -------------

export enum TimePeriod {
  ONE_DAY = 'ONE_DAY',
  SEVEN_DAYS = 'SEVEN_DAYS',
}

// ------------ API Defaults ------------
export const DEFAULT_API_TIMEOUT: number = 3_000;

export const MAX_MEMO_CHARACTERS: number = 256;

export const SHORT_BLOCK_WINDOW: number = 20;

export const SHORT_BLOCK_FORWARD: number = 3;

// Querying
export const PAGE_REQUEST: PageRequest = {
  key: new Uint8Array(),
  offset: Long.UZERO,
  limit: Long.MAX_UNSIGNED_VALUE,
  countTotal: true,
  reverse: false,
};

export class IndexerConfig {
    public restEndpoint: string;
    public websocketEndpoint: string;

    constructor(restEndpoint: string,
      websocketEndpoint: string) {
      this.restEndpoint = restEndpoint;
      this.websocketEndpoint = websocketEndpoint;
    }
}

export class ValidatorConfig {
  public restEndpoint: string;
  public chainId: string;
  public denoms: DenomConfig;
  public broadcastOptions?: BroadcastOptions;

  constructor(
    restEndpoint: string,
    chainId: string,
    denoms: DenomConfig,
    broadcastOptions?: BroadcastOptions,
  ) {
    this.restEndpoint = restEndpoint?.endsWith('/') ? restEndpoint.slice(0, -1) : restEndpoint;
    this.chainId = chainId;

    this.denoms = denoms;
    this.broadcastOptions = broadcastOptions;
  }
}

export class Network {
  constructor(
    public env: string,
    public indexerConfig: IndexerConfig,
    public validatorConfig: ValidatorConfig,
  ) {}

  static testnet(): Network {
    const indexerConfig = new IndexerConfig(
      IndexerApiHost.TESTNET,
      IndexerWSHost.TESTNET,
    );
    const validatorConfig = new ValidatorConfig(ValidatorApiHost.TESTNET, TESTNET_CHAIN_ID,
      {
        CHAINTOKEN_DENOM: 'adv4tnt',
        TDAI_DENOM: 'utdai',
        TDAI_GAS_DENOM: 'utdai',
        TDAI_DECIMALS: 6,
        CHAINTOKEN_DECIMALS: 18,
      });
    
    return new Network('testnet', indexerConfig, validatorConfig);
  }

  static local(): Network {
    const indexerConfig = new IndexerConfig(
      IndexerApiHost.LOCAL,
      IndexerWSHost.LOCAL,
    );
    const validatorConfig = new ValidatorConfig(ValidatorApiHost.LOCAL, LOCAL_CHAIN_ID,
      {
        CHAINTOKEN_DENOM: 'adv4tnt',
        TDAI_DENOM: 'utdai',
        TDAI_GAS_DENOM: 'utdai',
        TDAI_DECIMALS: 6,
        CHAINTOKEN_DECIMALS: 18,
      });
    console.log("LOGGING NEW NETWORK");
    return new Network('local', indexerConfig, validatorConfig);
  }

  // TODO: Add mainnet(): Network

  getString(): string {
    return this.env;
  }
}
