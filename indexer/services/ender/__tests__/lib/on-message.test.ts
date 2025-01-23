import { DateTime } from 'luxon';
import {
  assetRefresher,
  AssetTable,
  BlockFromDatabase,
  BlockTable,
  dbHelpers,
  IsoString,
  LiquidityTiersTable,
  liquidityTierRefresher,
  MarketTable,
  Ordering,
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  TendermintEventFromDatabase,
  TendermintEventTable,
  testConstants,
  testMocks,
  TransactionFromDatabase,
  TransactionTable,
} from '@dydxprotocol-indexer/postgres';
import {
  FundingEventV1,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketEventV1,
  PerpetualMarketCreateEventV2,
  PerpetualMarketCreateEventV3,
  SubaccountMessage,
  SubaccountUpdateEventV1,
  Timestamp,
  TransferEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { onMessage } from '../../src/lib/on-message';
import { KafkaMessage } from 'kafkajs';
import {
  createKafkaMessage, KafkaTopics, producer,
} from '@dydxprotocol-indexer/kafka';
import { MILLIS_IN_NANOS, SECONDS_IN_MILLIS } from '../../src/constants';
import { ConsolidatedKafkaEvent, DydxIndexerSubtypes } from '../../src/lib/types';
import config from '../../src/config';
import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  defaultFundingUpdateSampleEvent,
  defaultHeight,
  defaultMarketModify,
  defaultPerpetualMarketCreateEventV2,
  defaultPerpetualMarketCreateEventV3,
  defaultPreviousHeight,
  defaultSubaccountMessage,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import Long from 'long';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { expectPerpetualMarketMatchesEvent } from '../helpers/postgres-helpers';

describe('on-message', () => {
  let producerSendMock: jest.SpyInstance;
  const loggerError = jest.spyOn(logger, 'error');

  beforeEach(() => {
    producerSendMock = jest.spyOn(producer, 'send');
    producerSendMock.mockImplementation(() => {
    });
    updateBlockCache(defaultPreviousHeight);
  });

  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
    await createPostgresFunctions();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultDateTime: DateTime = DateTime.utc(2022, 6, 1, 12, 1, 1, 2);
  const defaultTime: Timestamp = {
    seconds: Long.fromValue(Math.floor(defaultDateTime.toSeconds()), true),
    nanos: (defaultDateTime.toMillis() % SECONDS_IN_MILLIS) * MILLIS_IN_NANOS,
  };
  const defaultTxHash: string = '0x32343534306431622d306461302d343831322d613730372d3965613162336162';
  const defaultTxHash2: string = '0x32363534306431622d306461302d343831322d613730372d3965613162336162';
  const defaultSubaccountUpdateEvent: SubaccountUpdateEventV1 = SubaccountUpdateEventV1
    .fromPartial({
      subaccountId: {
        owner: '',
        number: 0,
      },
    });
  const defaultSubaccountUpdateEventBinary: Uint8Array = Uint8Array.from(
    SubaccountUpdateEventV1.encode(
      defaultSubaccountUpdateEvent,
    ).finish(),
  );

  const defaultTransferEvent: TransferEventV1 = TransferEventV1.fromPartial({
    sender: {
      subaccountId: {
        owner: '',
        number: 0,
      },
    },
    recipient: {
      subaccountId: {
        owner: '',
        number: 0,
      },
    },
    assetId: 0,
    amount: Long.fromValue(100, true),
  });
  const defaultTransferEventBinary: Uint8Array = Uint8Array.from(TransferEventV1.encode(
    defaultTransferEvent,
  ).finish());

  const defaultFundingEventBinary: Uint8Array = Uint8Array.from(FundingEventV1.encode(
    defaultFundingUpdateSampleEvent,
  ).finish());

  const defaultMarketEventBinary: Uint8Array = Uint8Array.from(MarketEventV1.encode(
    defaultMarketModify,
  ).finish());

  const defaultPerpetualMarketEventV2Binary: Uint8Array = Uint8Array.from(
    PerpetualMarketCreateEventV2.encode(
      defaultPerpetualMarketCreateEventV2,
    ).finish(),
  );

  const defaultPerpetualMarketEventV3Binary: Uint8Array = Uint8Array.from(
    PerpetualMarketCreateEventV3.encode(
      defaultPerpetualMarketCreateEventV3,
    ).finish(),
  );

  it('successfully processes block with transaction event', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('successfully processes block with transaction event with unset version', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
        0,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('successfully processes block with transfer event', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      AssetTable.create(testConstants.defaultAsset),
    ]);
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        defaultTransferEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it.each([
    [
      'PerpetualMarketCreateV2',
      defaultPerpetualMarketCreateEventV2,
      defaultPerpetualMarketEventV2Binary,
    ],
    [
      'PerpetualMarketCreateV3',
      defaultPerpetualMarketCreateEventV3,
      defaultPerpetualMarketEventV3Binary,
    ],
  ])('successfully processes block with %s and its funding events', async (
    _name: string,
    marketCreateEvent: any,
    marketCreateEventBinary: Uint8Array,
  ) => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      MarketTable.create(testConstants.defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier2),
    ]);
    await Promise.all([
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket2),
    ]);
    await Promise.all([
      perpetualMarketRefresher.updatePerpetualMarkets(),
      liquidityTierRefresher.updateLiquidityTiers(),
    ]);

    const transactionIndex: number = -1;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.PERPETUAL_MARKET,
        marketCreateEventBinary,
        0,
        eventIndex,
      ),
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.FUNDING,
        marketCreateEventBinary,
        transactionIndex,
        eventIndex + 1,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), 0, eventIndex),
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex + 1),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    const newPerpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [], {
        orderBy: [[PerpetualMarketColumns.id, Ordering.ASC]],
      });
    expect(newPerpetualMarkets.length).toEqual(2);
    expectPerpetualMarketMatchesEvent(marketCreateEvent, newPerpetualMarkets[0]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('successfully processes block with funding event', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      MarketTable.create(testConstants.defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier2),
    ]);
    await Promise.all([
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket),
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket2),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const transactionIndex: number = -1;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.FUNDING,
        defaultFundingEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('throws error while processing unparsable messages', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    // unparsable transfer event
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));
    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      new Error('Could not parse message TransferEvent must have either a sender subaccount id or sender wallet address'),
    );
  });

  it('skips over unknown events while processing', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      AssetTable.create(testConstants.defaultAsset),
    ]);
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const eventIndex1: number = 1;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        defaultTransferEventBinary,
        transactionIndex,
        eventIndex,
      ),
      createIndexerTendermintEvent(
        'unknown',
        defaultTransferEventBinary,
        transactionIndex,
        eventIndex1,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
      at: 'helpers#indexerTendermintEventToEventWithType',
      message: 'Unable to parse event subtype: unknown',
    }));
    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
  });

  it('successfully processes block with market event', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
    ]);
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.MARKET,
        defaultMarketEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('successfully processes block with block event', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      MarketTable.create(testConstants.defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier2),
    ]);
    await Promise.all([
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket),
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket2),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const transactionIndex: number = -1;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.FUNDING,
        defaultFundingEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });

  it('successfully processes block with transaction event and block event', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      MarketTable.create(testConstants.defaultMarket2),
    ]);
    await Promise.all([
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier2),
    ]);
    await Promise.all([
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket),
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket2),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const blockTransactionIndex: number = -1;
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      // MARKET is a transaction event.
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.MARKET,
        defaultMarketEventBinary,
        transactionIndex,
        eventIndex,
      ),
      // FUNDING is a block event.
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.FUNDING,
        defaultFundingEventBinary,
        blockTransactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTendermintEvent(defaultHeight.toString(), blockTransactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);
  });

  it('successfully processes block with multiple transactions', async () => {
    const transactionIndex0: number = 0;
    const transactionIndex1: number = 1;
    const eventIndex0: number = 0;
    const eventIndex1: number = 1;

    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex0,
        eventIndex0,
      ),
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex0,
        eventIndex1,
      ),
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex1,
        eventIndex0,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [
        defaultTxHash,
        defaultTxHash2,
      ],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex0, eventIndex0),
      expectTendermintEvent(defaultHeight.toString(), transactionIndex0, eventIndex1),
      expectTendermintEvent(defaultHeight.toString(), transactionIndex1, eventIndex0),
      expectTransactionWithHash([defaultTxHash, defaultTxHash2]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
    expect(stats.gauge).toHaveBeenCalledWith('ender.processing_block_height', expect.any(Number));
    expect(stats.timing).toHaveBeenCalledWith('ender.processed_block.timing',
      expect.any(Number), 1, { success: 'true' });
  });
  it('successfully batches up kafka messages', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      AssetTable.create(testConstants.defaultAsset),
    ]);
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        defaultTransferEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES = 1;
    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    expect(producerSendMock).toHaveBeenCalledTimes(3);
    // First message batch sent should contain the first message
    expect(producerSendMock.mock.calls[0][0].messages).toHaveLength(1);
    // Second message batch should contain the second message
    expect(producerSendMock.mock.calls[1][0].messages).toHaveLength(1);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
  });

  it('skips sending websocket messages if flag is set', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    // Mock the return of kafka messages that in total have size > max message size in bytes
    const kafkaMessages: ConsolidatedKafkaEvent[] = [];
    let totalSizeBytes: number = 0;
    const subaccountByteChange: number = Buffer.from(
      Uint8Array.from(SubaccountMessage.encode(defaultSubaccountMessage).finish()),
    ).byteLength;
    while (totalSizeBytes <= config.KAFKA_MAX_BATCH_WEBSOCKET_MESSAGE_SIZE_BYTES) {
      kafkaMessages.push({
        topic: KafkaTopics.TO_WEBSOCKETS_SUBACCOUNTS,
        message: defaultSubaccountMessage,
      });
      totalSizeBytes += subaccountByteChange;
    }

    config.SEND_WEBSOCKET_MESSAGES = false;
    await onMessage(kafkaMessage);
    await Promise.all([
      expectTendermintEvent(defaultHeight.toString(), transactionIndex, eventIndex),
      expectTransactionWithHash([defaultTxHash]),
      expectBlock(defaultHeight.toString(), defaultDateTime.toISO()),
    ]);

    // No messages should have been sent.
    expect(producerSendMock).toHaveBeenCalledTimes(0);

    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
  });

  it('skips processing a block if it is already in the database', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    // Update block cache with default height
    updateBlockCache(defaultHeight.toString());

    await onMessage(kafkaMessage);

    expect(stats.increment).toHaveBeenCalledWith(`${config.SERVICE_NAME}.block_already_parsed`, 1);
    expect(stats.increment).toHaveBeenCalledWith('ender.received_kafka_message', 1);
    expect(stats.timing).not.toHaveBeenCalledWith(
      'ender.message_time_in_queue', expect.any(Number), 1, { topic: KafkaTopics.TO_ENDER });
  });

  it('refreshes caches if transaction is rolled back', async () => {
    const transactionIndex: number = 0;
    const eventIndex: number = 0;
    const events: IndexerTendermintEvent[] = [
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
        defaultSubaccountUpdateEventBinary,
        transactionIndex,
        eventIndex,
      ),
    ];

    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [defaultTxHash],
    );
    const binaryBlock: Uint8Array = Uint8Array.from(IndexerTendermintBlock.encode(block).finish());
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    // Update block cache with default height
    updateBlockCache(defaultHeight.toString());
    await testMocks.seedData();

    // Initialize assetRefresher
    await assetRefresher.updateAssets();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await onMessage(kafkaMessage);

    expect(assetRefresher.getAssetsMap()).not.toEqual({});
    expect(perpetualMarketRefresher.getPerpetualMarketsMap()).not.toEqual({});
  });
});

async function expectTendermintEvent(
  blockHeight: string,
  transactionIndex: number,
  eventIndex: number,
) {
  const tendermintEventId: Buffer = TendermintEventTable.createEventId(
    blockHeight,
    transactionIndex,
    eventIndex,
  );
  const tendermintEvent: TendermintEventFromDatabase | undefined = await
  TendermintEventTable.findById(
    tendermintEventId,
    { readReplica: true },
  );

  expect(tendermintEvent).not.toEqual(undefined);
}

async function expectTransactionWithHash(transactionHash: string[]) {
  const transactions: TransactionFromDatabase[] = await TransactionTable.findAll(
    { transactionHash },
    [],
    { readReplica: true },
  );

  expect(transactions.length).toEqual(transactionHash.length);
}

async function expectBlock(
  height: string,
  time: IsoString,
) {
  const block: BlockFromDatabase | undefined = await BlockTable.findByBlockHeight(
    height,
    { readReplica: true },
  );

  expect(block).not.toEqual(undefined);
  expect(block?.blockHeight).toEqual(height);
  expect(block?.time).toEqual(time);
}
