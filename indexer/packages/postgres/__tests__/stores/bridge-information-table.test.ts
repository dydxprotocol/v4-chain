import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import * as BridgeInformationTable from '../../src/stores/bridge-information-table';

describe('BridgeInformation store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const defaultBridgeInfo1 = {
    from_address: '0x1234567890abcdef1234567890abcdef12345678',
    chain_id: 'ethereum',
    amount: '1000000',
    created_at: '2023-01-01T00:00:00.000Z',
  };

  const defaultBridgeInfo2 = {
    from_address: '0x9876543210fedcba9876543210fedcba98765432',
    chain_id: 'polygon',
    amount: '2000000',
    created_at: '2023-01-02T00:00:00.000Z',
  };

  const defaultBridgeInfo3 = {
    from_address: '0x1234567890abcdef1234567890abcdef12345678',
    chain_id: 'avalanche',
    amount: '3000000',
    created_at: '2023-01-03T00:00:00.000Z',
  };

  // Helper function to create unique transaction hashes for tests
  const createUniqueTransactionHash = () => {
    return `0x${Math.random().toString(16).substr(2, 64).padEnd(64, '0')}`;
  };

  describe('create', () => {
    it('Successfully creates a bridge information record without ID', async () => {
      const createdRecord = await BridgeInformationTable.create(defaultBridgeInfo1);

      expect(createdRecord).toEqual(expect.objectContaining(defaultBridgeInfo1));
      expect(createdRecord.id).toBeDefined();
      expect(typeof createdRecord.id).toBe('string');
      expect(createdRecord.id.length).toBeGreaterThan(0);
    });

    it('Successfully creates a bridge information record with explicit ID', async () => {
      const bridgeInfoWithId = {
        ...defaultBridgeInfo1,
        id: 'custom-id-123',
      };

      const createdRecord = await BridgeInformationTable.create(bridgeInfoWithId);

      expect(createdRecord).toEqual(expect.objectContaining(bridgeInfoWithId));
      expect(createdRecord.id).toBe('custom-id-123');
    });

    it('Auto-generates different IDs for multiple records', async () => {
      const record1 = await BridgeInformationTable.create(defaultBridgeInfo1);
      const record2 = await BridgeInformationTable.create(defaultBridgeInfo2);

      expect(record1.id).toBeDefined();
      expect(record2.id).toBeDefined();
      expect(record1.id).not.toBe(record2.id);
    });

    it('Handles null transaction_hash correctly', async () => {
      const createdRecord = await BridgeInformationTable.create(defaultBridgeInfo1);

      expect(createdRecord.transaction_hash).toBeNull();
    });

    it('Handles provided transaction_hash correctly', async () => {
      const uniqueTxHash = createUniqueTransactionHash();
      const bridgeInfoWithTx = {
        ...defaultBridgeInfo2,
        transaction_hash: uniqueTxHash,
      };

      const createdRecord = await BridgeInformationTable.create(bridgeInfoWithTx);

      expect(createdRecord.transaction_hash).toBe(uniqueTxHash);
    });
  });

  describe('findByFromAddressWithTransactionHashFilter', () => {
    beforeEach(async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1); // No tx hash
      await BridgeInformationTable.create(defaultBridgeInfo3); // No tx hash
      await BridgeInformationTable.create({
        ...defaultBridgeInfo1,
        from_address: '0x1234567890abcdef1234567890abcdef12345678',
        amount: '4000000',
        created_at: '2023-01-04T00:00:00.000Z',
        transaction_hash: createUniqueTransactionHash(),
      }); // Has tx hash
      await BridgeInformationTable.create({
        ...defaultBridgeInfo2,
        transaction_hash: createUniqueTransactionHash(),
      }); // Different address
    });

    it('Successfully finds records with transaction hash', async () => {
      const result = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        true, // hasTransactionHash = true
      );

      expect(result.results).toHaveLength(1);
      expect(result.results[0].transaction_hash).not.toBeNull();
      expect(result.results[0].transaction_hash).toBeDefined();
    });

    it('Successfully finds records without transaction hash', async () => {
      const result = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false, // hasTransactionHash = false
      );

      expect(result.results).toHaveLength(2);
      result.results.forEach((record) => {
        expect(record.transaction_hash).toBeNull();
      });
      // Should be in descending order
      expect(result.results[0].created_at).toBe(defaultBridgeInfo3.created_at);
      expect(result.results[1].created_at).toBe(defaultBridgeInfo1.created_at);
    });

    it('Supports pagination with limit and page', async () => {
      const result = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false,
        { limit: 1, page: 1 },
      );

      expect(result.results).toHaveLength(1);
      expect(result.results[0].created_at).toBe(defaultBridgeInfo3.created_at);
      expect(result.limit).toBe(1);
      expect(result.offset).toBe(0); // page 1 with limit 1 = offset 0

      const nextResult = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false,
        { limit: 1, page: 2 },
      );

      expect(nextResult.results).toHaveLength(1);
      expect(nextResult.results[0].created_at).toBe(defaultBridgeInfo1.created_at);
      expect(nextResult.offset).toBe(1); // page 2 with limit 1 = offset 1
    });
  });

  describe('searchBridgeInformation', () => {
    beforeEach(async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);
      await BridgeInformationTable.create({
        ...defaultBridgeInfo2,
        transaction_hash: createUniqueTransactionHash(),
      });
      await BridgeInformationTable.create(defaultBridgeInfo3);
      await BridgeInformationTable.create({
        ...defaultBridgeInfo1,
        from_address: '0x1234567890abcdef1234567890abcdef12345678',
        amount: '4000000',
        created_at: '2023-01-04T00:00:00.000Z',
        transaction_hash: createUniqueTransactionHash(),
      });
    });

    it('Searches by from_addresses filter', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation({
        from_addresses: [defaultBridgeInfo1.from_address],
      });

      expect(result.results).toHaveLength(3);
      result.results.forEach((record) => {
        expect(record.from_address).toBe(defaultBridgeInfo1.from_address);
      });
    });

    it('Searches by multiple from_addresses', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation({
        from_addresses: [defaultBridgeInfo1.from_address, defaultBridgeInfo2.from_address],
      });

      expect(result.results).toHaveLength(4);
    });

    it('Searches by chain_id filter', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation({
        chain_id: 'ethereum',
      });

      expect(result.results).toHaveLength(2);
      result.results.forEach((record) => {
        expect(record.chain_id).toBe('ethereum');
      });
    });

    it('Searches by transaction_hash filter', async () => {
      // First create a record with a specific transaction hash
      const specificTxHash = createUniqueTransactionHash();
      const testRecord = {
        ...defaultBridgeInfo1,
        from_address: '0xspecific1234567890abcdef1234567890abcdef12',
        transaction_hash: specificTxHash,
      };
      await BridgeInformationTable.create(testRecord);

      const result = await BridgeInformationTable.searchBridgeInformation({
        transaction_hash: specificTxHash,
      });

      expect(result.results).toHaveLength(1);
      expect(result.results[0].transaction_hash).toBe(specificTxHash);
    });

    it('Searches by has_transaction_hash filter', async () => {
      const resultWithTx = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: true,
      });

      const resultWithoutTx = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: false,
      });

      expect(resultWithTx.results).toHaveLength(2);
      expect(resultWithoutTx.results).toHaveLength(2);

      resultWithTx.results.forEach((record) => {
        expect(record.transaction_hash).not.toBeNull();
      });

      resultWithoutTx.results.forEach((record) => {
        expect(record.transaction_hash).toBeNull();
      });
    });

    it('Combines multiple filters', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation({
        from_addresses: [defaultBridgeInfo1.from_address],
        chain_id: 'ethereum',
        has_transaction_hash: true,
      });

      expect(result.results).toHaveLength(1);
      expect(result.results[0].from_address).toBe(defaultBridgeInfo1.from_address);
      expect(result.results[0].chain_id).toBe('ethereum');
      expect(result.results[0].transaction_hash).not.toBeNull();
    });

    it('Supports pagination with limit only', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation(
        {},
        { limit: 2 },
      );

      expect(result.results).toHaveLength(2);
    });

    it('Returns empty results when no records match filters', async () => {
      const result = await BridgeInformationTable.searchBridgeInformation({
        chain_id: 'nonexistent-chain',
      });

      expect(result.results).toHaveLength(0);
    });
  });

  describe('updateTransactionHash', () => {
    it('Successfully updates transaction hash', async () => {
      const createdRecord = await BridgeInformationTable.create(defaultBridgeInfo1);
      expect(createdRecord.transaction_hash).toBeNull();

      const newTxHash = createUniqueTransactionHash();
      const updatedRecord = await BridgeInformationTable.updateTransactionHash(
        createdRecord.id,
        newTxHash,
      );

      expect(updatedRecord).toBeDefined();
      expect(updatedRecord?.transaction_hash).toBe(newTxHash);
      expect(updatedRecord?.id).toBe(createdRecord.id);
    });

    it('Returns undefined when record not found', async () => {
      const updatedRecord = await BridgeInformationTable.updateTransactionHash(
        'nonexistent-id',
        createUniqueTransactionHash(),
      );

      expect(updatedRecord).toBeUndefined();
    });

    it('Updates existing transaction hash', async () => {
      const initialTxHash = createUniqueTransactionHash();
      const bridgeInfoWithTx = {
        ...defaultBridgeInfo2,
        transaction_hash: initialTxHash,
      };
      const createdRecord = await BridgeInformationTable.create(bridgeInfoWithTx);
      expect(createdRecord.transaction_hash).toBe(initialTxHash);

      const newTxHash = createUniqueTransactionHash();
      const updatedRecord = await BridgeInformationTable.updateTransactionHash(
        createdRecord.id,
        newTxHash,
      );

      expect(updatedRecord?.transaction_hash).toBe(newTxHash);
    });
  });

  describe('Edge cases and validation', () => {
    it('Handles empty from_addresses array', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const result = await BridgeInformationTable.searchBridgeInformation({
        from_addresses: [],
      });

      expect(result.results).toHaveLength(0);
    });

    it('Handles case sensitivity for addresses', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const upperCaseAddress = defaultBridgeInfo1.from_address.toUpperCase();
      const result = await BridgeInformationTable.searchBridgeInformation({
        from_addresses: [upperCaseAddress],
      });

      // Should not find the record as addresses are case-sensitive
      expect(result.results).toHaveLength(0);
    });

    it('Handles duplicate transaction hashes (should fail due to unique constraint)', async () => {
      const txHash = createUniqueTransactionHash();
      const firstRecord = {
        ...defaultBridgeInfo2,
        transaction_hash: txHash,
      };
      await BridgeInformationTable.create(firstRecord);

      const duplicateTxHash = {
        ...defaultBridgeInfo1,
        transaction_hash: txHash,
      };

      await expect(
        BridgeInformationTable.create(duplicateTxHash),
      ).rejects.toThrow();
    });

    it('Allows multiple records with null transaction_hash', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);
      await BridgeInformationTable.create(defaultBridgeInfo3);

      const result = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: false,
      });

      expect(result.results).toHaveLength(2);
    });

    it('Handles very long addresses and chain IDs', async () => {
      const longDataRecord = {
        from_address: `0x${'1234567890abcdef'.repeat(10)}`, // Very long address
        chain_id: `very-long-chain-id-${'x'.repeat(100)}`,
        amount: '999999999999999999999',
        created_at: '2023-01-01T00:00:00.000Z',
      };

      const createdRecord = await BridgeInformationTable.create(longDataRecord);
      expect(createdRecord).toEqual(expect.objectContaining(longDataRecord));
    });

    it('Handles pagination correctly', async () => {
      // Create multiple records
      await BridgeInformationTable.create(defaultBridgeInfo1);
      await BridgeInformationTable.create(defaultBridgeInfo2);
      await BridgeInformationTable.create(defaultBridgeInfo3);

      // Test simple limit without pagination
      const result = await BridgeInformationTable.searchBridgeInformation(
        {},
        { limit: 2 },
      );

      expect(result.results).toHaveLength(2);
      // When no page is provided, pagination info is not included
      expect(result.limit).toBeUndefined();
      expect(result.offset).toBeUndefined();
      expect(result.total).toBeUndefined();
    });

    it('Handles limit without pagination', async () => {
      // Create multiple records
      await BridgeInformationTable.create(defaultBridgeInfo1);
      await BridgeInformationTable.create(defaultBridgeInfo2);
      await BridgeInformationTable.create(defaultBridgeInfo3);

      // Test just limit without page
      const result = await BridgeInformationTable.searchBridgeInformation(
        {},
        { limit: 2 },
      );

      expect(result.results).toHaveLength(2);
      expect(result.limit).toBeUndefined(); // No pagination info when page not provided
      expect(result.offset).toBeUndefined();
      expect(result.total).toBeUndefined();
    });

    it('Returns all records when no limit is specified', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);
      await BridgeInformationTable.create(defaultBridgeInfo2);
      await BridgeInformationTable.create(defaultBridgeInfo3);

      const result = await BridgeInformationTable.searchBridgeInformation();

      expect(result.results).toHaveLength(3);
    });
  });
});
