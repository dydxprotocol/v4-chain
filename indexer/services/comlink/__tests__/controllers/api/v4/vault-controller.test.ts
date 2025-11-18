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
  VaultTable,
  MEGAVAULT_MODULE_ADDRESS,
  MEGAVAULT_SUBACCOUNT_ID,
  TransferTable,
  VaultPnlTicksView,
} from '@dydxprotocol-indexer/postgres';
import { PnlTicksResponseObject, RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import { DateTime, Settings } from 'luxon';
import Big from 'big.js';
import config from '../../../../src/config';
import { clearVaultStartPnl, startVaultStartPnlCache } from '../../../../src/caches/vault-start-pnl';
import { pnlTicksToResponseObject } from '../../../../src/request-helpers/request-transformer';

describe('vault-controller#V4', () => {
  const latestBlockHeight: string = '25';
  const currentHourBlockHeight: string = '10';
  const currentDayBlockHeight: string = '9';
  const twoHourBlockHeight: string = '7';
  const almostTwoDayBlockHeight: string = '5';
  const twoDayBlockHeight: string = '3';
  const currentDay: DateTime = DateTime.utc().startOf('day').minus({ hour: 5 });
  const currentHour: DateTime = currentDay.plus({ hour: 1 });
  const latestTime: DateTime = currentDay.plus({ minute: 90 });
  const twoHoursAgo: DateTime = currentDay.minus({ hour: 2 });
  const twoDaysAgo: DateTime = currentDay.minus({ day: 2 });
  const almostTwoDaysAgo: DateTime = currentDay.minus({ hour: 47 });
  const initialFundingIndex: string = '10000';
  const vault1Equity: number = 159500;
  const vault2Equity: number = 10000;
  const mainVaultEquity: number = 10000;
  const vaultPnlHistoryHoursPrev: number = config.VAULT_PNL_HISTORY_HOURS;
  const vaultPnlLastPnlWindowPrev: number = config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS;
  const vaultPnlStartDatePrev: string = config.VAULT_PNL_START_DATE;

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /v1', () => {
    beforeEach(async () => {
      // Get a week of data for hourly pnl ticks.
      config.VAULT_PNL_HISTORY_HOURS = 168;
      // Use last 48 hours to get latest pnl tick for tests.
      config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS = 48;
      // Use a time before all pnl ticks as the pnl start date.
      config.VAULT_PNL_START_DATE = '2020-01-01T00:00:00Z';
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
          time: currentDay.toISO(),
          blockHeight: currentDayBlockHeight,
        }),
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: latestTime.toISO(),
          blockHeight: latestBlockHeight,
        }),
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: almostTwoDaysAgo.toISO(),
          blockHeight: almostTwoDayBlockHeight,
        }),
        BlockTable.create({
          ...testConstants.defaultBlock,
          time: currentHour.toISO(),
          blockHeight: currentHourBlockHeight,
        }),
      ]);
      await SubaccountTable.create(testConstants.vaultSubaccount);
      await SubaccountTable.create({
        address: MEGAVAULT_MODULE_ADDRESS,
        subaccountNumber: 0,
        updatedAt: latestTime.toISO(),
        updatedAtHeight: latestBlockHeight,
      });
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
      Settings.now = () => latestTime.valueOf();
    });

    afterEach(async () => {
      await dbHelpers.clearData();
      await VaultPnlTicksView.refreshDailyView();
      await VaultPnlTicksView.refreshHourlyView();
      clearVaultStartPnl();
      config.VAULT_PNL_HISTORY_HOURS = vaultPnlHistoryHoursPrev;
      config.VAULT_LATEST_PNL_TICK_WINDOW_HOURS = vaultPnlLastPnlWindowPrev;
      config.VAULT_PNL_START_DATE = vaultPnlStartDatePrev;
      Settings.now = () => new Date().valueOf();
    });

    it.each([
      ['no resolution', '', [1, 2], 4, undefined],
      ['daily resolution', '?resolution=day', [1, 2], 4, undefined],
      ['hourly resolution', '?resolution=hour', [1, 2, 3, 4], 4, undefined],
      ['start date adjust PnL', '?resolution=hour', [1, 2, 3, 4], 4, twoDaysAgo.toISO()],
    ])('Get /megavault/historicalPnl with single vault subaccount (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
      finalTickIndex: number,
      startDate: string | undefined,
    ) => {
      if (startDate !== undefined) {
        config.VAULT_PNL_START_DATE = startDate;
      }
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultSubaccount.address,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
      const createdPnlTicksFromDatabase: PnlTicksFromDatabase[] = await createPnlTicks();
      const createdPnlTicks
      : PnlTicksResponseObject[] = createdPnlTicksFromDatabase.map(pnlTicksToResponseObject);
      // Adjust PnL by total pnl of start date
      if (startDate !== undefined) {
        for (const createdPnlTick of createdPnlTicks) {
          createdPnlTick.totalPnl = Big(createdPnlTick.totalPnl).sub('10000').toFixed();
        }
      }
      const finalTick: PnlTicksResponseObject = {
        ...createdPnlTicks[finalTickIndex],
        equity: Big(vault1Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/megavault/historicalPnl${queryParam}`,
      });

      expect(response.body.megavaultPnl).toHaveLength(expectedTicksIndex.length + 1);
      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining(createdPnlTicks[index]);
          }).concat([finalTick]),
        ),
      );
    });

    it.each([
      ['no resolution', '', [1, 2], [undefined, 7], [11, 12]],
      ['daily resolution', '?resolution=day', [1, 2], [undefined, 7], [11, 12]],
      [
        'hourly resolution',
        '?resolution=hour',
        [1, 2, 3, 4],
        [undefined, 7, 8, 9],
        [11, 12, 13, 14],
      ],
    ])('Get /megavault/historicalPnl with 2 vault subaccounts and main subaccount (%s), ' +
       'excludes tick with missing vault ticks', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex1: (number | undefined)[],
      expectedTicksIndex2: (number | undefined)[],
      expectedTicksIndexMain: (number | undefined)[],
    ) => {
      const expectedTicksArray: (number | undefined)[][] = [
        expectedTicksIndex1,
        expectedTicksIndex2,
        expectedTicksIndexMain,
      ];
      await Promise.all([
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.defaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
          createdAt: twoDaysAgo.toISO(),
        }),
        // Single tick for this vault will be excluded from result.
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.vaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
          createdAt: almostTwoDaysAgo.toISO(),
        }),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
        }),
        TransferTable.create({
          ...testConstants.defaultTransfer,
          recipientSubaccountId: MEGAVAULT_SUBACCOUNT_ID,
          createdAt: twoDaysAgo.toISO(),
        }),
      ]);

      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks(
        true, // createMainSubaccounPnlTicks
      );
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/megavault/historicalPnl${queryParam}`,
      });

      const expectedPnlTickBase: any = {
        equity: (parseFloat(testConstants.defaultPnlTick.equity) * 3).toString(),
        // total pnl should be fetched from latest hourly pnl tick.
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) * 4).toString(),
        netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) * 3).toString(),
      };
      const finalTick: PnlTicksFromDatabase = {
        ...expectedPnlTickBase,
        equity: Big(vault1Equity).add(vault2Equity).add(mainVaultEquity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };

      expect(response.body.megavaultPnl).toHaveLength(expectedTicksIndex1.length + 1);
      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining(
          expectedTicksIndex1.map((_: number | undefined, pos: number) => {
            const pnlTickBase: any = {
              equity: '0',
              totalPnl: '0',
              netTransfers: '0',
            };
            let expectedTick: PnlTicksFromDatabase;
            for (const expectedTicks of expectedTicksArray) {
              if (expectedTicks[pos] !== undefined) {
                expectedTick = createdPnlTicks[expectedTicks[pos]!];
                pnlTickBase.equity = Big(pnlTickBase.equity).add(expectedTick.equity).toFixed();
                pnlTickBase.totalPnl = Big(pnlTickBase.totalPnl)
                  .add(expectedTick.totalPnl)
                  .toFixed();
                pnlTickBase.netTransfers = Big(pnlTickBase.netTransfers)
                  .add(expectedTick.netTransfers)
                  .toFixed();
              }
            }
            return expect.objectContaining({
              ...pnlTickBase,
              createdAt: expectedTick!.createdAt,
              blockHeight: expectedTick!.blockHeight,
              blockTime: expectedTick!.blockTime,
            });
          }).concat([expect.objectContaining(finalTick)]),
        ),
      );
    });

    it('Get /vaults/historicalPnl with no vault subaccounts', async () => {
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/vaults/historicalPnl',
      });

      expect(response.body.vaultsPnl).toEqual([]);
    });

    it.each([
      ['no resolution', '', [1, 2], 4],
      ['daily resolution', '?resolution=day', [1, 2], 4],
      ['hourly resolution', '?resolution=hour', [1, 2, 3, 4], 4],
    ])('Get /vaults/historicalPnl with single vault subaccount (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
      currentTickIndex: number,
    ) => {
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultAddress,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
      const createdPnlTicksFromDatabase: PnlTicksFromDatabase[] = await createPnlTicks();
      const createdPnlTicks
      : PnlTicksResponseObject[] = createdPnlTicksFromDatabase.map(pnlTicksToResponseObject);
      const finalTick: PnlTicksResponseObject = {
        ...createdPnlTicks[currentTickIndex],
        equity: Big(vault1Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/vaults/historicalPnl${queryParam}`,
      });

      expect(response.body.vaultsPnl).toHaveLength(1);
      expect(response.body.vaultsPnl[0].historicalPnl).toHaveLength(expectedTicksIndex.length + 1);
      expect(response.body.vaultsPnl[0]).toEqual({
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining(createdPnlTicks[index]);
          }).concat(finalTick),
        ),
      });
    });

    it.each([
      ['no resolution', '', [1, 2], [6, 7], 4, 9],
      ['daily resolution', '?resolution=day', [1, 2], [6, 7], 4, 9],
      ['hourly resolution', '?resolution=hour', [1, 2, 3, 4], [6, 7, 8, 9], 4, 9],
    ])('Get /vaults/historicalPnl with 2 vault subaccounts (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex1: number[],
      expectedTicksIndex2: number[],
      currentTickIndex1: number,
      currentTickIndex2: number,
    ) => {
      await Promise.all([
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.defaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
        }),
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.vaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        }),
      ]);
      const createdPnlTicksFromDatabase: PnlTicksFromDatabase[] = await createPnlTicks();
      const createdPnlTicks
      : PnlTicksResponseObject[] = createdPnlTicksFromDatabase.map(pnlTicksToResponseObject);
      const finalTick1: PnlTicksResponseObject = {
        ...createdPnlTicks[currentTickIndex1],
        equity: Big(vault1Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };
      const finalTick2: PnlTicksResponseObject = {
        ...createdPnlTicks[currentTickIndex2],
        equity: Big(vault2Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };

      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/vaults/historicalPnl${queryParam}`,
      });

      const expectedVaultPnl: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket.ticker,
        historicalPnl: expectedTicksIndex1.map((index: number) => {
          return createdPnlTicks[index];
        }).concat(finalTick1),
      };

      const expectedVaultPnl2: VaultHistoricalPnl = {
        ticker: testConstants.defaultPerpetualMarket2.ticker,
        historicalPnl: expectedTicksIndex2.map((index: number) => {
          return createdPnlTicks[index];
        }).concat(finalTick2),
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
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: '/v4/vault/v1/megavault/positions',
      });

      expect(response.body).toEqual({
        positions: [],
      });
    });

    it('Get /megavault/positions with 1 vault subaccount', async () => {
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultAddress,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
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

    it('Get /megavault/positions with 2 vault subaccount, 1 with no perpetual, 1 invalid', async () => {
      await Promise.all([
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.defaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
        }),
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.vaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        }),
        VaultTable.create({
          ...testConstants.defaultVault,
          address: 'invalid',
          clobPairId: '999',
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
              realizedPnl: getFixedRepresentation('100'),
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

    it('Get /megavault/historicalPnl returns cached results within TTL', async () => {
      const originalCacheTtl = config.VAULT_CACHE_TTL_MS;
      config.VAULT_CACHE_TTL_MS = 60000; // 1 minute

      try {
        // Setup: Create a vault and some PnL ticks
        await VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.defaultSubaccount.address,
          clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
        });
        // We still need some initial PnL data for the endpoint to return something.
        await createPnlTicks();

        // First request - should populate the cache
        const response1: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: '/v4/vault/v1/megavault/historicalPnl?resolution=hour',
        });
        expect(response1.status).toBe(200);
        expect(response1.body.megavaultPnl.length).toBeGreaterThan(0);

        // Modify underlying data that affects current equity calculation
        const newAssetSize = '999999999';
        await AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition, // Use existing asset details
          subaccountId: testConstants.defaultSubaccountId, // Target the vault subaccount
          size: newAssetSize, // Change the size
        });

        // Second request - should hit the cache and return the OLD data
        const response2: request.Response = await sendRequest({
          type: RequestMethod.GET,
          path: '/v4/vault/v1/megavault/historicalPnl?resolution=hour',
        });
        expect(response2.status).toBe(200);

        // Assert that the second response is identical to the first (cached) response
        expect(response2.body).toEqual(response1.body);

        // Verify the FINAL tick's equity in the cached response does NOT reflect the change
        // The final tick represents the current state and would change if not cached.
        const finalCachedTick = response2.body.megavaultPnl[response2.body.megavaultPnl.length - 1];
        const originalFinalTick = response1.body.megavaultPnl[
          response1.body.megavaultPnl.length - 1];
        // Should match the original final tick equity
        expect(finalCachedTick.equity).toEqual(originalFinalTick.equity);
      } finally {
        // Restore original config value
        config.VAULT_CACHE_TTL_MS = originalCacheTtl;
      }
    });
  });

  async function createPnlTicks(
    createMainSubaccountPnlTicks: boolean = false,
  ): Promise<PnlTicksFromDatabase[]> {
    const createdTicks: PnlTicksFromDatabase[] = await Promise.all([
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
        blockTime: currentDay.toISO(),
        createdAt: currentDay.toISO(),
        blockHeight: currentDayBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        totalPnl: (2 * parseFloat(testConstants.defaultPnlTick.totalPnl)).toString(),
        blockTime: currentHour.toISO(),
        createdAt: currentHour.toISO(),
        blockHeight: currentHourBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
      }),
      // Invalid pnl tick to be excluded as only a single pnl tick but 2 pnl ticks should exist.
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
        blockTime: almostTwoDaysAgo.toISO(),
        createdAt: almostTwoDaysAgo.toISO(),
        blockHeight: almostTwoDayBlockHeight,
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
        blockTime: currentDay.toISO(),
        createdAt: currentDay.toISO(),
        blockHeight: currentDayBlockHeight,
      }),
      PnlTicksTable.create({
        ...testConstants.defaultPnlTick,
        subaccountId: testConstants.vaultSubaccountId,
        blockTime: currentHour.toISO(),
        createdAt: currentHour.toISO(),
        blockHeight: currentHourBlockHeight,
      }),
    ]);

    if (createMainSubaccountPnlTicks) {
      const mainSubaccountTicks: PnlTicksFromDatabase[] = await Promise.all([
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
        }),
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
          blockTime: twoDaysAgo.toISO(),
          createdAt: twoDaysAgo.toISO(),
          blockHeight: twoDayBlockHeight,
        }),
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
          blockTime: twoHoursAgo.toISO(),
          createdAt: twoHoursAgo.toISO(),
          blockHeight: twoHourBlockHeight,
        }),
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
          blockTime: currentDay.toISO(),
          createdAt: currentDay.toISO(),
          blockHeight: currentDayBlockHeight,
        }),
        PnlTicksTable.create({
          ...testConstants.defaultPnlTick,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
          blockTime: currentHour.toISO(),
          createdAt: currentHour.toISO(),
          blockHeight: currentHourBlockHeight,
        }),
      ]);
      createdTicks.push(...mainSubaccountTicks);
    }
    await VaultPnlTicksView.refreshDailyView();
    await VaultPnlTicksView.refreshHourlyView();
    await startVaultStartPnlCache();

    return createdTicks;
  }
});
