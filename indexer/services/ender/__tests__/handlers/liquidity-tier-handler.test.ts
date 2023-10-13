import { stats, STATS_FUNCTION_NAME } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  LiquidityTierUpsertEventV1,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  BlockTable,
  dbHelpers,
  LiquidityTiersColumns,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  Ordering,
  QUOTE_CURRENCY_ATOMIC_RESOLUTION,
  TendermintEventTable,
  testConstants,
  protocolTranslations,
  liquidityTierRefresher,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  MarketTable,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectPerpetualMarketKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { LiquidityTierHandler } from '../../src/handlers/liquidity-tier-handler';
import {
  defaultHeight, defaultLiquidityTierUpsertEvent, defaultPreviousHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { defaultLiquidityTier } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import _ from 'lodash';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('liquidityTierHandler', () => {
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
        DydxIndexerSubtypes.LIQUIDITY_TIER,
        LiquidityTierUpsertEventV1.encode(defaultLiquidityTierUpsertEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: LiquidityTierHandler = new LiquidityTierHandler(
        block,
        indexerTendermintEvent,
        0,
        defaultLiquidityTierUpsertEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('creates new liquidity tier', async () => {
    const transactionIndex: number = 0;
    const liquidityTierEvent: LiquidityTierUpsertEventV1 = defaultLiquidityTierUpsertEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
      liquidityTierEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    // Confirm there is no existing liquidity tier
    await expectNoExistingLiquidityTiers();
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newLiquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [], {
        orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
      });
    expect(newLiquidityTiers.length).toEqual(1);
    expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent);
    expectTimingStats();
    validateLiquidityTierRefresher(defaultLiquidityTierUpsertEvent);
    expectKafkaMessages(producerSendMock, liquidityTierEvent, 0);
  });

  it('updates existing liquidity tier', async () => {
    const transactionIndex: number = 0;
    const liquidityTierEvent: LiquidityTierUpsertEventV1 = defaultLiquidityTierUpsertEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
      liquidityTierEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    // Create existing liquidity tier
    await LiquidityTiersTable.upsert(defaultLiquidityTier);

    // create perpetual market with existing liquidity tier to test websockets
    await Promise.all([
      MarketTable.create(testConstants.defaultMarket),
      MarketTable.create(testConstants.defaultMarket2),
    ]);
    await Promise.all([
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket),
      PerpetualMarketTable.create(testConstants.defaultPerpetualMarket2),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newLiquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [], {
        orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
      });
    expect(newLiquidityTiers.length).toEqual(1);
    expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent);
    expectTimingStats();
    validateLiquidityTierRefresher(defaultLiquidityTierUpsertEvent);
    expectKafkaMessages(producerSendMock, liquidityTierEvent, 2);
  });
});

function expectTimingStats() {
  expectTimingStat('upsert_liquidity_tier');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    {
      className: 'LiquidityTierHandler',
      eventType: 'LiquidityTierUpsertEvent',
      fnName,
    },
  );
}

export function expectLiquidityTier(
  liquidityTierFromDb: LiquidityTiersFromDatabase,
  event: LiquidityTierUpsertEventV1,
): void {
  expect(liquidityTierFromDb).toEqual(expect.objectContaining({
    id: event.id,
    name: event.name,
    initialMarginPpm: event.initialMarginPpm.toString(),
    maintenanceFractionPpm: event.maintenanceFractionPpm.toString(),
    basePositionNotional: protocolTranslations.quantumsToHuman(
      event.basePositionNotional.toString(),
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).toFixed(),
  }));
}

function createKafkaMessageFromLiquidityTiersEvent({
  liquidityTierEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  liquidityTierEvent: LiquidityTierUpsertEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.LIQUIDITY_TIER,
      LiquidityTierUpsertEventV1.encode(liquidityTierEvent).finish(),
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

async function expectNoExistingLiquidityTiers() {
  // Confirm there is no existing liquidity tier
  const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
    {},
    [], {
      orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
    });

  expect(liquidityTiers.length).toEqual(0);
}

function validateLiquidityTierRefresher(
  liquidityTierEvent: LiquidityTierUpsertEventV1,
) {
  const liquidityTier:
  LiquidityTiersFromDatabase = liquidityTierRefresher.getLiquidityTierFromId(
    liquidityTierEvent.id,
  );

  expect(liquidityTier).toEqual({
    id: liquidityTierEvent.id,
    name: liquidityTierEvent.name,
    initialMarginPpm: liquidityTierEvent.initialMarginPpm.toString(),
    maintenanceFractionPpm: liquidityTierEvent.maintenanceFractionPpm.toString(),
    basePositionNotional: protocolTranslations.quantumsToHuman(
      liquidityTierEvent.basePositionNotional.toString(),
      QUOTE_CURRENCY_ATOMIC_RESOLUTION,
    ).toFixed(),
  });
}

function expectKafkaMessages(
  producerSendMock: jest.SpyInstance,
  liquidityTier: LiquidityTierUpsertEventV1,
  numPerpetualMarkets: number,
) {
  const perpetualMarkets: PerpetualMarketFromDatabase[] = _.filter(
    perpetualMarketRefresher.getPerpetualMarketsList(),
    (perpetualMarket: PerpetualMarketFromDatabase) => {
      return perpetualMarket.liquidityTierId === liquidityTier.id;
    },
  );
  expect(perpetualMarkets.length).toEqual(numPerpetualMarkets);

  if (perpetualMarkets.length === 0) {
    return;
  }
  expectPerpetualMarketKafkaMessage(producerSendMock, perpetualMarkets);
}
