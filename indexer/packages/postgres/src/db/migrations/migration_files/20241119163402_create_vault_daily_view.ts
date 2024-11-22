import * as Knex from 'knex';

const RAW_VAULTS_PNL_DAILY_QUERY: string = `
CREATE MATERIALIZED VIEW IF NOT EXISTS vaults_daily_pnl AS WITH vault_subaccounts AS
(
       SELECT subaccounts.id
       FROM   vaults,
              subaccounts
       WHERE  vaults.address = subaccounts.address
       AND    subaccounts."subaccountNumber" = 0), pnl_subaccounts AS
(
       SELECT *
       FROM   vault_subaccounts
       UNION
       SELECT id
       FROM   subaccounts
       WHERE  address = 'dydx18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylxqje4r'
       AND    "subaccountNumber" = 0)
SELECT   "id",
         "subaccountId",
         "equity",
         "totalPnl",
         "netTransfers",
         "createdAt",
         "blockHeight",
         "blockTime"
FROM     (
                  SELECT   pnl_ticks.*,
                           ROW_NUMBER() OVER ( partition BY "subaccountId", DATE_TRUNC( 'day', "blockTime" ) ORDER BY "blockTime" ) AS r
                  FROM     pnl_ticks
                  WHERE    "subaccountId" IN
                           (
                                  SELECT *
                                  FROM   pnl_subaccounts)
                  AND      "blockTime" >= NOW() - interval '7776000 second' ) AS pnl_intervals
WHERE    r = 1
ORDER BY "subaccountId";
`;

export async function up(knex: Knex): Promise<void> {
  await knex.raw(RAW_VAULTS_PNL_DAILY_QUERY);
  await knex.raw('CREATE UNIQUE INDEX ON vaults_daily_pnl (id);');
}

export async function down(knex: Knex): Promise<void> {
  await knex.raw('DROP MATERIALIZED VIEW IF EXISTS vaults_daily_pnl;');
}
