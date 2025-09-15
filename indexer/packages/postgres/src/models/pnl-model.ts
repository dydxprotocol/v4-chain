import path from 'path';

import { Model } from 'objection';

import { IntegerPattern, NumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class PnlModel extends Model {
  static get tableName() {
    return 'pnl';
  }

  static get idColumn() {
    return ['subaccountId', 'createdAt'];
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'pnl.subaccountId',
        to: 'subaccounts.id',
      },
    },
    block_height: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'block-model'),
      join: {
        from: 'pnl.createdAtHeight',
        to: 'blocks.blockHeight',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'subaccountId',
        'createdAt',
        'createdAtHeight',
        'equity',
        'netTransfers',
        'totalPnl',
      ],
      properties: {
        subaccountId: { type: 'string', format: 'uuid' },
        createdAt: { type: 'string', format: 'date-time' },
        createdAtHeight: { type: 'string', pattern: IntegerPattern },
        equity: { type: 'string', pattern: NumericPattern },
        netTransfers: { type: 'string', pattern: NumericPattern },
        totalPnl: { type: 'string', pattern: NumericPattern },
      },
    };
  }

  subaccountId!: string;
  createdAt!: IsoString;
  createdAtHeight!: string;
  equity!: string;
  netTransfers!: string;
  totalPnl!: string;
}
