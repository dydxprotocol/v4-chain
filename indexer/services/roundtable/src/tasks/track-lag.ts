import { logger, safeAxiosRequest } from "@dydxprotocol-indexer/base";
import { BlockFromDatabase, BlockTable, IsoString } from "@dydxprotocol-indexer/postgres";
import Big from "big.js";
import { DateTime, Duration } from "luxon";

const INDEXER_FULL_NODE_URL = 'http://52.194.155.74';
const VALIDATOR_URL = 'http://18.188.95.153/';
const VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX = ':26657/block';

type BlockData = {
  block: string;
  timestamp: IsoString;
}

export default async function runTask(): Promise<void> {
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
    getValidatorBlockData(INDEXER_FULL_NODE_URL),
    getValidatorBlockData(VALIDATOR_URL),
  ]);

  if (indexerBlock === undefined) {
    return;
  }

  const indexerBlockLag: string = Big(indexerFullNodeBlock.block).minus(indexerBlock.blockHeight).toString();
  const indexerTimeLag: Duration = DateTime.fromISO(indexerFullNodeBlock.timestamp).diff(DateTime.fromISO(indexerBlock.time))
  const validatorBlockLag: string = Big(validatorBlock.block).minus(indexerBlock.blockHeight).toString();
  const validatorTimeLag: Duration = DateTime.fromISO(validatorBlock.timestamp).diff(DateTime.fromISO(indexerFullNodeBlock.timestamp))
  logger.info({
    at: 'track-lag#runTask',
    message: 'Got block lag',
    indexerBlockLag,
    indexerTimeLagInMilliseconds: indexerTimeLag.milliseconds,
    validatorBlockLag,
    validatorTimeLagInMilliseconds: validatorTimeLag.milliseconds,
  });
}

async function getValidatorBlockData(url: string): Promise<BlockData> {
  const response = safeAxiosRequest({
    url: `${url}${VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX}`,
    method: 'GET',
    transformResponse: (res) => res,
  });
  logger.info({
    at: 'track-lag#getValidatorBlockData',
    message: 'Got validator block data',
    url,
    response: JSON.stringify(response),
  });

  return {
    block: '0',
    timestamp: '2021-01-01T00:00:00.000Z',
  };
}
