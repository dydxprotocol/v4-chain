import {
  delay,
  logger,
  stats,
} from '@dydxprotocol-indexer/base';
import { ComplianceClientResponse } from '@dydxprotocol-indexer/compliance';
import {
  ComplianceDataColumns,
  ComplianceDataCreateObject,
  ComplianceDataFromDatabase,
  ComplianceTable,
  IsoString,
  SubaccountColumns,
  SubaccountFromDatabase,
  SubaccountTable,
} from '@dydxprotocol-indexer/postgres';
import _ from 'lodash';
import { DateTime } from 'luxon';

import config from '../config';
import { ClientAndProvider } from '../helpers/compliance-clients';

const taskName: string = 'update_compliance_data';

/**
 * This task updates the compliance data in the indexer.
 * On each run of this task:
 * - for all addresses with a subaccount that was updated within the last day (by default), update
 *   the compliance data for the address if:
 *     - the address has no compliance data OR
 *     - the address has compliance data more than 1 day (default) old and isn't blocked
 * - for all addresses with compliance data that was updated more than 1 month (by default) ago,
 *   update the compliance data for the address if:
 *     - the address isn't blocked
 *
 * The task takes in a limit for the number of addresses queried on the compliance provider per
 * iteration, and will stop querying for addresses once the limit is hit.
 * @param complianceProvider
 */
export default async function runTask(
  complianceProvider: ClientAndProvider,
): Promise<void> {
  const startTime: DateTime = DateTime.utc();
  let remainingQueries: number = config.MAX_COMPLIANCE_DATA_QUERY_PER_LOOP;
  const activeAddressThreshold: IsoString = startTime.minus(
    { seconds: config.ACTIVE_ADDRESS_THRESHOLD_SECONDS },
  ).toISO();
  const ageThreshold: IsoString = startTime.minus(
    { seconds: config.MAX_COMPLIANCE_DATA_AGE_SECONDS },
  ).toISO();
  const activeAgeThreshold: IsoString = startTime.minus(
    { seconds: config.MAX_ACTIVE_COMPLIANCE_DATA_AGE_SECONDS },
  ).toISO();
  let addressesToQuery: string[] = [];

  try {
    const startActiveAddresses: number = Date.now();
    // Get addresses that had activity recently
    const activeSubaccounts: SubaccountFromDatabase[] = await SubaccountTable.findAll(
      { updatedOnOrAfter: activeAddressThreshold },
      [],
      { readReplica: true },
    );
    const activeAddresses: string[] = _.chain(activeSubaccounts)
      .map(SubaccountColumns.address)
      .uniq()
      .value();

    // Get corresponding compliance data for all active addresses
    const activeAddressCompliance: ComplianceDataFromDatabase[] = await ComplianceTable.findAll(
      // To handle new addresses, can't filter by the compliance data age
      { address: activeAddresses, provider: complianceProvider.provider },
      [],
      { readReplica: true },
    );
    const addressesWithCompliance: string[] = _.chain(activeAddressCompliance)
      .map(ComplianceDataColumns.address)
      .uniq()
      .value();

    // Add any address that does not have compliance data to the list of addresses to query
    // Note: The query for compliance data can't filter out blocked or new compliance data. If it
    // did, the below logic wouldn't be able to correctly get the list of active addresses that are
    // new (have no compliane data stored).
    const addressesWithoutCompliance: string[] = _.without(
      activeAddresses,
      ...addressesWithCompliance,
    );
    if (addressesWithoutCompliance.length > remainingQueries) {
      remainingQueries = 0;
      addressesToQuery.push(...addressesWithoutCompliance.slice(0, remainingQueries));
    } else {
      remainingQueries -= addressesWithoutCompliance.length;
      addressesToQuery.push(...addressesWithoutCompliance);
    }

    // Add any address that has compliance data that's over the age threshold for active addresses
    // and is not blocked
    let activeAddressesToQuery: number = 0;
    for (const addressCompliance of activeAddressCompliance) {
      if (remainingQueries <= 0) {
        break;
      }

      if (addressCompliance.blocked) {
        continue;
      }

      if (DateTime.fromISO(addressCompliance.updatedAt) > DateTime.fromISO(activeAgeThreshold)) {
        continue;
      }

      addressesToQuery.push(addressCompliance.address);
      remainingQueries -= 1;
      activeAddressesToQuery += 1;
    }

    stats.timing(
      `${config.SERVICE_NAME}.${taskName}.get_active_addresses`,
      Date.now() - startActiveAddresses,
      undefined,
      { provider: complianceProvider.provider },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${taskName}.num_active_addresses`,
      activeAddressesToQuery,
      undefined,
      { provider: complianceProvider.provider },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${taskName}.num_new_addresses`,
      addressesWithoutCompliance.length,
      undefined,
      { provider: complianceProvider.provider },
    );

    const startOldAddresses: number = Date.now();
    // Get old compliance data
    const oldAddressCompliance: ComplianceDataFromDatabase[] = await ComplianceTable.findAll(
      {
        blocked: false,
        provider: complianceProvider.provider,
        updatedBeforeOrAt: ageThreshold,
        limit: remainingQueries,
      },
      [],
      { readReplica: true },
    );
    addressesToQuery.push(...(
      _.chain(oldAddressCompliance).map(ComplianceDataColumns.address).uniq().value()
    ));
    // Ensure all addresses to query are unique
    addressesToQuery = _.sortedUniq(addressesToQuery);

    stats.timing(
      `${config.SERVICE_NAME}.${taskName}.get_old_addresses`,
      Date.now() - startOldAddresses,
      undefined,
      { provider: complianceProvider.provider },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${taskName}.num_old_addresses`,
      oldAddressCompliance.length,
      undefined,
      { provider: complianceProvider.provider },
    );

    // Get compliance data for addresses
    const startQueryProvider: number = Date.now();
    const complianceResponses: ComplianceClientResponse[] = await getComplianceData(
      addressesToQuery,
      complianceProvider,
    );
    const calculatedAt: string = DateTime.utc().toISO();
    const complianceCreateObjects: ComplianceDataCreateObject[] = complianceResponses.map(
      (complianceResponse: ComplianceClientResponse): ComplianceDataCreateObject => {
        return {
          ...complianceResponse,
          provider: complianceProvider.provider,
          updatedAt: calculatedAt,
        };
      },
    );

    stats.timing(
      `${config.SERVICE_NAME}.${taskName}.query_compliance_data`,
      Date.now() - startQueryProvider,
      undefined,
      { provider: complianceProvider.provider },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${taskName}.num_addresses_to_screen`,
      addressesToQuery.length,
      undefined,
      { provider: complianceProvider.provider },
    );

    // Upsert into database
    const startUpsert: number = Date.now();
    await ComplianceTable.bulkUpsert(
      complianceCreateObjects,
    );

    stats.timing(
      `${config.SERVICE_NAME}.${taskName}.upsert_compliance_data`,
      Date.now() - startUpsert,
      undefined,
      { provider: complianceProvider.provider },
    );
    stats.gauge(
      `${config.SERVICE_NAME}.${taskName}.num_upserted`,
      complianceCreateObjects.length,
      undefined,
      { provider: complianceProvider.provider },
    );
  } catch (error) {
    logger.error({
      at: taskName,
      message: 'Error occurred in task for updating compliance data',
      error,
    });
  }
}

async function getComplianceData(
  addresses: string[],
  complianceProvider: ClientAndProvider,
): Promise<ComplianceClientResponse[]> {
  const complianceResponses: ComplianceClientResponse[] = [];
  for (const complianceBatch of _.chunk(addresses, config.COMPLIANCE_PROVIDER_QUERY_BATCH_SIZE)) {
    const startBatch: number = Date.now();
    const responses: PromiseSettledResult<ComplianceClientResponse>[] = await Promise.allSettled(
      complianceBatch.map((address: string): Promise<ComplianceClientResponse> => {
        return complianceProvider.client.getComplianceResponse(address);
      }),
    );
    const successResponses: PromiseFulfilledResult<ComplianceClientResponse>[] = responses.filter(
      (result: PromiseSettledResult<ComplianceClientResponse>):
        result is PromiseFulfilledResult<ComplianceClientResponse> => {
        return result.status === 'fulfilled';
      },
    );
    const failedResponses: PromiseRejectedResult[] = responses.filter(
      (result: PromiseSettledResult<ComplianceClientResponse>):
        result is PromiseRejectedResult => {
        return result.status === 'rejected';
      },
    );
    complianceResponses.push(...successResponses.map(
      (result: PromiseFulfilledResult<ComplianceClientResponse>): ComplianceClientResponse => {
        return result.value;
      },
    ));

    if (failedResponses.length > 0) {
      const addressesWithoutResponses: string[] = _.without(
        addresses,
        ..._.map(complianceResponses, 'address'),
      );
      stats.increment(
        `${config.SERVICE_NAME}.${taskName}.get_compliance_data_fail`,
        1,
        undefined,
        { provider: complianceProvider.provider },
      );
      logger.error({
        at: 'updated-compliance-data#getComplianceData',
        message: 'Failed to retrieve compliance data for the addresses',
        addresses: addressesWithoutResponses,
        errors: failedResponses,
      });
    }
    stats.timing(
      `${config.SERVICE_NAME}.${taskName}.get_batch_compliance_data`,
      Date.now() - startBatch,
      undefined,
      { provider: complianceProvider.provider },
    );
    await delay(config.COMPLIANCE_PROVIDER_QUERY_DELAY_MS);
  }
  return complianceResponses;
}
