import BaseModel from './base-model';

export default class TurnkeyUserModel extends BaseModel {
  static get tableName() {
    return 'turnkey_users';
  }

  static get idColumn() {
    return 'suborgId';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'suborgId',
        'svmAddress',
        'evmAddress',
        'salt',
        'createdAt',
      ],
      properties: {
        suborgId: { type: 'string' },
        username: { type: ['string', 'null'] },
        email: { type: ['string', 'null'] },
        svmAddress: { type: 'string' },
        evmAddress: { type: 'string' },
        salt: { type: 'string' },
        dydxAddress: { type: ['string', 'null'] },
        createdAt: { type: 'string' },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   */
  static get sqlToJsonConversions() {
    return {
      suborgId: 'string',
      username: 'string',
      email: 'string',
      svmAddress: 'string',
      evmAddress: 'string',
      salt: 'string',
      dydxAddress: 'string',
      createdAt: 'string',
    };
  }

  suborgId!: string;

  username?: string;

  email?: string;

  svmAddress!: string;

  evmAddress!: string;

  salt!: string;

  dydxAddress?: string;

  createdAt!: string;
}
