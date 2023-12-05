import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  UpdateClobPairEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultUpdateClobPairEvent,
  defaultHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { UpdateClobPairValidator } from '../../src/validators/update-clob-pair-validator';

describe('update-clob-pair-validator', () => {
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
      const validator: UpdateClobPairValidator = new UpdateClobPairValidator(
        defaultUpdateClobPairEvent,
        createBlock(defaultUpdateClobPairEvent),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error if id does not correspond to an existing perpetual market', () => {
      const validator: UpdateClobPairValidator = new UpdateClobPairValidator(
        {
          ...defaultUpdateClobPairEvent,
          clobPairId: 20,
        },
        createBlock(defaultUpdateClobPairEvent),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(
        'UpdateClobPairEvent.clobPairId must correspond with an existing perpetual_market.clobPairId',
      ));
    });
  });
});

function createBlock(
  updateClobPairEvent: UpdateClobPairEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.UPDATE_CLOB_PAIR,
    UpdateClobPairEventV1.encode(updateClobPairEvent).finish(),
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
