import { logger, stats } from '@dydxprotocol-indexer/base';
import {
  FundingEventV1,
  FundingEventV1_Type,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  assetRefresher,
  dbHelpers,
  FundingIndexUpdatesColumns,
  FundingIndexUpdatesFromDatabase,
  FundingIndexUpdatesTable,
  PerpetualMarketTable,
  OraclePriceTable,
  Ordering,
  perpetualMarketRefresher,
  protocolTranslations,
  TendermintEventTable,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes, FundingEventMessage } from '../../src/lib/types';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { FundingHandler } from '../../src/handlers/funding-handler';
import {
  defaultFundingRateEvent,
  defaultFundingUpdateSampleEvent,
  defaultFundingUpdateSampleEventWithAdditionalMarket,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { indexerTendermintEventToTransactionIndex } from '../../src/lib/helper';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { expectNextFundingRate } from '../helpers/redis-helpers';
import { redis } from '@dydxprotocol-indexer/redis';
import Big from 'big.js';
import { redisClient } from '../../src/helpers/redis/redis-controller';
import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { defaultPerpetualMarket } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('fundingHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await Promise.all([
      OraclePriceTable.create(testConstants.defaultOraclePrice),
      OraclePriceTable.create(testConstants.defaultOraclePrice2),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();
    await assetRefresher.updateAssets();
    updateBlockCache(defaultPreviousHeight);

    jest.spyOn(logger, 'error');
    jest.resetAllMocks();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    await redis.deleteAllAsync(redisClient);
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
        DydxIndexerSubtypes.FUNDING,
        FundingEventV1.encode(defaultFundingUpdateSampleEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [],
      );

      const handler: FundingHandler = new FundingHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultFundingUpdateSampleEvent,
      );

      const id: string = FundingIndexUpdatesTable.uuid(
        block.txHashes[transactionIndex],
        TendermintEventTable.createEventId(
          block.height.toString(),
          indexerTendermintEventToTransactionIndex(indexerTendermintEvent),
          indexerTendermintEvent.eventIndex,
        ),
        defaultFundingUpdateSampleEvent.updates[0].perpetualId.toString(),
      );
      const expectedParallelizationId: string = `FundingEvent_${id}`;
      expect(handler.getParallelizationIds()).toEqual([expectedParallelizationId]);
    });
  });

  it('successfully processes single premium sample event', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [defaultFundingUpdateSampleEvent],
      height: defaultHeight,
      time: defaultTime,
    });

    await onMessage(kafkaMessage);

    await expectNextFundingRate(
      new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
        defaultFundingUpdateSampleEvent.updates[0].fundingValuePpm,
      )),
      'BTC-USD',
    );
  });

  it.each([
    [
      'Non-zero sample',
      'TEST-USD',
      '0.0001',
      120,
      new Big('0.000115'), // 0.000120 / 8 + default 0.0001
    ],
    [
      'Sample is zero',
      'TEST-USD',
      '0.0001',
      0,
      new Big('0.0001'), // 0 + default 0.0001
    ],
  ])('(%s) Non-zero default funding: successfully handle premium sample', async (
    _name: string,
    ticker: string,
    defaultFundingRate1H: string,
    fundingValuePpm: number,
    expectedNextFundingRate: Big,
  ) => {
    const testPerpetualMarket = await PerpetualMarketTable.create({
      ...defaultPerpetualMarket,
      id: '1000', // Different id than `defaultPerpeptualMarket` to avoid conflict
      ticker,
      defaultFundingRate1H,
    });

    const fundingUpdateSampleEvent: FundingEventV1 = {
      type: FundingEventV1_Type.TYPE_PREMIUM_SAMPLE,
      updates: [
        {
          perpetualId: parseInt(testPerpetualMarket.id, 10),
          fundingValuePpm,
          fundingIndex: bigIntToBytes(BigInt(0)),
        },
      ],
    };

    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [fundingUpdateSampleEvent],
      height: defaultHeight,
      time: defaultTime,
    });

    await onMessage(kafkaMessage);

    await expectNextFundingRate(
      expectedNextFundingRate,
      ticker,
      defaultFundingRate1H,
    );
  });

  it('successfully processes multiple premium sample event for different markets', async () => {
    const fundingUpdateSampleEvent2: FundingEventV1 = {
      type: FundingEventV1_Type.TYPE_PREMIUM_SAMPLE,
      updates: [
        {
          perpetualId: 0,
          fundingValuePpm: 100,
          fundingIndex: bigIntToBytes(BigInt(0)),
        },
        {
          perpetualId: 1,
          fundingValuePpm: 50,
          fundingIndex: bigIntToBytes(BigInt(0)),
        },
      ],
    };

    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [defaultFundingUpdateSampleEvent, fundingUpdateSampleEvent2],
      height: defaultHeight,
      time: defaultTime,
    });

    await onMessage(kafkaMessage);

    await expectNextFundingRate(
      new Big('0.000006875'),
      'BTC-USD',
    );
    await expectNextFundingRate(
      new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
        fundingUpdateSampleEvent2.updates[1].fundingValuePpm,
      )),
      'ETH-USD',
    );
  });

  it('successfully ignores funding rate and index for market with no oracle price', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [{
        ...defaultFundingRateEvent,
        updates: [
          {
            perpetualId: 2,
            fundingValuePpm: 10,
            fundingIndex: bigIntToBytes(BigInt(0)),
          },
        ],
      }],
      height: defaultHeight,
      time: defaultTime,
    });
    await onMessage(kafkaMessage);
    expect(logger.error).toHaveBeenCalledWith(expect.objectContaining({
      at: 'FundingHandler#handleFundingSample',
      message: 'oracle_price not found for marketId.',
    }));
  });

  it('successfully processes and clears cache for a new funding rate with both existing/non-existent market',
    async () => {
      const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
        fundingEvents: [defaultFundingUpdateSampleEventWithAdditionalMarket],
        height: defaultHeight,
        time: defaultTime,
      });

      await onMessage(kafkaMessage);

      await expectNextFundingRate(
        new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
          defaultFundingUpdateSampleEvent.updates[0].fundingValuePpm,
        )),
        'BTC-USD',
      );

      const kafkaMessage2: KafkaMessage = createKafkaMessageFromFundingEvents({
        fundingEvents: [defaultFundingRateEvent],
        height: 4,
        time: defaultTime,
      });

      await onMessage(kafkaMessage2);
      await expectNextFundingRate(
        undefined,
        'BTC-USD',
      );
      const fundingIndices: FundingIndexUpdatesFromDatabase[] = await
      FundingIndexUpdatesTable.findAll({}, [], {});

      expect(fundingIndices.length).toEqual(1);
      expect(fundingIndices[0]).toEqual(expect.objectContaining({
        perpetualId: '0',
        rate: '0.00000125',
        oraclePrice: '10000',
        fundingIndex: '0.1',
      }));
      expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update_event', 0.1, { ticker: 'BTC-USD' });
      expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update', 0.1, { ticker: 'BTC-USD' });
      expect(stats.timing).toHaveBeenCalledWith(
        'ender.handle_funding_event.sql_latency',
        expect.any(Number),
        {
          className: 'FundingHandler',
          eventType: 'FundingEvent',
        },
      );
    });

  it('successfully processes and clears cache for a new funding rate', async () => {
    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [defaultFundingUpdateSampleEvent],
      height: defaultHeight,
      time: defaultTime,
    });

    await onMessage(kafkaMessage);

    await expectNextFundingRate(
      new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
        defaultFundingUpdateSampleEvent.updates[0].fundingValuePpm,
      )),
      'BTC-USD',
    );

    const kafkaMessage2: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [defaultFundingRateEvent],
      height: 4,
      time: defaultTime,
    });

    await onMessage(kafkaMessage2);
    await expectNextFundingRate(
      undefined,
      'BTC-USD',
    );
    const fundingIndices: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll({}, [], {});

    expect(fundingIndices.length).toEqual(1);
    expect(fundingIndices[0]).toEqual(expect.objectContaining({
      perpetualId: '0',
      rate: '0.00000125',
      oraclePrice: '10000',
      fundingIndex: '0.1',
    }));
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update_event', 0.1, { ticker: 'BTC-USD' });
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update', 0.1, { ticker: 'BTC-USD' });
  });

  it('successfully processes and clears cache for multiple new funding rates', async () => {
    const fundingSampleEvent: FundingEventV1 = {
      type: FundingEventV1_Type.TYPE_PREMIUM_SAMPLE,
      updates: [
        {
          perpetualId: 0,
          fundingValuePpm: 100,
          fundingIndex: bigIntToBytes(BigInt(0)),
        },
        {
          perpetualId: 1,
          fundingValuePpm: 50,
          fundingIndex: bigIntToBytes(BigInt(0)),
        },
      ],
    };
    const kafkaMessage: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [fundingSampleEvent],
      height: defaultHeight,
      time: defaultTime,
    });

    await onMessage(kafkaMessage);

    await Promise.all([
      expectNextFundingRate(
        new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
          fundingSampleEvent.updates[0].fundingValuePpm,
        )),
        'BTC-USD',
      ),
      expectNextFundingRate(
        new Big(protocolTranslations.funding8HourValuePpmTo1HourRate(
          fundingSampleEvent.updates[1].fundingValuePpm,
        )),
        'ETH-USD',
      ),
    ]);

    const fundingRateEvent: FundingEventMessage = {
      type: FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX,
      updates: [
        {
          perpetualId: 0,
          fundingValuePpm: 10,
          fundingIndex: bigIntToBytes(BigInt(10)),
        },
        {
          perpetualId: 1,
          fundingValuePpm: 100,
          fundingIndex: bigIntToBytes(BigInt(100)),
        },
      ],
    };
    const kafkaMessage2: KafkaMessage = createKafkaMessageFromFundingEvents({
      fundingEvents: [fundingRateEvent],
      height: 4,
      time: defaultTime,
    });

    await onMessage(kafkaMessage2);
    await Promise.all([
      expectNextFundingRate(
        undefined,
        'BTC-USD',
      ),
      expectNextFundingRate(
        undefined,
        'ETH-USD',
      ),
    ]);
    const fundingIndices: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {},
      [],
      {
        orderBy: [[FundingIndexUpdatesColumns.perpetualId, Ordering.ASC]],
      },
    );

    expect(fundingIndices.length).toEqual(2);
    expect(fundingIndices[0]).toEqual(expect.objectContaining({
      perpetualId: '0',
      rate: '0.00000125',
      oraclePrice: '10000',
      // 1e1 * 1e-6 * 1e-6 / 1e-10 = 1e-1
      fundingIndex: '0.1',
    }));
    expect(fundingIndices[1]).toEqual(expect.objectContaining({
      perpetualId: '1',
      rate: '0.0000125',
      oraclePrice: '500',
      // 1e2 * 1e-6 * 1e-6 / 1e-18 = 1e8
      fundingIndex: '100000000',
    }));
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update_event', 0.1, { ticker: 'BTC-USD' });
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update', 0.1, { ticker: 'BTC-USD' });
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update_event', 100000000, { ticker: 'ETH-USD' });
    expect(stats.gauge).toHaveBeenCalledWith('ender.funding_index_update', 100000000, { ticker: 'ETH-USD' });
  });
});

function createKafkaMessageFromFundingEvents({
  fundingEvents,
  height,
  time,
}: {
  fundingEvents: FundingEventV1[],
  height: number,
  time: Timestamp,
}) {
  const events: IndexerTendermintEvent[] = [];
  let eventIndex: number = 0;
  const transactionIndex: number = -1;
  for (const fundingEvent of fundingEvents) {
    events.push(
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.FUNDING,
        FundingEventV1.encode(fundingEvent).finish(),
        transactionIndex,
        eventIndex,
      ),
    );
    eventIndex += 1;
  }

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}
