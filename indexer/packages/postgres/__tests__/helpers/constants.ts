import {
  ORDER_FLAG_CONDITIONAL,
  ORDER_FLAG_LONG_TERM,
  ORDER_FLAG_SHORT_TERM,
} from '@dydxprotocol-indexer/v4-proto-parser';
import { DateTime } from 'luxon';

import * as AssetPositionTable from '../../src/stores/asset-position-table';
import * as CandleTable from '../../src/stores/candle-table';
import * as FundingIndexUpdatesTable from '../../src/stores/funding-index-updates-table';
import * as OraclePriceTable from '../../src/stores/oracle-price-table';
import * as OrderTable from '../../src/stores/order-table';
import * as PerpetualPositionTable from '../../src/stores/perpetual-position-table';
import * as SubaccountTable from '../../src/stores/subaccount-table';
import * as TendermintEventTable from '../../src/stores/tendermint-event-table';
import * as TransactionTable from '../../src/stores/transaction-table';
import * as TransferTable from '../../src/stores/transfer-table';
import {
  AssetCreateObject,
  AssetPositionCreateObject,
  BlockCreateObject,
  CandleCreateObject,
  CandleResolution,
  ComplianceDataCreateObject,
  ComplianceProvider,
  FillCreateObject,
  FillType,
  FundingIndexUpdatesCreateObject,
  Liquidity,
  LiquidityTiersCreateObject,
  MarketCreateObject,
  OraclePriceCreateObject,
  OrderCreateObject,
  OrderSide,
  OrderStatus,
  OrderType,
  PerpetualMarketCreateObject,
  PerpetualMarketStatus,
  PerpetualMarketType,
  PerpetualPositionCreateObject,
  PerpetualPositionStatus,
  PnlTicksCreateObject,
  PositionSide,
  SubaccountCreateObject,
  TendermintEventCreateObject,
  TimeInForce,
  TransactionCreateObject,
  TransferCreateObject,
} from '../../src/types';

export const createdDateTime: DateTime = DateTime.utc();
export const createdHeight: string = '2';
export const invalidTicker: string = 'INVALID-INVALID';
export const dydxChain: string = 'dydx';
export const defaultAddress: string = 'dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf';
export const blockedAddress: string = 'dydx1f9k5qldwmqrnwy8hcgp4fw6heuvszt35egvtx2';

// ============== Subaccounts ==============

export const defaultSubaccount: SubaccountCreateObject = {
  address: defaultAddress,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const defaultSubaccount2: SubaccountCreateObject = {
  address: defaultAddress,
  subaccountNumber: 1,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const defaultSubaccount3: SubaccountCreateObject = {
  address: defaultAddress,
  subaccountNumber: 2,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const defaultWalletAddress: string = 'defaultWalletAddress';

export const defaultSubaccountId: string = SubaccountTable.uuid(
  defaultAddress,
  defaultSubaccount.subaccountNumber,
);
export const defaultSubaccountId2: string = SubaccountTable.uuid(
  defaultAddress,
  defaultSubaccount2.subaccountNumber,
);
export const defaultSubaccountId3: string = SubaccountTable.uuid(
  defaultAddress,
  defaultSubaccount3.subaccountNumber,
);

export const isolatedSubaccount: SubaccountCreateObject = {
  address: defaultAddress,
  subaccountNumber: 128,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const isolatedSubaccount2: SubaccountCreateObject = {
  address: defaultAddress,
  subaccountNumber: 256,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const isolatedSubaccountId: string = SubaccountTable.uuid(
  defaultAddress,
  isolatedSubaccount.subaccountNumber,
);
export const isolatedSubaccountId2: string = SubaccountTable.uuid(
  defaultAddress,
  isolatedSubaccount2.subaccountNumber,
);

// ============== Assets ==============

export const defaultAsset: AssetCreateObject = {
  id: '0',
  symbol: 'USDC',
  atomicResolution: -6,
  hasMarket: false,
};
export const defaultAsset2: AssetCreateObject = {
  id: '1',
  symbol: 'DYDX',
  atomicResolution: 0,
  hasMarket: true,
  marketId: 1,
};
export const defaultAsset3: AssetCreateObject = {
  id: '2',
  symbol: 'WBTC',
  atomicResolution: -8,
  hasMarket: false,
};

// ============== AssetPositions ==============

export const defaultAssetPosition: AssetPositionCreateObject = {
  subaccountId: defaultSubaccountId,
  assetId: '0',
  size: '10000',
  isLong: true,
};
export const defaultAssetPositionId: string = AssetPositionTable.uuid(
  defaultAssetPosition.subaccountId,
  defaultAssetPosition.assetId,
);
export const defaultAssetPosition2: AssetPositionCreateObject = {
  subaccountId: defaultSubaccountId2,
  assetId: '1',
  size: '10000',
  isLong: false,
};
export const defaultAssetPositionId2: string = AssetPositionTable.uuid(
  defaultAssetPosition2.subaccountId,
  defaultAssetPosition2.assetId,
);
export const isolatedSubaccountAssetPosition: AssetPositionCreateObject = {
  subaccountId: isolatedSubaccountId,
  assetId: '0',
  size: '5000',
  isLong: true,
};
export const isolatedSubaccountAssetPositionId: string = AssetPositionTable.uuid(
  isolatedSubaccountAssetPosition.subaccountId,
  isolatedSubaccountAssetPosition.assetId,
);

// ============== PerpetualMarkets ==============

export const defaultPerpetualMarket: PerpetualMarketCreateObject = {
  id: '0',
  clobPairId: '1',
  ticker: 'BTC-USD',
  marketId: 0,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '23',
  volume24H: '1000000',
  trades24H: 250,
  nextFundingRate: '10.2',
  openInterest: '400000',
  quantumConversionExponent: -8,
  atomicResolution: -10,
  dangerIndexPpm: 1000000,
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: '0',
  subticksPerTick: 100,
  stepBaseQuantums: 10,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
};
export const defaultPerpetualMarket2: PerpetualMarketCreateObject = {
  id: '1',
  clobPairId: '2',
  ticker: 'ETH-USD',
  marketId: 1,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '23',
  volume24H: '100000',
  trades24H: 200,
  nextFundingRate: '1.2',
  openInterest: '40000',
  quantumConversionExponent: -6,
  atomicResolution: -18,
  dangerIndexPpm: 1000000,
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: '0',
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
};
export const defaultPerpetualMarket3: PerpetualMarketCreateObject = {
  id: '2',
  clobPairId: '3',
  ticker: 'SHIB-USD',
  marketId: 2,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '0.000000001',
  volume24H: '10000000',
  trades24H: 200,
  nextFundingRate: '1.2',
  openInterest: '40000',
  quantumConversionExponent: -16,
  atomicResolution: -2,
  dangerIndexPpm: 1000000,
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: '0',
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
};

export const isolatedPerpetualMarket: PerpetualMarketCreateObject = {
  id: '3',
  clobPairId: '4',
  ticker: 'ISO-USD',
  marketId: 3,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '0.000000001',
  volume24H: '10000000',
  trades24H: 200,
  nextFundingRate: '1.2',
  openInterest: '40000',
  quantumConversionExponent: -16,
  atomicResolution: -2,
  dangerIndexPpm: 1000000,
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: '1000000',
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.ISOLATED,
  baseOpenInterest: '100000',
};

export const isolatedPerpetualMarket2: PerpetualMarketCreateObject = {
  id: '4',
  clobPairId: '5',
  ticker: 'ISO2-USD',
  marketId: 4,
  status: PerpetualMarketStatus.ACTIVE,
  priceChange24H: '0.000000001',
  volume24H: '10000000',
  trades24H: 200,
  nextFundingRate: '1.2',
  openInterest: '40000',
  quantumConversionExponent: -16,
  atomicResolution: -2,
  dangerIndexPpm: 1000000,
  isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock: '1000000',
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.ISOLATED,
  baseOpenInterest: '100000',
};

// ============== Orders ==============

export const defaultOrder: OrderCreateObject = {
  subaccountId: defaultSubaccountId,
  clientId: '1',
  clobPairId: '1',
  side: OrderSide.BUY,
  size: '25',
  totalFilled: '0',
  price: '20000',
  type: OrderType.LIMIT,
  status: OrderStatus.OPEN,
  timeInForce: TimeInForce.FOK,
  reduceOnly: false,
  goodTilBlock: '100',
  orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
  clientMetadata: '0',
  updatedAt: '2023-01-22T00:00:00.000Z',
  updatedAtHeight: '1',
};

export const isolatedMarketOrder: OrderCreateObject = {
  subaccountId: isolatedSubaccountId,
  clientId: '1',
  clobPairId: '4',
  side: OrderSide.BUY,
  size: '25',
  totalFilled: '0',
  price: '20000',
  type: OrderType.LIMIT,
  status: OrderStatus.OPEN,
  timeInForce: TimeInForce.FOK,
  reduceOnly: false,
  goodTilBlock: '100',
  orderFlags: ORDER_FLAG_SHORT_TERM.toString(),
  clientMetadata: '0',
  updatedAt: '2023-01-22T00:00:00.000Z',
  updatedAtHeight: '1',
};

export const defaultOrderGoodTilBlockTime: OrderCreateObject = {
  ...defaultOrder,
  clientId: '2',
  goodTilBlock: undefined,
  goodTilBlockTime: '2023-01-22T00:00:00.000Z',
  createdAtHeight: '1',
  orderFlags: ORDER_FLAG_LONG_TERM.toString(),
};

export const defaultConditionalOrder: OrderCreateObject = {
  ...defaultOrderGoodTilBlockTime,
  type: OrderType.STOP_LIMIT,
  clientId: '3',
  orderFlags: ORDER_FLAG_CONDITIONAL.toString(),
  triggerPrice: '19000',
};

export const defaultOrderId: string = OrderTable.uuid(
  defaultOrder.subaccountId,
  defaultOrder.clientId,
  defaultOrder.clobPairId,
  defaultOrder.orderFlags,
);

export const defaultOrderGoodTilBlockTimeId: string = OrderTable.uuid(
  defaultOrderGoodTilBlockTime.subaccountId,
  defaultOrderGoodTilBlockTime.clientId,
  defaultOrderGoodTilBlockTime.clobPairId,
  defaultOrderGoodTilBlockTime.orderFlags,
);

export const defaultConditionalOrderId: string = OrderTable.uuid(
  defaultConditionalOrder.subaccountId,
  defaultConditionalOrder.clientId,
  defaultConditionalOrder.clobPairId,
  defaultConditionalOrder.orderFlags,
);

export const isolatedMarketOrderId: string = OrderTable.uuid(
  isolatedMarketOrder.subaccountId,
  isolatedMarketOrder.clientId,
  isolatedMarketOrder.clobPairId,
  isolatedMarketOrder.orderFlags,
);

// ============== Blocks ==============

export const defaultBlock: BlockCreateObject = {
  blockHeight: '1',
  time: DateTime.utc(2022, 6, 1).toISO(),
};
export const defaultBlock2: BlockCreateObject = {
  blockHeight: '2',
  time: DateTime.utc(2022, 6, 2).toISO(),
};

// ============== TendermintEvents ==============

export const defaultTendermintEvent: TendermintEventCreateObject = {
  blockHeight: '1',
  transactionIndex: -1,
  eventIndex: 0,
};
export const defaultTendermintEvent2: TendermintEventCreateObject = {
  blockHeight: '2',
  transactionIndex: -1,
  eventIndex: 0,
};
export const defaultTendermintEvent3: TendermintEventCreateObject = {
  blockHeight: '2',
  transactionIndex: 0,
  eventIndex: 0,
};
export const defaultTendermintEvent4: TendermintEventCreateObject = {
  blockHeight: '2',
  transactionIndex: 1,
  eventIndex: 1,
};
export const defaultTendermintEventId: Buffer = TendermintEventTable.createEventId(
  defaultTendermintEvent.blockHeight,
  defaultTendermintEvent.transactionIndex,
  defaultTendermintEvent.eventIndex,
);
export const defaultTendermintEventId2: Buffer = TendermintEventTable.createEventId(
  defaultTendermintEvent2.blockHeight,
  defaultTendermintEvent2.transactionIndex,
  defaultTendermintEvent2.eventIndex,
);
export const defaultTendermintEventId3: Buffer = TendermintEventTable.createEventId(
  defaultTendermintEvent3.blockHeight,
  defaultTendermintEvent3.transactionIndex,
  defaultTendermintEvent3.eventIndex,
);
export const defaultTendermintEventId4: Buffer = TendermintEventTable.createEventId(
  defaultTendermintEvent4.blockHeight,
  defaultTendermintEvent4.transactionIndex,
  defaultTendermintEvent4.eventIndex,
);

// ============== Transactions ==============

export const defaultTransaction: TransactionCreateObject = {
  blockHeight: '1',
  transactionIndex: 0,
  transactionHash: '3ac776f8-1900-43de-ac38-7fb516f7d6d0',
};
export const defaultTransactionId: string = TransactionTable.uuid(
  defaultTransaction.blockHeight,
  defaultTransaction.transactionIndex,
);

// ============== PerpetualPositions ==============

export const defaultPerpetualPosition: PerpetualPositionCreateObject = {
  subaccountId: defaultSubaccountId,
  perpetualId: defaultPerpetualMarket.id,
  side: PositionSide.LONG,
  status: PerpetualPositionStatus.OPEN,
  size: '10',
  maxSize: '25',
  entryPrice: '20000',
  sumOpen: '10',
  sumClose: '0',
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  openEventId: defaultTendermintEventId,
  lastEventId: defaultTendermintEventId2,
  settledFunding: '200000',
};

export const defaultPerpetualPositionId: string = PerpetualPositionTable.uuid(
  defaultPerpetualPosition.subaccountId,
  defaultPerpetualPosition.openEventId,
);

export const isolatedPerpetualPosition: PerpetualPositionCreateObject = {
  subaccountId: isolatedSubaccountId,
  perpetualId: isolatedPerpetualMarket.id,
  side: PositionSide.LONG,
  status: PerpetualPositionStatus.OPEN,
  size: '10',
  maxSize: '25',
  entryPrice: '20000',
  sumOpen: '10',
  sumClose: '0',
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  openEventId: defaultTendermintEventId,
  lastEventId: defaultTendermintEventId2,
  settledFunding: '200000',
};

export const isolatedPerpetualPositionId: string = PerpetualPositionTable.uuid(
  isolatedPerpetualPosition.subaccountId,
  isolatedPerpetualPosition.openEventId,
);

// ============== Fills ==============

export const defaultFill: FillCreateObject = {
  subaccountId: defaultSubaccountId,
  side: OrderSide.BUY,
  liquidity: Liquidity.TAKER,
  type: FillType.LIMIT,
  clobPairId: '1',
  orderId: defaultOrderId,
  size: '10',
  price: '20000',
  quoteAmount: '200000',
  eventId: defaultTendermintEventId,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  clientMetadata: '0',
  fee: '1.1',
};

export const isolatedMarketFill: FillCreateObject = {
  subaccountId: isolatedSubaccountId,
  side: OrderSide.BUY,
  liquidity: Liquidity.TAKER,
  type: FillType.LIMIT,
  clobPairId: '4',
  orderId: isolatedMarketOrderId,
  size: '10',
  price: '20000',
  quoteAmount: '200000',
  eventId: defaultTendermintEventId2,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  clientMetadata: '0',
  fee: '1.1',
};

export const isolatedMarketFill2: FillCreateObject = {
  subaccountId: isolatedSubaccountId2,
  side: OrderSide.BUY,
  liquidity: Liquidity.TAKER,
  type: FillType.LIMIT,
  clobPairId: '4',
  orderId: isolatedMarketOrderId,
  size: '10',
  price: '20000',
  quoteAmount: '200000',
  eventId: defaultTendermintEventId3,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  clientMetadata: '0',
  fee: '1.1',
};

// ============== Transfers ==============

export const defaultTransfer: TransferCreateObject = {
  senderSubaccountId: defaultSubaccountId,
  recipientSubaccountId: defaultSubaccountId2,
  assetId: defaultAsset.id,
  size: '10',
  eventId: defaultTendermintEventId,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
};

export const defaultTransfer2: TransferCreateObject = {
  ...defaultTransfer,
  senderSubaccountId: defaultSubaccountId3,
  size: '5',
};

export const defaultTransfer3: TransferCreateObject = {
  ...defaultTransfer2,
  assetId: defaultAsset2.id,
};

export const defaultTransferId: string = TransferTable.uuid(
  defaultTransfer.eventId,
  defaultTransfer.assetId,
  defaultTransfer.senderSubaccountId,
  defaultTransfer.recipientSubaccountId,
  defaultTransfer.senderWalletAddress,
  defaultTransfer.recipientWalletAddress,
);

export const defaultWithdrawal: TransferCreateObject = {
  senderSubaccountId: defaultSubaccountId,
  recipientWalletAddress: defaultWalletAddress,
  assetId: defaultAsset.id,
  size: '10',
  eventId: defaultTendermintEventId,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
};

export const defaultWithdrawalId: string = TransferTable.uuid(
  defaultWithdrawal.eventId,
  defaultWithdrawal.assetId,
  defaultWithdrawal.senderSubaccountId,
  defaultWithdrawal.recipientSubaccountId,
  defaultWithdrawal.senderWalletAddress,
  defaultWithdrawal.recipientWalletAddress,
);

export const defaultDeposit: TransferCreateObject = {
  senderWalletAddress: defaultWalletAddress,
  recipientSubaccountId: defaultSubaccountId,
  assetId: defaultAsset.id,
  size: '10',
  eventId: defaultTendermintEventId,
  transactionHash: '', // TODO: Add a real transaction Hash
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
};

export const defaultDepositId: string = TransferTable.uuid(
  defaultDeposit.eventId,
  defaultDeposit.assetId,
  defaultDeposit.senderSubaccountId,
  defaultDeposit.recipientSubaccountId,
  defaultDeposit.senderWalletAddress,
  defaultDeposit.recipientWalletAddress,
);

// ============== Markets ==============

export const defaultMarket: MarketCreateObject = {
  id: 0,
  pair: 'BTC-USD',
  exponent: -5,
  minPriceChangePpm: 50,
  spotPrice: '15000',
  pnlPrice: '15000',
};

export const defaultMarket2: MarketCreateObject = {
  id: 1,
  pair: 'ETH-USD',
  exponent: -6,
  minPriceChangePpm: 50,
  spotPrice: '1000',
  pnlPrice: '1000',
};

export const defaultMarket3: MarketCreateObject = {
  id: 2,
  pair: 'SHIB-USD',
  exponent: -12,
  minPriceChangePpm: 50,
  spotPrice: '0.000000065',
  pnlPrice: '0.000000065',
};

export const isolatedMarket: MarketCreateObject = {
  id: 3,
  pair: 'ISO-USD',
  exponent: -12,
  minPriceChangePpm: 50,
  spotPrice: '0.000000075',
  pnlPrice: '0.000000075',
};

export const isolatedMarket2: MarketCreateObject = {
  id: 4,
  pair: 'ISO2-USD',
  exponent: -12,
  minPriceChangePpm: 50,
  spotPrice: '0.000000085',
  pnlPrice: '0.000000085',
};

// ============== LiquidityTiers ==============

export const defaultLiquidityTier: LiquidityTiersCreateObject = {
  id: 0,
  name: 'Large-Cap',
  initialMarginPpm: '50000',  // 5%
  maintenanceFractionPpm: '600000',  // 60%
};

export const defaultLiquidityTier2: LiquidityTiersCreateObject = {
  id: 1,
  name: 'Mid-Cap',
  initialMarginPpm: '100000',  // 10%
  maintenanceFractionPpm: '500000',  // 50%
  openInterestLowerCap: '0',
  openInterestUpperCap: '5000000',
};

// ============== OraclePrices ==============

export const defaultOraclePrice: OraclePriceCreateObject = {
  marketId: defaultMarket.id,
  spotPrice: '10000',
  pnlPrice: '10000',
  effectiveAt: createdDateTime.toISO(),
  effectiveAtHeight: createdHeight,
};

export const defaultOraclePriceId: string = OraclePriceTable.uuid(
  defaultOraclePrice.marketId,
  defaultOraclePrice.effectiveAtHeight,
);

export const defaultOraclePrice2: OraclePriceCreateObject = {
  marketId: defaultMarket2.id,
  spotPrice: '500',
  pnlPrice: '500',
  effectiveAt: createdDateTime.toISO(),
  effectiveAtHeight: createdHeight,
};

export const defaultOraclePriceId2: string = OraclePriceTable.uuid(
  defaultOraclePrice2.marketId,
  defaultOraclePrice2.effectiveAtHeight,
);

// ============== Candle ==============

export const defaultCandle: CandleCreateObject = {
  startedAt: createdDateTime.toISO(),
  ticker: defaultPerpetualMarket.ticker,
  resolution: CandleResolution.ONE_MINUTE,
  low: '10000',
  high: '12000',
  open: '11000',
  close: '11500',
  baseTokenVolume: '400000',
  usdVolume: '2200000',
  trades: 300,
  startingOpenInterest: '200000',
};

export const defaultCandleId: string = CandleTable.uuid(
  defaultCandle.startedAt,
  defaultCandle.ticker,
  defaultCandle.resolution,
);

// ============== Pnl Ticks ==============

export const defaultPnlTick: PnlTicksCreateObject = {
  subaccountId: defaultSubaccountId,
  equity: '100000',
  totalPnl: '10000',
  netTransfers: '1000',
  createdAt: createdDateTime.toISO(),
  blockHeight: createdHeight,
  blockTime: defaultBlock2.time,
};

// ========= Funding Index updates ==========

export const defaultFundingIndexUpdate: FundingIndexUpdatesCreateObject = {
  perpetualId: defaultPerpetualMarket.id,
  eventId: defaultTendermintEventId,
  rate: '0.0004',
  oraclePrice: '10000',
  fundingIndex: '10050',
  effectiveAt: createdDateTime.toISO(),
  effectiveAtHeight: createdHeight,
};

export const defaultFundingIndexUpdateId: string = FundingIndexUpdatesTable.uuid(
  defaultFundingIndexUpdate.effectiveAtHeight,
  defaultFundingIndexUpdate.eventId,
  defaultFundingIndexUpdate.perpetualId,
);

export const isolatedMarketFundingIndexUpdate: FundingIndexUpdatesCreateObject = {
  perpetualId: isolatedPerpetualMarket.id,
  eventId: defaultTendermintEventId,
  rate: '0.0004',
  oraclePrice: '10000',
  fundingIndex: '10200',
  effectiveAt: createdDateTime.toISO(),
  effectiveAtHeight: createdHeight,
};

export const isolatedMarketFundingIndexUpdateId: string = FundingIndexUpdatesTable.uuid(
  isolatedMarketFundingIndexUpdate.effectiveAtHeight,
  isolatedMarketFundingIndexUpdate.eventId,
  isolatedMarketFundingIndexUpdate.perpetualId,
);

// ========= Compliance Data ==========

export const blockedComplianceData: ComplianceDataCreateObject = {
  address: blockedAddress,
  provider: ComplianceProvider.ELLIPTIC,
  chain: dydxChain,
  blocked: true,
  riskScore: '100.00',
  updatedAt: createdDateTime.toISO(),
};

export const nonBlockedComplianceData: ComplianceDataCreateObject = {
  address: defaultAddress,
  provider: ComplianceProvider.ELLIPTIC,
  chain: dydxChain,
  blocked: false,
  riskScore: '10.00',
  updatedAt: createdDateTime.plus(1).toISO(),
};
