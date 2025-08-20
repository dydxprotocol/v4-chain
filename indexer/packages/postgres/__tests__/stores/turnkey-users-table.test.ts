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
    suborg_id: 'suborg-1',
    email: 'user1@example.com',
    svm_address: 'SVM123456789',
    evm_address: '0x1234567890abcdef1234567890abcdef12345678',
    salt: 'salt1',
    dydx_address: 'dydx1abc123',
    created_at: '2023-01-01T00:00:00.000Z',
  };

  const defaultTurnkeyUser2 = {
    suborg_id: 'suborg-2',
    email: 'user2@example.com',
    svm_address: 'SVM987654321',
    evm_address: '0x9876543210fedcba9876543210fedcba98765432',
    salt: 'salt2',
    dydx_address: 'dydx1xyz789',
    created_at: '2023-01-02T00:00:00.000Z',
  };

  const defaultTurnkeyUser3 = {
    suborg_id: 'suborg-3',
    svm_address: 'SVM987654321',
    evm_address: '0x9876543210fedcba9876543210fedcba98765433',
    salt: 'salt3',
    dydx_address: 'dydx1xyz789',
    created_at: '2023-01-02T00:00:00.000Z',
  };

  it('Successfully creates a TurnkeyUser', async () => {
    const createdUser = await TurnkeyUserTable.create(defaultTurnkeyUser1);
    expect(createdUser).toEqual(expect.objectContaining(defaultTurnkeyUser1));
  });

  it('Successfully upserts a TurnkeyUser', async () => {
    await TurnkeyUserTable.upsert(defaultTurnkeyUser1);
    const user = await TurnkeyUserTable.findByEvmAddress(defaultTurnkeyUser1.evm_address);
    expect(user).toEqual(expect.objectContaining(defaultTurnkeyUser1));

    // Update with upsert
    const updatedUser = {
      ...defaultTurnkeyUser1,
      email: 'updated_user1@example.com',
    };
    await TurnkeyUserTable.upsert(updatedUser);
    const updatedResult = await TurnkeyUserTable.findByEvmAddress(defaultTurnkeyUser1.evm_address);
    expect(updatedResult).toEqual(expect.objectContaining(updatedUser));
  });

  describe('findByEvmAddress', () => {
    it('Successfully finds a TurnkeyUser by EVM address', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);
      await TurnkeyUserTable.create(defaultTurnkeyUser2);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEvmAddress(
        defaultTurnkeyUser1.evm_address,
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
        defaultTurnkeyUser2.svm_address,
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

  describe('findByEmail', () => {
    it('Successfully finds a TurnkeyUser by email', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);
      await TurnkeyUserTable.create(defaultTurnkeyUser2);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEmail(
        defaultTurnkeyUser1.email,
      );

      expect(user).toEqual(expect.objectContaining(defaultTurnkeyUser1));
    });

    it('Returns undefined when TurnkeyUser not found by email', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEmail(
        'nonexistent@example.com',
      );

      expect(user).toBeUndefined();
    });

    it('Returns undefined when email is empty', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser3);

      const user: TurnkeyUserFromDatabase | undefined = await TurnkeyUserTable.findByEmail('');

      expect(user).toBeUndefined();
    });
  });

  describe('Edge cases', () => {
    it('Handles case sensitivity for EVM addresses', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      // Test with uppercase version of the address
      const upperCaseEvmAddress = defaultTurnkeyUser1.evm_address.toUpperCase();
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

    it('Handles email case sensitivity', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      // Test with uppercase version of the email
      const upperCaseEmail = defaultTurnkeyUser1.email.toUpperCase();
      const user = await TurnkeyUserTable.findByEmail(upperCaseEmail);

      // This should not find the user as emails are case-sensitive in the database
      expect(user).toBeUndefined();
    });

    it('Handles empty email search', async () => {
      await TurnkeyUserTable.create(defaultTurnkeyUser1);

      const userByEmptyEmail = await TurnkeyUserTable.findByEmail('');

      expect(userByEmptyEmail).toBeUndefined();
    });
  });
});
