import { logger } from '@dydxprotocol-indexer/base';
import {
  dbHelpers, MarketFromDatabase, MarketTable, testMocks,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes, MarketModifyEventMessage } from '../../../src/lib/types';
import {
  defaultHeight, defaultMarketModify, defaultPreviousHeight, defaultTime, defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromMarketEvent } from '../../helpers/kafka-helpers';
import { producer } from '@dydxprotocol-indexer/kafka';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { MarketEventV1, IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../../helpers/indexer-proto-helpers';
import { MarketModifyHandler } from '../../../src/handlers/markets/market-modify-handler';
import Long from 'long';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';

describe('marketModifyHandler', () => {

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

      const handler: MarketModifyHandler = new MarketModifyHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        marketEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([
        `${handler.eventType}_0`,
      ]);
    });
  });

  it('modifies existing market', async () => {
    const transactionIndex: number = 0;

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [defaultMarketModify],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const market: MarketFromDatabase = await MarketTable.findById(
      defaultMarketModify.marketId,
    ) as MarketFromDatabase;

    expectMarketMatchesEvent(defaultMarketModify as MarketModifyEventMessage, market);
  });

  it('modifies non-existent market', async () => {
    const transactionIndex: number = 0;

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [{
        ...defaultMarketModify,
        marketId: 5,
      }],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      'Market in MarketModify doesn\'t exist',
    );
    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: expect.stringContaining('PL/pgSQL function dydx_market_modify_handler('),
      message: expect.stringContaining('Market in MarketModify doesn\'t exist'),
    }));
    expect(producerSendMock.mock.calls.length).toEqual(0);
  });
});

function expectMarketMatchesEvent(
  event: MarketModifyEventMessage,
  market: MarketFromDatabase,
) {
  expect(market.id).toEqual(event.marketId);
  expect(market.pair).toEqual(event.marketModify.base!.pair!);
  expect(market.minPriceChangePpm).toEqual(event.marketModify.base!.minPriceChangePpm!);
}
