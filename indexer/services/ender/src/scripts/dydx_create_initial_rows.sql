/**
  Parameters:
    - block_height: the height of the block being processed.
    - block_time: the time of the block being processed.
    - tx_hashes: Array of transaction hashes from the IndexerTendermintBlock.
    - events: Array of IndexerTendermintEvent objects.
  Returns: void.
*/
CREATE OR REPLACE FUNCTION dydx_create_initial_rows(
    block_height text, block_time timestamp, tx_hashes text[], events jsonb[]) RETURNS void AS $$
BEGIN
    -- Create blocks entry
    INSERT INTO blocks ("blockHeight", "time") VALUES (block_height::bigint, block_time);

    -- Create Transactions
    FOR i IN 1..array_length(tx_hashes, 1) LOOP
        PERFORM dydx_create_transaction(tx_hashes[i], block_height, i);
    END LOOP;

    -- Create TendermintEvents
    FOR i IN 1..array_length(events, 1) LOOP
        PERFORM dydx_create_tendermint_event(events[i], block_height);
    END LOOP;
END;
$$ LANGUAGE plpgsql;
