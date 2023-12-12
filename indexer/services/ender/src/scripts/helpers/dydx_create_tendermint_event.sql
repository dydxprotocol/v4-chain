CREATE OR REPLACE FUNCTION dydx_create_tendermint_event(
    event jsonb, block_height int
) RETURNS jsonb AS $$
/**
  Parameters:
    - event: The IndexerTendermintEvent object.
    - block_height: the height of the block being processed.
  Returns: The inserted event.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    transaction_idx int;
    event_id bytea;
    inserted_event jsonb;
BEGIN
    transaction_idx := dydx_tendermint_event_to_transaction_index(event);
    event_id := dydx_event_id_from_parts(CAST(block_height AS int), transaction_idx, CAST(event->>'eventIndex' AS int));

    INSERT INTO tendermint_events ("id", "blockHeight", "transactionIndex", "eventIndex")
    VALUES (event_id, block_height, transaction_idx, (event->'eventIndex')::int)
    RETURNING to_jsonb(tendermint_events.*) INTO inserted_event;

    RETURN inserted_event;
END;
$$ LANGUAGE plpgsql;
