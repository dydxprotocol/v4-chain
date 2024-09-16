import * as VaultTable from '../../src/stores/vault-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import { defaultVault, defaultAddress } from '../helpers/constants';
import { VaultFromDatabase, VaultStatus } from '../../src/types';

describe('Vault store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a vault', async () => {
    await VaultTable.create(defaultVault);
  });

  it('Successfully finds all vaults', async () => {
    await Promise.all([
      VaultTable.create(defaultVault),
      VaultTable.create({
        ...defaultVault,
        address: defaultAddress,
        clobPairId: '1',
      }),
    ]);

    const vaults: VaultFromDatabase[] = await VaultTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(vaults.length).toEqual(2);
    expect(vaults[0]).toEqual(expect.objectContaining(defaultVault));
    expect(vaults[1]).toEqual(expect.objectContaining({
      ...defaultVault,
      address: defaultAddress,
      clobPairId: '1',
    }));
  });

  it('Succesfully upserts a vault', async () => {
    await VaultTable.create(defaultVault);

    let vaults: VaultFromDatabase[] = await VaultTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(vaults.length).toEqual(1);
    expect(vaults[0]).toEqual(expect.objectContaining(defaultVault));

    await VaultTable.upsert({
      ...defaultVault,
      status: VaultStatus.CLOSE_ONLY,
    });

    vaults = await VaultTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(vaults.length).toEqual(1);
    expect(vaults[0]).toEqual(expect.objectContaining({
      ...defaultVault,
      status: VaultStatus.CLOSE_ONLY,
    }));
  });
});
