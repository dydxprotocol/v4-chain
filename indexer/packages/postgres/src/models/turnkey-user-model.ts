import BaseModel from './base-model';

export default class TurnkeyUserModel extends BaseModel {
  static get tableName() {
    return 'turnkey_users';
  }

  static get idColumn() {
    return 'suborg_id';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'suborg_id',
        'svm_address',
        'evm_address',
        'salt',
        'created_at',
      ],
      properties: {
        suborg_id: { type: 'string' },
        email: { type: ['string', 'null'] },
        svm_address: { type: 'string' },
        evm_address: { type: 'string' },
        smart_account_address: { type: ['string', 'null'] },
        salt: { type: 'string' },
        dydx_address: { type: ['string', 'null'] },
        created_at: { type: 'string' },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   */
  static get sqlToJsonConversions() {
    return {
      suborg_id: 'string',
      email: 'stringOrNull',
      svm_address: 'string',
      evm_address: 'string',
      smart_account_address: 'stringOrNull',
      salt: 'string',
      dydx_address: 'stringOrNull',
      created_at: 'string',
    };
  }

  suborg_id!: string;

  email?: string | undefined;

  svm_address!: string;

  evm_address!: string;

  smart_account_address?: string | undefined;

  salt!: string;

  dydx_address?: string | undefined;

  created_at!: string;
}
