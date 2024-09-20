import path from 'path';

import { Model } from 'objection';

import {
  IntegerPattern,
  NonNegativeNumericPattern,
  NumericPattern,
} from '../lib/validators';
import {
  PerpetualMarketStatus, PerpetualMarketType,
} from '../types';

export default class PerpetualMarketModel extends Model {
  static get tableName() {
    return 'perpetual_markets';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    perpetualPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'perpetual_markets.id',
        to: 'perpetual_positions.perpetualId',
      },
    },
    market: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'market-model'),
      join: {
        from: 'perpetual_markets.marketId',
        to: 'markets.id',
      },
    },
    liquidityTiers: {
      relation: Model.HasOneRelation,
      modelClass: path.join(__dirname, 'liquidity-tiers-model'),
      join: {
        from: 'perpetual_markets.liquidityTierId',
        to: 'liquidity_tiers.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'clobPairId',
        'ticker',
        'marketId',
        'status',
        'quantumConversionExponent',
        'atomicResolution',
        'subticksPerTick',
        'stepBaseQuantums',
        'liquidityTierId',
        'marketType',
      ],
      properties: {
        id: { type: 'string', pattern: IntegerPattern },
        clobPairId: { type: 'string', pattern: IntegerPattern },
        ticker: { type: 'string' },
        marketId: { type: 'integer' },
        status: { type: 'string', enum: [...Object.values(PerpetualMarketStatus)] },
        priceChange24H: { type: 'string', pattern: NumericPattern, nullable: true }, // Optional
        volume24H: { type: 'string', pattern: NonNegativeNumericPattern, nullable: true }, // Optional
        trades24H: { type: 'integer', nullable: true }, // Optional
        nextFundingRate: { type: 'string', pattern: NumericPattern, nullable: true }, // Optional
        openInterest: { type: 'string', pattern: NumericPattern, nullable: true }, // Optional
        quantumConversionExponent: { type: 'integer' },
        atomicResolution: { type: 'integer' },
        subticksPerTick: { type: 'integer' },
        stepBaseQuantums: { type: 'integer' },
        liquidityTierId: { type: 'integer' },
        marketType: { type: 'string' },
        baseOpenInterest: { type: 'string', pattern: NumericPattern, nullable: true }, // Optional
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
      clobPairId: 'string',
      ticker: 'string',
      marketId: 'integer',
      status: 'string',
      priceChange24H: 'string',
      volume24H: 'string',
      trades24H: 'integer',
      nextFundingRate: 'string',
      openInterest: 'string',
      quantumConversionExponent: 'integer',
      atomicResolution: 'integer',
      subticksPerTick: 'integer',
      stepBaseQuantums: 'integer',
      liquidityTierId: 'integer',
      marketType: 'string',
      baseOpenInterest: 'string',
    };
  }

  id!: string;

  clobPairId!: string;

  ticker!: string;

  marketId!: number;

  status!: PerpetualMarketStatus;

  priceChange24H!: string;

  volume24H!: string;

  trades24H!: number;

  nextFundingRate!: string;

  openInterest!: string;

  quantumConversionExponent!: number;

  atomicResolution!: number;

  subticksPerTick!: number;

  stepBaseQuantums!: number;

  liquidityTierId!: number;

  marketType!: PerpetualMarketType;

  baseOpenInterest!: string;
}
