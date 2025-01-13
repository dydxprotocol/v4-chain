import { NonNegativeNumericPattern, NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class AffiliateRefereeStatsModel extends BaseModel {
  static get tableName() {
    return 'affiliate_referee_stats';
  }

  static get idColumn() {
    return 'refereeAddress';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'refereeAddress',
        'affiliateAddress',
        'affiliateEarnings',
        'referredMakerTrades',
        'referredTakerTrades',
        'referredTotalVolume',
        'referralBlockHeight',
        'referredTakerFees',
        'referredMakerFees',
        'referredMakerRebates',
        'referredLiquidationFees',
      ],
      properties: {
        refereeAddress: { type: 'string' },
        affiliateAddress: { type: 'string' },
        affiliateEarnings: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerTrades: { type: 'int' },
        referredTakerTrades: { type: 'int' },
        referredTotalVolume: { type: 'string', pattern: NonNegativeNumericPattern },
        referralBlockHeight: { type: 'string', pattern: NonNegativeNumericPattern },
        referredTakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerFees: { type: 'string', pattern: NonNegativeNumericPattern },
        referredMakerRebates: { type: 'string', pattern: NumericPattern },
        referredLiquidationFees: { type: 'string', pattern: NonNegativeNumericPattern },
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
      refereeAddress: 'string',
      affiliateAddress: 'string',
      affiliateEarnings: 'string',
      referredMakerTrades: 'int',
      referredTakerTrades: 'int',
      referredTotalVolume: 'string',
      referralBlockHeight: 'string',
      referredTakerFees: 'string',
      referredMakerFees: 'string',
      referredMakerRebates: 'string',
      referredLiquidationFees: 'string',
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  refereeAddress!: string;

  affiliateAddress!: string;

  affiliateEarnings!: string;

  referredMakerTrades!: number;

  referredTakerTrades!: number;

  referredTotalVolume!: string;

  referralBlockHeight!: string;

  referredMakerFees!: string;

  referredTakerFees!: string;

  referredMakerRebates!: string;

  referredLiquidationFees!: string;
}
