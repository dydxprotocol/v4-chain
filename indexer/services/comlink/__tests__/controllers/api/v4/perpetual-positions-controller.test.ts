import {
  BlockTable,
  dbHelpers,
  FundingIndexUpdatesTable,
  perpetualMarketRefresher,
  PerpetualPositionStatus,
  PerpetualPositionTable,
  PositionSide,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { PerpetualPositionResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, getQueryString, sendRequest } from '../../../helpers/helpers';

describe('perpetual-positions-controller#V4', () => {
  const latestHeight: string = '3';

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await BlockTable.create({
      ...testConstants.defaultBlock,
      blockHeight: latestHeight,
    });
    await Promise.all([
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        fundingIndex: '10000',
        effectiveAtHeight: testConstants.createdHeight,
      }),
      FundingIndexUpdatesTable.create({
        ...testConstants.defaultFundingIndexUpdate,
        eventId: testConstants.defaultTendermintEventId2,
        effectiveAtHeight: latestHeight,
      }),
    ]);
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  describe('GET', () => {
    const defaultSubaccountNumber: number = testConstants.defaultSubaccount.subaccountNumber;
    const defaultAddress: string = testConstants.defaultSubaccount.address;

    it('Get /perpetualPositions gets long position', async () => {
      await PerpetualPositionTable.create(testConstants.defaultPerpetualPosition);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expected: PerpetualPositionResponseObject = {
        market: testConstants.defaultPerpetualMarket.ticker,
        side: testConstants.defaultPerpetualPosition.side,
        status: testConstants.defaultPerpetualPosition.status,
        size: testConstants.defaultPerpetualPosition.size,
        maxSize: testConstants.defaultPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.defaultPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.defaultPerpetualPosition.sumOpen!,
        sumClose: testConstants.defaultPerpetualPosition.sumClose!,
        // For the calculation of the net funding (long position):
        // settled funding on position = 200_000, size = 10, latest funding index = 10050
        // last updated funding index = 10000
        // total funding = 200_000 + (10 * (10000 - 10050)) = 199_500
        netFunding: getFixedRepresentation('199500'),
        // sumClose=0, so realized Pnl is the same as the net funding of the position.
        // Unsettled funding is funding payments that already "happened" but not reflected
        // in the subaccount's balance yet, so it's considered a part of realizedPnl.
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (long position):
        // index price = 15_000, entry price = 20_000, size = 10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = 10 * (15_000 - 20_000)
        unrealizedPnl: getFixedRepresentation('-50000'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it('Get /perpetualPositions gets short position', async () => {
      await PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        side: PositionSide.SHORT,
        size: '-10',
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expected: PerpetualPositionResponseObject = {
        market: testConstants.defaultPerpetualMarket.ticker,
        side: PositionSide.SHORT,
        status: testConstants.defaultPerpetualPosition.status,
        size: '-10',
        maxSize: testConstants.defaultPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.defaultPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.defaultPerpetualPosition.sumOpen!,
        sumClose: testConstants.defaultPerpetualPosition.sumClose!,
        // For the calculation of the net funding (short position):
        // settled funding on position = 200_000, size = -10, latest funding index = 10050
        // last updated funding index = 10000
        // total funding = 200_000 + (-10 * (10000 - 10050)) = 200_500
        netFunding: getFixedRepresentation('200500'),
        // sumClose=0, so realized Pnl is the same as the net funding of the position.
        // Unsettled funding is funding payments that already "happened" but not reflected
        // in the subaccount's balance yet, so it's considered a part of realizedPnl.
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (short position):
        // index price = 15_000, entry price = 20_000, size = -10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = -10 * (15_000 - 20_000)
        unrealizedPnl: getFixedRepresentation('50000'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it('Get /perpetualPositions gets CLOSED position without adjusting funding', async () => {
      await PerpetualPositionTable.create({
        ...testConstants.defaultPerpetualPosition,
        status: PerpetualPositionStatus.CLOSED,
        side: PositionSide.SHORT,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expected: PerpetualPositionResponseObject = {
        market: testConstants.defaultPerpetualMarket.ticker,
        side: PositionSide.SHORT,
        status: PerpetualPositionStatus.CLOSED,
        size: testConstants.defaultPerpetualPosition.size,
        maxSize: testConstants.defaultPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.defaultPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.defaultPerpetualPosition.sumOpen!,
        sumClose: testConstants.defaultPerpetualPosition.sumClose!,
        // CLOSED position should not have funding adjusted
        netFunding: getFixedRepresentation(
          testConstants.defaultPerpetualPosition.settledFunding,
        ),
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (short position):
        // index price = 15_000, entry price = 20_000, size = 10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = 10 * (15_000 - 20_000)
        unrealizedPnl: getFixedRepresentation('-50000'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
        ]),
      );
    });

    it.each([
      [
        'invalid status',
        { address: defaultAddress, subaccountNumber: defaultSubaccountNumber, status: 'INVALID' },
        'status',
        'status must be a valid Position Status (OPEN, etc)',
      ],
      [
        'multiple invalid status',
        {
          address: defaultAddress,
          subaccountNumber: defaultSubaccountNumber,
          status: 'INVALID,INVALID',
        },
        'status',
        'status must be a valid Position Status (OPEN, etc)',
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        subaccountNumber?: number,
        status?: string,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions?${getQueryString(queryParams)}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });

    it('Get /perpetualPositions/parentSubaccountNumber gets long/short positions across subaccounts', async () => {
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create({
          ...testConstants.isolatedPerpetualPosition,
          side: PositionSide.SHORT,
          size: '-10',
        }),
      ]);
      await Promise.all([
        FundingIndexUpdatesTable.create({
          ...testConstants.isolatedMarketFundingIndexUpdate,
          fundingIndex: '10000',
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.isolatedMarketFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: latestHeight,
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expected: PerpetualPositionResponseObject = {
        market: testConstants.defaultPerpetualMarket.ticker,
        side: testConstants.defaultPerpetualPosition.side,
        status: testConstants.defaultPerpetualPosition.status,
        size: testConstants.defaultPerpetualPosition.size,
        maxSize: testConstants.defaultPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.defaultPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.defaultPerpetualPosition.sumOpen!,
        sumClose: testConstants.defaultPerpetualPosition.sumClose!,
        // For the calculation of the net funding (long position):
        // settled funding on position = 200_000, size = 10, latest funding index = 10050
        // last updated funding index = 10000
        // total funding = 200_000 + (10 * (10000 - 10050)) = 199_500
        netFunding: getFixedRepresentation('199500'),
        // sumClose=0, so realized Pnl is the same as the net funding of the position.
        // Unsettled funding is funding payments that already "happened" but not reflected
        // in the subaccount's balance yet, so it's considered a part of realizedPnl.
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (long position):
        // index price = 15_000, entry price = 20_000, size = 10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = 10 * (15_000 - 20_000)
        unrealizedPnl: getFixedRepresentation('-50000'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };
      // object for expected 2 which holds an isolated position in an isolated perpetual
      // in the isolated subaccount
      const expected2: PerpetualPositionResponseObject = {
        market: testConstants.isolatedPerpetualMarket.ticker,
        side: PositionSide.SHORT,
        status: testConstants.isolatedPerpetualPosition.status,
        size: '-10',
        maxSize: testConstants.isolatedPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.isolatedPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.isolatedPerpetualPosition.sumOpen!,
        sumClose: testConstants.isolatedPerpetualPosition.sumClose!,
        // For the calculation of the net funding (short position):
        // settled funding on position = 200_000, size = -10, latest funding index = 10200
        // last updated funding index = 10000
        // total funding = 200_000 + (-10 * (10000 - 10200)) = 202_000
        netFunding: getFixedRepresentation('202000'),
        // sumClose=0, so realized Pnl is the same as the net funding of the position.
        // Unsettled funding is funding payments that already "happened" but not reflected
        // in the subaccount's balance yet, so it's considered a part of realizedPnl.
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (short position):
        // index price = 1, entry price = 1.5, size = -10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = -10 * (1-1.5)
        unrealizedPnl: getFixedRepresentation('5'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
          expect.objectContaining({
            ...expected2,
          }),
        ]),
      );
    });

    it('Get /perpetualPositions/parentSubaccountNumber gets CLOSED position without adjusting funding', async () => {
      await Promise.all([
        PerpetualPositionTable.create({
          ...testConstants.defaultPerpetualPosition,
          status: PerpetualPositionStatus.CLOSED,
        }),
        PerpetualPositionTable.create({
          ...testConstants.isolatedPerpetualPosition,
          side: PositionSide.SHORT,
          size: '-10',
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expected: PerpetualPositionResponseObject = {
        market: testConstants.defaultPerpetualMarket.ticker,
        side: testConstants.defaultPerpetualPosition.side,
        status: PerpetualPositionStatus.CLOSED,
        size: testConstants.defaultPerpetualPosition.size,
        maxSize: testConstants.defaultPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.defaultPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.defaultPerpetualPosition.sumOpen!,
        sumClose: testConstants.defaultPerpetualPosition.sumClose!,
        // CLOSED position should not have funding adjusted
        netFunding: getFixedRepresentation(
          testConstants.defaultPerpetualPosition.settledFunding,
        ),
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (short position):
        // index price = 15_000, entry price = 20_000, size = 10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = 10 * (15_000 - 20_000)
        unrealizedPnl: getFixedRepresentation('-50000'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };
      const expected2: PerpetualPositionResponseObject = {
        market: testConstants.isolatedPerpetualMarket.ticker,
        side: PositionSide.SHORT,
        status: testConstants.isolatedPerpetualPosition.status,
        size: '-10',
        maxSize: testConstants.isolatedPerpetualPosition.maxSize,
        entryPrice: getFixedRepresentation(testConstants.isolatedPerpetualPosition.entryPrice!),
        exitPrice: null,
        sumOpen: testConstants.isolatedPerpetualPosition.sumOpen!,
        sumClose: testConstants.isolatedPerpetualPosition.sumClose!,
        // CLOSED position should not have funding adjusted
        netFunding: getFixedRepresentation(
          testConstants.isolatedPerpetualPosition.settledFunding,
        ),
        realizedPnl: getFixedRepresentation('100'),
        // For the calculation of the unrealized pnl (short position):
        // index price = 1, entry price = 1.5, size = -10
        // unrealizedPnl = size * (index price - entry price)
        // unrealizedPnl = -10 * (1-1.5)
        unrealizedPnl: getFixedRepresentation('5'),
        createdAt: testConstants.createdDateTime.toISO(),
        closedAt: null,
        createdAtHeight: testConstants.createdHeight,
        subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expected,
          }),
          expect.objectContaining({
            ...expected2,
          }),
        ]),
      );
    });

    it.each([
      [
        'invalid status',
        {
          address: defaultAddress,
          parentSubaccountNumber: defaultSubaccountNumber,
          status: 'INVALID',
        },
        'status',
        'status must be a valid Position Status (OPEN, etc)',
      ],
      [
        'multiple invalid status',
        {
          address: defaultAddress,
          parentSubaccountNumber: defaultSubaccountNumber,
          status: 'INVALID,INVALID',
        },
        'status',
        'status must be a valid Position Status (OPEN, etc)',
      ],
    ])('Returns 400 when validation fails: %s', async (
      _reason: string,
      queryParams: {
        address?: string,
        subaccountNumber?: number,
        status?: string,
      },
      fieldWithError: string,
      expectedErrorMsg: string,
    ) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/perpetualPositions/parentSubaccountNumber?${getQueryString(queryParams)}`,
        expectedStatus: 400,
      });

      expect(response.body).toEqual(expect.objectContaining({
        errors: expect.arrayContaining([
          expect.objectContaining({
            param: fieldWithError,
            msg: expectedErrorMsg,
          }),
        ]),
      }));
    });
  });
});
