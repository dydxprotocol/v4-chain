import {
  apiTranslations,
  AssetFromDatabase,
  AssetPositionFromDatabase,
  AssetPositionSubaccountMessageContents,
  AssetsMap,
  FillFromDatabase,
  FillSubaccountMessageContents,
  helpers,
  liquidityTierRefresher,
  LiquidityTiersFromDatabase,
  LiquidityTiersMap,
  MarketFromDatabase,
  MarketMessageContents,
  OraclePriceFromDatabase,
  OrderFromDatabase,
  OrderSubaccountMessageContents,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketsMap,
  PerpetualPositionFromDatabase,
  PerpetualPositionSubaccountMessageContents,
  PositionSide,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  TradingMarketMessageContents,
  TradingPerpetualMarketMessage,
  TransferFromDatabase,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { SubaccountId } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';

/**
 * Adds the positions to the contents for the positions message in the
 * subaccount kafka channel.
 *
 * @param contents
 * @param subaccountIdProto
 * @param perpetualPositions
 * @param perpetualMarketsMapping
 * @param assetPositions
 * @param assetsMap
 */
export function addPositionsToContents(
  contents: SubaccountMessageContents,
  subaccountIdProto: SubaccountId,
  updateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[],
  perpetualMarketsMapping: PerpetualMarketsMap,
  assetPositions: AssetPositionFromDatabase[],
  assetsMap: AssetsMap,
  blockHeight: string,
): SubaccountMessageContents {
  return {
    ...contents,
    perpetualPositions: updateObjects.length === 0 ? undefined : generatePerpetualPositionsContents(
      subaccountIdProto,
      updateObjects,
      perpetualMarketsMapping,
    ),
    assetPositions: assetPositions.length === 0 ? undefined : generateAssetPositionsContents(
      subaccountIdProto,
      assetPositions,
      assetsMap,
    ),
    blockHeight,
  };
}

export function generatePerpetualPositionsContents(
  subaccountIdProto: SubaccountId,
  perpetualPositions: UpdatedPerpetualPositionSubaccountKafkaObject[],
  perpetualMarketsMapping: PerpetualMarketsMap,
): PerpetualPositionSubaccountMessageContents[] {
  return _.map(
    perpetualPositions,
    (perpetualPosition: UpdatedPerpetualPositionSubaccountKafkaObject):
    PerpetualPositionSubaccountMessageContents => {
      return {
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: perpetualPosition.id,
        market: perpetualMarketsMapping[perpetualPosition.perpetualId].ticker,
        side: perpetualPosition.side,
        status: perpetualPosition.status,
        size: perpetualPosition.size,
        maxSize: perpetualPosition.maxSize,
        // When a perpetual position update is sent from the protocol, there is 0 unsettled
        // for the position as all funding is settled on updates to the subaccount, which includes
        // updates to the perpetual position
        netFunding: perpetualPosition.settledFunding,
        entryPrice: perpetualPosition.entryPrice,
        exitPrice: perpetualPosition.exitPrice,
        sumOpen: perpetualPosition.sumOpen,
        sumClose: perpetualPosition.sumClose,
        realizedPnl: perpetualPosition.realizedPnl,
        unrealizedPnl: perpetualPosition.unrealizedPnl,
      };
    },
  );
}

export function generateAssetPositionsContents(
  subaccountIdProto: SubaccountId,
  assetPositions: AssetPositionFromDatabase[],
  assetsMap: AssetsMap,
): AssetPositionSubaccountMessageContents[] {
  return _.map(
    assetPositions,
    (position: AssetPositionFromDatabase): AssetPositionSubaccountMessageContents => {
      return {
        address: subaccountIdProto.owner,
        subaccountNumber: subaccountIdProto.number,
        positionId: position.id,
        assetId: position.assetId,
        symbol: assetsMap[position.assetId].symbol,
        side: position.isLong ? PositionSide.LONG : PositionSide.SHORT,
        size: position.size,
      };
    },
  );
}

/**
 * Gets the realized and unrealized pnl for a perpetual position.
 *
 * @param updateObject
 * @param perpetualMarket
 * @param marketIdToMarket
 */
// TODO: Move this to a shared package so code is not duplicated with comlink.
export function getPnl(
  updateObject: UpdatedPerpetualPositionSubaccountKafkaObject,
  perpetualMarket: PerpetualMarketFromDatabase,
  market: MarketFromDatabase,
): { realizedPnl: string | undefined, unrealizedPnl: string | undefined } {
  let realizedPnl: string | undefined;
  let unrealizedPnl: string | undefined;
  if (updateObject !== undefined) {
    const priceDiff: Big = (updateObject.side === PositionSide.LONG)
      ? Big(updateObject.exitPrice ?? 0).minus(updateObject.entryPrice)
      : Big(updateObject.entryPrice).minus(updateObject.exitPrice ?? 0);
    realizedPnl = priceDiff
      .mul(updateObject.sumClose)
      .plus(updateObject.settledFunding)
      .toFixed();
    unrealizedPnl = helpers.getUnrealizedPnl(
      updateObject,
      perpetualMarket,
      market,
    );
  }
  return { realizedPnl, unrealizedPnl };
}

/**
 * Annotates a perpetual position update with the realized and unrealized pnl.
 *
 * @param updateObject
 * @param perpetualMarketMap
 * @param marketIdToMarket
 */
export function annotateWithPnl(
  updateObject: UpdatedPerpetualPositionSubaccountKafkaObject,
  perpetualMarketMap: PerpetualMarketsMap,
  market: MarketFromDatabase,
): UpdatedPerpetualPositionSubaccountKafkaObject {
  return {
    ...updateObject,
    ...getPnl(updateObject, perpetualMarketMap[updateObject.perpetualId], market),
  };
}

/**
 * Converts a perpetual position from the database to the message format
 * used in the subaccount kafka channel.
 *
 * @param position
 */
export function convertPerpetualPosition(
  position: PerpetualPositionFromDatabase,
): UpdatedPerpetualPositionSubaccountKafkaObject {
  const {
    id,
    perpetualId,
    side,
    status,
    size,
    maxSize,
    entryPrice,
    exitPrice,
    sumOpen,
    sumClose,
    closedAt,
    closedAtHeight,
    closeEventId,
    lastEventId,
    settledFunding,
  } = position;

  const updatedPosition: UpdatedPerpetualPositionSubaccountKafkaObject = {
    perpetualId,
    maxSize,
    side,
    entryPrice,
    exitPrice,
    sumOpen,
    sumClose,
    id,
    closedAt: closedAt || null,
    closedAtHeight: closedAtHeight || null,
    closeEventId: closeEventId || null,
    lastEventId,
    settledFunding,
    status,
    size,
  };

  return updatedPosition;
}

/**
 * Generates the transfer contents for the transfer message in the
 * subaccount kafka channel.
 *
 * @param contents
 * @param transfer
 * @param asset associated with the asset id in transfer
 * @param subaccountId to generate the websocket message for
 * @param senderSubaccountId
 * @param recipientSubaccountId
 * @param blockHeight: latest block height processed by Indexer
 */
export function generateTransferContents(
  transfer: TransferFromDatabase,
  asset: AssetFromDatabase,
  subaccountId: SubaccountId,
  senderSubaccountId?: SubaccountId,
  recipientSubaccountId?: SubaccountId,
  blockHeight?: string,
): SubaccountMessageContents {
  return {
    transfers: {
      sender: {
        address: transfer.senderWalletAddress ?? senderSubaccountId!.owner,
        subaccountNumber: transfer.senderWalletAddress ? undefined
          : senderSubaccountId!.number,
      },
      recipient: {
        address: transfer.recipientWalletAddress ?? recipientSubaccountId!.owner,
        subaccountNumber: transfer.recipientWalletAddress ? undefined
          : recipientSubaccountId!.number,
      },
      symbol: asset.symbol,
      size: transfer.size,
      type: helpers.getTransferType(
        transfer,
        SubaccountTable.uuid(subaccountId.owner, subaccountId.number),
      ),
      createdAt: transfer.createdAt,
      createdAtHeight: transfer.createdAtHeight,
      transactionHash: transfer.transactionHash,
    },
    blockHeight,
  };
}

export function generateOraclePriceContents(
  oraclePrice: OraclePriceFromDatabase,
  ticker: string,
): MarketMessageContents {
  return {
    oraclePrices: {
      [ticker]: {
        oraclePrice: oraclePrice.price,
        effectiveAt: oraclePrice.effectiveAt,
        effectiveAtHeight: oraclePrice.effectiveAtHeight,
        marketId: oraclePrice.marketId,
      },
    },
  };
}

export function generateFillSubaccountMessage(
  fill: FillFromDatabase,
  ticker: string,
): FillSubaccountMessageContents {
  return {
    ...fill!,
    eventId: fill!.eventId.toString('hex'),
    ticker,
  };
}

export function generateOrderSubaccountMessage(
  order: OrderFromDatabase,
  ticker: string,
): OrderSubaccountMessageContents {
  return {
    ...order,
    timeInForce: apiTranslations.orderTIFToAPITIF(order.timeInForce),
    postOnly: apiTranslations.isOrderTIFPostOnly(order.timeInForce),
    goodTilBlock: order.goodTilBlock,
    goodTilBlockTime: order.goodTilBlockTime,
    ticker,
  };
}

export function generatePerpetualMarketMessage(
  perpetualMarkets: PerpetualMarketFromDatabase[],
): MarketMessageContents {
  const liquidityTierMap: LiquidityTiersMap = liquidityTierRefresher.getLiquidityTiersMap();

  const tradingMarketMessageContents: TradingMarketMessageContents = _.chain(perpetualMarkets)
    .keyBy(PerpetualMarketColumns.ticker)
    .mapValues((perpetualMarket: PerpetualMarketFromDatabase): TradingPerpetualMarketMessage => {
      const liquidityTier:
      LiquidityTiersFromDatabase = liquidityTierMap[perpetualMarket.liquidityTierId];

      return {
        id: perpetualMarket.id,
        clobPairId: perpetualMarket.clobPairId.toString(),
        ticker: perpetualMarket.ticker,
        marketId: perpetualMarket.marketId,
        status: perpetualMarket.status,
        quantumConversionExponent: perpetualMarket.quantumConversionExponent,
        atomicResolution: perpetualMarket.atomicResolution,
        subticksPerTick: perpetualMarket.subticksPerTick,
        stepBaseQuantums: perpetualMarket.stepBaseQuantums,
        marketType: perpetualMarket.marketType,
        initialMarginFraction: helpers.ppmToString(Number(liquidityTier.initialMarginPpm)),
        maintenanceMarginFraction: helpers.ppmToString(
          helpers.getMaintenanceMarginPpm(
            Number(liquidityTier.initialMarginPpm),
            Number(liquidityTier.maintenanceFractionPpm),
          ),
        ),
        openInterestLowerCap: liquidityTier.openInterestLowerCap,
        openInterestUpperCap: liquidityTier.openInterestUpperCap,
        tickSize: protocolTranslations.getTickSize(perpetualMarket),
        stepSize: protocolTranslations.getStepSize(perpetualMarket),
        priceChange24H: perpetualMarket.priceChange24H,
        volume24H: perpetualMarket.volume24H,
        trades24H: perpetualMarket.trades24H,
        nextFundingRate: perpetualMarket.nextFundingRate,
        openInterest: perpetualMarket.openInterest,
        baseOpenInterest: perpetualMarket.baseOpenInterest,
        defaultFundingRate1H: perpetualMarket.defaultFundingRate1H,
      };
    })
    .value();

  return {
    trading: tradingMarketMessageContents,
  };
}
