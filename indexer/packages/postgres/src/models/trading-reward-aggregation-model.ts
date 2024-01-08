import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import { IsoString } from '../types';
import { TradingRewardAggregationPeriod } from '../types/trading-reward-aggregation-types';

export default class TradingRewardAggregationModel extends Model {
  static get tableName() {
    return 'trading_reward_aggregations';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    wallets: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'wallet-model'),
      join: {
        from: 'trading_reward_aggregations.address',
        to: 'wallets.address',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'trading_reward_aggregations.startedAtHeight',
        to: 'blocks.height',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id', // Generated from `address` and `startedAt` and `period`
        'address',
        'startedAt',
        'startedAtHeight',
        'period',
        'amount', // amount of token rewards earned by address in the period starting with startedAt
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        address: { type: 'string' },
        startedAt: { type: 'string', format: 'date-time' }, // Inclusive
        startedAtHeight: { type: 'string', pattern: IntegerPattern }, // Inclusive
        endedAt: { type: ['string', 'null'], format: 'date-time' }, // Exclusive
        endedAtHeight: { type: ['string', 'null'], pattern: IntegerPattern }, // Inclusive
        period: { type: 'string', enum: [...Object.values(TradingRewardAggregationPeriod)] },
        amount: { type: 'string', pattern: NonNegativeNumericPattern },
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
      startedAt: 'date-time',
      startedAtHeight: 'string',
      endedAt: 'date-time',
      endedAtHeight: 'string',
      period: 'string',
      amount: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  id!: string;

  address!: string;

  startedAt!: IsoString;

  startedAtHeight!: string;

  endedAt!: IsoString;

  endedAtHeight!: string;

  period!: TradingRewardAggregationPeriod;

  amount!: string;
}
