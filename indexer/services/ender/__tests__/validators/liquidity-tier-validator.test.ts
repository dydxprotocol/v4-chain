import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  LiquidityTierUpsertEventV2,
} from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultHeight, defaultLiquidityTierUpsertEventV2, defaultTime, defaultTxHash,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { LiquidityTierValidatorV2 } from '../../src/validators/liquidity-tier-validator';

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
      const validator: LiquidityTierValidatorV2 = new LiquidityTierValidatorV2(
        defaultLiquidityTierUpsertEventV2,
        createBlock(defaultLiquidityTierUpsertEventV2),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      [
        'throws error on liquidity tier upsert event missing initialMarginPpm',
        {
          ...defaultLiquidityTierUpsertEventV2,
          initialMarginPpm: 0,
        } as LiquidityTierUpsertEventV2,
        'LiquidityTierUpsertEventV2 initialMarginPpm is not populated',
      ],
      [
        'throws error on perpetual market create event missing maintenanceFractionPpm',
        {
          ...defaultLiquidityTierUpsertEventV2,
          maintenanceFractionPpm: 0,
        } as LiquidityTierUpsertEventV2,
        'LiquidityTierUpsertEventV2 maintenanceFractionPpm is not populated',
      ],
    ])('%s', (_description: string, event: LiquidityTierUpsertEventV2, expectedMessage: string) => {
      const validator: LiquidityTierValidatorV2 = new LiquidityTierValidatorV2(
        event,
        createBlock(event),
        0,
      );
      expect(() => validator.validate()).toThrow(new ParseMessageError(expectedMessage));
    });

    it.each([
      [
        'logs error on liquidity tier upsert event with empty name',
        {
          ...defaultLiquidityTierUpsertEventV2,
          name: '',
        } as LiquidityTierUpsertEventV2,
        'LiquidityTierUpsertEventV2 name is not populated',
      ],
      // ... other test cases here ...
    ])('%s', (_description: string, event: LiquidityTierUpsertEventV2, expectedMessage: string) => {
      const loggerError = jest.spyOn(logger, 'error');

      const validator: LiquidityTierValidatorV2 = new LiquidityTierValidatorV2(
        event,
        createBlock(event),
        0,
      );
      validator.validate();
      expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
        message: expectedMessage,
      }));
    });

  });
});

function createBlock(
  liquidityTierEvent: LiquidityTierUpsertEventV2,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.LIQUIDITY_TIER,
    LiquidityTierUpsertEventV2.encode(liquidityTierEvent).finish(),
    0,
    2,
  );

  return createIndexerTendermintBlock(
    defaultHeight,
    defaultTime,
    [event],
    [defaultTxHash],
  );
}
