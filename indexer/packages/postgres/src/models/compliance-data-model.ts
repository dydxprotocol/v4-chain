import { NumericPattern } from '../lib/validators';
import UpsertQueryBuilder from '../query-builders/upsert';
import { ComplianceProvider, IsoString } from '../types';
import BaseModel from './base-model';

export default class ComplianceDataModel extends BaseModel {
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
        'blocked',
        'updatedAt',
      ],
      properties: {
        address: { type: 'string' },
        provider: { type: 'string', enum: [...Object.values(ComplianceProvider)] },
        chain: { type: ['string', 'null'], default: null },
        blocked: { type: 'boolean' },
        riskScore: { type: ['string', 'null'], pattern: NumericPattern, default: null },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  provider!: string;

  chain?: string;

  blocked!: boolean;

  riskScore?: string;

  updatedAt!: IsoString;
}
