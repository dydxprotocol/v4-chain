import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'orders',
    'status',
    ['OPEN', 'FILLED', 'CANCELED', 'BEST_EFFORT_CANCELED', 'UNTRIGGERED'],
  ));
}

export async function down(knex: Knex): Promise<void> {
  return knex.raw(formatAlterTableEnumSql(
    'orders',
    'status',
    ['OPEN', 'FILLED', 'CANCELED', 'BEST_EFFORT_CANCELED'],
  ));
}
