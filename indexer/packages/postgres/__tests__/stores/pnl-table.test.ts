import {
  Ordering,
  PnlColumns,
} from '../../src/types';
import * as PnlTable from '../../src/stores/pnl-table';
import * as BlockTable from '../../src/stores/block-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import {
  defaultBlock,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultPnl,
  defaultPnl2,
} from '../helpers/constants';

describe('Pnl store', () => {
  beforeEach(async () => {
    await seedData();
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

  it('Successfully creates a Pnl record', async () => {
    await PnlTable.create(defaultPnl);
  });

  it('Successfully creates multiple Pnl records', async () => {
    await BlockTable.create({
      ...defaultBlock,
      blockHeight: '5',
    });
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create(defaultPnl2),
    ]);

    const { results: pnls } = await PnlTable.findAll({}, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.ASC]],
    });

    expect(pnls.length).toEqual(2);
    expect(pnls[0]).toEqual(expect.objectContaining(defaultPnl));
    expect(pnls[1]).toEqual(expect.objectContaining(defaultPnl2));
  });

  it('Successfully finds Pnl records with subaccountId', async () => {
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create({
        ...defaultPnl,
        createdAt: '2022-06-01T01:00:00.000Z',
      }),
      PnlTable.create({
        ...defaultPnl,
        subaccountId: defaultSubaccountId2,
        createdAt: '2022-06-01T00:00:00.000Z',
      }),
    ]);

    const { results: pnls } = await PnlTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(pnls.length).toEqual(2);
  });

  it('Successfully finds Pnl records using pagination', async () => {
    await Promise.all([
      PnlTable.create(defaultPnl),
      PnlTable.create({
        ...defaultPnl,
        createdAt: '2020-01-01T00:00:00.000Z',
        createdAtHeight: '1000',
      }),
    ]);

    const responsePageOne = await PnlTable.findAll({
      page: 1,
      limit: 1,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining({
      ...defaultPnl,
      createdAt: '2020-01-01T00:00:00.000Z',
      createdAtHeight: '1000',
    }));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo = await PnlTable.findAll({
      page: 2,
      limit: 1,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining(defaultPnl));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages = await PnlTable.findAll({
      page: 1,
      limit: 2,
    }, [], {
      orderBy: [[PnlColumns.createdAtHeight, Ordering.DESC]],
    });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining({
      ...defaultPnl,
      createdAt: '2020-01-01T00:00:00.000Z',
      createdAtHeight: '1000',
    }));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining(defaultPnl));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
  });

});
