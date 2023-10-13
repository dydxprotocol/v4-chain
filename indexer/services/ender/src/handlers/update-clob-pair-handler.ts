import assert from 'assert';

import {
  PerpetualMarketFromDatabase, PerpetualMarketTable, perpetualMarketRefresher, protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { UpdateClobPairEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { generatePerpetualMarketMessage } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class UpdateClobPairHandler extends Handler<UpdateClobPairEventV1> {
  eventType: string = 'UpdateClobPairEventV1';

  public getParallelizationIds(): string[] {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const perpetualMarket:
    PerpetualMarketFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.updateClobPair(),
      this.generateTimingStatsOptions('update_clob_pair'),
    );
    return [
      this.generateConsolidatedMarketKafkaEvent(
        JSON.stringify(generatePerpetualMarketMessage([perpetualMarket])),
      ),
    ];
  }

  private async updateClobPair(): Promise<PerpetualMarketFromDatabase> {
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
      stepBaseQuantums: Number(this.event.stepBaseQuantums),
    }, { txId: this.txId });

    if (perpetualMarket === undefined) {
      this.logAndThrowParseMessageError(
        'Could not find perpetual market with corresponding updatePerpetualEvent.id',
        { event: this.event },
      );
      // This assert should never be hit because a ParseMessageError should be thrown above.
      assert(false);
    }

    await perpetualMarketRefresher.upsertPerpetualMarket(perpetualMarket);
    return perpetualMarket;
  }
}
