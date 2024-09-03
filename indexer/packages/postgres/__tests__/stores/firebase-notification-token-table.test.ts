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
});
