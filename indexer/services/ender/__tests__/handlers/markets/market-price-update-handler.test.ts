import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  dbHelpers,
  MarketFromDatabase,
  MarketMessageContents,
  MarketTable,
  OraclePriceFromDatabase,
  OraclePriceTable,
  protocolTranslations,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { onMessage } from '../../../src/lib/on-message';
import { DydxIndexerSubtypes, MarketPriceUpdateEventMessage } from '../../../src/lib/types';
import {
  defaultHeight,
  defaultMarketPriceUpdate,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../../helpers/constants';
import { createKafkaMessageFromMarketEvent } from '../../helpers/kafka-helpers';
import { producer } from '@dydxprotocol-indexer/kafka';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectMarketKafkaMessage,
} from '../../helpers/indexer-proto-helpers';
import { generateOraclePriceContents } from '../../../src/helpers/kafka-helper';
import { updateBlockCache } from '../../../src/caches/block-cache';
import { MarketEventV1, IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { MarketPriceUpdateHandler } from '../../../src/handlers/markets/market-price-update-handler';
import Long from 'long';
import { getPrice } from '../../../src/caches/price-cache';
import { createPostgresFunctions } from '../../../src/helpers/postgres/postgres-functions';

describe('marketPriceUpdateHandler', () => {
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

      const handler: MarketPriceUpdateHandler = new MarketPriceUpdateHandler(
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

  it('fails when no market exists', async () => {
    const transactionIndex: number = 0;
    const marketPriceUpdate: MarketEventV1 = {
      marketId: 5,
      priceUpdate: {
        priceWithExponent: Long.fromValue(50000000, true),
      },
    };
    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [marketPriceUpdate],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      new ParseMessageError('MarketPriceUpdateEvent contains a non-existent market id'),
    );

    expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
      at: 'MarketPriceUpdateHandler#logAndThrowParseMessageError',
      message: 'MarketPriceUpdateEvent contains a non-existent market id',
    }));
    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
    }));
    expect(producerSendMock.mock.calls.length).toEqual(0);
  });

  it('successfully inserts new oracle price for existing market', async () => {
    const transactionIndex: number = 0;

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [defaultMarketPriceUpdate],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const { market, oraclePrice } = await getDbState(defaultMarketPriceUpdate);

    expectOraclePriceMatchesEvent(
      defaultMarketPriceUpdate as MarketPriceUpdateEventMessage,
      oraclePrice,
      market,
      defaultHeight,
    );

    expect(getPrice(oraclePrice.marketId)).toEqual(oraclePrice.price);

    const contents: MarketMessageContents = generateOraclePriceContents(
      oraclePrice,
      market.pair,
    );

    expectMarketKafkaMessage({
      producerSendMock,
      contents: JSON.stringify(contents),
    });
  });

  it('successfully inserts new oracle price for market created in same block', async () => {
    const transactionIndex: number = 0;
    const newMarketId: number = 3000;

    // Include an event to create the market
    const marketCreate: MarketEventV1 = {
      marketId: newMarketId,
      marketCreate: {
        base: {
          pair: 'NEWTOKEN-USD',
          minPriceChangePpm: 500,
        },
        exponent: -5,
      },
    };
    const marketPriceUpdate: MarketEventV1 = {
      marketId: newMarketId,
      priceUpdate: {
        priceWithExponent: Long.fromValue(50000000),
      },
    };

    const kafkaMessage: KafkaMessage = createKafkaMessageFromMarketEvent({
      marketEvents: [marketCreate, marketPriceUpdate],
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await onMessage(kafkaMessage);

    const { market, oraclePrice } = await getDbState(marketPriceUpdate);

    expectOraclePriceMatchesEvent(
      marketPriceUpdate as MarketPriceUpdateEventMessage,
      oraclePrice,
      market,
      defaultHeight,
    );

    expect(getPrice(oraclePrice.marketId)).toEqual(oraclePrice.price);

    const contents: MarketMessageContents = generateOraclePriceContents(
      oraclePrice,
      market.pair,
    );

    expectMarketKafkaMessage({
      producerSendMock,
      contents: JSON.stringify(contents),
    });
  });
});

async function getDbState(marketPriceUpdate: MarketEventV1): Promise<any> {
  const [market, oraclePrice]:
  [MarketFromDatabase, OraclePriceFromDatabase] = await Promise.all([
    MarketTable
      .findById(
        marketPriceUpdate.marketId,
      ) as Promise<MarketFromDatabase>,
    OraclePriceTable.findMostRecentMarketOraclePrice(
      marketPriceUpdate.marketId,
    ) as Promise<OraclePriceFromDatabase>,
  ]);

  return { market, oraclePrice };
}

function expectOraclePriceMatchesEvent(
  event: MarketPriceUpdateEventMessage,
  oraclePrice: OraclePriceFromDatabase,
  market: MarketFromDatabase,
  height: number,
) {
  const expectedHumanPrice: string = protocolTranslations.protocolPriceToHuman(
    event.priceUpdate.priceWithExponent.toString(),
    market!.exponent,
  );
  expect(market.id).toEqual(event.marketId);
  expect(market.oraclePrice).toEqual(expectedHumanPrice);

  expect(oraclePrice.marketId).toEqual(event.marketId);
  expect(oraclePrice.price).toEqual(expectedHumanPrice);
  expect(oraclePrice.effectiveAtHeight).toEqual(height.toString());
}
