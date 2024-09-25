import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  dbHelpers,
  liquidityTierRefresher,
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
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectPerpetualMarketKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { UpdatePerpetualHandler } from '../../src/handlers/update-perpetual-handler';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../src/lib/on-message';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('update-perpetual-handler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    updateBlockCache(defaultPreviousHeight);
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await liquidityTierRefresher.updateLiquidityTiers();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    perpetualMarketRefresher.clear();
    liquidityTierRefresher.clear();
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
        UpdatePerpetualEventV1.encode(defaultUpdatePerpetualEvent).finish(),
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
        0,
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
    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
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
      dangerIndexPpm: defaultUpdatePerpetualEvent.dangerIndexPpm,
      isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock:
        defaultUpdatePerpetualEvent.isolatedMarketMaxCumulativeInsuranceFundDeltaPerBlock,
    }));
    expect(perpetualMarket).toEqual(
      perpetualMarketRefresher.getPerpetualMarketFromId(
        defaultUpdatePerpetualEvent.id.toString()));
    expectPerpetualMarketKafkaMessage(producerSendMock, [perpetualMarket!]);
  });
});

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
      UpdatePerpetualEventV1.encode(updatePerpetualEvent).finish(),
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
