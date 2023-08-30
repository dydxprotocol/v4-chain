import path from 'path';

import { Model } from 'objection';

import UpsertQueryBuilder from '../query-builders/upsert';
import { IsoString } from '../types';
import BaseModel from './base-model';

export default class SubaccountModel extends BaseModel {
  static get tableName() {
    return 'subaccounts';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    orders: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'order-model'),
      join: {
        from: 'subaccounts.id',
        to: 'orders.subaccountId',
      },
    },
    height: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'subaccounts.updatedAtHeight',
        to: 'blocks.blockHeight',
      },
    },
    perpetualPositions: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'subaccounts.id',
        to: 'perpetual_positions.subaccountId',
      },
    },
    fills: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'fill-model'),
      join: {
        from: 'subaccounts.id',
        to: 'fills.subaccountId',
      },
    },
    assetPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'asset-position-model'),
      join: {
        from: 'subaccounts.id',
        to: 'asset_positions.subaccountId',
      },
    },
    fromTransfer: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'transfer-model'),
      join: {
        from: 'subaccounts.id',
        to: 'transfers.senderSubaccountId',
      },
    },
    toTransfer: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'transfer-model'),
      join: {
        from: 'subaccounts.id',
        to: 'transfers.recipientSubaccountId',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'address',
        'subaccountNumber',
        'updatedAt',
        'updatedAtHeight',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        address: { type: 'string' },
        subaccountNumber: { type: 'integer' },
        updatedAt: { type: 'string', format: 'date-time' },
        updatedAtHeight: { type: 'string' },
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
      address: 'string',
      subaccountNumber: 'integer',
      updatedAt: 'date-time',
      updatedAtHeight: 'string',
    };
  }

  id!: string;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  subaccountNumber!: number;

  updatedAt!: IsoString;

  updatedAtHeight!: string;
}
