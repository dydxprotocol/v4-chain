import { IntegerPattern, NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class LiquidityTiersModel extends BaseModel {
  static get tableName() {
    return 'liquidity_tiers';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'name',
        'initialMarginPpm',
        'maintenanceFractionPpm',
      ],
      properties: {
        id: { type: 'integer' },
        name: { type: 'string' },
        initialMarginPpm: { type: 'string', pattern: IntegerPattern },
        maintenanceFractionPpm: { type: 'string', pattern: IntegerPattern },
        // Uppper cap for open interest in human readable format(USDC)
        openInterestLowerCap: { type: ['string', 'null'], pattern: NumericPattern },
        // Lower cap for open interest in human readable format(USDC)
        openInterestUpperCap: { type: ['string', 'null'], pattern: NumericPattern },
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
      id: 'integer',
      name: 'string',
      initialMarginPpm: 'string',
      maintenanceFractionPpm: 'string',
      openInterestLowerCap: 'string',
      openInterestUpperCap: 'string',
    };
  }

  id!: number;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  name!: string;

  initialMarginPpm!: string;

  maintenanceFractionPpm!: string;

  openInterestLowerCap?: string;

  openInterestUpperCap?: string;
}
