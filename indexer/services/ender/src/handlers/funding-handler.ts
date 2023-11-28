import { logger } from '@dydxprotocol-indexer/base';
import {
  FundingIndexUpdatesTable,
  PerpetualMarketFromDatabase,
  TendermintEventTable,
  protocolTranslations,
  storeHelpers,
  PerpetualMarketModel,
} from '@dydxprotocol-indexer/postgres';
import { NextFundingCache } from '@dydxprotocol-indexer/redis';
import {
  FundingEventV1,
  FundingEventV1_Type,
  FundingUpdateV1,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import * as pg from 'pg';

import { redisClient } from '../helpers/redis/redis-controller';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent, FundingEventMessage } from '../lib/types';
import { Handler } from './handler';

export class FundingHandler extends Handler<FundingEventMessage> {
  eventType: string = 'FundingEvent';
  transactionIndex: number = indexerTendermintEventToTransactionIndex(
    this.indexerTendermintEvent,
  );

  transactionHash: string = this.block.txHashes[this.transactionIndex];
  eventId: Buffer = TendermintEventTable.createEventId(
    this.block.height.toString(),
    indexerTendermintEventToTransactionIndex(this.indexerTendermintEvent),
    this.indexerTendermintEvent.eventIndex,
  );

  public getParallelizationIds(): string[] {
    const ids: string[] = [];
    _.forEach(this.event.updates, (fundingIndexUpdate: FundingUpdateV1) => {
      const id: string = FundingIndexUpdatesTable.uuid(
        this.transactionHash,
        this.eventId,
        fundingIndexUpdate.perpetualId.toString(),
      );
      ids.push(`${this.eventType}_${id}`);
    });
    return ids;
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_funding_handler(
        ${this.block.height},
        '${this.block.time?.toISOString()}',
        '${JSON.stringify(FundingEventV1.decode(eventDataBinary))}',
        ${this.indexerTendermintEvent.eventIndex},
        ${transactionIndex}
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'FundingHandler#internalHandle',
        message: 'Failed to handle FundingEventV1',
        error,
      });

      throw error;
    });

    const perpetualMarkets:
    Map<string, PerpetualMarketFromDatabase> = new Map<string, PerpetualMarketFromDatabase>();
    for (const [key, perpetualMarket] of Object.entries(result.rows[0].result.perpetual_markets)) {
      perpetualMarkets.set(
        key,
        PerpetualMarketModel.fromJson(perpetualMarket as object) as PerpetualMarketFromDatabase,
      );
    }

    const promises: Promise<number>[] = new Array<Promise<number>>(this.event.updates.length);

    for (let i: number = 0; i < this.event.updates.length; i++) {
      const update: FundingUpdateV1 = this.event.updates[i];
      if (result.rows[0].result.errors[i] != null) {
        logger.error({
          at: 'FundingHandler#handleFundingSample',
          message: result.rows[0].result.errors[i],
          update,
        });
        continue;
      }

      const perpetualMarket:
      PerpetualMarketFromDatabase | undefined = perpetualMarkets.get(update.perpetualId.toString());
      if (perpetualMarket === undefined) {
        logger.error({
          at: 'FundingHandler#handleFundingSample',
          message: 'Received FundingUpdate with unknown perpetualId.',
          update,
        });
        continue;
      }

      switch (this.event.type) {
        case FundingEventV1_Type.TYPE_PREMIUM_SAMPLE:
          promises[i] = NextFundingCache.addFundingSample(
            perpetualMarket.ticker,
            new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(update.fundingValuePpm)),
            redisClient,
          );
          break;
        case FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX:
          // clear the cache for the predicted next funding rate
          promises[i] = NextFundingCache.clearFundingSamples(perpetualMarket.ticker, redisClient);
          break;
        default:
          logger.error({
            at: 'FundingHandler#handle',
            message: 'Received unknown FundingEvent type.',
            event: this.event,
          });
      }
    }

    await Promise.all(promises);
    return [];
  }
}
