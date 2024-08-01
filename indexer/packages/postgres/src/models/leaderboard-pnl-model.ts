import path from 'path';

import { Model } from 'objection';

import { NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class LeaderboardPnlModel extends BaseModel {

  static get tableName() {
    return 'leaderboard_pnl';
  }

  static get idColumn() {
    return ['address', 'timeSpan'];
  }

  static relationMappings = {
    wallets: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'wallet-model'),
      join: {
        from: 'leaderboard_pnl.address',
        to: 'wallets.address',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'timeSpan',
        'pnl',
        'currentEquity',
        'rank',
      ],
      properties: {
        address: { type: 'string' },
        timeSpan: { type: 'string' },
        pnl: { type: 'string', pattern: NumericPattern },
        currentEquity: { type: 'string', pattern: NumericPattern },
        rank: { type: 'integer' },
      },
    };
  }

  address!: string;

  timeSpan!: string;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  pnl!: string;

  currentEquity!: string;

  rank!: number;
}
