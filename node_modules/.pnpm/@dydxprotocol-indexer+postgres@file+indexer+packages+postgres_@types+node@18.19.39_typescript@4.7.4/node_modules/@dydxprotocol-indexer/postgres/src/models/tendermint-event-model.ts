import path from 'path';

import { Model } from 'objection';

import { IntegerPattern } from '../lib/validators';

export default class TendermintEventModel extends Model {
  static get tableName() {
    return 'tendermint_events';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    openPerpetualPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'tendermint_events.id',
        to: 'perpetual_positions.openEventId',
      },
    },
    closePerpetualPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'tendermint_events.id',
        to: 'perpetual_positions.closeEventId',
      },
    },
    lastPerpetualPosition: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'perpetual-position-model'),
      join: {
        from: 'tendermint_events.id',
        to: 'perpetual_positions.lastEventId',
      },
    },
    fills: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'fill-model'),
      join: {
        from: 'tendermint_events.id',
        to: 'fills.eventId',
      },
    },
    blocks: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'tendermint_events.blockHeight',
        to: 'blocks.blockHeight',
      },
    },
    transfers: {
      relation: Model.HasManyRelation,
      modelClass: path.join(__dirname, 'transfer-model'),
      join: {
        from: 'tendermint_events.id',
        to: 'transfers.eventId',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'blockHeight',
        'transactionIndex',
        'eventIndex',
      ],
      properties: {
        blockHeight: { type: 'string', pattern: IntegerPattern },
        transactionIndex: { type: 'integer' },
        eventIndex: { type: 'integer' },
      },
    };
  }

  id!: Buffer;

  blockHeight!: string;

  transactionIndex!: number;

  eventIndex!: number;
}
