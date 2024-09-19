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
        'referredTotalVolume',
      ],
      properties: {
        address: { type: 'string' },
        affiliateEarnings: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerTrades: { type: 'int' },
        referredTakerTrades: { type: 'int' },
        totalReferredFees: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredUsers: { type: 'int' },
        referredNetProtocolEarnings: { type: 'string', pattern: NonNegativeNumericPattern },
        firstReferralBlockHeight: { type: 'string', pattern: NonNegativeNumericPattern },
        referredTotalVolume: { type: 'string', pattern: NonNegativeNumericPattern },
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
      affiliateEarnings: 'string',
      referredMakerTrades: 'int',
      referredTakerTrades: 'int',
      totalReferredFees: 'string',
      totalReferredUsers: 'int',
      referredNetProtocolEarnings: 'string',
      firstReferralBlockHeight: 'string',
      referredTotalVolume: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  affiliateEarnings!: string;

  referredMakerTrades!: number;

  referredTakerTrades!: number;

  totalReferredFees!: string;

  totalReferredUsers!: number;

  referredNetProtocolEarnings!: string;

  firstReferralBlockHeight!: string;

  referredTotalVolume!: string;
}
