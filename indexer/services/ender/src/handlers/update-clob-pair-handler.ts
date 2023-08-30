import {
  PerpetualMarketFromDatabase, PerpetualMarketTable, perpetualMarketRefresher, protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { UpdateClobPairEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdateClobPairHandler extends Handler<UpdateClobPairEventV1> {
  eventType: string = 'UpdateClobPairEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      this.updateClobPair(),
      this.generateTimingStatsOptions('update_clob_pair'),
    );
    // TODO(IND-374): Send update to markets websocket channel.
    return [];
  }

  private async updateClobPair(): Promise<void> {
    // perpetualMarketRefresher.getPerpetualMarketFromClobPairId() cannot be undefined because it
    // is validated by UpdateClobPairValidator.
    const perpetualMarketId: string = perpetualMarketRefresher.getPerpetualMarketFromClobPairId(
      this.event.clobPairId.toString(),
    )!.id;
    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.update({
      id: perpetualMarketId,
      status: protocolTranslations.clobStatusToMarketStatus(this.event.status),
      quantumConversionExponent: this.event.quantumConversionExponent,
      subticksPerTick: this.event.subticksPerTick,
      minOrderBaseQuantums: Number(this.event.minOrderBaseQuantums),
      stepBaseQuantums: Number(this.event.stepBaseQuantums),
    }, { txId: this.txId });

    if (perpetualMarket === undefined) {
      return this.logAndThrowParseMessageError(
        'Could not find perpetual market with corresponding updatePerpetualEvent.id',
        { event: this.event },
      );
    }
    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
  }
}
