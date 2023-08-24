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
  liquidityTierRefresher,
  LiquidityTiersColumns,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  Ordering,
  QUOTE_CURRENCY_ATOMIC_RESOLUTION,
  TendermintEventTable,
  testConstants,
  protocolTranslations,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { LiquidityTierHandler } from '../../src/handlers/liquidity-tier-handler';
import {
  defaultHeight, defaultLiquidityTierUpsertEvent, defaultPreviousHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { defaultLiquidityTier } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('liquidityTierHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
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
        binaryToBase64String(
          LiquidityTierUpsertEventV1.encode(defaultLiquidityTierUpsertEvent).finish(),
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

    await onMessage(kafkaMessage);

    const newLiquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [], {
        orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
      });
    expect(newLiquidityTiers.length).toEqual(1);
    expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent);
    expectTimingStats();
    const liquidityTier:
    LiquidityTiersFromDatabase = liquidityTierRefresher.getLiquidityTierFromId(0);
    expect(liquidityTier).toBeDefined();
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
    // Confirm there is no existing liquidity tier
    await LiquidityTiersTable.upsert(defaultLiquidityTier);

    await onMessage(kafkaMessage);

    const newLiquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
      {},
      [], {
        orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
      });
    expect(newLiquidityTiers.length).toEqual(1);
    expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent);
    expectTimingStats();
    const liquidityTier:
    LiquidityTiersFromDatabase = liquidityTierRefresher.getLiquidityTierFromId(0);
    expect(liquidityTier).toBeDefined();
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
    ).toFixed(6),
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
      binaryToBase64String(
        LiquidityTierUpsertEventV1.encode(liquidityTierEvent).finish(),
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

async function expectNoExistingLiquidityTiers() {
  // Confirm there is no existing liquidity tier
  const liquidityTiers: LiquidityTiersFromDatabase[] = await LiquidityTiersTable.findAll(
    {},
    [], {
      orderBy: [[LiquidityTiersColumns.id, Ordering.ASC]],
    });

  expect(liquidityTiers.length).toEqual(0);
}
