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
    recipientWallet: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'wallet-model'),
      join: {
        from: 'transfers.recipientWalletAddress',
        to: 'wallets.address',
      },
    },
    senderWallet: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'wallet-model'),
      join: {
        from: 'transfers.senderWalletAddress',
        to: 'wallets.address',
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
        'senderWalletAddress',
        'recipientWalletAddress',
        'assetId',
        'size',
        'eventId',
        'transactionHash',
        'createdAt',
        'createdAtHeight',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        senderSubaccountId: { type: ['string', 'null'], default: null, format: 'uuid' },
        recipientSubaccountId: { type: ['string', 'null'], default: null, format: 'uuid' },
        senderWalletAddress: { type: ['string', 'null'], default: null },
        recipientWalletAddress: { type: ['string', 'null'], default: null },
        assetId: { type: 'string', pattern: IntegerPattern },
        size: { type: 'string', pattern: NonNegativeNumericPattern },
        transactionHash: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   *
   * TODO(IND-239): Ensure that jsonSchema() / sqlToJsonConversions() / model fields match.
   */
  static get sqlToJsonConversions() {
    return {
      id: 'string',
      senderSubaccountId: 'string',
      recipientSubaccountId: 'string',
      senderWalletAddress: 'string',
      recipientWalletAddress: 'string',
      assetId: 'string',
      size: 'string',
      eventId: 'hex-string',
      transactionHash: 'string',
      createdAt: 'date-time',
      createdAtHeight: 'string',
    };
  }

  id!: string;

  senderSubaccountId?: string;

  recipientSubaccountId?: string;

  senderWalletAddress?: string;

  recipientWalletAddress?: string;

  assetId!: string;

  size!: string;

  eventId!: Buffer;

  transactionHash!: string;

  createdAt!: IsoString;

  createdAtHeight!: string;
}
