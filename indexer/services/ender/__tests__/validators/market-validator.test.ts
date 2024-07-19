import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  MarketEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { MarketValidator } from '../../src/validators/market-validator';
import {
  defaultHeight,
  defaultMarketBase,
  defaultMarketCreate,
  defaultMarketModify,
  defaultMarketPriceUpdate,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';
import Long from 'long';

describe('market-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it.each([
      ['market create event', defaultMarketCreate],
      ['market modify event', defaultMarketModify],
      ['market price update event', defaultMarketPriceUpdate],
    ])('does not throw error on valid %s', (_message: string, event: MarketEventV1) => {
      const validator: MarketValidator = new MarketValidator(
        event,
        createBlock(event),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      // Base Validation Errors
      [
        'does not specify oneofKind',
        {
          marketId: 0,
          event: {
            oneofKind: undefined,
          },
        } as MarketEventV1,
        'One of marketCreate, marketModify, or priceUpdate must be defined in MarketEvent',
      ],

      // Market Create Validation Errors
      [
        'base field is undefined',
        {
          ...defaultMarketCreate,
          marketCreate: {
            base: undefined,
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketCreate, base field is undefined',
      ],
      [
        'pair is empty',
        {
          ...defaultMarketCreate,
          marketCreate: {
            base: {
              ...defaultMarketBase,
              pair: '',
            },
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketCreate, pair is empty',
      ],
      [
        'minPriceChangePpm is 0',
        {
          ...defaultMarketCreate,
          marketCreate: {
            base: {
              ...defaultMarketBase,
              minPriceChangePpm: 0,
            },
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketCreate, minPriceChangePpm is 0',
      ],

      // Market Modify Validation Errors
      [
        'base field is undefined',
        {
          ...defaultMarketModify,
          marketModify: {
            base: undefined,
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketModify, base field is undefined',
      ],
      [
        'pair is empty',
        {
          ...defaultMarketModify,
          marketModify: {
            base: {
              ...defaultMarketBase,
              pair: '',
            },
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketModify, pair is empty',
      ],
      [
        'minPriceChangePpm is 0',
        {
          ...defaultMarketModify,
          marketModify: {
            base: {
              ...defaultMarketBase,
              minPriceChangePpm: 0,
            },
            exponent: 10,
          },
        } as MarketEventV1,
        'Invalid MarketModify, minPriceChangePpm is 0',
      ],

      // Market Price Update Validation Errors
      [
        'has priceWithExponent = 0',
        {
          ...defaultMarketPriceUpdate,
          priceUpdate: {
            priceWithExponent: Long.fromValue(0, true),
          },
        } as MarketEventV1,
        'Invalid MarketPriceUpdate, priceWithExponent must be > 0',
      ],
    ])('throws error if event %s', (_message: string, event: MarketEventV1, message: string) => {
      const validator: MarketValidator = new MarketValidator(
        event,
        createBlock(event),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        MarketValidator.name,
        message,
        { event },
      );
    });
  });
});

function createBlock(
  marketEvent: MarketEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.MARKET,
    MarketEventV1.encode(marketEvent).finish(),
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
