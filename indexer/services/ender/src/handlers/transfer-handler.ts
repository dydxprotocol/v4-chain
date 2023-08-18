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
} from '@dydxprotocol-indexer/postgres';
import { TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';

import { generateTransferContents } from '../helpers/kafka-helper';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class TransferHandler extends Handler<TransferEventV1> {
  eventType: string = 'TransferEvent';

  public getParallelizationIds(): string[] {
    // All transfers can be processed in parallel as this handler only creates new rows in the
    // database and does not modify the subaccount's asset positions (done in
    // SubaccountUpdateEvent), so no parallelization ids
    return [];
  }

  public async internalHandle(
  ): Promise<ConsolidatedKafkaEvent[]> {
    // This is a temporary fix for the fact that protocol is sending transfer events for
    // withdrawals/deposits but Indexer is not yet ready to handle them.
    if (
      this.event.senderSubaccountId === undefined || this.event.recipientSubaccountId === undefined
    ) {
      return [];
    }
    await this.runFuncWithTimingStatAndErrorLogging(
      Promise.all([
        this.upsertRecipientSubaccount(),
        // This is a temporary fix for the fact that protocol is sending transfer events for
        // withdrawals/deposits but Indexer is not yet ready to handle them.
        this.upsertSenderSubaccount(),
      ]),
      this.generateTimingStatsOptions('upsert_subaccounts'),
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
    const senderSubaccountId: string = SubaccountTable.subaccountIdToUuid(
      this.event.senderSubaccountId!,
    );
    const recipientSubaccountId: string = SubaccountTable.subaccountIdToUuid(
      this.event.recipientSubaccountId!,
    );

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
    await SubaccountTable.upsert({
      address: this.event!.recipientSubaccountId!.owner,
      subaccountNumber: this.event!.recipientSubaccountId!.number,
      updatedAt: this.timestamp.toISO(),
      updatedAtHeight: this.block.height.toString(),
    }, { txId: this.txId });
  }

  // This is a temporary fix for the fact that protocol is sending transfer events for
  // withdrawals/deposits but Indexer is not yet ready to handle them.
  protected async upsertSenderSubaccount(): Promise<void> {
    await SubaccountTable.upsert({
      address: this.event!.senderSubaccountId!.owner,
      subaccountNumber: this.event!.senderSubaccountId!.number,
      updatedAt: this.timestamp.toISO(),
      updatedAtHeight: this.block.height.toString(),
    }, { txId: this.txId });
  }

  /** Generates 2 kafka websocket events, 1 for the sender subaccount and another for the recipient
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
    const contents: SubaccountMessageContents = generateTransferContents(
      this.event.senderSubaccountId!,
      this.event.recipientSubaccountId!,
      transfer,
      asset,
    );

    return [
      this.generateConsolidatedSubaccountKafkaEvent(
        JSON.stringify(contents),
        this.event.senderSubaccountId!,
      ),
      this.generateConsolidatedSubaccountKafkaEvent(
        JSON.stringify(contents),
        this.event.recipientSubaccountId!,
      ),
    ];
  }
}
