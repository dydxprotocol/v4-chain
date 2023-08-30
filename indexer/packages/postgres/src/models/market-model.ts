import path from 'path';

import { Model } from 'objection';

import { NonNegativeNumericPattern } from '../lib/validators';

export default class MarketModel extends Model {
  static get tableName() {
    return 'markets';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    perpetualMarket: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'perpetual-market-model'),
      join: {
        from: 'markets.id',
        to: 'perpetual_markets.marketId',
      },
    },
    oraclePrices: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'oracle-price-model'),
      join: {
        from: 'markets.id',
        to: 'oracle_prices.marketId',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'pair',
        'exponent',
        'minPriceChangePpm',
      ],
      properties: {
        id: { type: 'integer' },
        pair: { type: 'string' },
        exponent: { type: 'integer' },
        minPriceChangePpm: { type: 'integer' },
        oraclePrice: { type: ['string', 'null'], pattern: NonNegativeNumericPattern, default: null },
      },
    };
  }

  id!: number;

  pair!: string;

  exponent!: number;

  minPriceChangePpm!: number;

  oraclePrice?: string;
}
