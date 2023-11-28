import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetPositionFromDatabase,
  AssetPositionModel,
  assetRefresher,
  AssetsMap,
  MarketColumns,
  MarketFromDatabase,
  MarketsMap,
  MarketTable,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PerpetualPositionModel,
  storeHelpers,
  SubaccountMessageContents,
  SubaccountTable,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import * as pg from 'pg';

import { SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../constants';
import { addPositionsToContents, annotateWithPnl } from '../helpers/kafka-helper';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { SubaccountUpdate } from '../lib/translated-types';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class SubaccountUpdateHandler extends Handler<SubaccountUpdate> {
  eventType: string = 'SubaccountUpdateEvent';

  public getParallelizationIds(): string[] {
    // SubaccountUpdateEvents with the same subaccountId must be handled sequentially
    return [
      `${this.eventType}_${SubaccountTable.subaccountIdToUuid(this.event.subaccountId!)}`,
      // To ensure that SubaccountUpdateEvents and OrderFillEvents for the same subaccount are not
      // processed in parallel
      `${SUBACCOUNT_ORDER_FILL_EVENT_TYPE}_${SubaccountTable.subaccountIdToUuid(this.event.subaccountId!)}`,
    ];
  }

  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );

    const result: pg.QueryResult = await storeHelpers.rawQuery(`SELECT dydx_subaccount_update_handler(
      ${this.block.height}, 
      '${this.block.time?.toISOString()}', 
      '${JSON.stringify(this.event)}', 
      ${this.indexerTendermintEvent.eventIndex}, 
      ${transactionIndex}) AS result;`,
    { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'subaccountUpdateHandler#handleViaSqlFunction',
        message: 'Failed to handle SubaccountUpdateEventV1',
        error,
      });
      throw error;
    });
    const updateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[] = _.map(
      result.rows[0].result.perpetual_positions,
      (value) => PerpetualPositionModel.fromJson(
        value) as UpdatedPerpetualPositionSubaccountKafkaObject,
    );
    const updatedAssetPositions: AssetPositionFromDatabase[] = _.map(
      result.rows[0].result.asset_positions,
      (value) => AssetPositionModel.fromJson(value) as AssetPositionFromDatabase,
    );
    const markets: MarketFromDatabase[] = await MarketTable.findAll(
      {},
      [],
      { txId: this.txId },
    );
    const marketIdToMarket: MarketsMap = _.keyBy(
      markets,
      MarketColumns.id,
    );
    for (let i = 0; i < updateObjects.length; i++) {
      updateObjects[i] = annotateWithPnl(
        updateObjects[i],
        perpetualMarketRefresher.getPerpetualMarketsMap(),
        marketIdToMarket,
      );
    }

    return [
      this.generateConsolidatedKafkaEvent(
        updateObjects,
        perpetualMarketRefresher.getPerpetualMarketsMap(),
        updatedAssetPositions,
        assetRefresher.getAssetsMap(),
      ),
    ];
  }

  /**
   * Generate the ConsolidatedKafkaEvent generated from this event.
   * @param updatedPerpetualPositions
   * @returns
   */
  protected generateConsolidatedKafkaEvent(
    updateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[],
    perpetualMarketsMapping: PerpetualMarketsMap,
    updatedAssetPositions: AssetPositionFromDatabase[],
    assetsMap: AssetsMap,
  ): ConsolidatedKafkaEvent {
    const contents: SubaccountMessageContents = addPositionsToContents(
      {},
      this.event.subaccountId!,
      updateObjects,
      perpetualMarketsMapping,
      updatedAssetPositions,
      assetsMap,
    );

    return this.generateConsolidatedSubaccountKafkaEvent(
      JSON.stringify(contents),
      this.event.subaccountId!,
    );
  }
}
