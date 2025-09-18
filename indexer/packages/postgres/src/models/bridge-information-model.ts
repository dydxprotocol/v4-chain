import BaseModel from './base-model';

export default class BridgeInformationModel extends BaseModel {
  static get tableName() {
    return 'bridge_information';
  }

  static get idColumn() {
    return 'id';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'from_address',
        'chain_id',
        'amount',
        'created_at',
      ],
      properties: {
        id: { type: 'string' },
        from_address: { type: 'string' },
        chain_id: { type: 'string' },
        amount: { type: 'string' },
        transaction_hash: { type: ['string', 'null'] },
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
      id: 'string',
      from_address: 'string',
      chain_id: 'string',
      amount: 'string',
      transaction_hash: 'stringOrNull',
      created_at: 'string',
    };
  }

  id!: string;

  from_address!: string;

  chain_id!: string;

  amount!: string;

  transaction_hash?: string | undefined;

  created_at!: string;
}
