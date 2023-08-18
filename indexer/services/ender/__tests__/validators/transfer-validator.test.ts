import { logger } from '@dydxprotocol-indexer/base';
import { IndexerTendermintBlock, IndexerTendermintEvent, TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { TransferValidator } from '../../src/validators/transfer-validator';
import {
  defaultHeight, defaultTime, defaultTransferEvent, defaultTxHash,
} from '../helpers/constants';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
} from '../helpers/indexer-proto-helpers';
import { expectDidntLogError } from '../helpers/validator-helpers';

describe('transfer-validator', () => {
  beforeEach(() => {
    jest.spyOn(logger, 'error');
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  describe('validate', () => {
    it('does not throw error on valid transfer', () => {
      const validator: TransferValidator = new TransferValidator(
        defaultTransferEvent,
        createBlock(defaultTransferEvent),
      );

      validator.validate();
      expectDidntLogError();
    });

    // This is testing a temporary fix for the fact that protocol is sending transfer events for
    // withdrawals/deposits but Indexer is not yet ready to handle them.
    it.each([
      [
        'does not contain senderSubaccountId',
        {
          ...defaultTransferEvent,
          senderSubaccountId: undefined,
        },
      ],
      [
        'does not contain recipientSubaccountId',
        {
          ...defaultTransferEvent,
          recipientSubaccountId: undefined,
        },
      ],
      [
        'does not contain recipientSubaccountId or senderSubaccountId',
        {
          ...defaultTransferEvent,
          recipientSubaccountId: undefined,
          senderSubaccountId: undefined,
        },
      ],
    ])('doesnt throw error if event %s', (_message: string, event: TransferEventV1) => {
      const validator: TransferValidator = new TransferValidator(
        event,
        createBlock(event),
      );

      validator.validate();
      expectDidntLogError();
    });
  });
});

function createBlock(
  transferEvent: TransferEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.TRANSFER,
    binaryToBase64String(
      TransferEventV1.encode(transferEvent).finish(),
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
