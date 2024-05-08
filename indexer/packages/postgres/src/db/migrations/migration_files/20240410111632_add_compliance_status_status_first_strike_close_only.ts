import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'compliance_status',
    'status',
    ['COMPLIANT', 'FIRST_STRIKE_CLOSE_ONLY', 'FIRST_STRIKE', 'CLOSE_ONLY', 'BLOCKED'],
  ));
}

export async function down(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'compliance_status',
    'status',
    ['COMPLIANT', 'FIRST_STRIKE', 'CLOSE_ONLY', 'BLOCKED'],
  ));
}
