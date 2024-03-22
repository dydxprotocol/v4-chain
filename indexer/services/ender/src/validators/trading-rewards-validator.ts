import { bytesToBigInt } from '@dydxprotocol-indexer/v4-proto-parser';
import { IndexerTendermintEvent, TradingRewardsEventV1, AddressTradingReward } from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';

import { Handler } from '../handlers/handler';
import { TradingRewardsHandler } from '../handlers/trading-rewards-handler';
import { Validator } from './validator';

export class TradingRewardsValidator extends Validator<TradingRewardsEventV1> {
  public validate(): void {
    _.forEach(this.event.tradingRewards, (reward: AddressTradingReward, index: number) => {
      this.validateTradingReward(reward, index);
    });
  }

  private validateTradingReward(
    reward: AddressTradingReward,
    index: number,
  ): void {
    const denoms: bigint = bytesToBigInt(reward.denomAmount);
    if (denoms === BigInt(0)) {
      return this.logAndThrowParseMessageError(
        `TradingReward in TradingRewardEvent at index ${index} is missing denoms.`,
        { event: this.event, denoms },
      );
    }

    if (reward.owner === undefined) {
      return this.logAndThrowParseMessageError(
        `TradingReward in TradingRewardEvent at index ${index} is missing an owner.`,
        { event: this.event, address: reward.owner },
      );
    }
  }

  public createHandlers(
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    __: string,
  ): Handler<TradingRewardsEventV1>[] {
    return [
      new TradingRewardsHandler(
        this.block,
        this.blockEventIndex,
        indexerTendermintEvent,
        txId,
        this.event,
      ),
    ];
  }
}
