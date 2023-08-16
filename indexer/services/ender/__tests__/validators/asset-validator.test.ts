import { logger } from '@dydxprotocol-indexer/base';
import { AssetCreateEventV1, IndexerTendermintBlock, IndexerTendermintEvent } from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, marketRefresher, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  defaultAssetCreateEvent, defaultHeight, defaultTime, defaultTxHash,
} from '../helpers/constants';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';
import { AssetValidator } from '../../src/validators/asset-validator';

describe('asset-validator', () => {
  beforeEach(async () => {
    await testMocks.seedData();
    await marketRefresher.updateMarkets();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid asset create event', () => {
      const validator: AssetValidator = new AssetValidator(
        defaultAssetCreateEvent,
        createBlock(defaultAssetCreateEvent),
      );

      validator.validate();
      expectDidntLogError();
    });

    it('throws error on invalid asset create event', () => {
      const event: AssetCreateEventV1 = {
        ...defaultAssetCreateEvent,
        marketId: 1000,
      };
      const validator: AssetValidator = new AssetValidator(
        event,
        createBlock(event),
      );
      const message: string = 'Unable to find market with id: 1000';
      expect(() => validator.validate()).toThrow(new Error(message));
    });
  });
});

function createBlock(
  assetCreateEvent: AssetCreateEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.ASSET,
    binaryToBase64String(
      AssetCreateEventV1.encode(assetCreateEvent).finish(),
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
