import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  SubaccountUpdateEventV1,
} from '@dydxprotocol-indexer/v4-protos';

import { DydxIndexerSubtypes } from '../../src/lib/types';
import { SubaccountUpdateValidator } from '../../src/validators/subaccount-update-validator';
import {
  defaultEmptySubaccountUpdateEvent, defaultHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';

describe('subaccount-update-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid subaccountupdate', () => {
      const validator: SubaccountUpdateValidator = new SubaccountUpdateValidator(
        defaultEmptySubaccountUpdateEvent,
        createBlock(defaultEmptySubaccountUpdateEvent),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error if event does not contain subaccountId', () => {
      const invalidEvent: SubaccountUpdateEventV1 = {
        updatedAssetPositions: [],
        updatedPerpetualPositions: [],
      };
      const validator: SubaccountUpdateValidator = new SubaccountUpdateValidator(
        invalidEvent,
        createBlock(invalidEvent),
        0,
      );

      const message: string = 'SubaccountUpdateEvent must contain a subaccountId';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        SubaccountUpdateValidator.name,
        message,
        { event: invalidEvent },
      );
    });
  });
});

function createBlock(
  subaccountUpdateEvent: SubaccountUpdateEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.SUBACCOUNT_UPDATE,
    SubaccountUpdateEventV1.encode(subaccountUpdateEvent).finish(),
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
