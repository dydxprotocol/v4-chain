import { stats } from '@dydxprotocol-indexer/base';
import {
  AssetFromDatabase,
  AssetModel,
  SubaccountMessageContents,
  TransferFromDatabase,
  TransferModel,
} from '@dydxprotocol-indexer/postgres';
import { TransferEventV1 } from '@dydxprotocol-indexer/v4-protos';
import * as pg from 'pg';

import config from '../config';
import { generateTransferContents } from '../helpers/kafka-helper';
import { ConsolidatedKafkaEvent } from '../lib/types';
import { Handler } from './handler';

export class TransferHandler extends Handler<TransferEventV1> {
  eventType: string = 'TransferEvent';

  public getParallelizationIds(): string[] {
    // Must be handled sequentially with asset create events
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/require-await
  public async internalHandle(resultRow: pg.QueryResultRow): Promise<ConsolidatedKafkaEvent[]> {
    // Handle latency from resultRow
    stats.timing(
      `${config.SERVICE_NAME}.handle_transfer_event.sql_latency`,
      Number(resultRow.latency),
      this.generateTimingStatsOptions(),
    );
    const asset: AssetFromDatabase = AssetModel.fromJson(
      resultRow.asset) as AssetFromDatabase;
    const transfer: TransferFromDatabase = TransferModel.fromJson(
      resultRow.transfer) as TransferFromDatabase;
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
        this.block.height.toString(),
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
        this.block.height.toString(),
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
