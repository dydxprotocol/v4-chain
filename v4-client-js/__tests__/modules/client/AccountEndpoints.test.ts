import { Network } from '../../../src/clients/constants';
import { IndexerClient } from '../../../src/clients/indexer-client';
import { DYDX_TEST_ADDRESS } from './constants';

describe('IndexerClient', () => {
  const client = new IndexerClient(Network.testnet().indexerConfig);

  describe('Private Endpoints', () => {
    it('Subaccounts', async () => {
      const response = await client.account.getSubaccounts(DYDX_TEST_ADDRESS);
      const subaccounts = response.subaccounts;
      expect(subaccounts.size).not.toBeNull();
      const subaccount0 = subaccounts[0];
      const subaccountNumber = subaccount0.subaccountNumber;
      expect(subaccountNumber).not.toBeNull();
    });

    it('Subaccount 0', async () => {
      const response = await client.account.getSubaccount(DYDX_TEST_ADDRESS, 0);
      const subaccount = response.subaccount;
      expect(subaccount).not.toBeNull();
      const subaccountNumber = subaccount.subaccountNumber;
      expect(subaccountNumber).not.toBeNull();
    });

    it('Asset Positions', async () => {
      const response = await client.account.getSubaccountAssetPositions(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const positions = response.positions;
      expect(positions).not.toBeNull();
      if (positions.length > 0) {
        const position = positions[0];
        expect(position).not.toBeNull();
      }
    });

    it('Perpetual Positions', async () => {
      const response = await client.account.getSubaccountPerpetualPositions(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const positions = response.positions;
      expect(positions).not.toBeNull();
      if (positions.length > 0) {
        const position = positions[0];
        expect(position).not.toBeNull();
      }
    });

    it('Transfers', async () => {
      const response = await client.account.getSubaccountTransfers(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const transfers = response.transfers;
      expect(transfers).not.toBeNull();
      if (transfers.length > 0) {
        const transfer = transfers[0];
        expect(transfer).not.toBeNull();
      }
    });

    it('Orders', async () => {
      const response = await client.account.getSubaccountOrders(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const orders = response;
      expect(orders).not.toBeNull();
      if (orders.length > 0) {
        const order = orders[0];
        expect(order).not.toBeNull();
      }
    });

    it('Fills', async () => {
      const response = await client.account.getSubaccountFills(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const fills = response.fills;
      expect(fills).not.toBeNull();
      if (fills.length > 0) {
        const fill = fills[0];
        expect(fill).not.toBeNull();
      }
    });

    it('Historical PNL', async () => {
      const response = await client.account.getSubaccountHistoricalPNLs(DYDX_TEST_ADDRESS, 0);
      expect(response).not.toBeNull();
      const historicalPnl = response.historicalPnl;
      expect(historicalPnl).not.toBeNull();
      if (historicalPnl.length > 0) {
        const historicalPnl0 = historicalPnl[0];
        expect(historicalPnl0).not.toBeNull();
      }
    });
  });
});
