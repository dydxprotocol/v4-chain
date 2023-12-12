import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NonNegativeNumericPattern, NumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class FundingIndexUpdatesModel extends Model {
  static get tableName() {
    return 'funding_index_updates';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    openTendermintEvent: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'funding_index_updates.eventId',
        to: 'tendermint_events.id',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'funding_index_updates.effectiveAtHeight',
        to: 'blocks.blockHeight',
      },
    },
    perpetualMarkets: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'perpetual-market-model'),
      join: {
        from: 'funding_index_updates.perpetualId',
        to: 'perpetual_markets.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'perpetualId',
        'eventId',
        'rate',
        'oraclePrice',
        'fundingIndex',
        'effectiveAt',
        'effectiveAtHeight',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        perpetualId: { type: 'string', pattern: IntegerPattern },
        rate: { type: 'string', pattern: NumericPattern },  // rate is a decimal string. a 2% rate is "0.02".
        oraclePrice: { type: 'string', pattern: NonNegativeNumericPattern },
        fundingIndex: { type: 'string', pattern: NumericPattern },
        effectiveAt: { type: 'string', format: 'date-time' },
        effectiveAtHeight: { type: 'string', pattern: IntegerPattern },
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
      perpetualId: 'string',
      eventId: 'hex-string',
      rate: 'string',
      oraclePrice: 'string',
      fundingIndex: 'string',
      effectiveAt: 'date-time',
      effectiveAtHeight: 'string',
    };
  }

  id!: string;

  perpetualId!: string;

  eventId!: Buffer;

  rate!: string;

  oraclePrice!: string;

  fundingIndex!: string;

  effectiveAt!: IsoString;

  effectiveAtHeight!: string;
}
