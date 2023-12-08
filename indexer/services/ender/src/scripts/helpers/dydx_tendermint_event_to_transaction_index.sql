CREATE OR REPLACE FUNCTION dydx_tendermint_event_to_transaction_index(event jsonb) RETURNS int AS $$
/**
  Gets the transaction index from the IndexerTendermint event.

  Parameters:
    - event: The JSON.stringify of a IndexerTendermintEvent object (https://github.com/dydxprotocol/v4-chain/blob/9ed26bd/proto/dydxprotocol/indexer/indexer_manager/event.proto#L25).
  Returns: int.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    IF event->'transactionIndex' IS NOT NULL THEN
        RETURN (event->'transactionIndex')::int;
    ELSIF event->'blockEvent' IS NOT NULL THEN
        CASE event->'blockEvent'
            WHEN '1'::jsonb THEN RETURN -2; /** BLOCK_EVENT_BEGIN_BLOCK */
            WHEN '2'::jsonb THEN RETURN -1; /** BLOCK_EVENT_END_BLOCK */
            ELSE RAISE EXCEPTION 'Received V4 event with invalid block event type: %', event->'blockEvent';
            END CASE;
    END IF;

    RAISE EXCEPTION 'Either transactionIndex or blockEvent must be defined in IndexerTendermintEvent';
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
