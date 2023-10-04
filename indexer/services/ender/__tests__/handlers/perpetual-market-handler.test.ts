import { stats, STATS_FUNCTION_NAME } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketCreateEventV1,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  PerpetualMarketColumns,
  PerpetualMarketFromDatabase,
  PerpetualMarketTable,
  dbHelpers,
  MarketTable,
  Ordering,
  testConstants,
  marketRefresher,
  BlockTable,
  TendermintEventTable,
  perpetualMarketRefresher,
  LiquidityTiersTable,
  liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectPerpetualMarket,
  expectPerpetualMarketKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { PerpetualMarketCreationHandler } from '../../src/handlers/perpetual-market-handler';
import {
  defaultPerpetualMarketCreateEvent,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('perpetualMarketHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await Promise.all([
      BlockTable.create(testConstants.defaultBlock),
      BlockTable.create(testConstants.defaultBlock2),
    ]);
    await Promise.all([
      TendermintEventTable.create(testConstants.defaultTendermintEvent),
      TendermintEventTable.create(testConstants.defaultTendermintEvent2),
      TendermintEventTable.create(testConstants.defaultTendermintEvent3),
    ]);
    updateBlockCache(defaultPreviousHeight);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
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
        DydxIndexerSubtypes.PERPETUAL_MARKET,
        PerpetualMarketCreateEventV1.encode(defaultPerpetualMarketCreateEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: PerpetualMarketCreationHandler = new PerpetualMarketCreationHandler(
        block,
        indexerTendermintEvent,
        0,
        defaultPerpetualMarketCreateEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('fails when market doesnt exist for perpetual market', async () => {
    const transactionIndex: number = 0;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
      perpetualMarketEvent: defaultPerpetualMarketCreateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await expect(onMessage(kafkaMessage)).rejects.toThrowError();
  });

  it('fails when liquidity tier doesnt exist for perpetual market', async () => {
    await MarketTable.create(testConstants.defaultMarket);
    await marketRefresher.updateMarkets();
    const transactionIndex: number = 0;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
      perpetualMarketEvent: defaultPerpetualMarketCreateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await expect(onMessage(kafkaMessage)).rejects.toThrowError();
  });

  it('creates new perpetual market', async () => {
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
    ]);
    await liquidityTierRefresher.updateLiquidityTiers();
    await marketRefresher.updateMarkets();

    const transactionIndex: number = 0;

    const perpetualMarketEvent: PerpetualMarketCreateEventV1 = defaultPerpetualMarketCreateEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
      perpetualMarketEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    // Confirm there is no existing perpetualMarket.
    await expectNoExistingPerpetualMarkets();

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newPerpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
      {},
      [], {
        orderBy: [[PerpetualMarketColumns.id, Ordering.ASC]],
      });
    expect(newPerpetualMarkets.length).toEqual(1);
    expectPerpetualMarketMatchesEvent(perpetualMarketEvent, newPerpetualMarkets[0]);
    expectTimingStats();
    const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher.getPerpetualMarketFromId('0');
    expect(perpetualMarket).toBeDefined();
    expectPerpetualMarket(perpetualMarket!, perpetualMarketEvent);
    expectPerpetualMarketKafkaMessage(producerSendMock, [perpetualMarket!]);
  });
});

function expectTimingStats() {
  expectTimingStat('create_perpetual_market');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    {
      className: 'PerpetualMarketCreationHandler',
      eventType: 'PerpetualMarketCreateEvent',
      fnName,
    },
  );
}

function expectPerpetualMarketMatchesEvent(
  perpetual: PerpetualMarketCreateEventV1,
  perpetualMarket: PerpetualMarketFromDatabase,
) {
  expect(perpetualMarket).toEqual(expect.objectContaining({
    id: perpetual.id.toString(),
    clobPairId: perpetual.clobPairId.toString(),
    ticker: perpetual.ticker,
    marketId: perpetual.marketId,
    quantumConversionExponent: perpetual.quantumConversionExponent,
    atomicResolution: perpetual.atomicResolution,
    subticksPerTick: perpetual.subticksPerTick,
    stepBaseQuantums: Number(perpetual.stepBaseQuantums),
    liquidityTierId: perpetual.liquidityTier,
  }));
}

function createKafkaMessageFromPerpetualMarketEvent({
  perpetualMarketEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  perpetualMarketEvent: PerpetualMarketCreateEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.PERPETUAL_MARKET,
      PerpetualMarketCreateEventV1.encode(perpetualMarketEvent).finish(),
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

async function expectNoExistingPerpetualMarkets() {
  // Confirm there is no existing perpetual markets
  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {},
    [], {
      orderBy: [[PerpetualMarketColumns.id, Ordering.ASC]],
    });

  expect(perpetualMarkets.length).toEqual(0);
}
