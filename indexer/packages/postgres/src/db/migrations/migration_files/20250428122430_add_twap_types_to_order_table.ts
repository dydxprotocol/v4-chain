import * as Knex from 'knex';

import { formatAlterTableEnumSql } from '../helpers';

export async function up(knex: Knex): Promise<void> {
    return knex.raw(formatAlterTableEnumSql(
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
}

export async function down(knex: Knex): Promise<void> {
    return knex.raw(formatAlterTableEnumSql(
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
}
