import { axiosRequest, logger, stats } from '@dydxprotocol-indexer/base';
import { BlockFromDatabase, BlockTable, IsoString } from '@dydxprotocol-indexer/postgres';
import Big from 'big.js';
import { DateTime } from 'luxon';

import config from '../config';

const VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX = '/block';

type BlockData = {
  block: string,
  timestamp: IsoString,
};

type ValidatorHeader = {
  height: string,
  time: IsoString,
};

export default async function runTask(): Promise<void> {
  logger.info({
    at: 'track-lag#runTask',
    message: 'Running track lag task',
  });
  const [
    indexerBlockFromDatabase,
    indexerFullNodeBlock,
    validatorBlock,
    otherFullNodeBlock,
  ]: [
    BlockFromDatabase,
    BlockData,
    BlockData,
    BlockData,
  ] = await Promise.all([
    BlockTable.getLatest(),
    getValidatorBlockData(config.TRACK_LAG_INDEXER_FULL_NODE_URL, 'indexer_full_node'),
    getValidatorBlockData(config.TRACK_LAG_VALIDATOR_URL, 'validator'),
    getValidatorBlockData(config.TRACK_LAG_OTHER_FULL_NODE_URL, 'other_full_node'),
  ]);

  const indexerBlock: BlockData = {
    block: indexerBlockFromDatabase.blockHeight,
    timestamp: indexerBlockFromDatabase.time,
  };

  logAndStatLag(indexerFullNodeBlock, indexerBlock, 'indexer_full_node_to_indexer');
  logAndStatLag(validatorBlock, indexerFullNodeBlock, 'validator_to_indexer_full_node');
  logAndStatLag(validatorBlock, otherFullNodeBlock, 'validator_to_other_full_node');
  logAndStatLag(otherFullNodeBlock, indexerFullNodeBlock, 'other_full_node_to_indexer_full_node');
  logAndStatLag(validatorBlock, indexerBlock, 'validator_to_indexer');
}

async function getValidatorBlockData(urlPrefix: string, metricName: string): Promise<BlockData> {
  const start: number = Date.now();
  const url: string = `${urlPrefix}${VALIDATOR_BLOCK_HEIGHT_URL_SUFFIX}`;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const response: any = JSON.parse(await axiosRequest({
    url,
    method: 'GET',
    transformResponse: (res) => res,
  }) as string);
  stats.timing(
    `${config.SERVICE_NAME}.get_validator_block_data.timing`,
    Date.now() - start,
    { metricName },
  );
  const header: ValidatorHeader = response.result.block.header;
  stats.gauge(`${config.SERVICE_NAME}.block_height`, Number(header.height), { metricName });

  return {
    block: header.height,
    timestamp: header.time,
  };
}

function logAndStatLag(
  laterBlockData: BlockData,
  earlierBlockData: BlockData,
  lagType: string,
): void {
  const blockLag: string = Big(laterBlockData.block).minus(earlierBlockData.block).toString();
  const timeLagInMilliseconds: number = DateTime
    .fromISO(laterBlockData.timestamp)
    .diff(DateTime.fromISO(earlierBlockData.timestamp))
    .milliseconds;

  logger.info({
    at: 'track-lag#logAndStatLag',
    message: 'Got block lag',
    lagType,
    blockLag,
    timeLagInMilliseconds,
  });
  stats.timing(`${config.SERVICE_NAME}.block_lag`, Number(blockLag), { lagType });
  stats.timing(`${config.SERVICE_NAME}.time_lag`, timeLagInMilliseconds, { lagType });
}
