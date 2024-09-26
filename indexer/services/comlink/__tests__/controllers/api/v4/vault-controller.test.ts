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
} from '@dydxprotocol-indexer/postgres';
import { RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { getFixedRepresentation, sendRequest } from '../../../helpers/helpers';
import { DateTime } from 'luxon';
import Big from 'big.js';

describe('vault-controller#V4', () => {
  const latestBlockHeight: string = '25';
  const currentBlockHeight: string = '7';
  const twoHourBlockHeight: string = '5';
  const twoDayBlockHeight: string = '3';
  const currentTime: DateTime = DateTime.utc().startOf('day').minus({ hour: 5 });
  const latestTime: DateTime = currentTime.plus({ second: 5 });
  const twoHoursAgo: DateTime = currentTime.minus({ hour: 2 });
  const twoDaysAgo: DateTime = currentTime.minus({ day: 2 });
  const initialFundingIndex: string = '10000';
  const vault1Equity: number = 159500;
  const vault2Equity: number = 10000;

  beforeAll(async () => {
    await dbHelpers.migrate();
  });

  afterAll(async () => {
    await dbHelpers.teardown();
  });

  describe('GET /v1', () => {
    beforeEach(async () => {
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
      ]);
      await SubaccountTable.create(testConstants.vaultSubaccount);
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
      ['no resolution', '', [1, 2]],
      ['daily resolution', '?resolution=day', [1, 2]],
      ['hourly resolution', '?resolution=hour', [1, 2, 3]],
    ])('Get /megavault/historicalPnl with 2 vault subaccounts (%s)', async (
      _name: string,
      queryParam: string,
      expectedTicksIndex: number[],
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
      const response: request.Response = await sendRequest({
        type: RequestMethod.GET,
        path: `/v4/vault/v1/megavault/historicalPnl${queryParam}`,
      });

      const expectedPnlTickBase: any = {
        equity: (parseFloat(testConstants.defaultPnlTick.equity) * 2).toString(),
        totalPnl: (parseFloat(testConstants.defaultPnlTick.totalPnl) * 2).toString(),
        netTransfers: (parseFloat(testConstants.defaultPnlTick.netTransfers) * 2).toString(),
      };
      const finalTick: PnlTicksFromDatabase = {
        ...expectedPnlTickBase,
        equity: Big(vault1Equity).add(vault2Equity).toFixed(),
        blockHeight: latestBlockHeight,
        blockTime: latestTime.toISO(),
        createdAt: latestTime.toISO(),
      };

      expect(response.body.megavaultPnl).toHaveLength(expectedTicksIndex.length + 1);
      expect(response.body.megavaultPnl).toEqual(
        expect.arrayContaining(
          expectedTicksIndex.map((index: number) => {
            return expect.objectContaining({
              ...expectedPnlTickBase,
              createdAt: createdPnlTicks[index].createdAt,
              blockHeight: createdPnlTicks[index].blockHeight,
              blockTime: createdPnlTicks[index].blockTime,
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
