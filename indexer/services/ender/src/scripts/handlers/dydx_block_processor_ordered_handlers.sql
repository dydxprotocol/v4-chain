CREATE OR REPLACE FUNCTION dydx_block_processor_ordered_handlers(block jsonb) RETURNS jsonb[] AS $$
/**
  Processes each event that should be handled by the batched handler. This includes all synchronous types
  (https://github.com/dydxprotocol/v4-chain/blob/b5d4e8a7c5cc48c460731b21c47f22eabef8b2b7/indexer/services/ender/src/lib/sync-handlers.ts#L11).

  Parameters:
    - block: A 'DecodedIndexerTendermintBlock' converted to JSON format. Conversion to JSON is expected to be done by JSON.stringify.

  Returns:
    An array containing the results for each event or NULL if this event is not handled by this block processor.
    See each individual handler function for a description of the the inputs and outputs.

  (Note that no text should exist before the function declaration to ensure that exception line numbers are correct.)

  TODO(IND-514): Remove the batch and sync handlers completely by moving all redis updates into
  a pipeline similar to how we return kafka events and then batch and emit them.
*/
DECLARE
    block_height int = (block->'height')::int;
    block_time timestamp = (block->>'time')::timestamp;
    event_ jsonb;
    rval jsonb[];
    event_index int;
    transaction_index int;
    event_data jsonb;
    -- Latency tracking variables
    event_start_time timestamp;
    event_end_time timestamp;
    event_latency interval;
BEGIN
    rval = array_fill(NULL::jsonb, ARRAY[coalesce(jsonb_array_length(block->'events'), 0)]::integer[]);

    /** Note that arrays are 1-indexed in PostgreSQL and empty arrays return NULL for array_length. */
    FOR i in 1..coalesce(array_length(rval, 1), 0) LOOP
        event_start_time := clock_timestamp();
        event_ = jsonb_array_element(block->'events', i-1);
        transaction_index = dydx_tendermint_event_to_transaction_index(event_);
        event_index = (event_->'eventIndex')::int;
        event_data = event_->'dataBytes';
        CASE event_->'subtype'
            WHEN '"market"'::jsonb THEN
                IF event_data->'priceUpdate' IS NOT NULL THEN
                    rval[i] = dydx_market_price_update_handler(block_height, block_time, event_data);
                ELSIF event_data->'marketCreate' IS NOT NULL THEN
                    rval[i] = dydx_market_create_handler(event_data);
                ELSIF event_data->'marketModify' IS NOT NULL THEN
                    rval[i] = dydx_market_modify_handler(event_data);
                ELSE
                    RAISE EXCEPTION 'Unknown market event %', event_;
                END IF;
            WHEN '"asset"'::jsonb THEN
                rval[i] = dydx_asset_create_handler(event_data);
            WHEN '"perpetual_market"'::jsonb THEN
                CASE (event_->'version')::int
                    WHEN 1 THEN
                        rval[i] = dydx_perpetual_market_v1_handler(event_data);
                    WHEN 2 THEN
                        rval[i] = dydx_perpetual_market_v2_handler(event_data);
                    WHEN 3 THEN
                        rval[i] = dydx_perpetual_market_v3_handler(event_data);
                    ELSE
                        NULL;
                END CASE;
            WHEN '"liquidity_tier"'::jsonb THEN
                rval[i] = dydx_liquidity_tier_handler(event_data);
            WHEN '"update_perpetual"'::jsonb THEN
                CASE (event_->'version')::int
                    WHEN 1 THEN
                        rval[i] = dydx_update_perpetual_v1_handler(event_data);
                    WHEN 2 THEN
                        rval[i] = dydx_update_perpetual_v2_handler(event_data);
                    WHEN 3 THEN
                        rval[i] = dydx_update_perpetual_v3_handler(event_data);
                    ELSE
                        NULL;
                END CASE;
            WHEN '"update_clob_pair"'::jsonb THEN
                rval[i] = dydx_update_clob_pair_handler(event_data);
            WHEN '"funding_values"'::jsonb THEN
                rval[i] = dydx_funding_handler(block_height, block_time, event_data, event_index, transaction_index);
            WHEN '"upsert_vault"'::jsonb THEN
                rval[i] = dydx_vault_upsert_handler(block_time, event_data);
            WHEN '"skipped_event"'::jsonb THEN
                rval[i] = jsonb_build_object();
            ELSE
                NULL;
            END CASE;

            event_end_time := clock_timestamp();
            event_latency := event_end_time - event_start_time;

            -- Add the event latency in ms to the rval output for this event
            rval[i] := jsonb_set(
                rval[i],
                '{latency}',
                to_jsonb(EXTRACT(EPOCH FROM event_latency) * 1000)
            );
    END LOOP;

    RETURN rval;
END;
$$ LANGUAGE plpgsql;
