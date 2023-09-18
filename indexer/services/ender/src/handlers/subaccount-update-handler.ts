import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetPositionFromDatabase,
  AssetPositionModel,
  AssetPositionTable,
  assetRefresher,
  AssetsMap,
  MarketColumns,
  MarketFromDatabase,
  MarketsMap,
  MarketTable,
  perpetualMarketRefresher,
  PerpetualMarketsMap,
  PerpetualPositionColumns,
  PerpetualPositionCreateObject,
  PerpetualPositionFromDatabase,
  PerpetualPositionModel,
  PerpetualPositionsMap,
  PerpetualPositionStatus,
  PerpetualPositionSubaccountUpdateObject,
  PerpetualPositionTable,
  PositionSide,
  protocolTranslations,
  storeHelpers,
  SubaccountMessageContents,
  SubaccountTable,
  TendermintEventTable,
  UpdatedPerpetualPositionSubaccountKafkaObject,
} from '@dydxprotocol-indexer/postgres';
import { bytesToBigInt, getPositionIsLong } from '@dydxprotocol-indexer/v4-proto-parser';
import { IndexerAssetPosition, IndexerPerpetualPosition } from '@dydxprotocol-indexer/v4-protos';
import Big from 'big.js';
import _ from 'lodash';
import { DateTime } from 'luxon';
import * as pg from 'pg';

import config from '../config';
import { QUOTE_CURRENCY_ATOMIC_RESOLUTION, SUBACCOUNT_ORDER_FILL_EVENT_TYPE } from '../constants';
import { addPositionsToContents, annotateWithPnl, convertPerpetualPosition } from '../helpers/kafka-helper';
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
    if (config.USE_SUBACCOUNT_UPDATE_SQL_FUNCTION) {
      return this.handleViaSqlFunction();
    }
    return this.handleViaKnexQueries();
  }

  public async handleViaSqlFunction(): Promise<ConsolidatedKafkaEvent[]> {
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
    ).catch((error) => {
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

  public async handleViaKnexQueries(): Promise<ConsolidatedKafkaEvent[]> {
    const subaccountId: string = SubaccountTable.subaccountIdToUuid(this.event.subaccountId!);

    await this.runFuncWithTimingStatAndErrorLogging(
      this.upsertSubaccount(),
      this.generateTimingStatsOptions('upsert_subaccount'),
    );

    const perpetualMarketsMapping:
    PerpetualMarketsMap = perpetualMarketRefresher.getPerpetualMarketsMap();

    const perpetualPositionsMap: PerpetualPositionsMap = await
    this.runFuncWithTimingStatAndErrorLogging(
      this.getPerpetualPositionsMapFromEvent(subaccountId, perpetualMarketsMapping),
      this.generateTimingStatsOptions('get_existing_perpetual_positions'),
    );
    const marketIdToMarket: MarketsMap = await this.runFuncWithTimingStatAndErrorLogging(
      MarketTable.getMarketsMap(),
      this.generateTimingStatsOptions('get_markets'),
    );

    const updateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[] = await
    this.runFuncWithTimingStatAndErrorLogging(
      this.updatePerpetualPositionsFromEvent(
        subaccountId,
        perpetualMarketsMapping,
        perpetualPositionsMap,
        marketIdToMarket,
      ),
      this.generateTimingStatsOptions('update_perpetual_positions'),
    );

    const assetsMap: AssetsMap = assetRefresher.getAssetsMap();
    const updatedAssetPositions: AssetPositionFromDatabase[] = await
    this.runFuncWithTimingStatAndErrorLogging(
      this.updateAssetPositionsFromEvent(subaccountId, assetsMap),
      this.generateTimingStatsOptions('update_asset_positions'),
    );

    // TODO: Update perpetual_assets once protocol supports assets
    return [
      this.generateConsolidatedKafkaEvent(
        updateObjects,
        perpetualMarketsMapping,
        updatedAssetPositions,
        assetsMap,
      ),
    ];
  }

  /**
   * If subaccount does not exist, creates one.
   */
  protected async upsertSubaccount(): Promise<void> {
    await SubaccountTable.upsert({
      address: this.event.subaccountId!.owner,
      subaccountNumber: this.event.subaccountId!.number,
      updatedAt: this.timestamp.toISO(),
      updatedAtHeight: this.block.height.toString(),
    }, { txId: this.txId });
  }

  /**
   * Returns a list of asset ids that are missing from assets
   *
   * @param assetIds
   * @param assets
   * @protected
   */
  protected findMissingAssets(assetIds: string[], assets: AssetFromDatabase[]): string[] {
    const presentAssets: string[] = assets.map((asset) => {
      return asset.id;
    });
    return assetIds.filter((element) => !presentAssets.includes(element));
  }

  /**
   * Updates all asset positions in postgres for each 'updatedAssetPosition'
   * in the SubaccountUpdateEvent.
   * @param subaccountId
   * @param AssetsMap
   * @protected
   */
  protected async updateAssetPositionsFromEvent(
    subaccountId: string,
    assetsMap: AssetsMap,
  ): Promise<AssetPositionFromDatabase[]> {
    const assetPositions: AssetPositionFromDatabase[] = await
    Promise.all(
      _.map(
        this.event.updatedAssetPositions,
        async (assetPositionProto: IndexerAssetPosition) => {
          return this.upsertAssetPositionFromAssetPositionProto(
            subaccountId,
            assetPositionProto,
            assetsMap[assetPositionProto.assetId.toString()],
          );
        },
      ),
    );
    return assetPositions;
  }

  protected async getPerpetualPositionsMapFromEvent(
    subaccountId: string,
    perpetualMarketsMapping: PerpetualMarketsMap,
  ): Promise<PerpetualPositionsMap> {
    const perpetualPositions = await PerpetualPositionTable.findAll({
      subaccountId: [subaccountId],
      perpetualId: _.map(
        this.event.updatedPerpetualPositions,
        (perpetualPositionProto: IndexerPerpetualPosition) => {
          return perpetualMarketsMapping[perpetualPositionProto.perpetualId].id;
        },
      ),
      status: [PerpetualPositionStatus.OPEN],
    }, [], { txId: this.txId });
    return _.keyBy(perpetualPositions, PerpetualPositionColumns.perpetualId);
  }

  /**
   * Updates all perpetual positions in postgres for each 'updatedPerpetualPosition'
   * in the SubaccountUpdateEvent.
   * @returns a list of PerpetualPositionSubaccountUpdateKafkaObject
   */
  protected async updatePerpetualPositionsFromEvent(
    subaccountId: string,
    perpetualMarketsMapping: PerpetualMarketsMap,
    perpetualPositionsMap: PerpetualPositionsMap,
    marketsMap: MarketsMap,
  ): Promise<UpdatedPerpetualPositionSubaccountKafkaObject[]> {
    const positionUpdateObjects: UpdatedPerpetualPositionSubaccountKafkaObject[] = [];
    const positionCreateObjects: PerpetualPositionCreateObject[] = [];

    _.forEach(
      this.event.updatedPerpetualPositions,
      (perpetualPositionProto: IndexerPerpetualPosition) => {
        const [
          updateObject,
          createObject,
        ]: [
          UpdatedPerpetualPositionSubaccountKafkaObject | null,
          PerpetualPositionCreateObject | null,
        ] = this.generateUpdateAndCreateFromPerpetualPositionProto(
          subaccountId,
          perpetualPositionProto,
          perpetualMarketsMapping,
          perpetualPositionsMap[perpetualPositionProto.perpetualId],
          marketsMap,
        );

        if (updateObject !== null) {
          positionUpdateObjects.push(updateObject);
        }
        if (createObject !== null) {
          positionCreateObjects.push(createObject);
        }
      },
    );

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const [createdPositions, _ignore]: [PerpetualPositionFromDatabase[], void] = await Promise.all([
      PerpetualPositionTable.bulkCreate(positionCreateObjects, { txId: this.txId }),
      PerpetualPositionTable.bulkUpdateSubaccountFields(
        _.map(
          positionUpdateObjects,
          this.getPerpetualPositionSubaccountUpdateObject,
        ),
        { txId: this.txId },
      ),
    ]);

    const createdPositionsWithPnl:
    UpdatedPerpetualPositionSubaccountKafkaObject[] = createdPositions
      .map(
        (position) => {
          return annotateWithPnl(
            convertPerpetualPosition(position),
            perpetualMarketRefresher.getPerpetualMarketsMap(),
            marketsMap,
          );
        });
    // We can combine the two arrays because the PerpetualPositionFromDatabase extends
    // UpdatedPerpetualPositionSubaccountKafkaObject.
    return _.flatten([positionUpdateObjects, createdPositionsWithPnl]);
  }

  /**
   * Generates a PerpetualPositionSubaccountUpdateObject from
   * PerpetualPositionSubaccountUpdateKafkaObject by picking the relevant fields.
   */
  protected getPerpetualPositionSubaccountUpdateObject(
    kafkaObject: UpdatedPerpetualPositionSubaccountKafkaObject,
  ): PerpetualPositionSubaccountUpdateObject {
    return _.pick(kafkaObject, [
      PerpetualPositionColumns.id,
      PerpetualPositionColumns.closedAt,
      PerpetualPositionColumns.closedAtHeight,
      PerpetualPositionColumns.closeEventId,
      PerpetualPositionColumns.lastEventId,
      PerpetualPositionColumns.settledFunding,
      PerpetualPositionColumns.status,
      PerpetualPositionColumns.size,
    ]);
  }

  /**
   * Makes postgres updates for the asset position based on the 'assetPositionProto'.
   * If there is an existing position, update the existing asset position.
   * Else create a new asset position.
   * @param subaccountId
   * @param perpetualPositionProto
   * @param existingPosition
   * @param perpetualMarket
   * @returns
   */
  // eslint-disable-next-line @typescript-eslint/require-await
  protected async upsertAssetPositionFromAssetPositionProto(
    subaccountId: string,
    assetPositionProto: IndexerAssetPosition,
    assetFromDatabase: AssetFromDatabase,
  ): Promise<AssetPositionFromDatabase> {
    return AssetPositionTable.upsert({
      subaccountId,
      assetId: assetFromDatabase.id,
      // TODO(DEC-1597): deprecate `isLong` in asset and perpetual position tables.
      isLong: getPositionIsLong(assetPositionProto),
      // TODO(DEC-1597): use signed instead of absolute value after deprecating `isLong`.
      size: protocolTranslations.serializedQuantumsToAbsHumanFixedString(
        assetPositionProto.quantums,
        assetFromDatabase.atomicResolution,
      ),
    }, { txId: this.txId });
  }

  /**
   * Returns the SubaccountUpdate and Create objects for PerpetualPositions based on the
   * 'perpetualPositionProto'.
   * If there is no existing position, create the perpetual position.
   * If the updated position has size 0, close the existing position.
   * If the updated position has the same side as the existing position,
   * update the existing position.
   * If the updated position has the opposite side as the existing position,
   * close the existing position, and create a new perpetual position.
   */
  protected generateUpdateAndCreateFromPerpetualPositionProto(
    subaccountId: string,
    perpetualPositionProto: IndexerPerpetualPosition,
    perpetualMarketMap: PerpetualMarketsMap,
    existingPosition: PerpetualPositionFromDatabase | undefined,
    marketIdToMarket: MarketsMap,
  ): [
      UpdatedPerpetualPositionSubaccountKafkaObject | null,
      PerpetualPositionCreateObject | null,
    ] {
    let updateObject: UpdatedPerpetualPositionSubaccountKafkaObject | null = null;
    const size: string = protocolTranslations.serializedQuantumsToAbsHumanFixedString(
      perpetualPositionProto.quantums,
      perpetualMarketMap[perpetualPositionProto.perpetualId].atomicResolution,
    );
    const side: PositionSide = getPositionIsLong(perpetualPositionProto)
      ? PositionSide.LONG
      : PositionSide.SHORT;
    const eventId: Buffer = TendermintEventTable.createEventId(
      this.block.height.toString(),
      indexerTendermintEventToTransactionIndex(this.indexerTendermintEvent),
      this.indexerTendermintEvent.eventIndex,
    );
    const blockTime: string = DateTime.fromJSDate(this.block.time!).toISO();
    const latestFundingQuantums: string = bytesToBigInt(
      perpetualPositionProto.fundingPayment,
    ).toString();
    const latestSettledFunding: Big = protocolTranslations.quantumsToHuman(
      latestFundingQuantums,
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).times(-1);
    let priorSettledFunding: Big = new Big(0);
    if (existingPosition !== undefined) {
      priorSettledFunding = new Big(existingPosition.settledFunding);
    }
    const settledFunding: string = priorSettledFunding.plus(latestSettledFunding).toString();

    // Close existing position and do not create another if incoming size is 0.
    if (existingPosition !== undefined && size === '0') {
      return [
        annotateWithPnl(
          {
            ...PerpetualPositionTable.closePositionUpdateObject(
              existingPosition,
              {
                id: existingPosition.id,
                closedAt: blockTime,
                closedAtHeight: this.block.height.toString(),
                closeEventId: eventId,
                settledFunding,
              },
            ),
            perpetualId: perpetualPositionProto.perpetualId.toString(),
            maxSize: existingPosition.maxSize,
            side: existingPosition.side,
            entryPrice: existingPosition.entryPrice,
            exitPrice: existingPosition.exitPrice,
            sumOpen: existingPosition.sumOpen,
            sumClose: existingPosition.sumClose,
          },
          perpetualMarketMap,
          marketIdToMarket,
        ),
        null,
      ];
    }

    if (existingPosition !== undefined) {
      if (existingPosition.side === side) {
        return [
          annotateWithPnl(
            {
              id: existingPosition.id,
              size,
              status: PerpetualPositionStatus.OPEN,
              lastEventId: eventId,
              settledFunding,
              perpetualId: existingPosition.perpetualId,
              maxSize: Big(existingPosition.maxSize).gte(size) ? existingPosition.maxSize : size,
              side: existingPosition.side,
              entryPrice: existingPosition.entryPrice,
              exitPrice: existingPosition.exitPrice,
              sumOpen: existingPosition.sumOpen,
              sumClose: existingPosition.sumClose,
            },
            perpetualMarketMap,
            marketIdToMarket,
          ),
          null,
        ];
      } else {
        // Close the existing position if the existing position is of the opposite side of the
        // new position. New position will be created below.
        updateObject = annotateWithPnl(
          {
            ...PerpetualPositionTable.closePositionUpdateObject(
              existingPosition,
              {
                id: existingPosition.id,
                closedAt: blockTime,
                closedAtHeight: this.block.height.toString(),
                closeEventId: eventId,
                settledFunding,
              },
            ),
            perpetualId: existingPosition.perpetualId,
            maxSize: existingPosition.maxSize,
            side: existingPosition.side,
            entryPrice: existingPosition.entryPrice,
            exitPrice: existingPosition.exitPrice,
            sumOpen: existingPosition.sumOpen,
            sumClose: existingPosition.sumClose,
          },
          perpetualMarketMap,
          marketIdToMarket,
        );
      }
    }

    // should create a new perpetual position if none exist or if previous position has changed side
    // and is not 0.
    // if the previous position changed sides, the last funding payment will be applied to the
    // settled funding of the closed position and the new position will be created with 0 settled
    // funding.
    return [
      updateObject,
      {
        subaccountId,
        perpetualId: perpetualPositionProto.perpetualId.toString(),
        side,
        status: PerpetualPositionStatus.OPEN,
        size,
        maxSize: size,
        createdAt: blockTime,
        createdAtHeight: this.block.height.toString(),
        openEventId: eventId,
        lastEventId: eventId,
        settledFunding: updateObject === null ? settledFunding : '0',
      },
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
