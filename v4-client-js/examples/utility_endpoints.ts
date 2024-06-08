/**
 * Simple JS example demostrating accessing public data with Indexer REST endpoints
 */

import { Network } from '../src/clients/constants';
import { IndexerClient } from '../src/clients/indexer-client';

async function test(): Promise<void> {
  const client = new IndexerClient(Network.testnet().indexerConfig);

  // Get indexer server time
  try {
    const response = await client.utility.getTime();
    console.log(response);
    const timeIso = response.iso;
    const timeEpoch = response.epoch;
    console.log('time');
    console.log(timeIso);
    console.log(timeEpoch);
  } catch (error) {
    console.log(error.message);
  }

  // Get indexer server height
  try {
    const response = await client.utility.getHeight();
    console.log(response);
    const height = response.height;
    const heightTime = response.time;
    console.log('height');
    console.log(height);
    console.log(heightTime);
  } catch (error) {
    console.log(error.message);
  }
}

test().then(() => {
}).catch((error) => {
  console.log(error.message);
});
