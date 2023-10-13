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

    it.each([
      [
        'throws error on perpetual market create event missing ticker',
        {
          ...defaultPerpetualMarketCreateEvent,
          ticker: '',
        } as PerpetualMarketCreateEventV1,
        'PerpetualMarketCreateEvent ticker is not populated',
      ],
      [
        'throws error on perpetual market create event missing subticksPerTick',
        {
          ...defaultPerpetualMarketCreateEvent,
          subticksPerTick: 0,
        } as PerpetualMarketCreateEventV1,
        'PerpetualMarketCreateEvent subticksPerTick is not populated',
      ],
      [
        'throws error on perpetual market create event missing stepBaseQuantums',
        {
          ...defaultPerpetualMarketCreateEvent,
          stepBaseQuantums: Long.fromValue(0, true),
        } as PerpetualMarketCreateEventV1,
        'PerpetualMarketCreateEvent stepBaseQuantums is not populated',
      ],
    ])('%s', (_description: string, event: PerpetualMarketCreateEventV1, expectedMessage: string) => {
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(event),
      );
      expect(() => validator.validate()).toThrow(new ParseMessageError(expectedMessage));
    });
  });
});

function createBlock(
  perpetualMarketEvent: PerpetualMarketCreateEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.PERPETUAL_MARKET,
    PerpetualMarketCreateEventV1.encode(perpetualMarketEvent).finish(),
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
