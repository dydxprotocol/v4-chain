import { logger } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetModel,
  storeHelpers,
  SubaccountMessageContents,
  TransferFromDatabase,
  TransferModel,
} from '@dydxprotocol-indexer/postgres';
import { TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import { generateTransferContents } from '../helpers/kafka-helper';
import { indexerTendermintEventToTransactionIndex } from '../lib/helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class TransferHandler extends Handler<TransferEventV1> {
  eventType: string = 'TransferEvent';

  public getParallelizationIds(): string[] {
    // Must be handled sequentially with asset create events
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(): Promise<ConsolidatedKafkaEvent[]> {
    const transactionIndex: number = indexerTendermintEventToTransactionIndex(
      this.indexerTendermintEvent,
    );
    const eventDataBinary: Uint8Array = this.indexerTendermintEvent.dataBytes;
    const result: pg.QueryResult = await storeHelpers.rawQuery(
      `SELECT dydx_transfer_handler(
        ${this.block.height},
        '${this.block.time?.toISOString()}',
        '${JSON.stringify(TransferEventV1.decode(eventDataBinary))}',
        ${this.indexerTendermintEvent.eventIndex},
        ${transactionIndex},
        '${this.block.txHashes[transactionIndex]}'
      ) AS result;`,
      { txId: this.txId },
    ).catch((error: Error) => {
      logger.error({
        at: 'TransferHandler#internalHandle',
        message: 'Failed to handle TransferEventV1',
        error,
      });

      throw error;
    });

    const asset: AssetFromDatabase = AssetModel.fromJson(
      result.rows[0].result.asset) as AssetFromDatabase;
    const transfer: TransferFromDatabase = TransferModel.fromJson(
      result.rows[0].result.transfer) as TransferFromDatabase;
    return this.generateKafkaEvents(
      transfer,
      asset,
    );
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
