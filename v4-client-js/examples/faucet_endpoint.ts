/**
 * Simple JS example demostrating filling subaccount with Faucet API
 */

import { FaucetApiHost } from '../src/clients/constants';
import { FaucetClient } from '../src/clients/faucet-client';
import { DYDX_TEST_ADDRESS } from './constants';

async function test(): Promise<void> {
  const client = new FaucetClient(FaucetApiHost.TESTNET);
  const address = DYDX_TEST_ADDRESS;

  // Use faucet to fill subaccount
  const faucetResponse = await client?.fill(address, 0, 2000);
  console.log(faucetResponse);
  const status = faucetResponse.status;
  console.log(status);
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
