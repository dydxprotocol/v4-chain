import { STATS_FUNCTION_NAME, stats } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  dbHelpers,
  perpetualMarketRefresher,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { updateBlockCache } from '../../src/caches/block-cache';
import {
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
  defaultUpdatePerpetualEvent,
} from '../helpers/constants';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  UpdatePerpetualEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { binaryToBase64String, createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { UpdatePerpetualHandler } from '../../src/handlers/update-perpetual-handler';
import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../src/lib/on-message';

describe('update-perpetual-handler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    updateBlockCache(defaultPreviousHeight);
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    perpetualMarketRefresher.clear();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.UPDATE_PERPETUAL,
        binaryToBase64String(
          UpdatePerpetualEventV1.encode(defaultUpdatePerpetualEvent).finish(),
        ),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: UpdatePerpetualHandler = new UpdatePerpetualHandler(
        block,
        indexerTendermintEvent,
        0,
        defaultUpdatePerpetualEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('updates an existing perpetual market', async () => {
    const transactionIndex: number = 0;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromUpdatePerpetualEvent({
      updatePerpetualEvent: defaultUpdatePerpetualEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    await onMessage(kafkaMessage);

    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.findById(
      defaultUpdatePerpetualEvent.id.toString(),
    );
    expect(perpetualMarket).toEqual(expect.objectContaining({
      id: defaultUpdatePerpetualEvent.id.toString(),
      ticker: defaultUpdatePerpetualEvent.ticker,
      marketId: defaultUpdatePerpetualEvent.marketId,
      atomicResolution: defaultUpdatePerpetualEvent.atomicResolution,
      liquidityTierId: defaultUpdatePerpetualEvent.liquidityTier,
    }));
    expectTimingStats();
  });
});

function expectTimingStats() {
  expectTimingStat('update_perpetual');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    {
      className: 'UpdatePerpetualHandler',
      eventType: 'UpdatePerpetualEventV1',
      fnName,
    },
  );
}

function createKafkaMessageFromUpdatePerpetualEvent({
  updatePerpetualEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  updatePerpetualEvent: UpdatePerpetualEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.UPDATE_PERPETUAL,
      binaryToBase64String(
        UpdatePerpetualEventV1.encode(updatePerpetualEvent).finish(),
      ),
      transactionIndex,
      0,
    ),
  );

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}
