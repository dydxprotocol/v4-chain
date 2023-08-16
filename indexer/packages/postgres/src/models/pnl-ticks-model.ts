import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class PnlTicksModel extends Model {
  static get tableName() {
    return 'pnl_ticks';
  }

  static get idColumn() {
    return 'id';
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'pnl_ticks.subaccountId',
        to: 'subaccounts.id',
      },
    },
    block_height: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'pnl_ticks.blockHeight',
        to: 'blocks.blockHeight',
      },
    },
    block_time: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'pnl_ticks.blockTime',
        to: 'blocks.time',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'id',
        'subaccountId',
        'equity',
        'totalPnl',
        'netTransfers',
        'createdAt',
        'blockHeight',
        'blockTime',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        subaccountId: { type: 'string', format: 'uuid' },
        equity: { type: 'string', pattern: NumericPattern },
        totalPnl: { type: 'string', pattern: NumericPattern },
        netTransfers: { type: 'string', pattern: NumericPattern },
        createdAt: { type: 'string', format: 'date-time' },
        blockHeight: { type: 'string', pattern: IntegerPattern },
        blockTime: { type: 'string', format: 'date-time' },
      },
    };
  }

  id!: string;

  subaccountId!: string;

  equity!: string;

  totalPnl!: string;

  netTransfers!: string;

  createdAt!: IsoString;

  blockHeight!: string;

  blockTime!: IsoString;
}
