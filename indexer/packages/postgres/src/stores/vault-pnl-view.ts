import { DateTime } from 'luxon';
import { knexReadReplica } from '../helpers/knex';
import { rawQuery } from '../helpers/stores-helpers';
import { PnlFromDatabase, PnlInterval } from '../types';

const VAULT_HOURLY_PNL_VIEW: string = 'vaults_hourly_pnl_v2';
const VAULT_DAILY_PNL_VIEW: string = 'vaults_daily_pnl_v2';

/**
 * Refresh the hourly vault PNL materialized view.
 */
export async function refreshHourlyView(): Promise<void> {
  await rawQuery(
    `REFRESH MATERIALIZED VIEW CONCURRENTLY ${VAULT_HOURLY_PNL_VIEW}`,
    {
      readReplica: false,
    },
  );
}

/**
 * Refresh the daily vault PNL materialized view.
 */
export async function refreshDailyView(): Promise<void> {
  await rawQuery(
    `REFRESH MATERIALIZED VIEW CONCURRENTLY ${VAULT_DAILY_PNL_VIEW}`,
    {
      readReplica: false,
    },
  );
}

/**
 * Get vault PNL data for a given interval and time window.
 * 
 * @param interval - The PNL tick interval (hour or day)
 * @param timeWindowSeconds - The time window in seconds
 * @param earliestDate - The earliest date to fetch data from
 * @returns Array of vault PNL records
 */
export async function getVaultsPnl(
  interval: PnlInterval,
  timeWindowSeconds: number,
  earliestDate: DateTime,
): Promise<PnlFromDatabase[]> {
  const viewName: string = interval === PnlInterval.hour
    ? VAULT_HOURLY_PNL_VIEW
    : VAULT_DAILY_PNL_VIEW;

  const result: {
    rows: PnlFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
      SELECT
        "subaccountId",
        "equity",
        "totalPnl",
        "netTransfers",
        "createdAt",
        "createdAtHeight"
      FROM ${viewName}
      WHERE
        "createdAt" >= '${earliestDate.toUTC().toISO()}'::timestamp AND
        "createdAt" > NOW() - INTERVAL '${timeWindowSeconds} second'
      ORDER BY "subaccountId", "createdAt";
    `,
  ) as unknown as {
    rows: PnlFromDatabase[],
  };

  return result.rows;
}

/**
 * Get the latest vault PNL snapshot for each vault.
 * 
 * @returns Array of latest vault PNL records, one per vault
 */
export async function getLatestVaultPnl(): Promise<PnlFromDatabase[]> {
  const result: {
    rows: PnlFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
      SELECT DISTINCT ON ("subaccountId")
        "subaccountId",
        "equity",
        "totalPnl",
        "netTransfers",
        "createdAt",
        "createdAtHeight"
      FROM ${VAULT_HOURLY_PNL_VIEW}
      ORDER BY "subaccountId", "createdAt" DESC;
    `,
  ) as unknown as {
    rows: PnlFromDatabase[],
  };

  return result.rows;
}