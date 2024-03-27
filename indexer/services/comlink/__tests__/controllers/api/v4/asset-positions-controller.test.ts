import {
  AssetPositionTable,
  BlockTable,
  dbHelpers,
  FundingIndexUpdatesTable,
  PerpetualPositionTable,
  PositionSide,
  testConstants,
  testMocks,
} from '@dydxprotocol-indexer/postgres';
import { AssetPositionResponseObject, RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';

describe('asset-positions-controller#V4', () => {
  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET', () => {

    afterEach(async () => {
      await dbHelpers.clearData();
    });

    it('Get /assetPositions gets long asset positions', async () => {
      const size: string = '192.12421';
      await testMocks.seedData();
      await AssetPositionTable.upsert({
        ...testConstants.defaultAssetPosition,
        size,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions?address=${testConstants.defaultAddress}` +
            `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.LONG,
        size,
        assetId: testConstants.defaultAssetPosition.assetId,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedAssetPosition,
          }),
        ]),
      );
    });

    it('Get /assetPositions gets short asset positions', async () => {
      await testMocks.seedData();
      await AssetPositionTable.upsert({
        ...testConstants.defaultAssetPosition,
        isLong: false,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions?address=${testConstants.defaultAddress}` +
            `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      const expectedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.SHORT,
        size: testConstants.defaultAssetPosition.size,
        assetId: testConstants.defaultAssetPosition.assetId,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedAssetPosition,
          }),
        ]),
      );
    });

    it('Get /assetPositions does not get asset positions with 0 size', async () => {
      await testMocks.seedData();

      await Promise.all([
        await AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        await AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          subaccountId: testConstants.defaultSubaccountId,
          size: '0',
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions?address=${testConstants.defaultAddress}` +
            `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.positions).toEqual(
        [{
          symbol: testConstants.defaultAsset.symbol,
          size: testConstants.defaultAssetPosition.size,
          side: PositionSide.LONG,
          assetId: testConstants.defaultAssetPosition.assetId,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        }],
      );
    });

    it('Get /assetPositions gets USDC asset position adjusted by unsettled funding', async () => {
      await testMocks.seedData();
      await BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '3',
      });
      await Promise.all([
        PerpetualPositionTable.create(
          testConstants.defaultPerpetualPosition,
        ),
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition2,
          subaccountId: testConstants.defaultSubaccountId,
          size: '0',
        }),
        // Funding index at height 0 is 10000
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: '10000',
          effectiveAtHeight: testConstants.createdHeight,
        }),
        // Funding index at height 3 is 10050
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: '3',
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions?address=${testConstants.defaultAddress}` +
            `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.positions).toEqual(
        [{
          symbol: testConstants.defaultAsset.symbol,
          // funding index difference = 10050 (height 3) - 10000 (height 0) = 50
          // size = 10000 (initial size) - 50 (funding index diff) * 10(position size)
          size: '9500',
          side: PositionSide.LONG,
          assetId: testConstants.defaultAssetPosition.assetId,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        }],
      );
    });

    it('Get /assetPositions/parentSubaccountNumber gets long and short asset positions across subaccounts', async () => {
      await testMocks.seedData();
      await Promise.all([
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.isolatedSubaccountAssetPosition,
          isLong: false,
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expectedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.LONG,
        size: testConstants.defaultAssetPosition.size,
        assetId: testConstants.defaultAssetPosition.assetId,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };
      const expectedIsolatedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.SHORT,
        size: testConstants.isolatedSubaccountAssetPosition.size,
        assetId: testConstants.isolatedSubaccountAssetPosition.assetId,
        subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedAssetPosition,
          }),
          expect.objectContaining({
            ...expectedIsolatedAssetPosition,
          }),
        ]),
      );
    });

    it('Get /assetPositions/parentSubaccountNumber does not get asset positions with 0 size', async () => {
      await testMocks.seedData();

      await Promise.all([
        await AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        await AssetPositionTable.upsert({
          ...testConstants.isolatedSubaccountAssetPosition,
          size: '0',
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      expect(response.body.positions).toEqual(
        [{
          symbol: testConstants.defaultAsset.symbol,
          size: testConstants.defaultAssetPosition.size,
          side: PositionSide.LONG,
          assetId: testConstants.defaultAssetPosition.assetId,
          subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
        }],
      );
    });

    it('Get /assetPositions/parentSubaccountNumber gets USDC asset positions adjusted by unsettled funding', async () => {
      await testMocks.seedData();
      await BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight: '3',
      });
      await Promise.all([
        PerpetualPositionTable.create(testConstants.defaultPerpetualPosition),
        PerpetualPositionTable.create(testConstants.isolatedPerpetualPosition),
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert(testConstants.isolatedSubaccountAssetPosition),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: '10000',
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: '3',
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.isolatedMarketFundingIndexUpdate,
          fundingIndex: '10000',
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.isolatedMarketFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: '3',
        }),
      ]);

      const parentSubaccountNumber: number = 0;
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions/parentSubaccountNumber?address=${testConstants.defaultAddress}` +
            `&parentSubaccountNumber=${parentSubaccountNumber}`,
      });

      const expectedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.LONG,
        // funding index difference = 10050 (height 3) - 10000 (height 0) = 50
        // size = 10000 (initial size) - 50 (funding index diff) * 10(position size)
        size: '9500',
        assetId: testConstants.defaultAssetPosition.assetId,
        subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
      };
      const expectedIsolatedAssetPosition: AssetPositionResponseObject = {
        symbol: testConstants.defaultAsset.symbol,
        side: PositionSide.LONG,
        // funding index difference = 10200 (height 3) - 10000 (height 0) = 200
        // size = 5000 (initial size) - 200 (funding index diff) * 10(position size)
        size: '3000',
        assetId: testConstants.isolatedSubaccountAssetPosition.assetId,
        subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedAssetPosition,
          }),
          expect.objectContaining({
            ...expectedIsolatedAssetPosition,
          }),
        ]),
      );
    });
  });
});
