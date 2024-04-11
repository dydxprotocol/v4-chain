import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  OpenInterestUpdateEventV1,
  Timestamp,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers,
  PerpetualMarketFromDatabase,
  perpetualMarketRefresher,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectMarketKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import {
  defaultHeight,
  defaultOpenInterestUpdateEvent,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import _ from 'lodash';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import {
  bytesToBigInt,
} from '@dydxprotocol-indexer/v4-proto-parser';

describe('openInterestUpdateHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    updateBlockCache(defaultPreviousHeight);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
    perpetualMarketRefresher.clear();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  it('update open interest', async () => {
    const transactionIndex: number = 0;
    const openInterestUpdateEvent: OpenInterestUpdateEventV1 = defaultOpenInterestUpdateEvent;
    const kafkaMessage: KafkaMessage = createKafkaEventForOpenInterestUpdateEvent({
      openInterestUpdateEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });
    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);
    await perpetualMarketRefresher.updatePerpetualMarkets();

    const perpetualMarketsFromDB: PerpetualMarketFromDatabase[] = [];
    for (const openInterestUpdate
      of defaultOpenInterestUpdateEvent.openInterestUpdates) {
      const perpetualMarket:
      PerpetualMarketFromDatabase = perpetualMarketRefresher.getPerpetualMarketFromId(
        openInterestUpdate.perpetualId.toString())!;
      expect(Number(perpetualMarket.openInterest)).toEqual(
        Number(bytesToBigInt(openInterestUpdate.openInterest)));
      perpetualMarketsFromDB.push(perpetualMarket);
    }
    expectMarketKafkaMessage({
      producerSendMock,
      contents:
        JSON.stringify({
          trading:
          _.chain(perpetualMarketsFromDB)
            .keyBy('ticker')
            .mapValues((perpetualMarket) => {
              return {
                id: perpetualMarket.id,
                openInterest: perpetualMarket.openInterest,
              };
            })
            .value(),
        }),
    });
  });
});

function createKafkaEventForOpenInterestUpdateEvent({
  openInterestUpdateEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  openInterestUpdateEvent: OpenInterestUpdateEventV1,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  events.push(
    createIndexerTendermintEvent(
      DydxIndexerSubtypes.OPEN_INTEREST_UPDATE,
      OpenInterestUpdateEventV1.encode(openInterestUpdateEvent).finish(),
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
