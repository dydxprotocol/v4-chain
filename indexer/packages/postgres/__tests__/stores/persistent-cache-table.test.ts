import { PersistentCacheFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultKV, defaultKV2 } from '../helpers/constants';
import * as PersistentCacheTable from '../../src/stores/persistent-cache-table';

describe('Persistent cache store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a key value pair', async () => {
    await PersistentCacheTable.create(defaultKV);
  });

  it('Successfully upserts a kv pair multiple times', async () => {
    const newKv = {
      ...defaultKV,
      value: 'someOtherValue',
    };
    await PersistentCacheTable.upsert(newKv);
    let kv: PersistentCacheFromDatabase | undefined = await PersistentCacheTable.findById(
      defaultKV.key,
    );
    expect(kv).toEqual(expect.objectContaining(newKv));

    const newKv2 = {
      ...defaultKV,
      value: 'someOtherValue2',
    };
    await PersistentCacheTable.upsert(newKv2);
    kv = await PersistentCacheTable.findById(defaultKV.key);

    expect(kv).toEqual(expect.objectContaining(newKv2));
  });

  it('Successfully finds all kv pairs', async () => {
    await Promise.all([
      PersistentCacheTable.create(defaultKV),
      PersistentCacheTable.create(defaultKV2),
    ]);

    const kvs: PersistentCacheFromDatabase[] = await PersistentCacheTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(kvs.length).toEqual(2);
    expect(kvs).toEqual(expect.arrayContaining([
      expect.objectContaining(defaultKV),
      expect.objectContaining(defaultKV2),
    ]));
  });

  it('Successfully finds a kv pair', async () => {
    await PersistentCacheTable.create(defaultKV);

    const kv: PersistentCacheFromDatabase | undefined = await PersistentCacheTable.findById(
      defaultKV.key,
    );

    expect(kv).toEqual(expect.objectContaining(defaultKV));
  });
});
