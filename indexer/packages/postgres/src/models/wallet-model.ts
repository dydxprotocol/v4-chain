import { NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class WalletModel extends BaseModel {
  static get tableName() {
    return 'wallets';
  }

  static get idColumn() {
    return 'address';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'totalTradingRewards',
      ],
      properties: {
        address: { type: 'string' },
        totalTradingRewards: { type: 'string', pattern: NumericPattern },
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
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  totalTradingRewards!: string;
}
