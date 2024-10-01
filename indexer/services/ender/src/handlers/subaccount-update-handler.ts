import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetPositionFromDatabase,
  AssetPositionModel,
  assetRefresher,
  AssetsMap,
  MarketFromDatabase,
  MarketModel,
  MarketsMap,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PerpetualPositionModel,
  SubaccountMessageContents,
  SubaccountTable,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import * as pg from 'pg';

import config from '../config';
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

  // eslint-disable-next-line @typescript-eslint/require-await
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
    const marketIdToMarket: MarketsMap = _.mapValues(
      resultRow.markets,
      (value) => MarketModel.fromJson(value) as MarketFromDatabase,
    );

    for (let i = 0; i < updateObjects.length; i++) {
      const marketId: number = perpetualMarketRefresher.getPerpetualMarketsMap()[
        updateObjects[i].perpetualId
      ].marketId;
      updateObjects[i] = annotateWithPnl(
        updateObjects[i],
        perpetualMarketRefresher.getPerpetualMarketsMap(),
        marketIdToMarket[marketId],
      );
    }
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_subaccount_update_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );

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
      this.block.height.toString(),
    );

    return this.generateConsolidatedSubaccountKafkaEvent(
      JSON.stringify(contents),
      this.event.subaccountId!,
    );
  }
}
