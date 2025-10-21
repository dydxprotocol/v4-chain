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
  SubaccountTable,
  FirebaseNotificationTokenTable,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import { stats } from '@dydxprotocol-indexer/base';
import config from '../../../../src/config';
import * as complianceUtils from '../../../../src/helpers/compliance/compliance-utils';
import { Secp256k1 } from '@cosmjs/crypto';
import { toBech32 } from '@cosmjs/encoding';
import { DateTime } from 'luxon';
import { verifyADR36Amino } from '@keplr-wallet/cosmos';
import { defaultAddress3 } from '@dydxprotocol-indexer/postgres/build/__tests__/helpers/constants';

jest.mock('@cosmjs/crypto', () => ({
  ...jest.requireActual('@cosmjs/crypto'),
  Secp256k1: {
    verifySignature: jest.fn(),
  },
  ExtendedSecp256k1Signature: {
    fromFixedLength: jest.fn(),
  },
}));

jest.mock('@cosmjs/encoding', () => ({
  toBech32: jest.fn(),
}));

jest.mock('@keplr-wallet/cosmos', () => ({
  ...jest.requireActual('@keplr-wallet/cosmos'),
  verifyADR36Amino: jest.fn(),
}));

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
    jest.clearAllMocks();
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
          updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
          latestProcessedBlockHeight: latestHeight,
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
              realizedPnl: getFixedRepresentation('100'),
              // size * (index-entry) = 10*(15000-20000) = -50000
              unrealizedPnl: getFixedRepresentation(-50000),
              status: testConstants.defaultPerpetualPosition.status,
              sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
              sumClose: testConstants.defaultPerpetualPosition.sumClose,
              createdAt: testConstants.defaultPerpetualPosition.createdAt,
              createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
              exitPrice: null,
              closedAt: null,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
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
          updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
          latestProcessedBlockHeight: latestHeight,
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
        path: `/v4/addresses/${defaultAddress3}/subaccountNumber/` +
        `${testConstants.defaultSubaccount.subaccountNumber}`,
        expectedStatus: 404,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: `No subaccount found with address ${defaultAddress3} and ` +
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
            marginEnabled: true,
            updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
            latestProcessedBlockHeight: latestHeight,
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
                realizedPnl: getFixedRepresentation('100'),
                // size * (index-entry) = 10*(15000-20000) = -50000
                unrealizedPnl: getFixedRepresentation(-50000),
                status: testConstants.defaultPerpetualPosition.status,
                sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
                sumClose: testConstants.defaultPerpetualPosition.sumClose,
                createdAt: testConstants.defaultPerpetualPosition.createdAt,
                createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
                exitPrice: null,
                closedAt: null,
                subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
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
            updatedAtHeight: testConstants.defaultSubaccount2.updatedAtHeight,
            latestProcessedBlockHeight: latestHeight,
            openPerpetualPositions: {},
            assetPositions: {},
          },
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.isolatedSubaccount.subaccountNumber,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            updatedAtHeight: testConstants.isolatedSubaccount.updatedAtHeight,
            latestProcessedBlockHeight: latestHeight,
            openPerpetualPositions: {},
            assetPositions: {},
          },
          {
            address: testConstants.defaultAddress,
            subaccountNumber: testConstants.isolatedSubaccount2.subaccountNumber,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            updatedAtHeight: testConstants.isolatedSubaccount2.updatedAtHeight,
            latestProcessedBlockHeight: latestHeight,
            openPerpetualPositions: {},
            assetPositions: {},
          },
        ],
        totalTradingRewards: testConstants.defaultWallet.totalTradingRewards,
      });
      expect(stats.increment).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1,
        {
          path: '/:address',
          method: 'GET',
        });
    });

    it('returns 0 for totalTradingRewards if no wallet exists', async () => {
      await PerpetualPositionTable.create(
        testConstants.defaultPerpetualPosition,
      );

      await SubaccountTable.create({
        ...testConstants.defaultSubaccount,
        address: testConstants.defaultWalletAddress,
        subaccountNumber: 0,
      });

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/addresses/${testConstants.defaultWalletAddress}`,
      });

      expect(response.body).toEqual({
        subaccounts: [
          {
            address: testConstants.defaultWalletAddress,
            subaccountNumber: 0,
            equity: getFixedRepresentation(0),
            freeCollateral: getFixedRepresentation(0),
            marginEnabled: true,
            updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
            latestProcessedBlockHeight: latestHeight,
            assetPositions: {},
            openPerpetualPositions: {},
          },
        ],
        totalTradingRewards: '0',
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
            msg: 'No subaccounts found for address invalidAddress',
          },
        ],
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
              updatedAtHeight: testConstants.defaultSubaccount.updatedAtHeight,
              latestProcessedBlockHeight: latestHeight,
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
                  realizedPnl: getFixedRepresentation('100'),
                  // size * (index-entry) = 10*(15000-20000) = -50000
                  unrealizedPnl: getFixedRepresentation(-50000),
                  status: testConstants.defaultPerpetualPosition.status,
                  sumOpen: testConstants.defaultPerpetualPosition.sumOpen,
                  sumClose: testConstants.defaultPerpetualPosition.sumClose,
                  createdAt: testConstants.defaultPerpetualPosition.createdAt,
                  createdAtHeight: testConstants.defaultPerpetualPosition.createdAtHeight,
                  exitPrice: null,
                  closedAt: null,
                  subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
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
              equity: getFixedRepresentation(5000),
              freeCollateral: getFixedRepresentation(5000),
              marginEnabled: true,
              updatedAtHeight: testConstants.isolatedSubaccount.updatedAtHeight,
              latestProcessedBlockHeight: latestHeight,
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
              equity: getFixedRepresentation(0),
              freeCollateral: getFixedRepresentation(0),
              marginEnabled: true,
              updatedAtHeight: testConstants.isolatedSubaccount2.updatedAtHeight,
              latestProcessedBlockHeight: latestHeight,
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
      path: `/v4/addresses/${defaultAddress3}/parentSubaccountNumber/` +
          `${testConstants.defaultSubaccount.subaccountNumber}`,
      expectedStatus: 404,
    });

    expect(response.body).toEqual({
      errors: [
        {
          msg: `No subaccounts found for address ${defaultAddress3} and ` +
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
      path: `/v4/addresses/${testConstants.defaultAddress}/parentSubaccountNumber/${parentSubaccountNumber}`,
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

  describe('/:address/testNotification', () => {
    it('Post /:address/testNotification throws error in production', async () => {
      // Mock the config to simulate production environment
      const originalNodeEnv = config.NODE_ENV;
      config.NODE_ENV = 'production';

      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/testNotification`,
        expectedStatus: 404,
      });

      expect(response.statusCode).toEqual(404);
      // Restore the original NODE_ENV
      config.NODE_ENV = originalNodeEnv;
    });
  });

  describe('/:address/registerToken', () => {
    const validToken = 'validToken';
    const validLanguage = 'en';
    const validTimestamp = 1726076825;
    const validMessage = 'Valid message';
    const validSignedMessage = 'Valid signed message';
    const validPubKey = 'Valid public key';

    const verifySignatureMock = Secp256k1.verifySignature as jest.Mock;
    const verifyADR36AminoMock = verifyADR36Amino as jest.Mock;
    const toBech32Mock = toBech32 as jest.Mock;
    let statsSpy = jest.spyOn(stats, 'increment');

    beforeEach(() => {
      verifySignatureMock.mockResolvedValue(true);
      toBech32Mock.mockReturnValue(testConstants.defaultAddress);
      jest.spyOn(DateTime, 'now').mockReturnValue(DateTime.fromSeconds(validTimestamp)); // Mock current time
      statsSpy = jest.spyOn(stats, 'increment');
    });

    afterEach(() => {
      jest.clearAllMocks();
      jest.restoreAllMocks();
    });

    it('Post /:address/registerToken with valid params returns 200', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 200,
      });

      expect(response.body).toEqual({});
      expect(statsSpy).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1, {
        path: '/:address/registerToken',
        method: 'POST',
      });
    });

    it('should register a new token', async () => {
      // Register a new token
      const newToken = 'newToken';
      await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: newToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 200,
      });

      // Check that old tokens are deleted and new token is registered
      const remainingTokens = await FirebaseNotificationTokenTable.findAll({}, []);
      expect(remainingTokens.map((t) => t.token)).toContain(newToken);
    });

    it('Post /:address/registerToken with valid params calls TokenTable registerToken', async () => {
      const registerTokenSpy = jest.spyOn(FirebaseNotificationTokenTable, 'registerToken');
      const token = 'validToken';
      const language = 'en';
      await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 200,
      });
      expect(registerTokenSpy).toHaveBeenCalledWith(
        token, testConstants.defaultAddress, language,
      );
      expect(statsSpy).toHaveBeenCalledWith('comlink.addresses-controller.response_status_code.200', 1, {
        path: '/:address/registerToken',
        method: 'POST',
      });
    });

    it('Post /:address/registerToken with invalid address returns 400', async () => {
      toBech32Mock.mockReturnValue('InvalidAddress');
      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${invalidAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 400,
      });

      expect(response.body).toEqual({
        errors: [
          {
            msg: 'Address invalidAddress is not a valid dYdX V4 address',
          },
        ],
      });
    });

    it.each([
      ['validToken', '', 'Invalid language code', 'language'],
      ['validToken', 'qq', 'Invalid language code', 'language'],
    ])('Post /:address/registerToken with bad language params returns 400', async (token, language, errorMsg, errorParam) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token,
          language,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 400,
      });

      expect(response.body).toEqual({
        errors: [
          {
            location: 'body',
            msg: errorMsg,
            param: errorParam,
            value: language,
          },
        ],
      });
    });

    it.each([
      ['', 'en', 'Token cannot be empty', 'token'],
    ])('Post /:address/registerToken with bad token params returns 400', async (token, language, errorMsg, errorParam) => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: '',
          language,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 400,
      });

      expect(response.body).toEqual({
        errors: [
          {
            location: 'body',
            msg: errorMsg,
            param: errorParam,
            value: token,
          },
        ],
      });
    });

    it('Post /:address/registerToken with invalid signature returns 400', async () => {
      verifySignatureMock.mockResolvedValue(false);

      const response: request.Response = await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: 'Invalid signature',
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 400,
      });

      expect(response.body).toEqual({
        errors: [{ msg: 'Signature verification failed' }],
      });
    });

    it('Post /:address/registerToken with Keplr wallet calls validateSignatureKeplr', async () => {
      verifyADR36AminoMock.mockReturnValue(true);
      const validateSignatureKeplr = jest.spyOn(complianceUtils, 'validateSignatureKeplr');
      await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: true,
        },
        expectedStatus: 200,
      });

      expect(validateSignatureKeplr).toHaveBeenCalledWith(
        expect.anything(),
        testConstants.defaultAddress,
        validMessage,
        validSignedMessage,
        validPubKey,
      );
    });

    it('Post /:address/registerToken with non-Keplr wallet calls validateSignature', async () => {
      const validateSignature = jest.spyOn(complianceUtils, 'validateSignature');
      await sendRequest({
        type: RequestMethod.POST,
        path: `/v4/addresses/${testConstants.defaultAddress}/registerToken`,
        body: {
          token: validToken,
          language: validLanguage,
          timestamp: validTimestamp,
          message: validMessage,
          signedMessage: validSignedMessage,
          pubKey: validPubKey,
          walletIsKeplr: false,
        },
        expectedStatus: 200,
      });

      expect(validateSignature).toHaveBeenCalledWith(
        expect.anything(),
        complianceUtils.AccountVerificationRequiredAction.REGISTER_TOKEN,
        testConstants.defaultAddress,
        validTimestamp,
        validMessage,
        validSignedMessage,
        validPubKey,
        '',
      );
    });
  });
});
