import { NonNegativeNumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class AffiliateInfoModel extends BaseModel {
  static get tableName() {
    return 'affiliate_info';
  }

  static get idColumn() {
    return 'address';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'affiliateEarnings',
        'referredMakerTrades',
        'referredTakerTrades',
        'totalReferredFees',
        'totalReferredUsers',
        'referredNetProtocolEarnings',
        'firstReferralBlockHeight',
      ],
      properties: {
        address: { type: 'string' },
        affiliateEarnings: { type: 'int', pattern: NonNegativeNumericPattern },
        referredMakerTrades: { type: 'int', pattern: NonNegativeNumericPattern },
        referredTakerTrades: { type: 'int', pattern: NonNegativeNumericPattern },
        totalReferredFees: { type: 'int', pattern: NonNegativeNumericPattern },
        totalReferredUsers: { type: 'int', pattern: NonNegativeNumericPattern },
        referredNetProtocolEarnings: { type: 'int', pattern: NonNegativeNumericPattern },
        firstReferralBlockHeight: { type: 'int', pattern: NonNegativeNumericPattern },
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
      affiliateEarnings: 'int',
      referredMakerTrades: 'int',
      referredTakerTrades: 'int',
      totalReferredFees: 'int',
      totalReferredUsers: 'int',
      referredNetProtocolEarnings: 'int',
      firstReferralBlockHeight: 'int',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  affiliateEarnings!: number;

  referredMakerTrades!: number;

  referredTakerTrades!: number;

  totalReferredFees!: number;

  totalReferredUsers!: number;

  referredNetProtocolEarnings!: number;

  firstReferralBlockHeight!: number;
}
