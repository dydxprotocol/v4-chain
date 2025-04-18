import { Big } from 'big.js';
import { DateTime } from 'luxon';

import {
  FillColumns,
  FillCreateObject,
  FillFromDatabase,
  Liquidity,
  Market24HourTradeVolumes,
  OpenSizeWithFundingIndex,
  OrderedFillsWithFundingIndices,
  Ordering,
  OrderSide,
} from '../../src/types';
import * as BlockTable from '../../src/stores/block-table';
import * as FillTable from '../../src/stores/fill-table';
import * as OraclePriceTable from '../../src/stores/oracle-price-table';
import * as OrderTable from '../../src/stores/order-table';
import * as FundingIndexUpdatesTable from '../../src/stores/funding-index-updates-table';
import {
  clearData,
  migrate,
  teardown,
} from '../../src/helpers/db-helpers';
import {
  seedData,
} from '../helpers/mock-generators';
import {
  createdDateTime,
  createdHeight,
  defaultBlock,
  defaultFill,
  defaultFundingIndexUpdate,
  defaultOraclePrice,
  defaultOraclePrice2,
  defaultOrder,
  defaultPerpetualMarket,
  defaultSubaccountId2,
  defaultTendermintEventId,
  defaultTendermintEventId2,
  defaultTendermintEventId3,
  defaultTendermintEventId4,
} from '../helpers/constants';
import { checkLengthAndContains } from './helpers';
import * as SubaccountTable from '../../src/stores/subaccount-table';

const defaultDateTime: DateTime = DateTime.fromISO('2025-01-01T00:00:00.000Z');

describe('Fill store', () => {
  beforeEach(async () => {
    await seedData();
    await OrderTable.create(defaultOrder);
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

  it('Successfully creates a Fill', async () => {
    await FillTable.create(defaultFill);
  });

  it('Successfully finds all Fills, default ordered by createdAtHeight descending', async () => {
    await Promise.all([
      FillTable.create({
        ...defaultFill,
        createdAtHeight: '1',
        eventId: defaultTendermintEventId2,
      }),
      FillTable.create(defaultFill),
    ]);

    const { results: fills } = await FillTable.findAll({}, [], {});

    expect(fills.length).toEqual(2);
    expect(fills[0]).toEqual(expect.objectContaining(defaultFill));
    expect(fills[1]).toEqual(expect.objectContaining({
      ...defaultFill,
      createdAtHeight: '1',
      eventId: defaultTendermintEventId2,
    }));
  });

  it('Successfully finds all Fills with given ordering', async () => {
    await Promise.all([
      FillTable.create(defaultFill),
      FillTable.create({
        ...defaultFill,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const { results: fills } = await FillTable.findAll({}, [], {
      orderBy: [[FillColumns.eventId, Ordering.DESC]],
    });

    expect(fills.length).toEqual(2);
    expect(fills[0]).toEqual(expect.objectContaining({
      ...defaultFill,
      eventId: defaultTendermintEventId2,
    }));
    expect(fills[1]).toEqual(expect.objectContaining(defaultFill));
  });

  it('Successfully finds fills using pagination', async () => {
    await Promise.all([
      FillTable.create(defaultFill),
      FillTable.create({
        ...defaultFill,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const responsePageOne = await FillTable.findAll({
      page: 1,
      limit: 1,
    }, [], {
      orderBy: [[FillColumns.eventId, Ordering.DESC]],
    });

    expect(responsePageOne.results.length).toEqual(1);
    expect(responsePageOne.results[0]).toEqual(expect.objectContaining({
      ...defaultFill,
      eventId: defaultTendermintEventId2,
    }));
    expect(responsePageOne.offset).toEqual(0);
    expect(responsePageOne.total).toEqual(2);

    const responsePageTwo = await FillTable.findAll({
      page: 2,
      limit: 1,
    }, [], {
      orderBy: [[FillColumns.eventId, Ordering.DESC]],
    });

    expect(responsePageTwo.results.length).toEqual(1);
    expect(responsePageTwo.results[0]).toEqual(expect.objectContaining(defaultFill));
    expect(responsePageTwo.offset).toEqual(1);
    expect(responsePageTwo.total).toEqual(2);

    const responsePageAllPages = await FillTable.findAll({
      page: 1,
      limit: 2,
    }, [], {
      orderBy: [[FillColumns.eventId, Ordering.DESC]],
    });

    expect(responsePageAllPages.results.length).toEqual(2);
    expect(responsePageAllPages.results[0]).toEqual(expect.objectContaining({
      ...defaultFill,
      eventId: defaultTendermintEventId2,
    }));
    expect(responsePageAllPages.results[1]).toEqual(expect.objectContaining(defaultFill));
    expect(responsePageAllPages.offset).toEqual(0);
    expect(responsePageAllPages.total).toEqual(2);
  });

  it('Successfully finds Fill with eventId', async () => {
    await Promise.all([
      FillTable.create(defaultFill),
      FillTable.create({
        ...defaultFill,
        eventId: defaultTendermintEventId2,
      }),
    ]);

    const { results: fills } = await FillTable.findAll(
      {
        eventId: defaultFill.eventId,
      },
      [],
      { readReplica: true },
    );

    expect(fills.length).toEqual(1);
    expect(fills[0]).toEqual(expect.objectContaining(defaultFill));
  });

  it.each([
    [1, 1, defaultFill],
    [-1, 0, undefined],
  ])('Successfuly finds Fill with createdBeforeOrAt, delta %d seconds', async (
    deltaSeconds: number,
    expectedLength: number,
    expectedFill?: FillCreateObject,
  ) => {
    await FillTable.create(defaultFill);

    const { results: fills } = await FillTable.findAll(
      {
        createdBeforeOrAt: createdDateTime.plus({ seconds: deltaSeconds }).toISO(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(fills, expectedLength, expectedFill);
  });

  it.each([
    [1, 1, defaultFill],
    [-1, 0, undefined],
  ])('Successfuly finds Fill with createdBeforeOrAtHeight, delta %d blocks', async (
    deltaBlocks: number,
    expectedLength: number,
    expectedFill?: FillCreateObject,
  ) => {
    await FillTable.create(defaultFill);

    const { results: fills } = await FillTable.findAll(
      {
        createdBeforeOrAtHeight: Big(createdHeight).plus(deltaBlocks).toFixed(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(fills, expectedLength, expectedFill);
  });

  it.each([
    [1, 1, defaultFill],
    [0, 1, defaultFill],
    [-1, 0, undefined],
  ])('Successfuly finds Fill with createdOnOrAfter, delta %d seconds', async (
    deltaSeconds: number,
    expectedLength: number,
    expectedFill?: FillCreateObject,
  ) => {
    await FillTable.create(defaultFill);

    const { results: fills } = await FillTable.findAll(
      {
        createdOnOrAfter: createdDateTime.minus({ seconds: deltaSeconds }).toISO(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(fills, expectedLength, expectedFill);
  });

  it.each([
    [1, 1, defaultFill],
    [0, 1, defaultFill],
    [-1, 0, undefined],
  ])('Successfuly finds Fill with createdOnOrAfterHeight, delta %d blocks', async (
    deltaBlocks: number,
    expectedLength: number,
    expectedFill?: FillCreateObject,
  ) => {
    await FillTable.create(defaultFill);

    const { results: fills } = await FillTable.findAll(
      {
        createdOnOrAfterHeight: Big(createdHeight).minus(deltaBlocks).toFixed(),
      },
      [],
      { readReplica: true },
    );

    checkLengthAndContains(fills, expectedLength, expectedFill);
  });

  it('Successfully finds a Fill', async () => {
    await FillTable.create(defaultFill);

    const fill: FillFromDatabase | undefined = await FillTable.findById(
      FillTable.uuid(defaultFill.eventId, defaultFill.liquidity),
    );

    expect(fill).toEqual(expect.objectContaining(defaultFill));
  });

  // TODO: Add a bunch of tests for different search parameters
  it('Successfully updates an Fill', async () => {
    await FillTable.create(defaultFill);

    const fill: FillFromDatabase | undefined = await FillTable.update({
      id: FillTable.uuid(defaultFill.eventId, defaultFill.liquidity),
      size: '32.50',
    });

    expect(fill).toEqual(expect.objectContaining({
      ...defaultFill,
      size: '32.50',
    }));
  });

  describe('get24HourInformation', () => {
    it('Successfully gets 24 hour information with trades', async () => {
      await Promise.all([
        FillTable.create(defaultFill),
        FillTable.create({ // this fill should be ignored
          ...defaultFill,
          liquidity: Liquidity.MAKER,
        }),
      ]);

      // defaultFill.createdAt is the current time the object is created,
      // so which should be in the last 24 before this function is called
      const marketTradeVolumes:
      _.Dictionary<Market24HourTradeVolumes> = await FillTable.get24HourInformation(
        [defaultPerpetualMarket.clobPairId],
      );

      expect(marketTradeVolumes).toEqual({
        [defaultPerpetualMarket.clobPairId]: {
          clobPairId: defaultPerpetualMarket.clobPairId,
          trades24H: '1',
          volume24H: defaultFill.quoteAmount,
        },
      });
    });

    it('Successfully gets 24 hour information with no trades', async () => {
      const marketTradeVolumes:
      _.Dictionary<Market24HourTradeVolumes> = await FillTable.get24HourInformation(
        [defaultPerpetualMarket.clobPairId],
      );

      expect(marketTradeVolumes).toEqual({
        [defaultPerpetualMarket.clobPairId]: {
          clobPairId: defaultPerpetualMarket.clobPairId,
          trades24H: '0',
          volume24H: '0',
        },
      });
    });
  });

  describe('getPnlOfFills/getTotalValueOfOpenPositions', () => {

    beforeEach(async () => {
      await Promise.all([
        OraclePriceTable.create(defaultOraclePrice),
        OraclePriceTable.create(defaultOraclePrice2),
      ]);
    });

    it('Successfully getPnlOfFills/getTotalValueOfOpenPositions', async () => {
      await Promise.all([
        FillTable.create(defaultFill),
        FillTable.create({
          ...defaultFill,
          liquidity: Liquidity.MAKER,
        }),
        FillTable.create({
          ...defaultFill,
          eventId: defaultTendermintEventId2,
          liquidity: Liquidity.TAKER,
          side: OrderSide.SELL,
          size: '2',
        }),
      ]);

      const pnlOfFills: Big = await FillTable.getCostOfFills(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(pnlOfFills).toEqual(Big(-360_000));  // -20000*10 - 20000*10 + 2*20000 = -360000

      const totalValueOfOpenPositions: Big = await FillTable.getTotalValueOfOpenPositions(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(totalValueOfOpenPositions.eq(
        Big(180_000),
      )).toBe(true);  // 18 * 10000 = 180_000
    });

    it('getPnlOfFills/getTotalValueOfOpenPositions( respects height and subaccount id', async () => {
      await BlockTable.create({
        ...defaultBlock,
        blockHeight: '5',
      });
      await Promise.all([
        FillTable.create(defaultFill),
        FillTable.create({
          ...defaultFill,
          subaccountId: defaultSubaccountId2,
          eventId: defaultTendermintEventId2,
        }),
        FillTable.create({
          ...defaultFill,
          liquidity: Liquidity.MAKER,
        }),
        FillTable.create({
          ...defaultFill,
          eventId: defaultTendermintEventId2,
          liquidity: Liquidity.MAKER,
          side: OrderSide.SELL,
          size: '2',
          createdAtHeight: '5',
        }),
      ]);

      const pnlOfFills: Big = await FillTable.getCostOfFills(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(pnlOfFills).toEqual(Big(-400_000));  // -20000*10 - 20000*10 = -400000

      const totalValueOfOpenPositions: Big = await FillTable.getTotalValueOfOpenPositions(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(totalValueOfOpenPositions).toEqual(Big(200_000));  // 20 * 10000 = 200000
    });

    it('returns 0 for missing data', async () => {
      let result: Big = await FillTable.getCostOfFills(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(result).toEqual(Big(0));
      result = await FillTable.getTotalValueOfOpenPositions(
        defaultFill.subaccountId,
        defaultFill.createdAtHeight,
      );
      expect(result).toEqual(Big(0));
    });
  });

  it('Successfully getClobPairs', async () => {
    await Promise.all([
      FillTable.create({
        ...defaultFill,
        createdAtHeight: '1',
      }),
      FillTable.create({
        ...defaultFill,
        liquidity: Liquidity.MAKER,
        size: '2',
        clobPairId: '2',
      }),
    ]);
    let clobPairs: string[] = await
    FillTable.getClobPairs(defaultFill.subaccountId, '1');
    expect(clobPairs).toEqual([defaultPerpetualMarket.clobPairId]);
    clobPairs = await
    FillTable.getClobPairs(defaultFill.subaccountId, '2');
    expect(clobPairs).toEqual([defaultPerpetualMarket.clobPairId, '2']);
  });

  it('Successfully getFeesPaid', async () => {
    await Promise.all([
      FillTable.create({
        ...defaultFill,
        createdAtHeight: '1',
      }),
      FillTable.create({
        ...defaultFill,
        eventId: defaultTendermintEventId2,
        createdAtHeight: '1',
        fee: '-0.5',
      }),
      FillTable.create({
        ...defaultFill,
        liquidity: Liquidity.MAKER,
        size: '2',
        clobPairId: '2',
      }),
    ]);
    let feesPaid: Big = await
    FillTable.getFeesPaid(defaultFill.subaccountId, '1');
    expect(feesPaid).toEqual(Big(0.6));

    feesPaid = await
    FillTable.getFeesPaid(defaultFill.subaccountId, '2');
    expect(feesPaid).toEqual(Big(1.7));
  });

  describe('getOrderedFillsWithFundingIndices', () => {
    beforeEach(async () => {
      await Promise.all([
        OraclePriceTable.create(defaultOraclePrice),
        OraclePriceTable.create(defaultOraclePrice2),
      ]);
      const blockHeights: string[] = ['3', '4', '5', '6', '7'];

      await Promise.all(blockHeights.map((height) => BlockTable.create({
        ...defaultBlock,
        blockHeight: height,
      }),
      ));

      await Promise.all([
        FundingIndexUpdatesTable.create(defaultFundingIndexUpdate),
        FundingIndexUpdatesTable.create({
          ...defaultFundingIndexUpdate,
          effectiveAtHeight: '3',
          fundingIndex: '10100',
        }),
        FundingIndexUpdatesTable.create({
          ...defaultFundingIndexUpdate,
          effectiveAtHeight: '4',
          fundingIndex: '10150',
        }),
        FundingIndexUpdatesTable.create({
          ...defaultFundingIndexUpdate,
          effectiveAtHeight: '5',
          fundingIndex: '10200',
        }),
      ]);
    });

    it('getOrderedFillsWithFundingIndices/getOpenSizeWithFundingIndex returns empty list for no fills',
      async () => {
        const orderedFillsWithFundingIndices: OrderedFillsWithFundingIndices[] = await
        FillTable.getOrderedFillsWithFundingIndices(
          defaultFill.clobPairId,
          defaultFill.subaccountId,
          '7',
        );
        expect(orderedFillsWithFundingIndices.length).toEqual(0);
        const unrealizedFunding: OpenSizeWithFundingIndex[] = await
        FillTable.getOpenSizeWithFundingIndex(
          defaultFill.subaccountId,
          '7',
        );
        expect(unrealizedFunding.length).toEqual(0);
      });

    it('Successfully getOrderedFillsWithFundingIndices/getOpenSizeWithFundingIndex', async () => {
      await Promise.all([
        FillTable.create(defaultFill),
        FillTable.create({
          ...defaultFill,
          createdAtHeight: '3',
          liquidity: Liquidity.MAKER,
          size: '3',
        }),
        FillTable.create({
          ...defaultFill,
          eventId: defaultTendermintEventId2,
          liquidity: Liquidity.TAKER,
          side: OrderSide.SELL,
          size: '4',
          createdAtHeight: '4',
        }),
        FillTable.create({
          ...defaultFill,
          eventId: defaultTendermintEventId2,
          liquidity: Liquidity.MAKER,
          side: OrderSide.SELL,
          size: '5',
          createdAtHeight: '5',
        }),
      ]);

      const unrealizedFunding: OpenSizeWithFundingIndex[] = await
      FillTable.getOpenSizeWithFundingIndex(
        defaultFill.subaccountId,
        '7',
      );
      expect(unrealizedFunding).toEqual(
        expect.objectContaining([{
          clobPairId: '1',
          fundingIndex: '10200',
          fundingIndexHeight: '5',
          lastFillHeight: '5',
          openSize: '4',  // 10+3-4-5 = 4
        }]),
      );

      const orderedFillsWithFundingIndices: OrderedFillsWithFundingIndices[] = await
      FillTable.getOrderedFillsWithFundingIndices(
        defaultFill.clobPairId,
        defaultFill.subaccountId,
        '7',
      );
      expect(orderedFillsWithFundingIndices).toEqual(
        expect.objectContaining([{
          createdAtHeight: '3',
          fundingIndex: '10100',
          id: '898ec79f-8bbe-5bbc-8b3e-bed2a011df71',
          lastFillCreatedAtHeight: '2',
          lastFillFundingIndex: '10050',
          lastFillId: '1d43057d-f67e-534b-bd86-95c16f277d39',
          lastFillSide: 'BUY',
          lastFillSize: '10',
          side: 'BUY',
          size: '3',
          subaccountId: 'df91255d-5e17-5e2e-824a-87ddf3c5214a',
        },
        {
          createdAtHeight: '4',
          fundingIndex: '10150',
          id: '12b2281e-d57f-5c3d-a52b-0144eb304230',
          lastFillCreatedAtHeight: '3',
          lastFillFundingIndex: '10100',
          lastFillId: '898ec79f-8bbe-5bbc-8b3e-bed2a011df71',
          lastFillSide: 'BUY',
          lastFillSize: '3',
          side: 'SELL',
          size: '4',
          subaccountId: 'df91255d-5e17-5e2e-824a-87ddf3c5214a',
        },
        {
          createdAtHeight: '5',
          fundingIndex: '10200',
          id: 'e87bb90f-3595-5f34-9f34-1fae94cec1f0',
          lastFillCreatedAtHeight: '4',
          lastFillFundingIndex: '10150',
          lastFillId: '12b2281e-d57f-5c3d-a52b-0144eb304230',
          lastFillSide: 'SELL',
          lastFillSize: '4',
          side: 'SELL',
          size: '5',
          subaccountId: 'df91255d-5e17-5e2e-824a-87ddf3c5214a',
        },
        ]),
      );
      const paidFunding: Big = FillTable.getSettledFunding(orderedFillsWithFundingIndices);
      // 10 * (10100-10050) + 13 * (10150-10100) + 9 * (10200-10150) = 1600
      expect(paidFunding.eq(new Big(1600))).toBe(true);
    });
  });

  describe('findAll - parentSubaccount', () => {
    it('successfully gets fills for parent and child subaccounts', async () => {
      // Create fills for parent and child subaccounts
      const address = 'parent_address';
      const parentSubaccountNumber = 0;
      const childSubaccountNumbers = [0, 128, 3840]; // Parent and 3 child subaccounts
      // Create subaccounts first
      await Promise.all(childSubaccountNumbers.map((subaccountNum) => SubaccountTable.create({
        address,
        subaccountNumber: subaccountNum,
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: '1',
      })));

      // Add block creation before fills
      await Promise.all([10, 9, 8].map((height) => BlockTable.create({
        ...defaultBlock,
        blockHeight: height.toString(),
      })));

      const tendermintEventIds = [
        defaultTendermintEventId,
        defaultTendermintEventId2,
        defaultTendermintEventId3,
      ];

      // Then create fills as before
      const fillPromises = childSubaccountNumbers.map((subaccountNum, index) => {
        return FillTable.create({
          ...defaultFill,
          subaccountId: SubaccountTable.uuid(address, subaccountNum),
          createdAtHeight: (10 - index).toString(),
          eventId: tendermintEventIds[index],
        });
      });
      const createdFills = await Promise.all(fillPromises);

      // Create additional block for the different address fill
      await BlockTable.create({
        ...defaultBlock,
        blockHeight: '11',
      });

      // Add a fill for a different address that shouldn't be returned
      await SubaccountTable.create({
        address,
        subaccountNumber: 2, // not a child subaccount
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: '1',
      });
      await FillTable.create({
        ...defaultFill,
        subaccountId: SubaccountTable.uuid(address, 2),
        createdAtHeight: '11',
        eventId: defaultTendermintEventId4,
      });

      // Test with high limit
      const { results: fills } = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 100,
        },
        [],
      );

      expect(fills.length).toEqual(3);
      // Verify DESC ordering by createdAtHeight
      expect(fills.map((f: FillFromDatabase) => f.createdAtHeight)).toEqual(['10', '9', '8']);
      // Verify returned fills match created fills
      expect(fills).toEqual(expect.arrayContaining(createdFills));

      // Test with custom limit
      const { results: limitedFills } = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 2,
        },
        [],
      );

      expect(limitedFills.length).toEqual(2);
      expect(limitedFills.map((f: FillFromDatabase) => f.createdAtHeight)).toEqual(['10', '9']);
      // Verify limited fills are a subset of created fills
      expect(createdFills).toEqual(expect.arrayContaining(limitedFills));
    });

    it('successfully paginates fills for parent and child subaccounts', async () => {
      // Create fills for parent and child subaccounts
      const address = 'parent_address';
      const parentSubaccountNumber = 0;
      const childSubaccountNumbers = [0, 128, 3840]; // Parent and 3 child subaccounts
      // Create subaccounts first
      await Promise.all(childSubaccountNumbers.map((subaccountNum) => SubaccountTable.create({
        address,
        subaccountNumber: subaccountNum,
        updatedAt: defaultDateTime.toISO(),
        updatedAtHeight: '1',
      })));

      // Add block creation before fills
      await Promise.all([10, 9, 8].map((height) => BlockTable.create({
        ...defaultBlock,
        blockHeight: height.toString(),
      })));

      const tendermintEventIds = [
        defaultTendermintEventId,
        defaultTendermintEventId2,
        defaultTendermintEventId3,
      ];

      // Then create fills
      await Promise.all(childSubaccountNumbers.map((subaccountNum, index) => {
        return FillTable.create({
          ...defaultFill,
          subaccountId: SubaccountTable.uuid(address, subaccountNum),
          createdAtHeight: (10 - index).toString(),
          eventId: tendermintEventIds[index],
        });
      }));

      // Test with pagination
      const responsePageOne = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 1,
          page: 1,
        },
        [],
      );

      expect(responsePageOne.results.length).toEqual(1);
      expect(responsePageOne.results[0].createdAtHeight).toEqual('10');
      expect(responsePageOne.offset).toEqual(0);
      expect(responsePageOne.total).toEqual(3);

      const responsePageTwo = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 1,
          page: 2,
        },
        [],
      );

      expect(responsePageTwo.results.length).toEqual(1);
      expect(responsePageTwo.results[0].createdAtHeight).toEqual('9');
      expect(responsePageTwo.offset).toEqual(1);
      expect(responsePageTwo.total).toEqual(3);

      const responsePageThree = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 1,
          page: 3,
        },
        [],
      );

      expect(responsePageThree.results.length).toEqual(1);
      expect(responsePageThree.results[0].createdAtHeight).toEqual('8');
      expect(responsePageThree.offset).toEqual(2);
      expect(responsePageThree.total).toEqual(3);

      // Test getting all results in one page
      const responseAllPages = await FillTable.findAll(
        {
          parentSubaccount: {
            address,
            subaccountNumber: parentSubaccountNumber,
          },
          limit: 3,
          page: 1,
        },
        [],
      );

      expect(responseAllPages.results.length).toEqual(3);
      expect(responseAllPages.results.map((f) => f.createdAtHeight)).toEqual(['10', '9', '8']);
      expect(responseAllPages.offset).toEqual(0);
      expect(responseAllPages.total).toEqual(3);
    });

    it('returns empty array when no fills exist', async () => {
      const { results: fills } = await FillTable.findAll(
        {
          parentSubaccount: {
            address: 'nonexistent_address',
            subaccountNumber: 0,
          },
          limit: 100,
        },
        [],
      );

      expect(fills).toEqual([]);
    });
  });
});
