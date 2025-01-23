import {
  logger,
  stats,
  InfoObject,
} from '@dydxprotocol-indexer/base';
import Athena from 'aws-sdk/clients/athena';
import RDS from 'aws-sdk/clients/rds';
import S3 from 'aws-sdk/clients/s3';
import { DateTime } from 'luxon';

import config from '../config';
import {
  checkIfExportJobToS3IsOngoing,
  checkIfTableExistsInAthena,
  checkIfS3ObjectExists,
  getMostRecentDBSnapshotIdentifier,
  startExportTask,
  startAthenaQuery,
} from '../helpers/aws';
import { AthenaTableDDLQueries } from '../helpers/types';
import * as athenaAffiliateInfo from '../lib/athena-ddl-tables/affiliate_info';
import * as athenaAffiliateReferredUsers from '../lib/athena-ddl-tables/affiliate_referred_users';
import * as athenaAssetPositions from '../lib/athena-ddl-tables/asset_positions';
import * as athenaAssets from '../lib/athena-ddl-tables/assets';
import * as athenaBlocks from '../lib/athena-ddl-tables/blocks';
import * as athenaCandles from '../lib/athena-ddl-tables/candles';
import * as athenaFills from '../lib/athena-ddl-tables/fills';
import * as athenaFundingIndexUpdates from '../lib/athena-ddl-tables/funding_index_updates';
import * as athenaLiquidityTiers from '../lib/athena-ddl-tables/liquidity_tiers';
import * as athenaMarkets from '../lib/athena-ddl-tables/markets';
import * as athenaOraclePrices from '../lib/athena-ddl-tables/oracle_prices';
import * as athenaOrders from '../lib/athena-ddl-tables/orders';
import * as athenaPerpetualMarkets from '../lib/athena-ddl-tables/perpetual_markets';
import * as athenaPerpetualPositions from '../lib/athena-ddl-tables/perpetual_positions';
import * as athenaPnlTicks from '../lib/athena-ddl-tables/pnl_ticks';
import * as athenaSubaccountUsernames from '../lib/athena-ddl-tables/subaccount_usernames';
import * as athenaSubaccounts from '../lib/athena-ddl-tables/subaccounts';
import * as athenaTendermintEvents from '../lib/athena-ddl-tables/tendermint_events';
import * as athenaTradingRewardAggregations from '../lib/athena-ddl-tables/trading_reward_aggregations';
import * as athenaTradingRewards from '../lib/athena-ddl-tables/trading_rewards';
import * as athenaTransfers from '../lib/athena-ddl-tables/transfers';
import * as athenaVaults from '../lib/athena-ddl-tables/vaults';
import * as athenaWallets from '../lib/athena-ddl-tables/wallets';

export const tablesToAddToAthena: { [table: string]: AthenaTableDDLQueries } = {
  asset_positions: athenaAssetPositions,
  assets: athenaAssets,
  blocks: athenaBlocks,
  candles: athenaCandles,
  fills: athenaFills,
  funding_index_updates: athenaFundingIndexUpdates,
  markets: athenaMarkets,
  oracle_prices: athenaOraclePrices,
  orders: athenaOrders,
  perpetual_markets: athenaPerpetualMarkets,
  perpetual_positions: athenaPerpetualPositions,
  pnl_ticks: athenaPnlTicks,
  subaccounts: athenaSubaccounts,
  tendermint_events: athenaTendermintEvents,
  trading_rewards: athenaTradingRewards,
  trading_reward_aggregation: athenaTradingRewardAggregations,
  transfers: athenaTransfers,
  liquidity_tiers: athenaLiquidityTiers,
  wallets: athenaWallets,
  affiliate_info: athenaAffiliateInfo,
  affiliate_referred_users: athenaAffiliateReferredUsers,
  vaults: athenaVaults,
  subaccount_usernames: athenaSubaccountUsernames,
};

const statStart: string = `${config.SERVICE_NAME}.update_research_environment`;

export default async function runTask(): Promise<void> {
  const at: string = 'update-research-environment#runTask';

  const rds: RDS = new RDS();

  // get most recent rds snapshot
  const startDescribe: number = Date.now();
  const dateString: string = DateTime.utc().toFormat('yyyy-MM-dd');
  const mostRecentSnapshot: string = await getMostRecentDBSnapshotIdentifier(
    rds,
    undefined,
    config.FAST_SYNC_SNAPSHOT_IDENTIFIER_PREFIX,
  ) as string;
  stats.timing(`${statStart}.describe_rds_snapshots`, Date.now() - startDescribe);

  // dev example: rds:dev-indexer-apne1-db-2023-06-25-18-34
  const s3Date: string = mostRecentSnapshot.split(config.RDS_INSTANCE_NAME)[1].slice(1);
  const s3: S3 = new S3();

  // check if s3 object exists
  const startS3Check: number = Date.now();
  const s3ObjectExists: boolean = await checkIfS3ObjectExists(s3, s3Date);
  stats.timing(`${statStart}.checkS3Object`, Date.now() - startS3Check);

  const rdsExportIdentifier: string = `${config.RDS_INSTANCE_NAME}-${s3Date}`;

  // If the s3 object exists, attempt to add Athena tables or if we are skipping for test purposes
  if (s3ObjectExists || config.SKIP_TO_ATHENA_TABLE_WRITING) {
    logger.info({
      at,
      dateString,
      message: 'S3 object exists. Creating athena tables.',
    });
    const [year, month, day]: string[] = dateString.split('-');
    const athenaDate: string = `${year}_${month}_${day}`;
    await createAthenaTablesIfNotExists(athenaDate, rdsExportIdentifier);

    return;
  }

  // if we haven't created the object, check if it is being created
  const rdsExportCheck: number = Date.now();
  const exportJobOngoing: boolean = await checkIfExportJobToS3IsOngoing(rds, rdsExportIdentifier);
  stats.timing(`${statStart}.checkRdsExport`, Date.now() - rdsExportCheck);

  if (exportJobOngoing) {
    logger.info({
      at,
      dateString,
      message: 'Will wait for export job to finish',
    });
    return;
  }
  // start Export Job if S3 Object does not exist
  const startExport: number = Date.now();
  try {
    const exportData: RDS.ExportTask = await startExportTask(rds, rdsExportIdentifier);

    logger.info({
      at,
      message: 'Started an export task',
      exportData,
    });
  } catch (error) { // TODO handle this by finding the most recent snapshot earlier
    const message: InfoObject = {
      at,
      message: 'export to S3 failed',
      error,
    };

    if (error.name === 'DBSnapshotNotFound') {
      stats.increment(`${statStart}.no_s3_snapshot`, 1);

      logger.info(message);
      return;
    }

    logger.error(message);
  } finally {
    stats.timing(`${statStart}.rdsSnapshotExport`, Date.now() - startExport);
  }
}

async function createAthenaTablesIfNotExists(
  athenaDate: string,
  rdsExportIdentifier: string,
): Promise<void> {
  const athena: Athena = new Athena();

  const start: number = Date.now();
  const athenaTableNames: string[] = Object.keys(tablesToAddToAthena);
  for (let i = 0; i < athenaTableNames.length; i++) {
    await createAthenaTableIfNotExists({
      athenaDate,
      rdsExportIdentifier,
      athena,
      table: athenaTableNames[i],
    });
  }

  stats.timing(`${statStart}.check_for_and_write_athena_tables`, Date.now() - start);
}

async function createAthenaTableIfNotExists({
  athenaDate,
  rdsExportIdentifier,
  table,
  athena,
}: {
  athenaDate: string,
  rdsExportIdentifier: string,
  table: string,
  athena: Athena,
}): Promise<void> {
  const at: string = 'update-research-environment#potentiallyAddAthenaTables';
  try {
    // try to add raw table if it does not exist
    const rawTable: string = `${athenaDate}_raw_${table}`;
    const rawTableExists: boolean = await checkIfTableExistsInAthena(athena, rawTable);
    logger.info({
      at,
      message: 'checkIfTableExistsInAthena',
      rawTableExists,
      table,
    });
    if (!rawTableExists) {
      const rawTableCreationSql: string = tablesToAddToAthena[table].generateRawTable(
        athenaDate,
        rdsExportIdentifier,
      );
      logger.info({
        at,
        message: 'raw table does not exist. Creating raw athena table',
        table,
        rawTableCreationSql,
      });
      const data: Athena.StartQueryExecutionOutput = await startAthenaQuery(
        athena,
        {
          query: rawTableCreationSql,
          timestamp: athenaDate,
        },
      );

      logger.info({
        at,
        message: 'Added raw table',
        table,
        data,
      });
    }

    // try to add queryable table if it does not exist
    const queryableTable = `${athenaDate}_${table}`;
    const [
      tableExists,
      rawTableExistsAfterWriting,
    ]: boolean[] = await Promise.all([
      checkIfTableExistsInAthena(athena, queryableTable),
      checkIfTableExistsInAthena(athena, rawTable),
    ]);

    logger.info({
      at,
      message: 'Checking if queryable table exists',
      tableExists,
      rawTableExistsAfterWriting,
      table,
    });
    if (!tableExists && rawTableExistsAfterWriting) {
      const tableCreationSql: string = tablesToAddToAthena[table].generateTable(athenaDate);
      logger.info({
        at,
        message: 'table does not exist. Creating athena table',
        table,
        tableCreationSql,
      });
      const data: Athena.StartQueryExecutionOutput = await startAthenaQuery(
        athena,
        {
          query: tableCreationSql,
          timestamp: athenaDate,
        },
      );

      logger.info({
        at,
        message: 'Added queryable table',
        table,
        data,
      });
    }

  } catch (error) {
    logger.error({
      at,
      message: 'failed to check for tables or add tables',
      table,
      error,
    });
  }
}
