import { logger, ParseMessageError } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  TransferEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import { TransferValidator } from '../../src/validators/transfer-validator';
import {
  defaultHeight,
  defaultSubaccountId,
  defaultTime,
  defaultTransferEvent,
  defaultTxHash,
  defaultWalletAddress,
} from '../helpers/constants';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import { expectDidntLogError, expectLoggedParseMessageError } from '../helpers/validator-helpers';

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
        0,
      );

      validator.validate();
      expectDidntLogError();
    });

    it.each([
      [
        'does not contain sender',
        {
          ...defaultTransferEvent,
          sender: undefined,
        },
        'TransferEvent must have either a sender subaccount id or sender wallet address',
      ],
      [
        'contains 2 senders',
        {
          ...defaultTransferEvent,
          sender: {
            address: 'defaultAddress',
            subaccountId: defaultSubaccountId,
          },
        },
        'TransferEvent must have either a sender subaccount id or sender wallet address',
      ],
      [
        'does not contain recipient',
        {
          ...defaultTransferEvent,
          recipient: undefined,
        },
        'TransferEvent must have either a recipient subaccount id or recipient wallet address',
      ],
      [
        'contains 2 recipients',
        {
          ...defaultTransferEvent,
          recipient: {
            address: 'defaultAddress',
            subaccountId: defaultSubaccountId,
          },
        },
        'TransferEvent must have either a recipient subaccount id or recipient wallet address',
      ],
      [
        'contains both a sender and recipient wallet address',
        {
          ...defaultTransferEvent,
          sender: {
            address: defaultWalletAddress,
          },
          recipient: {
            address: defaultWalletAddress,
          },
        },
        'TransferEvent cannot have both a sender and recipient wallet address',
      ],
    ])('throws error if event %s', (_message: string, event: TransferEventV1, message: string) => {
      const validator: TransferValidator = new TransferValidator(
        event,
        createBlock(event),
        0,
      );

      expect(() => validator.validate()).toThrow(new ParseMessageError(message));
      expectLoggedParseMessageError(
        TransferValidator.name,
        message,
        { event },
      );
    });
  });
});

function createBlock(
  transferEvent: TransferEventV1,
): IndexerTendermintBlock {
  const event: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.TRANSFER,
    TransferEventV1.encode(transferEvent).finish(),
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
