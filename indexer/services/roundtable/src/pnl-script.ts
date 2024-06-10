import { logger } from '@dydxprotocol-indexer/base';
import {
  BlockFromDatabase,
  BlockTable,
  IsoString,
  PnlTicksCreateObject, PnlTicksTable,
  Transaction
} from '@dydxprotocol-indexer/postgres';

import { getPnlTicksCreateObjects } from './helpers/pnl-ticks-helper';

async function main() {

  // const blockHeight: string = '12428756';
  // const blockTime: IsoString = '2024-05-15T20:36:34.543Z';
  const txId: number = await Transaction.start();
  const [
    block,
    ,
  ]: [
    BlockFromDatabase,
    string,
  ] = await Promise.all([
    BlockTable.getLatest({ readReplica: true }),
    PnlTicksTable.findLatestProcessedBlocktime(),
  ]);
  const latestBlockTime: string = block.time;
  const latestBlockHeight: string = block.blockHeight;
  logger.info({
    at: 'create-pnl-ticks#runTask',
    message: 'Latest block time',
    latestBlockTime,
    latestBlockHeight,
    txId,
  });
  try {
    const newTicksToCreate: PnlTicksCreateObject[] = await
    getPnlTicksCreateObjects(latestBlockHeight, latestBlockTime, txId, '306fe2fb-a398-595a-88fc-1cfa180e7a3e');
    logger.info({
      at: 'create-pnl-ticks#runTask',
      message: 'New PNL ticks created',
      newTicksToCreate,
      txId,
    });
  } catch (error) {
    logger.error({
      at: 'create-pnl-ticks#runTask',
      message: 'Error when getting pnl ticks',
      error,
      txId,
    });
    return;
  } finally {
    await Transaction.rollback(txId);
  }
}

main().then(() => {
  console.log('Process completed.');
}).catch((error) => {
  console.error('Failed to run main function:', error);
});
