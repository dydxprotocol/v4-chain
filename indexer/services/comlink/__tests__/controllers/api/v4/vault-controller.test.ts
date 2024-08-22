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
} from '@dydxprotocol-indexer/postgres';
import { PnlTicksResponseObject, RequestMethod, VaultHistoricalPnl } from '../../../../src/types';
import request from 'supertest';
import { sendRequest } from '../../../helpers/helpers';
import config from '../../../../src/config';

describe('vault-controller#V4', () => {
  const experimentVaultsPrevVal: string = config.EXPERIMENT_VAULTS;
  const experimentVaultMarketsPrevVal: string = config.EXPERIMENT_VAULT_MARKETS;
  const blockHeight: string = '3';

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

    it('Get /vault/v1/vaults/historicalPnl with 2 vault subaccounts', async () => {
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
  });
});
