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
  SubaccountMessageContents,
  SubaccountTable,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import * as pg from 'pg';

import { SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../constants';
import { addPositionsToContents, annotateWithPnl } from '../helpers/kafka-helper';
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

  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const updateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[] = _.map(
      resultRow.perpetual_positions,
      (value) => PerpetualPositionModel.fromJson(
        value) as UpdatedPerpetualPositionSubaccountKafkaObject,
    );
    const updatedAssetPositions: AssetPositionFromDatabase[] = _.map(
      resultRow.asset_positions,
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
