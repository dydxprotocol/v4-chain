import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  LiquidityTierUpsertEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultHeight, defaultLiquidityTierUpsertEvent, defaultTime, defaultTxHash,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { LiquidityTierValidator } from '../../src/validators/liquidity-tier-validator';
import Long from 'long';

describe('liquidity-tier-validator', () => {
  beforeEach(async () => {
    await testMocks.seedData();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid liquidity tier upsert event', () => {
      const validator: LiquidityTierValidator = new LiquidityTierValidator(
        defaultLiquidityTierUpsertEvent,
        createBlock(defaultLiquidityTierUpsertEvent),
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      [
        'throws error on liquidity tier upsert event missing initialMarginPpm',
        {
          ...defaultLiquidityTierUpsertEvent,
          initialMarginPpm: 0,
        } as LiquidityTierUpsertEventV1,
        'LiquidityTierUpsertEventV1 initialMarginPpm is not populated',
      ],
      [
        'throws error on perpetual market create event missing maintenanceFractionPpm',
        {
          ...defaultLiquidityTierUpsertEvent,
          maintenanceFractionPpm: 0,
        } as LiquidityTierUpsertEventV1,
        'LiquidityTierUpsertEventV1 maintenanceFractionPpm is not populated',
      ],
    ])('%s', (_description: string, event: LiquidityTierUpsertEventV1, expectedMessage: string) => {
      const validator: LiquidityTierValidator = new LiquidityTierValidator(
        event,
        createBlock(event),
      );
      expect(() => validator.validate()).toThrow(new ParseMessageError(expectedMessage));
    });

    it.each([
      [
        'logs error on liquidity tier upsert event with empty name',
        {
          ...defaultLiquidityTierUpsertEvent,
          name: '',
        } as LiquidityTierUpsertEventV1,
        'LiquidityTierUpsertEventV1 name is not populated',
      ],
      [
        'logs error on liquidity tier upsert event with basePositionNotional equal to 0',
        {
          ...defaultLiquidityTierUpsertEvent,
          basePositionNotional: Long.fromValue(0, true),
        } as LiquidityTierUpsertEventV1,
        'LiquidityTierUpsertEventV1 basePositionNotional is not populated',
      ],

      // ... other test cases here ...
    ])('%s', (_description: string, event: LiquidityTierUpsertEventV1, expectedMessage: string) => {
      const loggerError = jest.spyOn(logger, 'error');

      const validator: LiquidityTierValidator = new LiquidityTierValidator(
        event,
        createBlock(event),
      );
      validator.validate();
      expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
        message: expectedMessage,
      }));
    });

  });
});

function createBlock(
  liquidityTierEvent: LiquidityTierUpsertEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.LIQUIDITY_TIER,
    LiquidityTierUpsertEventV1.encode(liquidityTierEvent).finish(),
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
