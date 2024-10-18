import {
  logger,
  ParseMessageError,
} from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  TransferEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  AssetTable,
  AssetFromDatabase,
  dbHelpers,
  Ordering,
  SubaccountFromDatabase,
  SubaccountTable,
  testMocks,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
  SubaccountCreateObject,
  protocolTranslations,
  SubaccountMessageContents,
  assetRefresher,
  WalletTable,
  WalletFromDatabase,
  testConstants,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  createIndexerTendermintBlock,
  createIndexerTendermintEvent, expectSubaccountKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { generateTransferContents } from '../../src/helpers/kafka-helper';
import _ from 'lodash';
import { TransferHandler } from '../../src/handlers/transfer-handler';
import {
  defaultDateTime,
  defaultDepositEvent,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTransferEvent,
  defaultTxHash,
  defaultWalletAddress,
  defaultWithdrawalEvent,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';
import { createPostgresFunctions } from '../../src/helpers/postgres/postgres-functions';

const defaultWallet = {
  ...testConstants.defaultWallet,
  address: defaultWalletAddress,
};

describe('transferHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    await createPostgresFunctions();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    asset = await AssetTable.findById(
      defaultTransferEvent.assetId.toString(),
    ) as AssetFromDatabase;
    await assetRefresher.updateAssets();
    updateBlockCache(defaultPreviousHeight);
  });

  afterEach(async () => {
    await dbHelpers.clearData();
    jest.clearAllMocks();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
    jest.resetAllMocks();
  });

  const defaultSenderSubaccountId: string = SubaccountTable.subaccountIdToUuid(
    defaultTransferEvent.sender!.subaccountId!,
  );

  const defaultRecipientSubaccountId: string = SubaccountTable.subaccountIdToUuid(
    defaultTransferEvent.recipient!.subaccountId!,
  );

  const defaultSenderSubaccount: SubaccountCreateObject = {
    address: defaultTransferEvent.sender!.subaccountId!.owner,
    subaccountNumber: defaultTransferEvent.sender!.subaccountId!.number,
    updatedAt: defaultDateTime.toISO(),
    updatedAtHeight: defaultPreviousHeight,
  };

  const defaultRecipientSubaccount: SubaccountCreateObject = {
    address: defaultTransferEvent.recipient!.subaccountId!.owner,
    subaccountNumber: defaultTransferEvent.recipient!.subaccountId!.number,
    updatedAt: defaultDateTime.toISO(),
    updatedAtHeight: defaultPreviousHeight,
  };

  let asset: AssetFromDatabase;

  describe('getParallelizationIds', () => {
    it('returns the correct parallelization ids', () => {
      const transactionIndex: number = 0;
      const eventIndex: number = 0;

      const indexerTendermintEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        TransferEventV1.encode(defaultTransferEvent).finish(),
        transactionIndex,
        eventIndex,
      );
      const block: IndexerTendermintBlock = createIndexerTendermintBlock(
        0,
        defaultTime,
        [indexerTendermintEvent],
        [defaultTxHash],
      );

      const handler: TransferHandler = new TransferHandler(
        block,
        0,
        indexerTendermintEvent,
        0,
        defaultTransferEvent,
      );

      expect(handler.getParallelizationIds()).toEqual([]);
    });
  });

  it('fails when TransferEvent does not contain sender subaccountId', async () => {
    const transactionIndex: number = 0;
    const transferEvent: TransferEventV1 = TransferEventV1.fromPartial({
      recipient: {
        subaccountId: {
          owner: '',
          number: 0,
        },
      },
      assetId: 0,
      amount: 100,
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const loggerCrit = jest.spyOn(logger, 'crit');
    const loggerError = jest.spyOn(logger, 'error');
    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      new ParseMessageError(
        'TransferEvent must have either a sender subaccount id or sender wallet address',
      ),
    );

    expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
      at: 'TransferValidator#logAndThrowParseMessageError',
      message: 'TransferEvent must have either a sender subaccount id or sender wallet address',
    }));
    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
    }));
  });

  it('fails when TransferEvent does not contain recipient subaccountId', async () => {
    const transactionIndex: number = 0;
    const transferEvent: TransferEventV1 = TransferEventV1.fromPartial({
      sender: {
        subaccountId: {
          owner: '',
          number: 0,
        },
      },
      assetId: 0,
      amount: 100,
    });
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    const loggerCrit = jest.spyOn(logger, 'crit');
    const loggerError = jest.spyOn(logger, 'error');
    await expect(onMessage(kafkaMessage)).rejects.toThrowError(
      new ParseMessageError(
        'TransferEvent must have either a recipient subaccount id or recipient wallet address',
      ),
    );

    expect(loggerError).toHaveBeenCalledWith(expect.objectContaining({
      at: 'TransferValidator#logAndThrowParseMessageError',
      message: 'TransferEvent must have either a recipient subaccount id or recipient wallet address',
    }));
    expect(loggerCrit).toHaveBeenCalledWith(expect.objectContaining({
      at: 'onMessage#onMessage',
      message: 'Error: Unable to parse message, this must be due to a bug in V4 node',
    }));
  });

  it('creates new transfer for existing subaccounts', async () => {
    const transactionIndex: number = 0;

    const transferEvent: TransferEventV1 = defaultTransferEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Create the subaccounts
    await Promise.all([
      SubaccountTable.upsert(defaultSenderSubaccount),
      SubaccountTable.upsert(defaultRecipientSubaccount),
    ]);

    // Confirm there are subaccounts
    const subaccountIds: string[] = [defaultSenderSubaccountId, defaultRecipientSubaccountId];
    _.each(subaccountIds, async (subaccountId) => {
      const existingSubaccount:
      SubaccountFromDatabase | undefined = await SubaccountTable.findById(
        subaccountId,
      );
      expect(existingSubaccount).toBeDefined();
    });

    // Confirm there is no existing transfer to or from the recipient/sender subaccounts
    await expectNoExistingTransfers([defaultRecipientSubaccountId, defaultSenderSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer({
      recipientSubaccountId: defaultRecipientSubaccountId,
      senderSubaccountId: defaultSenderSubaccountId,
    });

    expectTransferMatchesEvent(transferEvent, newTransfer, asset);

    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      transferEvent,
      newTransfer,
      asset,
    );
  });

  it('creates new deposit for existing subaccount', async () => {
    const transactionIndex: number = 0;

    const depositEvent: TransferEventV1 = defaultDepositEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent: depositEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Create the subaccounts
    await Promise.all([
      SubaccountTable.upsert(defaultRecipientSubaccount),
    ]);

    // Confirm there is a recipient subaccount
    const existingSubaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      defaultRecipientSubaccountId,
    );
    expect(existingSubaccount).toBeDefined();

    // Confirm there is no existing transfer to or from the recipient subaccount
    await expectNoExistingTransfers([defaultRecipientSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer(
      {
        recipientSubaccountId: defaultRecipientSubaccountId,
      },
    );

    expectTransferMatchesEvent(depositEvent, newTransfer, asset);

    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      depositEvent,
      newTransfer,
      asset,
    );
    // Confirm the wallet was created
    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultWalletAddress,
    );
    expect(wallet).toEqual(defaultWallet);
  });

  it('creates new deposit for previously non-existent subaccount (also non-existent recipient wallet)', async () => {
    const transactionIndex: number = 0;

    const depositEvent: TransferEventV1 = defaultDepositEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent: depositEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Confirm there is no recipient subaccount
    const existingSubaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      defaultRecipientSubaccountId,
    );
    expect(existingSubaccount).toBeUndefined();

    // Confirm there is no existing transfer to or from the recipient subaccount
    await expectNoExistingTransfers([defaultRecipientSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer(
      {
        recipientSubaccountId: defaultRecipientSubaccountId,
      },
    );

    expectTransferMatchesEvent(depositEvent, newTransfer, asset);
    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      depositEvent,
      newTransfer,
      asset,
    );
    // Confirm the wallet was created for the sender and recipient
    const walletSender: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultWalletAddress,
    );
    const walletRecipient: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultDepositEvent.recipient!.subaccountId!.owner,
    );
    const newRecipientSubaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(
      defaultRecipientSubaccountId,
    );
    expect(newRecipientSubaccount).toBeDefined();
    expect(walletSender).toEqual(defaultWallet);
    expect(walletRecipient).toEqual({
      ...defaultWallet,
      address: defaultDepositEvent.recipient!.subaccountId!.owner,
    });
  });

  it('creates new withdrawal for existing subaccount', async () => {
    const transactionIndex: number = 0;

    const withdrawalEvent: TransferEventV1 = defaultWithdrawalEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent: withdrawalEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    // Create the subaccounts
    await Promise.all([
      SubaccountTable.upsert(defaultSenderSubaccount),
    ]);

    // Confirm there is a sender subaccount
    const existingSubaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
      defaultSenderSubaccountId,
    );
    expect(existingSubaccount).toBeDefined();

    // Confirm there is no existing transfer to or from the sender subaccount
    await expectNoExistingTransfers([defaultSenderSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer(
      {
        senderSubaccountId: defaultSenderSubaccountId,
      },
    );

    expectTransferMatchesEvent(withdrawalEvent, newTransfer, asset);

    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      withdrawalEvent,
      newTransfer,
      asset,
    );
    // Confirm the wallet was created
    const wallet: WalletFromDatabase | undefined = await WalletTable.findById(
      defaultWalletAddress,
    );
    expect(wallet).toEqual(defaultWallet);
  });

  it('creates new transfer and the recipient subaccount', async () => {
    const transactionIndex: number = 0;

    const transferEvent: TransferEventV1 = defaultTransferEvent;
    const kafkaMessage: KafkaMessage = createKafkaMessageFromTransferEvent({
      transferEvent,
      transactionIndex,
      height: defaultHeight,
      time: defaultTime,
      txHash: defaultTxHash,
    });

    await SubaccountTable.upsert(defaultSenderSubaccount);

    // Confirm there is 1 subaccount
    const existingSenderSubaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(
      defaultSenderSubaccountId,
    );
    expect(existingSenderSubaccount).toBeDefined();
    const existingRecipientSubaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(
      defaultRecipientSubaccountId,
    );
    expect(existingRecipientSubaccount).toBeUndefined();

    // Confirm there is no existing transfers
    await expectNoExistingTransfers([defaultRecipientSubaccountId, defaultSenderSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer(
      {
        recipientSubaccountId: defaultRecipientSubaccountId,
        senderSubaccountId: defaultSenderSubaccountId,
      });

    expectTransferMatchesEvent(transferEvent, newTransfer, asset);
    const newRecipientSubaccount: SubaccountFromDatabase | undefined = await
    SubaccountTable.findById(
      defaultRecipientSubaccountId,
    );
    expect(newRecipientSubaccount).toBeDefined();

    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      transferEvent,
      newTransfer,
      asset,
    );
  });
});

function createKafkaMessageFromTransferEvent({
  transferEvent,
  transactionIndex,
  height,
  time,
  txHash,
}: {
  transferEvent: TransferEventV1 | undefined,
  transactionIndex: number,
  height: number,
  time: Timestamp,
  txHash: string,
}) {
  const events: IndexerTendermintEvent[] = [];
  let eventIndex: number = 0;
  if (transferEvent !== undefined) {
    events.push(
      createIndexerTendermintEvent(
        DydxIndexerSubtypes.TRANSFER,
        TransferEventV1.encode(transferEvent).finish(),
        transactionIndex,
        eventIndex,
      ),
    );
    eventIndex += 1;
  }

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    height,
    time,
    events,
    [txHash],
  );

  const binaryBlock: Uint8Array = IndexerTendermintBlock.encode(block).finish();
  return createKafkaMessage(Buffer.from(binaryBlock));
}

function expectTransferMatchesEvent(
  event: TransferEventV1,
  transfer: TransferFromDatabase,
  asset: AssetFromDatabase,
) {
  if (transfer.senderSubaccountId) {
    expect(transfer.senderSubaccountId).toEqual(
      SubaccountTable.uuid(
        event!.sender!.subaccountId!.owner,
        event!.sender!.subaccountId!.number,
      ),
    );
  }
  if (transfer.recipientSubaccountId) {
    expect(transfer.recipientSubaccountId).toEqual(
      SubaccountTable.uuid(
        event!.recipient!.subaccountId!.owner,
        event!.recipient!.subaccountId!.number,
      ),
    );
  }
  if (transfer.senderWalletAddress) {
    expect(transfer.senderWalletAddress).toEqual(event.sender!.address);
  }
  if (transfer.recipientWalletAddress) {
    expect(transfer.recipientWalletAddress).toEqual(event.recipient!.address);
  }
  expect(transfer.assetId).toEqual(event.assetId.toString());
  expect(transfer.size).toEqual(
    protocolTranslations.quantumsToHumanFixedString(
      event.amount.toString(),
      asset.atomicResolution,
    ));
}

async function expectNoExistingTransfers(
  subaccountIds: string[],
) {
  // Confirm there is no existing transfer to or from the subaccounts
  const { results: transfers } = await TransferTable.findAllToOrFromSubaccountId(
    {
      subaccountId: subaccountIds,
    },
    [], {
      orderBy: [[TransferColumns.id, Ordering.ASC]],
    });

  expect(transfers.length).toEqual(0);
}

async function expectAndReturnNewTransfer(
  {
    recipientSubaccountId,
    senderSubaccountId,
    recipientWalletAddress,
    senderWalletAddress,
  } :
  {
    recipientSubaccountId?: string,
    senderSubaccountId?: string,
    recipientWalletAddress?: string,
    senderWalletAddress?: string,
  },
): Promise<TransferFromDatabase> {
  // Confirm there is now a transfer to or from the recipient subaccount
  if (recipientSubaccountId) {
    const { results: newTransfersRelatedToRecipient } = await
    TransferTable.findAllToOrFromSubaccountId(
      {
        subaccountId: [
          recipientSubaccountId,
        ],
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });
    expect(newTransfersRelatedToRecipient.length).toEqual(1);
    return newTransfersRelatedToRecipient[0];
  }

  if (senderSubaccountId) {
    // Confirm there is now a transfer to or from the sender subaccount
    const { results: newTransfersRelatedToSender } = await
    TransferTable.findAllToOrFromSubaccountId(
      {
        subaccountId: [
          senderSubaccountId,
        ],
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });

    expect(newTransfersRelatedToSender.length).toEqual(1);
    return newTransfersRelatedToSender[0];
  }
  if (recipientWalletAddress) {
    const newWithdrawal: TransferFromDatabase[] = await
    TransferTable.findAll(
      {
        recipientWalletAddress: [
          recipientWalletAddress,
        ],
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });
    expect(newWithdrawal.length).toEqual(1);
    return newWithdrawal[0];
  }
  if (senderWalletAddress) {
    const newWithdrawal: TransferFromDatabase[] = await
    TransferTable.findAll(
      {
        senderWalletAddress: [
          senderWalletAddress,
        ],
      },
      [], {
        orderBy: [[TransferColumns.id, Ordering.ASC]],
      });
    expect(newWithdrawal.length).toEqual(1);
    return newWithdrawal[0];
  }
  throw new Error('No subaccount or wallet address provided');
}

function expectTransfersSubaccountKafkaMessage(
  producerSendMock: jest.SpyInstance,
  event: TransferEventV1,
  transfer: TransferFromDatabase,
  asset: AssetFromDatabase,
  blockHeight: string = '3',
  transactionIndex: number = 0,
  eventIndex: number = 0,
) {
  let senderContents: SubaccountMessageContents = {};
  let recipientContents: SubaccountMessageContents = {};
  if (event.sender!.subaccountId) {
    senderContents = generateTransferContents(
      transfer,
      asset,
      event.sender!.subaccountId!,
      event.sender!.subaccountId!,
      event.recipient!.subaccountId,
      blockHeight,
    );
  }

  if (event.recipient!.subaccountId) {
    recipientContents = generateTransferContents(
      transfer,
      asset,
      event.recipient!.subaccountId!,
      event.sender!.subaccountId,
      event.recipient!.subaccountId!,
      blockHeight,
    );
  }

  if (event.sender!.subaccountId) {
    expectSubaccountKafkaMessage({
      producerSendMock,
      blockHeight,
      transactionIndex,
      eventIndex,
      contents: JSON.stringify(senderContents),
      subaccountIdProto: event.sender!.subaccountId!,
    });
  }

  if (event.recipient!.subaccountId) {
    expectSubaccountKafkaMessage({
      producerSendMock,
      blockHeight,
      transactionIndex,
      eventIndex,
      contents: JSON.stringify(recipientContents),
      subaccountIdProto: event.recipient!.subaccountId!,
    });
  }
}
