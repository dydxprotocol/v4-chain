import { getInstanceId, logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
} from '@dydxprotocol-indexer/kafka';
import { OrderbookMessageContents, PerpetualMarketFromDatabase, protocolTranslations } from '@dydxprotocol-indexer/postgres';
import { OffChainUpdateV1, OrderbookMessage, RedisOrder } from '@dydxprotocol-indexer/v4-protos';
import { IHeaders } from 'kafkajs';
import { OrderbookSide } from 'src/lib/types';

import { orderSideToOrderbookSide } from './helpers';

export abstract class Handler {
  public txHash?: string;

  public constructor(txHash?: string) {
    this.txHash = txHash;
  }

  protected abstract handle(update: OffChainUpdateV1, headers: IHeaders): Promise<void>;

  // TODO(DEC-1251): Add stats for message handling.
  public async handleUpdate(update: OffChainUpdateV1, headers: IHeaders): Promise<void> {
    return this.handle(update, headers);
  }

  protected logAndThrowParseMessageError(
    message: string,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    additionalParameters: any = {},
  ): void {
    logger.error({
      at: `${this.constructor.name}#logAndThrowParseMessageError`,
      message,
      txhash: this.txHash,
      ...additionalParameters,
    });
    throw new ParseMessageError(message);
  }

  protected createOrderbookWebsocketMessage(
    redisOrder: RedisOrder,
    perpetualMarket: PerpetualMarketFromDatabase,
    updatedQuantums: number,
  ): Buffer {
    const orderbookSide: OrderbookSide = orderSideToOrderbookSide(redisOrder.order!.side);
    const humanSize: string = protocolTranslations.quantumsToHumanFixedString(
      updatedQuantums.toString(),
      perpetualMarket.atomicResolution,
    );
    const contents: OrderbookMessageContents = {
      [orderbookSide]: [[redisOrder.price, humanSize]],
    };

    const orderbookMessage: OrderbookMessage = OrderbookMessage.fromPartial({
      clobPairId: perpetualMarket.clobPairId,
      contents: JSON.stringify(contents),
      version: ORDERBOOKS_WEBSOCKET_MESSAGE_VERSION,
    });

    return Buffer.from(OrderbookMessage.encode(orderbookMessage).finish());
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected generateTimingStatsOptions(fnName: string): any {
    return {
      className: this.constructor.name,
      fnName,
      instance: getInstanceId(),
    };
  }
}
