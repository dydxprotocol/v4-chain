/**
  Parameters:
    - event: The IndexerTendermintEvent object.
    - block_height: the height of the block being processed.
  Returns: The inserted event.
*/
CREATE OR REPLACE FUNCTION dydx_create_tendermint_event(
    event jsonb, block_height text
) RETURNS jsonb AS $$
DECLARE
    transaction_idx int;
    event_id bytea;
    inserted_event jsonb;
BEGIN
    transaction_idx := dydx_tendermint_event_to_transaction_index(event);
    event_id := dydx_event_id_from_parts(CAST(block_height AS int), transaction_idx, CAST(event->>'eventIndex' AS int));

    INSERT INTO tendermint_events ("id", "blockHeight", "transactionIndex", "eventIndex")
    VALUES (event_id, block_height::bigint, transaction_idx, CAST(event->>'eventIndex' AS int))
    RETURNING to_jsonb(tendermint_events.*) INTO inserted_event;

    RETURN inserted_event;
END;
$$ LANGUAGE plpgsql;
