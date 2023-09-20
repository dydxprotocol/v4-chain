import yargs from 'yargs';

import { runAsyncScript } from './helpers/util';
import { BlockFromDatabase, BlockTable, IsoString } from '@dydxprotocol-indexer/postgres';
import { DateTime, Duration } from 'luxon';
import Big from 'big.js';
import { delay, safeAxiosRequest, wrapBackgroundTask } from '../../../packages/base/build';

const VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX = ':26657/block';

type BlockData = {
  block: string;
  timestamp: IsoString;
}

const args = yargs.options({
  full_node_url: {
    type: 'string',
    description: 'Indexer full node url such as http://52.194.155.74',
  },
  validator_url: {
    type: 'string',
    description: 'Validator url such as http://52.194.155.74',
  },
}).argv;

runAsyncScript(async () => {
  startLoop(
    trackLag,
    5_000, // 5 seconds
  );
});

export function startLoop(
  loopTask: () => Promise<unknown>,
  loopIntervalMs: number,
): void {
  console.log('Start of loop');
  wrapBackgroundTask(
    startLoopAsync(
      loopTask,
      loopIntervalMs,
    ),
    true,
    'taskName',
  );
}

async function startLoopAsync(
  loopTask: () => Promise<unknown>,
  loopIntervalMs: number,
): Promise<void> {
  for (;;) {
    // If lock was created, run the task.

    try {
      await loopTask();
    } catch (error) {
      console.log(`error: ${error}`);
    }
    await delay(loopIntervalMs);
  }
}

async function trackLag(): Promise<void> {
  console.log('Start of trackLag');
  const [
    indexerBlock,
    indexerFullNodeBlock,
    validatorBlock,
  ]: [
    BlockFromDatabase | undefined,
    BlockData,
    BlockData,
  ] = await Promise.all([
    BlockTable.getLatest(),
    getValidatorBlockData(args.full_node_url!),
    getValidatorBlockData(args.validator_url!),
  ]);

  if (indexerBlock === undefined) {
    return;
  }

  const indexerBlockLag: string = Big(indexerFullNodeBlock.block).minus(indexerBlock.blockHeight).toString();
  const indexerTimeLag: Duration = DateTime.fromISO(indexerFullNodeBlock.timestamp).diff(DateTime.fromISO(indexerBlock.time))
  const validatorBlockLag: string = Big(validatorBlock.block).minus(indexerBlock.blockHeight).toString();
  const validatorTimeLag: Duration = DateTime.fromISO(validatorBlock.timestamp).diff(DateTime.fromISO(indexerFullNodeBlock.timestamp))
  /*
  logger.info({
    at: 'track-lag#runTask',
    message: 'Got block lag',
    indexerBlockLag,
    indexerTimeLagInMilliseconds: indexerTimeLag.milliseconds,
    validatorBlockLag,
    validatorTimeLagInMilliseconds: validatorTimeLag.milliseconds,
  });
  */
}

async function getValidatorBlockData(url: string): Promise<BlockData> {
  const response = await safeAxiosRequest({
    url: `${url}${VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX}`,
    method: 'GET',
    transformResponse: (res) => res,
  });
  console.log(`url: ${url}\nresponse: ${response}`);

  return {
    block: '0',
    timestamp: '2021-01-01T00:00:00.000Z',
  };
}
