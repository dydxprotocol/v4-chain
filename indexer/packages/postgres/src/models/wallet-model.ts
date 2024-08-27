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
        'isWhitelistAffiliate',
      ],
      properties: {
        address: { type: 'string' },
        totalTradingRewards: { type: 'string', pattern: NonNegativeNumericPattern },
        totalVolume: { type: 'string', pattern: NonNegativeNumericPattern },
        isWhitelistAffiliate: { type: 'boolean'},
      },
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  totalTradingRewards!: string;

  totalVolume!: string;

  isWhitelistAffiliate!: boolean;
}
