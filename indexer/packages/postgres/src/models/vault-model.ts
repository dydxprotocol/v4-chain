import { IntegerPattern } from '../lib/validators';
import { IsoString, VaultStatus } from '../types';
import BaseModel from './base-model';

export default class VaultModel extends BaseModel {

  static get tableName() {
    return 'vaults';
  }

  static get idColumn() {
    return ['address'];
  }

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'clobPairId',
        'status',
        'createdAt',
        'updatedAt',
      ],
      properties: {
        address: { type: 'string' },
        clobPairId: { type: 'string', pattern: IntegerPattern },
        status: { type: 'string' },
        createdAt: { type: 'string', format: 'date-time' },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  /**
   * A mapping from column name to JSON conversion expected.
   * See getSqlConversionForDydxModelTypes for valid conversions.
   *
   * TODO(IND-239): Ensure that jsonSchema() / sqlToJsonConversions() / model fields match.
   */
  static get sqlToJsonConversions() {
    return {
      address: 'string',
      clobPairId: 'string',
      status: 'string',
      createdAt: 'date-time',
      updatedAt: 'date-time',
    };
  }

  address!: string;

  clobPairId!: string;

  status!: VaultStatus;

  createdAt!: IsoString;

  updatedAt!: IsoString;
}
