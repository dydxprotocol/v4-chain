import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  UpdatePerpetualEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultUpdatePerpetualEvent,
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

  describe('validate', () => {
    it('does not throw error on valid perpetual market create event', () => {
      const validator: UpdatePerpetualValidator = new UpdatePerpetualValidator(
        defaultUpdatePerpetualEvent,
        createBlock(defaultUpdatePerpetualEvent),
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error if id does not correspond to an existing perpetual market', () => {
      const validator: UpdatePerpetualValidator = new UpdatePerpetualValidator(
        {
          ...defaultUpdatePerpetualEvent,
          id: 20,
        },
        createBlock(defaultUpdatePerpetualEvent),
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(
        'UpdatePerpetualEvent.id must correspond with an existing perpetual_market.id',
      ));
    });
  });
});

function createBlock(
  updatePerpetualEvent: UpdatePerpetualEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.UPDATE_PERPETUAL,
    UpdatePerpetualEventV1.encode(updatePerpetualEvent).finish(),
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
