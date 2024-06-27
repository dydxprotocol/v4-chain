import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'fills',
    'type',
    ['LIMIT', 'LIQUIDATED', 'LIQUIDATION', 'DELEVERAGED', 'OFFSETTING'],
  ));
}

export async function down(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'fills',
    'type',
    ['MARKET', 'LIMIT', 'LIQUIDATED', 'LIQUIDATION', 'DELEVERAGED', 'OFFSETTING'],
  ));
}
