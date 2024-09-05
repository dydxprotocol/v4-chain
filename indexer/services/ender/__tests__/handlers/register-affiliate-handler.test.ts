import {
  RegisterAffiliateEventV1,
  IndexerTendermintBlock,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers,
  testMocks,
  AffiliateReferredUsersTable,
  AffiliateReferredUserFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import {
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

describe('registerAffiliateHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await dbHelpers.clearData();
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

  it('should register affiliates in single block', async () => {
    const events: RegisterAffiliateEventV1[] = [
      {
        affiliate: 'address1',
        referee: 'address2',
      },
      {
        affiliate: 'address3',
        referee: 'address4',
      },
    ];
    const block: IndexerTendermintBlock = createBlockFromEvents(
      defaultHeight,
      ...events,
    );
    const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
    const kafkaMessage: KafkaMessage = createKafkaMessage(Buffer.from(binaryBlock));

    await onMessage(kafkaMessage);

    let actualEntry: AffiliateReferredUserFromDatabase | undefined = await AffiliateReferredUsersTable.findByRefereeAddress('address2');
    expect(actualEntry).toEqual({
      affiliateAddress: 'address1',
      refereeAddress: 'address2',
      referredAtBlock: defaultHeight.toString(),
    });

    actualEntry = await AffiliateReferredUsersTable.findByRefereeAddress('address4');
    expect(actualEntry).toEqual({
      affiliateAddress: 'address3',
      refereeAddress: 'address4',
      referredAtBlock: defaultHeight.toString(),
    });
  });
});

function createBlockFromEvents(
  height: number,
  ...events: RegisterAffiliateEventV1[]
): IndexerTendermintBlock {
  const transactionIndex = 0;
  let eventIndex = 0;

  return createIndexerTendermintBlock(
    height,
    defaultTime,
    events.map((event) => {
      const indexerEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.REGISTER_AFFILIATE,
        RegisterAffiliateEventV1.encode(event).finish(),
        transactionIndex,
        eventIndex,
      );
      eventIndex += 1;
      return indexerEvent;
    }),
    [defaultTxHash],
  );
}
