import { stats, STATS_FUNCTION_NAME } from '@dydxprotocol-indexer/base';
import {
  IndexerTendermintBlock,
  IndexerTendermintEvent,
  Timestamp,
  TransferEventV1,
} from '@dydxprotocol-indexer/v4-protos';
import {
  AssetFromDatabase,
  assetRefresher,
  AssetTable,
  dbHelpers,
  Ordering,
  protocolTranslations,
  SubaccountCreateObject,
  SubaccountFromDatabase,
  SubaccountMessageContents,
  SubaccountTable,
  testMocks,
  TransferColumns,
  TransferFromDatabase,
  TransferTable,
} from '@dydxprotocol-indexer/postgres';
import { KafkaMessage } from 'kafkajs';
import { createKafkaMessage, producer } from '@dydxprotocol-indexer/kafka';
import { onMessage } from '../../src/lib/on-message';
import { DydxIndexerSubtypes } from '../../src/lib/types';
import {
  binaryToBase64String,
  createIndexerTendermintBlock,
  createIndexerTendermintEvent,
  expectSubaccountKafkaMessage,
} from '../helpers/indexer-proto-helpers';
import { generateTransferContents } from '../../src/helpers/kafka-helper';
import _ from 'lodash';
import { TransferHandler } from '../../src/handlers/transfer-handler';
import {
  defaultDateTime,
  defaultHeight,
  defaultPreviousHeight,
  defaultTime,
  defaultTransferEvent,
  defaultTxHash,
} from '../helpers/constants';
import { updateBlockCache } from '../../src/caches/block-cache';

describe('transferHandler', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
    jest.spyOn(stats, 'gauge');
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
    defaultTransferEvent.senderSubaccountId!,
  );

  const defaultRecipientSubaccountId: string = SubaccountTable.subaccountIdToUuid(
    defaultTransferEvent.recipientSubaccountId!,
  );

  const defaultSenderSubaccount: SubaccountCreateObject = {
    address: defaultTransferEvent.senderSubaccountId!.owner,
    subaccountNumber: defaultTransferEvent.senderSubaccountId!.number,
    updatedAt: defaultDateTime.toISO(),
    updatedAtHeight: defaultPreviousHeight,
  };

  const defaultRecipientSubaccount: SubaccountCreateObject = {
    address: defaultTransferEvent.recipientSubaccountId!.owner,
    subaccountNumber: defaultTransferEvent.recipientSubaccountId!.number,
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
        binaryToBase64String(
          TransferEventV1.encode(defaultTransferEvent).finish(),
        ),
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
      recipientSubaccountId: {
        owner: '',
        number: 0,
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

    // This is testing a temporary fix for the fact that protocol is sending transfer events for
    // withdrawals/deposits but Indexer is not yet ready to handle them.
    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    expect(producerSendMock.mock.calls.length).toEqual(0);
  });

  it('fails when TransferEvent does not contain recipient subaccountId', async () => {
    const transactionIndex: number = 0;
    const transferEvent: TransferEventV1 = TransferEventV1.fromPartial({
      senderSubaccountId: {
        owner: '',
        number: 0,
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

    // This is testing a temporary fix for the fact that protocol is sending transfer events for
    // withdrawals/deposits but Indexer is not yet ready to handle them.
    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    expect(producerSendMock.mock.calls.length).toEqual(0);
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
      const existingSubaccount: SubaccountFromDatabase | undefined = await SubaccountTable.findById(
        subaccountId,
      );
      expect(existingSubaccount).toBeDefined();
    });

    // Confirm there is no existing transfer to or from the recipient/sender subaccounts
    await expectNoExistingTransfers([defaultRecipientSubaccountId, defaultSenderSubaccountId]);

    const producerSendMock: jest.SpyInstance = jest.spyOn(producer, 'send');
    await onMessage(kafkaMessage);

    const newTransfer: TransferFromDatabase = await expectAndReturnNewTransfer(
      defaultRecipientSubaccountId, defaultSenderSubaccountId);

    expectTransferMatchesEvent(transferEvent, newTransfer, asset);

    await expectTransfersSubaccountKafkaMessage(
      producerSendMock,
      transferEvent,
      newTransfer,
      asset,
    );
    expectTimingStats();
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
      defaultRecipientSubaccountId, defaultSenderSubaccountId);

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
    expectTimingStats();
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
        binaryToBase64String(
          TransferEventV1.encode(transferEvent).finish(),
        ),
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

function expectTimingStats() {
  expectTimingStat('upsert_subaccounts');
  expectTimingStat('create_transfer_and_get_asset');
}

function expectTimingStat(fnName: string) {
  expect(stats.timing).toHaveBeenCalledWith(
    `ender.${STATS_FUNCTION_NAME}.timing`,
    expect.any(Number),
    {
      className: 'TransferHandler',
      eventType: 'TransferEvent',
      fnName,
    },
  );
}

function expectTransferMatchesEvent(
  event: TransferEventV1,
  transfer: TransferFromDatabase,
  asset: AssetFromDatabase,
) {
  expect(transfer.senderSubaccountId).toEqual(
    SubaccountTable.uuid(
      event!.senderSubaccountId!.owner,
      event!.senderSubaccountId!.number,
    ),
  );
  expect(transfer.recipientSubaccountId).toEqual(
    SubaccountTable.uuid(
      event!.recipientSubaccountId!.owner,
      event!.recipientSubaccountId!.number,
    ),
  );
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
  const transfers: TransferFromDatabase[] = await TransferTable.findAllToOrFromSubaccountId(
    {
      subaccountId: subaccountIds,
    },
    [], {
      orderBy: [[TransferColumns.id, Ordering.ASC]],
    });

  expect(transfers.length).toEqual(0);
}

async function expectAndReturnNewTransfer(
  recipientSubaccountId: string,
  senderSubaccountId: string,
): Promise<TransferFromDatabase> {
  // Confirm there is now a transfer to or from the recipient subaccount
  const newTransfersRelatedToRecipient: TransferFromDatabase[] = await
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

  // Confirm there is now a transfer to or from the sender subaccount
  const newTransfersRelatedToSender: TransferFromDatabase[] = await
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
  // Confirm the transfer is the same between recipient/sender.
  expect(newTransfersRelatedToSender).toEqual(newTransfersRelatedToRecipient);
  return newTransfersRelatedToSender[0];
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
  const contents: SubaccountMessageContents = generateTransferContents(
    event.senderSubaccountId!,
    event.recipientSubaccountId!,
    transfer,
    asset,
  );

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto: event.recipientSubaccountId!,
  });

  expectSubaccountKafkaMessage({
    producerSendMock,
    blockHeight,
    transactionIndex,
    eventIndex,
    contents: JSON.stringify(contents),
    subaccountIdProto: event.senderSubaccountId!,
  });
}
