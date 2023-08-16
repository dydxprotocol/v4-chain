/**
  Returns an event id from parts.

  Parameters:
    - block_height: the height of the block being processing.
    - transaction_index: The transaction_index of the IndexerTendermintEvent after the conversion that takes into
        account the block_event (https://github.com/dydxprotocol/indexer/blob/cc70982/services/ender/src/lib/helper.ts#L33)
    - event_index: The 'event_index' of the IndexerTendermintEvent.
*/
CREATE OR REPLACE FUNCTION dydx_event_id_from_parts(block_height int, transaction_index int, event_index int) RETURNS bytea AS $$
BEGIN
    /*
    int4send converts to network order (which is also big endian order).
    || is the byte string concatenation operator.

    transactionIndex is -2 for BEGIN_BLOCK events, and -1 for END_BLOCK events. Increment by 2 to ensure result is >= 0.
    See https://github.com/dydxprotocol/indexer/blob/6aafb97/packages/postgres/src/stores/tendermint-event-table.ts#L33
    */
    RETURN int4send(block_height) || int4send(transaction_index + 2) || int4send(event_index);
END;
$$ language plpgsql IMMUTABLE PARALLEL SAFE;
