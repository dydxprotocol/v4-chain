import { FirebaseNotificationTokenFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAddress2, defaultFirebaseNotificationToken, defaultWallet } from '../helpers/constants';
import * as FirebaseNotificationTokenTable from '../../src/stores/firebase-notification-token-table';
import * as WalletTable from '../../src/stores/wallet-table';

describe('FirebaseNotificationToken store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    // Default wallet is required in the DB for token creation
    // As token has a foreign key constraint on wallet
    await WalletTable.create(defaultWallet);
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a Token', async () => {
    await FirebaseNotificationTokenTable.create(defaultFirebaseNotificationToken);
    const token = await FirebaseNotificationTokenTable.findByToken(
      defaultFirebaseNotificationToken.token,
    );
    expect(token).toEqual(expect.objectContaining(defaultFirebaseNotificationToken));
  });

  it('Successfully upserts a Token multiple times', async () => {
    await FirebaseNotificationTokenTable.upsert(defaultFirebaseNotificationToken);
    let token: FirebaseNotificationTokenFromDatabase | undefined = await
    FirebaseNotificationTokenTable.findByToken(
      defaultFirebaseNotificationToken.token,
    );

    expect(token).toEqual(expect.objectContaining(defaultFirebaseNotificationToken));

    // Upsert again to test update functionality
    const updatedToken = { ...defaultFirebaseNotificationToken, updatedAt: new Date().toISOString(), language: 'es' };
    await FirebaseNotificationTokenTable.upsert(updatedToken);
    token = await FirebaseNotificationTokenTable.findByToken(
      defaultFirebaseNotificationToken.token,
    );

    expect(token).toEqual(expect.objectContaining(updatedToken));
  });

  it('Successfully finds all Tokens', async () => {
    await WalletTable.create({ ...defaultWallet, address: defaultAddress2 });
    const additionalToken = {
      token: 'fake_token',
      address: defaultAddress2,
      language: 'en',
      updatedAt: new Date().toISOString(),
    };

    await Promise.all([
      FirebaseNotificationTokenTable.create(defaultFirebaseNotificationToken),
      FirebaseNotificationTokenTable.create(additionalToken),
    ]);

    const tokens: FirebaseNotificationTokenFromDatabase[] = await FirebaseNotificationTokenTable
      .findAll(
        {},
        [],
        { readReplica: true },
      );

    expect(tokens.length).toEqual(2);
    expect(tokens[0]).toEqual(expect.objectContaining(defaultFirebaseNotificationToken));
    expect(tokens[1]).toEqual(expect.objectContaining(additionalToken));
  });

  it('Successfully finds a Token by token', async () => {
    await FirebaseNotificationTokenTable.create(defaultFirebaseNotificationToken);

    const token: FirebaseNotificationTokenFromDatabase | undefined = await
    FirebaseNotificationTokenTable.findByToken(
      defaultFirebaseNotificationToken.token,
    );

    expect(token).toEqual(expect.objectContaining(defaultFirebaseNotificationToken));
  });

  describe('deleteMany', () => {
    it('should delete multiple tokens successfully', async () => {
      const token1 = { ...defaultFirebaseNotificationToken, token: 'token1todelete' };
      const token2 = { ...defaultFirebaseNotificationToken, token: 'token2todelete' };
      const token3 = { ...defaultFirebaseNotificationToken, token: 'token3todelete' };

      await Promise.all([
        FirebaseNotificationTokenTable.create(token1),
        FirebaseNotificationTokenTable.create(token2),
        FirebaseNotificationTokenTable.create(token3),
      ]);

      // Delete the tokens
      const tokensToDelete = ['token1todelete', 'token2todelete', 'token3todelete'];
      await FirebaseNotificationTokenTable.deleteMany(tokensToDelete);

      // Check if the tokens were deleted
      const remainingTokens = await FirebaseNotificationTokenTable.findAll({}, []);
      expect(remainingTokens.length).toEqual(0);
    });

    it('should handle an empty array of token IDs', async () => {
      await FirebaseNotificationTokenTable.create(defaultFirebaseNotificationToken);
      const result = await FirebaseNotificationTokenTable.deleteMany([]);
      expect(result).toEqual(0);

      // Verify the token still exists
      const remainingTokens = await FirebaseNotificationTokenTable.findAll({}, []);
      expect(remainingTokens.length).toEqual(1);
      expect(remainingTokens[0]).toEqual(expect.objectContaining(defaultFirebaseNotificationToken));
    });
  });
});
