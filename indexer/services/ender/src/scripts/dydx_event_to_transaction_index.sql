/**
  Converts a JSON.stringify format of an IndexerTendermintEvent to a transaction index (https://github.com/dydxprotocol/v4-proto/blob/8d35c86/dydxprotocol/indexer/indexer_manager/event.proto#L25).
*/
CREATE OR REPLACE FUNCTION dydx_event_to_transaction_index(event jsonb) RETURNS int AS $$
BEGIN
    IF event->'transactionIndex' IS NOT NULL THEN
        RETURN (event->'transactionIndex')::int;
    ELSIF event->'blockEvent' IS NOT NULL THEN
        CASE event->'blockEvent'
            WHEN '1'::jsonb /* BLOCK_EVENT_BEGIN_BLOCK */
                THEN RETURN -2;
            WHEN '2'::jsonb /* BLOCK_EVENT_END_BLOCK */
                THEN RETURN -1;
            ELSE
                RAISE EXCEPTION 'Received V4 event with invalid block event type: %', event->'blockEvent';
            END CASE;
    END IF;

    RAISE EXCEPTION 'TendermintEventTable.orderingWithinBlock.oneOfKind cannot be undefined';
END;
$$ LANGUAGE plpgsql IMMUTABLE PARALLEL SAFE;
