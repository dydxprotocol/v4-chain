import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  FundingIndexUpdatesTable,
  PerpetualMarketFromDatabase,
  TendermintEventTable,
  protocolTranslations,
  PerpetualMarketModel,
  FundingIndexUpdatesFromDatabase,
  FundingIndexUpdatesModel,
} from '@dydxprotocol-indexer/postgres';
import { NextFundingCache } from '@dydxprotocol-indexer/redis';
import { bytesToBigInt } from '@dydxprotocol-indexer/v4-proto-parser';
import {
  FundingEventV1_Type,
  FundingUpdateV1,
} from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
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
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarkets:
    Map<string, PerpetualMarketFromDatabase> = new Map<string, PerpetualMarketFromDatabase>();
    for (const [key, perpetualMarket] of Object.entries(resultRow.perpetual_markets)) {
      perpetualMarkets.set(
        key,
        PerpetualMarketModel.fromJson(perpetualMarket as object) as PerpetualMarketFromDatabase,
      );
    }
    const fundingIndices:
    Map<string, FundingIndexUpdatesFromDatabase> = new
    Map<string, FundingIndexUpdatesFromDatabase>();
    for (const [key, fundingIndex] of Object.entries(resultRow.funding_index_updates)) {
      fundingIndices.set(
        key,
        FundingIndexUpdatesModel.fromJson(
          fundingIndex as object,
        ) as FundingIndexUpdatesFromDatabase,
      );
    }

    const promises: Promise<number>[] = new Array<Promise<number>>(this.event.updates.length);

    for (let i: number = 0; i < this.event.updates.length; i++) {
      const update: FundingUpdateV1 = this.event.updates[i];
      if (resultRow.errors[i] != null) {
        logger.error({
          at: 'FundingHandler#handleFundingSample',
          message: resultRow.errors[i],
          update,
        });
        stats.increment(`${config.SERVICE_NAME}.handle_funding_event.failure`, 1);
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
        stats.increment(`${config.SERVICE_NAME}.handle_funding_event.failure`, 1);
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
          stats.gauge(
            `${config.SERVICE_NAME}.funding_index_update_event`,
            Number(
              protocolTranslations.fundingIndexToHumanFixedString(
                bytesToBigInt(update.fundingIndex).toString(),
                perpetualMarket,
              ),
            ),
            {
              ticker: perpetualMarket.ticker,
            },
          );
          // eslint-disable-next-line no-case-declarations
          const fundingIndexUpdate: FundingIndexUpdatesFromDatabase = fundingIndices.get(
            update.perpetualId.toString(),
          )!;
          stats.gauge(
            `${config.SERVICE_NAME}.funding_index_update`,
            Number(fundingIndexUpdate.fundingIndex),
            {
              ticker: perpetualMarket.ticker,
            },
          );
          break;
        default:
          logger.error({
            at: 'FundingHandler#handle',
            message: 'Received unknown FundingEvent type.',
            event: this.event,
          });
          stats.increment(`${config.SERVICE_NAME}.handle_funding_event.failure`, 1);
      }

      // Handle latency from resultRow
      stats.timing(
        `${config.SERVICE_NAME}.handle_funding_event.sql_latency`,
        Number(resultRow.latency),
        this.generateTimingStatsOptions(),
      );
    }

    await Promise.all(promises);
    return [];
  }
}
