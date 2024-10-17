CREATE OR REPLACE FUNCTION dydx_block_processor_unordered_handlers(block jsonb) RETURNS jsonb[] AS $$
/**
  Processes each event that should be handled by the batched handler. This includes all supported non synchronous types
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
    USDC_ASSET_ID constant text = '0';

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
            WHEN '"order_fill"'::jsonb THEN
                /** If event_data.order is populated then this means it is not a liquidation order. */
                IF event_data->'order' IS NOT NULL THEN
                    rval[i] = jsonb_build_object(
                            'makerOrder',
                            dydx_order_fill_handler_per_order('makerOrder', block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index), 'MAKER', 'LIMIT', USDC_ASSET_ID, event_data->>'makerCanceledOrderStatus'),
                            'order',
                            dydx_order_fill_handler_per_order('order', block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index), 'TAKER', 'LIMIT', USDC_ASSET_ID, event_data->>'takerCanceledOrderStatus'));
                ELSE
                    rval[i] = jsonb_build_object(
                            'makerOrder',
                            dydx_liquidation_fill_handler_per_order('makerOrder', block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index), 'MAKER', 'LIQUIDATION', USDC_ASSET_ID),
                            'liquidationOrder',
                            dydx_liquidation_fill_handler_per_order('liquidationOrder', block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index), 'TAKER', 'LIQUIDATED', USDC_ASSET_ID));
                END IF;
            WHEN '"subaccount_update"'::jsonb THEN
                rval[i] = dydx_subaccount_update_handler(block_height, block_time, event_data, event_index, transaction_index);
            WHEN '"transfer"'::jsonb THEN
                rval[i] = dydx_transfer_handler(block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index));
            WHEN '"stateful_order"'::jsonb THEN
                rval[i] = dydx_stateful_order_handler(block_height, block_time, event_data);
            WHEN '"deleveraging"'::jsonb THEN
                rval[i] = dydx_deleveraging_handler(block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index));
            WHEN '"trading_reward"'::jsonb THEN
                rval[i] = dydx_trading_rewards_handler(block_height, block_time, event_data, event_index, transaction_index, jsonb_array_element_text(block->'txHashes', transaction_index));
            WHEN '"register_affiliate"'::jsonb THEN
                rval[i] = dydx_register_affiliate_handler(block_height, event_data);
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
