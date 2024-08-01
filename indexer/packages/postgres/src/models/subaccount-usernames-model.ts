import path from 'path';

import { Model } from 'objection';

export default class SubaccountUsernames extends Model {

  static get tableName() {
    return 'subaccount_usernames';
  }

  static get idColumn() {
    return 'subaccountId';
  }

  static relationMappings = {
    subaccount: {
      relation: Model.BelongsToOneRelation,
      modelClass: path.join(__dirname, 'subaccount-model'),
      join: {
        from: 'subaccount_usernames.subaccountId',
        to: 'subaccounts.id',
      },
    },
  };

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'username',
        'subaccountId'],
      properties: {
        username: { type: 'string' },
        subaccountId: { type: 'string' },
      },
    };
  }

  username!: string;

  subaccountId!: string;
}
