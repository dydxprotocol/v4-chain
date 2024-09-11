import UpsertQueryBuilder from '../query-builders/upsert';
import BaseModel from './base-model';

export default class PersistentCacheModel extends BaseModel {
  static get tableName() {
    return 'persistent_cache';
  }

  static get idColumn() {
    return 'key';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: ['key', 'value'],
      properties: {
        key: { type: 'string' },
        value: { type: 'string' },
      },
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  key!: string;

  value!: string;
}
