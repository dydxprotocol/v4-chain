import path from 'path';

import { Model } from 'objection';

import { IntegerPattern } from '../lib/validators';

export default class AssetModel extends Model {
  static get tableName() {
    return 'assets';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    assetPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'asset-position-model'),
      join: {
        from: 'assets.id',
        to: 'asset_positions.assetId',
      },
    },
    transfers: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'transfer-model'),
      join: {
        from: 'assets.id',
        to: 'transfers.assetId',
      },
    },
    markets: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'market-model'),
      join: {
        from: 'assets.marketId',
        to: 'markets.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'symbol',
        'atomicResolution',
        'hasMarket',
      ],
      properties: {
        id: { type: 'string', pattern: IntegerPattern },
        symbol: { type: 'string' },
        atomicResolution: { type: 'integer' },
        hasMarket: { type: 'boolean', default: false },
        marketId: { type: ['integer', 'null'] },
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
      symbol: 'string',
      atomicResolution: 'integer',
      hasMarket: 'boolean',
      marketId: 'integer',
    };
  }

  id!: string;

  symbol!: string;

  atomicResolution!: number;

  hasMarket!: boolean;

  marketId?: number;
}
