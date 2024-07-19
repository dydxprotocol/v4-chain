import { IndexerTendermintEvent, IndexerTendermintEvent_BlockEvent } from '@dydxprotocol-indexer/v4-protos';
import { ParseMessageError } from '@dydxprotocol-indexer/base';
import { indexerTendermintEventToTransactionIndex } from '../../src/lib/helper';

describe('helper', () => {
  it.each([
    [
      'blockEvent BEGIN_BLOCK',
      {
        blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_BEGIN_BLOCK,
      },
      -2,
      false,
    ],
    [
      'blockEvent END_BLOCK',
      {
        blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_END_BLOCK,
      },
      -1,
      false,
    ],
    [
      'blockEvent UNSPECIFIED',
      {
        blockEvent: IndexerTendermintEvent_BlockEvent.BLOCK_EVENT_UNSPECIFIED,
      },
      0,
      true,
    ],
    [
      'transactionIndex 5',
      {
        oneofKind: 'transactionIndex',
        transactionIndex: 5,
      },
      5,
      false,
    ],
  ])('returns the correct transaction index for %s', (
    _name: string,
    eventFields: Partial<IndexerTendermintEvent>,
    expectedTxnIndex: number,
    throwError: boolean,
  ) => {
    const event: IndexerTendermintEvent = {
      ...eventFields,
      subtype: 'order_fill',
      dataBytes: Uint8Array.from(Buffer.from('data')),
      eventIndex: 0,
      version: 1,
    };
    if (throwError) {
      expect(() => indexerTendermintEventToTransactionIndex(event))
        .toThrowError(
          new ParseMessageError(
            `Received V4 event with invalid block event type: ${event.blockEvent}`,
          ),
        );
    } else {
      expect(indexerTendermintEventToTransactionIndex(event))
        .toEqual(expectedTxnIndex);
    }
  });
});
