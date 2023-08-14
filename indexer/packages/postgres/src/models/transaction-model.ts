import { Model } from 'objection';

import { IntegerPattern } from '../lib/validators';

export default class TransactionModel extends Model {
  static get tableName() {
    return 'transactions';
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
        'blockHeight',
        'transactionIndex',
        'transactionHash',
      ],
      properties: {
        id: { type: 'string', format: 'uuid' },
        blockHeight: { type: 'string', pattern: IntegerPattern },
        transactionIndex: { type: 'integer' },
        transactionHash: { type: 'string' },
      },
    };
  }

  id!: string;

  blockHeight!: string;

  transactionIndex!: number;

  transactionHash!: string;
}
