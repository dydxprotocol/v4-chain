import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  dbHelpers, MarketFromDatabase, MarketTable, testMocks,
} from '@dydxprotocol-indexer/postgres';
import { IndexerTendermintBlock, IndexerTendermintEvent, MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes, MarketCreateEventMessage } from '../../../src/lib/types';
import {
  defaultHeight,
  defaultMarketCreate,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromMarketEvent } from '../../helpers/kafka-helpers';
import { producer } from '@dydxprotocol-indexer/kafka';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { MarketCreateHandler } from '../../../src/handlers/markets/market-create-handler';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../../helpers/indexer-proto-helpers';
import Long from 'long';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';

describe('marketCreateHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
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

  const loggerCrit = jest.spyOn(logger, 'crit');
  const loggerError = jest.spyOn(logger, 'error');
  const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const marketEvent: MarketEventV1 = {
        marketId: 0,
        priceUpdate: {
          priceWithExponent: Long.fromValue(1, true),
        },
      };
      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.MARKET,
        MarketEventV1.encode(marketEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: MarketCreateHandler = new MarketCreateHandler(
        block,
        indexerTendermintEvent,
        0,
        marketEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_0`,
      ]);
    });
  });

  it('creates new market', async () => {
    const transactionIndex: number = 0;

    const marketCreate: MarketEventV1 = {
      marketId: 3,
      marketCreate: {
        base: {
          pair: 'DYDX-USD',
          minPriceChangePpm: 500,
        },
        exponent: -5,
      },
    };

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [marketCreate],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const market: MarketFromDatabase = await MarketTable.findById(
      marketCreate.marketId,
    ) as MarketFromDatabase;

    expectMarketMatchesEvent(marketCreate as MarketCreateEventMessage, market);
  });

  it('errors when attempting to create an existing market', async () => {
    const transactionIndex: number = 0;

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [defaultMarketCreate],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      new ParseMessageError('Market in MarketCreate already exists'),
    );

    // Check that market in database is the old market.
    const market: MarketFromDatabase = await MarketTable.findById(
      defaultMarketCreate.marketId,
    ) as MarketFromDatabase;
    expect(market.minPriceChangePpm).toEqual(50);

    expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
      at: 'MarketCreateHandler#logAndThrowParseMessageError',
      message: 'Market in MarketCreate already exists',
    }));
    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
    }));
    expect(producerSendMock.mock.calls.length).toEqual(0);
  });
});

function expectMarketMatchesEvent(
  event: MarketCreateEventMessage,
  market: MarketFromDatabase,
) {
  expect(market.id).toEqual(event.marketId);
  expect(market.pair).toEqual(event.marketCreate.base!.pair!);
  expect(market.minPriceChangePpm).toEqual(event.marketCreate.base!.minPriceChangePpm!);
  expect(market.exponent).toEqual(event.marketCreate.exponent);
}
