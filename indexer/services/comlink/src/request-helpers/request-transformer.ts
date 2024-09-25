import {
  apiTranslations,
  AssetPositionFromDatabase,
  BestEffortOpenedStatus,
  CandleColumns,
  CandleFromDatabase,
  FillFromDatabase,
  fillTypeToTradeType,
  FundingIndexUpdatesFromDatabase,
  helpers,
  LiquidityTiersFromDatabase,
  MarketFromDatabase,
  MarketsMap,
  OrderFromDatabase,
  OrderType,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PnlTicksFromDatabase,
  PositionSide,
  protocolTranslations,
  SubaccountFromDatabase,
  SubaccountTable,
  TimeInForce,
  TransferFromDatabase,
  TransferType,
  parentSubaccountHelpers,
  YieldParamsFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { OrderbookLevels, PriceLevel } from '@dydxprotocol-indexer/redis';
import { RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';

import {
  AssetById,
  AssetPositionResponseObject,
  AssetPositionsMap,
  CandleResponseObject,
  FillResponseObject,
  HistoricalFundingResponseObject,
  MarketAndTypeByClobPairId,
  OrderbookResponseObject,
  ParentSubaccountTransferResponseObject,
  OrderbookResponsePriceLevel,
  OrderResponseObject,
  PerpetualMarketResponseObject,
  PerpetualPositionResponseObject,
  PerpetualPositionsMap,
  PerpetualPositionWithFunding,
  PnlTicksResponseObject,
  PostgresOrderMap,
  RedisOrderMap,
  SparklineResponseObject,
  SubaccountById,
  SubaccountResponseObject,
  TradeResponseObject,
  TransferResponseObject,
  YieldParamsResponseObject,
} from '../types';
import { oneAssetYieldIndex } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

/**
 * @description Converts perpetual position objects from the database into response objects.
 * Calculates realized using entry/exit/sum open/sum close values on the position in addition to the
 * net funding payments.
 * Calculates unrealized pnl using entry/size values on the position and index prices of markets.
 * @param position Perpetual position object with unsettled funding property
 * @param perpetualMarketsMap Map of perpetual ids to perpetual market database objects.
 * @returns Position response object.
 */
export function perpetualPositionToResponseObject(
  position: PerpetualPositionWithFunding,
  perpetualMarketsMap: PerpetualMarketsMap,
  marketsMap: MarketsMap,
): PerpetualPositionResponseObject {
  // Realized pnl is calculated from the difference in price between the average entry/exit price
  // (order depending on side of the position) multiplied by amount of the position that was closed
  // in addition to the funding payments.
  const priceDiff: Big = (position.side === PositionSide.LONG)
    ? Big(position.exitPrice ?? 0).minus(position.entryPrice)
    : Big(position.entryPrice).minus(position.exitPrice ?? 0);
  const netFunding: Big = Big(position.settledFunding).plus(position.unsettledFunding);
  const realizedPnl: string = priceDiff
    .mul(position.sumClose)
    .plus(netFunding)
    .toFixed();

  return {
    market: perpetualMarketsMap[position.perpetualId].ticker,
    status: position.status,
    side: position.side,
    size: position.size,
    maxSize: position.maxSize,
    entryPrice: Big(position.entryPrice).toFixed(),
    exitPrice: position.exitPrice && Big(position.exitPrice).toFixed(),
    realizedPnl,
    unrealizedPnl: helpers.getUnrealizedPnl(
      position, perpetualMarketsMap[position.perpetualId], marketsMap,
    ),
    createdAt: position.createdAt,
    createdAtHeight: position.createdAtHeight,
    closedAt: position.closedAt,
    sumOpen: position.sumOpen,
    sumClose: position.sumClose,
    netFunding: netFunding.toFixed(),
    perpYieldIndex: position.perpYieldIndex,
  };
}

/**
 * @description Converts asset position objects from the database into response objects.
 *
 * @param position Asset position object from the database.
 * @param assetMap Map of asset ids to asset database objects.
 * @returns Asset position response object.
 */
export function assetPositionToResponseObject(
  position: AssetPositionFromDatabase,
  assetMap: AssetById,
  subaccountNumber: number,
): AssetPositionResponseObject {

  return {
    symbol: assetMap[position.assetId].symbol,
    side: position.isLong ? PositionSide.LONG : PositionSide.SHORT,
    size: position.size,
    assetId: position.assetId,
    subaccountNumber,
  };
}

/**
 * @description Converts fill objects from the database into response objects.
 * @param fill Fill object from database.
 * @param marketsByClobPairId Map of market tickers and market types to clob pair ids.
 * @returns Fill response object.
 */
export function fillToResponseObject(
  fill: FillFromDatabase,
  marketsByClobPairId: MarketAndTypeByClobPairId,
  subaccountNumber: number,
): FillResponseObject {
  return {
    id: fill.id,
    side: fill.side,
    liquidity: fill.liquidity,
    type: fill.type,
    market: marketsByClobPairId[fill.clobPairId].market,
    marketType: marketsByClobPairId[fill.clobPairId].marketType,
    price: fill.price,
    size: fill.size,
    fee: fill.fee,
    createdAt: fill.createdAt,
    createdAtHeight: fill.createdAtHeight,
    orderId: fill.orderId,
    clientMetadata: fill.clientMetadata,
    subaccountNumber,
  };
}

/**
 *
 * Converts funding rate objects from the database into response objects.
 * Use the ticker param as the ticker. This should map to the perpetualId in
 * the funding object.
 *
 * @param funding
 * @param ticker
 */
export function historicalFundingToResponseObject(
  funding: FundingIndexUpdatesFromDatabase,
  ticker: string,
): HistoricalFundingResponseObject {
  return {
    ticker,
    rate: funding.rate,
    price: funding.oraclePrice,
    effectiveAtHeight: funding.effectiveAtHeight,
    effectiveAt: funding.effectiveAt,
  };
}

export function fillToTradeResponseObject(
  fill: FillFromDatabase,
): TradeResponseObject {
  return {
    id: fill.eventId.toString('hex'),
    side: fill.side,
    size: fill.size,
    price: fill.price,
    type: fillTypeToTradeType(fill.type),
    createdAt: fill.createdAt,
    createdAtHeight: fill.createdAtHeight,
  };
}

/**
 * Converts transfer from the database into API response.
 *
 * @param transfer
 * @param assetMap map of assetId to symbol.
 * @param subaccountMap map of subaccountId to subaccounts for all subaccounts involved in
 * transfers.
 * @param subaccountId represents the subaccountId in the query. This is used to determine the
 * transfer type.
 */
export function transferToResponseObject(
  transfer: TransferFromDatabase,
  assetMap: AssetById,
  subaccountMap: SubaccountById,
  subaccountId: string,
): TransferResponseObject {
  return {
    id: transfer.id,
    sender: {
      address: transfer.senderWalletAddress ?? subaccountMap[transfer.senderSubaccountId!].address,
      subaccountNumber: transfer.senderWalletAddress ? undefined
        : subaccountMap[transfer.senderSubaccountId!].subaccountNumber,
    },
    recipient: {
      address: transfer.recipientWalletAddress ?? subaccountMap[
        transfer.recipientSubaccountId!
      ].address,
      subaccountNumber: transfer.recipientWalletAddress ? undefined
        : subaccountMap[transfer.recipientSubaccountId!].subaccountNumber,
    },
    size: transfer.size,
    createdAt: transfer.createdAt,
    createdAtHeight: transfer.createdAtHeight,
    symbol: assetMap[transfer.assetId].symbol,
    type: helpers.getTransferType(transfer, subaccountId),
    transactionHash: transfer.transactionHash,
  };
}

export function transferToParentSubaccountResponseObject(
  transfer: TransferFromDatabase,
  assetMap: AssetById,
  subaccountMap: SubaccountById,
  parentSubaccountNumber: number,
): ParentSubaccountTransferResponseObject {

  const senderParentSubaccountNum = transfer.senderWalletAddress
    ? undefined
    : parentSubaccountHelpers.getParentSubaccountNum(
      subaccountMap[transfer.senderSubaccountId!].subaccountNumber,
    );

  const recipientParentSubaccountNum = transfer.recipientWalletAddress
    ? undefined
    : parentSubaccountHelpers.getParentSubaccountNum(
      subaccountMap[transfer.recipientSubaccountId!].subaccountNumber,
    );

  // Determine transfer type based on parent subaccount number.
  let transferType: TransferType = TransferType.TRANSFER_IN;
  if (senderParentSubaccountNum === parentSubaccountNumber) {
    if (transfer.recipientSubaccountId) {
      transferType = TransferType.TRANSFER_OUT;
    } else {
      transferType = TransferType.WITHDRAWAL;
    }
  } else if (recipientParentSubaccountNum === parentSubaccountNumber) {
    if (transfer.senderSubaccountId) {
      transferType = TransferType.TRANSFER_IN;
    } else {
      transferType = TransferType.DEPOSIT;
    }
  }

  return {
    id: transfer.id,
    sender: {
      address: transfer.senderWalletAddress ?? subaccountMap[transfer.senderSubaccountId!].address,
      parentSubaccountNumber: senderParentSubaccountNum,
    },
    recipient: {
      address: transfer.recipientWalletAddress ?? subaccountMap[
        transfer.recipientSubaccountId!
      ].address,
      parentSubaccountNumber: recipientParentSubaccountNum,
    },
    size: transfer.size,
    createdAt: transfer.createdAt,
    createdAtHeight: transfer.createdAtHeight,
    symbol: assetMap[transfer.assetId].symbol,
    type: transferType,
    transactionHash: transfer.transactionHash,
  };
}

export function pnlTicksToResponseObject(
  pnlTicks: PnlTicksFromDatabase,
): PnlTicksResponseObject {
  return {
    id: pnlTicks.id,
    subaccountId: pnlTicks.subaccountId,
    equity: pnlTicks.equity,
    totalPnl: pnlTicks.totalPnl,
    netTransfers: pnlTicks.netTransfers,
    createdAt: pnlTicks.createdAt,
    blockHeight: pnlTicks.blockHeight,
    blockTime: pnlTicks.blockTime,
  };
}

export function subaccountToResponseObject({
  subaccount,
  equity,
  freeCollateral,
  openPerpetualPositions = {},
  assetPositions = {},
}: {
  subaccount: SubaccountFromDatabase,
  equity: string,
  freeCollateral: string,
  openPerpetualPositions: PerpetualPositionsMap,
  assetPositions: AssetPositionsMap,
}): SubaccountResponseObject {
  return {
    address: subaccount.address,
    subaccountNumber: subaccount.subaccountNumber,
    equity: Big(equity).toFixed(),
    freeCollateral: Big(freeCollateral).toFixed(),
    openPerpetualPositions,
    assetPositions,
    // TODO(DEC-687): Track `marginEnabled` for subaccounts.
    marginEnabled: true,
    assetYieldIndex: subaccount.assetYieldIndex,
  };
}

export function perpetualMarketToResponseObject(
  perpetualMarket: PerpetualMarketFromDatabase,
  liquidityTier: LiquidityTiersFromDatabase,
  market: MarketFromDatabase,
): PerpetualMarketResponseObject {
  return {
    clobPairId: perpetualMarket.clobPairId,
    ticker: perpetualMarket.ticker,
    status: perpetualMarket.status,
    spotPrice: market.spotPrice!,
    pnlPrice: market.pnlPrice!,
    priceChange24H: perpetualMarket.priceChange24H,
    volume24H: perpetualMarket.volume24H,
    trades24H: perpetualMarket.trades24H,
    nextFundingRate: perpetualMarket.nextFundingRate,
    initialMarginFraction: helpers.ppmToString(Number(liquidityTier.initialMarginPpm)),
    maintenanceMarginFraction: helpers.ppmToString(
      helpers.getMaintenanceMarginPpm(
        Number(liquidityTier.initialMarginPpm),
        Number(liquidityTier.maintenanceFractionPpm),
      ),
    ),
    openInterest: perpetualMarket.openInterest,
    atomicResolution: perpetualMarket.atomicResolution,
    dangerIndexPpm: perpetualMarket.dangerIndexPpm,
    isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock:
      perpetualMarket.isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
    quantumConversionExponent: perpetualMarket.quantumConversionExponent,
    tickSize: protocolTranslations.getTickSize(perpetualMarket),
    stepSize: protocolTranslations.getStepSize(perpetualMarket),
    stepBaseQuantums: perpetualMarket.stepBaseQuantums,
    subticksPerTick: perpetualMarket.subticksPerTick,
    marketType: perpetualMarket.marketType,
    openInterestLowerCap: liquidityTier.openInterestLowerCap,
    openInterestUpperCap: liquidityTier.openInterestUpperCap,
    baseOpenInterest: perpetualMarket.baseOpenInterest,
    perpYieldIndex: perpetualMarket.perpYieldIndex,
  };
}

export function OrderbookLevelsToResponseObject(
  orderbookLevels: OrderbookLevels,
  perpetualMarket: PerpetualMarketFromDatabase,
): OrderbookResponseObject {
  return {
    bids: OrderbookPriceLevelsToResponsePriceLevels(orderbookLevels.bids, perpetualMarket),
    asks: OrderbookPriceLevelsToResponsePriceLevels(orderbookLevels.asks, perpetualMarket),
  };
}

function OrderbookPriceLevelsToResponsePriceLevels(
  priceLevels: PriceLevel[],
  perpetualMarket: PerpetualMarketFromDatabase,
): OrderbookResponsePriceLevel[] {
  return priceLevels.map((level: PriceLevel) => {
    return {
      price: level.humanPrice,
      size: protocolTranslations.quantumsToHumanFixedString(
        level.quantums,
        perpetualMarket.atomicResolution,
      ),
    };
  });
}

export function mergePostgresAndRedisOrdersToResponseObjects(
  postgresOrderMap: PostgresOrderMap,
  redisOrderMap: RedisOrderMap,
  subaccountIdToNumber: Record<string, number>,
): OrderResponseObject[] {
  const orderIds: string[] = _.uniq(
    Object.keys(redisOrderMap).concat(Object.keys(postgresOrderMap)),
  );

  return _.map(orderIds, (orderId: string) => {
    return postgresAndRedisOrderToResponseObject(
      postgresOrderMap[orderId],
      subaccountIdToNumber,
      redisOrderMap[orderId],
    ) as OrderResponseObject;
  });
}

/**
 * Returns undefined if postgres and redis are both undefined/null.
 * If only redis is defined, then generate the response object from the redisOrder.
 * If only postgres is defined, then generate the response object from the postgresOrder.
 * If both postgres and redis are defined, then generate the response object from postgresOrder
 * and override the size, price, and goodTilBlock fields with the redisOrder.
 * @param postgresOrder
 * @param subaccountIdToNumber
 * @param redisOrder
 * @returns
 */
export function postgresAndRedisOrderToResponseObject(
  postgresOrder: OrderFromDatabase | undefined,
  subaccountIdToNumber: Record<string, number>,
  redisOrder?: RedisOrder | null,
): OrderResponseObject | undefined {
  if (postgresOrder === undefined) {
    if (redisOrder === null || redisOrder === undefined) {
      return undefined;
    }

    return redisOrderToResponseObject(redisOrder);
  }

  const orderResponse: OrderResponseObject = postgresOrderToResponseObject(
    postgresOrder,
    subaccountIdToNumber[postgresOrder.subaccountId],
  );
  if (redisOrder === null || redisOrder === undefined) {
    return orderResponse;
  }
  const redisOrderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
    redisOrder.order!.timeInForce,
  );

  return {
    ...orderResponse,
    size: redisOrder.size,
    price: redisOrder.price,
    goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)?.toString() ?? undefined,
    goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!) ?? undefined,
    timeInForce: apiTranslations.orderTIFToAPITIF(redisOrderTIF),
    postOnly: apiTranslations.isOrderTIFPostOnly(redisOrderTIF),
    reduceOnly: redisOrder.order!.reduceOnly,
  };
}

export function postgresOrderToResponseObject(
  order: OrderFromDatabase,
  subaccountNumber: number,
): OrderResponseObject {
  return {
    ...order,
    timeInForce: apiTranslations.orderTIFToAPITIF(order.timeInForce),
    postOnly: apiTranslations.isOrderTIFPostOnly(order.timeInForce),
    goodTilBlock: order.goodTilBlock ?? undefined,
    goodTilBlockTime: order.goodTilBlockTime ?? undefined,
    createdAtHeight: order.createdAtHeight ?? undefined,
    ticker: perpetualMarketRefresher.getPerpetualMarketTicker(order.clobPairId)!,
    triggerPrice: order.triggerPrice ?? undefined,
    subaccountNumber,
  };
}

export function redisOrderToResponseObject(
  redisOrder: RedisOrder,
): OrderResponseObject {
  const clobPairId: string = redisOrder.order!.orderId!.clobPairId.toString();
  const orderTIF: TimeInForce = protocolTranslations.protocolOrderTIFToTIF(
    redisOrder.order!.timeInForce,
  );
  return {
    id: redisOrder.id,
    subaccountId: SubaccountTable.subaccountIdToUuid(redisOrder.order!.orderId!.subaccountId!),
    clientId: redisOrder.order!.orderId!.clientId.toString(),
    clobPairId,
    side: protocolTranslations.protocolOrderSideToOrderSide(redisOrder.order!.side),
    size: redisOrder.size,
    totalFilled: '0',
    price: redisOrder.price,
    type: OrderType.LIMIT,
    status: BestEffortOpenedStatus.BEST_EFFORT_OPENED,
    timeInForce: apiTranslations.orderTIFToAPITIF(orderTIF),
    postOnly: apiTranslations.isOrderTIFPostOnly(orderTIF),
    reduceOnly: redisOrder.order!.reduceOnly,
    goodTilBlock: protocolTranslations.getGoodTilBlock(redisOrder.order!)
      ?.toString() ?? undefined,
    goodTilBlockTime: protocolTranslations.getGoodTilBlockTime(redisOrder.order!) ?? undefined,
    ticker: perpetualMarketRefresher.getPerpetualMarketTicker(clobPairId)!,
    orderFlags: redisOrder.order!.orderId!.orderFlags.toString(),
    clientMetadata: redisOrder.order!.clientMetadata.toString(),
    subaccountNumber: redisOrder.order!.orderId!.subaccountId!.number,
  };
}

export function candleToResponseObject(
  candle: CandleFromDatabase,
): CandleResponseObject {
  return _.omit(candle, [CandleColumns.id]);
}

export function candlesToSparklineResponseObject(
  tickers: string[],
  unsortedTickerCandles: CandleFromDatabase[],
  limit: number,
): SparklineResponseObject {
  const response: SparklineResponseObject = _.fromPairs(
    _.map(tickers, (ticker: string) => [ticker, []]),
  );
  return _.reduce(
    unsortedTickerCandles,
    (accumulator: { [ticker: string]: string[] }, candle: CandleFromDatabase) => {
      if (accumulator[candle.ticker].length < limit) {
        accumulator[candle.ticker].push(candle[CandleColumns.close]);
      }
      // Do not add to accumulator if accumulator length is already at limit.
      // Since candles are sorted by startedAt in descending order, the first 'limit' candles
      // will be the most recent candles.
      return accumulator;
    }, response,
  );
}

export function yieldParamsToResponseObject(
  yieldParams: YieldParamsFromDatabase,
): YieldParamsResponseObject {
  return {
    id: yieldParams.id,
    sDAIPrice: yieldParams.sDAIPrice,
    assetYieldIndex: yieldParams.assetYieldIndex,
    createdAt: yieldParams.createdAt,
    createdAtHeight: yieldParams.createdAtHeight,
  }
}
