/* eslint-disable max-len */
import { logger, stats, STATS_NO_SAMPLING } from '@dydxprotocol-indexer/base';
import { BLOCK_HEIGHT_WEBSOCKET_MESSAGE_VERSION, KafkaTopics } from '@dydxprotocol-indexer/kafka';
import {
  storeHelpers,
} from '@dydxprotocol-indexer/postgres';
import {
  BlockHeightMessage,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
} from '@dydxprotocol-indexer/v4-protos';
import _ from 'lodash';
import * as pg from 'pg';
import { DatabaseError } from 'pg';

import config from '../config';
import { Handler } from '../handlers/handler';
import { AssetValidator } from '../validators/asset-validator';
import { DeleveragingValidator } from '../validators/deleveraging-validator';
import { FundingValidator } from '../validators/funding-validator';
import { LiquidityTierValidatorV2, LiquidityTierValidator } from '../validators/liquidity-tier-validator';
import { MarketValidator } from '../validators/market-validator';
import { OrderFillValidator } from '../validators/order-fill-validator';
import { PerpetualMarketValidator } from '../validators/perpetual-market-validator';
import { RegisterAffiliateValidator } from '../validators/register-affiliate-validator';
import { StatefulOrderValidator } from '../validators/stateful-order-validator';
import { SubaccountUpdateValidator } from '../validators/subaccount-update-validator';
import { TradingRewardsValidator } from '../validators/trading-rewards-validator';
import { TransferValidator } from '../validators/transfer-validator';
import { UpdateClobPairValidator } from '../validators/update-clob-pair-validator';
import { UpdatePerpetualValidator } from '../validators/update-perpetual-validator';
import { UpsertVaultValidator } from '../validators/upsert-vault-validator';
import { Validator, ValidatorInitializer } from '../validators/validator';
import { BatchedHandlers } from './batched-handlers';
import { indexerTendermintEventToEventProtoWithType, indexerTendermintEventToTransactionIndex } from './helper';
import { KafkaPublisher } from './kafka-publisher';
import { SyncHandlers, SYNCHRONOUS_SUBTYPES } from './sync-handlers';
import {
  ConsolidatedKafkaEvent,
  DydxIndexerSubtypes, EventMessage, EventProtoWithTypeAndVersion, GroupedEvents, SKIPPED_EVENT_SUBTYPE,
} from './types';

const TXN_EVENT_SUBTYPE_VERSION_TO_VALIDATOR_MAPPING: Record<string, ValidatorInitializer> = {
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.ORDER_FILL.toString(), 1)]: OrderFillValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.SUBACCOUNT_UPDATE.toString(), 1)]: SubaccountUpdateValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.TRANSFER.toString(), 1)]: TransferValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.MARKET.toString(), 1)]: MarketValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.STATEFUL_ORDER.toString(), 1)]: StatefulOrderValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.ASSET.toString(), 1)]: AssetValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.PERPETUAL_MARKET.toString(), 1)]: PerpetualMarketValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.PERPETUAL_MARKET.toString(), 2)]: PerpetualMarketValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.PERPETUAL_MARKET.toString(), 3)]: PerpetualMarketValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.LIQUIDITY_TIER.toString(), 1)]: LiquidityTierValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPDATE_PERPETUAL.toString(), 1)]: UpdatePerpetualValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPDATE_PERPETUAL.toString(), 2)]: UpdatePerpetualValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPDATE_PERPETUAL.toString(), 3)]: UpdatePerpetualValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPDATE_CLOB_PAIR.toString(), 1)]: UpdateClobPairValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.DELEVERAGING.toString(), 1)]: DeleveragingValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.LIQUIDITY_TIER.toString(), 2)]: LiquidityTierValidatorV2,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.REGISTER_AFFILIATE.toString(), 1)]: RegisterAffiliateValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPSERT_VAULT.toString(), 1)]: UpsertVaultValidator,
};

const BLOCK_EVENT_SUBTYPE_VERSION_TO_VALIDATOR_MAPPING: Record<string, ValidatorInitializer> = {
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.FUNDING.toString(), 1)]: FundingValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.TRADING_REWARD.toString(), 1)]: TradingRewardsValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.STATEFUL_ORDER.toString(), 1)]: StatefulOrderValidator,
  [serializeSubtypeAndVersion(DydxIndexerSubtypes.UPSERT_VAULT.toString(), 1)]: UpsertVaultValidator,
};

function serializeSubtypeAndVersion(
  subtype: string,
  version: number,
): string {
  return `${subtype}-${version}`;
}

type DecodedIndexerTendermintBlock = Omit<IndexerTendermintBlock, 'events'> & {
  events: DecodedIndexerTendermintEvent[],
};

type DecodedIndexerTendermintEvent = Omit<IndexerTendermintEvent, 'dataBytes'> & {
  /** Decoded tendermint event. */
  dataBytes: object,
};

export class BlockProcessor {
  block: IndexerTendermintBlock;
  sqlEventPromises: Promise<object>[];
  sqlBlock: DecodedIndexerTendermintBlock;
  txId: number;
  batchedHandlers: BatchedHandlers;
  syncHandlers: SyncHandlers;
  messageReceivedTimestamp: string;

  constructor(
    block: IndexerTendermintBlock,
    txId: number,
    messageReceivedTimestamp: string,
  ) {
    this.block = block;
    this.txId = txId;
    this.messageReceivedTimestamp = messageReceivedTimestamp;
    this.sqlBlock = {
      ...this.block,
      events: new Array(this.block.events.length),
    };
    this.sqlEventPromises = new Array(this.block.events.length);
    this.batchedHandlers = new BatchedHandlers();
    this.syncHandlers = new SyncHandlers();
  }

  /**
   * Saves all Tendermint event data to postgres. Performs these operations in the following steps:
   * 1. Group - Groups events by transaction events and block events by transactionIndex and
   *            eventIndex.
   * 2. Validation - Validates that all data from v4 has the required fields with
   *                 Handler.validate(). Validation failure will throw a ParseMessageError.
   * 3. Organize - Groups events into events that can be processed in parallel. Each handler
   *               will generate a list of ids that if matching another handler's ids cannot
   *               be processed in parallel. For example, two SubaccountUpdateEvents for the
   *               same subaccount cannot be processed in parallel because the subaccount's
   *               final balance should be the last SubaccountUpdateEvent's balance.
   * 4. Processing - Based on the groupings created above, process events in each batch in parallel.
   * @returns the kafka publisher which contains all the events to be published to the kafka
   */
  public async process(): Promise<KafkaPublisher> {
    const groupedEvents: GroupedEvents = this.groupEvents();
    this.validateAndOrganizeEvents(groupedEvents);
    return this.processEvents();
  }

  /**
   * Groups events into block events and events for each transactionIndex
   * @param block the IndexerTendermintBlock to group events
   * @returns
   */
  private groupEvents(): GroupedEvents {
    const groupedEvents: GroupedEvents = {
      transactionEvents: [],
      blockEvents: [],
    };

    for (let i: number = 0; i < this.block.txHashes.length; i++) {
      groupedEvents.transactionEvents.push([]);
    }

    for (let i: number = 0; i < this.block.events.length; i++) {
      const event: IndexerTendermintEvent = this.block.events[i];
      const transactionIndex: number = indexerTendermintEventToTransactionIndex(event);
      const eventProtoWithType:
      EventProtoWithTypeAndVersion | undefined = indexerTendermintEventToEventProtoWithType(
        i,
        event,
      );
      if (eventProtoWithType === undefined) {
        continue;
      }
      if (transactionIndex === -1) {
        groupedEvents.blockEvents.push(eventProtoWithType);
        continue;
      }

      groupedEvents.transactionEvents[transactionIndex].push(eventProtoWithType);
    }
    return groupedEvents;
  }

  /**
   * Organizes all events into batches that can be processed in parallel, and validates that
   * the blocks are valid. Any invalid block will throw a ParseMessageError, and will be handled
   * in onMessage.
   * Each event is validated by a validator and also returns a list of handlers that will process
   * the event.
   */
  private validateAndOrganizeEvents(groupedEvents: GroupedEvents): void {
    for (const eventsInTransaction of groupedEvents.transactionEvents) {
      for (const eventProtoWithType of eventsInTransaction) {
        this.validateAndAddHandlerForEvent(
          eventProtoWithType,
          TXN_EVENT_SUBTYPE_VERSION_TO_VALIDATOR_MAPPING,
        );
      }
    }
    for (const eventProtoWithType of groupedEvents.blockEvents) {
      this.validateAndAddHandlerForEvent(
        eventProtoWithType,
        BLOCK_EVENT_SUBTYPE_VERSION_TO_VALIDATOR_MAPPING,
      );
    }
  }

  private validateAndAddHandlerForEvent(
    eventProto: EventProtoWithTypeAndVersion,
    validatorMap: Record<string, ValidatorInitializer>,
  ): void {
    const Initializer:
    ValidatorInitializer | undefined = validatorMap[
      serializeSubtypeAndVersion(
        eventProto.type,
        eventProto.version,
      )
    ];

    if (Initializer === undefined) {
      const message: string = `cannot process subtype ${eventProto.type} and version ${eventProto.version}`;
      logger.error({
        at: 'onMessage#saveTendermintEventData',
        message,
        eventProto,
      });
      return;
    }

    const validator: Validator<EventMessage> = new Initializer(
      eventProto.eventProto,
      this.block,
      eventProto.blockEventIndex,
    );
    validator.validate();
    this.sqlEventPromises[eventProto.blockEventIndex] = validator.getEventForBlockProcessor();
    let handlers: Handler<EventMessage>[];

    if (validator.shouldSkipSql()) {
      // If sql processing should be skipped, set event subtype to `skipped_event`.
      this.block.events[eventProto.blockEventIndex] = {
        ...this.block.events[eventProto.blockEventIndex],
        subtype: SKIPPED_EVENT_SUBTYPE,
      };
      logger.info({
        at: 'onMessage#shouldSkipSql',
        message: 'Excluded event from sql processing',
        eventProto,
      });
    }
    if (validator.shouldSkipHandlers()) {
      // If handlers should be skipped, set handlers to empty.
      handlers = [];
      logger.info({
        at: 'onMessage#shouldSkipHandlers',
        message: 'Excluded event from handler processing',
        eventProto,
      });
    } else {
      handlers = validator.createHandlers(
        eventProto.indexerTendermintEvent,
        this.txId,
        this.messageReceivedTimestamp,
      );
    }

    _.map(handlers, (handler: Handler<EventMessage>) => {
      if (SYNCHRONOUS_SUBTYPES.includes(eventProto.type as DydxIndexerSubtypes)) {
        this.syncHandlers.addHandler(eventProto.type, handler);
      } else {
        this.batchedHandlers.addHandler(handler);
      }
    });
  }

  createBlockHeightMsg(): ConsolidatedKafkaEvent {
    const message: BlockHeightMessage = {
      blockHeight: String(this.block.height),
      version: BLOCK_HEIGHT_WEBSOCKET_MESSAGE_VERSION,
      time: this.block.time?.toISOString() ?? '',
    };
    return {
      topic: KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT,
      message,
    };
  }

  private async processEvents(): Promise<KafkaPublisher> {
    const kafkaPublisher: KafkaPublisher = new KafkaPublisher();

    await Promise.all(this.sqlEventPromises).then((values) => {
      for (let i: number = 0; i < this.block.events.length; i++) {
        const event: IndexerTendermintEvent = this.block.events[i];

        this.sqlBlock.events[i] = {
          ...event,
          // For skipped events, use an empty object since SQL block processor ignores them.
          // Otherwise, use the decoded event since SQL block processor doesn't know how to decode
          // protobuf.
          dataBytes: event.subtype === SKIPPED_EVENT_SUBTYPE ? {} : values[i],
        };
      }
    });

    const start: number = Date.now();
    let success = false;
    let resultRow: pg.QueryResultRow;
    try {
      const result: pg.QueryResult = await storeHelpers.rawQuery(
        'SELECT dydx_block_processor(?) AS result;',
        {
          txId: this.txId,
          bindings: [JSON.stringify(this.sqlBlock)],
          sqlOptions: { name: 'dydx_block_processor' },
        },
      ).catch((error: DatabaseError) => {
        logger.crit({
          at: `BlockProcessor#processEvents\n${error.where}`,
          message: error.message,
          error,
        });
        throw error;
      });
      resultRow = result.rows[0].result;
      success = true;
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.processed_block_sql.timing`,
        Date.now() - start,
        STATS_NO_SAMPLING,
        { success: success.toString() },
      );
    }

    // Create a block message from the current block
    kafkaPublisher.addEvent(this.createBlockHeightMsg());

    // in genesis, handle sync events first, then batched events.
    // in other blocks, handle batched events first, then sync events.
    if (this.block.height === 0) {
      await this.syncHandlers.process(kafkaPublisher, resultRow);
      await this.batchedHandlers.process(kafkaPublisher, resultRow);
    } else {
      await this.batchedHandlers.process(kafkaPublisher, resultRow);
      await this.syncHandlers.process(kafkaPublisher, resultRow);
    }
    return kafkaPublisher;
  }
}
