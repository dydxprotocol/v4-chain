import {
  dbHelpers,
  testMocks,
  testConstants,
  perpetualMarketRefresher,
  PerpetualPositionTable,
  AssetPositionTable,
  PositionSide,
  FundingIndexUpdatesTable,
  BlockTable,
  liquidityTierRefresher,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import { stats } from '@dydxprotocol-indexer/base';
import { defaultAddress } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

describe('addresses-controller#V4', () => {
  const latestHeight: string = '3';
  const initialFundingIndex: string = '10000';

  beforeAll(async () => {
    await dbHelpers.migrate();
    jest.spyOn(stats, 'increment');
    jest.spyOn(stats, 'timing');
  });

  beforeEach(async () => {
    await testMocks.seedData();
    await liquidityTierRefresher.updateLiquidityTiers();
    await BlockTable.create({
      ...testConstants.defaultBlock,
      blockHeight: latestHeight,
    });
    await perpetualMarketRefresher.updatePerpetualMarkets();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  afterEach(async () => {
    await dbHelpers.clearData();
  });

  const invalidAddress: string = 'invalidAddress';
  describe('/addresses/:address/subaccountNumber/:subaccountNumber', () => {
    it('Get / gets subaccount', async () => {
      await PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      );
      await Promise.all([
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          subaccountId: testConstants.defaultSubaccountId,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: initialFundingIndex,
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: latestHeight,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${testConstants.defaultAddress}/subaccountNumber/` +
        `${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body).toEqual({
        subaccount: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          equity: getFixedRepresentation(159500),
          freeCollateral: getFixedRepresentation(152000),
          marginEnabled: true,
          assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
          openPerpetualPositions: {
            [testConstants.defaultPerpetualMarket.ticker]: {
              market: testConstants.defaultPerpetualMarket.ticker,
              size: testConstants.defaultPerpetualPosition.size,
              side: testConstants.defaultPerpetualPosition.side,
              entryPrice: getFixedRepresentation(
                testConstants.defaultPerpetualPosition.entryPrice!,
              ),
              maxSize: testConstants.defaultPerpetualPosition.maxSize,
              // 200000 + 10*(10000-10050)=199500
              netFunding: getFixedRepresentation('199500'),
              // sumClose=0, so realized Pnl is the same as the net funding of the position.
              // Unsettled funding is funding payments that already "happened" but not reflected
              // in the subaccount's balance yet, so it's considered a part of realizedPnl.
              realizedPnl: getFixedRepresentation('199500'),
              // size * (index-entry) = 10*(15000-20000) = -50000
              unrealizedPnl: getFixedRepresentation(-50000),
              status: testConstants.defaultPerpetualPosition.status,
              sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
              sumClose: testConstants.defaultPerpetualPosition.sumClose,
              createdAt: testConstants.defaultPerpetualPosition.createdAt,
              createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
              exitPrice: null,
              closedAt: null,
              perpYieldIndex: testConstants.defaultPerpetualPosition.perpYieldIndex!,
            },
          },
          assetPositions: {
            [testConstants.defaultAsset.symbol]: {
              symbol: testConstants.defaultAsset.symbol,
              size: '9500',
              side: PositionSide.LONG,
              assetId: testConstants.defaultAssetPosition.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
            [testConstants.defaultAsset2.symbol]: {
              symbol: testConstants.defaultAsset2.symbol,
              size: testConstants.defaultAssetPosition2.size,
              side: PositionSide.SHORT,
              assetId: testConstants.defaultAssetPosition2.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
          },
        },
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1,
        {
          path: '/:address/subaccountNumber/:subaccountNumber',
          method: 'GET',
        });
    });

    it('Asset positions with 0 size are not returned', async () => {
      await Promise.all([
        AssetPositionTable.upsert(
          testConstants.defaultAssetPosition,
        ),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          size: '0',
        },
        ),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${testConstants.defaultAddress}/subaccountNumber/` +
          `${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body).toEqual({
        subaccount: {
          address: testConstants.defaultAddress,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
          equity: getFixedRepresentation(10000),
          freeCollateral: getFixedRepresentation(10000),
          marginEnabled: true,
          assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
          openPerpetualPositions: {},
          assetPositions: {
            [testConstants.defaultAsset.symbol]: {
              symbol: testConstants.defaultAsset.symbol,
              size: testConstants.defaultAssetPosition.size,
              side: PositionSide.LONG,
              assetId: testConstants.defaultAssetPosition.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
          },
        },
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1,
        {
          path: '/:address/subaccountNumber/:subaccountNumber',
          method: 'GET',
        });
    });

    it('Get / with non-existent address and subaccount number returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${invalidAddress}/subaccountNumber/` +
        `${testConstants.defaultSubaccount.subaccountNumber}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: `No subaccount found with address ${invalidAddress} and ` +
            `subaccountNumber ${testConstants.defaultSubaccount.subaccountNumber}`,
          },
        ],
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.404', 1,
        {
          path: '/:address/subaccountNumber/:subaccountNumber',
          method: 'GET',
        });
    });
  });

  describe('/addresses/:address', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get / gets all subaccounts', async () => {
      await PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      );

      await Promise.all([
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          subaccountId: testConstants.defaultSubaccountId,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: initialFundingIndex,
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: latestHeight,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${testConstants.defaultAddress}`,
      });

      expect(response.body).toEqual({
        subaccounts: [
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            equity: getFixedRepresentation(159500),
            freeCollateral: getFixedRepresentation(152000),
            assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
            marginEnabled: true,
            openPerpetualPositions: {
              [testConstants.defaultPerpetualMarket.ticker]: {
                market: testConstants.defaultPerpetualMarket.ticker,
                size: testConstants.defaultPerpetualPosition.size,
                side: testConstants.defaultPerpetualPosition.side,
                entryPrice: getFixedRepresentation(
                  testConstants.defaultPerpetualPosition.entryPrice!,
                ),
                maxSize: testConstants.defaultPerpetualPosition.maxSize,
                // 200000 + 10*(10000-10050)=199500
                netFunding: getFixedRepresentation('199500'),
                // sumClose=0, so realized Pnl is the same as the net funding of the position.
                // Unsettled funding is funding payments that already "happened" but not reflected
                // in the subaccount's balance yet, so it's considered a part of realizedPnl.
                realizedPnl: getFixedRepresentation('199500'),
                // size * (index-entry) = 10*(15000-20000) = -50000
                unrealizedPnl: getFixedRepresentation(-50000),
                status: testConstants.defaultPerpetualPosition.status,
                sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
                sumClose: testConstants.defaultPerpetualPosition.sumClose,
                createdAt: testConstants.defaultPerpetualPosition.createdAt,
                createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
                exitPrice: null,
                closedAt: null,
                perpYieldIndex: testConstants.defaultPerpetualPosition.perpYieldIndex!,
              },
            },
            assetPositions: {
              [testConstants.defaultAsset.symbol]: {
                symbol: testConstants.defaultAsset.symbol,
                size: '9500',
                side: PositionSide.LONG,
                assetId: testConstants.defaultAssetPosition.assetId,
                subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
              },
              [testConstants.defaultAsset2.symbol]: {
                symbol: testConstants.defaultAsset2.symbol,
                size: testConstants.defaultAssetPosition2.size,
                side: PositionSide.SHORT,
                assetId: testConstants.defaultAssetPosition2.assetId,
                subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
              },
            },
          },
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.defaultSubaccount2.subaccountNumber,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            assetYieldIndex: testConstants.defaultSubaccount2.assetYieldIndex,
            openPerpetualPositions: {},
            assetPositions: {},
          },
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
            assetYieldIndex: testConstants.isolatedSubaccount.assetYieldIndex,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            openPerpetualPositions: {},
            assetPositions: {},
          },
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.isolatedSubaccount2.subaccountNumber,
            assetYieldIndex: testConstants.isolatedSubaccount2.assetYieldIndex,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            openPerpetualPositions: {},
            assetPositions: {},
          },
        ],
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1,
        {
          path: '/:address',
          method: 'GET',
        });
    });

    it('Get / with non-existent address returns 404', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${invalidAddress}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: `No subaccounts found for address ${invalidAddress}`,
          },
        ],
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.404', 1,
        {
          path: '/:address',
          method: 'GET',
        });
    });
  });

  describe('/addresses/:address/parentSubaccountNumber/:parentSubaccountNumber', () => {
    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /:address/parentSubaccountNumber/ gets all subaccounts for the provided parent', async () => {
      await PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      );

      await Promise.all([
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          subaccountId: testConstants.defaultSubaccountId,
        }),
        AssetPositionTable.upsert(testConstants.isolatedSubaccountAssetPosition),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: initialFundingIndex,
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: latestHeight,
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${testConstants.defaultAddress}/parentSubaccountNumber/${parentSubaccountNumber}`,
      });

      expect(response.body).toEqual({
        subaccount: {
          address: testConstants.defaultAddress,
          parentSubaccountNumber,
          equity: getFixedRepresentation(164500),
          freeCollateral: getFixedRepresentation(157000),
          childSubaccounts: [
            {
              address: testConstants.defaultAddress,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
              equity: getFixedRepresentation(159500),
              freeCollateral: getFixedRepresentation(152000),
              marginEnabled: true,
              assetYieldIndex: testConstants.defaultSubaccount.assetYieldIndex,
              openPerpetualPositions: {
                [testConstants.defaultPerpetualMarket.ticker]: {
                  market: testConstants.defaultPerpetualMarket.ticker,
                  size: testConstants.defaultPerpetualPosition.size,
                  side: testConstants.defaultPerpetualPosition.side,
                  entryPrice: getFixedRepresentation(
                    testConstants.defaultPerpetualPosition.entryPrice!,
                  ),
                  maxSize: testConstants.defaultPerpetualPosition.maxSize,
                  // 200000 + 10*(10000-10050)=199500
                  netFunding: getFixedRepresentation('199500'),
                  // sumClose=0, so realized Pnl is the same as the net funding of the position.
                  // Unsettled funding is funding payments that already "happened" but not reflected
                  // in the subaccount's balance yet, so it's considered a part of realizedPnl.
                  realizedPnl: getFixedRepresentation('199500'),
                  // size * (index-entry) = 10*(15000-20000) = -50000
                  unrealizedPnl: getFixedRepresentation(-50000),
                  status: testConstants.defaultPerpetualPosition.status,
                  sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
                  sumClose: testConstants.defaultPerpetualPosition.sumClose,
                  createdAt: testConstants.defaultPerpetualPosition.createdAt,
                  createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
                  exitPrice: null,
                  closedAt: null,
                  perpYieldIndex: testConstants.defaultPerpetualPosition.perpYieldIndex!,
                },
              },
              assetPositions: {
                [testConstants.defaultAsset.symbol]: {
                  symbol: testConstants.defaultAsset.symbol,
                  size: '9500',
                  side: PositionSide.LONG,
                  assetId: testConstants.defaultAssetPosition.assetId,
                  subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
                },
                [testConstants.defaultAsset2.symbol]: {
                  symbol: testConstants.defaultAsset2.symbol,
                  size: testConstants.defaultAssetPosition2.size,
                  side: PositionSide.SHORT,
                  assetId: testConstants.defaultAssetPosition2.assetId,
                  subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
                },
              },
            },
            {
              address: testConstants.defaultAddress,
              subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
              assetYieldIndex: testConstants.isolatedSubaccount.assetYieldIndex,
              equity: getFixedRepresentation(5000),
              freeCollateral: getFixedRepresentation(5000),
              marginEnabled: true,
              openPerpetualPositions: {},
              assetPositions: {
                [testConstants.defaultAsset.symbol]: {
                  symbol: testConstants.defaultAsset.symbol,
                  size: testConstants.isolatedSubaccountAssetPosition.size,
                  side: PositionSide.LONG,
                  assetId: testConstants.isolatedSubaccountAssetPosition.assetId,
                  subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
                },
              },
            },
            {
              address: testConstants.defaultAddress,
              subaccountNumber: testConstants.isolatedSubaccount2.subaccountNumber,
              assetYieldIndex: testConstants.isolatedSubaccount2.assetYieldIndex,
              equity: getFixedRepresentation(0),
              freeCollateral: getFixedRepresentation(0),
              marginEnabled: true,
              openPerpetualPositions: {},
              assetPositions: {},
            },
          ],
        },
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1,
        {
          path: '/:address/parentSubaccountNumber/:parentSubaccountNumber',
          method: 'GET',
        });
    });
  });

  it('Get /:address/parentSubaccountNumber/ with non-existent address returns 404', async () => {
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/addresses/${invalidAddress}/parentSubaccountNumber/` +
          `${testConstants.defaultSubaccount.subaccountNumber}`,
      expectedStatus: 404,
    });

    expect(response.body).toEqual({
      errors: [
        {
          msg: `No subaccounts found for address ${invalidAddress} and ` +
              `parentSubaccountNumber ${testConstants.defaultSubaccount.subaccountNumber}`,
        },
      ],
    });
    expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.404', 1,
      {
        path: '/:address/parentSubaccountNumber/:parentSubaccountNumber',
        method: 'GET',
      });
  });

  it('Get /:address/parentSubaccountNumber/ with invalid parentSubaccount number returns 400', async () => {
    const parentSubaccountNumber: number = 128;
    const response: request.Response = await sendRequest({
      type: RequestMethod.GET,
      path: `/v4/addresses/${defaultAddress}/parentSubaccountNumber/${parentSubaccountNumber}`,
      expectedStatus: 400,
    });

    expect(response.body).toEqual({
      errors: [
        {
          location: 'params',
          msg: 'parentSubaccountNumber must be a non-negative integer less than 128',
          param: 'parentSubaccountNumber',
          value: '128',
        },
      ],
    });
  });
});
