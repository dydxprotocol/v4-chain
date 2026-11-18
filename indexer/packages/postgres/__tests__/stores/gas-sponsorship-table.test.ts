import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import * as GasSponsorshipTable from '../../src/stores/gas-sponsorship-table';
import {
  GasSponsorshipCreateObject,
  GasSponsorshipFromDatabase,
  GasSponsorshipStatus,
} from '../../src/types';

describe('GasSponsorship store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const chainId: string = 'injective-1';
  const granter: string = 'inj1xdrkxgke8qzey9ze4cx48heav6lzkkszuc224r';
  const grantee1: string = 'inj1frwsqymyg3wnfhy4mqrq9cghwwzst0tgsmfal7';
  const grantee2: string = 'inj1v4d4hf84ea7cf74pcgzghhf5ek8q6empm7gejp';

  const sponsorship1: GasSponsorshipCreateObject = {
    chain_id: chainId,
    granter_address: granter,
    grantee_address: grantee1,
    dydx_address: 'dydx1frwsqymyg3wnfhy4mqrq9cghwwzst0tg9l4v2f',
    spend_limit: '324000000000000',
    fee_amount: '60000000000000',
    status: GasSponsorshipStatus.CONFIRMED,
    transaction_hash: 'C554B3ACB4F98161CCFE7938399A73C4D452FBD9EF92F88C4F79E68896905BE4',
    expires_at: '2026-11-19T12:00:00.000Z',
    created_at: '2026-11-18T12:00:00.000Z',
  };

  const sponsorship2: GasSponsorshipCreateObject = {
    ...sponsorship1,
    grantee_address: grantee2,
    // Larger than Number.MAX_SAFE_INTEGER, to check the amounts are not handled as floats.
    spend_limit: '35200000000000000001',
    fee_amount: '1',
    status: GasSponsorshipStatus.PENDING,
    transaction_hash: undefined,
    expires_at: '2026-11-19T13:00:00.000Z',
    created_at: '2026-11-18T13:00:00.000Z',
  };

  const failedSponsorship: GasSponsorshipCreateObject = {
    ...sponsorship1,
    spend_limit: '999',
    fee_amount: '1',
    status: GasSponsorshipStatus.FAILED,
    transaction_hash: undefined,
    expires_at: '2026-11-19T14:00:00.000Z',
    created_at: '2026-11-18T14:00:00.000Z',
  };

  describe('create', () => {
    it('creates a sponsorship and generates an id', async () => {
      const created: GasSponsorshipFromDatabase = await GasSponsorshipTable.create(sponsorship1);

      expect(created).toEqual(expect.objectContaining(sponsorship1));
      expect(created.id).toEqual(expect.any(String));
      expect(created.id.length).toBeGreaterThan(0);
    });

    it('stores amounts larger than a double can represent exactly', async () => {
      const created: GasSponsorshipFromDatabase = await GasSponsorshipTable.create(sponsorship2);
      const found: GasSponsorshipFromDatabase | undefined = await GasSponsorshipTable.findById(
        created.id,
      );

      expect(found?.spend_limit).toEqual('35200000000000000001');
      expect(found?.transaction_hash).toBeNull();
    });

    it('rejects an amount that is not a non-negative integer', async () => {
      await expect(
        GasSponsorshipTable.create({ ...sponsorship1, spend_limit: '-5' }),
      ).rejects.toThrow();
      await expect(
        GasSponsorshipTable.create({ ...sponsorship1, fee_amount: '1.5' }),
      ).rejects.toThrow();
    });

    it('rejects an unknown status', async () => {
      await expect(
        GasSponsorshipTable.create({
          ...sponsorship1,
          status: 'UNKNOWN' as GasSponsorshipStatus,
        }),
      ).rejects.toThrow();
    });

    it('rejects a duplicate transaction hash', async () => {
      await GasSponsorshipTable.create(sponsorship1);
      await expect(
        GasSponsorshipTable.create({ ...sponsorship1, grantee_address: grantee2 }),
      ).rejects.toThrow();
    });
  });

  describe('findById', () => {
    it('returns undefined when there is no such sponsorship', async () => {
      expect(await GasSponsorshipTable.findById('does-not-exist')).toBeUndefined();
    });
  });

  describe('findAll and findLatest', () => {
    beforeEach(async () => {
      await GasSponsorshipTable.create(sponsorship1);
      await GasSponsorshipTable.create(sponsorship2);
      await GasSponsorshipTable.create(failedSponsorship);
    });

    it('returns everything newest first when there are no filters', async () => {
      const results: GasSponsorshipFromDatabase[] = await GasSponsorshipTable.findAll();

      expect(results.map((r) => r.created_at)).toEqual([
        failedSponsorship.created_at,
        sponsorship2.created_at,
        sponsorship1.created_at,
      ]);
    });

    it('filters by grantee and status', async () => {
      const results: GasSponsorshipFromDatabase[] = await GasSponsorshipTable.findAll({
        chain_id: chainId,
        grantee_address: grantee1,
        statuses: [GasSponsorshipStatus.CONFIRMED, GasSponsorshipStatus.PENDING],
      });

      expect(results).toHaveLength(1);
      expect(results[0]).toEqual(expect.objectContaining(sponsorship1));
    });

    it('filters by creation time, inclusive of the bound', async () => {
      const results: GasSponsorshipFromDatabase[] = await GasSponsorshipTable.findAll({
        createdOnOrAfter: sponsorship2.created_at,
      });

      expect(results).toHaveLength(2);
    });

    it('filters by chain', async () => {
      expect(await GasSponsorshipTable.findAll({ chain_id: 'other-chain' })).toHaveLength(0);
    });

    it('applies the limit', async () => {
      expect(await GasSponsorshipTable.findAll({}, 2)).toHaveLength(2);
    });

    it('findLatest returns the newest match', async () => {
      const latest: GasSponsorshipFromDatabase | undefined = await GasSponsorshipTable.findLatest({
        chain_id: chainId,
        grantee_address: grantee1,
      });

      expect(latest).toEqual(expect.objectContaining({
        status: GasSponsorshipStatus.FAILED,
        created_at: failedSponsorship.created_at,
      }));
    });

    it('findLatest returns undefined when nothing matches', async () => {
      expect(
        await GasSponsorshipTable.findLatest({ grantee_address: 'inj1nobody' }),
      ).toBeUndefined();
    });
  });

  describe('getTotalCommittedSince', () => {
    it('returns zero when there are no sponsorships', async () => {
      expect(
        await GasSponsorshipTable.getTotalCommittedSince(chainId, '2026-11-18T00:00:00.000Z'),
      ).toEqual('0');
    });

    it('sums spend limits and fees exactly, leaving out failed sponsorships', async () => {
      await GasSponsorshipTable.create(sponsorship1);
      await GasSponsorshipTable.create(sponsorship2);
      await GasSponsorshipTable.create(failedSponsorship);

      // 324000000000000 + 60000000000000 + 35200000000000000001 + 1
      expect(
        await GasSponsorshipTable.getTotalCommittedSince(chainId, '2026-11-18T00:00:00.000Z'),
      ).toEqual('35200384000000000002');
    });

    it('leaves out sponsorships created before the bound', async () => {
      await GasSponsorshipTable.create(sponsorship1);
      await GasSponsorshipTable.create(sponsorship2);

      expect(
        await GasSponsorshipTable.getTotalCommittedSince(chainId, sponsorship2.created_at),
      ).toEqual('35200000000000000002');
    });

    it('leaves out sponsorships on other chains', async () => {
      await GasSponsorshipTable.create(sponsorship1);

      expect(
        await GasSponsorshipTable.getTotalCommittedSince('other-chain', '2026-11-18T00:00:00.000Z'),
      ).toEqual('0');
    });
  });

  describe('update', () => {
    it('updates the status and transaction hash', async () => {
      const created: GasSponsorshipFromDatabase = await GasSponsorshipTable.create(sponsorship2);

      const updated: GasSponsorshipFromDatabase | undefined = await GasSponsorshipTable.update(
        created.id,
        { status: GasSponsorshipStatus.CONFIRMED, transaction_hash: 'ABCD' },
      );

      expect(updated).toEqual(expect.objectContaining({
        id: created.id,
        status: GasSponsorshipStatus.CONFIRMED,
        transaction_hash: 'ABCD',
        spend_limit: sponsorship2.spend_limit,
      }));
      expect(await GasSponsorshipTable.findById(created.id)).toEqual(updated);
    });

    it('clears the transaction hash, so that another sponsorship can have it', async () => {
      const created: GasSponsorshipFromDatabase = await GasSponsorshipTable.create(sponsorship1);

      const updated: GasSponsorshipFromDatabase | undefined = await GasSponsorshipTable.update(
        created.id,
        { status: GasSponsorshipStatus.FAILED, transaction_hash: null },
      );

      expect(updated?.transaction_hash).toBeNull();
      expect(updated?.status).toEqual(GasSponsorshipStatus.FAILED);
      const reused: GasSponsorshipFromDatabase = await GasSponsorshipTable.create({
        ...sponsorship2,
        transaction_hash: sponsorship1.transaction_hash,
      });
      expect(reused.transaction_hash).toEqual(sponsorship1.transaction_hash);
    });

    it('returns undefined when there is no such sponsorship', async () => {
      expect(
        await GasSponsorshipTable.update('does-not-exist', {
          status: GasSponsorshipStatus.FAILED,
        }),
      ).toBeUndefined();
    });
  });
});
