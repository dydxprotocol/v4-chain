import { logger } from '@dydxprotocol-indexer/base';
import { IndexerTendermintBlock, IndexerTendermintEvent, RegisterAffiliateEventV1 } from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { defaultHeight, defaultTime, defaultTxHash } from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { RegisterAffiliateValidator } from '../../src/validators/register-affiliate-validator';

describe('register-affiliate-validator', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid register affiliate event', () => {
      const event: RegisterAffiliateEventV1 = { affiliate: 'address1', referee: 'address2' };
      const validator: RegisterAffiliateValidator = new RegisterAffiliateValidator(
        event,
        createBlock(event),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });
  });
});

function createBlock(
  RegisterAffiliateEventevent: RegisterAffiliateEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.REGISTER_AFFILIATE,
    RegisterAffiliateEventV1.encode(RegisterAffiliateEventevent).finish(),
    0,
    0,
  );

  return createIndexerTendermintBlock(
    defaultHeight,
    defaultTime,
    [event],
    [defaultTxHash],
  );
}
