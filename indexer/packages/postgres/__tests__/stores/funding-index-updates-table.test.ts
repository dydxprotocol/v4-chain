import { FundingIndexMap, FundingIndexUpdatesCreateObject, FundingIndexUpdatesFromDatabase } from '../../src/types';
import * as FundingIndexUpdatesTable from '../../src/stores/funding-index-updates-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultBlock,
  defaultFundingIndexUpdate,
  defaultFundingIndexUpdateId,
  defaultPerpetualMarket,
  defaultPerpetualMarket2,
  defaultPerpetualMarket3,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
} from '../helpers/constants';
import * as BlockTable from '../../src/stores/block-table';
import Big from 'big.js';

describe('funding index update store', () => {
  const updatedHeight: string = '5';

  beforeEach(async () => {
    await seedData();
    await BlockTable.create({
      ...defaultBlock,
      blockHeight: updatedHeight,
    });
  });

  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  it('Successfully creates a funding index update', async () => {
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);
  });

  it('Successfully creates multiple funding index updates', async () => {
    const fundingIndexUpdate2: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      perpetualId: defaultPerpetualMarket2.id,
      eventId: defaultTendermintEventId2,
      rate: '0.00005',
    };
    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdate2),
    ]);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        effectiveAtHeight: defaultFundingIndexUpdate.effectiveAtHeight,
      },
      [],
      {},
    );

    expect(fundingIndexUpdates.length).toEqual(2);
    expect(fundingIndexUpdates[0]).toEqual(expect.objectContaining(fundingIndexUpdate2));
    expect(fundingIndexUpdates[1]).toEqual(expect.objectContaining(defaultFundingIndexUpdate));
  });

  it('Successfully finds all FundingIndexUpdates', async () => {
    const fundingIndexUpdate2: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      eventId: defaultTendermintEventId2,
    };

    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdate2),
    ]);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        effectiveAtHeight: defaultFundingIndexUpdate.effectiveAtHeight,
      },
      [],
      {},
    );

    expect(fundingIndexUpdates.length).toEqual(2);
    expect(fundingIndexUpdates[0]).toEqual(expect.objectContaining(fundingIndexUpdate2));
    expect(fundingIndexUpdates[1]).toEqual(expect.objectContaining(defaultFundingIndexUpdate));
  });

  it('Successfully finds FundingIndexUpdates with effectiveAtHeight', async () => {
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        effectiveAtHeight: defaultFundingIndexUpdate.effectiveAtHeight,
      },
      [],
      { readReplica: true },
    );

    expect(fundingIndexUpdates.length).toEqual(1);
    expect(fundingIndexUpdates[0]).toEqual(expect.objectContaining({
      ...defaultFundingIndexUpdate,
    }));
  });

  it('Successfully finds all FundingIndexUpdates effective before or after the height', async () => {
    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create({
        ...defaultFundingIndexUpdate,
        effectiveAtHeight: updatedHeight,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        effectiveBeforeOrAtHeight: defaultFundingIndexUpdate.effectiveAtHeight,
      },
      [],
      {},
    );

    expect(fundingIndexUpdates.length).toEqual(1);
    expect(fundingIndexUpdates[0]).toEqual(expect.objectContaining(defaultFundingIndexUpdate));
  });

  it('Successfully finds all FundingIndexUpdates effective before or after time', async () => {
    const fundingIndexUpdates2: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      effectiveAtHeight: updatedHeight,
      effectiveAt: '1982-05-25T00:00:00.000Z',
      eventId: defaultTendermintEventId2,
    };
    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdates2),
    ]);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase[] = await
    FundingIndexUpdatesTable.findAll(
      {
        effectiveBeforeOrAt: '2000-05-25T00:00:00.000Z',
      },
      [],
      {},
    );

    expect(fundingIndexUpdates.length).toEqual(1);
    expect(fundingIndexUpdates[0]).toEqual(expect.objectContaining(fundingIndexUpdates2));
  });

  it('Successfully finds a FundingIndexUpdate', async () => {
    await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase | undefined = await
    FundingIndexUpdatesTable.findById(defaultFundingIndexUpdateId);
    expect(fundingIndexUpdates).toEqual(expect.objectContaining(defaultFundingIndexUpdate));
  });

  it('Successfully finds latest funding index update for market id', async () => {
    const fundingIndexUpdates2: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      effectiveAtHeight: updatedHeight,
      effectiveAt: '1982-05-25T00:00:00.000Z',
      eventId: defaultTendermintEventId2,
    };
    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdates2),
    ]);

    const fundingIndexUpdates: FundingIndexUpdatesFromDatabase = await FundingIndexUpdatesTable
      .findMostRecentMarketFundingIndexUpdate(
        defaultPerpetualMarket.id,
      ) as FundingIndexUpdatesFromDatabase;

    expect(fundingIndexUpdates).toEqual(expect.objectContaining(fundingIndexUpdates2));
  });

  it('Successfully finds funding index map effectiveBeforeOrAtHeight', async () => {

    const fundingIndexUpdates2: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      fundingIndex: '124',
      effectiveAtHeight: updatedHeight,
      effectiveAt: '1982-05-25T00:00:00.000Z',
      eventId: defaultTendermintEventId2,
    };
    const fundingIndexUpdates3: FundingIndexUpdatesCreateObject = {
      ...defaultFundingIndexUpdate,
      eventId: defaultTendermintEventId3,
      perpetualId: defaultPerpetualMarket2.id,
    };
    await Promise.all([
      FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
      FundingIndexUpdatesTable.create(fundingIndexUpdates2),
      FundingIndexUpdatesTable.create(fundingIndexUpdates3),
    ]);

    const fundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(
        '3',
      );

    expect(fundingIndexMap).toEqual({
      [defaultFundingIndexUpdate.perpetualId]: Big(defaultFundingIndexUpdate.fundingIndex),
      [fundingIndexUpdates3.perpetualId]: Big(fundingIndexUpdates3.fundingIndex),
      [defaultPerpetualMarket3.id]: Big(0),
    });
  });

  it('Gets default funding index of 0 in funding index map if no funding indexes', async () => {
    const fundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
      .findFundingIndexMap(
        '3',
      );

    expect(fundingIndexMap).toEqual({
      [defaultPerpetualMarket.id]: Big(0),
      [defaultPerpetualMarket2.id]: Big(0),
      [defaultPerpetualMarket3.id]: Big(0),
    });
  });

  it(
    'Gets default funding index of 0 in funding index map if no funding indexes for perpetual',
    async () => {
      await FundingIndexUpdatesTable.create(defaultFundingIndexUpdate);

      const fundingIndexMap: FundingIndexMap = await FundingIndexUpdatesTable
        .findFundingIndexMap(
          '3',
        );

      expect(fundingIndexMap).toEqual({
        [defaultPerpetualMarket.id]: Big(defaultFundingIndexUpdate.fundingIndex),
        [defaultPerpetualMarket2.id]: Big(0),
        [defaultPerpetualMarket3.id]: Big(0),
      });
    },
  );
});
