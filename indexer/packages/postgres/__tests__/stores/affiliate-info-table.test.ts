import { AffiliateInfoFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAffiliateInfo, defaultAffiliateInfo1 } from '../helpers/constants';
import * as AffiliateInfoTable from '../../src/stores/affiliate-info-table';

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

  it('Successfully creates affiliate info', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);
  });

  it('Cannot create duplicate infor for duplicate address', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);
    await expect(AffiliateInfoTable.create(defaultAffiliateInfo)).rejects.toThrowError();
  });

  it('Can upsert affiliate info multiple times', async () => {
    await AffiliateInfoTable.upsert(defaultAffiliateInfo);
    let info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));

    await AffiliateInfoTable.upsert(defaultAffiliateInfo1);
    info = await AffiliateInfoTable.findById(defaultAffiliateInfo1.address);
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo1));
  });

  it('Successfully finds all affiliate infos', async () => {
    await Promise.all([
      AffiliateInfoTable.create(defaultAffiliateInfo),
      AffiliateInfoTable.create(defaultAffiliateInfo1),
    ]);

    const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(infos.length).toEqual(2);
    expect(infos).toEqual(expect.arrayContaining([
      expect.objectContaining(defaultAffiliateInfo),
      expect.objectContaining(defaultAffiliateInfo1),
    ]));
  });

  it('Successfully finds an affiliate info', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);

    const info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );

    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));
  });
});
