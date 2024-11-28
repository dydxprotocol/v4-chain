import { NonNegativeNumericPattern, NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class AffiliatePerRefereeStatsModel extends BaseModel {
  static get tableName() {
    return 'affiliate_per_refree_stats';
  }

  static get idColumn() {
    return ['affiliateAddress', 'refereeAddress'];
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'affiliateAddress',
        'refereeAddress',
        'affiliateEarnings',
        'referredMakerTrades',
        'referredTakerTrades',
        'referredTotalVolume',
        'firstReferralBlockHeight',
        'totalReferredTakerFees',
        'totalReferredMakerFees',
        'totalReferredMakerRebates',
        'totalReferredLiquidationfees',
      ],
      properties: {
        affiliateAddress: { type: 'string' },
        refereeAddress: { type: 'string' },
        affiliateEarnings: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerTrades: { type: 'int' },
        referredTakerTrades: { type: 'int' },
        referredTotalVolume: { type: 'string', pattern: NonNegativeNumericPattern },
        firstReferralBlockHeight: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredTakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredMakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        totalReferredMakerRebates: { type: 'string', pattern: NumericPattern },
        totalReferredLiquidationfees: { type: 'string', pattern: NonNegativeNumericPattern },
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
      affiliateAddress: 'string',
      refereeAddress: 'string',
      affiliateEarnings: 'string',
      referredMakerTrades: 'int',
      referredTakerTrades: 'int',
      referredTotalVolume: 'string',
      firstReferralBlockHeight: 'int',
      totalReferredTakerFees: 'string',
      totalReferredMakerFees: 'string',
      totalReferredMakerRebates: 'string',
      totalReferredLiquidationfees: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  affiliateAddress!: string;

  refereeAddress!: string;

  affiliateEarnings!: string;

  referredMakerTrades!: number;

  referredTakerTrades!: number;

  referredTotalVolume!: string;

  firstReferralBlockHeight!: number;

  totalReferredMakerFees!: string;

  totalReferredTakerFees!: string;

  totalReferredMakerRebates!: string;

  totalReferredLiquidationfees!: string;
}
