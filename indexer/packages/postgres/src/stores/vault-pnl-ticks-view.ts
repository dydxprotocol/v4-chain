import { DateTime } from 'luxon';

import { knexReadReplica } from '../helpers/knex';
import { rawQuery } from '../helpers/stores-helpers';
import {
  PnlTickInterval,
  PnlTicksFromDatabase,
} from '../types';

const VAULT_HOURLY_PNL_VIEW: string = 'vaults_hourly_pnl';
const VAULT_DAILY_PNL_VIEW: string = 'vaults_daily_pnl';

export async function refreshHourlyView(): Promise<void> {
  await rawQuery(
    `REFRESH MATERIALIZED VIEW CONCURRENTLY ${VAULT_HOURLY_PNL_VIEW}`,
    {
      readReplica: false,
    },
  );
}

export async function refreshDailyView(): Promise<void> {
  await rawQuery(
    `REFRESH MATERIALIZED VIEW CONCURRENTLY ${VAULT_DAILY_PNL_VIEW}`,
    {
      readReplica: false,
    },
  );
}

export async function getVaultsPnl(
  interval: PnlTickInterval,
  timeWindowSeconds: number,
  earliestDate: DateTime,
): Promise<PnlTicksFromDatabase[]> {
  let viewName: string = VAULT_DAILY_PNL_VIEW;
  if (interval === PnlTickInterval.hour) {
    viewName = VAULT_HOURLY_PNL_VIEW;
  }
  const result: {
    rows: PnlTicksFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT
      "id",
      "subaccountId",
      "equity",
      "totalPnl",
      "netTransfers",
      "createdAt",
      "blockHeight",
      "blockTime"
    FROM ${viewName}
    WHERE
      "blockTime" >= '${earliestDate.toUTC().toISO()}'::timestamp AND
      "blockTime" > NOW() - INTERVAL '${timeWindowSeconds} second'
    ORDER BY "subaccountId", "blockTime";
    `,
  ) as unknown as {
    rows: PnlTicksFromDatabase[],
  };

  return result.rows;
}

export async function getLatestVaultPnl(): Promise<PnlTicksFromDatabase[]> {
  const result: {
    rows: PnlTicksFromDatabase[],
  } = await knexReadReplica.getConnection().raw(
    `
    SELECT DISTINCT ON ("subaccountId")
      "id",
      "subaccountId",
      "equity",
      "totalPnl",
      "netTransfers",
      "createdAt",
      "blockHeight",
      "blockTime"
    FROM ${VAULT_HOURLY_PNL_VIEW}
    ORDER BY "subaccountId", "blockTime" DESC;
    `,
  ) as unknown as {
    rows: PnlTicksFromDatabase[],
  };

  return result.rows;
}
