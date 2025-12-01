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
import * as TradingRewardAggregationTable from '../../src/stores/trading-reward-aggregation-table';
import * as TransactionTable from '../../src/stores/transaction-table';
import * as TransferTable from '../../src/stores/transfer-table';
import {
  AffiliateInfoCreateObject,
  AffiliateReferredUsersCreateObject,
  AssetCreateObject,
  AssetPositionCreateObject,
  BlockCreateObject,
  CandleCreateObject,
  CandleResolution,
  ComplianceDataCreateObject,
  ComplianceProvider,
  ComplianceReason,
  ComplianceStatus,
  ComplianceStatusCreateObject,
  ComplianceStatusUpsertObject,
  FillCreateObject,
  FillType,
  FundingIndexUpdatesCreateObject,
  LeaderboardPnlCreateObject,
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
  PnlCreateObject,
  PnlTicksCreateObject,
  PositionSide,
  SubaccountCreateObject,
  SubaccountUsernamesCreateObject,
  TendermintEventCreateObject,
  TimeInForce,
  TradingRewardAggregationCreateObject,
  TradingRewardAggregationPeriod,
  TradingRewardCreateObject,
  TransactionCreateObject,
  TransferCreateObject,
  WalletCreateObject,
  PersistentCacheCreateObject,
  VaultCreateObject,
  VaultStatus,
  BridgeInformationCreateObject,
} from '../../src/types';
import { denomToHumanReadableConversion } from './conversion-helpers';

export const createdDateTime: DateTime = DateTime.utc();
export const createdHeight: string = '2';
export const invalidTicker: string = 'INVALID-INVALID';
export const dydxChain: string = 'dydx';
export const defaultAddress: string = 'dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc565lnf';
export const defaultAddress2: string = 'dydx1n88uc38xhjgxzw9nwre4ep2c8ga4fjxc575lnf';
export const defaultAddress3: string = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';
export const blockedAddress: string = 'dydx1f9k5qldwmqrnwy8hcgp4fw6heuvszt35egvtx2';
// Vault address for vault id 0 was generated using
// script protocol/scripts/vault/get_vault.go
export const vaultAddress: string = 'dydx1c0m5x87llaunl5sgv3q5vd7j5uha26d2q2r2q0';

// ============== Subaccounts ==============
export const defaultWalletAddress: string = 'dydx199tqg4wdlnu4qjlxchpd7seg454937hjrknju4';

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

export const defaultSubaccount2Num0: SubaccountCreateObject = {
  address: defaultAddress2,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const defaultSubaccount3Num0: SubaccountCreateObject = {
  address: defaultAddress3,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

// defaultWalletAddress belongs to defaultWallet2 and is different from defaultAddress
export const defaultSubaccountDefaultWalletAddress: SubaccountCreateObject = {
  address: defaultWalletAddress,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const defaultSubaccountWithAlternateAddress: SubaccountCreateObject = {
  address: defaultAddress2,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

export const vaultSubaccount: SubaccountCreateObject = {
  address: vaultAddress,
  subaccountNumber: 0,
  updatedAt: createdDateTime.toISO(),
  updatedAtHeight: createdHeight,
};

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
export const defaultSubaccountIdDefaultWalletAddress: string = SubaccountTable.uuid(
  defaultWalletAddress,
  defaultSubaccountDefaultWalletAddress.subaccountNumber,
);
export const defaultSubaccountIdWithAlternateAddress: string = SubaccountTable.uuid(
  defaultAddress2,
  defaultSubaccountWithAlternateAddress.subaccountNumber,
);
export const isolatedSubaccountId: string = SubaccountTable.uuid(
  defaultAddress,
  isolatedSubaccount.subaccountNumber,
);
export const isolatedSubaccountId2: string = SubaccountTable.uuid(
  defaultAddress,
  isolatedSubaccount2.subaccountNumber,
);

export const vaultSubaccountId: string = SubaccountTable.uuid(
  vaultAddress,
  vaultSubaccount.subaccountNumber,
);

// ============== Wallets ==============
export const defaultWallet: WalletCreateObject = {
  address: defaultAddress,
  totalTradingRewards: denomToHumanReadableConversion(0),
  totalVolume: '0',
};

export const defaultWallet2: WalletCreateObject = {
  address: defaultWalletAddress,
  totalTradingRewards: denomToHumanReadableConversion(1),
  totalVolume: '0',
};

export const vaultWallet: WalletCreateObject = {
  address: vaultAddress,
  totalTradingRewards: denomToHumanReadableConversion(0),
  totalVolume: '0',
};

export const defaultWallet3: WalletCreateObject = {
  address: defaultAddress2,
  totalTradingRewards: denomToHumanReadableConversion(0),
  totalVolume: '0',
};

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
  subticksPerTick: 100,
  stepBaseQuantums: 10,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
  defaultFundingRate1H: '0',
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
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
  defaultFundingRate1H: '0',
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
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.CROSS,
  baseOpenInterest: '100000',
  defaultFundingRate1H: '0',
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
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.ISOLATED,
  baseOpenInterest: '100000',
  defaultFundingRate1H: '0.0001',
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
  subticksPerTick: 10,
  stepBaseQuantums: 1,
  liquidityTierId: 0,
  marketType: PerpetualMarketType.ISOLATED,
  baseOpenInterest: '100000',
  defaultFundingRate1H: '0.0001',
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
  orderRouterAddress: '',
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
  orderRouterAddress: '',
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

export const isolatedMarketOrderId: string = OrderTable.uuid(
  isolatedMarketOrder.subaccountId,
  isolatedMarketOrder.clientId,
  isolatedMarketOrder.clobPairId,
  isolatedMarketOrder.orderFlags,
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
  totalRealizedPnl: '100',
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
  entryPrice: '1.5',
  sumOpen: '10',
  sumClose: '0',
  createdAt: createdDateTime.toISO(),
  createdAtHeight: createdHeight,
  openEventId: defaultTendermintEventId,
  lastEventId: defaultTendermintEventId2,
  settledFunding: '200000',
  totalRealizedPnl: '100',
};

export const isolatedPerpetualPositionId: string = PerpetualPositionTable.uuid(
  isolatedPerpetualPosition.subaccountId,
  isolatedPerpetualPosition.openEventId,
);

// ============== Fills ==============

export const noBuilderAddress: string = '';

export const noOrderRouterAddress: string = '';

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
  affiliateRevShare: '1.10',
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
  affiliateRevShare: '0',
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
  affiliateRevShare: '0',
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

export const defaultTransferWithAlternateAddress: TransferCreateObject = {
  ...defaultTransfer,
  senderSubaccountId: defaultSubaccountIdWithAlternateAddress,
  recipientSubaccountId: defaultSubaccountId,
};

export const defaultTransferId: string = TransferTable.uuid(
  defaultTransfer.eventId,
  defaultTransfer.assetId,
  defaultTransfer.senderSubaccountId,
  defaultTransfer.recipientSubaccountId,
  defaultTransfer.senderWalletAddress,
  defaultTransfer.recipientWalletAddress,
);

export const defaultTransferWithAlternateAddressId: string = TransferTable.uuid(
  defaultTransferWithAlternateAddress.eventId,
  defaultTransferWithAlternateAddress.assetId,
  defaultTransferWithAlternateAddress.senderSubaccountId,
  defaultTransferWithAlternateAddress.recipientSubaccountId,
  defaultTransferWithAlternateAddress.senderWalletAddress,
  defaultTransferWithAlternateAddress.recipientWalletAddress,
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
  oraclePrice: '15000',
};

export const defaultMarket2: MarketCreateObject = {
  id: 1,
  pair: 'ETH-USD',
  exponent: -6,
  minPriceChangePpm: 50,
  oraclePrice: '1000',
};

export const defaultMarket3: MarketCreateObject = {
  id: 2,
  pair: 'SHIB-USD',
  exponent: -12,
  minPriceChangePpm: 50,
  oraclePrice: '0.000000065',
};

export const isolatedMarket: MarketCreateObject = {
  id: 3,
  pair: 'ISO-USD',
  exponent: -12,
  minPriceChangePpm: 50,
  oraclePrice: '1.00',
};

export const isolatedMarket2: MarketCreateObject = {
  id: 4,
  pair: 'ISO2-USD',
  exponent: -20,
  minPriceChangePpm: 50,
  oraclePrice: '0.000000085',
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
  price: '10000',
  effectiveAt: createdDateTime.toISO(),
  effectiveAtHeight: createdHeight,
};

export const defaultOraclePriceId: string = OraclePriceTable.uuid(
  defaultOraclePrice.marketId,
  defaultOraclePrice.effectiveAtHeight,
);

export const defaultOraclePrice2: OraclePriceCreateObject = {
  marketId: defaultMarket2.id,
  price: '500',
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
  orderbookMidPriceOpen: '11500',
  orderbookMidPriceClose: '12500',
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

// ========= Compliance Status ==========

export const compliantStatusData: ComplianceStatusCreateObject = {
  address: defaultAddress,
  status: ComplianceStatus.COMPLIANT,
  createdAt: createdDateTime.toISO(),
  updatedAt: createdDateTime.toISO(),
};

export const noncompliantStatusData: ComplianceStatusCreateObject = {
  address: blockedAddress,
  status: ComplianceStatus.BLOCKED,
  reason: ComplianceReason.SANCTIONED_GEO,
  createdAt: createdDateTime.plus(1).toISO(),
  updatedAt: createdDateTime.plus(1).toISO(),
};

export const noncompliantStatusUpsertData: ComplianceStatusUpsertObject = {
  address: blockedAddress,
  status: ComplianceStatus.BLOCKED,
  reason: ComplianceReason.SANCTIONED_GEO,
  updatedAt: createdDateTime.toISO(),
};

// ========= Trading Reward Data ==========

export const defaultTradingReward: TradingRewardCreateObject = {
  address: defaultAddress,
  blockHeight: createdHeight,
  blockTime: createdDateTime.toISO(),
  amount: denomToHumanReadableConversion(1),
};

// ========= Trading Reward Aggregation Data ==========

export const defaultTradingRewardAggregation: TradingRewardAggregationCreateObject = {
  address: defaultAddress,
  startedAtHeight: createdHeight,
  startedAt: createdDateTime.toISO(),
  period: TradingRewardAggregationPeriod.DAILY,
  amount: denomToHumanReadableConversion(1),
};
export const defaultTradingRewardAggregationId: string = TradingRewardAggregationTable.uuid(
  defaultTradingRewardAggregation.address,
  defaultTradingRewardAggregation.period,
  defaultTradingRewardAggregation.startedAtHeight,
);

// ============== Subaccount Usernames ==============
export const defaultSubaccountUsername: SubaccountUsernamesCreateObject = {
  username: 'LyingRaisin32',
  subaccountId: defaultSubaccountId,
};

export const defaultSubaccountUsername2: SubaccountUsernamesCreateObject = {
  username: 'LyingRaisin33',
  subaccountId: defaultSubaccountId2,
};

export const duplicatedSubaccountUsername: SubaccountUsernamesCreateObject = {
  username: 'LyingRaisin32',
  subaccountId: defaultSubaccountId3,
};

// defaultWalletAddress belongs to defaultWallet2 and is different from defaultAddress
export const subaccountUsernameWithDefaultWalletAddress: SubaccountUsernamesCreateObject = {
  username: 'EvilRaisin11',
  subaccountId: defaultSubaccountIdDefaultWalletAddress,
};

export const subaccountUsernameWithAlternativeAddress: SubaccountUsernamesCreateObject = {
  username: 'HonestRaisin32',
  subaccountId: defaultSubaccountIdWithAlternateAddress,
};

// ============== Leaderboard pnl Data ==============

export const defaultLeaderboardPnlOneDay: LeaderboardPnlCreateObject = {
  address: defaultAddress,
  timeSpan: 'ONE_DAY',
  pnl: '10000',
  currentEquity: '1000',
  rank: 1,
};

export const defaultLeaderboardPnl2OneDay: LeaderboardPnlCreateObject = {
  address: defaultAddress2,
  timeSpan: 'ONE_DAY',
  pnl: '100',
  currentEquity: '10000',
  rank: 2,
};

export const defaultLeaderboardPnl1AllTime: LeaderboardPnlCreateObject = {
  address: defaultAddress,
  timeSpan: 'ALL_TIME',
  pnl: '10000',
  currentEquity: '1000',
  rank: 1,
};

export const defaultLeaderboardPnlOneDayToUpsert: LeaderboardPnlCreateObject = {
  address: defaultAddress,
  timeSpan: 'ONE_DAY',
  pnl: '100000',
  currentEquity: '1000',
  rank: 1,
};

// ============== Affiliate referred users data ==============
export const defaultAffiliateReferredUser: AffiliateReferredUsersCreateObject = {
  affiliateAddress: defaultAddress,
  refereeAddress: defaultAddress2,
  referredAtBlock: '1',
};

// ============== Persistent cache Data ==============

export const defaultKV: PersistentCacheCreateObject = {
  key: 'someKey',
  value: 'someValue',
};

export const defaultKV2: PersistentCacheCreateObject = {
  key: 'otherKey',
  value: 'otherValue',
};

// ============== Affiliate Info Data ==============

export const defaultAffiliateInfo: AffiliateInfoCreateObject = {
  address: defaultAddress,
  affiliateEarnings: '10',
  referredMakerTrades: 10,
  referredTakerTrades: 20,
  totalReferredMakerFees: '10',
  totalReferredTakerFees: '10',
  totalReferredMakerRebates: '-10',
  totalReferredUsers: 5,
  firstReferralBlockHeight: '1',
  referredTotalVolume: '1000',
};

export const defaultAffiliateInfo2: AffiliateInfoCreateObject = {
  address: defaultWalletAddress,
  affiliateEarnings: '11',
  referredMakerTrades: 11,
  referredTakerTrades: 21,
  totalReferredMakerFees: '11',
  totalReferredTakerFees: '11',
  totalReferredMakerRebates: '-11',
  totalReferredUsers: 5,
  firstReferralBlockHeight: '11',
  referredTotalVolume: '1000',
};

export const defaultAffiliateInfo3: AffiliateInfoCreateObject = {
  address: defaultAddress2,
  affiliateEarnings: '12',
  referredMakerTrades: 12,
  referredTakerTrades: 22,
  totalReferredMakerFees: '12',
  totalReferredTakerFees: '12',
  totalReferredMakerRebates: '-12',
  totalReferredUsers: 10,
  firstReferralBlockHeight: '12',
  referredTotalVolume: '1111111',
};

// ==============  Tokens  =============

export const defaultFirebaseNotificationToken = {
  token: 'DEFAULT_TOKEN',
  address: defaultAddress,
  language: 'en',
  updatedAt: createdDateTime.toISO(),
};

// ==============  Vaults  =============

export const defaultVaultAddress: string = 'dydx1pzaql7h3tkt9uet8yht80me5td6gh0aprf58yk';

export const defaultVault: VaultCreateObject = {
  address: defaultVaultAddress,
  clobPairId: '0',
  status: VaultStatus.QUOTING,
  createdAt: createdDateTime.toISO(),
  updatedAt: createdDateTime.toISO(),
};

// ============== Funding Payments ==============

export const defaultFundingPayment = {
  subaccountId: defaultSubaccountId,
  createdAt: DateTime.utc().toISO(),
  createdAtHeight: '1',
  perpetualId: defaultPerpetualMarket.id,
  ticker: defaultPerpetualMarket.ticker,
  oraclePrice: '50000',
  size: '1',
  side: PositionSide.LONG,
  rate: '0.0001',
  payment: '5',
  fundingIndex: '5',
};

export const defaultFundingPayment2 = {
  subaccountId: defaultSubaccountId2,
  createdAt: DateTime.utc().toISO(),
  createdAtHeight: '2',
  perpetualId: defaultPerpetualMarket2.id,
  ticker: defaultPerpetualMarket2.ticker,
  oraclePrice: '3000',
  size: '2',
  side: PositionSide.SHORT,
  rate: '0.0002',
  payment: '-1.2',
  fundingIndex: '5.6',
};

// ==============  PnL  =================

export const defaultPnl: PnlCreateObject = {
  subaccountId: defaultSubaccountId,
  createdAt: defaultBlock.time,
  createdAtHeight: defaultBlock.blockHeight,
  equity: '10031.25',
  netTransfers: '10000.00',
  totalPnl: '31.25',
};

export const defaultPnl2: PnlCreateObject = {
  subaccountId: defaultSubaccountId,
  createdAt: DateTime.utc(2022, 6, 5).toISO(),
  createdAtHeight: '5',
  equity: '10013.00',
  netTransfers: '10000.00',
  totalPnl: '13.00',
};

// ============== Bridge Information ==============
export const defaultBridgeInformation: BridgeInformationCreateObject = {
  from_address: '0x1234567890abcdef1234567890abcdef12345678',
  chain_id: 'ethereum',
  amount: '1000000',
  created_at: createdDateTime.toISO(),
};

export const defaultBridgeInformation2: BridgeInformationCreateObject = {
  from_address: '0x9876543210fedcba9876543210fedcba98765432',
  chain_id: 'polygon',
  amount: '2000000',
  transaction_hash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890',
  created_at: DateTime.utc(2023, 1, 2).toISO(),
};

export const defaultBridgeInformation3: BridgeInformationCreateObject = {
  from_address: '0x1234567890abcdef1234567890abcdef12345678',
  chain_id: 'avalanche',
  amount: '3000000',
  created_at: DateTime.utc(2023, 1, 3).toISO(),
};
