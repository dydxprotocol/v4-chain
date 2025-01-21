import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  UpdatePerpetualEventV1,
  UpdatePerpetualEventV2,
  UpdatePerpetualEventV3,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultUpdatePerpetualEventV1,
  defaultUpdatePerpetualEventV2,
  defaultUpdatePerpetualEventV3,
  defaultHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { UpdatePerpetualValidator } from '../../src/validators/update-perpetual-validator';

describe('update-perpetual-validator', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await perpetualMarketRefresher.clear();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe.each([
    [
      'UpdatePerpetualEventV1',
      UpdatePerpetualEventV1.encode(defaultUpdatePerpetualEventV1).finish(),
      defaultUpdatePerpetualEventV1,
    ],
    [
      'UpdatePerpetualEventV2',
      UpdatePerpetualEventV2.encode(defaultUpdatePerpetualEventV2).finish(),
      defaultUpdatePerpetualEventV2,
    ],
    [
      'UpdatePerpetualEventV3',
      UpdatePerpetualEventV3.encode(defaultUpdatePerpetualEventV3).finish(),
      defaultUpdatePerpetualEventV3,
    ],
  ])('%s', (
    _name: string,
    updatePerpetualEventBytes: Uint8Array,
    event: UpdatePerpetualEventV1 | UpdatePerpetualEventV2 | UpdatePerpetualEventV3,
  ) => {
    it('does not throw error on valid perpetual market create event', () => {
      const validator: UpdatePerpetualValidator = new UpdatePerpetualValidator(
        event,
        createBlock(updatePerpetualEventBytes),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error if id does not correspond to an existing perpetual market', () => {
      const validator: UpdatePerpetualValidator = new UpdatePerpetualValidator(
        {
          ...event,
          id: 20,
        },
        createBlock(updatePerpetualEventBytes),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(
        'UpdatePerpetualEvent.id must correspond with an existing perpetual_market.id',
      ));
    });
  });
});

function createBlock(
  updatePerpetualEventBytes: Uint8Array,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.UPDATE_PERPETUAL,
    updatePerpetualEventBytes,
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
