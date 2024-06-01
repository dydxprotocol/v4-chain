/**
 * Simple JS example demostrating accessing subaccount data with Indexer REST endpoints
 */

import { Network } from '../src/clients/constants';
import { IndexerClient } from '../src/clients/indexer-client';
import { DYDX_TEST_ADDRESS } from './constants';

async function test(): Promise<void> {
  const client = new IndexerClient(Network.testnet().indexerConfig);
  const address = DYDX_TEST_ADDRESS;

  // Get subaccounts
  try {
    const response = await client.account.getSubaccounts(address);
    console.log(response);
    const subaccounts = response.subaccounts;
    console.log(subaccounts);
    const subaccount = subaccounts[0];
    const subaccountNumber = subaccount.subaccountNumber;
    console.log(subaccountNumber);
  } catch (error) {
    console.log(error.message);
  }

  // Get subaccount 0
  try {
    const response = await client.account.getSubaccount(address, 0);
    console.log(response);
    const subaccount = response.subaccount;
    console.log(subaccount);
    const subaccountNumber = subaccount.subaccountNumber;
    console.log(subaccountNumber);
  } catch (error) {
    console.log(error.message);
  }

  // Get asset positions
  try {
    const response = await client.account.getSubaccountAssetPositions(address, 0);
    console.log(response);
    const positions = response.positions;
    console.log(positions);
    if (positions.length > 0) {
      const positions0 = positions[0];
      console.log(positions0);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get perp positions
  try {
    const response = await client.account.getSubaccountPerpetualPositions(address, 0);
    console.log(response);
    const positions = response.positions;
    console.log(positions);
    if (positions.length > 0) {
      const positions0 = positions[0];
      console.log(positions0);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get transfers
  try {
    const response = await client.account.getSubaccountTransfers(address, 0);
    console.log(response);
    const transfers = response.transfers;
    console.log(transfers);
    if (transfers.length > 0) {
      const transfer0 = transfers[0];
      console.log(transfer0);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get orders
  try {
    const response = await client.account.getSubaccountOrders(address, 0);
    console.log(response);
    const orders = response;
    console.log(orders);
    if (orders.length > 0) {
      const order0 = orders[0];
      console.log(order0);
      const order0Id = order0.id;
      console.log(order0Id);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get fills
  try {
    const response = await client.account.getSubaccountFills(address, 0);
    console.log(response);
    const fills = response.fills;
    console.log(fills);
    if (fills.length > 0) {
      const fill0 = fills[0];
      console.log(fill0);
    }
  } catch (error) {
    console.log(error.message);
  }

  // Get historical pnl
  try {
    const response = await client.account.getSubaccountHistoricalPNLs(address, 0);
    console.log(response);
    const historicalPnl = response.historicalPnl;
    console.log(historicalPnl);
    if (historicalPnl.length > 0) {
      const historicalPnl0 = historicalPnl[0];
      console.log(historicalPnl0);
    }
  } catch (error) {
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
