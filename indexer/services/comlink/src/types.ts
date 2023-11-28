import {
  APIOrderStatus,
  APITimeInForce,
  AssetFromDatabase,
  CandleColumns,
  CandleFromDatabase,
  CandleResolution,
  FillType,
  IsoString,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderType,
  PerpetualMarketStatus,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PositionSide,
  SubaccountFromDatabase,
  TradeType,
  TransferType,
} from '@dydxprotocol-indexer/postgres';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import express from 'express';

/* ------- GENERAL/UNCATEGORIZED TYPES ------- */

export interface ResponseWithBody extends express.Response {
  body: unknown
}

export enum RequestMethod {
  DELETE = 'DELETE',
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
}

/* ------- SUBACCOUNT TYPES ------- */

export interface SubaccountResponseObject {
  address: string,
  subaccountNumber: number,
  equity: string,
  freeCollateral: string,
  openPerpetualPositions: PerpetualPositionsMap,
  assetPositions: AssetPositionsMap,
  marginEnabled: boolean,
}

export type SubaccountById = {[id: string]: SubaccountFromDatabase};

/* ------- TIME TYPES ------- */

export interface TimeResponse {
  iso: IsoString,
  epoch: number,
}

/* ------- POSITION TYPES ------- */

export interface PerpetualPositionResponse {
  positions: PerpetualPositionResponseObject[];
}

export interface PerpetualPositionWithFunding extends PerpetualPositionFromDatabase {
  unsettledFunding: string;
}

export interface PerpetualPositionResponseObject {
  market: string;
  status: PerpetualPositionStatus;
  side: PositionSide;
  size: string;
  maxSize: string;
  entryPrice: string;
  realizedPnl: string;
  createdAt: IsoString;
  createdAtHeight: string;
  sumOpen: string;
  sumClose: string;
  netFunding: string;
  unrealizedPnl: string;
  closedAt?: IsoString | null;
  exitPrice?: string | null;
}

export type PerpetualPositionsMap = { [market: string]: PerpetualPositionResponseObject };

export interface AssetPositionResponse {
  positions: AssetPositionResponseObject[];
}

export interface AssetPositionResponseObject {
  symbol: string;
  side: PositionSide;
  size: string;
  assetId: string;
}

export type AssetPositionsMap = { [symbol: string]: AssetPositionResponseObject };

/* ------- FILL TYPES ------- */

export interface FillResponse {
  fills: FillResponseObject[],
}

export interface FillResponseObject {
  id: string,
  side: OrderSide,
  liquidity: Liquidity,
  type: FillType,
  market: string,
  marketType: MarketType,
  price: string,
  size: string,
  fee: string,
  createdAt: IsoString,
  createdAtHeight: string,
  orderId?: string,
  clientMetadata?: string,
}

/* ------- TRANSFER TYPES ------- */

export interface TransferResponse {
  transfers: TransferResponseObject[],
}

export interface TransferResponseObject {
  id: string,
  sender: {
    address: string,
    subaccountNumber?: number,
  },
  recipient: {
    address: string,
    subaccountNumber?: number,
  },
  size: string,
  createdAt: string,
  createdAtHeight: string,
  symbol: string,
  type: TransferType,
  transactionHash: string,
}

/* ------- PNL TICKS TYPES ------- */

export interface HistoricalPnlResponse {
  historicalPnl: PnlTicksResponseObject[],
}

export interface PnlTicksResponseObject {
  id: string,
  subaccountId: string,
  equity: string,
  totalPnl: string,
  netTransfers: string,
  createdAt: string,
  blockHeight: string,
  blockTime: IsoString,
}

/* ------- TRADE TYPES ------- */

export interface TradeResponse {
  trades: TradeResponseObject[],
}

export interface TradeResponseObject {
  id: string,
  side: OrderSide,
  size: string,
  price: string,
  type: TradeType,
  createdAt: IsoString,
  createdAtHeight: string,
}

/* ------- Height TYPES ------- */

export interface HeightResponse {
  height: string,
  time: IsoString,
}

/* ------- MARKET TYPES ------- */

export type AssetById = {[assetId: string]: AssetFromDatabase};

export interface MarketAndType {
  marketType: MarketType,
  market: string,
}

export type MarketAndTypeByClobPairId = {[clobPairId: string]: MarketAndType};

export enum MarketType {
  PERPETUAL = 'PERPETUAL',
  SPOT = 'SPOT',
}

export interface PerpetualMarketResponse {
  markets: {
    [ticker: string]: PerpetualMarketResponseObject,
  }
}

export interface PerpetualMarketResponseObject {
  clobPairId: string;
  ticker: string;
  status: PerpetualMarketStatus;
  oraclePrice: string;
  priceChange24H: string;
  volume24H: string;
  trades24H: number;
  nextFundingRate: string;
  initialMarginFraction: string;
  maintenanceMarginFraction: string;
  basePositionNotional: string;
  openInterest: string;
  atomicResolution: number;
  quantumConversionExponent: number;
  tickSize: string;
  stepSize: string;
  stepBaseQuantums: number;
  subticksPerTick: number;
}

/* ------- ORDERBOOK TYPES ------- */

export interface OrderbookResponseObject {
  bids: OrderbookResponsePriceLevel[],
  asks: OrderbookResponsePriceLevel[],
}

export interface OrderbookResponsePriceLevel {
  price: string,
  size: string,
}

/* ------- ORDER TYPES ------- */
// TimeInForce stored in the database is different from the TimeInForce expected in the API
// The omitted field name have to be literal strings for Typescript to parse them correctly
export interface OrderResponseObject extends Omit<OrderFromDatabase, 'timeInForce' | 'status' | 'updatedAt' | 'updatedAtHeight'> {
  timeInForce: APITimeInForce,
  status: APIOrderStatus,
  postOnly: boolean,
  ticker: string;
  updatedAt?: IsoString;
  updatedAtHeight?: string
}

export type RedisOrderMap = { [orderId: string]: RedisOrder };

export type PostgresOrderMap = { [orderId: string]: OrderFromDatabase };

/* ------- CANDLE TYPES ------- */

export interface CandleResponse {
  candles: CandleResponseObject[],
}

export interface CandleResponseObject extends Omit<CandleFromDatabase, CandleColumns.id> {}

/* ------- CANDLE TYPES ------- */

export interface SparklineResponseObject {
  [ticker: string]: string[],
}

export enum SparklineTimePeriod {
  ONE_DAY = 'ONE_DAY',
  SEVEN_DAYS = 'SEVEN_DAYS',
}

/* ------- FUNDING TYPES ------- */

export interface HistoricalFundingResponse {
  historicalFunding: HistoricalFundingResponseObject[],
}

export interface HistoricalFundingResponseObject {
  ticker: string,
  rate: string,
  price: string,
  effectiveAt: IsoString,
  effectiveAtHeight: string,
}

/* ------- GET REQUEST TYPES ------- */

export interface AddressRequest {
  address: string,
}

export interface SubaccountRequest extends AddressRequest {
  subaccountNumber: number,
}

export interface LimitRequest {
  limit: number,
}

export interface TickerRequest {
  ticker?: string,
}

export interface LimitAndCreatedBeforeRequest extends LimitRequest {
  createdBeforeOrAtHeight?: number,
  createdBeforeOrAt?: IsoString,
}

export interface LimitAndEffectiveBeforeRequest extends LimitRequest {
  effectiveBeforeOrAtHeight?: number,
  effectiveBeforeOrAt?: IsoString,
}

export interface LimitAndCreatedBeforeAndAfterRequest extends LimitAndCreatedBeforeRequest {
  createdOnOrAfterHeight?: number,
  createdOnOrAfter?: IsoString,
}

export interface PerpetualPositionRequest extends SubaccountRequest, LimitAndCreatedBeforeRequest {
  status: PerpetualPositionStatus[],
}

export interface AssetPositionRequest extends SubaccountRequest {}

export interface TransferRequest extends SubaccountRequest, LimitAndCreatedBeforeRequest {}

export interface FillRequest extends SubaccountRequest, LimitAndCreatedBeforeRequest {
  market: string,
  marketType: MarketType,
}

export interface TradeRequest extends LimitAndCreatedBeforeRequest {
  ticker: string,
}

export interface PerpetualMarketRequest extends LimitRequest, TickerRequest {}

export interface PnlTicksRequest extends SubaccountRequest, LimitAndCreatedBeforeAndAfterRequest {}

export interface OrderbookRequest {
  ticker: string,
}

export interface GetOrderRequest {
  orderId: string,
}

export interface ListOrderRequest extends SubaccountRequest, LimitRequest, TickerRequest {
  side?: OrderSide,
  type?: OrderType,
  status?: OrderStatus[],
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
  returnLatestOrders?: boolean,
}

export interface CandleRequest extends LimitRequest {
  ticker: string,
  resolution: CandleResolution,
  fromISO?: IsoString,
  toISO?: IsoString,
}

export interface SparklinesRequest {
  timePeriod: SparklineTimePeriod,
}

export interface HistoricalFundingRequest extends LimitAndEffectiveBeforeRequest {
  ticker: string,
}

/* ------- COLLATERALIZATION TYPES ------- */

export interface Risk {
  initial: Big;
  maintenance: Big;
}

/* ------- COMPLIANCE TYPES ------- */

export interface ComplianceResponse {
  restricted: boolean;
  reason?: string;
}

export interface ComplianceRequest extends AddressRequest {}

export enum BlockedCode {
  GEOBLOCKED = 'GEOBLOCKED',
  COMPLIANCE_BLOCKED = 'COMPLIANCE_BLOCKED',
}
