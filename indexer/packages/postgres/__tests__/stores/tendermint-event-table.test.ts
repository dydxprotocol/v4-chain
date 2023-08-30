import { TendermintEventFromDatabase } from 'packages/postgres/src/types';
import * as BlockTable from '../../src/stores/block-table';
import * as TendermintEventTable from '../../src/stores/tendermint-event-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  defaultBlock, defaultBlock2, defaultTendermintEvent, defaultTendermintEventId,
} from '../helpers/constants';

describe('TendermintEvent store', () => {
  beforeAll(async () => {
    await migrate();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  beforeEach(async () => {
    await Promise.all([
      BlockTable.create(defaultBlock),
      BlockTable.create(defaultBlock2),
    ]);
  });

  it('Successfully creates a TendermintEvent', async () => {
    await TendermintEventTable.create(defaultTendermintEvent);
  });

  it('Successfully finds all TendermintEvents', async () => {
    await Promise.all([
      TendermintEventTable.create(defaultTendermintEvent),
      TendermintEventTable.create({
        ...defaultTendermintEvent,
        blockHeight: '2',
      }),
    ]);

    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {},
      [],
      { readReplica: true },
    );

    expect(tendermintEvents.length).toEqual(2);
    expect(tendermintEvents[0]).toEqual(expect.objectContaining(defaultTendermintEvent));
    expect(tendermintEvents[1]).toEqual(expect.objectContaining({
      ...defaultTendermintEvent,
      blockHeight: '2',
    }));
  });

  it('Successfully finds TendermintEvent with block height', async () => {
    await Promise.all([
      TendermintEventTable.create(defaultTendermintEvent),
      TendermintEventTable.create({
        ...defaultTendermintEvent,
        blockHeight: '2',
      }),
    ]);

    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {
        blockHeight: ['1'],
      },
      [],
      { readReplica: true },
    );

    expect(tendermintEvents.length).toEqual(1);
    expect(tendermintEvents[0]).toEqual(expect.objectContaining(defaultTendermintEvent));
  });

  it('Successfully finds TendermintEvent with block height and transaction index', async () => {
    await Promise.all([
      TendermintEventTable.create(defaultTendermintEvent),
      TendermintEventTable.create({
        ...defaultTendermintEvent,
        blockHeight: '2',
      }),
    ]);

    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {
        blockHeight: ['1'],
        transactionIndex: [-1],
      },
      [],
      { readReplica: true },
    );

    expect(tendermintEvents.length).toEqual(1);
    expect(tendermintEvents[0]).toEqual(expect.objectContaining(defaultTendermintEvent));
  });

  it('Successfully finds a TendermintEvent', async () => {
    await TendermintEventTable.create(defaultTendermintEvent);

    const tendermintEvent: TendermintEventFromDatabase | undefined = await
    TendermintEventTable.findById(
      defaultTendermintEventId,
    );

    expect(tendermintEvent).toEqual(expect.objectContaining(defaultTendermintEvent));
  });

  it('Unable finds a TendermintEvent', async () => {
    const tendermintEvent: TendermintEventFromDatabase | undefined = await
    TendermintEventTable.findById(
      defaultTendermintEventId,
    );
    expect(tendermintEvent).toEqual(undefined);
  });

  it('Event ids are sorted', async () => {
    // Create tendermint events with 1 <= blockHeight <= 10, -2 <= transactionIndex <= 8,
    // 0 <= eventIndex <= 9.
    const promises = [];
    const expectedEventIds: string[] = [];
    for (let blockHeight = 1; blockHeight <= 10; blockHeight += 1) {
      if (blockHeight >= 3) {
        await BlockTable.create({
          ...defaultBlock,
          blockHeight: blockHeight.toString(),
        });
      }
      for (let transactionIndex = -2; transactionIndex <= 8; transactionIndex += 1) {
        for (let eventIndex = 0; eventIndex <= 9; eventIndex += 1) {
          promises.push(
            TendermintEventTable.create({
              blockHeight: blockHeight.toString(),
              transactionIndex,
              eventIndex,
            }),
          );
          expectedEventIds.push(
            TendermintEventTable.createEventId(
              blockHeight.toString(),
              transactionIndex,
              eventIndex,
            ).toString('hex'),
          );
        }
      }
    }
    await Promise.all(promises);
    // by default, order is by id asc.
    const tendermintEvents: TendermintEventFromDatabase[] = await TendermintEventTable.findAll(
      {},
      [],
      { readReplica: true },
    );
    const eventIds: string[] = tendermintEvents.map((event) => event.id.toString('hex'));
    expect(eventIds).toEqual(expectedEventIds);
  });

  it.each([
    ['different block height, a < b', '5', 2, 2, '6', 2, 2, -1],
    ['different transaction index, a < b', '6', 1, 2, '6', 2, 2, -1],
    ['different event index, a < b', '6', 1, 1, '6', 1, 2, -1],
    ['different block height, a > b', '7', 2, 2, '6', 2, 2, 1],
    ['different transaction index, a > b', '6', 3, 2, '6', 2, 2, 1],
    ['different event index, a > b', '6', 1, 3, '6', 1, 2, 1],
    ['a === b', '5', 3, 2, '5', 3, 2, 0],
  ])('Compares tendermint events: %s', (
    _name: string,
    blockHeightA: string,
    transactionIndexA: number,
    eventIndexA: number,
    blockHeightB: string,
    transactionIndexB: number,
    eventIndexB: number,
    expectedResult: number,
  ) => {
    const tendermintEventA: TendermintEventFromDatabase = {
      id: TendermintEventTable.createEventId(blockHeightA, transactionIndexA, eventIndexA),
      blockHeight: blockHeightA,
      transactionIndex: transactionIndexA,
      eventIndex: eventIndexA,
    };
    const tendermintEventB: TendermintEventFromDatabase = {
      id: TendermintEventTable.createEventId(blockHeightB, transactionIndexB, eventIndexB),
      blockHeight: blockHeightB,
      transactionIndex: transactionIndexB,
      eventIndex: eventIndexB,
    };
    expect(TendermintEventTable.compare(tendermintEventA, tendermintEventB))
      .toEqual(expectedResult);
  });
});
