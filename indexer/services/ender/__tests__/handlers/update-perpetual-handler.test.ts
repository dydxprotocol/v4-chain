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
  defaultUpdatePerpetualEventV1,
  defaultUpdatePerpetualEventV2,
  defaultUpdatePerpetualEventV3,
} from '../helpers/constants';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  UpdatePerpetualEventV1,
  UpdatePerpetualEventV2,
  UpdatePerpetualEventV3,
} from '@dydxprotocol-indexer/v4-protos';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  eventPerpetualMarketTypeToIndexerPerpetualMarketType,
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

  describe.each([
    [
      'UpdatePerpetualEventV1',
      UpdatePerpetualEventV1.encode(defaultUpdatePerpetualEventV1).finish(),
      defaultUpdatePerpetualEventV1,
    ],
    [
      'UpdatePerpetualEventV2',
      UpdatePerpetualEventV2.encode(defaultUpdatePerpetualEventV2).finish(),
      defaultUpdatePerpetualEventV2,
    ],
    [
      'UpdatePerpetualEventV3',
      UpdatePerpetualEventV3.encode(defaultUpdatePerpetualEventV3).finish(),
      defaultUpdatePerpetualEventV3,
    ],
  ])('%s', (
    _name: string,
    updatePerpetualEventBytes: Uint8Array,
    event: UpdatePerpetualEventV1 | UpdatePerpetualEventV2
    | UpdatePerpetualEventV3,
  ) => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.UPDATE_PERPETUAL,
        updatePerpetualEventBytes,
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
        event,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });

    it('updates an existing perpetual market', async () => {
      const transactionIndex: number = 0;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromUpdatePerpetualEvent({
        updatePerpetualEventBytes,
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
      });
      const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
      await onMessage(kafkaMessage);

      const perpetualMarket:
      PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.findById(
        event.id.toString(),
      );

      expect(perpetualMarket).toEqual(expect.objectContaining({
        id: event.id.toString(),
        ticker: event.ticker,
        marketId: event.marketId,
        atomicResolution: event.atomicResolution,
        liquidityTierId: event.liquidityTier,
        // Add V2-specific field expectations when testing V2 events
        ...('marketType' in event && {
          marketType: eventPerpetualMarketTypeToIndexerPerpetualMarketType(event.marketType),
        }),
      }));
      expect(perpetualMarket).toEqual(
        perpetualMarketRefresher.getPerpetualMarketFromId(
          event.id.toString()));
      expectPerpetualMarketKafkaMessage(producerSendMock, [perpetualMarket!]);
    });
  });
});

function createKafkaMessageFromUpdatePerpetualEvent({
  updatePerpetualEventBytes,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  updatePerpetualEventBytes: Uint8Array,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.UPDATE_PERPETUAL,
      updatePerpetualEventBytes,
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
