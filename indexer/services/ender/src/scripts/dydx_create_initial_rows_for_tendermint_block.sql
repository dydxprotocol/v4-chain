CREATE OR REPLACE FUNCTION dydx_create_initial_rows_for_tendermint_block(
    block_height int, block_time timestamp, tx_hashes jsonb, events jsonb) RETURNS void AS $$
/**
  Parameters:
    - block_height: the height of the block being processed.
    - block_time: the time of the block being processed.
    - tx_hashes: Array of transaction hashes from the IndexerTendermintBlock.
    - events: Array of IndexerTendermintEvent objects.
  Returns: void.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
BEGIN
    -- Create block.
    INSERT INTO blocks ("blockHeight", "time") VALUES (block_height, block_time);

    -- Create transactions.
    IF tx_hashes IS NOT NULL AND jsonb_array_length(tx_hashes) > 0 THEN
        FOR i IN 0..jsonb_array_length(tx_hashes)-1 LOOP
            PERFORM dydx_create_transaction(jsonb_array_element_text(tx_hashes, i), block_height, i);
        END LOOP;
    END IF;

    -- Create tendermint events.
    IF events IS NOT NULL AND jsonb_array_length(events) > 0 THEN
        FOR i IN 0..jsonb_array_length(events)-1 LOOP
            PERFORM dydx_create_tendermint_event(jsonb_array_element(events, i), block_height);
        END LOOP;
    END IF;
END;
$$ LANGUAGE plpgsql;
