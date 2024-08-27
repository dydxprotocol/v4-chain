import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class OraclePriceModel extends Model {
  static get tableName() {
    return 'oracle_prices';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    market: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'market-model'),
      join: {
        from: 'oracle_prices.marketId',
        to: 'markets.id',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'oracle_prices.effectiveAtHeight',
        to: 'blocks.blockHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'marketId',
        'spotPrice',
        'pnlPrice',
        'effectiveAt',
        'effectiveAtHeight',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        marketId: { type: 'integer' },
        spotPrice: { type: 'string', pattern: NonNegativeNumericPattern },
        pnlPrice: { type: 'string', pattern: NonNegativeNumericPattern },
        effectiveAt: { type: 'string', format: 'date-time' },
        effectiveAtHeight: { type: 'string', pattern: IntegerPattern },
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
      marketId: 'integer',
      spotPrice: 'string',
      pnlPrice: 'string',
      effectiveAt: 'date-time',
      effectiveAtHeight: 'string',
    };
  }

  id!: string;

  marketId!: number;

  spotPrice!: string;

  pnlPrice!: string;

  effectiveAt!: IsoString;

  effectiveAtHeight!: string;
}
