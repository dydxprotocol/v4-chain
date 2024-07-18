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
    return ['subaccountId', 'timeSpan'];
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'leaderboard_pnl.subaccountId',
        to: 'subaccounts.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'subaccountId',
        'timeSpan',
        'pnl',
        'currentEquity',
        'rank',
      ],
      properties: {
        subaccountId: { type: 'string' },
        timeSpan: { type: 'string' },
        pnl: { type: 'string', pattern: NumericPattern },
        currentEquity: { type: 'string', pattern: NumericPattern },
        rank: { type: 'integer' },
      },
    };
  }

  subaccountId!: string;

  timeSpan!: string;

  QueryBuilderType!: UpsertQueryBuilder<this>;

  pnl!: string;

  currentEquity!: string;

  rank!: number;
}
