import {
  logger,
  ParseMessageError,
  runFuncWithTimingStat,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  MARKETS_WEBSOCKET_MESSAGE_VERSION,
  SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
  TRADES_WEBSOCKET_MESSAGE_VERSION,
  KafkaTopics,
} from '@dydxprotocol-indexer/kafka';
import { SubaccountMessageContents } from '@dydxprotocol-indexer/postgres';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketMessage,
  OffChainUpdateV1,
  SubaccountId,
} from '@dydxprotocol-indexer/v4-protos';
import { IHeaders } from 'kafkajs';
import { DateTime } from 'luxon';
import * as pg from 'pg';

import config from '../config';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import {
  AnnotatedSubaccountMessage, ConsolidatedKafkaEvent, EventMessage, SingleTradeMessage,
} from '../lib/types';

export type HandlerInitializer = new (
  block: IndexerTendermintBlock,
  blockEventIndex: number,
  indexerTendermintEvent: IndexerTendermintEvent,
  txId: number,
  event: EventMessage,
  messageReceivedTimestamp?: string,
) => Handler<EventMessage>;

/**
 * Base class for all event handlers. Each event handler is responsible for processing a
 * specific type of event, with some handlers such as OrderFillHandler and MarketHandler
 * being used to triage the event to more specific handlers.
 */
export abstract class Handler<T> {
  block: IndexerTendermintBlock;
  indexerTendermintEvent: IndexerTendermintEvent;
  timestamp: DateTime;
  txId: number;
  // The index of the event in the block.
  blockEventIndex: number;
  event: T;
  abstract eventType: string;
  messageReceivedTimestamp?: string;

  constructor(
    block: IndexerTendermintBlock,
    blockEventIndex: number,
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    event: T,
    messageReceivedTimestamp?: string,
  ) {
    this.block = block;
    this.blockEventIndex = blockEventIndex;
    this.indexerTendermintEvent = indexerTendermintEvent;
    this.timestamp = DateTime.fromJSDate(block.time!);
    this.txId = txId;
    this.event = event;
    this.messageReceivedTimestamp = messageReceivedTimestamp;
  }

  /**
   * Returns the ids for the event where no other events with the same id can be processed in
   * parallel.
   */
  public abstract getParallelizationIds(): string[];

  /**
   * Performs side effects based upon the results returned from the SQL based handler
   * implementations and then returns all consolidated Kafka events to be written to Kafka.
   */
  public abstract internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]>;

  /**
   * Performs side effects based upon the results returned from the SQL based handler
   * implementations and then returns all consolidated Kafka events to be written to Kafka.
   *
   * Wraps internalHandle with timing information.
   */
  public async handle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    const start: number = Date.now();
    try {
      return await this.internalHandle(resultRow);
    } finally {
      stats.timing(
        `${config.SERVICE_NAME}.handle_event.timing`,
        Date.now() - start,
        this.generateTimingStatsOptions(),
      );
    }
  }

  protected logAndThrowParseMessageError(
    message: string,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    additionalParameters: any = {},
  ): void {
    logger.error({
      at: `${this.constructor.name}#logAndThrowParseMessageError`,
      message,
      blockHeight: this.block.height,
      ...additionalParameters,
    });
    throw new ParseMessageError(message);
  }

  protected generateConsolidatedSubaccountKafkaEvent(
    contents: string,
    subaccountId: SubaccountId,
    orderId?: string,
    isFill?: boolean,
    subaccountMessageContents?: SubaccountMessageContents,
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_subaccount_kafka_event`, 1);
    const subaccountMessage: AnnotatedSubaccountMessage = {
      blockHeight: this.block.height.toString(),
      transactionIndex: indexerTendermintEventToTransactionIndex(
        this.indexerTendermintEvent,
      ),
      eventIndex: this.indexerTendermintEvent.eventIndex,
      contents,
      subaccountId,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
      orderId,
      isFill,
      subaccountMessageContents,
    };

    return {
      topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
      message: subaccountMessage,
    };
  }

  protected generateConsolidatedMarketKafkaEvent(
    contents: string,
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_market_kafka_event`, 1);
    const marketMessage: MarketMessage = {
      contents,
      version: MARKETS_WEBSOCKET_MESSAGE_VERSION,
    };

    return {
      topic: KafkaTopics.TO_WEBSOCKETS_MARKETS,
      message: marketMessage,
    };
  }

  protected generateConsolidatedTradeKafkaEvent(
    contents: string,
    clobPairId: string,
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_trade_kafka_event`, 1);
    const tradeMessage: SingleTradeMessage = {
      blockHeight: this.block.height.toString(),
      transactionIndex: indexerTendermintEventToTransactionIndex(
        this.indexerTendermintEvent,
      ),
      eventIndex: this.indexerTendermintEvent.eventIndex,
      contents,
      clobPairId,
      version: TRADES_WEBSOCKET_MESSAGE_VERSION,
    };

    return {
      topic: KafkaTopics.TO_WEBSOCKETS_TRADES,
      message: tradeMessage,
    };
  }

  protected generateConsolidatedVulcanKafkaEvent(
    key: Buffer,
    offChainUpdate: OffChainUpdateV1,
    headers?: IHeaders,
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_vulcan_kafka_event`, 1);

    return {
      topic: KafkaTopics.TO_VULCAN,
      message: {
        key,
        value: offChainUpdate,
        headers,
      },
    };
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  protected generateTimingStatsOptions(fnName?: string): any {
    return {
      className: this.constructor.name,
      eventType: this.eventType,
      fnName,
    };
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  protected async runFuncWithTimingStatAndErrorLogging(
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    promise: Promise<any>,
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    options: any,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ): Promise<any> {
    try {
      const ret = await runFuncWithTimingStat(promise, options);
      return ret;
    } catch (error) {
      logger.error({
        at: `${this.constructor.name}#runFuncWithTimingStatAndErrorLogging`,
        message: `handlerError: ${error.message}`,
        options,
        event: JSON.stringify(this.event),
        eventType: this.eventType,
        handler: this.constructor.name,
        transactionIndex: indexerTendermintEventToTransactionIndex(
          this.indexerTendermintEvent,
        ),
        eventIndex: this.indexerTendermintEvent.eventIndex,
        stacktrace: error.stack,
      });
      throw error;
    }
  }
}
