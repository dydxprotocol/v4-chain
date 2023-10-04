/**
  Gets the transaction index from the IndexerTendermint event.
  Parameters:
    - event: The IndexerTendermintEvent object.
  Returns: int.
*/
CREATE OR REPLACE FUNCTION dydx_tendermint_event_to_transaction_index(event jsonb) RETURNS int AS $$
DECLARE
    transaction_index_text text;
    block_event_text text;
BEGIN
    transaction_index_text := jsonb_extract_path_text(event, 'transactionIndex');
    block_event_text := jsonb_extract_path_text(event, 'blockEvent');

    IF transaction_index_text IS NOT NULL THEN
        RETURN transaction_index_text::int;
    ELSIF block_event_text IS NOT NULL THEN
        CASE block_event_text
            WHEN '1' THEN RETURN -2; -- BLOCK_EVENT_BEGIN_BLOCK
            WHEN '2' THEN RETURN -1; -- BLOCK_EVENT_END_BLOCK
            ELSE RAISE EXCEPTION 'Received V4 event with invalid block event type: %', block_event_text;
        END CASE;
    END IF;

    RAISE EXCEPTION 'Either transactionIndex or blockEvent must be defined in IndexerTendermintEvent';
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
