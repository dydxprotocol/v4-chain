import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  UpdateYieldParamsEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  dbHelpers, testMocks, perpetualMarketRefresher,
} from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultHeight,
  defaultTime,
  defaultTxHash,
  defaultUpdateYieldParamsEvent1,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { YieldParamsValidator } from '../../src/validators/yield-params-validator';

describe('yield-params-validator', () => {
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
    it('does not throw error on valid update yield params event', () => {
      const validator: YieldParamsValidator = new YieldParamsValidator(
        defaultUpdateYieldParamsEvent1,
        createBlock(defaultUpdateYieldParamsEvent1),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error on undefined assetYieldIndex', () => {
        const validator: YieldParamsValidator = new YieldParamsValidator(
          {
            ...defaultUpdateYieldParamsEvent1,
            assetYieldIndex: undefined as unknown as string, // satisfy type checker
          },
          createBlock(defaultUpdateYieldParamsEvent1),
          0,
        );

        expect(() => validator.validate()).toThrow(new ParseMessageError(
            'UpdateYieldParamsEvent must have an assetYieldIndex that is defined and non-empty',
          ));
    });

    it('throws error on empty assetYieldIndex', () => {
        const validator: YieldParamsValidator = new YieldParamsValidator(
          {
            ...defaultUpdateYieldParamsEvent1,
            assetYieldIndex: '',
          },
          createBlock(defaultUpdateYieldParamsEvent1),
          0,
        );

        expect(() => validator.validate()).toThrow(new ParseMessageError(
            'UpdateYieldParamsEvent must have an assetYieldIndex that is defined and non-empty',
          ));
    });

    it('throws error on undefined sDAIPrice', () => {
        const validator: YieldParamsValidator = new YieldParamsValidator(
          {
            ...defaultUpdateYieldParamsEvent1,
            sdaiPrice: undefined as unknown as string, // satisfy type checker
          },
          createBlock(defaultUpdateYieldParamsEvent1),
          0,
        );

        expect(() => validator.validate()).toThrow(new ParseMessageError(
            'UpdateYieldParamsEvent must have an sDAIPrice that is defined and non-empty',
          ));
    });

    it('throws error on empty sDAIPrice', () => {
        const validator: YieldParamsValidator = new YieldParamsValidator(
          {
            ...defaultUpdateYieldParamsEvent1,
            sdaiPrice: '',
          },
          createBlock(defaultUpdateYieldParamsEvent1),
          0,
        );

        expect(() => validator.validate()).toThrow(new ParseMessageError(
            'UpdateYieldParamsEvent must have an sDAIPrice that is defined and non-empty',
          ));
    });
  });
});

function createBlock(
  updateYieldParamsEvent: UpdateYieldParamsEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.YIELD_PARAMS,
    UpdateYieldParamsEventV1.encode(updateYieldParamsEvent).finish(),
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

