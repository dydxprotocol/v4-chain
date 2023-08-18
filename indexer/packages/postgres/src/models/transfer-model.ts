import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class TransferModel extends Model {
  static get tableName() {
    return 'transfers';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    recipientSubaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'transfers.recipientSubaccountId',
        to: 'subaccounts.id',
      },
    },
    senderSubaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'transfers.senderSubaccountId',
        to: 'subaccounts.id',
      },
    },
    asset: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'asset-model'),
      join: {
        from: 'transfers.assetId',
        to: 'assets.id',
      },
    },
    tendermintEvents: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'transfers.eventId',
        to: 'tendermint_events.id',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'transfers.createdAtHeight',
        to: 'blocks.blockHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'senderSubaccountId',
        'recipientSubaccountId',
        'assetId',
        'size',
        'eventId',
        'transactionHash',
        'createdAt',
        'createdAtHeight',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        senderSubaccountId: { type: 'string', format: 'uuid' },
        recipientSubaccountId: { type: 'string', format: 'uuid' },
        assetId: { type: 'string', pattern: IntegerPattern },
        size: { type: 'string', pattern: NonNegativeNumericPattern },
        transactionHash: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
      },
    };
  }

  id!: string;

  senderSubaccountId!: string;

  recipientSubaccountId!: string;

  assetId!: string;

  size!: string;

  eventId!: Buffer;

  transactionHash!: string;

  createdAt!: IsoString;

  createdAtHeight!: string;
}
