import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  PerpetualMarketCreateEventV1,
  PerpetualMarketCreateEventV2,
  PerpetualMarketCreateEventV3,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultPerpetualMarketCreateEventV1,
  defaultPerpetualMarketCreateEventV2,
  defaultPerpetualMarketCreateEventV3,
  defaultHeight,
  defaultTime,
  defaultTxHash,
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
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    await perpetualMarketRefresher.clear();
    jest.clearAllMocks();
  });

  describe.each([
    [
      'PerpetualMarketCreateEventV1',
      1,
      PerpetualMarketCreateEventV1.encode(defaultPerpetualMarketCreateEventV1).finish(),
      defaultPerpetualMarketCreateEventV1,
    ],
    [
      'PerpetualMarketCreateEventV2',
      2,
      PerpetualMarketCreateEventV2.encode(defaultPerpetualMarketCreateEventV2).finish(),
      defaultPerpetualMarketCreateEventV2,
    ],
    [
      'PerpetualMarketCreateEventV3',
      3,
      PerpetualMarketCreateEventV3.encode(defaultPerpetualMarketCreateEventV3).finish(),
      defaultPerpetualMarketCreateEventV3,
    ],
  ])('validate %s', (
    _name: string,
    version: number,
    perpetualMarketCreateEventBytes: Uint8Array,
    event: PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
    | PerpetualMarketCreateEventV3,
  ) => {
    it('does not throw error on valid perpetual market create event', () => {
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(perpetualMarketCreateEventBytes),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error on existing perpetual market', async () => {
      await perpetualMarketRefresher.updatePerpetualMarkets();
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        event,
        createBlock(perpetualMarketCreateEventBytes),
        0,
      );
      const message: string = 'PerpetualMarketCreateEvent id already exists';
      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
    });

    it.each([
      [
        'throws error on perpetual market create event missing ticker',
        {
          ...event,
          ticker: '',
        } as PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
        | PerpetualMarketCreateEventV3,
        'PerpetualMarketCreateEvent ticker is not populated',
      ],
      [
        'throws error on perpetual market create event missing subticksPerTick',
        {
          ...event,
          subticksPerTick: 0,
        } as PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
        | PerpetualMarketCreateEventV3,
        'PerpetualMarketCreateEvent subticksPerTick is not populated',
      ],
      [
        'throws error on perpetual market create event missing stepBaseQuantums',
        {
          ...event,
          stepBaseQuantums: Long.fromValue(0, true),
        } as PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
        | PerpetualMarketCreateEventV3,
        'PerpetualMarketCreateEvent stepBaseQuantums is not populated',
      ],
    ])('%s', (
      _description: string,
      eventToTest: PerpetualMarketCreateEventV1 | PerpetualMarketCreateEventV2
      | PerpetualMarketCreateEventV3,
      expectedMessage: string,
    ) => {
      const validator: PerpetualMarketValidator = new PerpetualMarketValidator(
        eventToTest,
        createBlock(perpetualMarketCreateEventBytes),
        0,
      );
      expect(() => validator.validate()).toThrow(new ParseMessageError(expectedMessage));
    });
  });
});

function createBlock(
  perpetualMarketEventDataBytes: Uint8Array,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.PERPETUAL_MARKET,
    perpetualMarketEventDataBytes,
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
