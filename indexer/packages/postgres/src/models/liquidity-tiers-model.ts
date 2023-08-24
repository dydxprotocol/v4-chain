import { IntegerPattern, NonNegativeNumericPattern } from '../lib/validators';
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
        'basePositionNotional',
      ],
      properties: {
        id: { type: 'integer' },
        name: { type: 'string' },
        initialMarginPpm: { type: 'string', pattern: IntegerPattern },
        maintenanceFractionPpm: { type: 'string', pattern: IntegerPattern },
        basePositionNotional: { type: 'string', pattern: NonNegativeNumericPattern },
      },
    };
  }

  id!: number;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  name!: string;

  initialMarginPpm!: string;

  maintenanceFractionPpm!: string;

  basePositionNotional!: string;
}
