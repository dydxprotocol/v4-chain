import {
  ONE_DAY_IN_MILLISECONDS,
  stats,
} from '@dydxprotocol-indexer/base';
import {
  ComplianceStatusFromDatabase,
  ComplianceStatusTable,
  ComplianceStatus,
  ComplianceStatusUpsertObject,
} from '@dydxprotocol-indexer/postgres';

import config from '../config';

// eslint-disable-next-line max-len
const CLOSE_ONLY_TO_BLOCKED_DAYS_IN_MS: number = config.CLOSE_ONLY_TO_BLOCKED_DAYS * ONE_DAY_IN_MILLISECONDS;

export default async function runTask(): Promise<void> {
  const queryStart: number = Date.now();

  // Query for addresses with status CLOSE_ONLY and updatedAt less than NOW() - INTERVAL days
  const staleCloseOnlyAddresses: ComplianceStatusFromDatabase[] = await
  ComplianceStatusTable.findAll(
    {
      status: [ComplianceStatus.CLOSE_ONLY],
      updatedBeforeOrAt: new Date(
        queryStart - CLOSE_ONLY_TO_BLOCKED_DAYS_IN_MS,
      ).toISOString(),
    },
    [],
    {
      readReplica: true,
    },
  );
  stats.timing(`${config.SERVICE_NAME}.query_stale_close_only.timing`, Date.now() - queryStart);

  const updateStart: number = Date.now();
  const addressesToUpdate: string[] = staleCloseOnlyAddresses.map(
    (record: ComplianceStatusFromDatabase) => record.address,
  );

  // Update addresses status to BLOCKED
  const complianceStatusUpsertObjects: ComplianceStatusUpsertObject[] = addressesToUpdate
    .reduce(
      (acc: ComplianceStatusUpsertObject[], address: string) => {
        acc.push({
          address,
          status: ComplianceStatus.BLOCKED,
          updatedAt: new Date().toISOString(),
        });
        return acc;
      }, []);
  await ComplianceStatusTable.bulkUpsert(complianceStatusUpsertObjects);

  stats.timing(
    `${config.SERVICE_NAME}.update_stale_close_only.timing`,
    Date.now() - updateStart,
  );
  stats.gauge(`${config.SERVICE_NAME}.num_stale_close_only_updated.count`, addressesToUpdate.length);
}
