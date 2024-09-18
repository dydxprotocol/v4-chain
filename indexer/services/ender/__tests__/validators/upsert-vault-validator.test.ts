import { ParseMessageError, logger } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  UpsertVaultEventV1,
  VaultStatus,
} from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, testConstants, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { defaultHeight, defaultTime, defaultTxHash } from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';
import { UpsertVaultValidator } from '../../src/validators/upsert-vault-validator';

describe('upsert-vault-validator', () => {
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
    it('does not throw error on valid uspert vault event', () => {
      const event: UpsertVaultEventV1 = {
        address: testConstants.defaultVaultAddress,
        clobPairId: 0,
        status: VaultStatus.VAULT_STATUS_QUOTING,
      };
      const validator: UpsertVaultValidator = new UpsertVaultValidator(
        event,
        createBlock(event),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error if address in event is empty', () => {
      const event: UpsertVaultEventV1 = {
        address: '',
        clobPairId: 0,
        status: VaultStatus.VAULT_STATUS_QUOTING,
      };
      const validator: UpsertVaultValidator = new UpsertVaultValidator(
        event,
        createBlock(event),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(
        'UpsertVaultEvent address is not populated',
      ));
    });
  });
});

function createBlock(
  upsertVaultEvent: UpsertVaultEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.UPSERT_VAULT,
    UpsertVaultEventV1.encode(upsertVaultEvent).finish(),
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
