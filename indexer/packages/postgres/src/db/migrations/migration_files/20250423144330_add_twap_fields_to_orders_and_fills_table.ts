import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.string('duration').nullable();
    table.string('interval').nullable();
    table.string('priceTolerance').nullable();
  });

  await knex.raw(formatAlterTableEnumSql(
    'orders',
    'type',
    [
      'LIMIT',
      'MARKET',
      'STOP_LIMIT',
      'STOP_MARKET',
      'TRAILING_STOP',
      'TAKE_PROFIT',
      'TAKE_PROFIT_MARKET',
      'LIQUIDATED',
      'LIQUIDATION',
      'HARD_TRADE',
      'FAILED_HARD_TRADE',
      'TRANSFER_PLACEHOLDER',
      'TWAP',
      'TWAP_SUBORDER',
    ],
  ));

  await knex.raw(formatAlterTableEnumSql(
    'fills',
    'type',
    ['LIMIT', 'LIQUIDATED', 'LIQUIDATION', 'DELEVERAGED', 'OFFSETTING', 'TWAP_SUBORDER'],
  ));
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.alterTable('orders', (table) => {
    table.dropColumn('duration');
    table.dropColumn('interval');
    table.dropColumn('priceTolerance');
  });

  await knex.raw(formatAlterTableEnumSql(
    'orders',
    'type',
    [
      'LIMIT',
      'MARKET',
      'STOP_LIMIT',
      'STOP_MARKET',
      'TRAILING_STOP',
      'TAKE_PROFIT',
      'TAKE_PROFIT_MARKET',
      'LIQUIDATED',
      'LIQUIDATION',
      'HARD_TRADE',
      'FAILED_HARD_TRADE',
      'TRANSFER_PLACEHOLDER',
    ],
  ));

  await knex.raw(formatAlterTableEnumSql(
    'fills',
    'type',
    ['LIMIT', 'LIQUIDATED', 'LIQUIDATION', 'DELEVERAGED', 'OFFSETTING'],
  ));
}
