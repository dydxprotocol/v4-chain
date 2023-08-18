import { SyncHandlers } from '../../src/lib/synch-handlers';
import { IndexerTendermintBlock, IndexerTendermintEvent, MarketEventV1 } from '@dydxprotocol-indexer/v4-protos';
import { createIndexerTendermintBlock, createIndexerTendermintEvent } from '../helpers/indexer-proto-helpers';
import {
  defaultHeight,
  defaultMarketCreate,
  defaultMarketModify,
  defaultTime,
  defaultTxHash,
} from '../helpers/constants';
import { DydxIndexerSubtypes } from '../../src/lib/types';

const defaultMarketEventBinary: Uint8Array = Uint8Array.from(MarketEventV1.encode(
  defaultMarketCreate,
).finish());
const defaultMarketEventData: string = Buffer.from(
  defaultMarketEventBinary.buffer,
).toString('base64');

function generateFakeHandler(parallelizationIds: string[]): FakeHandler {
  const defaultTransactionIndex: number = 0;
  const defaultEventIndex: number = 0;
  const fakeTxId: number = 0;
  const defaultEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.MARKET,
    defaultMarketEventData,
    defaultTransactionIndex,
    defaultEventIndex,
  );

  const defaultEvent: IndexerTendermintEvent = createIndexerTendermintEvent(
    DydxIndexerSubtypes.ASSET,
    defaultAssetEventData,
    defaultTransactionIndex,
    defaultEventIndex,
  );

  const block: IndexerTendermintBlock = createIndexerTendermintBlock(
    defaultHeight,
    defaultTime,
    [defaultEvent],
    [defaultTxHash],
  );
  const handler: FakeHandler = new FakeHandler(
    block,
    defaultEvent,
    fakeTxId,
    {},
  );
  handler.setParallelizationIds(parallelizationIds);
  return handler;
}

describe('syncHandler', () => {
  describe('addHandler', () => {
    it('successfully adds handler', () => {
      const synchHandlers: SyncHandlers = new SyncHandlers();

      const handler: FakeHandler = generateFakeHandler(pIds);
      batchedHandlers.addHandler(handler);

      expect(batchedHandlers.batchedHandlers).toEqual([[handler]]);
      expect(batchedHandlers.pIdBatches).toEqual([new Set(pIds)]);
    });
  });
});
