import {
  logger,
  stats,
  wrapBackgroundTask,
} from '@dydxprotocol-indexer/base';
import { MARKETS_WEBSOCKET_MESSAGE_VERSION, producer, WebsocketTopics } from '@dydxprotocol-indexer/kafka';
import {
  LiquidityTiersFromDatabase,
  LiquidityTiersMap,
  MarketMessageContents,
  PerpetualMarketFromDatabase,
  PerpetualMarketsMap,
  TradingMarketMessageContents,
  helpers,
} from '@dydxprotocol-indexer/postgres';
import {
  MarketMessage,
  OffChainUpdateV1,
  IndexerOrderId,
  OrderRemovalReason,
  OrderRemoveV1_OrderRemovalStatus,
} from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';

import config from '../config';

export function getMarginFields(
  perpetualMarket: PerpetualMarketFromDatabase,
  liquidityTiers: LiquidityTiersMap,
): {
  initialMarginFraction: string,
  maintenanceMarginFraction: string,
} {

  if (liquidityTiers[perpetualMarket.liquidityTierId] === undefined) {
    throw new Error(`Liquidity tier ${perpetualMarket.liquidityTierId} for perpetual market ${perpetualMarket.ticker} not found`);
  }
  const liquidityTier: LiquidityTiersFromDatabase = liquidityTiers[perpetualMarket.liquidityTierId];
  return {
    initialMarginFraction: helpers.ppmToString(Number(liquidityTier.initialMarginPpm)),
    maintenanceMarginFraction: helpers.ppmToString(
      helpers.getMaintenanceMarginPpm(
        Number(liquidityTier.initialMarginPpm),
        Number(liquidityTier.maintenanceFractionPpm),
      ),
    ),
  };
}

export function getUpdatedMarkets(
  oldMarkets: PerpetualMarketsMap,
  newMarkets: PerpetualMarketsMap,
  liquidityTiers: LiquidityTiersMap,
): MarketMessageContents {
  // Get only modified markets and only modified fields from within
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const updatedMarkets: TradingMarketMessageContents = {};
  Object.keys(newMarkets).forEach((perpetualMarketId: string) => {
    const diffPairs = _.differenceWith(
      _.toPairs({
        ...newMarkets[perpetualMarketId],
        ...getMarginFields(newMarkets[perpetualMarketId], liquidityTiers),
      }),
      _.toPairs({
        ...oldMarkets[perpetualMarketId],
        ...getMarginFields(oldMarkets[perpetualMarketId], liquidityTiers),
      }),
      _.isEqual,
    );
    if (diffPairs.length !== 0) {
      updatedMarkets[newMarkets[perpetualMarketId].ticker] = _.fromPairs(diffPairs);
    }
  });
  return {
    trading: updatedMarkets,
  };
}

export function compareAndHandleMarketsWebsocketMessage({
  oldMarkets,
  newMarkets,
  liquidityTiers,
}: {
  oldMarkets: PerpetualMarketsMap,
  newMarkets: PerpetualMarketsMap,
  liquidityTiers: LiquidityTiersMap,
}): void {
  // Get only modified markets and only modified fields from within
  const updatedMarkets: MarketMessageContents = getUpdatedMarkets(
    oldMarkets,
    newMarkets,
    liquidityTiers,
  );

  if (!_.isEmpty(updatedMarkets.trading)) {
    logger.info({
      at: 'websocket#compareAndHandleMarketsWebsocketMessage',
      message: 'Sending markets websocket message',
      updatedMarkets,
    });

    const marketMessage: MarketMessage = {
      contents: JSON.stringify(updatedMarkets),
      version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
    };
    const buffer: Buffer = Buffer.from(
      MarketMessage.encode(marketMessage).finish(),
    );

    wrapBackgroundTask(publishToMarketsWebsocket(
      buffer,
    ), false, 'publishToMarketsWebsocket');
  }
}

async function publishToMarketsWebsocket(buffer: Buffer) {
  const startUpdate: number = Date.now();
  try {
    await producer.send({
      topic: WebsocketTopics.TO_WEBSOCKETS_MARKETS,
      messages: [{
        value: buffer,
      }],
    });
  } catch (error) {
    logger.error({
      at: 'websocket#publishToMarketsWebsocket',
      message: 'Publish to markets websocket failed',
      error,
    });
  }

  stats.timing(`${config.SERVICE_NAME}.publish_markets_websocket_updates_duration`,
    Date.now() - startUpdate);
  stats.gauge(`${config.SERVICE_NAME}.publish_markets_websocket`, 1);
}

export function getExpiredOffChainUpdateMessage(removedOrderId: IndexerOrderId): Buffer {
  const message: OffChainUpdateV1 = OffChainUpdateV1.fromPartial({
    orderRemove: {
      removedOrderId,
      reason: OrderRemovalReason.ORDER_REMOVAL_REASON_INDEXER_EXPIRED,
      removalStatus: OrderRemoveV1_OrderRemovalStatus.ORDER_REMOVAL_STATUS_CANCELED,
    },
  });
  return Buffer.from(OffChainUpdateV1.encode(message).finish());
}
