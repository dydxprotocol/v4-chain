import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  LiquidityTierUpsertEventV1,
  LiquidityTierUpsertEventV2,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  BlockTable,
  dbHelpers,
  LiquidityTiersColumns,
  LiquidityTiersFromDatabase,
  LiquidityTiersTable,
  Ordering,
  TendermintEventTable,
  testConstants,
  liquidityTierRefresher,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  PerpetualMarketTable,
  MarketTable,
  protocolTranslations,
  QUOTE_CURRENCY_ATOMIC_RESOLUTION,
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
  defaultHeight,
  defaultLiquidityTierUpsertEventV2,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
  defaultLiquidityTierUpsertEventV1,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { defaultLiquidityTier } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';
import _ from 'lodash';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('liquidityTierHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
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

  describe('liquidityTierHandlerV1', () => {
    describe('getParallelizationIds', () => {
      it('returns the correct parallelization ids', () => {
        const transactionIndex: number = 0;
        const eventIndex: number = 0;

        const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
          DydxIndexerSubtypes.LIQUIDITY_TIER,
          LiquidityTierUpsertEventV1.encode(defaultLiquidityTierUpsertEventV1).finish(),
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
          0,
          indexerTendermintEvent,
          0,
          defaultLiquidityTierUpsertEventV1,
        );

        expect(handler.getParallelizationIds()).toEqual([]);
      });
    });

    it('creates new liquidity tier', async () => {
      const transactionIndex: number = 0;
      const liquidityTierEvent: LiquidityTierUpsertEventV1 = defaultLiquidityTierUpsertEventV1;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
        liquidityTierEvent: LiquidityTierUpsertEventV1.encode(
          liquidityTierEvent,
        ).finish(),
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 1,
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
      expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent, false);
      validateLiquidityTierRefresherForV1(defaultLiquidityTierUpsertEventV1);
      expectKafkaMessages(producerSendMock, liquidityTierEvent, 0);
    });

    it('updates existing liquidity tier', async () => {
      const transactionIndex: number = 0;
      const liquidityTierEvent: LiquidityTierUpsertEventV1 = defaultLiquidityTierUpsertEventV1;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
        liquidityTierEvent: LiquidityTierUpsertEventV1.encode(
          liquidityTierEvent,
        ).finish(),
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 1,
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
      expectLiquidityTier(newLiquidityTiers[0], liquidityTierEvent, false);
      validateLiquidityTierRefresherForV1(defaultLiquidityTierUpsertEventV1);
      expectKafkaMessages(producerSendMock, liquidityTierEvent, 2);
    });

  });

  describe('liquidityTierHandlerV2', () => {
    describe('getParallelizationIds', () => {
      it('returns the correct parallelization ids', () => {
        const transactionIndex: number = 0;
        const eventIndex: number = 0;

        const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
          DydxIndexerSubtypes.LIQUIDITY_TIER,
          LiquidityTierUpsertEventV2.encode(defaultLiquidityTierUpsertEventV2).finish(),
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
          0,
          indexerTendermintEvent,
          0,
          defaultLiquidityTierUpsertEventV2,
        );

        expect(handler.getParallelizationIds()).toEqual([]);
      });
    });

    it('creates new liquidity tier', async () => {
      const transactionIndex: number = 0;
      const liquidityTierEvent: LiquidityTierUpsertEventV2 = defaultLiquidityTierUpsertEventV2;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
        liquidityTierEvent: LiquidityTierUpsertEventV2.encode(
          liquidityTierEvent,
        ).finish(),
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 2,
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
      validateLiquidityTierRefresherForV2(defaultLiquidityTierUpsertEventV2);
      expectKafkaMessages(producerSendMock, liquidityTierEvent, 0);
    });

    it('updates existing liquidity tier', async () => {
      const transactionIndex: number = 0;
      const liquidityTierEvent: LiquidityTierUpsertEventV2 = defaultLiquidityTierUpsertEventV2;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromLiquidityTiersEvent({
        liquidityTierEvent: LiquidityTierUpsertEventV2.encode(
          liquidityTierEvent,
        ).finish(),
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 2,
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
      validateLiquidityTierRefresherForV2(defaultLiquidityTierUpsertEventV2);
      expectKafkaMessages(producerSendMock, liquidityTierEvent, 2);
    });
  });
});

export function expectLiquidityTier(
  liquidityTierFromDb: LiquidityTiersFromDatabase,
  event: any,
  checkOICaps: boolean = true,
): void {
  expect(liquidityTierFromDb.id).toEqual(event.id);
  expect(liquidityTierFromDb.name).toEqual(event.name);
  expect(liquidityTierFromDb.initialMarginPpm).toEqual(event.initialMarginPpm.toString());
  expect(liquidityTierFromDb.maintenanceFractionPpm).toEqual(
    event.maintenanceFractionPpm.toString());
  if (checkOICaps) {
    expect(liquidityTierFromDb.openInterestLowerCap).toEqual(
      protocolTranslations.quantumsToHumanFixedString(
        event.openInterestLowerCap.toString(), QUOTE_CURRENCY_ATOMIC_RESOLUTION));
    expect(liquidityTierFromDb.openInterestUpperCap).toEqual(
      protocolTranslations.quantumsToHumanFixedString(
        event.openInterestUpperCap.toString(), QUOTE_CURRENCY_ATOMIC_RESOLUTION));
  }
}

function createKafkaMessageFromLiquidityTiersEvent({
  liquidityTierEvent,
  transactionIndex,
  height,
  time,
  txHash,
  version,
}: {
  liquidityTierEvent: Uint8Array,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
  version: number,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.LIQUIDITY_TIER,
      liquidityTierEvent,
      transactionIndex,
      0,
      version,
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

function validateLiquidityTierRefresherForV2(
  liquidityTierEvent: LiquidityTierUpsertEventV2,
) {
  const liquidityTier:
  LiquidityTiersFromDatabase = liquidityTierRefresher.getLiquidityTierFromId(
    liquidityTierEvent.id,
  );

  validateCommonLiquidityTierFields(liquidityTier, liquidityTierEvent);
  expect(liquidityTier.openInterestLowerCap).toEqual(
    protocolTranslations.quantumsToHumanFixedString(
      liquidityTierEvent.openInterestLowerCap.toString(), QUOTE_CURRENCY_ATOMIC_RESOLUTION));
  expect(liquidityTier.openInterestUpperCap).toEqual(
    protocolTranslations.quantumsToHumanFixedString(
      liquidityTierEvent.openInterestUpperCap.toString(), QUOTE_CURRENCY_ATOMIC_RESOLUTION));
}

function validateLiquidityTierRefresherForV1(
  liquidityTierEvent: LiquidityTierUpsertEventV1,
) {
  const liquidityTier:
  LiquidityTiersFromDatabase = liquidityTierRefresher.getLiquidityTierFromId(
    liquidityTierEvent.id,
  );

  validateCommonLiquidityTierFields(liquidityTier, liquidityTierEvent);
  expect(liquidityTier.openInterestLowerCap).toEqual(null);
  expect(liquidityTier.openInterestUpperCap).toEqual(null);
}

function validateCommonLiquidityTierFields(liquidityTier: LiquidityTiersFromDatabase,
  liquidityTierEvent: any) {
  expect(liquidityTier.id).toEqual(liquidityTierEvent.id);
  expect(liquidityTier.name).toEqual(liquidityTierEvent.name);
  expect(liquidityTier.initialMarginPpm).toEqual(
    liquidityTierEvent.initialMarginPpm.toString());
  expect(liquidityTier.maintenanceFractionPpm).toEqual(
    liquidityTierEvent.maintenanceFractionPpm.toString());
}

function expectKafkaMessages(
  producerSendMock: jest.SpyInstance,
  liquidityTier: any,
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
