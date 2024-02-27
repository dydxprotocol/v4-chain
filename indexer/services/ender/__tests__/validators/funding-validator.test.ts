import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  FundingEventV1,
  FundingEventV1_Type,
  IndexerTendermintBlock,
  IndexerTendermintEvent,
} from '@dydxprotocol-indexer/v4-protos';
import { dbHelpers, perpetualMarketRefresher, testMocks } from '@dydxprotocol-indexer/postgres';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { FundingValidator } from '../../src/validators/funding-validator';
import {
  defaultFundingUpdateSampleEvent,
  defaultFundingRateEvent,
  defaultHeight,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';
import { bigIntToBytes } from '@dydxprotocol-indexer/v4-proto-parser';

describe('funding-validator', () => {
  beforeEach(async () => {
    await testMocks.seedData();
    await perpetualMarketRefresher.updatePerpetualMarkets();
    jest.spyOn(logger, 'error');
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it.each([
      ['funding rate event', defaultFundingRateEvent],
      ['funding premium sample event', defaultFundingUpdateSampleEvent],
    ])('does not throw error on valid %s', (_message: string, event: FundingEventV1) => {
      const validator: FundingValidator = new FundingValidator(
        event,
        createBlock(event),
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it('does not throw error if perpetualId does not exist', () => {
      const event: FundingEventV1 = {
        type: FundingEventV1_Type.TYPE_FUNDING_RATE_AND_INDEX,
        updates: [
          {
            perpetualId: 10,
            fundingValuePpm: 10,
            fundingIndex: bigIntToBytes(BigInt(0)),
          },
        ],
      } as FundingEventV1;
      const validator: FundingValidator = new FundingValidator(
        event,
        createBlock(event),
        0,
      );

      const errMsg: string = 'Invalid FundingEvent, perpetualId does not exist';
      expect(() => validator.validate()).not.toThrow(new ParseMessageError(errMsg));
      expect(logger.error).toHaveBeenCalledWith({
        at: `${FundingValidator.name}#validate`,
        message: errMsg,
        blockHeight: defaultHeight,
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        event,
      });
    });

    it('throws error if event does not specify valid type', () => {
      const event = {
        type: FundingEventV1_Type.TYPE_UNSPECIFIED,
        updates: [
          {
            perpetualId: 0,
            fundingValuePpm: 10,
            fundingIndex: bigIntToBytes(BigInt(0)),
          },
        ],
      } as FundingEventV1;
      const message = 'Invalid FundingEvent, type must be TYPE_PREMIUM_SAMPLE or TYPE_FUNDING_RATE_AND_INDEX';
      const validator = new FundingValidator(event, createBlock(event), 0);

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(FundingValidator.name, message, { event });
    });
  });
});

function createBlock(
  fundingEvent: FundingEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.FUNDING,
    FundingEventV1.encode(fundingEvent).finish(),
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
