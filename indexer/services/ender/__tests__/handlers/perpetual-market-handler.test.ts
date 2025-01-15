import {
  PerpetualMarketCreateEventV1,
  PerpetualMarketCreateEventV2,
  PerpetualMarketCreateEventV3,
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
  expectPerpetualMarketV1,
  expectPerpetualMarketKafkaMessage,
  expectPerpetualMarketV2,
  expectPerpetualMarketV3,
} from '../helpers/indexer-proto-helpers';
import { PerpetualMarketCreationHandler } from '../../src/handlers/perpetual-market-handler';
import {
  defaultPerpetualMarketCreateEventV1,
  defaultPerpetualMarketCreateEventV2,
  defaultPerpetualMarketCreateEventV3,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { expectPerpetualMarketMatchesEvent } from '../helpers/postgres-helpers';

describe('perpetualMarketHandler', () => {
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
    perpetualMarketRefresher.clear();
    jest.clearAllMocks();
    liquidityTierRefresher.clear();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  describe.each([
    [
      'PerpetualMarketCreateEventV1',
      1,
      PerpetualMarketCreateEventV1.encode(defaultPerpetualMarketCreateEventV1).finish(),
      expectPerpetualMarketV1,
      defaultPerpetualMarketCreateEventV1,
    ],
    [
      'PerpetualMarketCreateEventV2',
      2,
      PerpetualMarketCreateEventV2.encode(defaultPerpetualMarketCreateEventV2).finish(),
      expectPerpetualMarketV2,
      defaultPerpetualMarketCreateEventV2,
    ],
    [
      'PerpetualMarketCreateEventV3',
      3,
      PerpetualMarketCreateEventV3.encode(defaultPerpetualMarketCreateEventV3).finish(),
      expectPerpetualMarketV3,
      defaultPerpetualMarketCreateEventV3,
    ],
  ])('%s', (
    _name: string,
    version: number,
    perpetualMarketCreateEventBytes: Uint8Array,
    expectPerpetualMarket: Function,
    event: PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
    | PerpetualMarketCreateEventV3,
  ) => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.PERPETUAL_MARKET,
        perpetualMarketCreateEventBytes,
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
        0,
        indexerTendermintEvent,
        0,
        event,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });

    it('fails when market doesnt exist for perpetual market', async () => {
      const transactionIndex: number = 0;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
        perpetualMarketEventBytes: perpetualMarketCreateEventBytes,
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 1,
      });

      await expect(onMessage(kafkaMessage)).rejects.toThrowError();
    });

    it('fails when liquidity tier doesnt exist for perpetual market', async () => {
      await MarketTable.create(testConstants.defaultMarket);
      const transactionIndex: number = 0;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
        perpetualMarketEventBytes: perpetualMarketCreateEventBytes,
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version: 1,
      });

      await expect(onMessage(kafkaMessage)).rejects.toThrowError();
    });

    it('creates new perpetual market with the event', async () => {
      await Promise.all([
        MarketTable.create(testConstants.defaultMarket),
        LiquidityTiersTable.create(testConstants.defaultLiquidityTier),
      ]);
      await liquidityTierRefresher.updateLiquidityTiers();

      const transactionIndex: number = 0;
      const kafkaMessage: KafkaMessage = createKafkaMessageFromPerpetualMarketEvent({
        perpetualMarketEventBytes: perpetualMarketCreateEventBytes,
        transactionIndex,
        height: defaultHeight,
        time: defaultTime,
        txHash: defaultTxHash,
        version,
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
      expectPerpetualMarketMatchesEvent(event, newPerpetualMarkets[0]);
      const perpetualMarket: PerpetualMarketFromDatabase | undefined = perpetualMarketRefresher.getPerpetualMarketFromId('0');
      expect(perpetualMarket).toBeDefined();
      expectPerpetualMarket(perpetualMarket!, event);
      expectPerpetualMarketKafkaMessage(producerSendMock, [perpetualMarket!]);
    });
  });
});

function createKafkaMessageFromPerpetualMarketEvent({
  perpetualMarketEventBytes,
  transactionIndex,
  height,
  time,
  txHash,
  version,
}: {
  perpetualMarketEventBytes: Uint8Array,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
  version: number,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.PERPETUAL_MARKET,
      perpetualMarketEventBytes,
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

async function expectNoExistingPerpetualMarkets() {
  // Confirm there is no existing perpetual markets
  const perpetualMarkets: PerpetualMarketFromDatabase[] = await PerpetualMarketTable.findAll(
    {},
    [], {
      orderBy: [[PerpetualMarketColumns.id, Ordering.ASC]],
    });

  expect(perpetualMarkets.length).toEqual(0);
}
