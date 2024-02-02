import * as Knex from 'knex';

export async function up(knex: Knex): Promise<void> {
  return knex.schema
    .createTable('compliance_status', (table) => {
      table.string('address').primary(); // dYdX V4 chain address
      table.enum(
        'status', // compliance of the address
        [
          'COMPLIANT', // the address is compliant
          'FIRST_STRIKE', // the address has a single-strike
          'CLOSE_ONLY', // the address is in close-only mode
          'BLOCKED', // the address is blocked
        ],
      ).notNullable();
      table.enum(
        'reason',
        [
          'MANUAL', // the address was manually set to the status
          'US_GEO', // the address was set to the status due to connecting from US (restricted geography)
          'CA_GEO', // the address was set to the status due to connecting from CA (restricted geography)
          'SANCTIONED_GEO', // the address was set to the status due to connection from sanctioned geography
          'COMPLIANCE_PROVIDER', // the address was set to the status due to being flagged by a compliance provider
        ],
      ).nullable().defaultTo(null); // null for COMPLIANT addresses
      table.timestamp('createdAt').notNullable().defaultTo(knex.fn.now());
      table.timestamp('updatedAt').notNullable().defaultTo(knex.fn.now());

      table.index(['status']); // needed to search for CLOSE_ONLY addresses
      table.index(['updatedAt']); // needed to search for CLOSE_ONLY addresses
    });
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('compliance_status');
}
