import { Model } from 'objection';

import UpsertQueryBuilder from '../query-builders/upsert';

export default class BaseModel extends Model {
  static get QueryBuilder() {
    return UpsertQueryBuilder;
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;
}
