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
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import { DateTime } from 'luxon';
import Big from 'big.js';
import config from '../../../../src/config';

describe('vault-controller#V4', () => {
  const latestBlockHeight: string = '25';
  const currentBlockHeight: string = '9';
  const twoHourBlockHeight: string = '7';
  const almostTwoDayBlockHeight: string = '5';
  const twoDayBlockHeight: string = '3';
  const currentTime: DateTime = DateTime.utc().startOf('day').minus({ hour: 5 });
  const latestTime: DateTime = currentTime.plus({ second: 5 });
  const twoHoursAgo: DateTime = currentTime.minus({ hour: 2 });
  const twoDaysAgo: DateTime = currentTime.minus({ day: 2 });
  const almostTwoDaysAgo: DateTime = currentTime.minus({ hour: 47 });
  const initialFundingIndex: string = '10000';
  const vault1Equity: number = 159500;
  const vault2Equity: number = 10000;
  const mainVaultEquity: number = 10000;
  const vaultPnlHistoryHoursPrev: number = config.VAULT_PNL_HISTORY_HOURS;

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
    });

    afterEach(async () => {
      await dbHelpers.clearData();
      config.VAULT_PNL_HISTORY_HOURS = vaultPnlHistoryHoursPrev;
    });

    it('Get /megavault/historicalPnl with no vault subaccounts', async () => {
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
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultSubaccount.address,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();
      const finalTick: PnlTicksFromDatabase = {
        ...createdPnlTicks[expectedTicksIndex[expectedTicksIndex.length - 1]],
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
      ['no resolution', '', [1, 2], [undefined, 6], [9, 10]],
      ['daily resolution', '?resolution=day', [1, 2], [undefined, 6], [9, 10]],
      [
        'hourly resolution',
        '?resolution=hour',
        [1, undefined, 2, 3],
        [undefined, 5, 6, 7],
        [9, undefined, 10, 11],
      ],
    ])('Get /megavault/historicalPnl with 2 vault subaccounts and main subaccount (%s)', async (
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
        }),
        VaultTable.create({
          ...testConstants.defaultVault,
          address: testConstants.vaultAddress,
          clobPairId: testConstants.defaultPerpetualMarket2.clobPairId,
        }),
        AssetPositionTable.upsert({
          ...testConstants.defaultAssetPosition,
          subaccountId: MEGAVAULT_SUBACCOUNT_ID,
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
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) * 3).toString(),
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
      ['no resolution', '', [1, 2]],
      ['daily resolution', '?resolution=day', [1, 2]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3]],
    ])('Get /vaults/historicalPnl with single vault subaccount (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
    ) => {
      await VaultTable.create({
        ...testConstants.defaultVault,
        address: testConstants.defaultAddress,
        clobPairId: testConstants.defaultPerpetualMarket.clobPairId,
      });
      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();
      const finalTick: PnlTicksFromDatabase = {
        ...createdPnlTicks[expectedTicksIndex[expectedTicksIndex.length - 1]],
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
      ['no resolution', '', [1, 2], [5, 6]],
      ['daily resolution', '?resolution=day', [1, 2], [5, 6]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3], [5, 6, 7]],
    ])('Get /vaults/historicalPnl with 2 vault subaccounts (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex1: number[],
      expectedTicksIndex2: number[],
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
      const createdPnlTicks: PnlTicksFromDatabase[] = await createPnlTicks();
      const finalTick1: PnlTicksFromDatabase = {
        ...createdPnlTicks[expectedTicksIndex1[expectedTicksIndex1.length - 1]],
        equity: Big(vault1Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };
      const finalTick2: PnlTicksFromDatabase = {
        ...createdPnlTicks[expectedTicksIndex2[expectedTicksIndex2.length - 1]],
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
        blockTime: currentTime.toISO(),
        createdAt: currentTime.toISO(),
        blockHeight: currentBlockHeight,
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
          blockTime: currentTime.toISO(),
          createdAt: currentTime.toISO(),
          blockHeight: currentBlockHeight,
        }),
      ]);
      createdTicks.push(...mainSubaccountTicks);
    }

    return createdTicks;
  }
});
