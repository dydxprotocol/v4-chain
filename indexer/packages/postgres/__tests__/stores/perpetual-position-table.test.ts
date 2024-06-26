import { Big } from 'big.js';
import { DateTime } from 'luxon';

import {
  IsoString,
  MarketOpenInterest,
  OrderColumns,
  Ordering,
  PerpetualPositionCreateObject,
  PerpetualPositionFromDatabase,
  PerpetualPositionStatus,
  PerpetualPositionSubaccountUpdateObject,
  PositionSide,
  SubaccountToPerpetualPositionsMap,
  TendermintEventCreateObject,
} from '../../src/types';
import * as PerpetualPositionTable from '../../src/stores/perpetual-position-table';
import * as PerpetualMarketTable from '../../src/stores/perpetual-market-table';
import { clearData, migrate, teardown } from '../../src/helpers/db-helpers';
import { seedData } from '../helpers/mock-generators';
import { ValidationError } from '../../src/lib/errors';
import {
  createdDateTime,
  createdHeight,
  defaultBlock2,
  defaultPerpetualMarket,
  defaultPerpetualPosition,
  defaultSubaccountId,
  defaultSubaccountId2,
  defaultTendermintEvent3,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
} from '../helpers/constants';
import { checkLengthAndContains } from './helpers';
import _ from 'lodash';
import { TendermintEventTable } from '../../src';

describe('PerpetualPosition store', () => {
  beforeAll(async () => {
    await migrate();
  });

  beforeEach(async () => {
    await seedData();
  });

  afterEach(async () => {
    await clearData();
  });

  afterAll(async () => {
    await teardown();
  });

  const defaultPerpetualPositionId: string = PerpetualPositionTable.uuid(
    defaultPerpetualPosition.subaccountId,
    defaultPerpetualPosition.openEventId,
  );

  it('Successfully creates a PerpetualPosition', async () => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);
  });

  it('Successfully creates a PerpetualPosition without optional fields', async () => {
    await PerpetualPositionTable.create({
      subaccountId: defaultSubaccountId,
      perpetualId: defaultPerpetualMarket.id,
      side: PositionSide.LONG,
      status: PerpetualPositionStatus.OPEN,
      size: '10',
      maxSize: '25',
      createdAt: createdDateTime.toISO(),
      createdAtHeight: createdHeight,
      openEventId: defaultTendermintEventId,
      lastEventId: defaultTendermintEventId2,
      settledFunding: '200000',
    });
  });

  it('Successfully finds all PerpetualPositions', async () => {
    await PerpetualMarketTable.create({
      ...defaultPerpetualMarket,
      id: '100',
    });
    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: '100',
        openEventId: defaultTendermintEventId2,
      }),
    ]);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {},
      [],
      {
        readReplica: true,
        orderBy: [[OrderColumns.perpetualId, Ordering.ASC],
          [OrderColumns.openEventId, Ordering.ASC]],
      },
    );

    expect(perpetualPositions.length).toEqual(2);
    expect(perpetualPositions[0]).toEqual(expect.objectContaining(defaultPerpetualPosition));
    expect(perpetualPositions[1]).toEqual(expect.objectContaining({
      ...defaultPerpetualPosition,
      perpetualId: '100',
      openEventId: defaultTendermintEventId2,
    }));
  });

  it('Successfully finds PerpetualPosition with perpetualId', async () => {
    await PerpetualMarketTable.create({
      ...defaultPerpetualMarket,
      id: '100',
    });
    await Promise.all([
      PerpetualPositionTable.create(defaultPerpetualPosition),
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        perpetualId: '100',
        openEventId: defaultTendermintEventId2,
      }),
    ]);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {
        perpetualId: [defaultPerpetualMarket.id],
      },
      [],
      { readReplica: true },
    );

    expect(perpetualPositions.length).toEqual(1);
    expect(perpetualPositions[0]).toEqual(expect.objectContaining(defaultPerpetualPosition));
  });

  it('Successfully finds PerpetualPosition by Subaccount', async () => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {
        subaccountId: [defaultSubaccountId],
      },
      [],
      { readReplica: true },
    );

    expect(perpetualPositions.length).toEqual(1);
    expect(perpetualPositions[0]).toEqual(expect.objectContaining(defaultPerpetualPosition));

    const noPerpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {
        subaccountId: ['6fa6b369-4107-4f0c-bc57-6cbc08cb15a5'],
      },
      [],
      { readReplica: true },
    );

    expect(noPerpetualPositions.length).toEqual(0);
  });

  it.each([
    [1, 1, defaultPerpetualPosition],
    [-1, 0, undefined],
  ])('Successfully finds PerpetualPosition by createdBeforeOrAt, delta %d seconds', async (
    deltaSeconds: number,
    expectedLength: number,
    expectedPosition?: PerpetualPositionCreateObject,
  ) => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {
        createdBeforeOrAt: createdDateTime.plus({ seconds: deltaSeconds }).toISO(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(perpetualPositions, expectedLength, expectedPosition);
  });

  it('Successfully finds PerpetualPositions sorted by openEventId', async () => {
    const earlierPosition: PerpetualPositionCreateObject = {
      ...defaultPerpetualPosition,
      openEventId: defaultTendermintEventId3,
      lastEventId: defaultTendermintEventId3,
    };
    const nextTendermintEvent: TendermintEventCreateObject = {
      blockHeight: defaultTendermintEvent3.blockHeight,
      transactionIndex: defaultTendermintEvent3.transactionIndex,
      eventIndex: defaultTendermintEvent3.eventIndex + 1,
    };
    const nextTendermintEventId: Buffer = TendermintEventTable.createEventId(
      nextTendermintEvent.blockHeight,
      nextTendermintEvent.transactionIndex,
      nextTendermintEvent.eventIndex,
    );
    const laterPosition: PerpetualPositionCreateObject = {
      ...defaultPerpetualPosition,
      openEventId: nextTendermintEventId,
      lastEventId: nextTendermintEventId,
    };
    await TendermintEventTable.create(nextTendermintEvent);
    await Promise.all([
      await PerpetualPositionTable.create(earlierPosition),
      await PerpetualPositionTable.create(laterPosition),
    ]);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll({}, [], { readReplica: true });

    expect(perpetualPositions.length).toEqual(2);
    expect(perpetualPositions[0]).toEqual(expect.objectContaining(laterPosition));
    expect(perpetualPositions[1]).toEqual(expect.objectContaining(earlierPosition));
  });

  it.each([
    [1, 1, defaultPerpetualPosition],
    [-1, 0, undefined],
  ])('Successfully finds PerpetualPosition by createdBeforeOrAtHeight, delta %d blocks', async (
    deltaHeight: number,
    expectedLength: number,
    expectedPosition?: PerpetualPositionCreateObject,
  ) => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const perpetualPositions: PerpetualPositionFromDatabase[] = await
    PerpetualPositionTable.findAll(
      {
        createdBeforeOrAtHeight: Big(createdHeight).plus(deltaHeight).toFixed(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(perpetualPositions, expectedLength, expectedPosition);
  });

  it('Successfully finds a PerpetualPosition', async () => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(defaultPerpetualPositionId);

    expect(perpetualPosition).toEqual(expect.objectContaining(defaultPerpetualPosition));
  });

  it('Unable finds a PerpetualPosition', async () => {
    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.findById(defaultPerpetualPositionId);
    expect(perpetualPosition).toEqual(undefined);
  });

  it('Successfully updates a perpetualPosition', async () => {
    await PerpetualPositionTable.create(defaultPerpetualPosition);

    const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
    PerpetualPositionTable.update({
      id: defaultPerpetualPositionId,
      size: '20',
    });

    expect(perpetualPosition).toEqual(expect.objectContaining({
      ...defaultPerpetualPosition,
      size: '20',
    }));
  });

  describe('findOpenPositionForSubaccountPerpetual', () => {
    it('Successfully gets the open position for a subaccountId and perpetualId', async () => {
      await PerpetualPositionTable.create(defaultPerpetualPosition);

      const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
      PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        defaultPerpetualPosition.subaccountId,
        defaultPerpetualPosition.perpetualId,
      );

      expect(perpetualPosition).toEqual(expect.objectContaining(defaultPerpetualPosition));
    });

    it('Successfully gets no open positions for a subaccountId and perpetualId', async () => {
      await PerpetualPositionTable.create(defaultPerpetualPosition);

      const otherPerpetualId: string = '3';
      const perpetualPosition: PerpetualPositionFromDatabase | undefined = await
      PerpetualPositionTable.findOpenPositionForSubaccountPerpetual(
        defaultPerpetualPosition.subaccountId,
        otherPerpetualId,
      );

      expect(perpetualPosition).toEqual(undefined);
    });
  });

  describe('findOpenPositionsForSubaccount', () => {
    it('Successfully gets the open positions for subaccountIds', async () => {
      await Promise.all([
        PerpetualMarketTable.create({
          ...defaultPerpetualMarket,
          id: '100',
        }),
        PerpetualMarketTable.create({
          ...defaultPerpetualMarket,
          id: '101',
        }),
      ]);
      const perpetualPosition2: PerpetualPositionCreateObject = {
        ...defaultPerpetualPosition,
        perpetualId: '100',
        openEventId: defaultTendermintEventId2,
      };
      const perpetualPosition3: PerpetualPositionCreateObject = {
        ...defaultPerpetualPosition,
        subaccountId: defaultSubaccountId2,
        perpetualId: '101',
        openEventId: defaultTendermintEventId2,
      };
      const perpetualPosition4: PerpetualPositionCreateObject = {
        ...defaultPerpetualPosition,
        subaccountId: defaultSubaccountId2,
        perpetualId: '100',
        openEventId: defaultTendermintEventId,
        status: PerpetualPositionStatus.CLOSED,
      };
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create(perpetualPosition2),
        PerpetualPositionTable.create(perpetualPosition3),
        PerpetualPositionTable.create(perpetualPosition4),
      ]);

      const perpetualPositions: SubaccountToPerpetualPositionsMap = await
      PerpetualPositionTable.findOpenPositionsForSubaccounts(
        [
          defaultSubaccountId,
          defaultSubaccountId2,
        ],
        {},
      );

      expect(perpetualPositions).toEqual(expect.objectContaining({
        [defaultSubaccountId]: {
          [defaultPerpetualMarket.id]: expect.objectContaining(defaultPerpetualPosition),
          [perpetualPosition2.perpetualId]: expect.objectContaining(perpetualPosition2),
        },
        [defaultSubaccountId2]: {
          [perpetualPosition3.perpetualId]: expect.objectContaining(perpetualPosition3),
        },
      }));
    });

    it('Successfully gets no open positions for a subaccountId', async () => {
      const perpetualPositions: SubaccountToPerpetualPositionsMap = await
      PerpetualPositionTable.findOpenPositionsForSubaccounts(
        [defaultPerpetualPosition.subaccountId],
      );

      expect(perpetualPositions).toEqual({});
    });
  });

  describe('closePosition', () => {
    it('Successfully able to close position', async () => {
      const perpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create(defaultPerpetualPosition);

      const closedAt: IsoString = DateTime.utc().toISO();
      const closedAtHeight: string = '2';
      const closeEventId: Buffer = defaultTendermintEventId3;
      const settledFunding: string = '300000';
      const closedPerpetualPosition: PerpetualPositionFromDatabase | undefined = await
      PerpetualPositionTable.closePosition(
        perpetualPosition,
        {
          id: defaultPerpetualPositionId,
          closedAt,
          closedAtHeight,
          closeEventId,
          settledFunding,
        },
      );

      expect(closedPerpetualPosition).toEqual(
        {
          ...perpetualPosition,
          status: PerpetualPositionStatus.CLOSED,
          closedAt,
          closedAtHeight,
          closeEventId,
          lastEventId: closeEventId,
          settledFunding,
          size: '0',
        },
      );
    });

    it(
      'Successfully able to close position, with fixed-point notation size/price values',
      async () => {
        // These values will be converted to strings in exponential notation by the big.js library
        // if `toString` is used instead of `toFixed`
        const tinySize: string = '0.0000001';
        const tinyMaxSize: string = '0.00000025';
        const tinyPrice: string = '0.00000003';
        const perpetualPosition: PerpetualPositionFromDatabase = await
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          size: tinySize,
          maxSize: tinyMaxSize,
          sumOpen: tinySize,
          entryPrice: tinyPrice,
        });

        const closedAt: IsoString = DateTime.utc().toISO();
        const closedAtHeight: string = '2';
        const closeEventId: Buffer = defaultTendermintEventId3;
        const settledFunding: string = '0.000000035';
        const closedPerpetualPosition: PerpetualPositionFromDatabase | undefined = await
        PerpetualPositionTable.closePosition(
          perpetualPosition,
          {
            id: defaultPerpetualPositionId,
            closedAt,
            closedAtHeight,
            closeEventId,
            settledFunding,
          },
        );

        expect(closedPerpetualPosition).toEqual(
          {
            ...perpetualPosition,
            status: PerpetualPositionStatus.CLOSED,
            closedAt,
            closedAtHeight,
            closeEventId,
            lastEventId: closeEventId,
            settledFunding,
            size: '0',
          },
        );
      });

    it('Unable to close position when position is already closed', async () => {
      const perpetualPosition: PerpetualPositionFromDatabase = await
      PerpetualPositionTable.create({
        ...defaultPerpetualPosition,
        status: PerpetualPositionStatus.CLOSED,
      });

      await expect(PerpetualPositionTable.closePosition(
        perpetualPosition,
        {
          id: defaultPerpetualPositionId,
          closedAt: defaultPerpetualPosition.createdAt,
          closedAtHeight: defaultPerpetualPosition.createdAtHeight,
          closeEventId: defaultPerpetualPosition.openEventId,
          settledFunding: defaultPerpetualPosition.settledFunding,
        },
      )).rejects.toThrow(new ValidationError('Unable to close because position is closed'));
    });
  });

  describe('getOpenInterestLong', () => {
    it('Successfully gets open interest long with positions', async () => {
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({ // this position should be ignored
          ...defaultPerpetualPosition,
          side: PositionSide.SHORT,
          openEventId: defaultTendermintEventId2,
        }),
        PerpetualPositionTable.create({ // this position should be ignored
          ...defaultPerpetualPosition,
          perpetualId: '1',
          side: PositionSide.SHORT,
          openEventId: defaultTendermintEventId3,
        }),
      ]);

      // defaultPerpetualPosition.createdAt is the current time the object is created,
      // so which should be in the last 24 before this function is called
      const marketOpenInterest:
      _.Dictionary<MarketOpenInterest> = await PerpetualPositionTable.getOpenInterestLong(
        [defaultPerpetualPosition.perpetualId],
      );

      expect(marketOpenInterest).toEqual({
        [defaultPerpetualPosition.perpetualId]: {
          perpetualMarketId: defaultPerpetualPosition.perpetualId,
          openInterest: defaultPerpetualPosition.size,
        },
      });
    });

    it('Gets default data when there are no matching positions', async () => {
      const fakePerpetualId = '2';
      await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({ // this position should be ignored
          ...defaultPerpetualPosition,
          side: PositionSide.SHORT,
          openEventId: defaultTendermintEventId2,
        }),
        PerpetualPositionTable.create({ // this position should be ignored
          ...defaultPerpetualPosition,
          perpetualId: '1',
          side: PositionSide.SHORT,
          openEventId: defaultTendermintEventId3,
        }),
      ]);

      const marketOpenInterest:
      _.Dictionary<MarketOpenInterest> = await PerpetualPositionTable.getOpenInterestLong([
        fakePerpetualId,
      ]);

      expect(marketOpenInterest).toEqual({
        [fakePerpetualId]: { perpetualMarketId: fakePerpetualId, openInterest: '0' },
      });
    });

    it('Successfully gets open interest long with no positions', async () => {
      const marketOpenInterest:
      _.Dictionary<MarketOpenInterest> = await PerpetualPositionTable.getOpenInterestLong(
        [defaultPerpetualMarket.id],
      );

      expect(marketOpenInterest).toEqual({
        [defaultPerpetualMarket.id]: {
          perpetualMarketId: defaultPerpetualMarket.id,
          openInterest: '0',
        },
      });
    });
  });

  describe('bulkCreate', () => {
    it('Successfully creates multiple positions', async () => {
      const createdPositions:
      PerpetualPositionFromDatabase[] = await PerpetualPositionTable.bulkCreate([
        defaultPerpetualPosition,
        {
          ...defaultPerpetualPosition,
          side: PositionSide.SHORT,
          openEventId: defaultTendermintEventId2,
        },
      ]);

      expect(createdPositions).toHaveLength(2);
      for (let i = 0; i < createdPositions.length; i += 1) {
        const position: PerpetualPositionFromDatabase = createdPositions[i];
        expect(
          await PerpetualPositionTable.findById(position.id),
        ).toEqual(position);
      }
    });
  });

  describe('bulkUpdateSubaccountFields', () => {
    it.each([
      [
        'with no maxSize update',
        {
          id: defaultPerpetualPositionId,
          lastEventId: defaultPerpetualPosition.lastEventId,
          settledFunding: '0',
          status: PerpetualPositionStatus.CLOSED,
          size: defaultPerpetualPosition.maxSize,
        },
      ],
      [
        'with maxSize updated',
        {
          id: defaultPerpetualPositionId,
          lastEventId: defaultPerpetualPosition.lastEventId,
          settledFunding: '0',
          status: PerpetualPositionStatus.CLOSED,
          size: Big(defaultPerpetualPosition.maxSize).plus(10).toString(),
        },
      ],
      [
        'with all fields',
        {
          id: defaultPerpetualPositionId,
          lastEventId: defaultPerpetualPosition.lastEventId,
          settledFunding: '0',
          status: PerpetualPositionStatus.CLOSED,
          size: Big(defaultPerpetualPosition.maxSize).plus(10).toString(),
          closedAtHeight: defaultBlock2.blockHeight,
          closedAt: defaultBlock2.time,
          closeEventId: defaultPerpetualPosition.lastEventId,
        },
      ],
    ])('Successfully updates a position %s', async (
      _name: string,
      updateObject: PerpetualPositionSubaccountUpdateObject,
    ) => {
      const position: PerpetualPositionFromDatabase = await PerpetualPositionTable.create(
        defaultPerpetualPosition,
      );
      await PerpetualPositionTable.bulkUpdateSubaccountFields([updateObject]);

      expect(await PerpetualPositionTable.findById(position.id)).toEqual(
        expect.objectContaining(updateObject),
      );
    });

    it('Successfully processes no updates', async () => {
      expect(await PerpetualPositionTable.bulkUpdateSubaccountFields([])).toBeUndefined();
    });

    it('Successfully updates multiple positions', async () => {
      const [position, secondPosition]: PerpetualPositionFromDatabase[] = await Promise.all([
        PerpetualPositionTable.create(defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...defaultPerpetualPosition,
          openEventId: defaultTendermintEventId2,
        }),
      ]);

      const updateObject: PerpetualPositionSubaccountUpdateObject = {
        id: position.id,
        lastEventId: position.lastEventId,
        settledFunding: '0',
        status: PerpetualPositionStatus.CLOSED,
        size: position.maxSize,
      };
      const secondUpdateObject: PerpetualPositionSubaccountUpdateObject = {
        ...updateObject,
        id: secondPosition.id,
        size: Big(position.maxSize).plus(10).toString(),
      };

      await PerpetualPositionTable.bulkUpdateSubaccountFields([updateObject, secondUpdateObject]);
      expect(await PerpetualPositionTable.findById(position.id)).toEqual(
        expect.objectContaining(updateObject),
      );
      expect(await PerpetualPositionTable.findById(secondPosition.id)).toEqual(
        expect.objectContaining(secondUpdateObject),
      );
    });
  });
});
