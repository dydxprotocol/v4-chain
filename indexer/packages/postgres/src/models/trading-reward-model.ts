import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import { IsoString } from '../types';

export default class TradingRewardModel extends Model {
  static get tableName() {
    return 'trading_rewards';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    wallet: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'wallet-model'),
      join: {
        from: 'trading_rewards.address',
        to: 'wallets.address',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'address',
        'blockTime',
        'blockHeight',
        'amount',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        address: { type: 'string' },
        blockTime: { type: 'string', format: 'date-time' },
        blockHeight: { type: 'string', pattern: IntegerPattern },
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
      blockTime: 'date-time',
      blockHeight: 'string',
      amount: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  id!: string;

  address!: string;

  blockTime!: IsoString;

  blockHeight!: string;

  amount!: string;
}
