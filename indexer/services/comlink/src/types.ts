import {
  APIOrderStatus,
  APITimeInForce,
  AssetFromDatabase,
  CandleColumns,
  CandleFromDatabase,
  CandleResolution,
  ComplianceReason,
  ComplianceStatus,
  FillType,
  IsoString,
  Liquidity,
  OrderFromDatabase,
  OrderSide,
  OrderStatus,
  OrderType,
  PerpetualMarketStatus,
  PerpetualMarketType,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PnlFromDatabase,
  PnlTickInterval,
  PnlTicksFromDatabase,
  PositionSide,
  SubaccountFromDatabase,
  TradeType,
  TradingRewardAggregationPeriod,
  TransferType,
  VaultFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import { TurnkeyApiTypes } from '@turnkey/sdk-server';
import Big from 'big.js';
import express from 'express';

/* ------- GENERAL/UNCATEGORIZED TYPES ------- */

export interface ResponseWithBody extends express.Response {
  body: unknown,
}

export enum RequestMethod {
  DELETE = 'DELETE',
  GET = 'GET',
  POST = 'POST',
  PUT = 'PUT',
}

/* ------- Pagination ------- */
export interface PaginationResponse {
  /**
   * @isInt
   */
  pageSize?: number,
  /**
   * @isInt
   */
  totalResults?: number,
  /**
   * @isInt
   */
  offset?: number,
}

/* ------- SUBACCOUNT TYPES ------- */

export interface AddressResponse {
  subaccounts: SubaccountResponseObject[],
  totalTradingRewards: string,
}

export interface SubaccountResponseObject {
  address: string,
  /**
   * @isInt
   */
  subaccountNumber: number,
  equity: string,
  freeCollateral: string,
  openPerpetualPositions: PerpetualPositionsMap,
  assetPositions: AssetPositionsMap,
  marginEnabled: boolean,
  updatedAtHeight: string,
  latestProcessedBlockHeight: string,
}

export interface ParentSubaccountResponse {
  address: string,
  /**
   * @isInt
   */
  parentSubaccountNumber: number,
  equity: string, // aggregated over all child subaccounts
  freeCollateral: string, // aggregated over all child subaccounts
  childSubaccounts: SubaccountResponseObject[],
}

export type SubaccountById = { [id: string]: SubaccountFromDatabase };

/* ------- TIME TYPES ------- */

export interface TimeResponse {
  iso: IsoString,
  epoch: number,
}

/* ------- POSITION TYPES ------- */

export interface PerpetualPositionResponse {
  positions: PerpetualPositionResponseObject[],
}

export interface PerpetualPositionWithFunding extends PerpetualPositionFromDatabase {
  unsettledFunding: string,
}

export interface PerpetualPositionResponseObject {
  market: string,
  status: PerpetualPositionStatus,
  side: PositionSide,
  size: string,
  maxSize: string,
  entryPrice: string,
  realizedPnl: string,
  createdAt: IsoString,
  createdAtHeight: string,
  sumOpen: string,
  sumClose: string,
  netFunding: string,
  unrealizedPnl: string,
  closedAt?: IsoString | null,
  exitPrice?: string | null,
  /**
   * @isInt
   */
  subaccountNumber: number,
}

export type PerpetualPositionsMap = { [market: string]: PerpetualPositionResponseObject };

export interface AssetPositionResponse {
  positions: AssetPositionResponseObject[],
}

export interface AssetPositionResponseObject {
  symbol: string,
  side: PositionSide,
  size: string,
  assetId: string,
  /**
   * @isInt
   */
  subaccountNumber: number,
}

export type AssetPositionsMap = { [symbol: string]: AssetPositionResponseObject };

/* ------- FILL TYPES ------- */

export interface FillResponse extends PaginationResponse {
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
  affiliateRevShare: string,
  createdAt: IsoString,
  createdAtHeight: string,
  orderId?: string,
  clientMetadata?: string,
  /**
   * @isInt
   */
  subaccountNumber: number,
  builderFee?: string,
  builderAddress?: string,
  orderRouterAddress?: string,
  orderRouterFee?: string,
  positionSizeBefore?: string,
  entryPriceBefore?: string,
  positionSideBefore?: string,
}

/* ------- TRANSFER TYPES ------- */

export interface TransferResponse extends PaginationResponse {
  transfers: TransferResponseObject[],
}

export interface TransferResponseObject {
  id: string,
  sender: {
    address: string,
    /**
     * @isInt
     */
    subaccountNumber?: number,
  },
  recipient: {
    address: string,
    /**
     * @isInt
     */
    subaccountNumber?: number,
  },
  size: string,
  createdAt: string,
  createdAtHeight: string,
  symbol: string,
  type: TransferType,
  transactionHash: string,
}

export interface ParentSubaccountTransferResponse extends PaginationResponse {
  transfers: TransferResponseObject[],
}

export interface ParentSubaccountTransferResponseObject {
  id: string,
  sender: {
    address: string,
    /**
     * @isInt
     */
    parentSubaccountNumber?: number,
  },
  recipient: {
    address: string,
    /**
     * @isInt
     */
    parentSubaccountNumber?: number,
  },
  size: string,
  createdAt: string,
  createdAtHeight: string,
  symbol: string,
  type: TransferType,
  transactionHash: string,
}

export interface TransferBetweenResponse extends PaginationResponse {
  // Indexer will return data in descending order with the first transfer
  // being the most recent transfer. Will always return up to 100 transfers.
  // Transfers are categorized from the perspective of the source subaccount
  transfersSubset: TransferResponseObject[],

  // Given that source subaccount is the trader and the recipient subaccount
  // is the vault, total net transfer should always be positive
  totalNetTransfers: string,
}

/* ------- PNL TICKS TYPES ------- */

export interface HistoricalPnlResponse extends PaginationResponse {
  historicalPnl: PnlTicksResponseObject[],
}

export interface PnlTicksResponseObject {
  equity: string,
  totalPnl: string,
  netTransfers: string,
  createdAt: string,
  blockHeight: string,
  blockTime: IsoString,
}

export interface AggregatedPnlTick {
  pnlTick: PnlTicksFromDatabase,
  numTicks: number,
}

/* ------- TRADE TYPES ------- */

export interface TradeResponse extends PaginationResponse {
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

export type AssetById = { [assetId: string]: AssetFromDatabase };

export interface MarketAndType {
  marketType: MarketType,
  market: string,
}

export type MarketAndTypeByClobPairId = { [clobPairId: string]: MarketAndType };

export enum MarketType {
  PERPETUAL = 'PERPETUAL',
  SPOT = 'SPOT',
}

export interface PerpetualMarketResponse {
  markets: {
    [ticker: string]: PerpetualMarketResponseObject,
  },
}

export interface PerpetualMarketResponseObject {
  clobPairId: string,
  ticker: string,
  status: PerpetualMarketStatus,
  oraclePrice: string,
  priceChange24H: string,
  volume24H: string,
  /**
   * @isInt
   */
  trades24H: number,
  nextFundingRate: string,
  initialMarginFraction: string,
  maintenanceMarginFraction: string,
  openInterest: string,
  /**
   * @isInt
   */
  atomicResolution: number,
  /**
   * @isInt
   */
  quantumConversionExponent: number,
  tickSize: string,
  stepSize: string,
  /**
   * @isInt
   */
  stepBaseQuantums: number,
  /**
   * @isInt
   */
  subticksPerTick: number,
  marketType: PerpetualMarketType,
  openInterestLowerCap?: string,
  openInterestUpperCap?: string,
  baseOpenInterest: string,
  defaultFundingRate1H?: string,
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
  ticker: string,
  updatedAt?: IsoString,
  updatedAtHeight?: string,
  /**
   * @isInt
   */
  subaccountNumber: number,
}

export type RedisOrderMap = { [orderId: string]: RedisOrder };

export type PostgresOrderMap = { [orderId: string]: OrderFromDatabase };

/* ------- CANDLE TYPES ------- */

export interface CandleResponse {
  candles: CandleResponseObject[],
}

export interface CandleResponseObject extends Omit<CandleFromDatabase, CandleColumns.id> { }

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

export interface ParentSubaccountRequest extends AddressRequest {
  parentSubaccountNumber: number,
}

export interface PaginationRequest {
  page?: number,
}

export interface LimitRequest {
  limit: number,
}

export interface TickerRequest {
  ticker?: string,
}

interface CreatedBeforeRequest {
  createdBeforeOrAtHeight?: number,
  createdBeforeOrAt?: IsoString,
}

export interface LimitAndCreatedBeforeRequest extends LimitRequest, CreatedBeforeRequest { }

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

export interface ParentSubaccountPerpetualPositionRequest extends ParentSubaccountRequest,
  LimitAndCreatedBeforeRequest {
  status: PerpetualPositionStatus[],
}

export interface AssetPositionRequest extends SubaccountRequest { }

export interface ParentSubaccountAssetPositionRequest extends ParentSubaccountRequest {
}

export interface TransferRequest
  extends SubaccountRequest, LimitAndCreatedBeforeRequest, PaginationRequest { }

export interface ParentSubaccountTransferRequest
  extends ParentSubaccountRequest, LimitAndCreatedBeforeRequest, PaginationRequest {
}

export interface TransferBetweenRequest extends CreatedBeforeRequest {
  sourceAddress: string,
  sourceSubaccountNumber: number,
  recipientAddress: string,
  recipientSubaccountNumber: number,
}

export interface FillRequest
  extends SubaccountRequest, LimitAndCreatedBeforeRequest, PaginationRequest {
  market: string,
  marketType: MarketType,
  includeTypes?: FillType[],
  excludeTypes?: FillType[],
}

export interface ParentSubaccountFillRequest
  extends ParentSubaccountRequest, LimitAndCreatedBeforeRequest, PaginationRequest {
  market: string,
  marketType: MarketType,
  includeTypes?: FillType[],
  excludeTypes?: FillType[],
}

export interface TradeRequest extends LimitAndCreatedBeforeRequest, PaginationRequest {
  ticker: string,
}

export interface PerpetualMarketRequest extends LimitRequest, TickerRequest {
  market?: string,
}

export interface PnlTicksRequest
  extends SubaccountRequest, LimitAndCreatedBeforeAndAfterRequest, PaginationRequest { }

export interface ParentSubaccountPnlTicksRequest
  extends ParentSubaccountRequest, LimitAndCreatedBeforeAndAfterRequest {
}

export interface PnlRequest
  extends SubaccountRequest, LimitAndCreatedBeforeAndAfterRequest, PaginationRequest {
  daily?: boolean,
}

export interface ParentSubaccountPnlRequest
  extends ParentSubaccountRequest, LimitAndCreatedBeforeAndAfterRequest {
}

export interface OrderbookRequest {
  ticker: string,
}

export interface GetOrderRequest {
  orderId: string,
}

export interface ListOrderRequest extends SubaccountRequest, LimitRequest, TickerRequest {
  side?: OrderSide,
  type?: OrderType,
  includeTypes?: OrderType[],
  excludeTypes?: OrderType[],
  status?: OrderStatus[],
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockAfter?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
  goodTilBlockTimeAfter?: IsoString,
  returnLatestOrders?: boolean,
}

export interface ParentSubaccountListOrderRequest
  extends ParentSubaccountRequest, LimitRequest, TickerRequest {
  side?: OrderSide,
  type?: OrderType,
  includeTypes?: OrderType[],
  excludeTypes?: OrderType[],
  status?: OrderStatus[],
  goodTilBlockBeforeOrAt?: number,
  goodTilBlockAfter?: number,
  goodTilBlockTimeBeforeOrAt?: IsoString,
  goodTilBlockTimeAfter?: IsoString,
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

export interface RegisterTokenRequest {
  address: string,
  token: string,
  language: string,
  message: string,
  timestamp: number,
  signedMessage: string,
  pubKey: string,
  walletIsKeplr: boolean,
}

/* ------- COLLATERALIZATION TYPES ------- */

export interface Risk {
  initial: Big,
  maintenance: Big,
}

/* ------- COMPLIANCE TYPES ------- */

export interface ComplianceResponse {
  restricted: boolean,
  reason?: string,
}

export interface ComplianceRequest extends AddressRequest { }

export interface SetComplianceStatusRequest extends AddressRequest {
  status: ComplianceStatus,
  reason?: ComplianceReason,
}

export enum BlockedCode {
  GEOBLOCKED = 'GEOBLOCKED',
  COMPLIANCE_BLOCKED = 'COMPLIANCE_BLOCKED',
}

export interface ComplianceV2Response {
  status: ComplianceStatus,
  reason?: ComplianceReason,
  updatedAt?: string,
}

/* ------- HISTORICAL TRADING REWARD TYPES ------- */

export interface HistoricalTradingRewardAggregationRequest extends AddressRequest, LimitRequest {
  period: TradingRewardAggregationPeriod,
  startingBeforeOrAt: IsoString,
  startingBeforeOrAtHeight: string,
}

export interface HistoricalTradingRewardAggregationsResponse {
  // Indexer will not fill in empty periods, if there is no data after this period,
  // Indexer will return an empty list. Will return in descending order, the most
  // recent at the start
  rewards: HistoricalTradingRewardAggregation[],
}

export interface HistoricalTradingRewardAggregation {
  tradingReward: string, // i.e. '100.1' for 100.1 token earned through trading rewards
  startedAt: IsoString, // Start of the aggregation period, inclusive
  startedAtHeight: string, // first block included in the aggregation, inclusive
  endedAt?: IsoString, // End of the aggregation period, exclusive
  endedAtHeight?: string, // last block included in the aggregation, inclusive
  period: TradingRewardAggregationPeriod,
}

/* ------- HISTORICAL BLOCK TRADING REWARD TYPES ------- */
export interface HistoricalBlockTradingRewardRequest extends AddressRequest, LimitRequest {
  startingBeforeOrAt: IsoString,
  startingBeforeOrAtHeight: string,
}

export interface HistoricalBlockTradingRewardsResponse {
  // Indexer will not fill in empty periods, if there is no data after this period,
  // Indexer will return an empty list. Will return in descending order, the most
  // recent at the start
  rewards: HistoricalBlockTradingReward[],
}

export interface HistoricalBlockTradingReward {
  tradingReward: string, // i.e. '100.1' for 100.1 token earned through trading rewards
  createdAt: IsoString,
  createdAtHeight: string,
}

/* ------- Social Trading Types ------- */

export interface TraderSearchResponse {
  result?: TraderSearchResponseObject,
}

export interface TraderSearchRequest {
  searchParam: string,
}

export interface TraderSearchResponseObject {
  address: string,
  subaccountNumber: number,
  subaccountId: string,
  username: string,
}

/* ------- Vault Types ------- */

export interface VaultHistoricalPnl {
  ticker: string,
  historicalPnl: PnlTicksResponseObject[],
}

export interface MegavaultHistoricalPnlResponse {
  megavaultPnl: PnlTicksResponseObject[],
}

export interface VaultsHistoricalPnlResponse {
  vaultsPnl: VaultHistoricalPnl[],
}

export interface VaultPosition {
  ticker: string,
  assetPosition: AssetPositionResponseObject,
  perpetualPosition?: PerpetualPositionResponseObject,
  equity: string,
}

export interface MegavaultPositionResponse {
  positions: VaultPosition[],
}

export interface MegavaultHistoricalPnlRequest {
  resolution: PnlTickInterval,
}

export interface VaultsHistoricalPnlRequest extends MegavaultHistoricalPnlRequest { }

export interface VaultMapping {
  [subaccountId: string]: VaultFromDatabase,
}

/* ------- Affiliates Types ------- */
export interface AffiliateMetadataRequest {
  address: string,
}

export interface AffiliateAddressRequest {
  referralCode: string,
}

export interface AffiliateSnapshotRequest {
  addressFilter?: string[],
  limit?: number,
  offset?: number,
  sortByAffiliateEarning?: boolean,
}

export interface AffiliateTotalVolumeRequest {
  address: string,
}

export interface AffiliateMetadataResponse {
  referralCode: string,
  isVolumeEligible: boolean,
  isAffiliate: boolean,
}

export interface AffiliateAddressResponse {
  address: string,
}

export interface CreateReferralCodeRequest {
  address: string,
  newCode: string,
  signedMessage: string, // base64 encoded
  pubKey: string, // base64 encoded
  timestamp: number, // UNIX timestamp in seconds
}
export interface CreateReferralCodeResponse {
  referralCode: string,
}

export interface AffiliateSnapshotResponse {
  affiliateList: AffiliateSnapshotResponseObject[],
  total: number,
  currentOffset: number,
}

export interface AffiliateSnapshotResponseObject {
  affiliateAddress: string,
  affiliateReferralCode: string,
  affiliateEarnings: number,
  affiliateReferredTrades: number,
  affiliateTotalReferredFees: number,
  affiliateReferredUsers: number,
  affiliateReferredNetProtocolEarnings: number,
  affiliateReferredTotalVolume: number,
  affiliateReferredMakerFees: number,
  affiliateReferredTakerFees: number,
  affiliateReferredMakerRebates: number,
}

export interface AffiliateTotalVolumeResponse {
  totalVolume: number | null,
}

/* ------- FUNDING PAYMENTS TYPES ------- */
export interface FundingPaymentResponseObject {
  createdAt: IsoString,
  createdAtHeight: string,
  perpetualId: string,
  ticker: string,
  oraclePrice: string,
  size: string,
  side: string,
  rate: string,
  payment: string,
  subaccountNumber: string,
  fundingIndex: string,
}

export interface FundingPaymentResponse extends PaginationResponse {
  fundingPayments: FundingPaymentResponseObject[],
}

/* ------- TURNKEY TYPES ------- */
export interface TurnkeyAuthResponse {
  dydxAddress?: string,
  organizationId?: string,
  apiKeyId?: string,
  userId?: string,
  session?: string,
  salt: string,
  alreadyExists?: boolean,
}

export interface TurnkeyCreateSuborgResponse {
  subOrgId: string,
  apiKeyId?: string,
  userId?: string,
  dydxAddress?: string,
  salt: string,
}

export interface CreateSuborgParams {
  email?: string,
  providerName?: string,
  oidcToken?: string,
  authenticatorName?: string,
  challenge?: string,
  attestation?: TurnkeyApiTypes['v1Attestation'],
}
export interface GetSuborgParams {
  email?: string,
  oidcToken?: string,
  credentialId?: string,
}

export enum SigninMethod {
  EMAIL = 'email',
  SOCIAL = 'social',
  PASSKEY = 'passkey',
}

/* ------- APPLE LOGIN TYPES ------- */

export interface AppleLoginRedirectRequest {
  state: string, // public key from client
  code: string,  // auth code from Apple
}

export interface AppleTokenResponse {
  access_token: string,
  token_type: string,
  expires_in: number,
  refresh_token?: string,
  id_token: string,
}

export interface AppleJWTClaims {
  iss: string,   // Team ID
  iat: number,   // Issued at
  exp: number,   // Expires at
  aud: string,   // Audience (https://appleid.apple.com)
  sub: string,   // Service ID (client_id)
}

export interface AppleLoginResponse {
  success: boolean,
  encodedPayload?: string,
  error?: string,
}

/* ------- PNL TYPES ------- */

export interface PnlResponse extends PaginationResponse {
  pnl: PnlResponseObject[],
}

export interface PnlResponseObject {
  equity: string,
  netTransfers: string,
  totalPnl: string,
  createdAt: string,
  createdAtHeight: string,
}

export interface AggregatedPnl {
  pnl: PnlFromDatabase,
  numPnls: number,
}
