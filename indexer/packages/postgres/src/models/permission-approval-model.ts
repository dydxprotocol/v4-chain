import BaseModel from './base-model';

export default class PermissionApprovalModel extends BaseModel {
  static get tableName() {
    return 'permission_approval';
  }

  static get idColumn() {
    return 'suborg_id';
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'suborg_id',
      ],
      properties: {
        suborg_id: { type: 'string' },
        arbitrum_approval: { type: ['string', 'null'] },
        base_approval: { type: ['string', 'null'] },
        avalanche_approval: { type: ['string', 'null'] },
        optimism_approval: { type: ['string', 'null'] },
        ethereum_approval: { type: ['string', 'null'] },
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
      arbitrum_approval: 'stringOrNull',
      base_approval: 'stringOrNull',
      avalanche_approval: 'stringOrNull',
      optimism_approval: 'stringOrNull',
      ethereum_approval: 'stringOrNull',
    };
  }

  suborg_id!: string;

  arbitrum_approval?: string | undefined;

  base_approval?: string | undefined;

  avalanche_approval?: string | undefined;

  optimism_approval?: string | undefined;

  ethereum_approval?: string | undefined;
}
