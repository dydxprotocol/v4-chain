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
      };

      expect(response.body.positions).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedAssetPosition,
          }),
        ]),
      );
    });

    it('Get /assetPositions gets short asset and perpetual positions', async () => {
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
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/assetPositions?address=${testConstants.defaultAddress}` +
          `&subaccountNumber=${testConstants.defaultSubaccount.subaccountNumber}`,
      });

      expect(response.body.positions).toEqual(
        [{
          symbol: testConstants.defaultAsset.symbol,
          size: '9500',
          side: PositionSide.LONG,
          assetId: testConstants.defaultAssetPosition.assetId,
        }],
      );
    });
  });
});
