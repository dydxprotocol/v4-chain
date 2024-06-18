import { DateTime } from 'luxon';
import { dbHelpers, Transaction } from '@dydxprotocol-indexer/postgres';
import {
  AssetCreateEventV1,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketEventV1,
  SubaccountUpdateEventV1,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { MILLIS_IN_NANOS, SECONDS_IN_MILLIS } from '../../src/constants';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultAssetCreateEvent,
  defaultHeight,
  defaultMarketCreate,
  defaultPreviousHeight,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import Long from 'long';
import { BlockProcessor } from '../../src/lib/block-processor';
import { BatchedHandlers } from '../../src/lib/batched-handlers';
import { SyncHandlers } from '../../src/lib/sync-handlers';
import { mock, MockProxy } from 'jest-mock-extended';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { BLOCK_HEIGHT_WEBSOCKET_MESSAGE_VERSION, KafkaTopics } from '@dydxprotocol-indexer/kafka';

describe('block-processor', () => {
  let batchedHandlers: MockProxy<BatchedHandlers>;
  let syncHandlers: MockProxy<SyncHandlers>;

  beforeEach(() => {
    batchedHandlers = mock<BatchedHandlers>();
    syncHandlers = mock<SyncHandlers>();
    updateBlockCache(defaultPreviousHeight);
  });

  beforeAll(async () => {
    await dbHelpers.migrate();
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
      // updatedPerpetualPositions: [],
      // updatedAssetPositions: [],
    });
  const defaultSubaccountUpdateEventBinary: Uint8Array = Uint8Array.from(
    SubaccountUpdateEventV1.encode(
      defaultSubaccountUpdateEvent,
    ).finish(),
  );

  const defaultMarketEventBinary: Uint8Array = Uint8Array.from(MarketEventV1.encode(
    defaultMarketCreate,
  ).finish());

  const defaultAssetEventBinary: Uint8Array = Uint8Array.from(AssetCreateEventV1.encode(
    defaultAssetCreateEvent,
  ).finish());

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
      DydxIndexerSubtypes.MARKET,
      defaultMarketEventBinary,
      transactionIndex1,
      eventIndex0,
    ),
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.ASSET,
      defaultAssetEventBinary,
      transactionIndex0,
      eventIndex1,
    ),
  ];

  it('batched handlers called before sync handlers for normal blocks', async () => {
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [
        defaultTxHash,
        defaultTxHash2,
      ],
    );

    const txId: number = await Transaction.start();
    const blockProcessor: BlockProcessor = new BlockProcessor(
      block,
      txId,
      defaultDateTime.toString(),
    );
    blockProcessor.batchedHandlers = batchedHandlers;
    blockProcessor.syncHandlers = syncHandlers;
    await blockProcessor.process();
    await Transaction.commit(txId);
    expect(syncHandlers.addHandler).toHaveBeenCalledTimes(2);
    expect(batchedHandlers.addHandler).toHaveBeenCalledTimes(1);
    expect(batchedHandlers.process.mock.invocationCallOrder[0]).toBeLessThan(
      syncHandlers.process.mock.invocationCallOrder[0],
    );
  });

  it('sync handlers called before batched handlers for genesis block', async () => {
    updateBlockCache('-1');
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      0,
      defaultTime,
      events,
      [
        defaultTxHash,
        defaultTxHash2,
      ],
    );

    const txId: number = await Transaction.start();
    const blockProcessor: BlockProcessor = new BlockProcessor(
      block,
      txId,
      defaultDateTime.toString(),
    );
    blockProcessor.batchedHandlers = batchedHandlers;
    blockProcessor.syncHandlers = syncHandlers;
    await blockProcessor.process();
    await Transaction.commit(txId);
    expect(syncHandlers.addHandler).toHaveBeenCalledTimes(2);
    expect(batchedHandlers.addHandler).toHaveBeenCalledTimes(1);
    expect(syncHandlers.process.mock.invocationCallOrder[0]).toBeLessThan(
      batchedHandlers.process.mock.invocationCallOrder[0],
    );
  });

  it('Adds a block height message to the Kafka publisher', async () => {
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [
        defaultTxHash,
        defaultTxHash2,
      ],
    );

    const txId: number = await Transaction.start();
    const blockProcessor: BlockProcessor = new BlockProcessor(
      block,
      txId,
      defaultDateTime.toString(),
    );
    const processor = await blockProcessor.process();
    await Transaction.commit(txId);
    expect(processor.blockHeightMessages).toHaveLength(1);
    expect(processor.blockHeightMessages[0].blockHeight).toEqual(String(defaultHeight));
    expect(processor.blockHeightMessages[0].version)
      .toEqual(BLOCK_HEIGHT_WEBSOCKET_MESSAGE_VERSION);
    expect(processor.blockHeightMessages[0].time).toEqual(defaultDateTime.toString());
  });

  it('createBlockHeightMsg creates a BlockHeightMessage', async () => {
    const block: IndexerTendermintBlock = createIndexerTendermintBlock(
      defaultHeight,
      defaultTime,
      events,
      [
        defaultTxHash,
        defaultTxHash2,
      ],
    );

    const txId: number = await Transaction.start();
    const blockProcessor: BlockProcessor = new BlockProcessor(
      block,
      txId,
      defaultDateTime.toString(),
    );
    await Transaction.commit(txId);

    const msg = blockProcessor.createBlockHeightMsg();
    expect(msg).toEqual({
      topic: KafkaTopics.TO_WEBSOCKETS_BLOCK_HEIGHT,
      message: {
        blockHeight: String(defaultHeight),
        time: defaultDateTime.toString(),
        version: BLOCK_HEIGHT_WEBSOCKET_MESSAGE_VERSION,
      },
    });
  });
});
