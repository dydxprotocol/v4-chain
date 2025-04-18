import { DateTime } from 'luxon';

import { MAX_PARENT_SUBACCOUNTS } from '../../src/constants';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  getParentSubaccountNum,
  getSubaccountQueryForParent,
} from '../../src/lib/parent-subaccount-helpers';
import * as SubaccountTable from '../../src/stores/subaccount-table';

const defaultDateTime: DateTime = DateTime.fromISO('2025-01-01T00:00:00.000Z');

describe('getParentSubaccountNum', () => {
  it('Gets the parent subaccount number from a child subaccount number', () => {
    expect(getParentSubaccountNum(0)).toEqual(0);
    expect(getParentSubaccountNum(128)).toEqual(0);
    expect(getParentSubaccountNum(128 * 999 - 1)).toEqual(127);
  });

  it('Throws an error if the child subaccount number is greater than the max child subaccount number', () => {
    expect(() => getParentSubaccountNum(128001)).toThrowError('Child subaccount number must be less than or equal to 128000');
  });
});

describe('getSubaccountQueryForParent', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('returns the correct subaccount IDs for a parent subaccount', async () => {
    // Setup - create parent and child subaccounts
    const address = 'test_parent_address';
    const parentSubaccountNumber = 0;

    // Create several subaccounts:
    // 1. Parent subaccount (0)
    // 2. Child subaccounts (128, 256, 384, 512)
    // 3. Subaccount with same address but not a child (1) - should not be returned
    // 4. Subaccount with different address but same number (0) - should not be returned
    const childSubaccountNumbers = [0, 128, 256, 384, 512];
    const nonChildSubaccountNumbers = [1, 2, 3];
    const differentAddress = 'different_address';

    // Create parent and child subaccounts
    await Promise.all(childSubaccountNumbers.map((subaccountNum) => SubaccountTable.create({
      address,
      subaccountNumber: subaccountNum,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '1',
    })));

    // Create non-child subaccounts
    await Promise.all(nonChildSubaccountNumbers.map((subaccountNum) => SubaccountTable.create({
      address,
      subaccountNumber: subaccountNum,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '1',
    })));

    // Create subaccounts with different address
    await Promise.all([0, 128].map((subaccountNum) => SubaccountTable.create({
      address: differentAddress,
      subaccountNumber: subaccountNum,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '1',
    })));

    // Execute query
    const subaccountQuery = getSubaccountQueryForParent({
      address,
      subaccountNumber: parentSubaccountNumber,
    });

    // Get results by executing the query
    const results = await subaccountQuery;

    // Assert expected results - should only include parent and child subaccounts
    expect(results.length).toEqual(5); // Parent + 4 children

    // Get all the subaccount IDs that should be returned
    const expectedSubaccountIds = childSubaccountNumbers.map(
      (num) => SubaccountTable.uuid(address, num),
    );

    // Convert results to array of subaccountIds
    const resultSubaccountIds = results.map((row: { subaccountId: string }) => row.subaccountId);

    // Verify all expected IDs are in the results
    expect(resultSubaccountIds.sort()).toEqual(expectedSubaccountIds.sort());

    // Verify non-child subaccounts are not included
    const nonChildSubaccountIds = nonChildSubaccountNumbers.map(
      (num) => SubaccountTable.uuid(address, num),
    );
    for (const id of nonChildSubaccountIds) {
      expect(resultSubaccountIds).not.toContain(id);
    }

    // Verify different address subaccounts are not included
    const differentAddressSubaccountIds = [0, 128].map(
      (num) => SubaccountTable.uuid(differentAddress, num),
    );
    for (const id of differentAddressSubaccountIds) {
      expect(resultSubaccountIds).not.toContain(id);
    }
  });

  it('returns the correct subaccount IDs for a non-zero parent subaccount', async () => {
    // Setup with a non-zero parent subaccount number
    const address = 'test_parent_address';
    const parentSubaccountNumber = 5; // Non-zero parent

    // Create subaccounts with parent and child numbers
    const childSubaccountNumbers = [
      parentSubaccountNumber,
      parentSubaccountNumber + MAX_PARENT_SUBACCOUNTS,
      parentSubaccountNumber + (2 * MAX_PARENT_SUBACCOUNTS),
    ];

    await Promise.all(childSubaccountNumbers.map((subaccountNum) => SubaccountTable.create({
      address,
      subaccountNumber: subaccountNum,
      updatedAt: defaultDateTime.toISO(),
      updatedAtHeight: '1',
    })));

    // Execute query
    const subaccountQuery = getSubaccountQueryForParent({
      address,
      subaccountNumber: parentSubaccountNumber,
    });

    // Get results
    const results = await subaccountQuery;

    // Assert
    expect(results.length).toEqual(3);

    const expectedSubaccountIds = childSubaccountNumbers.map(
      (num) => SubaccountTable.uuid(address, num),
    );
    const resultSubaccountIds = results.map((row: { subaccountId: string }) => row.subaccountId);

    expect(resultSubaccountIds.sort()).toEqual(expectedSubaccountIds.sort());
  });

  it('returns empty array when no matching subaccounts exist', async () => {
    // Execute query for non-existent address
    const subaccountQuery = getSubaccountQueryForParent({
      address: 'nonexistent_address',
      subaccountNumber: 0,
    });

    // Get results
    const results = await subaccountQuery;

    // Assert
    expect(results).toEqual([]);
  });
});
