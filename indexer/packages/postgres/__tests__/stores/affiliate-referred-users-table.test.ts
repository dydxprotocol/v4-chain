import { AffiliateReferredUserFromDatabase, AffiliateReferredUsersCreateObject } from '../../src/types';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { defaultAffiliateReferredUser } from '../helpers/constants';
import * as AffiliateReferredUsersTable from '../../src/stores/affiliate-referred-users-table';

describe('AffiliateReferredUsers store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates affiliate referee pairs', async () => {
    await AffiliateReferredUsersTable.create(defaultAffiliateReferredUser);
    await AffiliateReferredUsersTable.create({
      ...defaultAffiliateReferredUser,
      refereeAddress: 'fake_address',
    });
  });

  it('Should not allow duplicate refree address', async () => {
    await AffiliateReferredUsersTable.create(defaultAffiliateReferredUser);

    // Second creation should fail due to the duplicate refereeAddress
    await expect(
      AffiliateReferredUsersTable.create({
        ...defaultAffiliateReferredUser,
        affiliateAddress: 'another_affiliate_address',
      }),
    ).rejects.toThrow();
  });

  it('Successfully finds all entries', async () => {
    const entry1: AffiliateReferredUsersCreateObject = {
      ...defaultAffiliateReferredUser,
      refereeAddress: 'referee_address1',
    };
    const entry2: AffiliateReferredUsersCreateObject = {
      ...defaultAffiliateReferredUser,
      affiliateAddress: 'affiliate_address1',
      refereeAddress: 'referee_address2',
    };

    await Promise.all([
      AffiliateReferredUsersTable.create(defaultAffiliateReferredUser),
      AffiliateReferredUsersTable.create(entry1),
      AffiliateReferredUsersTable.create(entry2),
    ]);

    const entries: AffiliateReferredUserFromDatabase[] = await AffiliateReferredUsersTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(entries.length).toEqual(3);
    expect(entries).toEqual(
      expect.arrayContaining([
        expect.objectContaining(defaultAffiliateReferredUser),
        expect.objectContaining(entry1),
        expect.objectContaining(entry2),
      ]),
    );
  });

  it('Successfully finds entries by affiliate address', async () => {
    const entry1: AffiliateReferredUsersCreateObject = {
      affiliateAddress: 'affiliate_address1',
      refereeAddress: 'referee_address1',
      referredAtBlock: '1',
    };
    const entry2: AffiliateReferredUsersCreateObject = {
      affiliateAddress: 'affiliate_address1',
      refereeAddress: 'referee_address2',
      referredAtBlock: '20',
    };

    await AffiliateReferredUsersTable.create(entry1);
    await AffiliateReferredUsersTable.create(entry2);

    const entries: AffiliateReferredUserFromDatabase[] | undefined = await AffiliateReferredUsersTable.findByAffiliateAddress('affiliate_address1');

    if (entries) {
      expect(entries.length).toEqual(2);
      expect(entries).toEqual(
        expect.arrayContaining([
          expect.objectContaining(entry1),
          expect.objectContaining(entry2),
        ]),
      );
    } else {
      throw new Error('findByAffiliateAddress returned undefined, expected an array');
    }
  });

  it('Successfully finds entry by referee address', async () => {
    const entry1: AffiliateReferredUsersCreateObject = {
      affiliateAddress: 'affiliate_address1',
      refereeAddress: 'referee_address1',
      referredAtBlock: '1',
    };
    const entry2: AffiliateReferredUsersCreateObject = {
      affiliateAddress: 'affiliate_address1',
      refereeAddress: 'referee_address2',
      referredAtBlock: '20',
    };

    await AffiliateReferredUsersTable.create(entry1);
    await AffiliateReferredUsersTable.create(entry2);

    const entry: AffiliateReferredUserFromDatabase | undefined = await AffiliateReferredUsersTable.findByRefereeAddress('referee_address1');

    expect(entry).toEqual(expect.objectContaining(entry1));
  });
});
