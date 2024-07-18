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

  static relationMappings = {};

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
