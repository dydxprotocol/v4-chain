import { logger } from '@dydxprotocol-indexer/base';
import {
  FundingIndexUpdatesTable,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  TendermintEventTable,
  FundingIndexUpdatesCreateObject,
  FundingIndexUpdatesFromDatabase,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { NextFundingCache } from '@dydxprotocol-indexer/redis';
import { bytesToBigInt } from '@dydxprotocol-indexer/v4-proto-parser';
import { FundingEventV1_Type, FundingUpdateV1 } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';

import { getPrice } from '../caches/price-cache';
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
    logger.info({
      at: 'FundingHandler#handle',
      message: 'Received FundingEvent.',
      event: this.event,
    });
    const castedFundingEvent: FundingEventMessage = this.event as FundingEventMessage;
    switch (castedFundingEvent.type) {
      case FundingEventV1_Type.TYPE_PREMIUM_SAMPLE:
        await this.runFuncWithTimingStatAndErrorLogging(
          this.handleFundingSample(castedFundingEvent.updates),
          this.generateTimingStatsOptions('handle_premium_sample'),
        );
        break;
      case FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX:
        await this.runFuncWithTimingStatAndErrorLogging(
          this.handleFundingRate(castedFundingEvent.updates),
          this.generateTimingStatsOptions('handle_funding_rate'),
        );
        break;
      default:
        logger.error({
          at: 'FundingHandler#handle',
          message: 'Received unknown FundingEvent type.',
          event: this.event,
        });
    }
    return [];
  }

  public async handleFundingSample(samples: FundingUpdateV1[]): Promise<void> {
    await Promise.all(
      _.map(samples, (sample: FundingUpdateV1) => {
        const perpetualMarket:
        PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher.getPerpetualMarketFromId(
          sample.perpetualId.toString(),
        );
        if (perpetualMarket === undefined) {
          logger.error({
            at: 'FundingHandler#handleFundingSample',
            message: 'Received FundingUpdate with unknown perpetualId.',
            sample,
          });
          return;
        }
        const ticker: string = perpetualMarket.ticker;
        const rate: string = protocolTranslations.funding8HourValuePpmTo1HourRate(
          sample.fundingValuePpm,
        );
        return NextFundingCache.addFundingSample(ticker, new Big(rate), redisClient);
      }),
    );
  }

  public async handleFundingRate(updates: FundingUpdateV1[]): Promise<void> {
    // clear the cache for the predicted next funding rate
    await Promise.all(
      _.map(updates, (update: FundingUpdateV1) => {
        const perpetualMarket:
        PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher.getPerpetualMarketFromId(
          update.perpetualId.toString(),
        );
        if (perpetualMarket === undefined) {
          logger.error({
            at: 'FundingHandler#handleFundingRate',
            message: 'Received FundingUpdate with unknown perpetualId.',
            update,
          });
          return;
        }
        const ticker: string = perpetualMarket.ticker;
        const numCleared:
        Promise<number> = NextFundingCache.clearFundingSamples(ticker, redisClient);
        const fundingIndexUpdatesCreateObject: FundingIndexUpdatesCreateObject = {
          perpetualId: update.perpetualId.toString(),
          eventId: this.eventId,
          rate: protocolTranslations.funding8HourValuePpmTo1HourRate(update.fundingValuePpm),
          oraclePrice: getPrice(perpetualMarket.marketId),
          fundingIndex: protocolTranslations.fundingIndexToHumanFixedString(
            bytesToBigInt(update.fundingIndex).toString(),
            perpetualMarket,
          ),
          effectiveAt: this.timestamp.toISO(),
          effectiveAtHeight: this.block.height.toString(),
        };
        const fundingIndexUpdatesFromDatabase:
        Promise<FundingIndexUpdatesFromDatabase> = FundingIndexUpdatesTable
          .create(
            fundingIndexUpdatesCreateObject,
            { txId: this.txId },
          );
        return [numCleared, fundingIndexUpdatesFromDatabase];
      })
        // flatten nested promise arrays
        .map(Promise.all, Promise),
    );
  }
}
