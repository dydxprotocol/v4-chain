import { delay, logger } from '@dydxprotocol-indexer/base';

export async function startUpdateLoop(
  updateFunction: () => Promise<void>,
  delayMs: number,
  name: string,
): Promise<void> {
  for (;;) {
    try {
      await updateFunction();
    } catch (error) {
      logger.error({
        at: name,
        message: 'Failed to run update',
        error,
      });
    }
    await delay(delayMs);
  }
}
