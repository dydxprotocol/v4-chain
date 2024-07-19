import Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex
    .schema
    .createTable('perpetual_markets', (table) => {
      table.bigInteger('id').primary();
      table.bigInteger('clobPairId').notNullable();
      table.string('ticker').notNullable();
      table.integer('marketId').notNullable();
      table.enum(
        'status',
        [
          'ACTIVE',
          'PAUSED',
          'CANCEL_ONLY',
          'POST_ONLY',
        ],
      ).notNullable();
      table.string('baseAsset').notNullable();
      table.string('quoteAsset').notNullable();
      table.decimal('lastPrice', null).notNullable();
      table.decimal('priceChange24H', null).notNullable();
      table.decimal('volume24H', null).notNullable();
      table.integer('trades24H').notNullable();
      table.decimal('nextFundingRate', null).notNullable();
      table.decimal('initialMarginFraction', null).notNullable();
      table.decimal('incrementalInitialMarginFraction', null).notNullable();
      table.decimal('maintenanceMarginFraction', null).notNullable();
      table.decimal('basePositionSize', null).notNullable();
      table.decimal('incrementalPositionSize', null).notNullable();
      table.decimal('maxPositionSize', null).notNullable();
      table.decimal('openInterest', null).notNullable();
      table.integer('quantumConversionExponent').notNullable();
      table.integer('atomicResolution').notNullable();
      table.integer('subticksPerTick').notNullable();
      table.integer('minOrderBaseQuantums').notNullable();
      table.integer('stepBaseQuantums').notNullable();
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('perpetual_markets');
}
