import { logger } from '@dydxprotocol-indexer/base';
import { IsoString, PnlTicksCreateObject, Transaction } from '@dydxprotocol-indexer/postgres';

import { getPnlTicksCreateObjects } from './helpers/pnl-ticks-helper';

async function main() {

  const blockHeight: string = '12428756';
  const blockTime: IsoString = '2024-05-15T20:36:34.543Z';
  const txId: number = await Transaction.start();
  try {
    const newTicksToCreate: PnlTicksCreateObject[] = await
    getPnlTicksCreateObjects(blockHeight, blockTime, txId);
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
