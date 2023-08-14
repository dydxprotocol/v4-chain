import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class AssetPositionModel extends BaseModel {
  static get tableName() {
    return 'asset_positions';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'asset_positions.subaccountId',
        to: 'subaccounts.id',
      },
    },
    asset: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'asset-model'),
      join: {
        from: 'asset_positions.assetId',
        to: 'assets.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'subaccountId',
        'assetId',
        'size',
        'isLong',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        subaccountId: { type: 'string', format: 'uuid' },
        assetId: { type: 'string', pattern: IntegerPattern },
        size: { type: 'string', pattern: NonNegativeNumericPattern },
        isLong: { type: 'boolean', default: true },
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
      subaccountId: 'string',
      assetId: 'string',
      size: 'string',
      isLong: 'boolean',
    };
  }

  id!: string;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  subaccountId!: string;

  assetId!: string;

  size!: string;

  isLong!: boolean;
}
