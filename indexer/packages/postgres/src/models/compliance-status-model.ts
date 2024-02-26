import { Model } from 'objection';

import UpsertQueryBuilder from '../query-builders/upsert';
import { ComplianceReason, ComplianceStatus } from '../types';

export default class ComplianceStatusModel extends Model {
  static get tableName() {
    return 'compliance_status';
  }

  static get idColumn() {
    return 'address';
  }

  static relationMappings = {};

  static get jsonSchema() {
    return {
      type: 'object',
      required: ['address', 'status'],
      properties: {
        address: { type: 'string' },
        status: {
          type: 'string',
          enum: [...Object.values(ComplianceStatus)],
        },
        reason: {
          type: ['string', 'null'],
          enum: [...Object.values(ComplianceReason), null],
          default: null,
        },
        createdAt: { type: 'string', format: 'date-time' },
        updatedAt: { type: 'string', format: 'date-time' },
      },
    };
  }

  QueryBuilderType!: UpsertQueryBuilder<this>;

  address!: string;

  status!: ComplianceStatus;

  reason?: ComplianceReason;

  createdAt!: string;

  updatedAt!: string;
}
