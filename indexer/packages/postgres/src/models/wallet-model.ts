import path from 'path';

import { Model } from 'objection';

import { NonNegativeNumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class WalletModel extends BaseModel {
  static get tableName() {
    return 'wallets';
  }

  static get idColumn() {
    return 'address';
  }

  static relationMappings = {
    tradingRewardAggregations: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'trading-reward-aggregation-model'),
      join: {
        from: 'wallets.address',
        to: 'trading_reward_aggregations.address',
      },
    },
    tradingRewards: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'trading-reward-model'),
      join: {
        from: 'wallets.address',
        to: 'trading_rewards.address',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'totalTradingRewards',
        'totalVolume',
      ],
      properties: {
        address: { type: 'string' },
        totalTradingRewards: { type: 'string', pattern: NonNegativeNumericPattern },
        totalVolume: { type: 'string', pattern: NonNegativeNumericPattern },
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
      address: 'string',
      totalTradingRewards: 'string',
      totalVolume: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  totalTradingRewards!: string;

  totalVolume!: string;
}
