CREATE OR REPLACE FUNCTION dydx_block_processor(block jsonb) RETURNS jsonb AS $$
/**
  Processes an entire block by creating the initial tendermint rows for the block and then processes each event
  individually through their respective handlers.

  Parameters:
    - block: A 'DecodedIndexerTendermintBlock' converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.

  Returns:
    An array containing the results for each event or NULL if this event is not handled by this block processor.
    See each individual handler function for a description of the the inputs and outputs.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)
*/
DECLARE
    block_height int = (block->'height')::int;
    block_time timestamp = (block->>'time')::timestamp;
    rval jsonb[];
    rval_to_merge jsonb[];
BEGIN
    PERFORM dydx_create_initial_rows_for_tendermint_block(block_height, block_time, block->'txHashes', block->'events');

    /** In genesis, handle ordered events first, then unordered events. In other blocks, handle unordered events first, then ordered events. */
    IF NOT block_height = 0 THEN
        rval = dydx_block_processor_unordered_handlers(block);
        rval_to_merge = dydx_block_processor_ordered_handlers(block);
    ELSE
        rval = dydx_block_processor_ordered_handlers(block);
        rval_to_merge = dydx_block_processor_unordered_handlers(block);
    END IF;

    /**
      Merge the results of the two handlers together by taking the first non-null result of each.

      Note that arrays are 1-indexed in PostgreSQL and empty arrays return NULL for array_length.
    */
    FOR i in 1..coalesce(array_length(rval, 1), 0) LOOP
        rval[i] = coalesce(rval[i], rval_to_merge[i]);
    END LOOP;

    RETURN to_jsonb(rval);
END;
$$ LANGUAGE plpgsql;
