import { stats } from '@dydxprotocol-indexer/base';
import {
  SubaccountMessageContents,
  TradingRewardFromDatabase,
  TradingRewardModel,
  TradingRewardSubaccountMessageContents,
} from '@dydxprotocol-indexer/postgres';
import { TradingRewardsEventV1 } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class TradingRewardsHandler extends Handler<TradingRewardsEventV1> {
  eventType: string = 'TradingRewardEvent';

  public getParallelizationIds(): string[] {
    // Can be handled in any order for trading reward events
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_trading_rewards_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    const tradingRewards: TradingRewardFromDatabase[] = _.map(
      resultRow.trading_rewards,
      (tradingReward: object) => {
        return TradingRewardModel.fromJson(tradingReward) as TradingRewardFromDatabase;
      },
    );
    return this.generateKafkaEvents(
      tradingRewards,
    );
  }

  /** Generates a kafka websocket event for each address that receives a trading reward.
   *
   * @param tradingRewards
   * @protected
   */
  protected generateKafkaEvents(
    tradingRewards: TradingRewardFromDatabase[],
  ): ConsolidatedKafkaEvent[] {
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    _.forEach(tradingRewards, (tradingReward: TradingRewardFromDatabase) => {
      const tradingRewardSubaccountMessageContents: TradingRewardSubaccountMessageContents = {
        tradingReward: tradingReward.amount,
        createdAtHeight: tradingReward.blockHeight,
        createdAt: tradingReward.blockTime,
      };

      const subaccountMessageContents: SubaccountMessageContents = {
        tradingReward: tradingRewardSubaccountMessageContents,
        blockHeight: this.block.height.toString(),
      };

      kafkaEvents.push(
        this.generateConsolidatedSubaccountKafkaEvent(
          JSON.stringify(subaccountMessageContents),
          {
            owner: tradingReward.address,
            number: 0,
          },
        ),
      );
    });

    return kafkaEvents;
  }
}
