import { TokenFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAddress2, defaultToken, defaultWallet } from '../helpers/constants';
import * as TokenTable from '../../src/stores/token-table';
import * as WalletTable from '../../src/stores/wallet-table';

describe('Token store', () => {
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
    await TokenTable.create(defaultToken);
    const token = await TokenTable.findByToken(defaultToken.token);
    expect(token).toEqual(expect.objectContaining(defaultToken));
  });

  it('Successfully upserts a Token multiple times', async () => {
    await TokenTable.upsert(defaultToken);
    let token: TokenFromDatabase | undefined = await TokenTable.findByToken(
      defaultToken.token,
    );

    expect(token).toEqual(expect.objectContaining(defaultToken));

    // Upsert again to test update functionality
    const updatedToken = { ...defaultToken, updatedAt: new Date().toISOString(), language: 'es' };
    await TokenTable.upsert(updatedToken);
    token = await TokenTable.findByToken(defaultToken.token);

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
      TokenTable.create(defaultToken),
      TokenTable.create(additionalToken),
    ]);

    const tokens: TokenFromDatabase[] = await TokenTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(tokens.length).toEqual(2);
    expect(tokens[0]).toEqual(expect.objectContaining(defaultToken));
    expect(tokens[1]).toEqual(expect.objectContaining(additionalToken));
  });

  it('Successfully finds a Token by token', async () => {
    await TokenTable.create(defaultToken);

    const token: TokenFromDatabase | undefined = await TokenTable.findByToken(
      defaultToken.token,
    );

    expect(token).toEqual(expect.objectContaining(defaultToken));
  });
});
