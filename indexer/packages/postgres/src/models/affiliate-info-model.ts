import { NonNegativeNumericPattern, NumericPattern } from '../lib/validators';
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
        'totalReferredMakerFees',
        'totalReferredTakerFees',
        'totalReferredMakerRebates',
        'totalReferredUsers',
        'firstReferralBlockHeight',
        'referredTotalVolume',
      ],
      properties: {
        address: { type: 'string' },
        affiliateEarnings: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerTrades: { type: 'int' },
        referredTakerTrades: { type: 'int' },
        totalReferredMakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredTakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredMakerRebates: { type: 'string', pattern: NumericPattern },
        totalReferredUsers: { type: 'int' },
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
      totalReferredMakerFees: 'string',
      totalReferredTakerFees: 'string',
      totalReferredMakerRebates: 'string',
      totalReferredUsers: 'int',
      firstReferralBlockHeight: 'string',
      referredTotalVolume: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  affiliateEarnings!: string;

  referredMakerTrades!: number;

  referredTakerTrades!: number;

  totalReferredMakerFees!: string;

  totalReferredTakerFees!: string;

  totalReferredMakerRebates!: string;

  totalReferredUsers!: number;

  firstReferralBlockHeight!: string;

  referredTotalVolume!: string;
}
