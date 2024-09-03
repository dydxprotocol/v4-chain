import { NonNegativeNumericPattern } from '../lib/validators';
import BaseModel from './base-model';

export default class AffiliateReferredUsersModel extends BaseModel {
  static get tableName() {
    return 'affiliate_referred_users';
  }

  static get idColumn() {
    return 'refereeAddress';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'affiliateAddress',
        'refereeAddress',
        'referredAtBlock',
      ],
      properties: {
        affiliateAddress: { type: 'string' },
        refereeAddress: { type: 'string' },
        referredAtBlock: { type: 'string', pattern: NonNegativeNumericPattern },
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
      referredAtBlock: 'string',
    };
  }

  affiliateAddress!: string;

  refereeAddress!: string;

  referredAtBlock!: string;
}
