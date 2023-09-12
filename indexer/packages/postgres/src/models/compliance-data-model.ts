import { Model } from 'objection';

import { NumericPattern } from '../lib/validators';
import { IsoString } from '../types';

export default class ComplianceDataModel extends Model {
  static get tableName() {
    return 'compliance_data';
  }

  static get idColumn() {
    return 'address';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'sanctioned',
        'updatedAt',
      ],
      properties: {
        address: { type: 'string' },
        chain: { type: ['string', 'null'], default: null },
        sanctioned: { type: 'boolean' },
        riskScore: { type: ['string', 'null'], pattern: NumericPattern, default: null },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  address!: string;

  chain?: string;

  sanctioned!: boolean;

  riskScore?: string;

  updatedAt!: IsoString;
}
