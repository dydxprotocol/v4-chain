import {
  dbHelpers,
  testConstants,
  testMocks,
  PnlTicksCreateObject,
  PnlTicksTable,
  perpetualMarketRefresher,
  BlockTable,
  liquidityTierRefresher,
  SubaccountTable,
  PositionSide,
  PerpetualPositionTable,
  AssetPositionTable,
  FundingIndexUpdatesTable,
} from '@dydxprotocol-indexer/postgres';
import { PnlTicksResponseObject, RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import config from '../../../../src/config';

describe('vault-controller#V4', () => {
  const experimentVaultsPrevVal: string = config.EXPERIMENT_VAULTS;
  const experimentVaultMarketsPrevVal: string = config.EXPERIMENT_VAULT_MARKETS;
  const blockHeight: string = '3';
  const initialFundingIndex: string = '10000';

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /v1', () => {
    beforeEach(async () => {
      config.EXPERIMENT_VAULTS = testConstants.defaultPnlTick.subaccountId;
      config.EXPERIMENT_VAULT_MARKETS = testConstants.defaultPerpetualMarket.clobPairId;
      await testMocks.seedData();
      await perpetualMarketRefresher.updatePerpetualMarkets();
      await liquidityTierRefresher.updateLiquidityTiers();
      await BlockTable.create({
        ...testConstants.defaultBlock,
        blockHeight,
      });
      await SubaccountTable.create(testConstants.vaultSubaccount);
    });

    afterEach(async () => {
      config.EXPERIMENT_VAULTS = experimentVaultsPrevVal;
      config.EXPERIMENT_VAULT_MARKETS = experimentVaultMarketsPrevVal;
      await dbHelpers.clearData();
    });

    it('Get /megavault/historicalPnl with no vault subaccounts', async () => {
      config.EXPERIMENT_VAULTS = '';
      config.EXPERIMENT_VAULT_MARKETS = '';

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/historicalPnl',
      });

      expect(response.body.megavaultPnl).toEqual([]);
    });

    it('Get /megavault/historicalPnl with single vault subaccount', async () => {
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/historicalPnl',
      });

      const expectedPnlTickResponse: PnlTicksResponseObject = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };

      const expectedPnlTick2Response: any = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          createdAt,
        ),
      };

      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTick2Response,
          }),
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );
    });

    it('Get /megavault/historicalPnl with 2 vault subaccounts', async () => {
      config.EXPERIMENT_VAULTS = [
        testConstants.defaultPnlTick.subaccountId,
        testConstants.vaultSubaccountId,
      ].join(',');
      config.EXPERIMENT_VAULT_MARKETS = [
        testConstants.defaultPerpetualMarket.clobPairId,
        testConstants.defaultPerpetualMarket2.clobPairId,
      ].join(',');

      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/historicalPnl',
      });

      const expectedPnlTickResponse: any = {
        // id and subaccountId don't matter
        equity: (parseFloat(testConstants.defaultPnlTick.equity) +
            parseFloat(pnlTick2.equity)).toString(),
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) +
            parseFloat(pnlTick2.totalPnl)).toString(),
        netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) +
            parseFloat(pnlTick2.netTransfers)).toString(),
        createdAt: testConstants.defaultPnlTick.createdAt,
        blockHeight: testConstants.defaultPnlTick.blockHeight,
        blockTime: testConstants.defaultPnlTick.blockTime,
      };

      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      );
    });

    it('Get /vaults/historicalPnl with no vault subaccounts', async () => {
      config.EXPERIMENT_VAULTS = '';
      config.EXPERIMENT_VAULT_MARKETS = '';

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/vaults/historicalPnl',
      });

      expect(response.body.vaultsPnl).toEqual([]);
    });

    it('Get /vaults/historicalPnl with single vault subaccount', async () => {
      const createdAt: string = '2000-05-25T00:00:00.000Z';
      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/vaults/historicalPnl',
      });

      const expectedPnlTickResponse: PnlTicksResponseObject = {
        ...testConstants.defaultPnlTick,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          testConstants.defaultPnlTick.createdAt,
        ),
      };

      const expectedPnlTick2Response: any = {
        ...testConstants.defaultPnlTick,
        createdAt,
        blockHeight,
        id: PnlTicksTable.uuid(
          testConstants.defaultPnlTick.subaccountId,
          createdAt,
        ),
      };

      expect(response.body.vaultsPnl).toHaveLength(1);

      expect(response.body.vaultsPnl[0]).toEqual({
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: expect.arrayContaining([
          expect.objectContaining({
            ...expectedPnlTick2Response,
          }),
          expect.objectContaining({
            ...expectedPnlTickResponse,
          }),
        ]),
      });
    });

    it('Get /vaults/historicalPnl with 2 vault subaccounts', async () => {
      config.EXPERIMENT_VAULTS = [
        testConstants.defaultPnlTick.subaccountId,
        testConstants.vaultSubaccountId,
      ].join(',');
      config.EXPERIMENT_VAULT_MARKETS = [
        testConstants.defaultPerpetualMarket.clobPairId,
        testConstants.defaultPerpetualMarket2.clobPairId,
      ].join(',');

      const pnlTick2: PnlTicksCreateObject = {
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
      };
      await Promise.all([
        PnlTicksTable.create(testConstants.defaultPnlTick),
        PnlTicksTable.create(pnlTick2),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/vaults/historicalPnl',
      });

      const expectedVaultPnl: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: [
          {
            ...testConstants.defaultPnlTick,
            id: PnlTicksTable.uuid(
              testConstants.defaultPnlTick.subaccountId,
              testConstants.defaultPnlTick.createdAt,
            ),
          },
        ],
      };

      const expectedVaultPnl2: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket2.ticker,
        historicalPnl: [
          {
            ...pnlTick2,
            id: PnlTicksTable.uuid(
              pnlTick2.subaccountId,
              pnlTick2.createdAt,
            ),
          },
        ],
      };

      expect(response.body.vaultsPnl).toEqual(
        expect.arrayContaining([
          expect.objectContaining({
            ...expectedVaultPnl,
          }),
          expect.objectContaining({
            ...expectedVaultPnl2,
          }),
        ]),
      );
    });

    it('Get /megavault/positions with no vault subaccount', async () => {
      config.EXPERIMENT_VAULTS = '';
      config.EXPERIMENT_VAULT_MARKETS = '';

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/positions',
      });

      expect(response.body).toEqual({
        positions: [],
      });
    });

    it('Get /megavault/positions with 1 vault subaccount', async () => {
      await Promise.all([
        PerpetualPositionTable.create(
          testConstants.defaultPerpetualPosition,
        ),
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: initialFundingIndex,
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: blockHeight,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/positions',
      });

      expect(response.body).toEqual({
        positions: [
          {
            equity: getFixedRepresentation(159500),
            perpetualPosition: {
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
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
            assetPosition: {
              symbol: testConstants.defaultAsset.symbol,
              size: '9500',
              side: PositionSide.LONG,
              assetId: testConstants.defaultAssetPosition.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
            ticker: testConstants.defaultPerpetualMarket.ticker,
          },
        ],
      });
    });

    it('Get /megavault/positions with 2 vault subaccount, 1 with no perpetual', async () => {
      config.EXPERIMENT_VAULTS = [
        testConstants.defaultPnlTick.subaccountId,
        testConstants.vaultSubaccountId,
      ].join(',');
      config.EXPERIMENT_VAULT_MARKETS = [
        testConstants.defaultPerpetualMarket.clobPairId,
        testConstants.defaultPerpetualMarket2.clobPairId,
      ].join(',');

      await Promise.all([
        PerpetualPositionTable.create(
          testConstants.defaultPerpetualPosition,
        ),
        AssetPositionTable.upsert(testConstants.defaultAssetPosition),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition,
          subaccountId: testConstants.vaultSubaccountId,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          fundingIndex: initialFundingIndex,
          effectiveAtHeight: testConstants.createdHeight,
        }),
        FundingIndexUpdatesTable.create({
          ...testConstants.defaultFundingIndexUpdate,
          eventId: testConstants.defaultTendermintEventId2,
          effectiveAtHeight: blockHeight,
        }),
      ]);

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/positions',
      });

      expect(response.body).toEqual({
        positions: [
          // Same position as test with a single vault subaccount.
          {
            equity: getFixedRepresentation(159500),
            perpetualPosition: {
              market: testConstants.defaultPerpetualMarket.ticker,
              size: testConstants.defaultPerpetualPosition.size,
              side: testConstants.defaultPerpetualPosition.side,
              entryPrice: getFixedRepresentation(
                testConstants.defaultPerpetualPosition.entryPrice!,
              ),
              maxSize: testConstants.defaultPerpetualPosition.maxSize,
              netFunding: getFixedRepresentation('199500'),
              realizedPnl: getFixedRepresentation('199500'),
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
            assetPosition: {
              symbol: testConstants.defaultAsset.symbol,
              size: '9500',
              side: PositionSide.LONG,
              assetId: testConstants.defaultAssetPosition.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
            ticker: testConstants.defaultPerpetualMarket.ticker,
          },
          {
            equity: getFixedRepresentation(10000),
            perpetualPosition: undefined,
            assetPosition: {
              symbol: testConstants.defaultAsset.symbol,
              size: testConstants.defaultAssetPosition.size,
              side: PositionSide.LONG,
              assetId: testConstants.defaultAssetPosition.assetId,
              subaccountNumber: testConstants.defaultSubaccount.subaccountNumber,
            },
            ticker: testConstants.defaultPerpetualMarket2.ticker,
          },
        ],
      });
    });
  });
});
