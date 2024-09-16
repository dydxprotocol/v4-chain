import { AffiliateInfoFromDatabase } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAffiliateInfo, defaultAffiliateInfo2 } from '../helpers/constants';
import * as AffiliateInfoTable from '../../src/stores/affiliate-info-table';

describe('Affiliate info store', () => {
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

  it('Cannot create duplicate info for duplicate address', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);
    await expect(AffiliateInfoTable.create(defaultAffiliateInfo)).rejects.toThrowError();
  });

  it('Can upsert affiliate info multiple times', async () => {
    await AffiliateInfoTable.upsert(defaultAffiliateInfo);
    let info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));

    await AffiliateInfoTable.upsert(defaultAffiliateInfo2);
    info = await AffiliateInfoTable.findById(defaultAffiliateInfo2.address);
    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo2));
  });

  it('Successfully finds all affiliate infos', async () => {
    await Promise.all([
      AffiliateInfoTable.create(defaultAffiliateInfo),
      AffiliateInfoTable.create(defaultAffiliateInfo2),
    ]);

    const infos: AffiliateInfoFromDatabase[] = await AffiliateInfoTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(infos.length).toEqual(2);
    expect(infos).toEqual(expect.arrayContaining([
      expect.objectContaining(defaultAffiliateInfo),
      expect.objectContaining(defaultAffiliateInfo2),
    ]));
  });

  it('Successfully finds an affiliate info', async () => {
    await AffiliateInfoTable.create(defaultAffiliateInfo);

    const info: AffiliateInfoFromDatabase | undefined = await AffiliateInfoTable.findById(
      defaultAffiliateInfo.address,
    );

    expect(info).toEqual(expect.objectContaining(defaultAffiliateInfo));
  });

  describe('paginatedFindWithAddressFilter', () => {
    beforeEach(async () => {
      await migrate();
      for (let i = 0; i < 10; i++) {
        await AffiliateInfoTable.create({
          ...defaultAffiliateInfo,
          address: `address_${i}`,
          affiliateEarnings: i.toString(),
        });
      }
    });

    it('Successfully filters by address', async () => {
      // eslint-disable-next-line max-len
      const infos: AffiliateInfoFromDatabase[] | undefined = await AffiliateInfoTable.paginatedFindWithAddressFilter(
        ['address_0'],
        0,
        10,
        false,
      );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(1);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_0',
        affiliateEarnings: '0',
      }));
    });

    it('Successfully sorts by affiliate earning', async () => {
      // eslint-disable-next-line max-len
      const infos: AffiliateInfoFromDatabase[] | undefined = await AffiliateInfoTable.paginatedFindWithAddressFilter(
        [],
        0,
        10,
        true,
      );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(10);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_9',
        affiliateEarnings: '9',
      }));
      expect(infos![9]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_0',
        affiliateEarnings: '0',
      }));
    });

    it('Successfully uses offset and limit', async () => {
      // eslint-disable-next-line max-len
      const infos: AffiliateInfoFromDatabase[] | undefined = await AffiliateInfoTable.paginatedFindWithAddressFilter(
        [],
        5,
        2,
        false,
      );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(2);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_5',
        affiliateEarnings: '5',
      }));
      expect(infos![1]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_6',
        affiliateEarnings: '6',
      }));
    });

    it('Successfully filters, sorts, offsets, and limits', async () => {
      // eslint-disable-next-line max-len
      const infos: AffiliateInfoFromDatabase[] | undefined = await AffiliateInfoTable.paginatedFindWithAddressFilter(
        [],
        3,
        2,
        true,
      );
      expect(infos).toBeDefined();
      expect(infos!.length).toEqual(2);
      expect(infos![0]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_6',
        affiliateEarnings: '6',
      }));
      expect(infos![1]).toEqual(expect.objectContaining({
        ...defaultAffiliateInfo,
        address: 'address_5',
        affiliateEarnings: '5',
      }));
    });
  });
});
