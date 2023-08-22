import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import { PerpetualMarketCreateEventV1, IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, marketRefresher, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultPerpetualMarketCreateEvent, defaultHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { PerpetualMarketValidator } from '../../src/validators/perpetual-market-validator';
import Long from 'long';

describe('perpetual-market-validator', () => {
  beforeEach(async () => {
    await testMocks.seedData();
    await marketRefresher.updateMarkets();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await perpetualMarketRefresher.clear();
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid perpetual market create event', () => {
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        defaultPerpetualMarketCreateEvent,
        createBlock(defaultPerpetualMarketCreateEvent),
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error on existing perpetual market', async () => {
      await perpetualMarketRefresher.updatePerpetualMarkets();
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        defaultPerpetualMarketCreateEvent,
        createBlock(defaultPerpetualMarketCreateEvent),
      );
      const message: string = 'PerpetualMarketCreateEvent id already exists';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });

    it('throws error on perpetual market create event missing ticker', () => {
      const event: PerpetualMarketCreateEventV1 = {
        ...defaultPerpetualMarketCreateEvent,
        ticker: '',
      };
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(event),
      );
      const message: string = 'PerpetualMarketCreateEvent ticker is not populated';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });

    it('throws error on perpetual market create event missing subticksPerTick', () => {
      const event: PerpetualMarketCreateEventV1 = {
        ...defaultPerpetualMarketCreateEvent,
        subticksPerTick: 0,
      };
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(event),
      );
      const message: string = 'PerpetualMarketCreateEvent subticksPerTick is not populated';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });

    it('throws error on perpetual market create event missing minOrderBaseQuantums', () => {
      const event: PerpetualMarketCreateEventV1 = {
        ...defaultPerpetualMarketCreateEvent,
        minOrderBaseQuantums: Long.fromValue(0),
      };
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(event),
      );
      const message: string = 'PerpetualMarketCreateEvent minOrderBaseQuantums is not populated';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });

    it('throws error on perpetual market create event missing stepBaseQuantums', () => {
      const event: PerpetualMarketCreateEventV1 = {
        ...defaultPerpetualMarketCreateEvent,
        stepBaseQuantums: Long.fromValue(0),
      };
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(event),
      );
      const message: string = 'PerpetualMarketCreateEvent stepBaseQuantums is not populated';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });
  });
});

function createBlock(
  assetCreateEvent: PerpetualMarketCreateEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.PERPETUAL_MARKET,
    binaryToBase64String(
      PerpetualMarketCreateEventV1.encode(assetCreateEvent).finish(),
    ),
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
