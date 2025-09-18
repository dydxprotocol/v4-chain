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
    transaction_hash: '0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890',
    created_at: '2023-01-02T00:00:00.000Z',
  };

  const defaultBridgeInfo3 = {
    from_address: '0x1234567890abcdef1234567890abcdef12345678',
    chain_id: 'avalanche',
    amount: '3000000',
    created_at: '2023-01-03T00:00:00.000Z',
  };

  const defaultBridgeInfo4 = {
    from_address: '0x1234567890abcdef1234567890abcdef12345678',
    chain_id: 'ethereum',
    amount: '4000000',
    transaction_hash: '0x1111111111111111111111111111111111111111111111111111111111111111',
    created_at: '2023-01-04T00:00:00.000Z',
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

  describe('upsert', () => {
    it('Successfully upserts a bridge information record', async () => {
      const upsertedRecord = await BridgeInformationTable.upsert(defaultBridgeInfo1);
      expect(upsertedRecord).toEqual(expect.objectContaining(defaultBridgeInfo1));
      expect(upsertedRecord.id).toBeDefined();

      // Update with upsert
      const updatedRecord = {
        id: upsertedRecord.id,
        ...defaultBridgeInfo1,
        amount: '5000000',
      };
      const secondUpsert = await BridgeInformationTable.upsert(updatedRecord);
      expect(secondUpsert.amount).toBe('5000000');
      expect(secondUpsert.id).toBe(upsertedRecord.id);
    });

    it('Generates ID when not provided in upsert', async () => {
      const upsertedRecord = await BridgeInformationTable.upsert(defaultBridgeInfo1);

      expect(upsertedRecord.id).toBeDefined();
      expect(typeof upsertedRecord.id).toBe('string');
    });
  });

  describe('findById', () => {
    it('Successfully finds a bridge information record by ID', async () => {
      const createdRecord = await BridgeInformationTable.create(defaultBridgeInfo1);

      const foundRecord = await BridgeInformationTable.findById(createdRecord.id);

      expect(foundRecord).toEqual(expect.objectContaining(defaultBridgeInfo1));
      expect(foundRecord?.id).toBe(createdRecord.id);
    });

    it('Returns undefined when record not found by ID', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const foundRecord = await BridgeInformationTable.findById('nonexistent-id');

      expect(foundRecord).toBeUndefined();
    });
  });

  describe('findByFromAddress', () => {
    it('Successfully finds bridge information records by from_address in descending order', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1); // 2023-01-01
      await BridgeInformationTable.create(defaultBridgeInfo3); // 2023-01-03
      await BridgeInformationTable.create({
        ...defaultBridgeInfo4,
        transaction_hash: createUniqueTransactionHash(),
      }); // 2023-01-04
      await BridgeInformationTable.create({
        ...defaultBridgeInfo2,
        transaction_hash: createUniqueTransactionHash(),
      }); // Different address

      const records = await BridgeInformationTable.findByFromAddress(
        defaultBridgeInfo1.from_address,
      );

      expect(records).toHaveLength(3);
      // Should be in descending order by created_at
      expect(records[0].created_at).toBe(defaultBridgeInfo4.created_at);
      expect(records[1].created_at).toBe(defaultBridgeInfo3.created_at);
      expect(records[2].created_at).toBe(defaultBridgeInfo1.created_at);
    });

    it('Returns empty array when no records found by from_address', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const records = await BridgeInformationTable.findByFromAddress(
        '0xnonexistent1234567890abcdef1234567890abcdef',
      );

      expect(records).toHaveLength(0);
    });
  });

  describe('findByFromAddressWithTransactionHashFilter', () => {
    beforeEach(async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1); // No tx hash
      await BridgeInformationTable.create(defaultBridgeInfo3); // No tx hash
      await BridgeInformationTable.create({
        ...defaultBridgeInfo4,
        transaction_hash: createUniqueTransactionHash(),
      }); // Has tx hash
      await BridgeInformationTable.create({
        ...defaultBridgeInfo2,
        transaction_hash: createUniqueTransactionHash(),
      }); // Different address
    });

    it('Successfully finds records with transaction hash', async () => {
      const records = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        true, // hasTransactionHash = true
      );

      expect(records).toHaveLength(1);
      expect(records[0].transaction_hash).not.toBeNull();
      expect(records[0].transaction_hash).toBeDefined();
    });

    it('Successfully finds records without transaction hash', async () => {
      const records = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false, // hasTransactionHash = false
      );

      expect(records).toHaveLength(2);
      records.forEach((record) => {
        expect(record.transaction_hash).toBeNull();
      });
      // Should be in descending order
      expect(records[0].created_at).toBe(defaultBridgeInfo3.created_at);
      expect(records[1].created_at).toBe(defaultBridgeInfo1.created_at);
    });

    it('Supports pagination with limit and offset', async () => {
      const records = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false,
        { limit: 1, offset: 0 },
      );

      expect(records).toHaveLength(1);
      expect(records[0].created_at).toBe(defaultBridgeInfo3.created_at);

      const nextRecords = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false,
        { limit: 1, offset: 1 },
      );

      expect(nextRecords).toHaveLength(1);
      expect(nextRecords[0].created_at).toBe(defaultBridgeInfo1.created_at);
    });

    it('Supports custom ordering by amount', async () => {
      const records = await BridgeInformationTable.findByFromAddressWithTransactionHashFilter(
        defaultBridgeInfo1.from_address,
        false,
        { orderBy: 'amount', orderDirection: 'ASC' },
      );

      expect(records).toHaveLength(2);
      expect(records[0].amount).toBe(defaultBridgeInfo1.amount); // 1000000
      expect(records[1].amount).toBe(defaultBridgeInfo3.amount); // 3000000
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
        ...defaultBridgeInfo4,
        transaction_hash: createUniqueTransactionHash(),
      });
    });

    it('Searches by from_address filter', async () => {
      const records = await BridgeInformationTable.searchBridgeInformation({
        from_address: defaultBridgeInfo1.from_address,
      });

      expect(records).toHaveLength(3);
      records.forEach((record) => {
        expect(record.from_address).toBe(defaultBridgeInfo1.from_address);
      });
    });

    it('Searches by chain_id filter', async () => {
      const records = await BridgeInformationTable.searchBridgeInformation({
        chain_id: 'ethereum',
      });

      expect(records).toHaveLength(2);
      records.forEach((record) => {
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

      const records = await BridgeInformationTable.searchBridgeInformation({
        transaction_hash: specificTxHash,
      });

      expect(records).toHaveLength(1);
      expect(records[0].transaction_hash).toBe(specificTxHash);
    });

    it('Searches by has_transaction_hash filter', async () => {
      const recordsWithTx = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: true,
      });

      const recordsWithoutTx = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: false,
      });

      expect(recordsWithTx).toHaveLength(2);
      expect(recordsWithoutTx).toHaveLength(2);

      recordsWithTx.forEach((record) => {
        expect(record.transaction_hash).not.toBeNull();
      });

      recordsWithoutTx.forEach((record) => {
        expect(record.transaction_hash).toBeNull();
      });
    });

    it('Combines multiple filters', async () => {
      const records = await BridgeInformationTable.searchBridgeInformation({
        from_address: defaultBridgeInfo1.from_address,
        chain_id: 'ethereum',
        has_transaction_hash: true,
      });

      expect(records).toHaveLength(1);
      expect(records[0].from_address).toBe(defaultBridgeInfo1.from_address);
      expect(records[0].chain_id).toBe('ethereum');
      expect(records[0].transaction_hash).not.toBeNull();
    });

    it('Supports custom ordering and pagination', async () => {
      const records = await BridgeInformationTable.searchBridgeInformation(
        {},
        {
          orderBy: 'amount',
          orderDirection: 'DESC',
          limit: 2,
          offset: 1,
        },
      );

      expect(records).toHaveLength(2);
      // Should skip the highest amount and return the next 2
      expect(records[0].amount).toBe('3000000');
      expect(records[1].amount).toBe('2000000');
    });

    it('Returns empty array when no records match filters', async () => {
      const records = await BridgeInformationTable.searchBridgeInformation({
        chain_id: 'nonexistent-chain',
      });

      expect(records).toHaveLength(0);
    });
  });

  describe('updateTransactionHash', () => {
    it('Successfully updates transaction hash', async () => {
      const createdRecord = await BridgeInformationTable.create(defaultBridgeInfo1);
      expect(createdRecord.transaction_hash).toBeNull();

      const newTxHash = '0x9999999999999999999999999999999999999999999999999999999999999999';
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
        '0x1234567890abcdef',
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
    it('Handles empty string from_address search', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const records = await BridgeInformationTable.findByFromAddress('');

      expect(records).toHaveLength(0);
    });

    it('Handles case sensitivity for addresses', async () => {
      await BridgeInformationTable.create(defaultBridgeInfo1);

      const upperCaseAddress = defaultBridgeInfo1.from_address.toUpperCase();
      const records = await BridgeInformationTable.findByFromAddress(upperCaseAddress);

      // Should not find the record as addresses are case-sensitive
      expect(records).toHaveLength(0);
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

      const records = await BridgeInformationTable.searchBridgeInformation({
        has_transaction_hash: false,
      });

      expect(records).toHaveLength(2);
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
  });
});
