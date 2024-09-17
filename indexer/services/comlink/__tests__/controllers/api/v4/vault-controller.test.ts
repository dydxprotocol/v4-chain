import {
  dbHelpers,
  testConstants,
  testMocks,
  PnlTicksTable,
  perpetualMarketRefresher,
  BlockTable,
  liquidityTierRefresher,
  SubaccountTable,
  PositionSide,
  PerpetualPositionTable,
  AssetPositionTable,
  FundingIndexUpdatesTable,
  PnlTicksFromDatabase,
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import config from '../../../../src/config';
import { DateTime } from 'luxon';

describe('vault-controller#V4', () => {
  const experimentVaultsPrevVal: string = config.EXPERIMENT_VAULTS;
  const experimentVaultMarketsPrevVal: string = config.EXPERIMENT_VAULT_MARKETS;
  const currentBlockHeight: string = '7';
  const twoHourBlockHeight: string = '5';
  const twoDayBlockHeight: string = '3';
  const currentTime: DateTime = DateTime.utc();
  const twoHoursAgo: DateTime = currentTime.minus({ hour: 2 });
  const twoDaysAgo: DateTime = currentTime.minus({ day: 2 });
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
      await Promise.all([
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: twoDaysAgo.toISO(),
          blockHeight: twoDayBlockHeight,
        }),
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: twoHoursAgo.toISO(),
          blockHeight: twoHourBlockHeight,
        }),
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: currentTime.toISO(),
          blockHeight: currentBlockHeight,
        }),
      ]);
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

    it.each([
      ['no resolution', '', [1, 2]],
      ['daily resolution', '?resolution=day', [1, 2]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3]],
    ])('Get /megavault/historicalPnl with single vault subaccount (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
    ) => {
      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/megavault/historicalPnl${queryParam}`,
      });

      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining(createdPnlTicks[index]);
          }),
        ),
      );
    });

    it.each([
      ['no resolution', '', [1, 2]],
      ['daily resolution', '?resolution=day', [1, 2]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3]],
    ])('Get /megavault/historicalPnl with 2 vault subaccounts (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
    ) => {
      config.EXPERIMENT_VAULTS = [
        testConstants.defaultPnlTick.subaccountId,
        testConstants.vaultSubaccountId,
      ].join(',');
      config.EXPERIMENT_VAULT_MARKETS = [
        testConstants.defaultPerpetualMarket.clobPairId,
        testConstants.defaultPerpetualMarket2.clobPairId,
      ].join(',');

      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/megavault/historicalPnl${queryParam}`,
      });

      const expectedPnlTickBase: any = {
        equity: (parseFloat(testConstants.defaultPnlTick.equity) * 2).toString(),
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) * 2).toString(),
        netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) * 2).toString(),
      };

      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining({
              ...expectedPnlTickBase,
              createdAt: createdPnlTicks[index].createdAt,
              blockHeight: createdPnlTicks[index].blockHeight,
              blockTime: createdPnlTicks[index].blockTime,
            });
          }),
        ),
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

    it.each([
      ['no resolution', '', [1, 2]],
      ['daily resolution', '?resolution=day', [1, 2]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3]],
    ])('Get /vaults/historicalPnl with single vault subaccount (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
    ) => {
      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/vaults/historicalPnl${queryParam}`,
      });

      expect(response.body.vaultsPnl).toHaveLength(1);

      expect(response.body.vaultsPnl[0]).toEqual({
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining(createdPnlTicks[index]);
          }),
        ),
      });
    });

    it.each([
      ['no resolution', '', [1, 2], [5, 6]],
      ['daily resolution', '?resolution=day', [1, 2], [5, 6]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3], [5, 6, 7]],
    ])('Get /vaults/historicalPnl with 2 vault subaccounts (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex1: number[],
      expectedTicksIndex2: number[],
    ) => {
      config.EXPERIMENT_VAULTS = [
        testConstants.defaultPnlTick.subaccountId,
        testConstants.vaultSubaccountId,
      ].join(',');
      config.EXPERIMENT_VAULT_MARKETS = [
        testConstants.defaultPerpetualMarket.clobPairId,
        testConstants.defaultPerpetualMarket2.clobPairId,
      ].join(',');

      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/vaults/historicalPnl${queryParam}`,
      });

      const expectedVaultPnl: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: expectedTicksIndex1.map((index: number) => {
          return createdPnlTicks[index];
        }),
      };

      const expectedVaultPnl2: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket2.ticker,
        historicalPnl: expectedTicksIndex2.map((index: number) => {
          return createdPnlTicks[index];
        }),
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
          effectiveAtHeight: twoDayBlockHeight,
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
          effectiveAtHeight: twoDayBlockHeight,
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

  async function createPnlTicks(): Promise<PnlTicksFromDatabase[]> {
    return Promise.all([
      PnlTicksTable.create(testConstants.defaultPnlTick),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        blockTime: twoDaysAgo.toISO(),
        createdAt: twoDaysAgo.toISO(),
        blockHeight: twoDayBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        blockTime: twoHoursAgo.toISO(),
        createdAt: twoHoursAgo.toISO(),
        blockHeight: twoHourBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        blockTime: currentTime.toISO(),
        createdAt: currentTime.toISO(),
        blockHeight: currentBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
        blockTime: twoDaysAgo.toISO(),
        createdAt: twoDaysAgo.toISO(),
        blockHeight: twoDayBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
        blockTime: twoHoursAgo.toISO(),
        createdAt: twoHoursAgo.toISO(),
        blockHeight: twoHourBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
        blockTime: currentTime.toISO(),
        createdAt: currentTime.toISO(),
        blockHeight: currentBlockHeight,
      }),
    ]);
  }
});
