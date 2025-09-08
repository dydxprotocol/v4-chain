import BaseModel from './base-model';

export default class PermissionApprovalModel extends BaseModel {
  static get tableName() {
    return 'permission_approval';
  }

  static get idColumn() {
    return ['suborg_id', 'chain_id'];
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'suborg_id',
        'chain_id',
        'approval',
      ],
      properties: {
        suborg_id: { type: 'string' },
        chain_id: { type: 'string' },
        approval: { type: 'string' },
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
      chain_id: 'string',
      approval: 'string',
    };
  }

  suborg_id!: string;

  chain_id!: string;

  approval!: string;
}
