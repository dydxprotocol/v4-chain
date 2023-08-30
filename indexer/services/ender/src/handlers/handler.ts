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
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketMessage,
  OffChainUpdateV1,
  SubaccountId,
  SubaccountMessage,
} from '@dydxprotocol-indexer/v4-protos';
import { DateTime } from 'luxon';

import config from '../config';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent, EventMessage, SingleTradeMessage } from '../lib/types';

export type HandlerInitializer = new (
  block: IndexerTendermintBlock,
  indexerTendermintEvent: IndexerTendermintEvent,
  txId: number,
  event: EventMessage,
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
  event: T;
  abstract eventType: string;

  constructor(
    block: IndexerTendermintBlock,
    indexerTendermintEvent: IndexerTendermintEvent,
    txId: number,
    event: T,
  ) {
    this.block = block;
    this.indexerTendermintEvent = indexerTendermintEvent;
    this.timestamp = DateTime.fromJSDate(block.time!);
    this.txId = txId;
    this.event = event;
  }

  /**
   * Returns the ids for the event where no other events with the same id can be processed in
   * parallel.
   */
  public abstract getParallelizationIds(): string[];

  /**
   * Processes the event and updates Postgres in the transaction specified by
   * txId provided in the constructor, then returns all consolidated kafka events to be
   * written to Kafka.
   */
  public abstract internalHandle(): Promise<ConsolidatedKafkaEvent[]>;

  /**
   * Handle the event and export timing stats
   */
  public async handle(): Promise<ConsolidatedKafkaEvent[]> {
    const start: number = Date.now();
    try {
      return await this.internalHandle();
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
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_subaccount_kafka_event`, 1);
    const subaccountMessage: SubaccountMessage = {
      blockHeight: this.block.height.toString(),
      transactionIndex: indexerTendermintEventToTransactionIndex(
        this.indexerTendermintEvent,
      ),
      eventIndex: this.indexerTendermintEvent.eventIndex,
      contents,
      subaccountId,
      version: SUBACCOUNTS_WEBSOCKET_MESSAGE_VERSION,
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
  ): ConsolidatedKafkaEvent {
    stats.increment(`${config.SERVICE_NAME}.create_vulcan_kafka_event`, 1);

    return {
      topic: KafkaTopics.TO_VULCAN,
      message: {
        key,
        value: offChainUpdate,
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
