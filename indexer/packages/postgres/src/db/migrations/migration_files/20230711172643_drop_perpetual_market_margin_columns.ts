import * as Knex from 'knex';

import { getPerpetualMarketMarginRestoreSql } from '../../helpers';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.dropColumn('initialMarginFraction');
    table.dropColumn('incrementalInitialMarginFraction');
    table.dropColumn('maintenanceMarginFraction');
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('perpetual_markets', (table) => {
    table.decimal('initialMarginFraction', null).notNullable().defaultTo('0');
    table.decimal('incrementalInitialMarginFraction', null).notNullable().defaultTo('0');
    table.decimal('maintenanceMarginFraction', null).notNullable().defaultTo('0');
  });
  // Update perpetual_markets table with margin values from genesis.
  const updateSql: string[] = getPerpetualMarketMarginRestoreSql();
  for (const sql of updateSql) {
    await knex.raw(sql);
  }
}
