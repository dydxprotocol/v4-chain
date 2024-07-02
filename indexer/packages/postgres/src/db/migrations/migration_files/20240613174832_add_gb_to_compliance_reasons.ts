import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'compliance_status',
    'reason',
    ['MANUAL', 'US_GEO', 'CA_GEO', 'GB_GEO', 'SANCTIONED_GEO', 'COMPLIANCE_PROVIDER'],
  ));
}

export async function down(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'compliance_status',
    'reason',
    ['MANUAL', 'US_GEO', 'CA_GEO', 'SANCTIONED_GEO', 'COMPLIANCE_PROVIDER'],
  ));
}
