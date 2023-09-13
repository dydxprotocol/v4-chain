import { Model } from 'objection';

import { NumericPattern } from '../lib/validators';
import { ComplianceProvider, IsoString } from '../types';

export default class ComplianceDataModel extends Model {
  static get tableName() {
    return 'compliance_data';
  }

  static get idColumn() {
    return ['address', 'provider'];
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: [
        'address',
        'provider',
        'sanctioned',
        'updatedAt',
      ],
      properties: {
        address: { type: 'string' },
        provider: { type: 'string', enum: [...Object.values(ComplianceProvider)] },
        chain: { type: ['string', 'null'], default: null },
        sanctioned: { type: 'boolean' },
        riskScore: { type: ['string', 'null'], pattern: NumericPattern, default: null },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  address!: string;

  provider!: string;

  chain?: string;

  sanctioned!: boolean;

  riskScore?: string;

  updatedAt!: IsoString;
}
