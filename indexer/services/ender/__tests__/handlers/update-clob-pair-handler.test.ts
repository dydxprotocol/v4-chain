import {
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  dbHelpers,
  liquidityTierRefresher,
  perpetualMarketRefresher,
  protocolTranslations,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { updateBlockCache } from '../../src/caches/block-cache';
import {
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
  defaultUpdateClobPairEvent,
} from '../helpers/constants';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  UpdateClobPairEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectPerpetualMarketKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { UpdateClobPairHandler } from '../../src/handlers/update-clob-pair-handler';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../src/lib/on-message';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('update-clob-pair-handler', () => {
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
        DydxIndexerSubtypes.UPDATE_CLOB_PAIR,
        UpdateClobPairEventV1.encode(defaultUpdateClobPairEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: UpdateClobPairHandler = new UpdateClobPairHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultUpdateClobPairEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('updates an existing perpetual market', async () => {
    const transactionIndex: number = 0;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromUpdateClobPairEvent({
      updatePerpetualEvent: defaultUpdateClobPairEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const perpetualMarketId: string = perpetualMarketRefresher.getPerpetualMarketFromClobPairId(
      defaultUpdateClobPairEvent.clobPairId.toString(),
    )!.id;
    const perpetualMarket:
    PerpetualMarketFromDatabase | undefined = await PerpetualMarketTable.findById(
      perpetualMarketId,
    );
    expect(perpetualMarket).toEqual(expect.objectContaining({
      clobPairId: defaultUpdateClobPairEvent.clobPairId.toString(),
      status: protocolTranslations.clobStatusToMarketStatus(defaultUpdateClobPairEvent.status),
      quantumConversionExponent: defaultUpdateClobPairEvent.quantumConversionExponent,
      subticksPerTick: defaultUpdateClobPairEvent.subticksPerTick,
      stepBaseQuantums: defaultUpdateClobPairEvent.stepBaseQuantums.toNumber(),
    }));
    expect(perpetualMarket).toEqual(
      perpetualMarketRefresher.getPerpetualMarketFromId(perpetualMarketId));
    expectPerpetualMarketKafkaMessage(producerSendMock, [perpetualMarket!]);
  });
});

function createKafkaMessageFromUpdateClobPairEvent({
  updatePerpetualEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  updatePerpetualEvent: UpdateClobPairEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.UPDATE_CLOB_PAIR,
      UpdateClobPairEventV1.encode(updatePerpetualEvent).finish(),
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
