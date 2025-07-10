import { TurnkeyUserFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import * as TurnkeyUserTable from '../../src/stores/turnkey-users-table';

describe('TurnkeyUser store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const defaultTurnkeyUser1 = {
    suborgId: 'suborg-1',
    username: 'user1',
    email: 'user1@example.com',
    svmAddress: 'SVM123456789',
    evmAddress: '0x1234567890abcdef1234567890abcdef12345678',
    salt: 'salt1',
    dydxAddress: 'dydx1abc123',
    createdAt: '2023-01-01T00:00:00.000Z',
  };

  const defaultTurnkeyUser2 = {
    suborgId: 'suborg-2',
    username: 'user2',
    email: 'user2@example.com',
    svmAddress: 'SVM987654321',
    evmAddress: '0x9876543210fedcba9876543210fedcba98765432',
    salt: 'salt2',
    dydxAddress: 'dydx1xyz789',
    createdAt: '2023-01-02T00:00:00.000Z',
  };

  it('Successfully creates a TurnkeyUser', async () => {
    const createdUser = await TurnkeyUserTable.create(defaultTurnkeyUser1);
    expect(createdUser).toEqual(expect.objectContaining(defaultTurnkeyUser1));
  });

  it('Successfully upserts a TurnkeyUser', async () => {
    await TurnkeyUserTable.upsert(defaultTurnkeyUser1);
    const user = await TurnkeyUserTable.findByEvmAddress(defaultTurnkeyUser1.evmAddress);
    expect(user).toEqual(expect.objectContaining(defaultTurnkeyUser1));

    // Update with upsert
    const updatedUser = {
      ...defaultTurnkeyUser1,
      username: 'updated_user1',
      email: 'updated_user1@example.com',
    };
    await TurnkeyUserTable.upsert(updatedUser);
    const updatedResult = await TurnkeyUserTable.findByEvmAddress(defaultTurnkeyUser1.evmAddress);
    expect(updatedResult).toEqual(expect.objectContaining(updatedUser));
  });

  describe('findByEvmAddress', () => {
    it('Successfully finds a TurnkeyUser by EVM address', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);
      await TurnkeyUserTable.create(defaultTurnkeyUser2);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEvmAddress(
        defaultTurnkeyUser1.evmAddress,
      );

      expect(user).toEqual(expect.objectContaining(defaultTurnkeyUser1));
    });

    it('Returns undefined when TurnkeyUser not found by EVM address', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEvmAddress(
        '0xnonexistent1234567890abcdef1234567890abcdef',
      );

      expect(user).toBeUndefined();
    });
  });

  describe('findBySvmAddress', () => {
    it('Successfully finds a TurnkeyUser by SVM address', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);
      await TurnkeyUserTable.create(defaultTurnkeyUser2);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findBySvmAddress(
        defaultTurnkeyUser2.svmAddress,
      );

      expect(user).toEqual(expect.objectContaining(defaultTurnkeyUser2));
    });

    it('Returns undefined when TurnkeyUser not found by SVM address', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findBySvmAddress(
        'SVMnonexistent123456789',
      );

      expect(user).toBeUndefined();
    });
  });

  describe('Edge cases', () => {
    it('Handles case sensitivity for EVM addresses', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      // Test with uppercase version of the address
      const upperCaseEvmAddress = defaultTurnkeyUser1.evmAddress.toUpperCase();
      const user = await TurnkeyUserTable.findByEvmAddress(upperCaseEvmAddress);

      // This should not find the user as EVM addresses are case-sensitive in the database
      expect(user).toBeUndefined();
    });

    it('Handles empty string searches', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      const userByEmptyEvm = await TurnkeyUserTable.findByEvmAddress('');
      const userByEmptySvm = await TurnkeyUserTable.findBySvmAddress('');

      expect(userByEmptyEvm).toBeUndefined();
      expect(userByEmptySvm).toBeUndefined();
    });
  });
});
