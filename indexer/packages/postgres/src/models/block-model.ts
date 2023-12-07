import path from 'path';

import { Model } from 'objection';

import { IntegerPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class BlockModel extends Model {
  static get tableName() {
    return 'blocks';
  }

  static get idColumn() {
    return 'blockHeight';
  }

  static relationMappings = {
    fills: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'fill-model'),
      join: {
        from: 'blocks.blockHeight',
        to: 'fills.createdAtHeight',
      },
    },
    transfers: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'transfer-model'),
      join: {
        from: 'blocks.blockHeight',
        to: 'transfers.createdAtHeight',
      },
    },
    tendermintEvents: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'tendermint-event-model'),
      join: {
        from: 'blocks.blockHeight',
        to: 'tendermint_events.blockHeight',
      },
    },
    oraclePrices: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'oracle-price-model'),
      join: {
        from: 'blocks.blockHeight',
        to: 'oracle_prices.effectiveAtHeight',
      },
    },
    tradingRewardAggregations: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'trading-reward-aggregation-model'),
      join: {
        from: 'blocks.blockHeight',
        to: 'trading_reward_aggregations.startedAtHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'blockHeight',
        'time',
      ],
      properties: {
        blockHeight: { type: 'string', pattern: IntegerPattern },
        time: { type: 'string', format: 'date-time' },
      },
    };
  }

  blockHeight!: string;

  time!: IsoString;
}
