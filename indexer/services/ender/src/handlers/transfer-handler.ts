import {
  AssetFromDatabase,
  assetRefresher,
  protocolTranslations,
  SubaccountMessageContents,
  SubaccountTable,
  TendermintEventTable,
  TransferCreateObject,
  TransferFromDatabase,
  TransferTable,
  WalletTable,
} from '@dydxprotocol-indexer/postgres';
import { TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { generateTransferContents } from '../helpers/kafka-helper';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent, TransferEventType } from '../lib/types';
import { Handler } from './handler';

export class TransferHandler extends Handler<TransferEventV1> {
  eventType: string = 'TransferEvent';

  public getParallelizationIds(): string[] {
    // Must be handled sequentially with asset create events
    return [];
  }

  public async internalHandle(
  ): Promise<ConsolidatedKafkaEvent[]> {
    await this.runFuncWithTimingStatAndErrorLogging(
      Promise.all([
        this.upsertRecipientSubaccount(),
        this.upsertWallets(),
      ]),
      this.generateTimingStatsOptions('upsert_recipient_subaccount_and_wallets'),
    );

    const asset: AssetFromDatabase = assetRefresher.getAssetFromId(
      this.event.assetId.toString(),
    );
    const transfer: TransferFromDatabase = await this.runFuncWithTimingStatAndErrorLogging(
      this.createTransferFromEvent(asset),
      this.generateTimingStatsOptions('create_transfer_and_get_asset'),
    );

    return this.generateKafkaEvents(
      transfer,
      asset,
    );
  }

  protected async createTransferFromEvent(asset: AssetFromDatabase): Promise<TransferFromDatabase> {
    const eventId: Buffer = TendermintEventTable.createEventId(
      this.block.height.toString(),
      indexerTendermintEventToTransactionIndex(this.indexerTendermintEvent),
      this.indexerTendermintEvent.eventIndex,
    );
    const senderWalletAddress: string | undefined = this.event.sender!.address;
    const recipientWalletAddress: string | undefined = this.event.recipient!.address;
    const senderSubaccountId: string | undefined = this.event.sender!.subaccountId
      ? SubaccountTable.subaccountIdToUuid(this.event.sender!.subaccountId!)
      : undefined;
    const recipientSubaccountId: string | undefined = this.event.recipient!.subaccountId
      ? SubaccountTable.subaccountIdToUuid(this.event.recipient!.subaccountId!)
      : undefined;

    const size: string = protocolTranslations.quantumsToHumanFixedString(
      this.event.amount.toString(),
      asset.atomicResolution,
    );
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );

    const transferToCreate: TransferCreateObject = {
      senderSubaccountId,
      recipientSubaccountId,
      senderWalletAddress,
      recipientWalletAddress,
      assetId: this.event.assetId.toString(),
      size,
      eventId,
      transactionHash: this.block.txHashes[transactionIndex],
      createdAt: this.timestamp.toISO(),
      createdAtHeight: this.block.height.toString(),
    };

    const transferFromDatabase: TransferFromDatabase = await TransferTable.create(
      transferToCreate,
      { txId: this.txId },
    );

    return transferFromDatabase;
  }

  protected async upsertRecipientSubaccount(): Promise<void> {
    if (this.event!.recipient!.subaccountId) {
      await SubaccountTable.upsert({
        address: this.event!.recipient!.subaccountId!.owner,
        subaccountNumber: this.event!.recipient!.subaccountId!.number,
        updatedAt: this.timestamp.toISO(),
        updatedAtHeight: this.block.height.toString(),
      }, { txId: this.txId });
    }
  }

  protected async upsertWallets(): Promise<void> {
    const promises = [];
    if (this.event!.sender!.address) {
      promises.push(
        WalletTable.upsert({
          address: this.event!.sender!.address,
        }, { txId: this.txId }),
      );
    }
    if (this.event!.recipient!.address) {
      promises.push(
        WalletTable.upsert({
          address: this.event!.recipient!.address,
        }, { txId: this.txId }),
      );
    }
    await Promise.all(promises);
  }

  protected getTransferType(): TransferEventType {
    if (this.event!.sender!.address) {
      return TransferEventType.DEPOSIT;
    }
    if (this.event!.recipient!.address) {
      return TransferEventType.WITHDRAWAL;
    }
    return TransferEventType.TRANSFER;
  }

  /** Generates a kafka websocket event for each subaccount involved in the transfer.
   *
   * If the transfer is between 2 subaccounts, 1 event for the sender subaccount and another
   * for the recipient will be generated.
   *
   * If the transfer is between a subaccount and a wallet, 1 event will be generated for the
   * subaccount.
   *
   * @param transfer
   * @param asset
   * @protected
   */
  protected generateKafkaEvents(
    transfer: TransferFromDatabase,
    asset: AssetFromDatabase,
  ): ConsolidatedKafkaEvent[] {
    const kafkaEvents: ConsolidatedKafkaEvent[] = [];

    if (this.event.sender!.subaccountId) {
      const senderContents: SubaccountMessageContents = generateTransferContents(
        transfer,
        asset,
        this.event.sender!.subaccountId!,
        this.event.sender!.subaccountId,
        this.event.recipient!.subaccountId,
      );

      kafkaEvents.push(
        this.generateConsolidatedSubaccountKafkaEvent(
          JSON.stringify(senderContents),
          this.event.sender!.subaccountId!,
        ),
      );
    }

    if (this.event.recipient!.subaccountId) {
      const recipientContents: SubaccountMessageContents = generateTransferContents(
        transfer,
        asset,
        this.event.recipient!.subaccountId!,
        this.event.sender!.subaccountId,
        this.event.recipient!.subaccountId,
      );

      kafkaEvents.push(
        this.generateConsolidatedSubaccountKafkaEvent(
          JSON.stringify(recipientContents),
          this.event.recipient!.subaccountId!,
        ),
      );
    }

    return kafkaEvents;
  }
}
