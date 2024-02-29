import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import { DeleveragingEventV1, IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { DeleveragingValidator } from '../../src/validators/deleveraging-validator';
import {
  defaultDeleveragingEvent, defaultHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';
import Long from 'long';

describe('deleveraging-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid deleveraging', () => {
      const validator: DeleveragingValidator = new DeleveragingValidator(
        defaultDeleveragingEvent,
        createBlock(defaultDeleveragingEvent),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      [
        'does not contain liquidated',
        {
          ...defaultDeleveragingEvent,
          liquidated: undefined,
        },
        'DeleveragingEvent must have a liquidated subaccount id',
      ],
      [
        'does not contain offsetting',
        {
          ...defaultDeleveragingEvent,
          offsetting: undefined,
        },
        'DeleveragingEvent must have an offsetting subaccount id',
      ],
      [
        'has fillAmount of 0',
        {
          ...defaultDeleveragingEvent,
          fillAmount: new Long(0),
        },
        'DeleveragingEvent fillAmount cannot equal 0',
      ],
      [
        'has totalQuoteQuantums of 0',
        {
          ...defaultDeleveragingEvent,
          totalQuoteQuantums: new Long(0),
        },
        'DeleveragingEvent totalQuoteQuantums cannot equal 0',
      ],
    ])('throws error if event %s', (_message: string, event: DeleveragingEventV1, message: string) => {
      const validator: DeleveragingValidator = new DeleveragingValidator(
        event,
        createBlock(event),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        DeleveragingValidator.name,
        message,
        { event },
      );
    });
  });
});

function createBlock(
  deleveragingEvent: DeleveragingEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.DELEVERAGING,
    DeleveragingEventV1.encode(deleveragingEvent).finish(),
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
